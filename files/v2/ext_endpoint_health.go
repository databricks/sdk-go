package files

// Endpoint health: an endpoint that keeps failing transfers that a later
// endpoint then completes is left out of new transfers for a while, so that
// they stop paying a failed attempt on it, and is tried again after that.

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"
	"time"

	"github.com/databricks/sdk-go/core/apierr"
	"github.com/databricks/sdk-go/core/apiretry"
	"github.com/databricks/sdk-go/files/v2/internal/cloudstorage"
)

// endpointHealthKey keys a Client's endpoint health in Client.extensions.
// Uploads and downloads keep their own, as an endpoint may fail one and serve
// the other.
type endpointHealthKey struct{ upload bool }

// endpointHealth returns the client's endpoint health for uploads or
// downloads, built on first use with this engine's tunables.
func (e *engine) endpointHealth(upload bool) *endpointHealth {
	key := endpointHealthKey{upload: upload}
	if v, ok := e.c.extensions.Load(key); ok {
		return v.(*endpointHealth)
	}
	v, _ := e.c.extensions.LoadOrStore(key, newEndpointHealth(e.tun, e.c.logger))
	return v.(*endpointHealth)
}

// endpointHealth tracks which endpoints of one client keep failing.
//
// Only a failure that a later endpoint then made up for counts: if every
// endpoint fails, as for a missing file, the request is at fault, not the
// endpoint. After skipAfter such transfers in a row, the endpoint is skipped
// for a cool-off: transfers leave it out. Then one transfer probes it. A
// success clears the endpoint, and a failed probe doubles its cool-off, up to
// maxCoolOff. It is shared by every transfer of the client.
type endpointHealth struct {
	skipAfter                  int
	initialCoolOff, maxCoolOff time.Duration
	logger                     *slog.Logger

	mu          sync.Mutex
	states      map[endpoint]*endpointState
	nextProbeID int64
}

type endpointState struct {
	// failures is how many transfers in a row failed on the endpoint and then
	// completed on a later one.
	failures int
	coolOff  time.Duration
	// skipped is whether the endpoint is skipped until skipUntil.
	skipped   bool
	skipUntil time.Time
	// probe is the id of the probe one transfer was let through with after the
	// cool-off, or 0. A probe that hasn't reported within a cool-off is issued
	// again with a new id, so that a transfer that outlives it can't end the
	// newer one.
	probe int64
}

// endpointOrder is the endpoints one transfer walks, in order, and the probes
// it was given, by endpoint.
type endpointOrder struct {
	endpoints []endpoint
	probes    map[endpoint]int64
}

func newEndpointHealth(tun tunables, logger *slog.Logger) *endpointHealth {
	return &endpointHealth{
		skipAfter:      tun.endpointSkipAfter,
		initialCoolOff: tun.endpointSkipCoolOff,
		maxCoolOff:     tun.endpointSkipMaxCoolOff,
		logger:         logger,
		states:         map[endpoint]*endpointState{},
	}
}

// order returns the endpoints a new transfer walks: all but the skipped ones.
// The last one is never skipped, so that a transfer always has an endpoint and
// keeps the retries the last endpoint gets. The transfer must report the
// probes it was given, through recordFailure, recordSuccess or returnProbe.
func (h *endpointHealth) order(now time.Time, endpoints []endpoint) endpointOrder {
	order := endpointOrder{probes: map[endpoint]int64{}}
	h.mu.Lock()
	defer h.mu.Unlock()
	for i, ep := range endpoints {
		if s := h.states[ep]; i+1 < len(endpoints) && s != nil && s.skipped {
			if now.Before(s.skipUntil) {
				continue
			}
			// The cool-off is over: this transfer probes the endpoint, and
			// others keep skipping it until the probe is in. A probe that never
			// reports is issued again after another cool-off.
			h.nextProbeID++
			s.probe = h.nextProbeID
			s.skipUntil = now.Add(s.coolOff)
			order.probes[ep] = s.probe
		}
		order.endpoints = append(order.endpoints, ep)
	}
	return order
}

// returnProbe hands back a probe the transfer never tried the endpoint with,
// for example because an earlier endpoint took it, so that the next transfer
// probes it instead of waiting another cool-off.
func (h *endpointHealth) returnProbe(now time.Time, ep endpoint, probe int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if s := h.states[ep]; s != nil && s.probe == probe {
		s.probe = 0
		s.skipUntil = now
	}
}

func (h *endpointHealth) recordSuccess(ctx context.Context, ep endpoint) {
	h.mu.Lock()
	_, ok := h.states[ep]
	delete(h.states, ep)
	h.mu.Unlock()
	if ok {
		h.logger.DebugContext(ctx, "endpoint completed a transfer, so it is no longer skipped", "endpoint", ep.String())
	}
}

// recordFailure counts a transfer that failed on ep and then completed on a
// later endpoint. probe is the probe the transfer was given for ep, or 0.
func (h *endpointHealth) recordFailure(ctx context.Context, now time.Time, ep endpoint, probe int64) {
	if h.skipAfter <= 0 {
		return
	}
	h.mu.Lock()
	s := h.states[ep]
	if s == nil {
		s = &endpointState{}
		h.states[ep] = s
	}
	switch {
	case probe != 0 && s.probe == probe:
		s.probe = 0
		s.coolOff = min(2*s.coolOff, h.maxCoolOff)
		s.skipUntil = now.Add(s.coolOff)
		coolOff := s.coolOff
		h.mu.Unlock()
		h.logger.InfoContext(ctx, "endpoint still fails, skipping it", "endpoint", ep.String(), "for", coolOff)
		return
	case s.skipped:
		// Tried while skipped: by a transfer that started before the skip or
		// holds a probe since issued again, while another transfer held the
		// probe, or by an upload only a skipped endpoint can take.
		h.mu.Unlock()
		return
	}
	s.failures++
	if s.failures < h.skipAfter {
		h.mu.Unlock()
		return
	}
	s.skipped = true
	s.coolOff = h.initialCoolOff
	s.skipUntil = now.Add(s.coolOff)
	failures, coolOff := s.failures, s.coolOff
	h.mu.Unlock()
	h.logger.WarnContext(ctx, "endpoint failed transfers in a row that another endpoint then completed, skipping it",
		"endpoint", ep.String(), "failures", failures, "for", coolOff)
}

// transferDone feeds a transfer that succeeded back to the health that
// ordered it: walkedPast are the endpoints it moved on from on a failure that
// counts against them (see countsAgainstEndpoint), and tried every endpoint it
// tried, in order, ending with the one that took it.
func (e *engine) transferDone(ctx context.Context, h *endpointHealth, order endpointOrder, walkedPast, tried []endpoint) {
	now := e.clock()
	for _, ep := range walkedPast {
		h.recordFailure(ctx, now, ep, order.probes[ep])
	}
	h.recordSuccess(ctx, tried[len(tried)-1])
	for ep, probe := range order.probes {
		if !slices.Contains(tried, ep) {
			h.returnProbe(now, ep, probe)
		}
	}
}

// countsAgainstEndpoint reports whether a failure the walk moved on from says
// something about the endpoint, rather than a passing fault or the caller's
// own cancellation. A control call that got no answer counts, though it was
// retried: an endpoint that accepts connections and never answers would
// otherwise cost every transfer those retries.
func countsAgainstEndpoint(ctx context.Context, err error) bool {
	if ctx.Err() != nil {
		return false
	}
	return errors.Is(err, errNoAnswer) || !retryableFailure(ctx, err)
}

// retryableFailure reports whether err is a failure the endpoint's retries
// were for, which walks only because they or the retry budget ran out: a
// throttled or unavailable answer (408, 429, 502, 503 or 504, or a Retry-After
// hint), a dropped connection or a timeout, a stalled transfer, or the spent
// budget. The endpoint may be fine, so such a failure does not count against
// it. A timeout on a POST, such as initiate or a part-URL request, is
// retryable too, because the API client retries a transient network error on
// any method. An error's message alone never makes it retryable.
func retryableFailure(ctx context.Context, err error) bool {
	if aerr, ok := errors.AsType[*apierr.APIError](err); ok {
		if slices.Contains(controlPlaneRetryStatusCodes, aerr.HTTPStatusCode()) || apiretry.RetryDurationHint(err) > 0 {
			return true
		}
	}
	switch {
	case apiretry.IsTransientNetworkError(err), unanswered(err), cloudstorage.IsRetriable(err):
		return true
	case errors.Is(err, errUploadStalled), errors.Is(err, errDownloadStalled):
		return true
	}
	return errors.Is(err, context.DeadlineExceeded) && ctx.Err() == nil
}

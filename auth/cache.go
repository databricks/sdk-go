package auth

import (
	"context"
	"sync"
	"time"
)

const (
	// Default duration for the stale period. The number has been set
	// arbitrarily and might be changed in the future.
	defaultStaleDuration = 3 * time.Minute
)

type Option func(*cachedTokenProvider)

// WithCachedToken sets the initial token to be used by a cached token source.
func WithCachedToken(t *Token) Option {
	return func(cts *cachedTokenProvider) {
		cts.cachedToken = t
	}
}

// WithAsyncRefresh enables or disables the asynchronous token refresh.
func WithAsyncRefresh(b bool) Option {
	return func(cts *cachedTokenProvider) {
		cts.disableAsync = !b
	}
}

// NewCachedTokenProvider wraps a [TokenProvider] to cache the tokens it returns.
// By default, the cache will refresh tokens asynchronously a few minutes before
// they expire.
//
// The token cache is safe for concurrent use by multiple goroutines and will
// guarantee that only one token refresh is triggered at a time.
//
// The token cache does not take care of retries in case the token source
// returns and error; it is the responsibility of the provided token source to
// handle retries appropriately.
//
// If the TokenProvider is already a cached token source (obtained by calling this
// function), it is returned as is.
func NewCachedTokenProvider(ts TokenProvider, opts ...Option) TokenProvider {
	// This is meant as a niche optimization to avoid double caching of the
	// token source in situations where the caller needs caching guarantees
	// but does not know if the token source is already cached.
	if cts, ok := ts.(*cachedTokenProvider); ok {
		return cts
	}

	cts := &cachedTokenProvider{
		TokenProvider: ts,
		staleDuration: defaultStaleDuration,
		timeNow:       time.Now,
	}

	for _, opt := range opts {
		opt(cts)
	}

	return cts
}

type cachedTokenProvider struct {
	// The token source to obtain tokens from.
	TokenProvider TokenProvider

	// If true, only refresh the token with a blocking call when it is expired.
	disableAsync bool

	// Duration during which a token is considered stale, see tokenState.
	staleDuration time.Duration

	mu          sync.Mutex
	cachedToken *Token

	// Indicates that an async refresh is in progress. This is used to prevent
	// multiple async refreshes from being triggered at the same time.
	isRefreshing bool

	// Error returned by the last refresh. Async refreshes are disabled if this
	// value is not nil so that the cache does not continue sending request to
	// a potentially failing server. The next blocking call will re-enable async
	// refreshes by setting this value to nil if it succeeds, or return the
	// error if it fails.
	refreshErr error

	timeNow func() time.Time // for testing
}

// Token returns a token from the cache or fetches a new one if the current
// token is expired.
func (cts *cachedTokenProvider) Token(ctx context.Context) (*Token, error) {
	if cts.disableAsync {
		return cts.blockingToken(ctx)
	}
	return cts.asyncToken(ctx)
}

// tokenState represents the state of the token. Each token can be in one of
// the following three states:
//   - fresh: The token is valid.
//   - stale: The token is valid but will expire soon.
//   - expired: The token has expired and cannot be used.
//
// Token state through time:
//
//	issue time     expiry time
//	    v               v
//	    | fresh | stale | expired -> time
//	    |     valid     |
type tokenState int

const (
	fresh   tokenState = iota // The token is valid.
	stale                     // The token is valid but will expire soon.
	expired                   // The token has expired and cannot be used.
)

// tokenState returns the state of the token. The function is not thread-safe
// and should be called with the lock held.
func (c *cachedTokenProvider) tokenState() tokenState {
	if c.cachedToken == nil {
		return expired
	}
	switch lifeSpan := c.cachedToken.Expiry.Sub(c.timeNow()); {
	case lifeSpan <= 0:
		return expired
	case lifeSpan <= c.staleDuration:
		return stale
	default:
		return fresh
	}
}

func (cts *cachedTokenProvider) asyncToken(ctx context.Context) (*Token, error) {
	cts.mu.Lock()
	ts := cts.tokenState()
	t := cts.cachedToken
	cts.mu.Unlock()

	switch ts {
	case fresh:
		return t, nil
	case stale:
		cts.triggerAsyncRefresh(ctx)
		return t, nil
	default: // expired
		return cts.blockingToken(ctx)
	}
}

func (cts *cachedTokenProvider) blockingToken(ctx context.Context) (*Token, error) {
	cts.mu.Lock()

	// The lock is kept for the entire operation to ensure that only one
	// blockingToken operation is running at a time.
	defer cts.mu.Unlock()

	// This is important to recover from potential previous failed attempts
	// to refresh the token asynchronously, see declaration of refreshErr for
	// more information.
	cts.isRefreshing = false
	cts.refreshErr = nil

	// It's possible that the token got refreshed (either by a blockingToken or
	// an asyncRefresh call) while this particular call was waiting to acquire
	// the mutex. This check avoids refreshing the token again in such cases.
	if ts := cts.tokenState(); ts != expired { // fresh or stale
		return cts.cachedToken, nil
	}

	t, err := cts.TokenProvider.Token(ctx)
	if err != nil {
		return nil, err
	}
	cts.cachedToken = t
	return t, nil
}

func (cts *cachedTokenProvider) triggerAsyncRefresh(ctx context.Context) {
	cts.mu.Lock()
	defer cts.mu.Unlock()
	if !cts.isRefreshing && cts.refreshErr == nil {
		cts.isRefreshing = true

		go func() {
			t, err := cts.TokenProvider.Token(ctx)

			cts.mu.Lock()
			defer cts.mu.Unlock()
			cts.isRefreshing = false
			if err != nil {
				cts.refreshErr = err
				return
			}
			cts.cachedToken = t
		}()
	}
}

package v2

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/databricks/sdk-go/core/apierr"
)

type httpCallOptions struct {
	req    *http.Request
	client *http.Client
	logger *slog.Logger
}

// lazyRequest implements slog.LogValuer for lazy evaluation of request
// logging.
type lazyRequest struct{ req *http.Request }

func (r lazyRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("method", r.req.Method),
		slog.String("url", r.req.URL.String()),
	)
}

// lazyResponse implements slog.LogValuer for lazy evaluation of response
// logging.
type lazyResponse struct {
	resp *http.Response
	body []byte
}

func (r lazyResponse) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("status", r.resp.StatusCode),
		slog.String("body", string(r.body)),
	)
}

// executeHTTPCall executes an HTTP call and returns the response body or
// and error. This function takes care of logging the request and response.
func executeHTTPCall(opts httpCallOptions) ([]byte, error) {
	opts.logger.Debug("HTTP request", "request", lazyRequest{opts.req})
	resp, err := opts.client.Do(opts.req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		opts.logger.Debug("HTTP response", "response", lazyResponse{resp, nil})
		return nil, err
	}
	opts.logger.Debug("HTTP response", "response", lazyResponse{resp, body})
	if err := apierr.FromHTTPError(resp.StatusCode, resp.Header, body); err != nil {
		return nil, err
	}
	return body, nil
}

package internal

import (
	"net/http"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/log"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientOptions struct {
	HTTPClient  httpClient
	Credentials auth.Credentials
	Timeout     time.Duration
	Logger      log.Logger
	Debug       bool
}

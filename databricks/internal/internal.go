package internal

import (
	"net/http"
	"time"

	"github.com/databricks/sdk-go/auth"
	"github.com/databricks/sdk-go/core/log"
)

type ClientOptions struct {
	HTTPClient  *http.Client
	Credentials auth.Credentials
	Timeout     time.Duration
	Logger      log.Logger
	Debug       bool
}

package credentials

import (
	"context"
	"errors"
	"fmt"

	"github.com/databricks/sdk-go/auth"
)

var errTokenRequired = errors.New("token is required")

// NewPATCredentials returns a Credentials that can be used to authenticate with
// a Personal Access Token.
func NewPATCredentials(token string) (auth.Credentials, error) {
	if token == "" {
		return nil, errTokenRequired
	}
	return &patCredentials{token: token}, nil
}

type patCredentials struct {
	token string
}

func (c *patCredentials) AuthHeaders(ctx context.Context) ([]auth.Header, error) {
	return []auth.Header{
		{Key: "Authorization", Value: fmt.Sprintf("Bearer %s", c.token)},
	}, nil
}

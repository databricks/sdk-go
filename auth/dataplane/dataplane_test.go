package dataplane

import (
	"context"
	"fmt"
	"testing"

	"github.com/databricks/sdk-go/auth"
)

type mockClient func(context.Context, string, *auth.Token) (*auth.Token, error)

func (m mockClient) GetOAuthToken(ctx context.Context, authDetails string, t *auth.Token) (*auth.Token, error) {
	return m(ctx, authDetails, t)
}

func TestDataPlaneTokenProvider_Token(t *testing.T) {
	testErr := fmt.Errorf("test error")
	testToken := &auth.Token{Value: "access token", Type: "Bearer"}

	testCases := []struct {
		desc      string
		apiClient OAuthClient
		cpts      auth.TokenProvider
		wantToken *auth.Token
		wantErr   error
	}{
		{
			desc: "Failing control plane token source",
			cpts: auth.TokenProviderFn(func(context.Context) (*auth.Token, error) {
				return testToken, testErr
			}),
			wantErr: testErr,
		},
		{
			desc: "Failing oauth endpoint",
			cpts: auth.TokenProviderFn(func(context.Context) (*auth.Token, error) {
				return testToken, nil
			}),
			apiClient: mockClient(func(context.Context, string, *auth.Token) (*auth.Token, error) {
				return nil, testErr
			}),
			wantErr: testErr,
		},
		{
			desc: "Successful token retrieval",
			cpts: auth.TokenProviderFn(func(context.Context) (*auth.Token, error) {
				return &auth.Token{Value: "control plane test token", Type: "Bearer"}, nil
			}),
			apiClient: mockClient(func(context.Context, string, *auth.Token) (*auth.Token, error) {
				return testToken, nil
			}),
			wantToken: testToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			dpts := NewEndpointTokenProvider(tc.apiClient, tc.cpts)

			gotToken, gotErr := dpts.Token(context.Background(), "endpoint", "authDetails")

			if gotErr != tc.wantErr {
				t.Errorf("Token(): got error %v, want %v", gotErr, tc.wantErr)
			}
			if gotToken != tc.wantToken {
				t.Errorf("Token(): got token %v, want %v", gotToken, tc.wantToken)
			}
		})
	}
}

// Package oidc provides utilities for working with OIDC ID tokens.
//
// This package is experimental and subject to change.
package oidc

import (
	"context"
	"fmt"
	"os"
)

// IDToken represents an OIDC ID token that can be exchanged for a Databricks
// access token.
type IDToken struct {
	Value string
}

// IDTokenProvider is anything that returns an IDToken given an audience.
type IDTokenProvider interface {
	IDToken(ctx context.Context, audience string) (*IDToken, error)
}

// IDTokenProviderFn is an adapter to allow the use of ordinary functions as
// IDTokenProvider.
//
// Example:
//
//	   ts := IDTokenProviderFn(func(ctx context.Context, aud string) (*IDToken, error) {
//			return &IDToken{}, nil
//	   })
type IDTokenProviderFn func(ctx context.Context, audience string) (*IDToken, error)

func (fn IDTokenProviderFn) IDToken(ctx context.Context, audience string) (*IDToken, error) {
	return fn(ctx, audience)
}

// NewEnvIDTokenProvider returns an IDTokenProvider that reads the IDtoken from
// environment variable env.
//
// Note that the IDTokenProvider does not cache the token and will read the token
// from environment variable env each time.
func NewEnvIDTokenProvider(env string) IDTokenProvider {
	return IDTokenProviderFn(func(ctx context.Context, _ string) (*IDToken, error) {
		t := os.Getenv(env)
		if t == "" {
			return nil, fmt.Errorf("missing env var %q", env)
		}
		return &IDToken{Value: t}, nil
	})
}

// NewFileTokenProvider returns an IDTokenProvider that reads the ID token from a
// file. The file should contain a single line with the token.
func NewFileTokenProvider(path string) IDTokenProvider {
	return IDTokenProviderFn(func(ctx context.Context, _ string) (*IDToken, error) {
		if path == "" {
			return nil, fmt.Errorf("missing path")
		}
		t, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("file %q does not exist", path)
			}
			return nil, err
		}
		if len(t) == 0 {
			return nil, fmt.Errorf("file %q is empty", path)
		}
		return &IDToken{Value: string(t)}, nil
	})
}

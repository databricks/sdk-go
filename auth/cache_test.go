package auth

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// withTimeNow returns an Option that sets the time source used by the cached
// token source. This is intended for testing purposes only.
func withTimeNow(f func() time.Time) Option {
	return func(cts *cachedTokenProvider) {
		cts.timeNow = f
	}
}

func TestNewCachedTokenProvider_noCaching(t *testing.T) {
	want := &cachedTokenProvider{}
	got := NewCachedTokenProvider(want, nil)
	if got != want {
		t.Errorf("NewCachedTokenProvider() = %v, want %v", got, want)
	}
}

func TestNewCachedTokenProvider_default(t *testing.T) {
	ts := TokenProviderFn(func(_ context.Context) (*Token, error) {
		return nil, nil
	})

	got, ok := NewCachedTokenProvider(ts).(*cachedTokenProvider)
	if !ok {
		t.Fatalf("NewCachedTokenProvider() = %T, want *cachedTokenProvider", got)
	}

	if got.staleDuration != maxStaleDuration {
		t.Errorf("NewCachedTokenProvider() staleDuration = %v, want %v", got.staleDuration, maxStaleDuration)
	}
	if got.disableAsync != false {
		t.Errorf("NewCachedTokenProvider() disableAsync = %v, want %v", got.disableAsync, false)
	}
	if got.cachedToken != nil {
		t.Errorf("NewCachedTokenProvider() cachedToken = %v, want nil", got.cachedToken)
	}
}

func TestNewCachedTokenProvider_options(t *testing.T) {
	ts := TokenProviderFn(func(_ context.Context) (*Token, error) {
		return nil, nil
	})

	wantDisableAsync := false
	wantCachedToken := &Token{Expiry: time.Unix(42, 0)}

	opts := []Option{
		WithAsyncRefresh(!wantDisableAsync),
		WithCachedToken(wantCachedToken),
	}

	got, ok := NewCachedTokenProvider(ts, opts...).(*cachedTokenProvider)
	if !ok {
		t.Fatalf("NewCachedTokenProvider() = %T, want *cachedTokenProvider", got)
	}

	if got.disableAsync != wantDisableAsync {
		t.Errorf("NewCachedTokenProvider(): disableAsync = %v, want %v", got.disableAsync, wantDisableAsync)
	}
	if got.cachedToken != wantCachedToken {
		t.Errorf("NewCachedTokenProvider(): cachedToken = %v, want %v", got.cachedToken, wantCachedToken)
	}
}

func TestNewCachedTokenProvider_dynamicStaleDuration(t *testing.T) {
	now := time.Unix(1337, 0)

	testCases := []struct {
		name              string
		tokenTTL          time.Duration
		wantStaleDuration time.Duration
		advanceForStale   time.Duration
	}{
		{
			name:              "standard OAuth token with 60-minute TTL",
			tokenTTL:          60 * time.Minute,
			wantStaleDuration: 20 * time.Minute,
			advanceForStale:   41 * time.Minute,
		},
		{
			name:              "short-lived token with 10-minute TTL",
			tokenTTL:          10 * time.Minute,
			wantStaleDuration: 5 * time.Minute,
			advanceForStale:   6 * time.Minute,
		},
		{
			name:              "very short token with 90-second TTL",
			tokenTTL:          90 * time.Second,
			wantStaleDuration: 45 * time.Second,
			advanceForStale:   50 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fetchTime := now
			refreshed := false
			ts := TokenProviderFn(func(_ context.Context) (*Token, error) {
				if refreshed {
					t.Fatal("unexpected second token fetch")
				}
				refreshed = true
				return &Token{
					Expiry: fetchTime.Add(tc.tokenTTL),
				}, nil
			})

			initialToken := &Token{
				Expiry: now.Add(tc.tokenTTL),
			}
			cts, ok := NewCachedTokenProvider(ts,
				withTimeNow(func() time.Time { return now }),
				WithCachedToken(initialToken),
			).(*cachedTokenProvider)
			if !ok {
				t.Fatalf("NewCachedTokenSource() = %T, want *cachedTokenSource", cts)
			}

			if cts.staleDuration != tc.wantStaleDuration {
				t.Errorf("initial staleDuration = %v, want %v", cts.staleDuration, tc.wantStaleDuration)
			}
			if state := cts.tokenState(); state != fresh {
				t.Errorf("initial tokenState() = %v, want fresh", state)
			}

			cts.timeNow = func() time.Time { return now.Add(tc.advanceForStale) }
			if state := cts.tokenState(); state != stale {
				t.Errorf("tokenState() after %v = %v, want stale", tc.advanceForStale, state)
			}

			fetchTime = now.Add(tc.tokenTTL + 1*time.Second)
			cts.timeNow = func() time.Time { return fetchTime }

			_, err := cts.Token(context.Background())
			if err != nil {
				t.Fatalf("Token() error = %v", err)
			}

			if cts.staleDuration != tc.wantStaleDuration {
				t.Errorf("refreshed staleDuration = %v, want %v", cts.staleDuration, tc.wantStaleDuration)
			}
			if state := cts.tokenState(); state != fresh {
				t.Errorf("tokenState() after refresh = %v, want fresh", state)
			}

			cts.timeNow = func() time.Time { return fetchTime.Add(tc.advanceForStale) }
			if state := cts.tokenState(); state != stale {
				t.Errorf("tokenState() after refresh + %v = %v, want stale", tc.advanceForStale, state)
			}
		})
	}
}

func TestCachedTokenProvider_tokenState(t *testing.T) {
	now := time.Unix(1337, 0) // mock value for time.Now()

	testCases := []struct {
		token         *Token
		staleDuration time.Duration
		want          tokenState
	}{
		{
			token:         nil,
			staleDuration: 10 * time.Minute,
			want:          expired,
		},
		{
			token: &Token{
				Expiry: now.Add(-1 * time.Second),
			},
			staleDuration: 10 * time.Minute,
			want:          expired,
		},
		{
			token: &Token{
				Expiry: now.Add(1 * time.Hour),
			},
			staleDuration: 10 * time.Minute,
			want:          fresh,
		},
		{
			token: &Token{
				Expiry: now.Add(5 * time.Minute),
			},
			staleDuration: 10 * time.Minute,
			want:          stale,
		},
	}

	for _, tc := range testCases {
		cts := &cachedTokenProvider{
			cachedToken:   tc.token,
			staleDuration: tc.staleDuration,
			disableAsync:  false,
			timeNow:       func() time.Time { return now },
		}

		got := cts.tokenState()

		if got != tc.want {
			t.Errorf("tokenState() = %v, want %v", got, tc.want)
		}
	}
}

func TestCachedTokenProvider_Token(t *testing.T) {
	now := time.Unix(1337, 0) // mock value for time.Now()
	nTokenCalls := 10         // number of goroutines calling Token()
	testCases := []struct {
		desc         string // description of the test case
		cachedToken  *Token // token cached before calling Token()
		disableAsync bool   // whether are disabled or not
		refreshErr   error  // whether the cache was in error state

		returnedToken *Token // token returned by the token source
		returnedError error  // error returned by the token source

		wantCalls int    // expected number of calls to the token source
		wantToken *Token // expected token in the cache
	}{
		{
			desc:          "[Blocking] no cached token",
			disableAsync:  true,
			returnedToken: &Token{Expiry: now.Add(1 * time.Hour)},
			wantCalls:     1,
			wantToken:     &Token{Expiry: now.Add(1 * time.Hour)},
		},
		{
			desc:          "[Blocking] expired cached token",
			disableAsync:  true,
			cachedToken:   &Token{Expiry: now.Add(-1 * time.Second)},
			returnedToken: &Token{Expiry: now.Add(1 * time.Hour)},
			wantCalls:     1,
			wantToken:     &Token{Expiry: now.Add(1 * time.Hour)},
		},
		{
			desc:         "[Blocking] fresh cached token",
			disableAsync: true,
			cachedToken:  &Token{Expiry: now.Add(1 * time.Hour)},
			wantCalls:    0,
			wantToken:    &Token{Expiry: now.Add(1 * time.Hour)},
		},
		{
			desc:         "[Blocking] stale cached token",
			disableAsync: true,
			cachedToken:  &Token{Expiry: now.Add(1 * time.Minute)},
			wantCalls:    0,
			wantToken:    &Token{Expiry: now.Add(1 * time.Minute)},
		},
		{
			desc:          "[Blocking] refresh error",
			disableAsync:  true,
			returnedError: fmt.Errorf("test error"),
			wantCalls:     10,
		},
		{
			desc:          "[Blocking] recover from error",
			disableAsync:  true,
			refreshErr:    fmt.Errorf("refresh error"),
			cachedToken:   &Token{Expiry: now.Add(-1 * time.Minute)},
			returnedToken: &Token{Expiry: now.Add(-1 * time.Hour)},
			wantCalls:     10,
			wantToken:     &Token{Expiry: now.Add(-1 * time.Hour)},
		},
		{
			desc:          "[Async] no cached token",
			returnedToken: &Token{Expiry: now.Add(1 * time.Hour)},
			wantCalls:     1,
			wantToken:     &Token{Expiry: now.Add(1 * time.Hour)},
		},
		{
			desc:          "[Async] no cached token",
			returnedToken: &Token{Expiry: now.Add(1 * time.Hour)},
			wantCalls:     1,
			wantToken:     &Token{Expiry: now.Add(1 * time.Hour)},
		},
		{
			desc:          "[Async] expired cached token",
			cachedToken:   &Token{Expiry: now.Add(-1 * time.Second)},
			returnedToken: &Token{Expiry: now.Add(1 * time.Hour)},
			wantCalls:     1,
			wantToken:     &Token{Expiry: now.Add(1 * time.Hour)},
		},
		{
			desc:        "[Async] fresh cached token",
			cachedToken: &Token{Expiry: now.Add(1 * time.Hour)},
			wantCalls:   0,
			wantToken:   &Token{Expiry: now.Add(1 * time.Hour)},
		},
		{
			desc:          "[Async] stale cached token",
			cachedToken:   &Token{Expiry: now.Add(1 * time.Minute)},
			returnedToken: &Token{Expiry: now.Add(1 * time.Hour)},
			wantCalls:     1,
			wantToken:     &Token{Expiry: now.Add(1 * time.Hour)},
		},
		{
			desc:          "[Async] refresh error",
			cachedToken:   &Token{Expiry: now.Add(1 * time.Minute)},
			returnedError: fmt.Errorf("test error"),
			wantCalls:     1,
			wantToken:     &Token{Expiry: now.Add(1 * time.Minute)},
		},
		{
			desc:          "[Async] stale cached token, expired token returned",
			cachedToken:   &Token{Expiry: now.Add(1 * time.Minute)},
			returnedToken: &Token{Expiry: now.Add(-1 * time.Second)},
			wantCalls:     1,
			wantToken:     &Token{Expiry: now.Add(-1 * time.Second)},
		},
		{
			desc:          "[Async] recover from error",
			refreshErr:    fmt.Errorf("refresh error"),
			cachedToken:   &Token{Expiry: now.Add(-1 * time.Minute)},
			returnedToken: &Token{Expiry: now.Add(-1 * time.Hour)},
			wantCalls:     10,
			wantToken:     &Token{Expiry: now.Add(-1 * time.Hour)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var gotCalls atomic.Int32
			cts := &cachedTokenProvider{
				disableAsync:  tc.disableAsync,
				staleDuration: 10 * time.Minute,
				cachedToken:   tc.cachedToken,
				timeNow:       func() time.Time { return now },
				TokenProvider: TokenProviderFn(func(_ context.Context) (*Token, error) {
					gotCalls.Add(1)
					time.Sleep(10 * time.Millisecond)
					return tc.returnedToken, tc.returnedError
				}),
			}

			wg := sync.WaitGroup{}
			for i := 0; i < nTokenCalls; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, _ = cts.Token(context.Background())
				}()
			}

			wg.Wait()

			// Wait for async refreshes to finish. This part is a little brittle
			// but necessary to ensure that the async refresh is done before
			// checking the results. The use of the mutex is to avoid a potential
			// data race if the timeout is not long enough.
			time.Sleep(20 * time.Millisecond)
			cts.mu.Lock()
			gotToken := cts.cachedToken
			cts.mu.Unlock()

			if int(gotCalls.Load()) != tc.wantCalls {
				t.Errorf("want %d calls to cts.TokenProvider.Token(), got %d", tc.wantCalls, gotCalls.Load())
			}
			if !reflect.DeepEqual(tc.wantToken, gotToken) {
				t.Errorf("want cached token %v, got %v", tc.wantToken, gotToken)
			}
		})
	}
}

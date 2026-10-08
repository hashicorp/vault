// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package gcp

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/hashicorp/vault/api"
)

const (
	testRole         = "test-role"
	testVaultAddress = "https://vault.example.com"
	testJWT          = "header.payload.signature"

	identityPath = "/computeMetadata/v1/" + identityMetadataSuffix
)

// capturedRequest records what the fake metadata server was asked for so a test
// can assert on the request the metadata client actually built. The mutex is
// needed because the client retries some statuses, so the handler runs more
// than once, on goroutines the test does not join.
type capturedRequest struct {
	mu       sync.Mutex
	path     string
	query    url.Values
	flavor   string
	attempts int
}

func (c *capturedRequest) snapshot() (path string, query url.Values, flavor string, attempts int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.path, c.query, c.flavor, c.attempts
}

// startFakeMetadataService stands up an httptest server and points the metadata
// client at it via GCE_METADATA_HOST. Setting that variable also makes
// metadata.OnGCE report true without probing the network, which is what makes
// this code path testable off GCE.
//
// Note that OnGCE memoizes its result for the lifetime of the process, so no
// test here may assert the not-on-GCE branch: a false would leak into every
// later subtest and make this file order dependent.
func startFakeMetadataService(t *testing.T, handler http.HandlerFunc) *capturedRequest {
	t.Helper()

	captured := &capturedRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured.mu.Lock()
		captured.path = r.URL.Path
		captured.query = r.URL.Query()
		captured.flavor = r.Header.Get("Metadata-Flavor")
		captured.attempts++
		captured.mu.Unlock()
		handler(w, r)
	}))
	t.Cleanup(srv.Close)

	// The metadata client builds "http://" + host, so hand it a bare host:port.
	t.Setenv("GCE_METADATA_HOST", strings.TrimPrefix(srv.URL, "http://"))

	return captured
}

func testGCEAuth(t *testing.T) *GCPAuth {
	t.Helper()

	a, err := NewGCPAuth(testRole, WithGCEAuth())
	if err != nil {
		t.Fatalf("NewGCPAuth: unexpected error: %v", err)
	}

	return a
}

func TestGetJWTFromMetadataServiceSuccess(t *testing.T) {
	captured := startFakeMetadataService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testJWT))
	})

	jwt, err := testGCEAuth(t).getJWTFromMetadataService(context.Background(), testVaultAddress)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if jwt != testJWT {
		t.Errorf("jwt = %q, want %q", jwt, testJWT)
	}

	path, query, flavor, _ := captured.snapshot()
	if path != identityPath {
		t.Errorf("request path = %q, want %q", path, identityPath)
	}
	wantAudience := testVaultAddress + "/vault/" + testRole
	if got := query.Get("audience"); got != wantAudience {
		t.Errorf("audience = %q, want %q", got, wantAudience)
	}
	if got := query.Get("format"); got != "full" {
		t.Errorf("format = %q, want %q", got, "full")
	}
	if flavor != "Google" {
		t.Errorf("Metadata-Flavor = %q, want %q", flavor, "Google")
	}
}

// TestGetJWTFromMetadataServiceNon2xx covers the reported defect: a metadata
// endpoint that declines to issue an identity token used to have its error body
// returned as the JWT, so the caller saw a JWS parse error from Vault rather
// than the reason issuance failed.
func TestGetJWTFromMetadataServiceNon2xx(t *testing.T) {
	const responseBody = "Unable to generate token; IAM returned 403"

	startFakeMetadataService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(responseBody))
	})

	jwt, err := testGCEAuth(t).getJWTFromMetadataService(context.Background(), testVaultAddress)
	if err == nil {
		t.Fatalf("expected an error, got jwt %q", jwt)
	}
	if jwt != "" {
		t.Errorf("jwt = %q, want empty on error", jwt)
	}

	var metadataErr *metadata.Error
	if !errors.As(err, &metadataErr) {
		t.Fatalf("error is not a *metadata.Error: %v", err)
	}
	if metadataErr.Code != http.StatusForbidden {
		t.Errorf("status code = %d, want %d", metadataErr.Code, http.StatusForbidden)
	}
	if !strings.Contains(metadataErr.Message, responseBody) {
		t.Errorf("message = %q, want it to contain %q", metadataErr.Message, responseBody)
	}
}

func TestGetJWTFromMetadataServiceNotFound(t *testing.T) {
	startFakeMetadataService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	})

	jwt, err := testGCEAuth(t).getJWTFromMetadataService(context.Background(), testVaultAddress)
	if err == nil {
		t.Fatalf("expected an error, got jwt %q", jwt)
	}
	if jwt != "" {
		t.Errorf("jwt = %q, want empty on error", jwt)
	}

	var notDefined metadata.NotDefinedError
	if !errors.As(err, &notDefined) {
		t.Errorf("error is not a metadata.NotDefinedError: %v", err)
	}
}

// TestGetJWTFromMetadataServiceContextCancelled guards the switch to a context
// aware request. The metadata call previously ignored the context that Login
// receives, so it could not be cancelled.
func TestGetJWTFromMetadataServiceContextCancelled(t *testing.T) {
	startFakeMetadataService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testJWT))
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	jwt, err := testGCEAuth(t).getJWTFromMetadataService(ctx, testVaultAddress)
	if err == nil {
		t.Fatalf("expected an error, got jwt %q", jwt)
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("error = %v, want it to wrap context.Canceled", err)
	}
}

// TestLoginPropagatesContextToMetadata guards the call site in Login rather
// than the helper. Reverting Login to pass context.Background() leaves every
// other test in this file green, so without this the change users actually
// experience has no regression guard.
func TestLoginPropagatesContextToMetadata(t *testing.T) {
	startFakeMetadataService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(testJWT))
	})

	// An explicit config so VAULT_ADDR cannot redirect this. The address is
	// never dialled: if the context is honoured, Login fails before the write.
	client, err := api.NewClient(&api.Config{Address: "http://127.0.0.1:1"})
	if err != nil {
		t.Fatalf("api.NewClient: unexpected error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = testGCEAuth(t).Login(ctx, client)
	if err == nil {
		t.Fatal("expected an error")
	}

	// errors.Is(err, context.Canceled) holds either way, because the write to
	// Vault also respects the context. Only the wrapper distinguishes which
	// step failed, so assert on that.
	const want = "unable to retrieve JWT from GCE metadata service"
	if !strings.Contains(err.Error(), want) {
		t.Errorf("error = %q, want it to contain %q", err.Error(), want)
	}
}

// TestGetJWTFromMetadataServiceServerErrorRetries pins the largest behaviour
// change here: the metadata client retries 5xx and 429, which the hand-rolled
// request did not. The retry budget is bounded (5 attempts, 100ms doubling) and
// the backoff honours the context, which is what keeps the library's documented
// 15 second worst case from reaching a caller that sets a deadline.
func TestGetJWTFromMetadataServiceServerErrorRetries(t *testing.T) {
	captured := startFakeMetadataService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("metadata server is unwell"))
	})

	jwt, err := testGCEAuth(t).getJWTFromMetadataService(context.Background(), testVaultAddress)
	if err == nil {
		t.Fatalf("expected an error, got jwt %q", jwt)
	}
	if jwt != "" {
		t.Errorf("jwt = %q, want empty on error", jwt)
	}

	// Bounded, so a broken metadata service cannot retry forever.
	if _, _, _, attempts := captured.snapshot(); attempts != 6 {
		t.Errorf("attempts = %d, want 6 (initial + 5 retries)", attempts)
	}
}

// TestGetJWTFromMetadataServiceRetryHonoursDeadline is the other half: without
// it, the retry loop above would be an unbounded wait from the caller's point
// of view. The pre-fix code used a bare &http.Client{} with no timeout and no
// context, so a hung metadata service blocked Login indefinitely.
func TestGetJWTFromMetadataServiceRetryHonoursDeadline(t *testing.T) {
	startFakeMetadataService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, err := testGCEAuth(t).getJWTFromMetadataService(ctx, testVaultAddress)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected an error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("error = %v, want it to wrap context.DeadlineExceeded", err)
	}
	// Generous upper bound: the point is that it aborts on the deadline rather
	// than running the full retry budget, not the exact timing.
	if elapsed > 2*time.Second {
		t.Errorf("took %v, want the deadline to cut the retry loop short", elapsed)
	}
}

// TestGetJWTFromMetadataServiceEmptyBodyIsNotRejected records a gap rather than
// a guarantee. A 200 with an empty body is passed through as an empty token,
// which Vault then rejects at login, reaching the same confusing-error symptom
// by a second route. Issue #32068 reports the non-2xx case only, so rejecting
// it here would exceed that scope. This behaviour is not promised: a later
// change may reject empty tokens and should update this test rather than treat
// it as a contract.
func TestGetJWTFromMetadataServiceEmptyBodyIsNotRejected(t *testing.T) {
	startFakeMetadataService(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	jwt, err := testGCEAuth(t).getJWTFromMetadataService(context.Background(), testVaultAddress)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if jwt != "" {
		t.Errorf("jwt = %q, want empty", jwt)
	}
}

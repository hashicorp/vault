// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/command/client/auth"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/stretchr/testify/require"
)

// mockAuthMethodWithTracking is a mock auth method that tracks how many times
// Authenticate is called, which helps us verify the bug behavior
type mockAuthMethodWithTracking struct {
	authenticateCalls int
	authCalled        chan struct{}
	mu                sync.Mutex
}

func (m *mockAuthMethodWithTracking) Authenticate(_ context.Context, _ *api.Client) (string, http.Header, map[string]interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.authenticateCalls++
	m.authCalled <- struct{}{}

	// Return a valid auth response
	return "auth/approle/login", nil, map[string]interface{}{
		"role_id":   "test-role-id",
		"secret_id": "test-secret-id",
	}, nil
}

func (m *mockAuthMethodWithTracking) NewCreds() chan struct{} {
	return nil
}

func (m *mockAuthMethodWithTracking) CredSuccess() {}

func (m *mockAuthMethodWithTracking) Shutdown() {}

func (m *mockAuthMethodWithTracking) GetAuthenticateCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.authenticateCalls
}

// mockVaultServer is a lightweight mock HTTP server that simulates Vault API endpoints
// needed for testing AuthHandler without requiring a full Vault cluster
type mockVaultServer struct {
	mu                sync.Mutex
	statusCode        int
	errorMsg          string
	failCount         int
	lookupSelfCalls   int
	lookupSelfSuccess chan struct{}
	server            *httptest.Server
}

func newMockVaultServer(statusCode int, errorMsg string, failCount int) *mockVaultServer {
	m := &mockVaultServer{
		statusCode: statusCode,
		errorMsg:   errorMsg,
		failCount:  failCount,

		lookupSelfSuccess: make(chan struct{}),
	}

	m.server = httptest.NewServer(http.HandlerFunc(m.handler))
	return m
}

func (m *mockVaultServer) handler(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasSuffix(r.URL.Path, "/v1/auth/token/lookup-self"):
		m.handleLookupSelf(w)
	case strings.HasSuffix(r.URL.Path, "/v1/auth/token/create"):
		m.handleTokenCreate(w)
	case strings.HasSuffix(r.URL.Path, "/v1/auth/approle/login"):
		m.handleApproleLogin(w)
	default:
		http.Error(w, "endpoint not implemented in mock", http.StatusNotFound)
	}
}

func (m *mockVaultServer) handleLookupSelf(w http.ResponseWriter) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lookupSelfCalls++
	callNum := m.lookupSelfCalls
	shouldFail := callNum <= m.failCount

	if shouldFail {
		// Return configured error
		w.WriteHeader(m.statusCode)
		fmt.Fprintf(w, `{"errors":["%s"]}`, m.errorMsg)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"id":        "test-token-123",
			"ttl":       json.Number("3600"),
			"renewable": true,
			"policies":  []string{"default"},
			"type":      "service",
		},
	}
	json.NewEncoder(w).Encode(response)
	m.lookupSelfSuccess <- struct{}{}
}

func (m *mockVaultServer) handleTokenCreate(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token":   "test-token-123",
			"policies":       []string{"default"},
			"lease_duration": 3600,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (m *mockVaultServer) handleApproleLogin(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token":   "new-token-456",
			"policies":       []string{"default"},
			"lease_duration": 3600,
			"renewable":      true,
		},
	}
	json.NewEncoder(w).Encode(response)
}

func (m *mockVaultServer) URL() string {
	return m.server.URL
}

func (m *mockVaultServer) Close() {
	m.server.Close()
}

func (m *mockVaultServer) GetLookupSelfCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.lookupSelfCalls
}

func waitForPrecondition(precondition <-chan struct{}, timeout time.Duration) error {
	select {
	case <-precondition:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for precondition success")
	}
}

// TestAuthHandler_PreloadedTokenErrors tests various error scenarios during
// preloaded token lookup to ensure transient errors trigger retries while
// permanent errors trigger re-authentication.
//
// This test covers the bug where Vault Agent incorrectly treats transient errors
// (500/503) during initial token lookup-self as permanent failures, causing it to
// discard the cached token and re-authenticate instead of retrying the lookup.
//
// Expected behavior:
// - Transient errors (5xx, 429): Retry lookup-self with exponential backoff
// - Permanent errors (4xx): Discard token and re-authenticate
func TestAuthHandler_PreloadedTokenErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		statusCode     int
		errorMsg       string
		isTransient    bool
		failCount      int
		minLookupCalls int
		maxAuthCalls   int
		description    string
	}{
		{
			name:           "transient_500_retries",
			statusCode:     http.StatusInternalServerError,
			errorMsg:       "local node not active but active cluster node not found",
			isTransient:    true,
			failCount:      1,
			minLookupCalls: 2,
			maxAuthCalls:   0,
			description:    "500 error should trigger retry, not re-auth",
		},
		{
			name:           "transient_503_retries",
			statusCode:     http.StatusServiceUnavailable,
			errorMsg:       "service unavailable",
			isTransient:    true,
			failCount:      1,
			minLookupCalls: 2,
			maxAuthCalls:   0,
			description:    "503 error should trigger retry, not re-auth",
		},
		{
			name:           "transient_429_retries",
			statusCode:     http.StatusTooManyRequests,
			errorMsg:       "rate limit exceeded",
			isTransient:    true,
			failCount:      1,
			minLookupCalls: 2,
			maxAuthCalls:   0,
			description:    "429 error should trigger retry, not re-auth",
		},
		{
			name:           "permanent_403_reauths",
			statusCode:     http.StatusForbidden,
			errorMsg:       "permission denied",
			isTransient:    false,
			failCount:      1,
			minLookupCalls: 1,
			maxAuthCalls:   1,
			description:    "403 error should trigger re-auth, not retry",
		},
		{
			name:           "permanent_404_reauths",
			statusCode:     http.StatusNotFound,
			errorMsg:       "token not found",
			isTransient:    false,
			failCount:      1,
			minLookupCalls: 1,
			maxAuthCalls:   1,
			description:    "404 error should trigger re-auth, not retry",
		},
		{
			name:           "permanent_400_reauths",
			statusCode:     http.StatusBadRequest,
			errorMsg:       "bad request",
			isTransient:    false,
			failCount:      1,
			minLookupCalls: 1,
			maxAuthCalls:   1,
			description:    "400 error should trigger re-auth, not retry",
		},
		{
			name:           "multiple_transient_retries",
			statusCode:     http.StatusInternalServerError,
			errorMsg:       "internal server error",
			isTransient:    true,
			failCount:      2,
			minLookupCalls: 3,
			maxAuthCalls:   0,
			description:    "Multiple 500 errors should retry multiple times",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock server with configured error behavior
			mockServer := newMockVaultServer(tt.statusCode, tt.errorMsg, tt.failCount)
			defer mockServer.Close()

			// Create API client pointing to mock server
			config := api.DefaultConfig()
			config.Address = mockServer.URL()
			client, err := api.NewClient(config)
			require.NoError(t, err)

			// Get preloaded token
			preloadedToken := createMockToken(t)

			// Create mock auth method with channel for synchronization
			mockAuth := &mockAuthMethodWithTracking{
				authCalled: make(chan struct{}),
			}

			// Configure and start auth handler
			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			ah := auth.NewAuthHandler(&auth.AuthHandlerConfig{
				Logger:     logging.NewVaultLogger(hclog.Debug).Named("auth.handler"),
				Client:     client,
				Token:      preloadedToken,
				MinBackoff: 100 * time.Millisecond,
				MaxBackoff: 500 * time.Millisecond,
			})

			errCh := make(chan error, 1)
			go func() {
				errCh <- ah.Run(ctx, mockAuth)
			}()

			// We're only simulating errors here, so default to expecting authCalled, unless the error is transient, in which case we should do a lookup-self.
			precondition := mockAuth.authCalled
			if tt.isTransient {
				precondition = mockServer.lookupSelfSuccess
			}

			err = waitForPrecondition(precondition, 20*time.Second)
			require.NoError(t, err, "%s: precondition not met in time", tt.description)

			// Verify call counts
			lookupCalls := mockServer.GetLookupSelfCalls()
			authCalls := mockAuth.GetAuthenticateCalls()

			require.GreaterOrEqual(t, lookupCalls, tt.minLookupCalls,
				"%s: expected at least %d lookup-self calls, got %d",
				tt.description, tt.minLookupCalls, lookupCalls)

			if tt.isTransient {
				require.Equal(t, tt.maxAuthCalls, authCalls,
					"%s: expected %d authenticate calls (should retry, not re-auth), got %d",
					tt.description, tt.maxAuthCalls, authCalls)
			} else {
				require.GreaterOrEqual(t, authCalls, 1,
					"%s: expected at least 1 authenticate call (should re-auth), got %d",
					tt.description, authCalls)
			}

			cancelFunc()
			select {
			case <-errCh:
			case <-time.After(2 * time.Second):
				t.Fatal("timeout waiting for auth handler to stop")
			}
		})
	}
}

// createMockToken returns a test token for use with the mock server
func createMockToken(t *testing.T) string {
	t.Helper()
	// The mock server will validate this token when lookup-self is called
	return "test-token-123"
}

// headerTrackingAuthMethod is a mock auth method that returns a distinct
// header value on each Authenticate call, simulating Kerberos SPNEGO where
// a fresh ticket is generated per re-auth.
type headerTrackingAuthMethod struct {
	mu         sync.Mutex
	callCount  int
	headerKey  string
	authCalled chan struct{}
}

func (m *headerTrackingAuthMethod) Authenticate(_ context.Context, _ *api.Client) (string, http.Header, map[string]interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callCount++
	// Each call returns a unique header value, like a fresh Kerberos ticket.
	h := make(http.Header)
	h.Set(m.headerKey, fmt.Sprintf("token-cycle-%d", m.callCount))
	m.authCalled <- struct{}{}
	return "auth/approle/login", h, map[string]interface{}{
		"role_id":   "test-role-id",
		"secret_id": "test-secret-id",
	}, nil
}

func (m *headerTrackingAuthMethod) NewCreds() chan struct{} { return nil }
func (m *headerTrackingAuthMethod) CredSuccess()            {}
func (m *headerTrackingAuthMethod) Shutdown()               {}

func (m *headerTrackingAuthMethod) CallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

// mockVaultServerWithHeaderCapture extends mockVaultServer to capture the
// header values received on each login request.
type mockVaultServerWithHeaderCapture struct {
	mockVaultServer
	capturedHeaders []string
	headerKey       string
	loginCalled     chan struct{}
}

func newMockVaultServerWithHeaderCapture(headerKey string) *mockVaultServerWithHeaderCapture {
	m := &mockVaultServerWithHeaderCapture{
		mockVaultServer: mockVaultServer{
			lookupSelfSuccess: make(chan struct{}, 10),
		},
		headerKey:   headerKey,
		loginCalled: make(chan struct{}, 10),
	}
	m.server = httptest.NewServer(http.HandlerFunc(m.captureHandler))
	return m
}

func (m *mockVaultServerWithHeaderCapture) captureHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasSuffix(r.URL.Path, "/v1/auth/approle/login"):
		m.mu.Lock()
		// Capture all values for the tracked header key.
		m.capturedHeaders = append(m.capturedHeaders, r.Header[m.headerKey]...)
		m.mu.Unlock()
		// Return a very short, non-renewable lease so the LifetimeWatcher
		// fires DoneCh quickly and triggers the next re-auth cycle.
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token":   "new-token-456",
				"policies":       []string{"default"},
				"lease_duration": 1,
				"renewable":      false,
			},
		})
		m.loginCalled <- struct{}{}
	default:
		m.mockVaultServer.handler(w, r)
	}
}

func (m *mockVaultServerWithHeaderCapture) CapturedHeaders() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]string, len(m.capturedHeaders))
	copy(result, m.capturedHeaders)
	return result
}

// TestAuthHandler_HeaderNotAccumulatedAcrossReauthCycles is a regression test
// for the Kerberos auto-auth bug where AddHeader was called on every re-auth
// cycle, causing duplicate Authorization: Negotiate headers to accumulate on
// the shared client. Vault reads only the first value via http.Header.Get, so
// stale SPNEGO tickets were silently sent after the first renewal, breaking
// Kerberos authentication.
//
// This test verifies that after multiple re-auth cycles:
//  1. The server receives exactly one value per header key per request.
//  2. The value received is the fresh one from the current Authenticate call,
//     not a stale one from a previous cycle.
func TestAuthHandler_HeaderNotAccumulatedAcrossReauthCycles(t *testing.T) {
	t.Parallel()

	const headerKey = "Authorization"
	const cycles = 3

	mockServer := newMockVaultServerWithHeaderCapture(headerKey)
	defer mockServer.Close()

	config := api.DefaultConfig()
	config.Address = mockServer.URL()
	client, err := api.NewClient(config)
	require.NoError(t, err)

	mockAuth := &headerTrackingAuthMethod{
		headerKey:  headerKey,
		authCalled: make(chan struct{}, cycles+1),
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	ah := auth.NewAuthHandler(&auth.AuthHandlerConfig{
		Logger:     logging.NewVaultLogger(hclog.Debug).Named("auth.handler"),
		Client:     client,
		MinBackoff: 10 * time.Millisecond,
		MaxBackoff: 50 * time.Millisecond,
	})

	errCh := make(chan error, 1)
	go func() {
		errCh <- ah.Run(ctx, mockAuth)
	}()

	// Drain OutputCh so the auth handler is never blocked waiting for a
	// consumer. Without this the channel fills after cycle 1 and the handler
	// deadlocks on the second OutputCh send.
	go func() {
		for range ah.OutputCh {
		}
	}()

	// Wait for the desired number of re-auth cycles.
	timeout := time.After(15 * time.Second)
	for i := 0; i < cycles; i++ {
		select {
		case <-mockServer.loginCalled:
		case <-timeout:
			t.Fatalf("timed out waiting for re-auth cycle %d", i+1)
		}
	}

	cancelFunc()
	select {
	case <-errCh:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for auth handler to stop")
	}

	captured := mockServer.CapturedHeaders()
	require.GreaterOrEqual(t, len(captured), cycles,
		"expected at least %d captured header values, got %d", cycles, len(captured))

	// Each login request must have received exactly one header value (not
	// an ever-growing list of accumulated stale values), and each value
	// must match the fresh token generated for that specific cycle.
	for i, val := range captured {
		expected := fmt.Sprintf("token-cycle-%d", i+1)
		require.Equalf(t, expected, val,
			"cycle %d: expected header value %q, got %q - stale header may have been sent", i+1, expected, val)
	}
}

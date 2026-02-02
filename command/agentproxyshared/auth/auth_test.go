// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package auth

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
	"github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

type userpassTestMethod struct{}

func newUserpassTestMethod(t *testing.T, client *api.Client) AuthMethod {
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
		Config: api.AuthConfigInput{
			DefaultLeaseTTL: "1s",
			MaxLeaseTTL:     "3s",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	return &userpassTestMethod{}
}

func (u *userpassTestMethod) Authenticate(_ context.Context, client *api.Client) (string, http.Header, map[string]interface{}, error) {
	_, err := client.Logical().Write("auth/userpass/users/foo", map[string]interface{}{
		"password": "bar",
	})
	if err != nil {
		return "", nil, nil, err
	}
	return "auth/userpass/login/foo", nil, map[string]interface{}{
		"password": "bar",
	}, nil
}

func (u *userpassTestMethod) NewCreds() chan struct{} {
	return nil
}

func (u *userpassTestMethod) CredSuccess() {
}

func (u *userpassTestMethod) Shutdown() {
}

func TestAuthHandler(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	ctx, cancelFunc := context.WithCancel(context.Background())

	ah := NewAuthHandler(&AuthHandlerConfig{
		Logger: logging.NewVaultLogger(hclog.Trace).Named("auth.handler"),
		Client: client,
	})

	am := newUserpassTestMethod(t, client)
	errCh := make(chan error)
	go func() {
		errCh <- ah.Run(ctx, am)
	}()

	// Consume tokens so we don't block
	stopTime := time.Now().Add(5 * time.Second)
	closed := false
consumption:
	for {
		select {
		case err := <-errCh:
			if err != nil {
				t.Fatal(err)
			}
			break consumption
		case <-ah.OutputCh:
		case <-ah.TemplateTokenCh:
		// Nothing
		case <-time.After(stopTime.Sub(time.Now())):
			if !closed {
				cancelFunc()
				closed = true
			}
		}
	}
}

func TestAgentBackoff(t *testing.T) {
	max := 1024 * time.Second
	backoff := newAutoAuthBackoff(consts.DefaultMinBackoff, max, false)

	// Test initial value
	if backoff.backoff.Current() > consts.DefaultMinBackoff || backoff.backoff.Current() < consts.DefaultMinBackoff*3/4 {
		t.Fatalf("expected 1s initial backoff, got: %v", backoff.backoff.Current())
	}

	// Test that backoffSleep values are in expected range (75-100% of 2*previous)
	next, _ := backoff.backoff.Next()
	for i := 0; i < 9; i++ {
		old := next
		next, _ = backoff.backoff.Next()

		expMax := 2 * old
		expMin := 3 * expMax / 4

		if next < expMin || next > expMax {
			t.Fatalf("expected backoffSleep in range %v to %v, got: %v", expMin, expMax, backoff)
		}
	}

	// Test that backoffSleep is capped
	for i := 0; i < 100; i++ {
		_, _ = backoff.backoff.Next()
		if backoff.backoff.Current() > max {
			t.Fatalf("backoff exceeded max of 100s: %v", backoff)
		}
	}

	// Test reset
	backoff.backoff.Reset()
	if backoff.backoff.Current() > consts.DefaultMinBackoff || backoff.backoff.Current() < consts.DefaultMinBackoff*3/4 {
		t.Fatalf("expected 1s backoff after reset, got: %v", backoff.backoff.Current())
	}
}

func TestAgentMinBackoffCustom(t *testing.T) {
	type test struct {
		minBackoff time.Duration
		want       time.Duration
	}

	tests := []test{
		{minBackoff: 0 * time.Second, want: 1 * time.Second},
		{minBackoff: 1 * time.Second, want: 1 * time.Second},
		{minBackoff: 5 * time.Second, want: 5 * time.Second},
		{minBackoff: 10 * time.Second, want: 10 * time.Second},
	}

	for _, test := range tests {
		max := 1024 * time.Second
		backoff := newAutoAuthBackoff(test.minBackoff, max, false)

		// Test initial value
		if backoff.backoff.Current() > test.want || backoff.backoff.Current() < test.want*3/4 {
			t.Fatalf("expected %d initial backoffSleep, got: %v", test.want, backoff.backoff.Current())
		}

		// Test that backoffSleep values are in expected range (75-100% of 2*previous)
		next, _ := backoff.backoff.Next()
		for i := 0; i < 5; i++ {
			old := next
			next, _ = backoff.backoff.Next()

			expMax := 2 * old
			expMin := 3 * expMax / 4

			if next < expMin || next > expMax {
				t.Fatalf("expected backoffSleep in range %v to %v, got: %v", expMin, expMax, backoff)
			}
		}

		// Test that backoffSleep is capped
		for i := 0; i < 100; i++ {
			next, _ = backoff.backoff.Next()
			if next > max {
				t.Fatalf("backoffSleep exceeded max of 100s: %v", backoff)
			}
		}

		// Test reset
		backoff.backoff.Reset()
		if backoff.backoff.Current() > test.want || backoff.backoff.Current() < test.want*3/4 {
			t.Fatalf("expected %d backoffSleep after reset, got: %v", test.want, backoff.backoff.Current())
		}
	}
}

// mockAuthMethodWithTracking is a mock auth method that tracks how many times
// Authenticate is called, which helps us verify the bug behavior
type mockAuthMethodWithTracking struct {
	authenticateCalls int
	authCalled        chan struct{}
	mu                sync.Mutex
}

func (m *mockAuthMethodWithTracking) Authenticate(ctx context.Context, client *api.Client) (string, http.Header, map[string]interface{}, error) {
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
		m.handleLookupSelf(w, r)
	case strings.HasSuffix(r.URL.Path, "/v1/auth/token/create"):
		m.handleTokenCreate(w, r)
	case strings.HasSuffix(r.URL.Path, "/v1/auth/approle/login"):
		m.handleApproleLogin(w, r)
	default:
		http.Error(w, "endpoint not implemented in mock", http.StatusNotFound)
	}
}

func (m *mockVaultServer) handleLookupSelf(w http.ResponseWriter, r *http.Request) {
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

func (m *mockVaultServer) handleTokenCreate(w http.ResponseWriter, r *http.Request) {
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

func (m *mockVaultServer) handleApproleLogin(w http.ResponseWriter, r *http.Request) {
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

func waitForPrecondition(precondition *chan struct{}, timeout time.Duration) error {
	select {
	case <-*precondition:
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

			ah := NewAuthHandler(&AuthHandlerConfig{
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

			err = waitForPrecondition(&precondition, 20*time.Second)
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

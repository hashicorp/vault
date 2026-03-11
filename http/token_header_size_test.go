// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// sendTokenHeaderRequest builds and dispatches an HTTP GET to addr, placing
// token in the named header. For "Authorization", the value is wrapped as a
// Bearer token per RFC 6750.
func sendTokenHeaderRequest(t *testing.T, client *http.Client, addr, headerName, token string) (*http.Response, error) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, addr+"/v1/auth/token/lookup-self", nil)
	require.NoError(t, err)
	if headerName == "Authorization" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		req.Header.Set(headerName, token)
	}
	return client.Do(req)
}

// newTestClusterForTokenHeader creates a NewTestCluster with default settings
// and returns the cluster, an HTTP client configured with TLS, and the address.
func newTestClusterForTokenHeader(t *testing.T, opts *vault.TestClusterOptions) (*http.Client, string) {
	t.Helper()
	if opts == nil {
		opts = &vault.TestClusterOptions{}
	}
	opts.HandlerFunc = Handler
	cluster := vault.NewTestCluster(t, nil, opts)
	cluster.Start()

	core := cluster.Cores[0]
	transport := cleanhttp.DefaultPooledTransport()
	transport.TLSClientConfig = core.TLSConfig()
	httpClient := &http.Client{Transport: transport, Timeout: 15 * time.Second}
	return httpClient, core.Client.Address()
}

// TestTokenHeader_ExceedsDefaultLimit_IsRejected verifies that tokens larger than
// DefaultMaxTokenHeaderSize are rejected with 400 before token validation occurs.
func TestTokenHeader_ExceedsDefaultLimit_IsRejected(t *testing.T) {
	t.Parallel()

	client, addr := newTestClusterForTokenHeader(t, nil)

	cases := []struct {
		name       string
		headerSize int
		headerName string
	}{
		{
			name:       "just-over-limit-x-vault-token",
			headerSize: DefaultMaxTokenHeaderSize + 1,
			headerName: consts.AuthHeaderName,
		},
		{
			name:       "just-over-limit-authorization-bearer",
			headerSize: DefaultMaxTokenHeaderSize + 1,
			headerName: "Authorization",
		},
		{
			name:       "100kb-x-vault-token",
			headerSize: 100 * 1024,
			headerName: consts.AuthHeaderName,
		},
		{
			name:       "100kb-authorization-bearer",
			headerSize: 100 * 1024,
			headerName: "Authorization",
		},
		{
			name:       "900kb-x-vault-token",
			headerSize: 900 * 1024,
			headerName: consts.AuthHeaderName,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			token := strings.Repeat("x", tc.headerSize)
			resp, err := sendTokenHeaderRequest(t, client, addr, tc.headerName, token)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusBadRequest, resp.StatusCode,
				"token of %d bytes must be rejected by Vault middleware (want 400, not 403)",
				tc.headerSize)
		})
	}
}

// TestTokenHeader_WithinDefaultLimit_IsProcessed verifies that tokens at or
// below DefaultMaxTokenHeaderSize reach token validation normally (403, not 400).
func TestTokenHeader_WithinDefaultLimit_IsProcessed(t *testing.T) {
	t.Parallel()

	client, addr := newTestClusterForTokenHeader(t, nil)

	cases := []struct {
		name       string
		headerSize int
		headerName string
	}{
		{
			name:       "512b-x-vault-token",
			headerSize: 512,
			headerName: consts.AuthHeaderName,
		},
		{
			name:       "512b-authorization-bearer",
			headerSize: 512,
			headerName: "Authorization",
		},
		{
			name:       "4kb-x-vault-token",
			headerSize: 4 * 1024,
			headerName: consts.AuthHeaderName,
		},
		{
			name:       "at-limit-x-vault-token",
			headerSize: DefaultMaxTokenHeaderSize,
			headerName: consts.AuthHeaderName,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			token := strings.Repeat("x", tc.headerSize)
			resp, err := sendTokenHeaderRequest(t, client, addr, tc.headerName, token)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusForbidden, resp.StatusCode,
				"token of %d bytes is within the limit and must reach validation", tc.headerSize)
		})
	}
}

// TestTokenHeader_ConfigurableLimit_Enforced verifies that an operator can
// lower the token header size limit via the listener configuration, and that
// Vault enforces the configured value rather than DefaultMaxTokenHeaderSize.
func TestTokenHeader_ConfigurableLimit_Enforced(t *testing.T) {
	t.Parallel()

	const customLimit = 1024 // 1 KB — well below the default 8 KB
	client, addr := newTestClusterForTokenHeader(t, &vault.TestClusterOptions{
		DefaultHandlerProperties: vault.HandlerProperties{
			ListenerConfig: &configutil.Listener{
				CustomMaxTokenHeaderSize: customLimit,
			},
		},
	})

	oversized := strings.Repeat("x", customLimit+1)
	resp, err := sendTokenHeaderRequest(t, client, addr, consts.AuthHeaderName, oversized)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode,
		"token exceeding custom limit of %d bytes must be rejected with 400", customLimit)

	atLimit := strings.Repeat("x", customLimit)
	resp2, err := sendTokenHeaderRequest(t, client, addr, consts.AuthHeaderName, atLimit)
	require.NoError(t, err)
	defer resp2.Body.Close()

	require.Equal(t, http.StatusForbidden, resp2.StatusCode,
		"token at the custom limit (%d bytes) must reach validation", customLimit)
}

// TestTokenHeader_MultipleAuthorizationHeaders_BypassPrevented verifies that an
// oversized Bearer token is rejected even when a non-Bearer Authorization header
// precedes it in the request.
func TestTokenHeader_MultipleAuthorizationHeaders_BypassPrevented(t *testing.T) {
	t.Parallel()

	client, addr := newTestClusterForTokenHeader(t, nil)

	oversizedBearer := strings.Repeat("x", DefaultMaxTokenHeaderSize+1)

	req, err := http.NewRequest(http.MethodGet, addr+"/v1/auth/token/lookup-self", nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Basic dXNlcjpwYXNz")
	req.Header.Add("Authorization", "Bearer "+oversizedBearer)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode,
		"oversized Bearer token must be rejected even when a non-Bearer Authorization header precedes it")
}

// TestTokenHeader_DisabledLimit_AllowsOversizedToken verifies that setting
// max_token_header_size = -1 in the listener config fully disables the check,
// allowing tokens larger than DefaultMaxTokenHeaderSize to reach validation.
func TestTokenHeader_DisabledLimit_AllowsOversizedToken(t *testing.T) {
	t.Parallel()

	client, addr := newTestClusterForTokenHeader(t, &vault.TestClusterOptions{
		DefaultHandlerProperties: vault.HandlerProperties{
			ListenerConfig: &configutil.Listener{
				CustomMaxTokenHeaderSize: -1,
			},
		},
	})

	oversized := strings.Repeat("x", DefaultMaxTokenHeaderSize*2)
	resp, err := sendTokenHeaderRequest(t, client, addr, consts.AuthHeaderName, oversized)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusForbidden, resp.StatusCode,
		"with max_token_header_size = -1 the size guard must not fire; token must reach validation")
}

// TestTokenHeader_ErrorResponseFormat verifies that oversized-token rejections
// return a valid Vault JSON error envelope with an "errors" array.
func TestTokenHeader_ErrorResponseFormat(t *testing.T) {
	t.Parallel()

	client, addr := newTestClusterForTokenHeader(t, nil)

	token := strings.Repeat("x", DefaultMaxTokenHeaderSize+1)
	resp, err := sendTokenHeaderRequest(t, client, addr, consts.AuthHeaderName, token)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	require.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var envelope struct {
		Errors []string `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(body, &envelope),
		"response body must be valid JSON: %s", string(body))
	require.NotEmpty(t, envelope.Errors,
		"response must contain at least one error message")
	require.Contains(t, envelope.Errors[0], "authentication token",
		"error message must describe the token size violation")
}

// TestTokenHeader_NoAuthHeader_Unaffected verifies that requests with no
// authentication header pass through wrapTokenHeaderSizeHandler unchanged.
func TestTokenHeader_NoAuthHeader_Unaffected(t *testing.T) {
	t.Parallel()

	client, addr := newTestClusterForTokenHeader(t, nil)

	req, err := http.NewRequest(http.MethodGet, addr+"/v1/auth/token/lookup-self", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusForbidden, resp.StatusCode,
		"requests with no auth header must not be rejected by the size guard")
}

// BenchmarkTokenHeader_ProcessingCost measures the per-request overhead of
// wrapTokenHeaderSizeHandler at increasing token sizes.
func BenchmarkTokenHeader_ProcessingCost(b *testing.B) {
	core, _, _ := vault.TestCoreUnsealed(b)
	ln, addr := TestServer(b, core)
	defer ln.Close()

	sizes := []struct {
		label string
		bytes int
	}{
		{"1KB", 1 * 1024},
		{"8KB", 8 * 1024},   // at default limit
		{"64KB", 64 * 1024}, // above default limit — fast-rejected by wrapTokenHeaderSizeHandler
		{"512KB", 512 * 1024},
	}

	client := cleanhttp.DefaultClient()
	client.Timeout = 30 * time.Second

	for _, sz := range sizes {
		token := strings.Repeat("x", sz.bytes)
		b.Run(fmt.Sprintf("size=%s", sz.label), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				req, _ := http.NewRequest(http.MethodGet, addr+"/v1/auth/token/lookup-self", nil)
				req.Header.Set(consts.AuthHeaderName, token)
				resp, err := client.Do(req)
				if err == nil {
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
			}
		})
	}
}

// TestTokenHeaderMaxBytes_NilConfig verifies that a nil listener config
// returns DefaultMaxTokenHeaderSize.
func TestTokenHeaderMaxBytes_NilConfig(t *testing.T) {
	t.Parallel()
	require.Equal(t, DefaultMaxTokenHeaderSize, TokenHeaderMaxBytes(nil))
}

// TestTokenHeaderMaxBytes_ZeroCustomSize verifies that a zero CustomMaxTokenHeaderSize
// falls back to DefaultMaxTokenHeaderSize.
func TestTokenHeaderMaxBytes_ZeroCustomSize(t *testing.T) {
	t.Parallel()
	lnConfig := &configutil.Listener{} // CustomMaxTokenHeaderSize == 0
	require.Equal(t, DefaultMaxTokenHeaderSize, TokenHeaderMaxBytes(lnConfig))
}

// TestTokenHeaderMaxBytes_CustomSize verifies that a positive CustomMaxTokenHeaderSize
// is returned verbatim.
func TestTokenHeaderMaxBytes_CustomSize(t *testing.T) {
	t.Parallel()
	lnConfig := &configutil.Listener{CustomMaxTokenHeaderSize: 4096}
	require.Equal(t, 4096, TokenHeaderMaxBytes(lnConfig))
}

// TestTokenHeaderMaxBytes_Disabled verifies that max_token_header_size = -1
// returns 0, leaving http.Server.MaxHeaderBytes at the Go stdlib default.
func TestTokenHeaderMaxBytes_Disabled(t *testing.T) {
	t.Parallel()
	lnConfig := &configutil.Listener{CustomMaxTokenHeaderSize: -1}
	require.Equal(t, 0, TokenHeaderMaxBytes(lnConfig))
}

// TestTokenHeader_StdlibBackstop_NonAuthHeaderRejected verifies that setting
// MaxHeaderBytes on http.Server rejects oversized non-authentication headers that
// wrapTokenHeaderSizeHandler does not inspect.
func TestTokenHeader_StdlibBackstop_NonAuthHeaderRejected(t *testing.T) {
	t.Parallel()

	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	oversizedValue := strings.Repeat("x", DefaultMaxTokenHeaderSize*2)

	req, err := http.NewRequest(http.MethodGet, addr+"/v1/sys/health", nil)
	require.NoError(t, err)
	req.Header.Set("X-Evil-Header", oversizedValue)

	client := cleanhttp.DefaultClient()
	client.Timeout = 15 * time.Second
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
		require.NotEqual(t, http.StatusOK, resp.StatusCode,
			"oversized non-auth header must not reach the application layer")
	}
}

// TestTokenHeaderMaxBytes_ServerUsesCorrectDefault verifies that an http.Server
// built with TokenHeaderMaxBytes(nil) enforces the limit on any header.
func TestTokenHeaderMaxBytes_ServerUsesCorrectDefault(t *testing.T) {
	t.Parallel()

	reached := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Handler:        inner,
		MaxHeaderBytes: TokenHeaderMaxBytes(nil),
	}

	// Start on a random port.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			t.Errorf("srv.Serve: %v", err)
		}
	}()
	t.Cleanup(func() { srv.Close() })

	addr := "http://" + ln.Addr().String()

	t.Run("within_limit_reaches_handler", func(t *testing.T) {
		reached = false
		req, err := http.NewRequest(http.MethodGet, addr+"/", nil)
		require.NoError(t, err)
		req.Header.Set("X-Test", strings.Repeat("a", DefaultMaxTokenHeaderSize-200))
		resp, err := cleanhttp.DefaultClient().Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.True(t, reached, "handler should have been called for a within-limit request")
	})

	t.Run("exceeds_limit_rejected_by_stdlib", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, addr+"/", nil)
		require.NoError(t, err)
		// 2× to exceed the stdlib's effective limit (MaxHeaderBytes + 4096).
		req.Header.Set("X-Evil", strings.Repeat("x", DefaultMaxTokenHeaderSize*2))
		resp, err := cleanhttp.DefaultClient().Do(req)
		if err == nil {
			defer resp.Body.Close()
			require.NotEqual(t, http.StatusOK, resp.StatusCode,
				"stdlib must reject a request whose headers exceed MaxHeaderBytes")
		}
	})
}

// TestTokenHeader_ClusterListener_SkipsCheck verifies that the cluster listener
// does not re-enforce the token header size limit on forwarded requests.
// This prevents a regression where a user raising CustomMaxTokenHeaderSize above
// the default would see forwarded requests rejected by the active node's cluster
// listener, which only knows the default limit (same failure mode as the JSON
// limits regression fixed by DisableJSONLimitParsing).
func TestTokenHeader_ClusterListener_SkipsCheck(t *testing.T) {
	// A token that exceeds the default limit but would be allowed by a custom
	// API-listener limit of 16 KB. The cluster listener must pass it through
	// without re-checking.
	oversizedForDefault := strings.Repeat("a", DefaultMaxTokenHeaderSize+100)

	// Cluster listener is configured identically to how server.go sets it up:
	// DisableTokenHeaderSizeParsing = true, no CustomMaxTokenHeaderSize override.
	clusterProps := &vault.HandlerProperties{
		ListenerConfig: &configutil.Listener{
			DisableTokenHeaderSizeParsing: true,
		},
	}

	reached := false
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	})

	clusterHandler := wrapTokenHeaderSizeHandler(inner, clusterProps)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	srv := &http.Server{Handler: clusterHandler}
	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			t.Errorf("srv.Serve: %v", err)
		}
	}()
	t.Cleanup(func() { srv.Close() })

	addr := "http://" + ln.Addr().String()

	t.Run("oversized_token_passes_cluster_listener", func(t *testing.T) {
		reached = false
		req, err := http.NewRequest(http.MethodGet, addr+"/", nil)
		require.NoError(t, err)
		req.Header.Set(consts.AuthHeaderName, oversizedForDefault)
		resp, err := cleanhttp.DefaultClient().Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode,
			"cluster listener must not re-enforce the token size limit on forwarded requests")
		require.True(t, reached, "inner handler should be reached on the cluster listener")
	})

	t.Run("api_listener_still_enforces_default_limit", func(t *testing.T) {
		// Confirm the same token is rejected by a normal API-listener handler
		// (no DisableTokenHeaderSizeParsing), so we know the test token really does
		// exceed the default limit.
		apiProps := &vault.HandlerProperties{
			ListenerConfig: &configutil.Listener{},
		}
		apiHandler := wrapTokenHeaderSizeHandler(inner, apiProps)

		apiLn, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		apiSrv := &http.Server{Handler: apiHandler}
		go func() {
			if err := apiSrv.Serve(apiLn); err != nil && err != http.ErrServerClosed {
				t.Errorf("apiSrv.Serve: %v", err)
			}
		}()
		t.Cleanup(func() { apiSrv.Close() })

		req, err := http.NewRequest(http.MethodGet, "http://"+apiLn.Addr().String()+"/", nil)
		require.NoError(t, err)
		req.Header.Set(consts.AuthHeaderName, oversizedForDefault)
		resp, err := cleanhttp.DefaultClient().Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusBadRequest, resp.StatusCode,
			"API listener must reject a token that exceeds the default size limit")
	})
}

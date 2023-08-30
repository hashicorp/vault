// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
)

func init() {
	// Ensure our special envvars are not present
	os.Setenv("VAULT_ADDR", "")
	os.Setenv("VAULT_TOKEN", "")
}

func TestDefaultConfig_envvar(t *testing.T) {
	os.Setenv("VAULT_ADDR", "https://vault.mycompany.com")
	defer os.Setenv("VAULT_ADDR", "")

	config := DefaultConfig()
	if config.Address != "https://vault.mycompany.com" {
		t.Fatalf("bad: %s", config.Address)
	}

	os.Setenv("VAULT_TOKEN", "testing")
	defer os.Setenv("VAULT_TOKEN", "")

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if token := client.Token(); token != "testing" {
		t.Fatalf("bad: %s", token)
	}
}

func TestClientDefaultHttpClient(t *testing.T) {
	_, err := NewClient(&Config{
		HttpClient: http.DefaultClient,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientNilConfig(t *testing.T) {
	client, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("expected a non-nil client")
	}
}

func TestClientDefaultHttpClient_unixSocket(t *testing.T) {
	os.Setenv("VAULT_AGENT_ADDR", "unix:///var/run/vault.sock")
	defer os.Setenv("VAULT_AGENT_ADDR", "")

	client, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	if client == nil {
		t.Fatal("expected a non-nil client")
	}
	if client.addr.Scheme != "http" {
		t.Fatalf("bad: %s", client.addr.Scheme)
	}
	if client.addr.Host != "/var/run/vault.sock" {
		t.Fatalf("bad: %s", client.addr.Host)
	}
}

func TestClientSetAddress(t *testing.T) {
	client, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Start with TCP address using HTTP
	if err := client.SetAddress("http://172.168.2.1:8300"); err != nil {
		t.Fatal(err)
	}
	if client.addr.Host != "172.168.2.1:8300" {
		t.Fatalf("bad: expected: '172.168.2.1:8300' actual: %q", client.addr.Host)
	}
	// Test switching to Unix Socket address from TCP address
	if err := client.SetAddress("unix:///var/run/vault.sock"); err != nil {
		t.Fatal(err)
	}
	if client.addr.Scheme != "http" {
		t.Fatalf("bad: expected: 'http' actual: %q", client.addr.Scheme)
	}
	if client.addr.Host != "/var/run/vault.sock" {
		t.Fatalf("bad: expected: '/var/run/vault.sock' actual: %q", client.addr.Host)
	}
	if client.addr.Path != "" {
		t.Fatalf("bad: expected '' actual: %q", client.addr.Path)
	}
	if client.config.HttpClient.Transport.(*http.Transport).DialContext == nil {
		t.Fatal("bad: expected DialContext to not be nil")
	}
	// Test switching to TCP address from Unix Socket address
	if err := client.SetAddress("http://172.168.2.1:8300"); err != nil {
		t.Fatal(err)
	}
	if client.addr.Host != "172.168.2.1:8300" {
		t.Fatalf("bad: expected: '172.168.2.1:8300' actual: %q", client.addr.Host)
	}
	if client.addr.Scheme != "http" {
		t.Fatalf("bad: expected: 'http' actual: %q", client.addr.Scheme)
	}
}

func TestClientToken(t *testing.T) {
	tokenValue := "foo"
	handler := func(w http.ResponseWriter, req *http.Request) {}

	config, ln := testHTTPServer(t, http.HandlerFunc(handler))
	defer ln.Close()

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	client.SetToken(tokenValue)

	// Verify the token is set
	if v := client.Token(); v != tokenValue {
		t.Fatalf("bad: %s", v)
	}

	client.ClearToken()

	if v := client.Token(); v != "" {
		t.Fatalf("bad: %s", v)
	}
}

func TestClientHostHeader(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(req.Host))
	}
	config, ln := testHTTPServer(t, http.HandlerFunc(handler))
	defer ln.Close()

	config.Address = strings.ReplaceAll(config.Address, "127.0.0.1", "localhost")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Set the token manually
	client.SetToken("foo")

	resp, err := client.RawRequest(client.NewRequest(http.MethodPut, "/"))
	if err != nil {
		t.Fatal(err)
	}

	// Copy the response
	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)

	// Verify we got the response from the primary
	if buf.String() != strings.ReplaceAll(config.Address, "http://", "") {
		t.Fatalf("Bad address: %s", buf.String())
	}
}

func TestClientBadToken(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {}

	config, ln := testHTTPServer(t, http.HandlerFunc(handler))
	defer ln.Close()

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	client.SetToken("foo")
	_, err = client.RawRequest(client.NewRequest(http.MethodPut, "/"))
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken("foo\u007f")
	_, err = client.RawRequest(client.NewRequest(http.MethodPut, "/"))
	if err == nil || !strings.Contains(err.Error(), "printable") {
		t.Fatalf("expected error due to bad token")
	}
}

func TestClientDisableRedirects(t *testing.T) {
	tests := map[string]struct {
		statusCode       int
		expectedNumReqs  int
		disableRedirects bool
	}{
		"Disabled redirects: Moved permanently":  {statusCode: 301, expectedNumReqs: 1, disableRedirects: true},
		"Disabled redirects: Found":              {statusCode: 302, expectedNumReqs: 1, disableRedirects: true},
		"Disabled redirects: Temporary Redirect": {statusCode: 307, expectedNumReqs: 1, disableRedirects: true},
		"Enable redirects: Moved permanently":    {statusCode: 301, expectedNumReqs: 2, disableRedirects: false},
	}

	for name, tc := range tests {
		test := tc
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			numReqs := 0
			var config *Config

			respFunc := func(w http.ResponseWriter, req *http.Request) {
				// Track how many requests the server has handled
				numReqs++
				// Send back the relevant status code and generate a location
				w.Header().Set("Location", fmt.Sprintf(config.Address+"/reqs/%v", numReqs))
				w.WriteHeader(test.statusCode)
			}

			config, ln := testHTTPServer(t, http.HandlerFunc(respFunc))
			config.DisableRedirects = test.disableRedirects
			defer ln.Close()

			client, err := NewClient(config)
			if err != nil {
				t.Fatalf("%s: error %v", name, err)
			}

			req := client.NewRequest("GET", "/")
			resp, err := client.rawRequestWithContext(context.Background(), req)
			if err != nil {
				t.Fatalf("%s: error %v", name, err)
			}

			if numReqs != test.expectedNumReqs {
				t.Fatalf("%s: expected %v request(s) but got %v", name, test.expectedNumReqs, numReqs)
			}

			if resp.StatusCode != test.statusCode {
				t.Fatalf("%s: expected status code %v got %v", name, test.statusCode, resp.StatusCode)
			}

			location, err := resp.Location()
			if err != nil {
				t.Fatalf("%s error %v", name, err)
			}
			if req.URL.String() == location.String() {
				t.Fatalf("%s: expected request URL %v to be different from redirect URL %v", name, req.URL, resp.Request.URL)
			}
		})
	}
}

func TestClientRedirect(t *testing.T) {
	primary := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("test"))
	}
	config, ln := testHTTPServer(t, http.HandlerFunc(primary))
	defer ln.Close()

	standby := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Location", config.Address)
		w.WriteHeader(307)
	}
	config2, ln2 := testHTTPServer(t, http.HandlerFunc(standby))
	defer ln2.Close()

	client, err := NewClient(config2)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Set the token manually
	client.SetToken("foo")

	// Do a raw "/" request
	resp, err := client.RawRequest(client.NewRequest(http.MethodPut, "/"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Copy the response
	var buf bytes.Buffer
	io.Copy(&buf, resp.Body)

	// Verify we got the response from the primary
	if buf.String() != "test" {
		t.Fatalf("Bad: %s", buf.String())
	}
}

func TestDefaulRetryPolicy(t *testing.T) {
	cases := map[string]struct {
		resp      *http.Response
		err       error
		expect    bool
		expectErr error
	}{
		"retry on error": {
			err:    fmt.Errorf("error"),
			expect: true,
		},
		"don't retry connection failures": {
			err: &url.Error{
				Err: x509.UnknownAuthorityError{},
			},
		},
		"don't retry on 200": {
			resp: &http.Response{
				StatusCode: http.StatusOK,
			},
		},
		"don't retry on 4xx": {
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
			},
		},
		"don't retry on 501": {
			resp: &http.Response{
				StatusCode: http.StatusNotImplemented,
			},
		},
		"retry on 500": {
			resp: &http.Response{
				StatusCode: http.StatusInternalServerError,
			},
			expect: true,
		},
		"retry on 5xx": {
			resp: &http.Response{
				StatusCode: http.StatusGatewayTimeout,
			},
			expect: true,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			retry, err := DefaultRetryPolicy(context.Background(), test.resp, test.err)
			if retry != test.expect {
				t.Fatalf("expected to retry request: '%t', but actual result was: '%t'", test.expect, retry)
			}
			if err != test.expectErr {
				t.Fatalf("expected error from retry policy: %q, but actual result was: %q", err, test.expectErr)
			}
		})
	}
}

func TestClientEnvSettings(t *testing.T) {
	cwd, _ := os.Getwd()

	caCertBytes, err := os.ReadFile(cwd + "/test-fixtures/keys/cert.pem")
	if err != nil {
		t.Fatalf("error reading %q cert file: %v", cwd+"/test-fixtures/keys/cert.pem", err)
	}

	oldCACert := os.Getenv(EnvVaultCACert)
	oldCACertBytes := os.Getenv(EnvVaultCACertBytes)
	oldCAPath := os.Getenv(EnvVaultCAPath)
	oldClientCert := os.Getenv(EnvVaultClientCert)
	oldClientKey := os.Getenv(EnvVaultClientKey)
	oldSkipVerify := os.Getenv(EnvVaultSkipVerify)
	oldMaxRetries := os.Getenv(EnvVaultMaxRetries)
	oldDisableRedirects := os.Getenv(EnvVaultDisableRedirects)

	os.Setenv(EnvVaultCACert, cwd+"/test-fixtures/keys/cert.pem")
	os.Setenv(EnvVaultCACertBytes, string(caCertBytes))
	os.Setenv(EnvVaultCAPath, cwd+"/test-fixtures/keys")
	os.Setenv(EnvVaultClientCert, cwd+"/test-fixtures/keys/cert.pem")
	os.Setenv(EnvVaultClientKey, cwd+"/test-fixtures/keys/key.pem")
	os.Setenv(EnvVaultSkipVerify, "true")
	os.Setenv(EnvVaultMaxRetries, "5")
	os.Setenv(EnvVaultDisableRedirects, "true")

	defer func() {
		os.Setenv(EnvVaultCACert, oldCACert)
		os.Setenv(EnvVaultCACertBytes, oldCACertBytes)
		os.Setenv(EnvVaultCAPath, oldCAPath)
		os.Setenv(EnvVaultClientCert, oldClientCert)
		os.Setenv(EnvVaultClientKey, oldClientKey)
		os.Setenv(EnvVaultSkipVerify, oldSkipVerify)
		os.Setenv(EnvVaultMaxRetries, oldMaxRetries)
		os.Setenv(EnvVaultDisableRedirects, oldDisableRedirects)
	}()

	config := DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		t.Fatalf("error reading environment: %v", err)
	}

	tlsConfig := config.HttpClient.Transport.(*http.Transport).TLSClientConfig
	if len(tlsConfig.RootCAs.Subjects()) == 0 {
		t.Fatalf("bad: expected a cert pool with at least one subject")
	}
	if tlsConfig.GetClientCertificate == nil {
		t.Fatalf("bad: expected client tls config to have a certificate getter")
	}
	if tlsConfig.InsecureSkipVerify != true {
		t.Fatalf("bad: %v", tlsConfig.InsecureSkipVerify)
	}
	if config.DisableRedirects != true {
		t.Fatalf("bad: expected disable redirects to be true: %v", config.DisableRedirects)
	}
}

func TestClientDeprecatedEnvSettings(t *testing.T) {
	oldInsecure := os.Getenv(EnvVaultInsecure)
	os.Setenv(EnvVaultInsecure, "true")
	defer os.Setenv(EnvVaultInsecure, oldInsecure)

	config := DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		t.Fatalf("error reading environment: %v", err)
	}

	tlsConfig := config.HttpClient.Transport.(*http.Transport).TLSClientConfig
	if tlsConfig.InsecureSkipVerify != true {
		t.Fatalf("bad: %v", tlsConfig.InsecureSkipVerify)
	}
}

func TestClientEnvNamespace(t *testing.T) {
	var seenNamespace string
	handler := func(w http.ResponseWriter, req *http.Request) {
		seenNamespace = req.Header.Get(NamespaceHeaderName)
	}
	config, ln := testHTTPServer(t, http.HandlerFunc(handler))
	defer ln.Close()

	oldVaultNamespace := os.Getenv(EnvVaultNamespace)
	defer os.Setenv(EnvVaultNamespace, oldVaultNamespace)
	os.Setenv(EnvVaultNamespace, "test")

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	_, err = client.RawRequest(client.NewRequest(http.MethodGet, "/"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if seenNamespace != "test" {
		t.Fatalf("Bad: %s", seenNamespace)
	}
}

func TestParsingRateAndBurst(t *testing.T) {
	var (
		correctFormat                    = "400:400"
		observedRate, observedBurst, err = parseRateLimit(correctFormat)
		expectedRate, expectedBurst      = float64(400), 400
	)
	if err != nil {
		t.Error(err)
	}
	if expectedRate != observedRate {
		t.Errorf("Expected rate %v but found %v", expectedRate, observedRate)
	}
	if expectedBurst != observedBurst {
		t.Errorf("Expected burst %v but found %v", expectedBurst, observedBurst)
	}
}

func TestParsingRateOnly(t *testing.T) {
	var (
		correctFormat                    = "400"
		observedRate, observedBurst, err = parseRateLimit(correctFormat)
		expectedRate, expectedBurst      = float64(400), 400
	)
	if err != nil {
		t.Error(err)
	}
	if expectedRate != observedRate {
		t.Errorf("Expected rate %v but found %v", expectedRate, observedRate)
	}
	if expectedBurst != observedBurst {
		t.Errorf("Expected burst %v but found %v", expectedBurst, observedBurst)
	}
}

func TestParsingErrorCase(t *testing.T) {
	incorrectFormat := "foobar"
	_, _, err := parseRateLimit(incorrectFormat)
	if err == nil {
		t.Error("Expected error, found no error")
	}
}

func TestClientTimeoutSetting(t *testing.T) {
	oldClientTimeout := os.Getenv(EnvVaultClientTimeout)
	os.Setenv(EnvVaultClientTimeout, "10")
	defer os.Setenv(EnvVaultClientTimeout, oldClientTimeout)
	config := DefaultConfig()
	config.ReadEnvironment()
	_, err := NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (rt roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

func TestClientNonTransportRoundTripper(t *testing.T) {
	client := &http.Client{
		Transport: roundTripperFunc(http.DefaultTransport.RoundTrip),
	}

	_, err := NewClient(&Config{
		HttpClient: client,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientNonTransportRoundTripperUnixAddress(t *testing.T) {
	client := &http.Client{
		Transport: roundTripperFunc(http.DefaultTransport.RoundTrip),
	}

	_, err := NewClient(&Config{
		HttpClient: client,
		Address:    "unix:///var/run/vault.sock",
	})
	if err == nil {
		t.Fatal("bad: expected error got nil")
	}
}

func TestClone(t *testing.T) {
	type fields struct{}
	tests := []struct {
		name    string
		config  *Config
		headers *http.Header
		token   string
	}{
		{
			name:   "default",
			config: DefaultConfig(),
		},
		{
			name: "cloneHeaders",
			config: &Config{
				CloneHeaders: true,
			},
			headers: &http.Header{
				"X-foo": []string{"bar"},
				"X-baz": []string{"qux"},
			},
		},
		{
			name: "preventStaleReads",
			config: &Config{
				ReadYourWrites: true,
			},
		},
		{
			name: "cloneToken",
			config: &Config{
				CloneToken: true,
			},
			token: "cloneToken",
		},
		{
			name: "cloneTLSConfig-enabled",
			config: &Config{
				CloneTLSConfig: true,
				clientTLSConfig: &tls.Config{
					ServerName: "foo.bar.baz",
				},
			},
		},
		{
			name: "cloneTLSConfig-disabled",
			config: &Config{
				CloneTLSConfig: false,
				clientTLSConfig: &tls.Config{
					ServerName: "foo.bar.baz",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent, err := NewClient(tt.config)
			if err != nil {
				t.Fatalf("NewClient failed: %v", err)
			}

			// Set all of the things that we provide setter methods for, which modify config values
			err = parent.SetAddress("http://example.com:8080")
			if err != nil {
				t.Fatalf("SetAddress failed: %v", err)
			}

			clientTimeout := time.Until(time.Now().AddDate(0, 0, 1))
			parent.SetClientTimeout(clientTimeout)

			checkRetry := func(ctx context.Context, resp *http.Response, err error) (bool, error) {
				return true, nil
			}
			parent.SetCheckRetry(checkRetry)

			parent.SetLogger(hclog.NewNullLogger())

			parent.SetLimiter(5.0, 10)
			parent.SetMaxRetries(5)
			parent.SetOutputCurlString(true)
			parent.SetOutputPolicy(true)
			parent.SetSRVLookup(true)

			if tt.headers != nil {
				parent.SetHeaders(*tt.headers)
			}

			if tt.token != "" {
				parent.SetToken(tt.token)
			}

			clone, err := parent.Clone()
			if err != nil {
				t.Fatalf("Clone failed: %v", err)
			}

			if parent.Address() != clone.Address() {
				t.Fatalf("addresses don't match: %v vs %v", parent.Address(), clone.Address())
			}
			if parent.ClientTimeout() != clone.ClientTimeout() {
				t.Fatalf("timeouts don't match: %v vs %v", parent.ClientTimeout(), clone.ClientTimeout())
			}
			if parent.CheckRetry() != nil && clone.CheckRetry() == nil {
				t.Fatal("checkRetry functions don't match. clone is nil.")
			}
			if (parent.Limiter() != nil && clone.Limiter() == nil) || (parent.Limiter() == nil && clone.Limiter() != nil) {
				t.Fatalf("limiters don't match: %v vs %v", parent.Limiter(), clone.Limiter())
			}
			if parent.Limiter().Limit() != clone.Limiter().Limit() {
				t.Fatalf("limiter limits don't match: %v vs %v", parent.Limiter().Limit(), clone.Limiter().Limit())
			}
			if parent.Limiter().Burst() != clone.Limiter().Burst() {
				t.Fatalf("limiter bursts don't match: %v vs %v", parent.Limiter().Burst(), clone.Limiter().Burst())
			}
			if parent.MaxRetries() != clone.MaxRetries() {
				t.Fatalf("maxRetries don't match: %v vs %v", parent.MaxRetries(), clone.MaxRetries())
			}
			if parent.OutputCurlString() == clone.OutputCurlString() {
				t.Fatalf("outputCurlString was copied over when it shouldn't have been: %v and %v", parent.OutputCurlString(), clone.OutputCurlString())
			}
			if parent.SRVLookup() != clone.SRVLookup() {
				t.Fatalf("SRVLookup doesn't match: %v vs %v", parent.SRVLookup(), clone.SRVLookup())
			}
			if tt.config.CloneHeaders {
				if !reflect.DeepEqual(parent.Headers(), clone.Headers()) {
					t.Fatalf("Headers() don't match: %v vs %v", parent.Headers(), clone.Headers())
				}
				if parent.config.CloneHeaders != clone.config.CloneHeaders {
					t.Fatalf("config.CloneHeaders doesn't match: %v vs %v", parent.config.CloneHeaders, clone.config.CloneHeaders)
				}
				if tt.headers != nil {
					if !reflect.DeepEqual(*tt.headers, clone.Headers()) {
						t.Fatalf("expected headers %v, actual %v", *tt.headers, clone.Headers())
					}
				}
			}
			if tt.config.ReadYourWrites && parent.replicationStateStore == nil {
				t.Fatalf("replicationStateStore is nil")
			}
			if tt.config.CloneToken {
				if tt.token == "" {
					t.Fatalf("test requires a non-empty token")
				}
				if parent.config.CloneToken != clone.config.CloneToken {
					t.Fatalf("config.CloneToken doesn't match: %v vs %v", parent.config.CloneToken, clone.config.CloneToken)
				}
				if parent.token != clone.token {
					t.Fatalf("tokens do not match: %v vs %v", parent.token, clone.token)
				}
			} else {
				// assumes `VAULT_TOKEN` is unset or has an empty value.
				expected := ""
				if clone.token != expected {
					t.Fatalf("expected clone's token %q, actual %q", expected, clone.token)
				}
			}
			if !reflect.DeepEqual(parent.replicationStateStore, clone.replicationStateStore) {
				t.Fatalf("expected replicationStateStore %v, actual %v", parent.replicationStateStore,
					clone.replicationStateStore)
			}
			if tt.config.CloneTLSConfig {
				if !reflect.DeepEqual(parent.config.TLSConfig(), clone.config.TLSConfig()) {
					t.Fatalf("config.clientTLSConfig doesn't match: %v vs %v",
						parent.config.TLSConfig(), clone.config.TLSConfig())
				}
			} else if tt.config.clientTLSConfig != nil {
				if reflect.DeepEqual(parent.config.TLSConfig(), clone.config.TLSConfig()) {
					t.Fatalf("config.clientTLSConfig should not match: %v vs %v",
						parent.config.TLSConfig(), clone.config.TLSConfig())
				}
			} else {
				if !reflect.DeepEqual(parent.config.TLSConfig(), clone.config.TLSConfig()) {
					t.Fatalf("config.clientTLSConfig doesn't match: %v vs %v",
						parent.config.TLSConfig(), clone.config.TLSConfig())
				}
			}
		})
	}
}

// TestCloneWithHeadersNoDeadlock confirms that the cloning of the client doesn't cause
// a deadlock.
// Raised in https://github.com/hashicorp/vault/issues/22393 -- there was a
// potential deadlock caused by running the problematicFunc() function in
// multiple goroutines.
func TestCloneWithHeadersNoDeadlock(t *testing.T) {
	client, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	wg := &sync.WaitGroup{}

	problematicFunc := func() {
		wg.Add(1)
		client.SetCloneToken(true)
		_, err := client.CloneWithHeaders()
		if err != nil {
			t.Fatal(err)
		}
		wg.Done()
	}

	for i := 0; i < 1000; i++ {
		go problematicFunc()
	}
	wg.Wait()
}

// TestCloneNoDeadlock is like TestCloneWithHeadersNoDeadlock but with
// Clone instead of CloneWithHeaders
func TestCloneNoDeadlock(t *testing.T) {
	client, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	wg := &sync.WaitGroup{}

	problematicFunc := func() {
		wg.Add(1)
		client.SetCloneToken(true)
		_, err := client.Clone()
		if err != nil {
			t.Fatal(err)
		}
		wg.Done()
	}

	for i := 0; i < 1000; i++ {
		go problematicFunc()
	}
	wg.Wait()
}

func TestSetHeadersRaceSafe(t *testing.T) {
	client, err1 := NewClient(nil)
	if err1 != nil {
		t.Fatalf("NewClient failed: %v", err1)
	}

	start := make(chan interface{})
	done := make(chan interface{})

	testPairs := map[string]string{
		"soda":    "rootbeer",
		"veggie":  "carrots",
		"fruit":   "apples",
		"color":   "red",
		"protein": "egg",
	}

	for key, value := range testPairs {
		tmpKey := key
		tmpValue := value
		go func() {
			<-start
			// This test fails if here, you replace client.AddHeader(tmpKey, tmpValue) with:
			// 	headerCopy := client.Header()
			// 	headerCopy.AddHeader(tmpKey, tmpValue)
			// 	client.SetHeader(headerCopy)
			client.AddHeader(tmpKey, tmpValue)
			done <- true
		}()
	}

	// Start everyone at once.
	close(start)

	// Wait until everyone is done.
	for i := 0; i < len(testPairs); i++ {
		<-done
	}

	// Check that all the test pairs are in the resulting
	// headers.
	resultingHeaders := client.Headers()
	for key, value := range testPairs {
		if resultingHeaders.Get(key) != value {
			t.Fatal("expected " + value + " for " + key)
		}
	}
}

func TestMergeReplicationStates(t *testing.T) {
	type testCase struct {
		name     string
		old      []string
		new      string
		expected []string
	}

	testCases := []testCase{
		{
			name:     "empty-old",
			old:      nil,
			new:      "v1:cid:1:0:",
			expected: []string{"v1:cid:1:0:"},
		},
		{
			name:     "old-smaller",
			old:      []string{"v1:cid:1:0:"},
			new:      "v1:cid:2:0:",
			expected: []string{"v1:cid:2:0:"},
		},
		{
			name:     "old-bigger",
			old:      []string{"v1:cid:2:0:"},
			new:      "v1:cid:1:0:",
			expected: []string{"v1:cid:2:0:"},
		},
		{
			name:     "mixed-single",
			old:      []string{"v1:cid:1:0:"},
			new:      "v1:cid:0:1:",
			expected: []string{"v1:cid:0:1:", "v1:cid:1:0:"},
		},
		{
			name:     "mixed-single-alt",
			old:      []string{"v1:cid:0:1:"},
			new:      "v1:cid:1:0:",
			expected: []string{"v1:cid:0:1:", "v1:cid:1:0:"},
		},
		{
			name:     "mixed-double",
			old:      []string{"v1:cid:0:1:", "v1:cid:1:0:"},
			new:      "v1:cid:2:0:",
			expected: []string{"v1:cid:0:1:", "v1:cid:2:0:"},
		},
		{
			name:     "newer-both",
			old:      []string{"v1:cid:0:1:", "v1:cid:1:0:"},
			new:      "v1:cid:2:1:",
			expected: []string{"v1:cid:2:1:"},
		},
	}

	b64enc := func(ss []string) []string {
		var ret []string
		for _, s := range ss {
			ret = append(ret, base64.StdEncoding.EncodeToString([]byte(s)))
		}
		return ret
	}
	b64dec := func(ss []string) []string {
		var ret []string
		for _, s := range ss {
			d, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				t.Fatal(err)
			}
			ret = append(ret, string(d))
		}
		return ret
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := b64dec(MergeReplicationStates(b64enc(tc.old), base64.StdEncoding.EncodeToString([]byte(tc.new))))
			if diff := deep.Equal(out, tc.expected); len(diff) != 0 {
				t.Errorf("got=%v, expected=%v, diff=%v", out, tc.expected, diff)
			}
		})
	}
}

func TestReplicationStateStore_recordState(t *testing.T) {
	b64enc := func(s string) string {
		return base64.StdEncoding.EncodeToString([]byte(s))
	}

	tests := []struct {
		name     string
		expected []string
		resp     []*Response
	}{
		{
			name: "single",
			resp: []*Response{
				{
					Response: &http.Response{
						Header: map[string][]string{
							HeaderIndex: {
								b64enc("v1:cid:1:0:"),
							},
						},
					},
				},
			},
			expected: []string{
				b64enc("v1:cid:1:0:"),
			},
		},
		{
			name: "empty",
			resp: []*Response{
				{
					Response: &http.Response{
						Header: map[string][]string{},
					},
				},
			},
			expected: nil,
		},
		{
			name: "multiple",
			resp: []*Response{
				{
					Response: &http.Response{
						Header: map[string][]string{
							HeaderIndex: {
								b64enc("v1:cid:0:1:"),
							},
						},
					},
				},
				{
					Response: &http.Response{
						Header: map[string][]string{
							HeaderIndex: {
								b64enc("v1:cid:1:0:"),
							},
						},
					},
				},
			},
			expected: []string{
				b64enc("v1:cid:0:1:"),
				b64enc("v1:cid:1:0:"),
			},
		},
		{
			name: "duplicates",
			resp: []*Response{
				{
					Response: &http.Response{
						Header: map[string][]string{
							HeaderIndex: {
								b64enc("v1:cid:1:0:"),
							},
						},
					},
				},
				{
					Response: &http.Response{
						Header: map[string][]string{
							HeaderIndex: {
								b64enc("v1:cid:1:0:"),
							},
						},
					},
				},
			},
			expected: []string{
				b64enc("v1:cid:1:0:"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &replicationStateStore{}

			var wg sync.WaitGroup
			for _, r := range tt.resp {
				wg.Add(1)
				go func(r *Response) {
					defer wg.Done()
					w.recordState(r)
				}(r)
			}
			wg.Wait()

			if !reflect.DeepEqual(tt.expected, w.store) {
				t.Errorf("recordState(): expected states %v, actual %v", tt.expected, w.store)
			}
		})
	}
}

func TestReplicationStateStore_requireState(t *testing.T) {
	tests := []struct {
		name     string
		states   []string
		req      []*Request
		expected []string
	}{
		{
			name:   "empty",
			states: []string{},
			req: []*Request{
				{
					Headers: make(http.Header),
				},
			},
			expected: nil,
		},
		{
			name: "basic",
			states: []string{
				"v1:cid:0:1:",
				"v1:cid:1:0:",
			},
			req: []*Request{
				{
					Headers: make(http.Header),
				},
			},
			expected: []string{
				"v1:cid:0:1:",
				"v1:cid:1:0:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &replicationStateStore{
				store: tt.states,
			}

			var wg sync.WaitGroup
			for _, r := range tt.req {
				wg.Add(1)
				go func(r *Request) {
					defer wg.Done()
					store.requireState(r)
				}(r)
			}

			wg.Wait()

			var actual []string
			for _, r := range tt.req {
				if values := r.Headers.Values(HeaderIndex); len(values) > 0 {
					actual = append(actual, values...)
				}
			}
			sort.Strings(actual)
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("requireState(): expected states %v, actual %v", tt.expected, actual)
			}
		})
	}
}

func TestClient_ReadYourWrites(t *testing.T) {
	b64enc := func(s string) string {
		return base64.StdEncoding.EncodeToString([]byte(s))
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set(HeaderIndex, strings.TrimLeft(req.URL.Path, "/"))
	})

	tests := []struct {
		name       string
		handler    http.Handler
		wantStates []string
		values     [][]string
		clone      bool
	}{
		{
			name:    "multiple_duplicates",
			clone:   false,
			handler: handler,
			wantStates: []string{
				b64enc("v1:cid:0:4:"),
			},
			values: [][]string{
				{
					b64enc("v1:cid:0:4:"),
					b64enc("v1:cid:0:2:"),
				},
				{
					b64enc("v1:cid:0:4:"),
					b64enc("v1:cid:0:2:"),
				},
			},
		},
		{
			name:    "basic_clone",
			clone:   true,
			handler: handler,
			wantStates: []string{
				b64enc("v1:cid:0:4:"),
			},
			values: [][]string{
				{
					b64enc("v1:cid:0:4:"),
				},
				{
					b64enc("v1:cid:0:3:"),
				},
			},
		},
		{
			name:    "multiple_clone",
			clone:   true,
			handler: handler,
			wantStates: []string{
				b64enc("v1:cid:0:4:"),
			},
			values: [][]string{
				{
					b64enc("v1:cid:0:4:"),
					b64enc("v1:cid:0:2:"),
				},
				{
					b64enc("v1:cid:0:3:"),
					b64enc("v1:cid:0:1:"),
				},
			},
		},
		{
			name:    "multiple_duplicates_clone",
			clone:   true,
			handler: handler,
			wantStates: []string{
				b64enc("v1:cid:0:4:"),
			},
			values: [][]string{
				{
					b64enc("v1:cid:0:4:"),
					b64enc("v1:cid:0:2:"),
				},
				{
					b64enc("v1:cid:0:4:"),
					b64enc("v1:cid:0:2:"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRequest := func(client *Client, val string) {
				req := client.NewRequest(http.MethodGet, "/"+val)
				req.Headers.Set(HeaderIndex, val)
				resp, err := client.RawRequestWithContext(context.Background(), req)
				if err != nil {
					t.Fatal(err)
				}

				// validate that the server provided a valid header value in its response
				actual := resp.Header.Get(HeaderIndex)
				if actual != val {
					t.Errorf("expected header value %v, actual %v", val, actual)
				}
			}

			config, ln := testHTTPServer(t, handler)
			defer ln.Close()

			config.ReadYourWrites = true
			config.Address = fmt.Sprintf("http://%s", ln.Addr())
			parent, err := NewClient(config)
			if err != nil {
				t.Fatal(err)
			}

			var wg sync.WaitGroup
			for i := 0; i < len(tt.values); i++ {
				var c *Client
				if tt.clone {
					c, err = parent.Clone()
					if err != nil {
						t.Fatal(err)
					}
				} else {
					c = parent
				}

				for _, val := range tt.values[i] {
					wg.Add(1)
					go func(val string) {
						defer wg.Done()
						testRequest(c, val)
					}(val)
				}
			}

			wg.Wait()

			if !reflect.DeepEqual(tt.wantStates, parent.replicationStateStore.states()) {
				t.Errorf("expected states %v, actual %v", tt.wantStates, parent.replicationStateStore.states())
			}
		})
	}
}

func TestClient_SetReadYourWrites(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		calls  []bool
	}{
		{
			name:   "false",
			config: &Config{},
			calls:  []bool{false},
		},
		{
			name:   "true",
			config: &Config{},
			calls:  []bool{true},
		},
		{
			name:   "multi-false",
			config: &Config{},
			calls:  []bool{false, false},
		},
		{
			name:   "multi-true",
			config: &Config{},
			calls:  []bool{true, true},
		},
		{
			name:   "multi-mix",
			config: &Config{},
			calls:  []bool{false, true, false, true},
		},
	}

	assertSetReadYourRights := func(t *testing.T, c *Client, v bool, s *replicationStateStore) {
		t.Helper()
		c.SetReadYourWrites(v)
		if c.config.ReadYourWrites != v {
			t.Fatalf("expected config.ReadYourWrites %#v, actual %#v", v, c.config.ReadYourWrites)
		}
		if !reflect.DeepEqual(s, c.replicationStateStore) {
			t.Fatalf("expected replicationStateStore %#v, actual %#v", s, c.replicationStateStore)
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config: tt.config,
			}
			for i, v := range tt.calls {
				var expectStateStore *replicationStateStore
				if v {
					if c.replicationStateStore == nil {
						c.replicationStateStore = &replicationStateStore{
							store: []string{},
						}
					}
					c.replicationStateStore.store = append(c.replicationStateStore.store,
						fmt.Sprintf("%s-%d", tt.name, i))
					expectStateStore = c.replicationStateStore
				}
				assertSetReadYourRights(t, c, v, expectStateStore)
			}
		})
	}
}

func TestClient_SetCloneToken(t *testing.T) {
	tests := []struct {
		name  string
		calls []bool
	}{
		{
			name:  "false",
			calls: []bool{false},
		},
		{
			name:  "true",
			calls: []bool{true},
		},
		{
			name:  "multi",
			calls: []bool{true, false, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config: &Config{},
			}

			var expected bool
			for _, v := range tt.calls {
				actual := c.CloneToken()
				if expected != actual {
					t.Fatalf("expected %v, actual %v", expected, actual)
				}

				expected = v
				c.SetCloneToken(expected)
				actual = c.CloneToken()
				if actual != expected {
					t.Fatalf("SetCloneToken(): expected %v, actual %v", expected, actual)
				}
			}
		})
	}
}

func TestClientWithNamespace(t *testing.T) {
	var ns string
	handler := func(w http.ResponseWriter, req *http.Request) {
		ns = req.Header.Get(NamespaceHeaderName)
	}
	config, ln := testHTTPServer(t, http.HandlerFunc(handler))
	defer ln.Close()

	// set up a client with a namespace
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	ogNS := "test"
	client.SetNamespace(ogNS)
	_, err = client.rawRequestWithContext(
		context.Background(),
		client.NewRequest(http.MethodGet, "/"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if ns != ogNS {
		t.Fatalf("Expected namespace: %q, got %q", ogNS, ns)
	}

	// make a call with a temporary namespace
	newNS := "new-namespace"
	_, err = client.WithNamespace(newNS).rawRequestWithContext(
		context.Background(),
		client.NewRequest(http.MethodGet, "/"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if ns != newNS {
		t.Fatalf("Expected new namespace: %q, got %q", newNS, ns)
	}
	// ensure client has not been modified
	_, err = client.rawRequestWithContext(
		context.Background(),
		client.NewRequest(http.MethodGet, "/"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if ns != ogNS {
		t.Fatalf("Expected original namespace: %q, got %q", ogNS, ns)
	}

	// make call with empty ns
	_, err = client.WithNamespace("").rawRequestWithContext(
		context.Background(),
		client.NewRequest(http.MethodGet, "/"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if ns != "" {
		t.Fatalf("Expected no namespace, got %q", ns)
	}

	// ensure client has not been modified
	if client.Namespace() != ogNS {
		t.Fatalf("Expected original namespace: %q, got %q", ogNS, client.Namespace())
	}
}

func TestVaultProxy(t *testing.T) {
	const NoProxy string = "NO_PROXY"

	tests := map[string]struct {
		name                     string
		vaultHttpProxy           string
		vaultProxyAddr           string
		noProxy                  string
		requestUrl               string
		expectedResolvedProxyUrl string
	}{
		"VAULT_HTTP_PROXY used when NO_PROXY env var doesn't include request host": {
			vaultHttpProxy: "https://hashicorp.com",
			vaultProxyAddr: "",
			noProxy:        "terraform.io",
			requestUrl:     "https://vaultproject.io",
		},
		"VAULT_HTTP_PROXY used when NO_PROXY env var includes request host": {
			vaultHttpProxy: "https://hashicorp.com",
			vaultProxyAddr: "",
			noProxy:        "terraform.io,vaultproject.io",
			requestUrl:     "https://vaultproject.io",
		},
		"VAULT_PROXY_ADDR used when NO_PROXY env var doesn't include request host": {
			vaultHttpProxy: "",
			vaultProxyAddr: "https://hashicorp.com",
			noProxy:        "terraform.io",
			requestUrl:     "https://vaultproject.io",
		},
		"VAULT_PROXY_ADDR used when NO_PROXY env var includes request host": {
			vaultHttpProxy: "",
			vaultProxyAddr: "https://hashicorp.com",
			noProxy:        "terraform.io,vaultproject.io",
			requestUrl:     "https://vaultproject.io",
		},
		"VAULT_PROXY_ADDR used when VAULT_HTTP_PROXY env var also supplied": {
			vaultHttpProxy:           "https://hashicorp.com",
			vaultProxyAddr:           "https://terraform.io",
			noProxy:                  "",
			requestUrl:               "https://vaultproject.io",
			expectedResolvedProxyUrl: "https://terraform.io",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.vaultHttpProxy != "" {
				oldVaultHttpProxy := os.Getenv(EnvHTTPProxy)
				os.Setenv(EnvHTTPProxy, tc.vaultHttpProxy)
				defer os.Setenv(EnvHTTPProxy, oldVaultHttpProxy)
			}

			if tc.vaultProxyAddr != "" {
				oldVaultProxyAddr := os.Getenv(EnvVaultProxyAddr)
				os.Setenv(EnvVaultProxyAddr, tc.vaultProxyAddr)
				defer os.Setenv(EnvVaultProxyAddr, oldVaultProxyAddr)
			}

			if tc.noProxy != "" {
				oldNoProxy := os.Getenv(NoProxy)
				os.Setenv(NoProxy, tc.noProxy)
				defer os.Setenv(NoProxy, oldNoProxy)
			}

			c := DefaultConfig()
			if c.Error != nil {
				t.Fatalf("Expected no error reading config, found error %v", c.Error)
			}

			r, _ := http.NewRequest("GET", tc.requestUrl, nil)
			proxyUrl, err := c.HttpClient.Transport.(*http.Transport).Proxy(r)
			if err != nil {
				t.Fatalf("Expected no error resolving proxy, found error %v", err)
			}
			if proxyUrl == nil || proxyUrl.String() == "" {
				t.Fatalf("Expected proxy to be resolved but no proxy returned")
			}
			if tc.expectedResolvedProxyUrl != "" && proxyUrl.String() != tc.expectedResolvedProxyUrl {
				t.Fatalf("Expected resolved proxy URL to be %v but was %v", tc.expectedResolvedProxyUrl, proxyUrl.String())
			}
		})
	}
}

func TestParseAddressWithUnixSocket(t *testing.T) {
	address := "unix:///var/run/vault.sock"
	config := DefaultConfig()

	u, err := config.ParseAddress(address)
	if err != nil {
		t.Fatal("Error not expected")
	}
	if u.Scheme != "http" {
		t.Fatal("Scheme not changed to http")
	}
	if u.Host != "/var/run/vault.sock" {
		t.Fatal("Host not changed to socket name")
	}
	if u.Path != "" {
		t.Fatal("Path expected to be blank")
	}
	if config.HttpClient.Transport.(*http.Transport).DialContext == nil {
		t.Fatal("DialContext function not set in config.HttpClient.Transport")
	}
}

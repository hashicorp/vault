package api

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/consts"
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

func TestClientSetAddress(t *testing.T) {
	client, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.SetAddress("http://172.168.2.1:8300"); err != nil {
		t.Fatal(err)
	}
	if client.addr.Host != "172.168.2.1:8300" {
		t.Fatalf("bad: expected: '172.168.2.1:8300' actual: %q", client.addr.Host)
	}
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

	resp, err := client.RawRequest(client.NewRequest("PUT", "/"))
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
	_, err = client.RawRequest(client.NewRequest("PUT", "/"))
	if err != nil {
		t.Fatal(err)
	}

	client.SetToken("foo\u007f")
	_, err = client.RawRequest(client.NewRequest("PUT", "/"))
	if err == nil || !strings.Contains(err.Error(), "printable") {
		t.Fatalf("expected error due to bad token")
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
	resp, err := client.RawRequest(client.NewRequest("PUT", "/"))
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
				t.Fatalf("expected error from retry policy: '%s', but actual result was: '%s'", err, test.expectErr)
			}
		})
	}
}

func TestClientEnvSettings(t *testing.T) {
	cwd, _ := os.Getwd()
	oldCACert := os.Getenv(EnvVaultCACert)
	oldCAPath := os.Getenv(EnvVaultCAPath)
	oldClientCert := os.Getenv(EnvVaultClientCert)
	oldClientKey := os.Getenv(EnvVaultClientKey)
	oldSkipVerify := os.Getenv(EnvVaultSkipVerify)
	oldMaxRetries := os.Getenv(EnvVaultMaxRetries)
	os.Setenv(EnvVaultCACert, cwd+"/test-fixtures/keys/cert.pem")
	os.Setenv(EnvVaultCAPath, cwd+"/test-fixtures/keys")
	os.Setenv(EnvVaultClientCert, cwd+"/test-fixtures/keys/cert.pem")
	os.Setenv(EnvVaultClientKey, cwd+"/test-fixtures/keys/key.pem")
	os.Setenv(EnvVaultSkipVerify, "true")
	os.Setenv(EnvVaultMaxRetries, "5")
	defer os.Setenv(EnvVaultCACert, oldCACert)
	defer os.Setenv(EnvVaultCAPath, oldCAPath)
	defer os.Setenv(EnvVaultClientCert, oldClientCert)
	defer os.Setenv(EnvVaultClientKey, oldClientKey)
	defer os.Setenv(EnvVaultSkipVerify, oldSkipVerify)
	defer os.Setenv(EnvVaultMaxRetries, oldMaxRetries)

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
		seenNamespace = req.Header.Get(consts.NamespaceHeaderName)
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

	_, err = client.RawRequest(client.NewRequest("GET", "/"))
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
	client1, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	// Set all of the things that we provide setter methods for, which modify config values
	err = client1.SetAddress("http://example.com:8080")
	if err != nil {
		t.Fatalf("SetAddress failed: %v", err)
	}

	clientTimeout := time.Until(time.Now().AddDate(0, 0, 1))
	client1.SetClientTimeout(clientTimeout)

	checkRetry := func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		return true, nil
	}
	client1.SetCheckRetry(checkRetry)

	client1.SetLogger(hclog.NewNullLogger())

	client1.SetLimiter(5.0, 10)
	client1.SetMaxRetries(5)
	client1.SetOutputCurlString(true)
	client1.SetSRVLookup(true)

	client2, err := client1.Clone()
	if err != nil {
		t.Fatalf("Clone failed: %v", err)
	}

	if client1.Address() != client2.Address() {
		t.Fatalf("addresses don't match: %v vs %v", client1.Address(), client2.Address())
	}
	if client1.ClientTimeout() != client2.ClientTimeout() {
		t.Fatalf("timeouts don't match: %v vs %v", client1.ClientTimeout(), client2.ClientTimeout())
	}
	if client1.CheckRetry() != nil && client2.CheckRetry() == nil {
		t.Fatal("checkRetry functions don't match. client2 is nil.")
	}
	if (client1.Limiter() != nil && client2.Limiter() == nil) || (client1.Limiter() == nil && client2.Limiter() != nil) {
		t.Fatalf("limiters don't match: %v vs %v", client1.Limiter(), client2.Limiter())
	}
	if client1.Limiter().Limit() != client2.Limiter().Limit() {
		t.Fatalf("limiter limits don't match: %v vs %v", client1.Limiter().Limit(), client2.Limiter().Limit())
	}
	if client1.Limiter().Burst() != client2.Limiter().Burst() {
		t.Fatalf("limiter bursts don't match: %v vs %v", client1.Limiter().Burst(), client2.Limiter().Burst())
	}
	if client1.MaxRetries() != client2.MaxRetries() {
		t.Fatalf("maxRetries don't match: %v vs %v", client1.MaxRetries(), client2.MaxRetries())
	}
	if client1.OutputCurlString() != client2.OutputCurlString() {
		t.Fatalf("outputCurlString doesn't match: %v vs %v", client1.OutputCurlString(), client2.OutputCurlString())
	}
	if client1.SRVLookup() != client2.SRVLookup() {
		t.Fatalf("SRVLookup doesn't match: %v vs %v", client1.SRVLookup(), client2.SRVLookup())
	}
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

func TestParseAddressWithUnixSocket(t *testing.T) {
	address := "unix:///var/run/vault.sock"
	config := DefaultConfig()

	u, err := ParseAddress(address, config)
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

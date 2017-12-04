package api

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
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

func TestClientEnvSettings(t *testing.T) {
	cwd, _ := os.Getwd()
	oldCACert := os.Getenv(EnvVaultCACert)
	oldCAPath := os.Getenv(EnvVaultCAPath)
	oldClientCert := os.Getenv(EnvVaultClientCert)
	oldClientKey := os.Getenv(EnvVaultClientKey)
	oldSkipVerify := os.Getenv(EnvVaultInsecure)
	oldMaxRetries := os.Getenv(EnvVaultMaxRetries)
	os.Setenv(EnvVaultCACert, cwd+"/test-fixtures/keys/cert.pem")
	os.Setenv(EnvVaultCAPath, cwd+"/test-fixtures/keys")
	os.Setenv(EnvVaultClientCert, cwd+"/test-fixtures/keys/cert.pem")
	os.Setenv(EnvVaultClientKey, cwd+"/test-fixtures/keys/key.pem")
	os.Setenv(EnvVaultInsecure, "true")
	os.Setenv(EnvVaultMaxRetries, "5")
	defer os.Setenv(EnvVaultCACert, oldCACert)
	defer os.Setenv(EnvVaultCAPath, oldCAPath)
	defer os.Setenv(EnvVaultClientCert, oldClientCert)
	defer os.Setenv(EnvVaultClientKey, oldClientKey)
	defer os.Setenv(EnvVaultInsecure, oldSkipVerify)
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

func TestClientTimeoutSetting(t *testing.T) {
	oldClientTimeout := os.Getenv(EnvVaultClientTimeout)
	os.Setenv(EnvVaultClientTimeout, "10")
	defer os.Setenv(EnvVaultClientTimeout, oldClientTimeout)
	config := DefaultConfig()
	config.ReadEnvironment()
	client, err := NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	_ = client.NewRequest("PUT", "/")
	if client.config.HttpClient.Timeout != time.Second*10 {
		t.Fatalf("error setting client timeout using env variable")
	}

	// Setting custom client timeout for a new request
	client.SetClientTimeout(time.Second * 20)
	_ = client.NewRequest("PUT", "/")
	if client.config.HttpClient.Timeout != time.Second*20 {
		t.Fatalf("error setting client timeout using SetClientTimeout")
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

func TestClone(t *testing.T) {
	client1, err1 := NewClient(nil)
	if err1 != nil {
		t.Fatalf("NewClient failed: %v", err1)
	}
	client2, err2 := client1.Clone()
	if err2 != nil {
		t.Fatalf("Clone failed: %v", err2)
	}

	_ = client2
}

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

func TestClientToken(t *testing.T) {
	tokenValue := "foo"
	handler := func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    AuthCookieName,
			Value:   tokenValue,
			Expires: time.Now().Add(time.Hour),
		})
	}

	config, ln := testHTTPServer(t, http.HandlerFunc(handler))
	defer ln.Close()

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Should have no token initially
	if v := client.Token(); v != "" {
		t.Fatalf("bad: %s", v)
	}

	// Do a raw "/" request to set the cookie
	if _, err := client.RawRequest(client.NewRequest("GET", "/")); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the token is set
	if v := client.Token(); v != tokenValue {
		t.Fatalf("bad: %s", v)
	}

	client.ClearToken()

	if v := client.Token(); v != "" {
		t.Fatalf("bad: %s", v)
	}
}

func TestClientSetToken(t *testing.T) {
	var tokenValue string
	handler := func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie(AuthCookieName)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		tokenValue = cookie.Value
	}

	config, ln := testHTTPServer(t, http.HandlerFunc(handler))
	defer ln.Close()

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Should have no token initially
	if v := client.Token(); v != "" {
		t.Fatalf("bad: %s", v)
	}

	// Set the cookie manually
	client.SetToken("foo")

	// Do a raw "/" request to get the cookie
	if _, err := client.RawRequest(client.NewRequest("GET", "/")); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the token is set
	if v := client.Token(); v != "foo" {
		t.Fatalf("bad: %s", v)
	}
	if v := tokenValue; v != "foo" {
		t.Fatalf("bad: %s", v)
	}

	client.ClearToken()

	if v := client.Token(); v != "" {
		t.Fatalf("bad: %s", v)
	}
}

func TestClientRedirect(t *testing.T) {
	primary := func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie(AuthCookieName)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if cookie.Value != "foo" {
			t.Fatalf("Bad: %#v", cookie)
		}

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

	// Set the cookie manually
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

package api

import (
	"net/http"
	"testing"
	"time"

	vaultHttp "github.com/hashicorp/vault/http"
)

func TestClientToken(t *testing.T) {
	tokenValue := "foo"
	handler := func(w http.ResponseWriter, req *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    vaultHttp.AuthCookieName,
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
		cookie, err := req.Cookie(vaultHttp.AuthCookieName)
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

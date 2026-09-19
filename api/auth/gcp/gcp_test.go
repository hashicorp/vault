// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package gcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRequestIdentityToken checks how the gce login path talks to the
// metadata server: the token comes back on success, a failed request is
// reported with its status and body instead of being handed on as a jwt,
// and the caller's context is honoured.
func TestRequestIdentityToken(t *testing.T) {
	auth := &GCPAuth{roleName: "my-role"}

	t.Run("returns the token on success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Metadata-Flavor") != "Google" {
				t.Errorf("missing Metadata-Flavor header")
			}
			if got, want := r.URL.Query().Get("audience"), "https://vault.example.com/vault/my-role"; got != want {
				t.Errorf("audience = %q, want %q", got, want)
			}
			if got := r.URL.Query().Get("format"); got != "full" {
				t.Errorf("format = %q, want full", got)
			}
			w.Write([]byte("header.payload.signature"))
		}))
		defer srv.Close()

		jwt, err := auth.requestIdentityToken(context.Background(), srv.Client(), srv.URL, "https://vault.example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if jwt != "header.payload.signature" {
			t.Fatalf("jwt = %q, want the response body", jwt)
		}
	})

	t.Run("reports a failed request instead of returning its body", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "identity token not available for the requested audience", http.StatusNotFound)
		}))
		defer srv.Close()

		jwt, err := auth.requestIdentityToken(context.Background(), srv.Client(), srv.URL, "https://vault.example.com")
		if err == nil {
			t.Fatalf("expected an error, got jwt %q", jwt)
		}
		if jwt != "" {
			t.Fatalf("expected an empty jwt on failure, got %q", jwt)
		}
		for _, want := range []string{"404", "identity token not available for the requested audience"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not mention %q", err, want)
			}
		}
	})

	t.Run("honours a cancelled context", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("header.payload.signature"))
		}))
		defer srv.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if _, err := auth.requestIdentityToken(ctx, srv.Client(), srv.URL, "https://vault.example.com"); err == nil {
			t.Fatal("expected an error from the cancelled context")
		}
	})
}

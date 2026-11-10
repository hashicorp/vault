// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package ibmcloud

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

type fakeIAM struct {
	srv       *httptest.Server
	key       *rsa.PrivateKey
	kid       string
	issuer    string
	discovery int32
	keyFetch  int32
	failKeys  atomic.Bool
}

func newFakeIAM(t *testing.T) *fakeIAM {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	// Deliberately not derived from the endpoint URL, to prove the issuer
	// comes from the discovery document.
	f := &fakeIAM{key: key, kid: "kid-1", issuer: "https://issuer.example.test/identity"}

	mux := http.NewServeMux()
	mux.HandleFunc(discoveryPath, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&f.discovery, 1)
		_ = json.NewEncoder(w).Encode(map[string]string{"issuer": f.issuer})
	})
	mux.HandleFunc(jwksPath, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&f.keyFetch, 1)
		if f.failKeys.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		set := jose.JSONWebKeySet{Keys: []jose.JSONWebKey{{
			Key: &key.PublicKey, KeyID: f.kid, Algorithm: "RS256", Use: "sig",
		}}}
		_ = json.NewEncoder(w).Encode(set)
	})
	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeIAM) sign(t *testing.T, alg jose.SignatureAlgorithm, kid string, claims jwt.Claims) string {
	t.Helper()
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: alg, Key: f.key},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", kid),
	)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := jwt.Signed(signer).Claims(claims).CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func (f *fakeIAM) validClaims() jwt.Claims {
	return jwt.Claims{
		Issuer: f.issuer,
		Expiry: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
}

func TestJWKSVerifier_Verify(t *testing.T) {
	f := newFakeIAM(t)
	v := newJWKSVerifier(f.srv.URL, f.srv.Client())
	ctx := context.Background()

	if err := v.verify(ctx, f.sign(t, jose.RS256, f.kid, f.validClaims())); err != nil {
		t.Fatalf("valid token rejected: %v", err)
	}
	if got := atomic.LoadInt32(&f.discovery); got != 1 {
		t.Fatalf("expected 1 discovery fetch, got %d", got)
	}

	// Cached: no further fetches.
	if err := v.verify(ctx, f.sign(t, jose.RS256, f.kid, f.validClaims())); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&f.keyFetch) != 1 || atomic.LoadInt32(&f.discovery) != 1 {
		t.Fatal("expected cached keys and issuer on second verify")
	}

	cases := map[string]struct {
		alg    jose.SignatureAlgorithm
		kid    string
		mutate func(*jwt.Claims)
	}{
		"wrong issuer": {jose.RS256, f.kid, func(c *jwt.Claims) { c.Issuer = "https://evil.example.test/identity" }},
		"expired":      {jose.RS256, f.kid, func(c *jwt.Claims) { c.Expiry = jwt.NewNumericDate(time.Now().Add(-time.Hour)) }},
		"missing exp":  {jose.RS256, f.kid, func(c *jwt.Claims) { c.Expiry = nil }},
		"not yet":      {jose.RS256, f.kid, func(c *jwt.Claims) { c.NotBefore = jwt.NewNumericDate(time.Now().Add(time.Hour)) }},
		"wrong alg":    {jose.RS384, f.kid, nil},
		"no kid":       {jose.RS256, "", nil},
		"unknown kid":  {jose.RS256, "other", nil},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			claims := f.validClaims()
			if tc.mutate != nil {
				tc.mutate(&claims)
			}
			err := v.verify(ctx, f.sign(t, tc.alg, tc.kid, claims))
			if !errors.Is(err, errInvalidToken) {
				t.Fatalf("expected errInvalidToken, got %v", err)
			}
		})
	}

	t.Run("tampered signature", func(t *testing.T) {
		raw := f.sign(t, jose.RS256, f.kid, f.validClaims())
		raw = raw[:len(raw)-4] + "AAAA"
		if err := v.verify(ctx, raw); !errors.Is(err, errInvalidToken) {
			t.Fatalf("expected errInvalidToken, got %v", err)
		}
	})

	t.Run("malformed", func(t *testing.T) {
		if err := v.verify(ctx, "not-a-jwt"); !errors.Is(err, errInvalidToken) {
			t.Fatalf("expected errInvalidToken, got %v", err)
		}
	})
}

func TestJWKSVerifier_UnknownKidRefetchRateLimited(t *testing.T) {
	f := newFakeIAM(t)
	v := newJWKSVerifier(f.srv.URL, f.srv.Client())
	ctx := context.Background()

	tok := f.sign(t, jose.RS256, "other", f.validClaims())
	for i := 0; i < 3; i++ {
		if err := v.verify(ctx, tok); !errors.Is(err, errInvalidToken) {
			t.Fatalf("expected errInvalidToken, got %v", err)
		}
	}
	if got := atomic.LoadInt32(&f.keyFetch); got != 1 {
		t.Fatalf("expected 1 key fetch within the refetch interval, got %d", got)
	}
}

func TestJWKSVerifier_FetchFailureIsNotInvalidToken(t *testing.T) {
	f := newFakeIAM(t)
	f.failKeys.Store(true)
	v := newJWKSVerifier(f.srv.URL, f.srv.Client())
	ctx := context.Background()

	tok := f.sign(t, jose.RS256, f.kid, f.validClaims())
	err := v.verify(ctx, tok)
	if err == nil || errors.Is(err, errInvalidToken) {
		t.Fatalf("expected infrastructure error, got %v", err)
	}

	// Backoff: an immediate retry must not hit IAM again, and must still be
	// reported as an infrastructure error.
	err = v.verify(ctx, tok)
	if err == nil || errors.Is(err, errInvalidToken) {
		t.Fatalf("expected infrastructure error on retry, got %v", err)
	}
	if got := atomic.LoadInt32(&f.keyFetch); got != 1 {
		t.Fatalf("expected 1 key fetch during backoff, got %d", got)
	}

	// After the interval, recovery works.
	f.failKeys.Store(false)
	v.minRefetchInterval = 0
	if err := v.verify(ctx, tok); err != nil {
		t.Fatalf("expected recovery, got %v", err)
	}
}

func TestIBMCloudIAMClient_ValidateTokenStatusMapping(t *testing.T) {
	for _, tc := range []struct {
		name        string
		status      int
		body        string
		wantInvalid bool
		wantErr     bool
	}{
		{"active", http.StatusOK, `{"active":true,"sub_type":"User"}`, false, false},
		{"inactive", http.StatusOK, `{"active":false}`, true, true},
		{"missing active", http.StatusOK, `{}`, true, true},
		{"bad request", http.StatusBadRequest, `{}`, true, true},
		{"unauthorized", http.StatusUnauthorized, `{}`, true, true},
		{"server error", http.StatusInternalServerError, `{}`, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			c, err := NewIBMCloudIAMClient(srv.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			_, err = c.ValidateToken(context.Background(), "tok")
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if errors.Is(err, errInvalidToken) != tc.wantInvalid {
				t.Fatalf("errInvalidToken = %v, want %v (err %v)", errors.Is(err, errInvalidToken), tc.wantInvalid, err)
			}
		})
	}
}

func TestRoleTypeForSubjectType(t *testing.T) {
	for in, want := range map[string]string{
		"User":      roleTypeAPIKey,
		"ServiceId": roleTypeAPIKey,
		"Profile":   roleTypeTrustedProfile,
	} {
		got, err := roleTypeForSubjectType(in)
		if err != nil || got != want {
			t.Fatalf("%s: got %q, %v; want %q", in, got, err, want)
		}
	}
	if _, err := roleTypeForSubjectType("Other"); err == nil {
		t.Fatal("expected error for unsupported sub_type")
	}
}

func TestMatchCRNPattern(t *testing.T) {
	cases := []struct {
		name    string
		crn     string
		pattern string
		want    bool
	}{
		{
			name:    "Exact match",
			crn:     "crn:v1:bluemix:public:is:us-south:a/abc123::instance:i-12345",
			pattern: "crn:v1:bluemix:public:is:us-south:a/abc123::instance:i-12345",
			want:    true,
		},
		{
			name:    "Wildcard at end",
			crn:     "crn:v1:bluemix:public:is:us-south:a/abc123::instance:i-12345",
			pattern: "crn:v1:bluemix:public:is:us-south:a/abc123::instance:*",
			want:    true,
		},
		{
			name:    "Wildcard in middle",
			crn:     "crn:v1:bluemix:public:is:us-south:a/abc123::instance:i-12345",
			pattern: "crn:v1:*:public:is:us-south:a/abc123::instance:i-12345",
			want:    true,
		},
		{
			name:    "Global wildcard",
			crn:     "crn:v1:bluemix:public:is:us-south:a/abc123::instance:i-12345",
			pattern: "*",
			want:    true,
		},
		{
			name:    "Mismatch exact",
			crn:     "crn:v1:bluemix:public:is:us-south:a/abc123::instance:i-12345",
			pattern: "crn:v1:bluemix:public:is:us-south:a/abc123::instance:i-54321",
			want:    false,
		},
		{
			name:    "Mismatch wildcard prefix",
			crn:     "crn:v1:bluemix:public:is:us-south:a/abc123::instance:i-12345",
			pattern: "crn:v2:*",
			want:    false,
		},
		{
			name:    "Special regex chars handled as literals",
			crn:     "crn:v1:bluemix:public:is:us-south:a/abc123::instance+foo:i-12345",
			pattern: "crn:v1:bluemix:public:is:us-south:a/abc123::instance+foo:*",
			want:    true,
		},
		{
			name:    "Special regex dots handled as literals (no wildcard bypass)",
			crn:     "crn:v1:bluemix:public:is:us-south:a/abc123::instanceXfoo:i-12345",
			pattern: "crn:v1:bluemix:public:is:us-south:a/abc123::instance.foo:*",
			want:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchCRNPattern(tc.crn, tc.pattern)
			if got != tc.want {
				t.Errorf("matchCRNPattern(%q, %q) = %v; want %v", tc.crn, tc.pattern, got, tc.want)
			}
		})
	}
}

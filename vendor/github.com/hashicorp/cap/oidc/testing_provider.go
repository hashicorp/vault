package oidc

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/cap/oidc/internal/strutils"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"
)

// TestProvider is a local http server that supports test provider capabilities
// which makes writing tests much easier.  Much of this TestProvider
// design/implementation comes from Consul's oauthtest package. A big thanks to
// the original package's contributors.
//
// It's important to remember that the TestProvider is stateful (see any of its
// receiver functions that begin with Set*).
//
// Once you've started a TestProvider http server with StartTestProvider(...),
// the following test endpoints are supported:
//
//    * GET /.well-known/openid-configuration    OIDC Discovery
//
//    * GET or POST  /authorize                  OIDC authorization supporting both
//                                               the authorization code flow (with
//                                               optional PKCE) and the implicit
//                                               flow with form_post.
//
//    * POST /token                              OIDC token
//
//    * GET /userinfo                            OAuth UserInfo
//
//    * GET /.well-known/jwks.json               JWKs used to verify issued JWT tokens
//
//  Making requests to these endpoints are facilitated by
//    * TestProvider.HTTPClient which returns an http.Client for making requests.
//    * TestProvider.CACert which the pem-encoded CA certificate used by the HTTPS server.
//
// Runtime Configuration:
//  * Issuer: Addr() returns the the current base URL for the test provider's
//  running webserver, which can be used as an OIDC Issuer for discovery and
//  is also used for the iss claim when issuing JWTs.
//
//  * Relying Party ClientID/ClientSecret: SetClientCreds(...) updates the
//  creds and they are empty by default.
//
//  * Now: SetNowFunc(...) updates the provider's "now" function and time.Now
//  is the default.
//
//  * Expiry: SetExpectedExpiry( exp time.Duration) updates the expiry and
//    now + 5 * time.Second is the default.
//
//  * Signing keys: SetSigningKeys(...) updates the keys and a ECDSA P-256 pair
//  of priv/pub keys are the default with a signing algorithm of ES256
//
//  * Authorization Code: SetExpectedAuthCode(...) updates the auth code
//  required by the /authorize endpoint and the code is empty by default.
//
//  * Authorization Nonce: SetExpectedAuthNonce(...) updates the nonce required
//  by the /authorize endpont and the nonce is empty by default.
//
//  * Allowed RedirectURIs: SetAllowedRedirectURIs(...) updates the allowed
//  redirect URIs and "https://example.com" is the default.
//
//  * Custom Claims: SetCustomClaims(...) updates custom claims added to JWTs issued
//  and the custom claims are empty by default.
//
//  * Audiences: SetCustomAudience(...) updates the audience claim of JWTs issued
//  and the ClientID is the default.
//
//  * Authentication Time (auth_time): SetOmitAuthTimeClaim(...) allows you to
//  turn off/on the inclusion of an auth_time claim in issued JWTs and the claim
//  is included by default.
//
//  * Issuing id_tokens: SetOmitIDTokens(...) allows you to turn off/on the issuing of
//  id_tokens from the /token endpoint.  id_tokens are issued by default.
//
//  * Issuing access_tokens: SetOmitAccessTokens(...) allows you to turn off/on
//  the issuing of access_tokens from the /token endpoint. access_tokens are issued
//  by default.
//
//  * Authorization State: SetExpectedState sets the value for the state parameter
//  returned from the /authorized endpoint
//
//  * Token Responses: SetDisableToken disables the /token endpoint, causing
//  it to return a 401 http status.
//
//  * Implicit Flow Responses: SetDisableImplicit disables implicit flow responses,
//  causing them to return a 401 http status.
//
//  * PKCE verifier: SetPKCEVerifier(oidc.CodeVerifier) sets the PKCE code_verifier
//  and PKCEVerifier() returns the current verifier.
//
//  * UserInfo: SetUserInfoReply sets the UserInfo endpoint response and
//  UserInfoReply() returns the current response.
type TestProvider struct {
	httpServer *httptest.Server
	caCert     string

	jwks                *jose.JSONWebKeySet
	allowedRedirectURIs []string
	replySubject        string
	replyUserinfo       interface{}
	replyExpiry         time.Duration

	mu                sync.Mutex
	clientID          string
	clientSecret      string
	expectedAuthCode  string
	expectedAuthNonce string
	expectedState     string
	customClaims      map[string]interface{}
	customAudiences   []string
	omitAuthTimeClaim bool
	omitIDToken       bool
	omitAccessToken   bool
	disableUserInfo   bool
	disableJWKs       bool
	disableToken      bool
	disableImplicit   bool
	invalidJWKs       bool
	nowFunc           func() time.Time
	pkceVerifier      CodeVerifier

	// privKey *ecdsa.PrivateKey
	privKey crypto.PrivateKey
	pubKey  crypto.PublicKey
	keyID   string
	alg     Alg

	t *testing.T

	client *http.Client
}

// Stop stops the running TestProvider.
func (p *TestProvider) Stop() {
	p.httpServer.Close()
	if p.client != nil {
		p.client.CloseIdleConnections()
	}
}

// StartTestProvider creates and starts a running TestProvider http server.  The
// WithPort option is supported.  The TestProvider will be shutdown when the
// test and all it's subtests complete via a registered function with
// t.Cleanup(...).
func StartTestProvider(t *testing.T, opt ...Option) *TestProvider {
	t.Helper()
	require := require.New(t)
	opts := getTestProviderOpts(opt...)

	v, err := NewCodeVerifier()
	require.NoError(err)
	p := &TestProvider{
		t:            t,
		nowFunc:      time.Now,
		pkceVerifier: v,
		customClaims: map[string]interface{}{},
		replyExpiry:  5 * time.Second,

		allowedRedirectURIs: []string{
			"https://example.com",
		},
		replySubject: "alice@example.com",
		replyUserinfo: map[string]interface{}{
			"sub":           "alice@example.com",
			"dob":           "1978",
			"friend":        "bob",
			"nickname":      "A",
			"advisor":       "Faythe",
			"nosy-neighbor": "Eve",
		},
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(err)
	p.pubKey, p.privKey = &priv.PublicKey, priv
	p.alg = ES256
	p.keyID = strconv.Itoa(int(time.Now().Unix()))
	p.jwks = &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:   p.pubKey,
				KeyID: p.keyID,
			},
		},
	}
	p.httpServer = httptestNewUnstartedServerWithPort(t, p, opts.withPort)
	p.httpServer.Config.ErrorLog = log.New(ioutil.Discard, "", 0)
	p.httpServer.StartTLS()
	t.Cleanup(p.Stop)

	cert := p.httpServer.Certificate()

	var buf bytes.Buffer
	err = pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	require.NoError(err)
	p.caCert = buf.String()

	return p
}

// testProviderOptions is the set of available options for TestProvider
// functions
type testProviderOptions struct {
	withPort     int
	withAtHashOf string
	withCHashOf  string
}

// testProviderDefaults is a handy way to get the defaults at runtime and during unit
// tests.
func testProviderDefaults() testProviderOptions {
	return testProviderOptions{}
}

// getTestProviderOpts gets the test provider defaults and applies the opt
// overrides passed in
func getTestProviderOpts(opt ...Option) testProviderOptions {
	opts := testProviderDefaults()
	ApplyOpts(&opts, opt...)
	return opts
}

// WithTestPort provides an optional port for the test provider.
//
// Valid for: TestProvider.StartTestProvider
func WithTestPort(port int) Option {
	return func(o interface{}) {
		if o, ok := o.(*testProviderOptions); ok {
			o.withPort = port
		}
	}
}

// withTestAtHash provides an option to request the at_hash claim. Valid for:
// TestProvider.issueSignedJWT
func withTestAtHash(accessToken string) Option {
	return func(o interface{}) {
		if o, ok := o.(*testProviderOptions); ok {
			o.withAtHashOf = accessToken
		}
	}
}

// withTestCHash provides an option to request the c_hash claim. Valid for:
// TestProvider.issueSignedJWT
func withTestCHash(authorizationCode string) Option {
	return func(o interface{}) {
		if o, ok := o.(*testProviderOptions); ok {
			o.withCHashOf = authorizationCode
		}
	}
}

// HTTPClient returns an http.Client for the test provider. The returned client
// uses a pooled transport (so it can reuse connections) that uses the
// test provider's CA certificate. This client's idle connections are closed in
// TestProvider.Done()
func (p *TestProvider) HTTPClient() *http.Client {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.client != nil {
		return p.client
	}
	p.t.Helper()
	require := require.New(p.t)

	// use the cleanhttp package to create a "pooled" transport that's better
	// configured for requests that re-use the same provider host.  Among other
	// things, this transport supports better concurrency when making requests
	// to the same host.  On the downside, this transport can leak file
	// descriptors over time, so we'll be sure to call
	// client.CloseIdleConnections() in the TestProvider.Done() to stave that off.
	tr := cleanhttp.DefaultPooledTransport()

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(p.caCert))
	require.True(ok)

	tr.TLSClientConfig = &tls.Config{
		RootCAs: certPool,
	}

	p.client = &http.Client{
		Transport: tr,
	}
	return p.client
}

// SetExpectedExpiry is for configuring the expected expiry for any JWTs issued
// by the provider (the default is 5 seconds)
func (p *TestProvider) SetExpectedExpiry(exp time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.replyExpiry = exp
}

// SetClientCreds is for configuring the relying party client ID and client
// secret information required for the OIDC workflows.
func (p *TestProvider) SetClientCreds(clientID, clientSecret string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clientID = clientID
	p.clientSecret = clientSecret
}

// ClientCreds returns the relying party client information required for the
// OIDC workflows.
func (p *TestProvider) ClientCreds() (clientID, clientSecret string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.clientID, p.clientSecret
}

// SetExpectedAuthCode configures the auth code to return from /auth and the
// allowed auth code for /token.
func (p *TestProvider) SetExpectedAuthCode(code string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.expectedAuthCode = code
}

// SetExpectedAuthNonce configures the nonce value required for /auth.
func (p *TestProvider) SetExpectedAuthNonce(nonce string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.expectedAuthNonce = nonce
}

// SetAllowedRedirectURIs allows you to configure the allowed redirect URIs for
// the OIDC workflow. If not configured a sample of "https://example.com" is
// used.
func (p *TestProvider) SetAllowedRedirectURIs(uris []string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.allowedRedirectURIs = uris
}

// SetCustomClaims lets you set claims to return in the JWT issued by the OIDC
// workflow.
func (p *TestProvider) SetCustomClaims(customClaims map[string]interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.customClaims = customClaims
}

// SetCustomAudience configures what audience value to embed in the JWT issued
// by the OIDC workflow.
func (p *TestProvider) SetCustomAudience(customAudiences ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.customAudiences = customAudiences
}

// SetNowFunc configures how the test provider will determine the current time.  The
// default is time.Now()
func (p *TestProvider) SetNowFunc(n func() time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.t.Helper()
	require := require.New(p.t)
	require.NotNilf(n, "TestProvider.SetNowFunc: time func is nil")
	p.nowFunc = n
}

// SetOmitAuthTimeClaim turn on/off the omitting of an auth_time claim from
// id_tokens from the /token endpoint.  If set to true, the test provider will
// not include the auth_time claim in issued id_tokens from the /token endpoint.
func (p *TestProvider) SetOmitAuthTimeClaim(omitAuthTime bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.omitAuthTimeClaim = omitAuthTime
}

// SetOmitIDTokens turn on/off the omitting of id_tokens from the /token
// endpoint.  If set to true, the test provider will not omit (issue) id_tokens
// from the /token endpoint.
func (p *TestProvider) SetOmitIDTokens(omitIDTokens bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.omitIDToken = omitIDTokens
}

// OmitAccessTokens turn on/off the omitting of access_tokens from the /token
// endpoint.  If set to true, the test provider will not omit (issue)
// access_tokens from the /token endpoint.
func (p *TestProvider) SetOmitAccessTokens(omitAccessTokens bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.omitAccessToken = omitAccessTokens
}

// SetDisableUserInfo makes the userinfo endpoint return 404 and omits it from the
// discovery config.
func (p *TestProvider) SetDisableUserInfo(disable bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.disableUserInfo = disable
}

// SetDisableJWKs makes the JWKs endpoint return 404
func (p *TestProvider) SetDisableJWKs(disable bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.disableJWKs = disable
}

// SetInvalidJWKS makes the JWKs endpoint return an invalid response
func (p *TestProvider) SetInvalidJWKS(invalid bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.invalidJWKs = invalid
}

// SetExpectedState sets the value for the state parameter returned from
// /authorized
func (p *TestProvider) SetExpectedState(s string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.expectedState = s
}

// SetDisableToken makes the /token endpoint return 401
func (p *TestProvider) SetDisableToken(disable bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.disableToken = disable
}

// SetDisableImplicit makes implicit flow responses return 401
func (p *TestProvider) SetDisableImplicit(disable bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.disableImplicit = disable
}

// SetPKCEVerifier sets the PKCE oidc.CodeVerifier
func (p *TestProvider) SetPKCEVerifier(verifier CodeVerifier) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.t.Helper()
	require.NotNil(p.t, verifier)
	p.pkceVerifier = verifier
}

// PKCEVerifier returns the PKCE oidc.CodeVerifier
func (p *TestProvider) PKCEVerifier() CodeVerifier {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.pkceVerifier
}

// SetUserInfoReply sets the UserInfo endpoint response.
func (p *TestProvider) SetUserInfoReply(resp interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.replyUserinfo = resp
}

// SetUserInfoReply sets the UserInfo endpoint response.
func (p *TestProvider) UserInfoReply() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.replyUserinfo
}

// Addr returns the current base URL for the test provider's running webserver,
// which can be used as an OIDC issuer for discovery and is also used for the
// iss claim when issuing JWTs.
func (p *TestProvider) Addr() string { return p.httpServer.URL }

// CACert returns the pem-encoded CA certificate used by the test provider's
// HTTPS server.
func (p *TestProvider) CACert() string { return p.caCert }

// SigningKeys returns the test provider's keys used to sign JWTs, its Alg and
// Key ID.
func (p *TestProvider) SigningKeys() (crypto.PrivateKey, crypto.PublicKey, Alg, string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.privKey, p.pubKey, p.alg, p.keyID
}

// SetSigningKeys sets the test provider's keys and alg used to sign JWTs.
func (p *TestProvider) SetSigningKeys(privKey crypto.PrivateKey, pubKey crypto.PublicKey, alg Alg, KeyID string) {
	const op = "TestProvider.SetSigningKeys"
	p.mu.Lock()
	defer p.mu.Unlock()
	p.t.Helper()
	require := require.New(p.t)
	require.NotNilf(privKey, "%s: private key is nil")
	require.NotNilf(pubKey, "%s: public key is empty")
	require.NotEmptyf(alg, "%s: alg is empty")
	require.NotEmptyf(KeyID, "%s: key id is empty")
	p.privKey = privKey
	p.pubKey = pubKey
	p.alg = alg
	p.keyID = KeyID
	p.jwks = &jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{
			{
				Key:   p.pubKey,
				KeyID: p.keyID,
			},
		},
	}
}

func (p *TestProvider) writeJSON(w http.ResponseWriter, out interface{}) error {
	const op = "TestProvider.writeJSON"
	p.t.Helper()
	require := require.New(p.t)
	require.NotNilf(w, "%s: http.ResponseWriter is nil")
	enc := json.NewEncoder(w)
	return enc.Encode(out)
}

// writeImplicitResponse will write the required form data response for an
// implicit flow response to the OIDC authorize endpoint
func (p *TestProvider) writeImplicitResponse(w http.ResponseWriter, state, redirectURL string) error {
	p.t.Helper()
	require := require.New(p.t)
	require.NotNilf(w, "%s: http.ResponseWriter is nil")

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	const respForm = `
<!DOCTYPE html>
<html lang="en">
<head><title>Submit This Form</title></head>
<body onload="javascript:document.forms[0].submit()">
<form method="post" action="%s">
<input type="hidden" name="state" id="state" value="%s"/>
%s
</form>
</body>
</html>`
	const tokenField = `<input type="hidden" name="%s" id="%s" value="%s"/>
`
	accessToken := p.issueSignedJWT()
	idToken := p.issueSignedJWT(withTestAtHash(accessToken))
	var respTokens strings.Builder
	if !p.omitAccessToken {
		respTokens.WriteString(fmt.Sprintf(tokenField, "access_token", "access_token", accessToken))
	}
	if !p.omitIDToken {
		respTokens.WriteString(fmt.Sprintf(tokenField, "id_token", "id_token", idToken))
	}
	if _, err := w.Write([]byte(fmt.Sprintf(respForm, redirectURL, state, respTokens.String()))); err != nil {
		return err
	}
	return nil
}

func (p *TestProvider) issueSignedJWT(opt ...Option) string {
	opts := getTestProviderOpts(opt...)

	claims := map[string]interface{}{
		"sub":       p.replySubject,
		"iss":       p.Addr(),
		"nbf":       float64(p.nowFunc().Add(-p.replyExpiry).Unix()),
		"exp":       float64(p.nowFunc().Add(p.replyExpiry).Unix()),
		"auth_time": float64(p.nowFunc().Unix()),
		"iat":       float64(p.nowFunc().Unix()),
		"aud":       []string{p.clientID},
	}
	if len(p.customAudiences) != 0 {
		claims["aud"] = append(claims["aud"].([]string), p.customAudiences...)
	}
	if p.expectedAuthNonce != "" {
		p.customClaims["nonce"] = p.expectedAuthNonce
	}
	for k, v := range p.customClaims {
		claims[k] = v
	}
	if opts.withAtHashOf != "" {
		claims["at_hash"] = p.testHash(opts.withAtHashOf)
	}
	if opts.withCHashOf != "" {
		claims["c_hash"] = p.testHash(opts.withCHashOf)
	}
	return TestSignJWT(p.t, p.privKey, string(p.alg), claims, nil)
}

// testHash will generate an hash using a signature algorithm. It is used to
// test at_hash and c_hash id_token claims. This is helpful internally, but
// intentionally not exported.
func (p *TestProvider) testHash(data string) string {
	p.t.Helper()
	require := require.New(p.t)
	require.NotEmptyf(data, "testHash: data to hash is empty")
	var h hash.Hash
	switch p.alg {
	case RS256, ES256, PS256:
		h = sha256.New()
	case RS384, ES384, PS384:
		h = sha512.New384()
	case RS512, ES512, PS512:
		h = sha512.New()
	case EdDSA:
		return "EdDSA-hash"
	default:
		require.FailNowf("", "testHash: unsupported signing algorithm %s", string(p.alg))
	}
	require.NotNil(h)
	_, _ = h.Write([]byte(string(data))) // hash documents that Write will never return an error
	sum := h.Sum(nil)[:h.Size()/2]
	actual := base64.RawURLEncoding.EncodeToString(sum)
	return actual
}

// writeAuthErrorResponse writes a standard OIDC authentication error response.
// See: https://openid.net/specs/openid-connect-core-1_0.html#AuthError
func (p *TestProvider) writeAuthErrorResponse(w http.ResponseWriter, req *http.Request, redirectURL, state, errorCode, errorMessage string) {
	p.t.Helper()
	require := require.New(p.t)
	require.NotNilf(w, "%s: http.ResponseWriter is nil")
	require.NotNilf(req, "%s: http.Request is nil")
	require.NotEmptyf(errorCode, "%s: errorCode is empty")

	// state and error are required error response parameters
	redirectURI := redirectURL +
		"?state=" + url.QueryEscape(state) +
		"&error=" + url.QueryEscape(errorCode)

	if errorMessage != "" {
		// add optional error response parameter
		redirectURI += "&error_description=" + url.QueryEscape(errorMessage)
	}

	http.Redirect(w, req, redirectURI, http.StatusFound)
}

// writeTokenErrorResponse writes a standard OIDC token error response.
// See: https://openid.net/specs/openid-connect-core-1_0.html#TokenErrorResponse
func (p *TestProvider) writeTokenErrorResponse(w http.ResponseWriter, statusCode int, errorCode, errorMessage string) error {
	require := require.New(p.t)
	require.NotNilf(w, "%s: http.ResponseWriter is nil")
	require.NotEmptyf(errorCode, "%s: errorCode is empty")
	require.NotEmptyf(statusCode, "%s: statusCode is empty")

	body := struct {
		Code string `json:"error"`
		Desc string `json:"error_description,omitempty"`
	}{
		Code: errorCode,
		Desc: errorMessage,
	}

	w.WriteHeader(statusCode)
	return p.writeJSON(w, &body)
}

// ServeHTTP implements the test provider's http.Handler.
func (p *TestProvider) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// define all the endpoints supported
	const (
		openidConfiguration = "/.well-known/openid-configuration"
		authorize           = "/authorize"
		token               = "/token"
		userInfo            = "/userinfo"
		wellKnownJwks       = "/.well-known/jwks.json"
	)
	p.mu.Lock()
	defer p.mu.Unlock()

	p.t.Helper()
	require := require.New(p.t)
	require.NotNilf(w, "%s: http.ResponseWriter is nil")
	require.NotNilf(req, "%s: http.Request is nil")

	// set a default Content-Type which will be overridden as needed.
	w.Header().Set("Content-Type", "application/json")

	switch req.URL.Path {
	case openidConfiguration:
		// OIDC Discovery endpoint request
		// See: https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderConfigurationResponse
		if req.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		reply := struct {
			Issuer           string `json:"issuer"`
			AuthEndpoint     string `json:"authorization_endpoint"`
			TokenEndpoint    string `json:"token_endpoint"`
			JWKSURI          string `json:"jwks_uri"`
			UserinfoEndpoint string `json:"userinfo_endpoint,omitempty"`
		}{
			Issuer:           p.Addr(),
			AuthEndpoint:     p.Addr() + authorize,
			TokenEndpoint:    p.Addr() + token,
			JWKSURI:          p.Addr() + wellKnownJwks,
			UserinfoEndpoint: p.Addr() + userInfo,
		}
		if p.disableUserInfo {
			reply.UserinfoEndpoint = ""
		}

		err := p.writeJSON(w, &reply)
		require.NoErrorf(err, "%s: internal error: %w", openidConfiguration, err)

		return
	case authorize:
		// Supports both the authorization code and implicit flows
		// See: https://openid.net/specs/openid-connect-core-1_0.html#AuthorizationEndpoint
		if !strutils.StrListContains([]string{"POST", "GET"}, req.Method) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		err := req.ParseForm()
		require.NoErrorf(err, "%s: internal error: %w", authorize, err)

		respType := req.FormValue("response_type")
		scopes := req.Form["scope"]
		state := req.FormValue("state")
		redirectURI := req.FormValue("redirect_uri")
		respMode := req.FormValue("response_mode")

		if respType != "code" && !strings.Contains(respType, "id_token") {
			p.writeAuthErrorResponse(w, req, redirectURI, state, "unsupported_response_type", "")
			return
		}
		if !strutils.StrListContains(scopes, "openid") {
			p.writeAuthErrorResponse(w, req, redirectURI, state, "invalid_scope", "")
			return
		}

		if p.expectedAuthCode == "" {
			p.writeAuthErrorResponse(w, req, redirectURI, state, "access_denied", "")
			return
		}

		nonce := req.FormValue("nonce")
		if p.expectedAuthNonce != "" && p.expectedAuthNonce != nonce {
			p.writeAuthErrorResponse(w, req, redirectURI, state, "access_denied", "")
			return
		}

		if state == "" {
			p.writeAuthErrorResponse(w, req, redirectURI, state, "invalid_request", "missing state parameter")
			return
		}

		if redirectURI == "" {
			p.writeAuthErrorResponse(w, req, redirectURI, state, "invalid_request", "missing redirect_uri parameter")
			return
		}

		var s string
		switch {
		case p.expectedState != "":
			s = p.expectedState
		default:
			s = state
		}

		if strings.Contains(respType, "id_token") {
			if respMode != "form_post" {
				p.writeAuthErrorResponse(w, req, redirectURI, state, "unsupported_response_mode", "must be form_post")
			}
			if p.disableImplicit {
				p.writeAuthErrorResponse(w, req, redirectURI, state, "access_denied", "")
			}
			err := p.writeImplicitResponse(w, s, redirectURI)
			require.NoErrorf(err, "%s: internal error: %w", token, err)
			return
		}

		redirectURI += "?state=" + url.QueryEscape(s) +
			"&code=" + url.QueryEscape(p.expectedAuthCode)

		http.Redirect(w, req, redirectURI, http.StatusFound)

		return

	case wellKnownJwks:
		if p.disableJWKs {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if p.invalidJWKs {
			_, err := w.Write([]byte("It's not a keyset!"))
			require.NoErrorf(err, "%s: internal error: %w", wellKnownJwks, err)
			return
		}
		if req.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		err := p.writeJSON(w, p.jwks)
		require.NoErrorf(err, "%s: internal error: %w", wellKnownJwks, err)
		return
	case token:
		if p.disableToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if req.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		switch {
		case req.FormValue("grant_type") != "authorization_code":
			_ = p.writeTokenErrorResponse(w, http.StatusBadRequest, "invalid_request", "bad grant_type")
			return
		case !strutils.StrListContains(p.allowedRedirectURIs, req.FormValue("redirect_uri")):
			_ = p.writeTokenErrorResponse(w, http.StatusBadRequest, "invalid_request", "redirect_uri is not allowed")
			return
		case req.FormValue("code") != p.expectedAuthCode:
			_ = p.writeTokenErrorResponse(w, http.StatusUnauthorized, "invalid_grant", "unexpected auth code")
			return
		case req.FormValue("code_verifier") != "" && req.FormValue("code_verifier") != p.pkceVerifier.Verifier():
			_ = p.writeTokenErrorResponse(w, http.StatusUnauthorized, "invalid_verifier", "unexpected verifier")
			return
		}

		accessToken := p.issueSignedJWT()
		idToken := p.issueSignedJWT(withTestAtHash(accessToken), withTestCHash(p.expectedAuthCode))
		reply := struct {
			AccessToken string `json:"access_token,omitempty"`
			IDToken     string `json:"id_token,omitempty"`
		}{
			AccessToken: accessToken,
			IDToken:     idToken,
		}
		if p.omitIDToken {
			reply.IDToken = ""
		}
		if p.omitAccessToken {
			reply.AccessToken = ""
		}

		if err := p.writeJSON(w, &reply); err != nil {
			require.NoErrorf(err, "%s: internal error: %w", token, err)
			return
		}
		return
	case userInfo:
		if p.disableUserInfo {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if req.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err := p.writeJSON(w, p.replyUserinfo); err != nil {
			require.NoErrorf(err, "%s: internal error: %w", userInfo, err)
			return
		}
		return

	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// httptestNewUnstartedServerWithPort is roughly the same as
// httptest.NewUnstartedServer() but allows the caller to explicitly choose the
// port if desired.
func httptestNewUnstartedServerWithPort(t *testing.T, handler http.Handler, port int) *httptest.Server {
	t.Helper()
	require := require.New(t)
	require.NotNil(handler)
	if port == 0 {
		return httptest.NewUnstartedServer(handler)
	}
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	l, err := net.Listen("tcp", addr)
	require.NoError(err)

	return &httptest.Server{
		Listener: l,
		Config:   &http.Server{Handler: handler},
	}
}

package jwt

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/oauth2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// KeySet represents a set of keys that can be used to verify the signatures of JWTs.
// A KeySet is expected to be backed by a set of local or remote keys.
type KeySet interface {

	// VerifySignature parses the given JWT, verifies its signature, and returns the claims in its payload.
	// The given JWT must be of the JWS compact serialization form.
	VerifySignature(ctx context.Context, token string) (claims map[string]interface{}, err error)
}

// jsonWebKeySet verifies JWT signatures using keys obtained from a JWKS URL.
type jsonWebKeySet struct {
	remoteJWKS oidc.KeySet
}

// staticKeySet verifies JWT signatures using local public keys.
type staticKeySet struct {
	publicKeys []crypto.PublicKey
}

// NewOIDCDiscoveryKeySet returns a KeySet that verifies JWT signatures using keys from the
// JSON Web Key Set (JWKS) published in the discovery document at the given issuer URL.
// The client used to obtain the remote keys will verify server certificates using the root
// certificates provided by issuerCAPEM. If issuerCAPEM is not provided, system certificates
// are used.
func NewOIDCDiscoveryKeySet(ctx context.Context, issuer string, issuerCAPEM string) (KeySet, error) {
	if issuer == "" {
		return nil, errors.New("issuer must not be empty")
	}

	// Configure an http client with the given certificates
	caCtx, err := createCAContext(ctx, issuerCAPEM)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	if c, ok := caCtx.Value(oauth2.HTTPClient).(*http.Client); ok {
		client = c
	}

	// Create and send the http request for the OIDC discovery document
	wellKnown := strings.TrimSuffix(issuer, "/") + "/.well-known/openid-configuration"
	req, err := http.NewRequest(http.MethodGet, wellKnown, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req.WithContext(caCtx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body and status code
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %s", resp.Status, body)
	}

	// Unmarshal the response body to obtain the issuer and JWKS URL
	var p struct {
		Issuer  string `json:"issuer"`
		JWKSURL string `json:"jwks_uri"`
	}
	if err := unmarshalResp(resp, body, &p); err != nil {
		return nil, fmt.Errorf("failed to decode OIDC discovery document: %w", err)
	}

	// Ensure that the returned issuer matches what was given by issuer
	if p.Issuer != issuer {
		return nil, fmt.Errorf("issuer did not match the returned issuer, expected %q got %q",
			issuer, p.Issuer)
	}

	return &jsonWebKeySet{
		remoteJWKS: oidc.NewRemoteKeySet(caCtx, p.JWKSURL),
	}, nil
}

// NewJSONWebKeySet returns a KeySet that verifies JWT signatures using keys from the JSON Web
// Key Set (JWKS) at the given jwksURL. The client used to obtain the remote JWKS will verify
// server certificates using the root certificates provided by jwksCAPEM. If jwksCAPEM is not
// provided, system certificates are used.
func NewJSONWebKeySet(ctx context.Context, jwksURL string, jwksCAPEM string) (KeySet, error) {
	if jwksURL == "" {
		return nil, errors.New("jwksURL must not be empty")
	}

	caCtx, err := createCAContext(ctx, jwksCAPEM)
	if err != nil {
		return nil, err
	}

	return &jsonWebKeySet{
		remoteJWKS: oidc.NewRemoteKeySet(caCtx, jwksURL),
	}, nil
}

// VerifySignature parses the given JWT, verifies its signature using JWKS keys, and returns
// the claims in its payload. The given JWT must be of the JWS compact serialization form.
func (ks *jsonWebKeySet) VerifySignature(ctx context.Context, token string) (map[string]interface{}, error) {
	payload, err := ks.remoteJWKS.VerifySignature(ctx, token)
	if err != nil {
		return nil, err
	}

	// Unmarshal payload into a set of all received claims
	allClaims := map[string]interface{}{}
	if err := json.Unmarshal(payload, &allClaims); err != nil {
		return nil, err
	}

	return allClaims, nil
}

// NewStaticKeySet returns a KeySet that verifies JWT signatures using the given publicKeys.
func NewStaticKeySet(publicKeys []crypto.PublicKey) (KeySet, error) {
	if len(publicKeys) == 0 {
		return nil, errors.New("publicKeys must not be empty")
	}

	return &staticKeySet{
		publicKeys: publicKeys,
	}, nil
}

// VerifySignature parses the given JWT, verifies its signature using local public keys, and
// returns the claims in its payload. The given JWT must be of the JWS compact serialization form.
func (ks *staticKeySet) VerifySignature(_ context.Context, token string) (map[string]interface{}, error) {
	parsedJWT, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	var valid bool
	allClaims := map[string]interface{}{}
	for _, key := range ks.publicKeys {
		if err := parsedJWT.Claims(key, &allClaims); err == nil {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.New("no known key successfully validated the token signature")
	}

	return allClaims, nil
}

// ParsePublicKeyPEM is used to parse RSA and ECDSA public keys from PEMs. The given
// data must be of PEM-encoded x509 certificate or PKIX public key forms. It returns
// an *rsa.PublicKey or *ecdsa.PublicKey.
func ParsePublicKeyPEM(data []byte) (crypto.PublicKey, error) {
	block, data := pem.Decode(data)
	if block != nil {
		var rawKey interface{}
		var err error
		if rawKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
				rawKey = cert.PublicKey
			} else {
				return nil, err
			}
		}

		if rsaPublicKey, ok := rawKey.(*rsa.PublicKey); ok {
			return rsaPublicKey, nil
		}
		if ecPublicKey, ok := rawKey.(*ecdsa.PublicKey); ok {
			return ecPublicKey, nil
		}
	}

	return nil, errors.New("data does not contain any valid RSA or ECDSA public keys")
}

// createCAContext returns a context with a custom TLS client that's configured with the root
// certificates from caPEM. If no certificates are configured, the original context is returned.
func createCAContext(ctx context.Context, caPEM string) (context.Context, error) {
	if caPEM == "" {
		return ctx, nil
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(caPEM)); !ok {
		return nil, errors.New("could not parse CA PEM value successfully")
	}

	tr := cleanhttp.DefaultPooledTransport()
	tr.TLSClientConfig = &tls.Config{
		RootCAs: certPool,
	}
	tc := &http.Client{
		Transport: tr,
	}

	caCtx := context.WithValue(ctx, oauth2.HTTPClient, tc)

	return caCtx, nil
}

// unmarshalResp JSON unmarshals the given body into the value pointed to by v.
// If it is unable to JSON unmarshal body into v, then it returns an appropriate
// error based on the Content-Type header of r.
func unmarshalResp(r *http.Response, body []byte, v interface{}) error {
	err := json.Unmarshal(body, &v)
	if err == nil {
		return nil
	}
	ct := r.Header.Get("Content-Type")
	mediaType, _, parseErr := mime.ParseMediaType(ct)
	if parseErr == nil && mediaType == "application/json" {
		return fmt.Errorf("got Content-Type = application/json, but could not unmarshal as JSON: %v", err)
	}
	return fmt.Errorf("expected Content-Type = application/json, got %q: %v", ct, err)
}

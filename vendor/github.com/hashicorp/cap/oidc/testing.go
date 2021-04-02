package oidc

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"hash"
	"math/big"
	"net"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// TestGenerateKeys will generate a test ECDSA P-256 pub/priv key pair.
func TestGenerateKeys(t *testing.T) (crypto.PublicKey, crypto.PrivateKey) {
	t.Helper()
	require := require.New(t)
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(err)
	return &priv.PublicKey, priv
}

// TestSignJWT will bundle the provided claims into a test signed JWT.
func TestSignJWT(t *testing.T, key crypto.PrivateKey, alg string, claims interface{}, keyID []byte) string {
	t.Helper()
	require := require.New(t)

	hdr := map[jose.HeaderKey]interface{}{}
	if keyID != nil {
		hdr["key_id"] = string(keyID)
	}

	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.SignatureAlgorithm(alg), Key: key},
		(&jose.SignerOptions{ExtraHeaders: hdr}).WithType("JWT"),
	)
	require.NoError(err)

	raw, err := jwt.Signed(sig).
		Claims(claims).
		CompactSerialize()
	require.NoError(err)
	return raw
}

// TestGenerateCA will generate a test x509 CA cert, along with it encoded in a
// PEM format.
func TestGenerateCA(t *testing.T, hosts []string) (*x509.Certificate, string) {
	t.Helper()
	require := require.New(t)

	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(err)

	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	keyUsage := x509.KeyUsageDigitalSignature

	validFor := 2 * time.Minute
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	require.NoError(err)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	require.NoError(err)

	c, err := x509.ParseCertificate(derBytes)
	require.NoError(err)

	return c, string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))
}

// testHash will generate an hash using a signature algorithm. It is used to
// test at_hash and c_hash id_token claims. This is helpful internally, but
// intentionally not exported.
func testHash(t *testing.T, signatureAlg Alg, data string) string {
	t.Helper()
	require := require.New(t)
	var h hash.Hash
	switch signatureAlg {
	case RS256, ES256, PS256:
		h = sha256.New()
	case RS384, ES384, PS384:
		h = sha512.New384()
	case RS512, ES512, PS512:
		h = sha512.New()
	case EdDSA:
		return "EdDSA-hash"
	default:
		require.FailNowf("", "testHash: unsupported signing algorithm %s", string(signatureAlg))
	}
	require.NotNil(h)
	_, _ = h.Write([]byte(string(data))) // hash documents that Write will never return an error
	sum := h.Sum(nil)[:h.Size()/2]
	actual := base64.RawURLEncoding.EncodeToString(sum)
	return actual
}

// testDefaultJWT creates a default test JWT and is internally helpful, but for now we won't export it.
func testDefaultJWT(t *testing.T, privKey crypto.PrivateKey, expireIn time.Duration, nonce string, additionalClaims map[string]interface{}) string {
	t.Helper()
	now := float64(time.Now().Unix())
	claims := map[string]interface{}{
		"iss":   "https://example.com/",
		"iat":   now,
		"nbf":   now,
		"exp":   float64(time.Now().Unix()),
		"aud":   []string{"www.example.com"},
		"sub":   "alice@example.com",
		"nonce": nonce,
	}
	for k, v := range additionalClaims {
		claims[k] = v
	}
	testJWT := TestSignJWT(t, privKey, string(ES256), claims, nil)
	return testJWT
}

// testNewConfig creates a new config from the TestProvider. It will set the
// TestProvider's client ID/secret and use the TestProviders signing algorithm
// when building the configuration. This is helpful internally, but
// intentionally not exported.
func testNewConfig(t *testing.T, clientID, clientSecret, allowedRedirectURL string, tp *TestProvider) *Config {
	const op = "testNewConfig"
	t.Helper()
	require := require.New(t)

	require.NotEmptyf(clientID, "%s: client id is empty", op)
	require.NotEmptyf(clientSecret, "%s: client secret is empty", op)
	require.NotEmptyf(allowedRedirectURL, "%s: redirect URL is empty", op)

	tp.SetClientCreds(clientID, clientSecret)
	_, _, alg, _ := tp.SigningKeys()
	c, err := NewConfig(
		tp.Addr(),
		clientID,
		ClientSecret(clientSecret),
		[]Alg{alg},
		[]string{allowedRedirectURL},
		nil,
		WithProviderCA(tp.CACert()),
	)
	require.NoError(err)
	return c
}

// testNewProvider creates a new Provider.  It uses the TestProvider (tp) to properly
// construct the provider's configuration (see testNewConfig). This is helpful internally, but
// intentionally not exported.
func testNewProvider(t *testing.T, clientID, clientSecret, redirectURL string, tp *TestProvider) *Provider {
	const op = "testNewProvider"
	t.Helper()
	require := require.New(t)
	require.NotEmptyf(clientID, "%s: client id is empty", op)
	require.NotEmptyf(clientSecret, "%s: client secret is empty", op)
	require.NotEmptyf(redirectURL, "%s: redirect URL is empty", op)

	tc := testNewConfig(t, clientID, clientSecret, redirectURL, tp)
	p, err := NewProvider(tc)
	require.NoError(err)
	t.Cleanup(p.Done)
	return p
}

// testAssertEqualFunc gives you a way to assert that two functions (passed as
// interface{}) are equal.  This is helpful internally, but intentionally not
// exported.
func testAssertEqualFunc(t *testing.T, wantFunc, gotFunc interface{}, format string, args ...interface{}) {
	t.Helper()
	want := runtime.FuncForPC(reflect.ValueOf(wantFunc).Pointer()).Name()
	got := runtime.FuncForPC(reflect.ValueOf(gotFunc).Pointer()).Name()
	assert.Equalf(t, want, got, format, args...)
}

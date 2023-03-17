package pki

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"hash"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// Setup helpers
func createBackendWithStorage(t testing.TB) (*backend, logical.Storage) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	var err error
	b := Backend(config)
	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}
	// Assume for our tests we have performed the migration already.
	b.pkiStorageVersion.Store(1)
	return b, config.StorageView
}

func mountPKIEndpoint(t testing.TB, client *api.Client, path string) {
	var err error
	err = client.Sys().Mount(path, &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	require.NoError(t, err, "failed mounting pki endpoint")
}

// Signing helpers
func requireSignedBy(t *testing.T, cert *x509.Certificate, key crypto.PublicKey) {
	switch key.(type) {
	case *rsa.PublicKey:
		requireRSASignedBy(t, cert, key.(*rsa.PublicKey))
	case *ecdsa.PublicKey:
		requireECDSASignedBy(t, cert, key.(*ecdsa.PublicKey))
	case ed25519.PublicKey:
		requireED25519SignedBy(t, cert, key.(ed25519.PublicKey))
	default:
		require.Fail(t, "unknown public key type %#v", key)
	}
}

func requireRSASignedBy(t *testing.T, cert *x509.Certificate, key *rsa.PublicKey) {
	require.Contains(t, []x509.SignatureAlgorithm{x509.SHA256WithRSA, x509.SHA512WithRSA},
		cert.SignatureAlgorithm, "only sha256 signatures supported")

	var hasher hash.Hash
	var hashAlgo crypto.Hash

	switch cert.SignatureAlgorithm {
	case x509.SHA256WithRSA:
		hasher = sha256.New()
		hashAlgo = crypto.SHA256
	case x509.SHA512WithRSA:
		hasher = sha512.New()
		hashAlgo = crypto.SHA512
	}

	hasher.Write(cert.RawTBSCertificate)
	hashData := hasher.Sum(nil)

	err := rsa.VerifyPKCS1v15(key, hashAlgo, hashData, cert.Signature)
	require.NoError(t, err, "the certificate was not signed by the expected public rsa key.")
}

func requireECDSASignedBy(t *testing.T, cert *x509.Certificate, key *ecdsa.PublicKey) {
	require.Contains(t, []x509.SignatureAlgorithm{x509.ECDSAWithSHA256, x509.ECDSAWithSHA512},
		cert.SignatureAlgorithm, "only ecdsa signatures supported")

	var hasher hash.Hash
	switch cert.SignatureAlgorithm {
	case x509.ECDSAWithSHA256:
		hasher = sha256.New()
	case x509.ECDSAWithSHA512:
		hasher = sha512.New()
	}

	hasher.Write(cert.RawTBSCertificate)
	hashData := hasher.Sum(nil)

	verify := ecdsa.VerifyASN1(key, hashData, cert.Signature)
	require.True(t, verify, "the certificate was not signed by the expected public ecdsa key.")
}

func requireED25519SignedBy(t *testing.T, cert *x509.Certificate, key ed25519.PublicKey) {
	require.Equal(t, x509.PureEd25519, cert.SignatureAlgorithm)
	ed25519.Verify(key, cert.RawTBSCertificate, cert.Signature)
}

// Certificate helper
func parseCert(t *testing.T, pemCert string) *x509.Certificate {
	block, _ := pem.Decode([]byte(pemCert))
	require.NotNil(t, block, "failed to decode PEM block")

	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)
	return cert
}

func requireMatchingPublicKeys(t *testing.T, cert *x509.Certificate, key crypto.PublicKey) {
	certPubKey := cert.PublicKey
	areEqual, err := certutil.ComparePublicKeysAndType(certPubKey, key)
	require.NoError(t, err, "failed comparing public keys: %#v", err)
	require.True(t, areEqual, "public keys mismatched: got: %v, expected: %v", certPubKey, key)
}

func getSelfSigned(t *testing.T, subject, issuer *x509.Certificate, key *rsa.PrivateKey) (string, *x509.Certificate) {
	t.Helper()
	selfSigned, err := x509.CreateCertificate(rand.Reader, subject, issuer, key.Public(), key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(selfSigned)
	if err != nil {
		t.Fatal(err)
	}
	pemSS := strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: selfSigned,
	})))
	return pemSS, cert
}

// CRL related helpers
func getCrlCertificateList(t *testing.T, client *api.Client, mountPoint string) pkix.TBSCertificateList {
	path := fmt.Sprintf("/v1/%s/crl", mountPoint)
	return getParsedCrlAtPath(t, client, path).TBSCertList
}

func parseCrlPemBytes(t *testing.T, crlPem []byte) pkix.TBSCertificateList {
	certList, err := x509.ParseCRL(crlPem)
	require.NoError(t, err)
	return certList.TBSCertList
}

func requireSerialNumberInCRL(t *testing.T, revokeList pkix.TBSCertificateList, serialNum string) bool {
	serialsInList := make([]string, 0, len(revokeList.RevokedCertificates))
	for _, revokeEntry := range revokeList.RevokedCertificates {
		formattedSerial := certutil.GetHexFormatted(revokeEntry.SerialNumber.Bytes(), ":")
		serialsInList = append(serialsInList, formattedSerial)
		if formattedSerial == serialNum {
			return true
		}
	}

	if t != nil {
		t.Fatalf("the serial number %s, was not found in the CRL list containing: %v", serialNum, serialsInList)
	}

	return false
}

func getParsedCrl(t *testing.T, client *api.Client, mountPoint string) *pkix.CertificateList {
	path := fmt.Sprintf("/v1/%s/crl", mountPoint)
	return getParsedCrlAtPath(t, client, path)
}

func getParsedCrlForIssuer(t *testing.T, client *api.Client, mountPoint string, issuer string) *pkix.CertificateList {
	path := fmt.Sprintf("/v1/%v/issuer/%v/crl/der", mountPoint, issuer)
	crl := getParsedCrlAtPath(t, client, path)

	// Now fetch the issuer as well and verify the certificate
	path = fmt.Sprintf("/v1/%v/issuer/%v/der", mountPoint, issuer)
	req := client.NewRequest("GET", path)
	resp, err := client.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	certBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(certBytes) == 0 {
		t.Fatalf("expected certificate in response body")
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		t.Fatal(err)
	}
	if cert == nil {
		t.Fatalf("expected parsed certificate")
	}

	if err := cert.CheckCRLSignature(crl); err != nil {
		t.Fatalf("expected valid signature on CRL for issuer %v: %v", issuer, crl)
	}

	return crl
}

func getParsedCrlAtPath(t *testing.T, client *api.Client, path string) *pkix.CertificateList {
	req := client.NewRequest("GET", path)
	resp, err := client.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	crlBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(crlBytes) == 0 {
		t.Fatalf("expected CRL in response body")
	}

	crl, err := x509.ParseDERCRL(crlBytes)
	if err != nil {
		t.Fatal(err)
	}
	return crl
}

func getParsedCrlFromBackend(t *testing.T, b *backend, s logical.Storage, path string) *pkix.CertificateList {
	resp, err := CBRead(b, s, path)
	if err != nil {
		t.Fatal(err)
	}

	crl, err := x509.ParseDERCRL(resp.Data[logical.HTTPRawBody].([]byte))
	if err != nil {
		t.Fatal(err)
	}
	return crl
}

// Direct storage backend helpers (b, s := createBackendWithStorage(t)) which
// are mostly compatible with client.Logical() operations. The main difference
// is that the JSON round-tripping hasn't occurred, so values are as the
// backend returns them (e.g., []string instead of []interface{}).
func CBReq(b *backend, s logical.Storage, operation logical.Operation, path string, data map[string]interface{}) (*logical.Response, error) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  operation,
		Path:       path,
		Data:       data,
		Storage:    s,
		MountPoint: "pki/",
	})
	if err != nil || resp == nil {
		return resp, err
	}

	if msg, ok := resp.Data["error"]; ok && msg != nil && len(msg.(string)) > 0 {
		return resp, fmt.Errorf("%s", msg)
	}

	return resp, nil
}

func CBRead(b *backend, s logical.Storage, path string) (*logical.Response, error) {
	return CBReq(b, s, logical.ReadOperation, path, make(map[string]interface{}))
}

func CBWrite(b *backend, s logical.Storage, path string, data map[string]interface{}) (*logical.Response, error) {
	return CBReq(b, s, logical.UpdateOperation, path, data)
}

func CBList(b *backend, s logical.Storage, path string) (*logical.Response, error) {
	return CBReq(b, s, logical.ListOperation, path, make(map[string]interface{}))
}

func CBDelete(b *backend, s logical.Storage, path string) (*logical.Response, error) {
	return CBReq(b, s, logical.DeleteOperation, path, make(map[string]interface{}))
}

func CBPatch(b *backend, s logical.Storage, path string, data map[string]interface{}) (*logical.Response, error) {
	return CBReq(b, s, logical.PatchOperation, path, data)
}

func requireSuccessNonNilResponse(t *testing.T, resp *logical.Response, err error, msgAndArgs ...interface{}) {
	require.NoError(t, err, msgAndArgs...)
	require.False(t, resp.IsError(), msgAndArgs...)
	require.NotNil(t, resp, msgAndArgs...)
}

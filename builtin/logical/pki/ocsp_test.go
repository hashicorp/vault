package pki

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ocsp"
)

func TestOcsp_Disabled(t *testing.T) {
	t.Parallel()
	type testArgs struct {
		reqType string
	}
	var tests []testArgs
	for _, reqType := range []string{"get", "post"} {
		tests = append(tests, testArgs{
			reqType: reqType,
		})
	}
	for _, tt := range tests {
		localTT := tt
		t.Run(localTT.reqType, func(t *testing.T) {
			b, s, testEnv := setupOcspEnv(t, "rsa")
			resp, err := CBWrite(b, s, "config/crl", map[string]interface{}{
				"ocsp_disable": "true",
			})
			requireSuccessNilResponse(t, resp, err)
			resp, err = sendOcspRequest(t, b, s, localTT.reqType, testEnv.leafCertIssuer1, testEnv.issuer1, crypto.SHA1)
			require.NoError(t, err)
			requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
			require.Equal(t, 401, resp.Data["http_status_code"])
			require.Equal(t, "application/ocsp-response", resp.Data["http_content_type"])
			respDer := resp.Data["http_raw_body"].([]byte)

			require.Equal(t, ocsp.UnauthorizedErrorResponse, respDer)
		})
	}
}

// If we can't find the issuer within the request return an UnauthorizedErrorResponse according to the
// RFC, similar to if we are disabled (lack of authority)
func TestOcsp_UnknownIssuer(t *testing.T) {
	t.Parallel()

	_, _, testEnv := setupOcspEnv(t, "ec")
	// Create another completely empty mount so the created issuer/certificate above is unknown
	b, s := createBackendWithStorage(t)

	resp, err := sendOcspRequest(t, b, s, "get", testEnv.leafCertIssuer1, testEnv.issuer1, crypto.SHA1)
	require.NoError(t, err)
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 401, resp.Data["http_status_code"])
	require.Equal(t, "application/ocsp-response", resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	require.Equal(t, ocsp.UnauthorizedErrorResponse, respDer)
}

func TestOcsp_WrongIssuerInRequest(t *testing.T) {
	// If the issuers do exist, but the request coming in associates the serial with the
	// wrong issuer return an error.
	t.Parallel()

	b, s, testEnv := setupOcspEnv(t, "ec")
	serial := serialFromCert(testEnv.leafCertIssuer1)
	resp, err := CBWrite(b, s, "revoke", map[string]interface{}{
		"serial_number": serial,
	})
	requireSuccessNonNilResponse(t, resp, err, "revoke")

	resp, err = sendOcspRequest(t, b, s, "get", testEnv.leafCertIssuer1, testEnv.issuer2, crypto.SHA1)
	require.NoError(t, err)
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 401, resp.Data["http_status_code"])
	require.Equal(t, "application/ocsp-response", resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	require.Equal(t, ocsp.UnauthorizedErrorResponse, respDer)
}

func TestOcsp_MalformedRequests(t *testing.T) {
	t.Parallel()
	type testArgs struct {
		reqType string
	}
	var tests []testArgs
	for _, reqType := range []string{"get", "post"} {
		tests = append(tests, testArgs{
			reqType: reqType,
		})
	}
	for _, tt := range tests {
		localTT := tt
		t.Run(localTT.reqType, func(t *testing.T) {
			b, s, _ := setupOcspEnv(t, "rsa")
			badReq := []byte("this is a bad request")
			var resp *logical.Response
			var err error
			switch localTT.reqType {
			case "get":
				resp, err = sendOcspGetRequest(b, s, badReq)
			case "post":
				resp, err = sendOcspPostRequest(b, s, badReq)
			default:
				t.Fatalf("bad request type")
			}
			require.NoError(t, err)
			requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
			require.Equal(t, 400, resp.Data["http_status_code"])
			require.Equal(t, "application/ocsp-response", resp.Data["http_content_type"])
			respDer := resp.Data["http_raw_body"].([]byte)

			require.Equal(t, ocsp.MalformedRequestErrorResponse, respDer)
		})
	}
}

func TestOcsp_InvalidIssuerIdInRevocationEntry(t *testing.T) {
	// Validate that we properly handle a revocation entry that contains an issuer ID that no longer exists
	t.Parallel()

	b, s, testEnv := setupOcspEnv(t, "ec")
	ctx := context.Background()

	// Revoke the entry
	serial := serialFromCert(testEnv.leafCertIssuer1)
	resp, err := CBWrite(b, s, "revoke", map[string]interface{}{
		"serial_number": serial,
	})
	requireSuccessNonNilResponse(t, resp, err, "revoke")

	// Twiddle the entry so that the issuer id is no longer valid.
	storagePath := revokedPath + normalizeSerial(serial)
	var revInfo revocationInfo
	revEntry, err := s.Get(ctx, storagePath)
	require.NoError(t, err, "failed looking up storage path: %s", storagePath)
	err = revEntry.DecodeJSON(&revInfo)
	require.NoError(t, err, "failed decoding storage entry: %v", revEntry)
	revInfo.CertificateIssuer = "00000000-0000-0000-0000-000000000000"
	revEntry, err = logical.StorageEntryJSON(storagePath, revInfo)
	require.NoError(t, err, "failed re-encoding revocation info: %v", revInfo)
	err = s.Put(ctx, revEntry)
	require.NoError(t, err, "failed writing out new revocation entry: %v", revEntry)

	// Send the request
	resp, err = sendOcspRequest(t, b, s, "get", testEnv.leafCertIssuer1, testEnv.issuer1, crypto.SHA1)
	require.NoError(t, err)
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 401, resp.Data["http_status_code"])
	require.Equal(t, "application/ocsp-response", resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	require.Equal(t, ocsp.UnauthorizedErrorResponse, respDer)
}

func TestOcsp_ValidRequests(t *testing.T) {
	t.Parallel()
	type testArgs struct {
		reqType   string
		caKeyType string
		reqHash   crypto.Hash
	}
	var tests []testArgs
	for _, reqType := range []string{"get", "post"} {
		for _, caKeyType := range []string{"rsa", "ec"} { // "ed25519" is not supported at the moment in x/crypto/ocsp
			for _, requestHash := range []crypto.Hash{crypto.SHA1, crypto.SHA256} {
				tests = append(tests, testArgs{
					reqType:   reqType,
					caKeyType: caKeyType,
					reqHash:   requestHash,
				})
			}
		}
	}
	for _, tt := range tests {
		localTT := tt
		testName := fmt.Sprintf("%s-%s-%s", localTT.reqType, localTT.caKeyType, localTT.reqHash)
		t.Run(testName, func(t *testing.T) {
			runOcspRequestTest(t, localTT.reqType, localTT.caKeyType, localTT.reqHash)
		})
	}
}

func runOcspRequestTest(t *testing.T, requestType string, caKeyType string, requestHash crypto.Hash) {
	b, s, testEnv := setupOcspEnv(t, caKeyType)

	// Non-revoked cert
	resp, err := sendOcspRequest(t, b, s, requestType, testEnv.leafCertIssuer1, testEnv.issuer1, requestHash)
	requireSuccessNonNilResponse(t, resp, err, "ocsp get request")
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 200, resp.Data["http_status_code"])
	require.Equal(t, "application/ocsp-response", resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	ocspResp, err := ocsp.ParseResponse(respDer, testEnv.issuer1)
	require.NoError(t, err, "parsing ocsp get response")

	require.Equal(t, ocsp.Good, ocspResp.Status)
	require.Equal(t, requestHash, ocspResp.IssuerHash)
	require.Equal(t, testEnv.issuer1, ocspResp.Certificate)
	require.Equal(t, 0, ocspResp.RevocationReason)
	require.Equal(t, testEnv.leafCertIssuer1.SerialNumber, ocspResp.SerialNumber)

	requireOcspSignatureAlgoForKey(t, testEnv.issuer1.PublicKey, ocspResp.SignatureAlgorithm)
	requireOcspResponseSignedBy(t, ocspResp, testEnv.issuer1.PublicKey)

	// Now revoke it
	resp, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"serial_number": serialFromCert(testEnv.leafCertIssuer1),
	})
	requireSuccessNonNilResponse(t, resp, err, "revoke")

	resp, err = sendOcspRequest(t, b, s, requestType, testEnv.leafCertIssuer1, testEnv.issuer1, requestHash)
	requireSuccessNonNilResponse(t, resp, err, "ocsp get request with revoked")
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 200, resp.Data["http_status_code"])
	require.Equal(t, "application/ocsp-response", resp.Data["http_content_type"])
	respDer = resp.Data["http_raw_body"].([]byte)

	ocspResp, err = ocsp.ParseResponse(respDer, testEnv.issuer1)
	require.NoError(t, err, "parsing ocsp get response with revoked")

	require.Equal(t, ocsp.Revoked, ocspResp.Status)
	require.Equal(t, requestHash, ocspResp.IssuerHash)
	require.Equal(t, testEnv.issuer1, ocspResp.Certificate)
	require.Equal(t, 0, ocspResp.RevocationReason)
	require.Equal(t, testEnv.leafCertIssuer1.SerialNumber, ocspResp.SerialNumber)

	requireOcspSignatureAlgoForKey(t, testEnv.issuer1.PublicKey, ocspResp.SignatureAlgorithm)
	requireOcspResponseSignedBy(t, ocspResp, testEnv.issuer1.PublicKey)

	// Request status for our second issuer
	resp, err = sendOcspRequest(t, b, s, requestType, testEnv.leafCertIssuer2, testEnv.issuer2, requestHash)
	requireSuccessNonNilResponse(t, resp, err, "ocsp get request")
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 200, resp.Data["http_status_code"])
	require.Equal(t, "application/ocsp-response", resp.Data["http_content_type"])
	respDer = resp.Data["http_raw_body"].([]byte)

	ocspResp, err = ocsp.ParseResponse(respDer, testEnv.issuer2)
	require.NoError(t, err, "parsing ocsp get response")

	require.Equal(t, ocsp.Good, ocspResp.Status)
	require.Equal(t, requestHash, ocspResp.IssuerHash)
	require.Equal(t, testEnv.issuer2, ocspResp.Certificate)
	require.Equal(t, 0, ocspResp.RevocationReason)
	require.Equal(t, testEnv.leafCertIssuer2.SerialNumber, ocspResp.SerialNumber)

	requireOcspSignatureAlgoForKey(t, testEnv.issuer2.PublicKey, ocspResp.SignatureAlgorithm)
	requireOcspResponseSignedBy(t, ocspResp, testEnv.issuer2.PublicKey)
}

func requireOcspSignatureAlgoForKey(t *testing.T, key crypto.PublicKey, algorithm x509.SignatureAlgorithm) {
	switch key.(type) {
	case *rsa.PublicKey:
		require.Equal(t, x509.SHA256WithRSA, algorithm)
	case *ecdsa.PublicKey:
		require.Equal(t, x509.ECDSAWithSHA256, algorithm)
	case ed25519.PublicKey:
		require.Equal(t, x509.PureEd25519, algorithm)
	default:
		t.Fatalf("unsupported public key type %T", key)
	}
}

type ocspTestEnv struct {
	issuer1 *x509.Certificate
	issuer2 *x509.Certificate

	issuerId1 issuerID
	issuerId2 issuerID

	leafCertIssuer1 *x509.Certificate
	leafCertIssuer2 *x509.Certificate
}

func setupOcspEnv(t *testing.T, keyType string) (*backend, logical.Storage, *ocspTestEnv) {
	b, s := createBackendWithStorage(t)
	var issuerCerts []*x509.Certificate
	var leafCerts []*x509.Certificate
	var issuerIds []issuerID

	for i := 0; i < 2; i++ {
		resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
			"key_type":    keyType,
			"ttl":         "40h",
			"common_name": "example-ocsp.com",
		})
		requireSuccessNonNilResponse(t, resp, err, "root/generate/internal")
		requireFieldsSetInResp(t, resp, "issuer_id")
		issuerId := resp.Data["issuer_id"].(issuerID)

		resp, err = CBWrite(b, s, "roles/test"+strconv.FormatInt(int64(i), 10), map[string]interface{}{
			"allow_bare_domains": true,
			"allow_subdomains":   true,
			"allowed_domains":    "foobar.com",
			"no_store":           false,
			"generate_lease":     false,
			"issuer_ref":         issuerId,
			"key_type":           keyType,
		})
		requireSuccessNilResponse(t, resp, err, "roles/test"+strconv.FormatInt(int64(i), 10))

		resp, err = CBWrite(b, s, "issue/test"+strconv.FormatInt(int64(i), 10), map[string]interface{}{
			"common_name": "test.foobar.com",
		})
		requireSuccessNonNilResponse(t, resp, err, "roles/test"+strconv.FormatInt(int64(i), 10))
		requireFieldsSetInResp(t, resp, "certificate", "issuing_ca", "serial_number")
		leafCert := parseCert(t, resp.Data["certificate"].(string))
		issuingCa := parseCert(t, resp.Data["issuing_ca"].(string))

		issuerCerts = append(issuerCerts, issuingCa)
		leafCerts = append(leafCerts, leafCert)
		issuerIds = append(issuerIds, issuerId)
	}

	testEnv := &ocspTestEnv{
		issuerId1:       issuerIds[0],
		issuer1:         issuerCerts[0],
		leafCertIssuer1: leafCerts[0],

		issuerId2:       issuerIds[1],
		issuer2:         issuerCerts[1],
		leafCertIssuer2: leafCerts[1],
	}

	return b, s, testEnv
}

func sendOcspRequest(t *testing.T, b *backend, s logical.Storage, getOrPost string, cert, issuer *x509.Certificate, requestHash crypto.Hash) (*logical.Response, error) {
	ocspRequest := generateRequest(t, requestHash, cert, issuer)

	switch strings.ToLower(getOrPost) {
	case "get":
		return sendOcspGetRequest(b, s, ocspRequest)
	case "post":
		return sendOcspPostRequest(b, s, ocspRequest)
	default:
		t.Fatalf("unsupported value for sendOcspRequest getOrPost arg: %s", getOrPost)
	}
	return nil, nil
}

func sendOcspGetRequest(b *backend, s logical.Storage, ocspRequest []byte) (*logical.Response, error) {
	urlEncoded := base64.StdEncoding.EncodeToString(ocspRequest)
	return CBRead(b, s, "ocsp/"+urlEncoded)
}

func sendOcspPostRequest(b *backend, s logical.Storage, ocspRequest []byte) (*logical.Response, error) {
	reader := io.NopCloser(bytes.NewReader(ocspRequest))
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation:  logical.UpdateOperation,
		Path:       "ocsp",
		Storage:    s,
		MountPoint: "pki/",
		HTTPRequest: &http.Request{
			Body: reader,
		},
	})

	return resp, err
}

func generateRequest(t *testing.T, requestHash crypto.Hash, cert *x509.Certificate, issuer *x509.Certificate) []byte {
	opts := &ocsp.RequestOptions{Hash: requestHash}
	ocspRequestDer, err := ocsp.CreateRequest(cert, issuer, opts)
	require.NoError(t, err, "Failed generating OCSP request")
	return ocspRequestDer
}

func requireOcspResponseSignedBy(t *testing.T, ocspResp *ocsp.Response, key crypto.PublicKey) {
	require.Contains(t, []x509.SignatureAlgorithm{x509.SHA256WithRSA, x509.ECDSAWithSHA256}, ocspResp.SignatureAlgorithm)

	hasher := sha256.New()
	hashAlgo := crypto.SHA256
	hasher.Write(ocspResp.TBSResponseData)
	hashData := hasher.Sum(nil)

	switch key.(type) {
	case *rsa.PublicKey:
		err := rsa.VerifyPKCS1v15(key.(*rsa.PublicKey), hashAlgo, hashData, ocspResp.Signature)
		require.NoError(t, err, "the ocsp response was not signed by the expected public rsa key.")
	case *ecdsa.PublicKey:
		verify := ecdsa.VerifyASN1(key.(*ecdsa.PublicKey), hashData, ocspResp.Signature)
		require.True(t, verify, "the certificate was not signed by the expected public ecdsa key.")
	}
}

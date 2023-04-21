package pki

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ocsp"
)

// If the ocsp_disabled flag is set to true in the crl configuration make sure we always
// return an Unauthorized error back as we assume an end-user disabling the feature does
// not want us to act as the OCSP authority and the RFC specifies this is the appropriate response.
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
			requireSuccessNonNilResponse(t, resp, err)
			resp, err = SendOcspRequest(t, b, s, localTT.reqType, testEnv.leafCertIssuer1, testEnv.issuer1, crypto.SHA1)
			require.NoError(t, err)
			requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
			require.Equal(t, 401, resp.Data["http_status_code"])
			require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
			respDer := resp.Data["http_raw_body"].([]byte)

			require.Equal(t, ocsp.UnauthorizedErrorResponse, respDer)
		})
	}
}

// If we can't find the issuer within the request and have no default issuer to sign an Unknown response
// with return an UnauthorizedErrorResponse/according to/the RFC, similar to if we are disabled (lack of authority)
// This behavior differs from CRLs when an issuer is removed from a mount.
func TestOcsp_UnknownIssuerWithNoDefault(t *testing.T) {
	t.Parallel()

	_, _, testEnv := setupOcspEnv(t, "ec")
	// Create another completely empty mount so the created issuer/certificate above is unknown
	b, s := CreateBackendWithStorage(t)

	resp, err := SendOcspRequest(t, b, s, "get", testEnv.leafCertIssuer1, testEnv.issuer1, crypto.SHA1)
	require.NoError(t, err)
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 401, resp.Data["http_status_code"])
	require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	require.Equal(t, ocsp.UnauthorizedErrorResponse, respDer)
}

// If the issuer in the request does exist, but the request coming in associates the serial with the
// wrong issuer return an Unknown response back to the caller.
func TestOcsp_WrongIssuerInRequest(t *testing.T) {
	t.Parallel()

	b, s, testEnv := setupOcspEnv(t, "ec")
	serial := serialFromCert(testEnv.leafCertIssuer1)
	resp, err := CBWrite(b, s, "revoke", map[string]interface{}{
		"serial_number": serial,
	})
	requireSuccessNonNilResponse(t, resp, err, "revoke")

	resp, err = SendOcspRequest(t, b, s, "get", testEnv.leafCertIssuer1, testEnv.issuer2, crypto.SHA1)
	require.NoError(t, err)
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 200, resp.Data["http_status_code"])
	require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	ocspResp, err := ocsp.ParseResponse(respDer, testEnv.issuer1)
	require.NoError(t, err, "parsing ocsp get response")

	require.Equal(t, ocsp.Unknown, ocspResp.Status)
}

// Verify that requests we can't properly decode result in the correct response of MalformedRequestError
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
			require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
			respDer := resp.Data["http_raw_body"].([]byte)

			require.Equal(t, ocsp.MalformedRequestErrorResponse, respDer)
		})
	}
}

// Validate that we properly handle a revocation entry that contains an issuer ID that no longer exists,
// the best we can do in this use case is to respond back with the default issuer that we don't know
// the issuer that they are requesting (we can't guarantee that the client is actually requesting a serial
// from that issuer)
func TestOcsp_InvalidIssuerIdInRevocationEntry(t *testing.T) {
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
	resp, err = SendOcspRequest(t, b, s, "get", testEnv.leafCertIssuer1, testEnv.issuer1, crypto.SHA1)
	require.NoError(t, err)
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 200, resp.Data["http_status_code"])
	require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	ocspResp, err := ocsp.ParseResponse(respDer, testEnv.issuer1)
	require.NoError(t, err, "parsing ocsp get response")

	require.Equal(t, ocsp.Unknown, ocspResp.Status)
}

// Validate that we properly handle an unknown issuer use-case but that the default issuer
// does not have the OCSP usage flag set, we can't do much else other than reply with an
// Unauthorized response.
func TestOcsp_UnknownIssuerIdWithDefaultHavingOcspUsageRemoved(t *testing.T) {
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

	// Update our issuers to no longer have the OcspSigning usage
	resp, err = CBPatch(b, s, "issuer/"+testEnv.issuerId1.String(), map[string]interface{}{
		"usage": "read-only,issuing-certificates,crl-signing",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed resetting usage flags on issuer1")
	resp, err = CBPatch(b, s, "issuer/"+testEnv.issuerId2.String(), map[string]interface{}{
		"usage": "read-only,issuing-certificates,crl-signing",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed resetting usage flags on issuer2")

	// Send the request
	resp, err = SendOcspRequest(t, b, s, "get", testEnv.leafCertIssuer1, testEnv.issuer1, crypto.SHA1)
	require.NoError(t, err)
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 401, resp.Data["http_status_code"])
	require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	require.Equal(t, ocsp.UnauthorizedErrorResponse, respDer)
}

// Verify that if we do have a revoked certificate entry for the request, that matches an
// issuer but that issuer does not have the OcspUsage flag set that we return an Unauthorized
// response back to the caller
func TestOcsp_RevokedCertHasIssuerWithoutOcspUsage(t *testing.T) {
	b, s, testEnv := setupOcspEnv(t, "ec")

	// Revoke our certificate
	resp, err := CBWrite(b, s, "revoke", map[string]interface{}{
		"serial_number": serialFromCert(testEnv.leafCertIssuer1),
	})
	requireSuccessNonNilResponse(t, resp, err, "revoke")

	// Update our issuer to no longer have the OcspSigning usage
	resp, err = CBPatch(b, s, "issuer/"+testEnv.issuerId1.String(), map[string]interface{}{
		"usage": "read-only,issuing-certificates,crl-signing",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed resetting usage flags on issuer")
	requireFieldsSetInResp(t, resp, "usage")

	// Do not assume a specific ordering for usage...
	usages, err := NewIssuerUsageFromNames(strings.Split(resp.Data["usage"].(string), ","))
	require.NoError(t, err, "failed parsing usage return value")
	require.True(t, usages.HasUsage(IssuanceUsage))
	require.True(t, usages.HasUsage(CRLSigningUsage))
	require.False(t, usages.HasUsage(OCSPSigningUsage))

	// Request an OCSP request from it, we should get an Unauthorized response back
	resp, err = SendOcspRequest(t, b, s, "get", testEnv.leafCertIssuer1, testEnv.issuer1, crypto.SHA1)
	requireSuccessNonNilResponse(t, resp, err, "ocsp get request")
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 401, resp.Data["http_status_code"])
	require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	require.Equal(t, ocsp.UnauthorizedErrorResponse, respDer)
}

// Verify if our matching issuer for a revocation entry has no key associated with it that
// we bail with an Unauthorized response.
func TestOcsp_RevokedCertHasIssuerWithoutAKey(t *testing.T) {
	b, s, testEnv := setupOcspEnv(t, "ec")

	// Revoke our certificate
	resp, err := CBWrite(b, s, "revoke", map[string]interface{}{
		"serial_number": serialFromCert(testEnv.leafCertIssuer1),
	})
	requireSuccessNonNilResponse(t, resp, err, "revoke")

	// Delete the key associated with our issuer
	resp, err = CBRead(b, s, "issuer/"+testEnv.issuerId1.String())
	requireSuccessNonNilResponse(t, resp, err, "failed reading issuer")
	requireFieldsSetInResp(t, resp, "key_id")
	keyId := resp.Data["key_id"].(keyID)

	// This is a bit naughty but allow me to delete the key...
	sc := b.makeStorageContext(context.Background(), s)
	issuer, err := sc.fetchIssuerById(testEnv.issuerId1)
	require.NoError(t, err, "failed to get issuer from storage")
	issuer.KeyID = ""
	err = sc.writeIssuer(issuer)
	require.NoError(t, err, "failed to write issuer update")

	resp, err = CBDelete(b, s, "key/"+keyId.String())
	requireSuccessNonNilResponse(t, resp, err, "failed deleting key")

	// Request an OCSP request from it, we should get an Unauthorized response back
	resp, err = SendOcspRequest(t, b, s, "get", testEnv.leafCertIssuer1, testEnv.issuer1, crypto.SHA1)
	requireSuccessNonNilResponse(t, resp, err, "ocsp get request")
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 401, resp.Data["http_status_code"])
	require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	require.Equal(t, ocsp.UnauthorizedErrorResponse, respDer)
}

// Verify if for some reason an end-user has rotated an existing certificate using the same
// key so our algo matches multiple issuers and one has OCSP usage disabled. We expect that
// even if a prior issuer issued the certificate, the new matching issuer can respond and sign
// the response to the caller on its behalf.
//
// NOTE: This test is a bit at the mercy of iteration order of the issuer ids.
//
//	If it becomes flaky, most likely something is wrong in the code
//	and not the test.
func TestOcsp_MultipleMatchingIssuersOneWithoutSigningUsage(t *testing.T) {
	b, s, testEnv := setupOcspEnv(t, "ec")

	// Create a matching issuer as issuer1 with the same backing key
	resp, err := CBWrite(b, s, "root/rotate/existing", map[string]interface{}{
		"key_ref":     testEnv.keyId1,
		"ttl":         "40h",
		"common_name": "example-ocsp.com",
	})
	requireSuccessNonNilResponse(t, resp, err, "rotate issuer failed")
	requireFieldsSetInResp(t, resp, "issuer_id")
	rotatedCert := parseCert(t, resp.Data["certificate"].(string))

	// Remove ocsp signing from our issuer
	resp, err = CBPatch(b, s, "issuer/"+testEnv.issuerId1.String(), map[string]interface{}{
		"usage": "read-only,issuing-certificates,crl-signing",
	})
	requireSuccessNonNilResponse(t, resp, err, "failed resetting usage flags on issuer")
	requireFieldsSetInResp(t, resp, "usage")
	// Do not assume a specific ordering for usage...
	usages, err := NewIssuerUsageFromNames(strings.Split(resp.Data["usage"].(string), ","))
	require.NoError(t, err, "failed parsing usage return value")
	require.True(t, usages.HasUsage(IssuanceUsage))
	require.True(t, usages.HasUsage(CRLSigningUsage))
	require.False(t, usages.HasUsage(OCSPSigningUsage))

	// Request an OCSP request from it, we should get a Good response back, from the rotated cert
	resp, err = SendOcspRequest(t, b, s, "get", testEnv.leafCertIssuer1, testEnv.issuer1, crypto.SHA1)
	requireSuccessNonNilResponse(t, resp, err, "ocsp get request")
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 200, resp.Data["http_status_code"])
	require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	ocspResp, err := ocsp.ParseResponse(respDer, testEnv.issuer1)
	require.NoError(t, err, "parsing ocsp get response")

	require.Equal(t, ocsp.Good, ocspResp.Status)
	require.Equal(t, crypto.SHA1, ocspResp.IssuerHash)
	require.Equal(t, 0, ocspResp.RevocationReason)
	require.Equal(t, testEnv.leafCertIssuer1.SerialNumber, ocspResp.SerialNumber)

	requireOcspSignatureAlgoForKey(t, rotatedCert.SignatureAlgorithm, ocspResp.SignatureAlgorithm)
	requireOcspResponseSignedBy(t, ocspResp, rotatedCert)
}

// Make sure OCSP GET/POST requests work through the entire stack, and not just
// through the quicker backend layer the other tests are doing.
func TestOcsp_HigherLevel(t *testing.T) {
	t.Parallel()
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client
	mountPKIEndpoint(t, client, "pki")
	resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"key_type":    "ec",
		"common_name": "root-ca.com",
		"ttl":         "600h",
	})

	require.NoError(t, err, "error generating root ca: %v", err)
	require.NotNil(t, resp, "expected ca info from root")

	issuerCert := parseCert(t, resp.Data["certificate"].(string))

	resp, err = client.Logical().Write("pki/roles/example", map[string]interface{}{
		"allowed_domains":  "example.com",
		"allow_subdomains": "true",
		"no_store":         "false", // make sure we store this cert
		"max_ttl":          "1h",
		"key_type":         "ec",
	})
	require.NoError(t, err, "error setting up pki role: %v", err)

	resp, err = client.Logical().Write("pki/issue/example", map[string]interface{}{
		"common_name": "test.example.com",
		"ttl":         "15m",
	})
	require.NoError(t, err, "error issuing certificate: %v", err)
	require.NotNil(t, resp, "got nil response from issuing request")
	certToRevoke := parseCert(t, resp.Data["certificate"].(string))
	serialNum := resp.Data["serial_number"].(string)

	// Revoke the certificate
	resp, err = client.Logical().Write("pki/revoke", map[string]interface{}{
		"serial_number": serialNum,
	})
	require.NoError(t, err, "error revoking certificate: %v", err)
	require.NotNil(t, resp, "got nil response from revoke")

	// Make sure that OCSP handler responds properly
	ocspReq := generateRequest(t, crypto.SHA256, certToRevoke, issuerCert)
	ocspPostReq := client.NewRequest(http.MethodPost, "/v1/pki/ocsp")
	ocspPostReq.Headers.Set("Content-Type", "application/ocsp-request")
	ocspPostReq.BodyBytes = ocspReq
	rawResp, err := client.RawRequest(ocspPostReq)
	require.NoError(t, err, "failed sending ocsp post request")

	require.Equal(t, 200, rawResp.StatusCode)
	require.Equal(t, ocspResponseContentType, rawResp.Header.Get("Content-Type"))
	bodyReader := rawResp.Body
	respDer, err := io.ReadAll(bodyReader)
	bodyReader.Close()
	require.NoError(t, err, "failed reading response body")

	ocspResp, err := ocsp.ParseResponse(respDer, issuerCert)
	require.NoError(t, err, "parsing ocsp get response")

	require.Equal(t, ocsp.Revoked, ocspResp.Status)
	require.Equal(t, certToRevoke.SerialNumber, ocspResp.SerialNumber)

	// Test OCSP Get request for ocsp
	urlEncoded := base64.StdEncoding.EncodeToString(ocspReq)
	ocspGetReq := client.NewRequest(http.MethodGet, "/v1/pki/ocsp/"+urlEncoded)
	ocspGetReq.Headers.Set("Content-Type", "application/ocsp-request")
	rawResp, err = client.RawRequest(ocspGetReq)
	require.NoError(t, err, "failed sending ocsp get request")

	require.Equal(t, 200, rawResp.StatusCode)
	require.Equal(t, ocspResponseContentType, rawResp.Header.Get("Content-Type"))
	bodyReader = rawResp.Body
	respDer, err = io.ReadAll(bodyReader)
	bodyReader.Close()
	require.NoError(t, err, "failed reading response body")

	ocspResp, err = ocsp.ParseResponse(respDer, issuerCert)
	require.NoError(t, err, "parsing ocsp get response")

	require.Equal(t, ocsp.Revoked, ocspResp.Status)
	require.Equal(t, certToRevoke.SerialNumber, ocspResp.SerialNumber)
}

func TestOcsp_ValidRequests(t *testing.T) {
	type caKeyConf struct {
		keyType string
		keyBits int
		sigBits int
	}
	t.Parallel()
	type testArgs struct {
		reqType string
		keyConf caKeyConf
		reqHash crypto.Hash
	}
	var tests []testArgs
	for _, reqType := range []string{"get", "post"} {
		for _, keyConf := range []caKeyConf{
			{"rsa", 0, 0},
			{"rsa", 0, 384},
			{"rsa", 0, 512},
			{"ec", 0, 0},
			{"ec", 521, 0},
		} {
			// "ed25519" is not supported at the moment in x/crypto/ocsp
			for _, requestHash := range []crypto.Hash{crypto.SHA1, crypto.SHA256, crypto.SHA384, crypto.SHA512} {
				tests = append(tests, testArgs{
					reqType: reqType,
					keyConf: keyConf,
					reqHash: requestHash,
				})
			}
		}
	}
	for _, tt := range tests {
		localTT := tt
		testName := fmt.Sprintf("%s-%s-keybits-%d-sigbits-%d-reqHash-%s", localTT.reqType, localTT.keyConf.keyType,
			localTT.keyConf.keyBits,
			localTT.keyConf.sigBits,
			localTT.reqHash)
		t.Run(testName, func(t *testing.T) {
			runOcspRequestTest(t, localTT.reqType, localTT.keyConf.keyType, localTT.keyConf.keyBits,
				localTT.keyConf.sigBits, localTT.reqHash)
		})
	}
}

func runOcspRequestTest(t *testing.T, requestType string, caKeyType string, caKeyBits int, caKeySigBits int, requestHash crypto.Hash) {
	b, s, testEnv := setupOcspEnvWithCaKeyConfig(t, caKeyType, caKeyBits, caKeySigBits)

	// Non-revoked cert
	resp, err := SendOcspRequest(t, b, s, requestType, testEnv.leafCertIssuer1, testEnv.issuer1, requestHash)
	requireSuccessNonNilResponse(t, resp, err, "ocsp get request")
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 200, resp.Data["http_status_code"])
	require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
	respDer := resp.Data["http_raw_body"].([]byte)

	ocspResp, err := ocsp.ParseResponse(respDer, testEnv.issuer1)
	require.NoError(t, err, "parsing ocsp get response")

	require.Equal(t, ocsp.Good, ocspResp.Status)
	require.Equal(t, requestHash, ocspResp.IssuerHash)
	require.Equal(t, 0, ocspResp.RevocationReason)
	require.Equal(t, testEnv.leafCertIssuer1.SerialNumber, ocspResp.SerialNumber)

	requireOcspSignatureAlgoForKey(t, testEnv.issuer1.SignatureAlgorithm, ocspResp.SignatureAlgorithm)
	requireOcspResponseSignedBy(t, ocspResp, testEnv.issuer1)

	// Now revoke it
	resp, err = CBWrite(b, s, "revoke", map[string]interface{}{
		"serial_number": serialFromCert(testEnv.leafCertIssuer1),
	})
	requireSuccessNonNilResponse(t, resp, err, "revoke")

	resp, err = SendOcspRequest(t, b, s, requestType, testEnv.leafCertIssuer1, testEnv.issuer1, requestHash)
	requireSuccessNonNilResponse(t, resp, err, "ocsp get request with revoked")
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 200, resp.Data["http_status_code"])
	require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
	respDer = resp.Data["http_raw_body"].([]byte)

	ocspResp, err = ocsp.ParseResponse(respDer, testEnv.issuer1)
	require.NoError(t, err, "parsing ocsp get response with revoked")

	require.Equal(t, ocsp.Revoked, ocspResp.Status)
	require.Equal(t, requestHash, ocspResp.IssuerHash)
	require.Equal(t, 0, ocspResp.RevocationReason)
	require.Equal(t, testEnv.leafCertIssuer1.SerialNumber, ocspResp.SerialNumber)

	requireOcspSignatureAlgoForKey(t, testEnv.issuer1.SignatureAlgorithm, ocspResp.SignatureAlgorithm)
	requireOcspResponseSignedBy(t, ocspResp, testEnv.issuer1)

	// Request status for our second issuer
	resp, err = SendOcspRequest(t, b, s, requestType, testEnv.leafCertIssuer2, testEnv.issuer2, requestHash)
	requireSuccessNonNilResponse(t, resp, err, "ocsp get request")
	requireFieldsSetInResp(t, resp, "http_content_type", "http_status_code", "http_raw_body")
	require.Equal(t, 200, resp.Data["http_status_code"])
	require.Equal(t, ocspResponseContentType, resp.Data["http_content_type"])
	respDer = resp.Data["http_raw_body"].([]byte)

	ocspResp, err = ocsp.ParseResponse(respDer, testEnv.issuer2)
	require.NoError(t, err, "parsing ocsp get response")

	require.Equal(t, ocsp.Good, ocspResp.Status)
	require.Equal(t, requestHash, ocspResp.IssuerHash)
	require.Equal(t, 0, ocspResp.RevocationReason)
	require.Equal(t, testEnv.leafCertIssuer2.SerialNumber, ocspResp.SerialNumber)

	// Verify that our thisUpdate and nextUpdate fields are updated as expected
	thisUpdate := ocspResp.ThisUpdate
	nextUpdate := ocspResp.NextUpdate
	require.True(t, thisUpdate.Before(nextUpdate),
		fmt.Sprintf("thisUpdate %s, should have been before nextUpdate: %s", thisUpdate, nextUpdate))
	nextUpdateDiff := nextUpdate.Sub(thisUpdate)
	expectedDiff, err := time.ParseDuration(defaultCrlConfig.OcspExpiry)
	require.NoError(t, err, "failed to parse default ocsp expiry value")
	require.Equal(t, expectedDiff, nextUpdateDiff,
		fmt.Sprintf("the delta between thisUpdate %s and nextUpdate: %s should have been around: %s but was %s",
			thisUpdate, nextUpdate, defaultCrlConfig.OcspExpiry, nextUpdateDiff))

	requireOcspSignatureAlgoForKey(t, testEnv.issuer2.SignatureAlgorithm, ocspResp.SignatureAlgorithm)
	requireOcspResponseSignedBy(t, ocspResp, testEnv.issuer2)
}

func requireOcspSignatureAlgoForKey(t *testing.T, expected x509.SignatureAlgorithm, actual x509.SignatureAlgorithm) {
	t.Helper()

	require.Equal(t, expected.String(), actual.String())
}

type ocspTestEnv struct {
	issuer1 *x509.Certificate
	issuer2 *x509.Certificate

	issuerId1 issuerID
	issuerId2 issuerID

	leafCertIssuer1 *x509.Certificate
	leafCertIssuer2 *x509.Certificate

	keyId1 keyID
	keyId2 keyID
}

func setupOcspEnv(t *testing.T, keyType string) (*backend, logical.Storage, *ocspTestEnv) {
	return setupOcspEnvWithCaKeyConfig(t, keyType, 0, 0)
}

func setupOcspEnvWithCaKeyConfig(t *testing.T, keyType string, caKeyBits int, caKeySigBits int) (*backend, logical.Storage, *ocspTestEnv) {
	b, s := CreateBackendWithStorage(t)
	var issuerCerts []*x509.Certificate
	var leafCerts []*x509.Certificate
	var issuerIds []issuerID
	var keyIds []keyID

	for i := 0; i < 2; i++ {
		resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
			"key_type":       keyType,
			"key_bits":       caKeyBits,
			"signature_bits": caKeySigBits,
			"ttl":            "40h",
			"common_name":    "example-ocsp.com",
		})
		requireSuccessNonNilResponse(t, resp, err, "root/generate/internal")
		requireFieldsSetInResp(t, resp, "issuer_id", "key_id")
		issuerId := resp.Data["issuer_id"].(issuerID)
		keyId := resp.Data["key_id"].(keyID)

		resp, err = CBWrite(b, s, "roles/test"+strconv.FormatInt(int64(i), 10), map[string]interface{}{
			"allow_bare_domains": true,
			"allow_subdomains":   true,
			"allowed_domains":    "foobar.com",
			"no_store":           false,
			"generate_lease":     false,
			"issuer_ref":         issuerId,
			"key_type":           keyType,
		})
		requireSuccessNonNilResponse(t, resp, err, "roles/test"+strconv.FormatInt(int64(i), 10))

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
		keyIds = append(keyIds, keyId)
	}

	testEnv := &ocspTestEnv{
		issuerId1:       issuerIds[0],
		issuer1:         issuerCerts[0],
		leafCertIssuer1: leafCerts[0],
		keyId1:          keyIds[0],

		issuerId2:       issuerIds[1],
		issuer2:         issuerCerts[1],
		leafCertIssuer2: leafCerts[1],
		keyId2:          keyIds[1],
	}

	return b, s, testEnv
}

func SendOcspRequest(t *testing.T, b *backend, s logical.Storage, getOrPost string, cert, issuer *x509.Certificate, requestHash crypto.Hash) (*logical.Response, error) {
	t.Helper()

	ocspRequest := generateRequest(t, requestHash, cert, issuer)

	switch strings.ToLower(getOrPost) {
	case "get":
		return sendOcspGetRequest(b, s, ocspRequest)
	case "post":
		return sendOcspPostRequest(b, s, ocspRequest)
	default:
		t.Fatalf("unsupported value for SendOcspRequest getOrPost arg: %s", getOrPost)
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
	t.Helper()

	opts := &ocsp.RequestOptions{Hash: requestHash}
	ocspRequestDer, err := ocsp.CreateRequest(cert, issuer, opts)
	require.NoError(t, err, "Failed generating OCSP request")
	return ocspRequestDer
}

func requireOcspResponseSignedBy(t *testing.T, ocspResp *ocsp.Response, issuer *x509.Certificate) {
	t.Helper()

	err := ocspResp.CheckSignatureFrom(issuer)
	require.NoError(t, err, "Failed signature verification of ocsp response: %w", err)
}

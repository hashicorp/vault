package pki

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestResignCrls_ForbidSigningOtherIssuerCRL(t *testing.T) {
	t.Parallel()

	// Some random CRL from another issuer
	pem1 := "-----BEGIN X509 CRL-----\nMIIBvjCBpwIBATANBgkqhkiG9w0BAQsFADAbMRkwFwYDVQQDExByb290LWV4YW1w\nbGUuY29tFw0yMjEwMjYyMTI5MzlaFw0yMjEwMjkyMTI5MzlaMCcwJQIUSnVf8wsd\nHjOt9drCYFhWxS9QqGoXDTIyMTAyNjIxMjkzOVqgLzAtMB8GA1UdIwQYMBaAFHki\nZ0XDUQVSajNRGXrg66OaIFlYMAoGA1UdFAQDAgEDMA0GCSqGSIb3DQEBCwUAA4IB\nAQBGIdtqTwemnLZF5AoP+jzvKZ26S3y7qvRIzd7f4A0EawzYmWXSXfwqo4TQ4DG3\nnvT+AaA1zCCOlH/1U+ufN9gSSN0j9ax58brSYMnMskMCqhLKIp0qnvS4jr/gopmF\nv8grbvLHEqNYTu1T7umMLdNQUsWT3Qc+EIjfoKj8xD2FHsZwJ+EMbytwl8Unipjr\nhz4rmcES/65vavfdFpOI6YXfi+UAaHBdkTqmHgg4BdpuXfYtlf+iotFSOkygD5fl\n0D+RVFW9uJv2WfbQ7kRt1X/VcFk/onw0AQqxZRVUzvjoMw+EMcxSq3UKOlXcWDxm\nEFz9rFQQ66L388EP8RD7Dh3X\n-----END X509 CRL-----"

	b, s := createBackendWithStorage(t)
	resp, err := CBWrite(b, s, "root/generate/internal", map[string]interface{}{
		"common_name": "test.com",
		"key_type":    "ec",
	})
	requireSuccessNonNilResponse(t, resp, err)

	resp, err = CBWrite(b, s, "issuer/default/resign-crls", map[string]interface{}{
		"crl_number":  "2",
		"next_update": "1h",
		"format":      "pem",
		"crls":        []string{pem1},
	})
	require.ErrorContains(t, err, "was not signed by requested issuer")
}

func TestResignCrls_NormalCrl(t *testing.T) {
	t.Parallel()
	b1, s1 := createBackendWithStorage(t)
	b2, s2 := createBackendWithStorage(t)

	// Setup two backends, with the same key material/certificate with a different leaf in each that is revoked.
	caCert, serial1, serial2, crl1, crl2 := setupResignCrlMounts(t, b1, s1, b2, s2)

	// Attempt to combine the CRLs
	resp, err := CBWrite(b1, s1, "issuer/default/resign-crls", map[string]interface{}{
		"crl_number":  "2",
		"next_update": "1h",
		"format":      "pem",
		"crls":        []string{crl1, crl2},
	})
	requireSuccessNonNilResponse(t, resp, err)
	requireFieldsSetInResp(t, resp, "crl")
	pemCrl := resp.Data["crl"].(string)
	combinedCrl, err := decodePemCrl(pemCrl)
	require.NoError(t, err, "failed decoding combined CRL")
	serials := extractSerialsFromCrl(t, combinedCrl)

	require.Contains(t, serials, serial1)
	require.Contains(t, serials, serial2)
	require.Equal(t, 2, len(serials), "serials contained more serials than expected")

	require.Equal(t, big.NewInt(int64(2)), combinedCrl.Number)
	require.Equal(t, combinedCrl.ThisUpdate.Add(1*time.Hour), combinedCrl.NextUpdate)

	extensions := combinedCrl.Extensions
	requireExtensionOid(t, []int{2, 5, 29, 20}, extensions) // CRL Number Extension
	requireExtensionOid(t, []int{2, 5, 29, 35}, extensions) // akidOid
	require.Equal(t, 2, len(extensions))

	err = combinedCrl.CheckSignatureFrom(caCert)
	require.NoError(t, err, "failed signature check of CRL")
}

func TestResignCrls_EliminateDuplicates(t *testing.T) {
	t.Parallel()
	b1, s1 := createBackendWithStorage(t)
	b2, s2 := createBackendWithStorage(t)

	// Setup two backends, with the same key material/certificate with a different leaf in each that is revoked.
	_, serial1, _, crl1, _ := setupResignCrlMounts(t, b1, s1, b2, s2)

	// Pass in the same CRLs to make sure we do not duplicate entries
	resp, err := CBWrite(b1, s1, "issuer/default/resign-crls", map[string]interface{}{
		"crl_number":  "2",
		"next_update": "1h",
		"format":      "pem",
		"crls":        []string{crl1, crl1},
	})
	requireSuccessNonNilResponse(t, resp, err)
	requireFieldsSetInResp(t, resp, "crl")
	pemCrl := resp.Data["crl"].(string)
	combinedCrl, err := decodePemCrl(pemCrl)
	require.NoError(t, err, "failed decoding combined CRL")

	// Technically this will die if we have duplicates.
	serials := extractSerialsFromCrl(t, combinedCrl)

	// We should have no warnings about collisions if they have the same revoked time
	require.Empty(t, resp.Warnings, "expected no warnings in response")

	require.Contains(t, serials, serial1)
	require.Equal(t, 1, len(serials), "serials contained more serials than expected")
}

func TestResignCrls_ConflictingExpiry(t *testing.T) {
	t.Parallel()
	b1, s1 := createBackendWithStorage(t)
	b2, s2 := createBackendWithStorage(t)

	// Setup two backends, with the same key material/certificate with a different leaf in each that is revoked.
	_, serial1, serial2, crl1, _ := setupResignCrlMounts(t, b1, s1, b2, s2)

	timeAfterMountSetup := time.Now()

	// Read in serial1 from mount 1
	resp, err := CBRead(b1, s1, "cert/"+serial1)
	requireSuccessNonNilResponse(t, resp, err, "failed reading serial 1's certificate")
	requireFieldsSetInResp(t, resp, "certificate")
	cert1 := resp.Data["certificate"].(string)

	// Wait until at least we have rolled over to the next second to match sure the generated CRL time
	// on backend 2 for the serial 1 will be different
	for {
		if time.Now().After(timeAfterMountSetup.Add(1 * time.Second)) {
			break
		}
	}

	// Use BYOC to revoke the same certificate on backend 2 now
	resp, err = CBWrite(b2, s2, "revoke", map[string]interface{}{
		"certificate": cert1,
	})
	requireSuccessNonNilResponse(t, resp, err, "failed revoking serial 1 on backend 2")

	// Fetch the new CRL from backend2 now
	resp, err = CBRead(b2, s2, "cert/crl")
	requireSuccessNonNilResponse(t, resp, err, "error fetch crl from backend 2")
	requireFieldsSetInResp(t, resp, "certificate")
	crl2 := resp.Data["certificate"].(string)

	// Attempt to combine the CRLs
	resp, err = CBWrite(b1, s1, "issuer/default/resign-crls", map[string]interface{}{
		"crl_number":  "2",
		"next_update": "1h",
		"format":      "pem",
		"crls":        []string{crl2, crl1}, // Make sure we don't just grab the first colliding entry...
	})
	requireSuccessNonNilResponse(t, resp, err)
	requireFieldsSetInResp(t, resp, "crl")
	pemCrl := resp.Data["crl"].(string)
	combinedCrl, err := decodePemCrl(pemCrl)
	require.NoError(t, err, "failed decoding combined CRL")
	combinedSerials := extractSerialsFromCrl(t, combinedCrl)

	require.Contains(t, combinedSerials, serial1)
	require.Contains(t, combinedSerials, serial2)
	require.Equal(t, 2, len(combinedSerials), "serials contained more serials than expected")

	// Make sure we issued a warning about the time collision
	require.NotEmpty(t, resp.Warnings, "expected at least one warning")
	require.Contains(t, resp.Warnings[0], "different revocation times detected")

	// Make sure we have the initial revocation time from backend 1 within the combined CRL.
	decodedCrl1, err := decodePemCrl(crl1)
	require.NoError(t, err, "failed decoding crl from backend 1")
	serialsFromBackend1 := extractSerialsFromCrl(t, decodedCrl1)

	require.Equal(t, serialsFromBackend1[serial1], combinedSerials[serial1])

	// Make sure we have the initial revocation time from backend 1 does not match with backend 2's time
	decodedCrl2, err := decodePemCrl(crl2)
	require.NoError(t, err, "failed decoding crl from backend 2")
	serialsFromBackend2 := extractSerialsFromCrl(t, decodedCrl2)

	require.NotEqual(t, serialsFromBackend1[serial1], serialsFromBackend2[serial1])
}

func TestResignCrls_DeltaCrl(t *testing.T) {
	t.Parallel()

	b1, s1 := createBackendWithStorage(t)
	b2, s2 := createBackendWithStorage(t)

	// Setup two backends, with the same key material/certificate with a different leaf in each that is revoked.
	caCert, serial1, serial2, crl1, crl2 := setupResignCrlMounts(t, b1, s1, b2, s2)

	resp, err := CBWrite(b1, s1, "issuer/default/resign-crls", map[string]interface{}{
		"crl_number":       "5",
		"delta_crl_number": "4",
		"next_update":      "12h",
		"format":           "pem",
		"crls":             []string{crl1, crl2},
	})
	requireSuccessNonNilResponse(t, resp, err)
	requireFieldsSetInResp(t, resp, "crl")
	pemCrl := resp.Data["crl"].(string)
	combinedCrl, err := decodePemCrl(pemCrl)
	require.NoError(t, err, "failed decoding combined CRL")
	serials := extractSerialsFromCrl(t, combinedCrl)

	require.Contains(t, serials, serial1)
	require.Contains(t, serials, serial2)
	require.Equal(t, 2, len(serials), "serials contained more serials than expected")

	require.Equal(t, big.NewInt(int64(5)), combinedCrl.Number)
	require.Equal(t, combinedCrl.ThisUpdate.Add(12*time.Hour), combinedCrl.NextUpdate)

	extensions := combinedCrl.Extensions
	requireExtensionOid(t, []int{2, 5, 29, 27}, extensions) // Delta CRL Extension
	requireExtensionOid(t, []int{2, 5, 29, 20}, extensions) // CRL Number Extension
	requireExtensionOid(t, []int{2, 5, 29, 35}, extensions) // akidOid
	require.Equal(t, 3, len(extensions))

	err = combinedCrl.CheckSignatureFrom(caCert)
	require.NoError(t, err, "failed signature check of CRL")
}

func setupResignCrlMounts(t *testing.T, b1 *backend, s1 logical.Storage, b2 *backend, s2 logical.Storage) (*x509.Certificate, string, string, string, string) {
	t.Helper()

	// Setup two mounts with the same CA/key material
	resp, err := CBWrite(b1, s1, "root/generate/exported", map[string]interface{}{
		"common_name": "test.com",
	})
	requireSuccessNonNilResponse(t, resp, err)
	requireFieldsSetInResp(t, resp, "certificate", "private_key")
	pemCaCert := resp.Data["certificate"].(string)
	caCert := parseCert(t, pemCaCert)
	privKey := resp.Data["private_key"].(string)

	// Import the above key/cert into another mount
	resp, err = CBWrite(b2, s2, "config/ca", map[string]interface{}{
		"pem_bundle": pemCaCert + "\n" + privKey,
	})
	requireSuccessNonNilResponse(t, resp, err, "error setting up CA on backend 2")

	// Create the same role in both mounts
	resp, err = CBWrite(b1, s1, "roles/test", map[string]interface{}{
		"allowed_domains":  "test.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
	})
	requireSuccessNilResponse(t, resp, err, "error setting up pki role on backend 1")

	resp, err = CBWrite(b2, s2, "roles/test", map[string]interface{}{
		"allowed_domains":  "test.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
	})
	requireSuccessNilResponse(t, resp, err, "error setting up pki role on backend 2")

	// Issue and revoke a cert in backend 1
	resp, err = CBWrite(b1, s1, "issue/test", map[string]interface{}{
		"common_name": "test1.test.com",
	})
	requireSuccessNonNilResponse(t, resp, err, "error issuing cert from backend 1")
	requireFieldsSetInResp(t, resp, "serial_number")
	serial1 := resp.Data["serial_number"].(string)

	resp, err = CBWrite(b1, s1, "revoke", map[string]interface{}{"serial_number": serial1})
	requireSuccessNonNilResponse(t, resp, err, "error revoking cert from backend 2")

	// Issue and revoke a cert in backend 2
	resp, err = CBWrite(b2, s2, "issue/test", map[string]interface{}{
		"common_name": "test1.test.com",
	})
	requireSuccessNonNilResponse(t, resp, err, "error issuing cert from backend 2")
	requireFieldsSetInResp(t, resp, "serial_number")
	serial2 := resp.Data["serial_number"].(string)

	resp, err = CBWrite(b2, s2, "revoke", map[string]interface{}{"serial_number": serial2})
	requireSuccessNonNilResponse(t, resp, err, "error revoking cert from backend 2")

	// Fetch PEM CRLs from each
	resp, err = CBRead(b1, s1, "cert/crl")
	requireSuccessNonNilResponse(t, resp, err, "error fetch crl from backend 1")
	requireFieldsSetInResp(t, resp, "certificate")
	crl1 := resp.Data["certificate"].(string)

	resp, err = CBRead(b2, s2, "cert/crl")
	requireSuccessNonNilResponse(t, resp, err, "error fetch crl from backend 2")
	requireFieldsSetInResp(t, resp, "certificate")
	crl2 := resp.Data["certificate"].(string)
	return caCert, serial1, serial2, crl1, crl2
}

func requireExtensionOid(t *testing.T, identifier asn1.ObjectIdentifier, extensions []pkix.Extension, msgAndArgs ...interface{}) {
	t.Helper()

	found := false
	var oidsInExtensions []string
	for _, extension := range extensions {
		oidsInExtensions = append(oidsInExtensions, extension.Id.String())
		if extension.Id.Equal(identifier) {
			found = true
			break
		}
	}

	if !found {
		msg := fmt.Sprintf("Failed to find matching asn oid %s out of %v", identifier.String(), oidsInExtensions)
		require.Fail(t, msg, msgAndArgs)
	}
}

func extractSerialsFromCrl(t *testing.T, crl *x509.RevocationList) map[string]time.Time {
	serials := map[string]time.Time{}

	for _, revokedCert := range crl.RevokedCertificates {
		serial := serialFromBigInt(revokedCert.SerialNumber)
		if _, exists := serials[serial]; exists {
			t.Fatalf("Serial number %s was duplicated in CRL", serial)
		}
		serials[serial] = revokedCert.RevocationTime
	}
	return serials
}

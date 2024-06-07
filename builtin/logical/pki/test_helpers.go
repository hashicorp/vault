// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io"
	"math"
	"math/big"
	http2 "net/http"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ocsp"
)

// Setup helpers
func CreateBackendWithStorage(t testing.TB) (*backend, logical.Storage) {
	t.Helper()

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
	t.Helper()

	err := client.Sys().Mount(path, &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	require.NoError(t, err, "failed mounting pki endpoint")
}

func mountCertEndpoint(t testing.TB, client *api.Client, path string) {
	t.Helper()

	err := client.Sys().EnableAuthWithOptions(path, &api.MountInput{
		Type: "cert",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	require.NoError(t, err, "failed mounting cert endpoint")
}

// Signing helpers
func requireSignedBy(t *testing.T, cert *x509.Certificate, signingCert *x509.Certificate) {
	t.Helper()

	if err := cert.CheckSignatureFrom(signingCert); err != nil {
		t.Fatalf("signature verification failed: %v", err)
	}
}

func requireSignedByAtPath(t *testing.T, client *api.Client, leaf *x509.Certificate, path string) {
	t.Helper()

	resp, err := client.Logical().Read(path)
	require.NoError(t, err, "got unexpected error fetching parent certificate")
	require.NotNil(t, resp, "missing response when fetching parent certificate")
	require.NotNil(t, resp.Data, "missing data from parent certificate response")
	require.NotNil(t, resp.Data["certificate"], "missing certificate field on parent read response")

	parentCert := resp.Data["certificate"].(string)
	parent := parseCert(t, parentCert)

	requireSignedBy(t, leaf, parent)
}

// Certificate helper
func parseCert(t *testing.T, pemCert string) *x509.Certificate {
	t.Helper()

	block, _ := pem.Decode([]byte(pemCert))
	require.NotNil(t, block, "failed to decode PEM block")

	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)
	return cert
}

func requireMatchingPublicKeys(t *testing.T, cert *x509.Certificate, key crypto.PublicKey) {
	t.Helper()

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
	t.Helper()

	path := fmt.Sprintf("/v1/%s/crl", mountPoint)
	return getParsedCrlAtPath(t, client, path).TBSCertList
}

func parseCrlPemBytes(t *testing.T, crlPem []byte) pkix.TBSCertificateList {
	t.Helper()

	certList, err := x509.ParseCRL(crlPem)
	require.NoError(t, err)
	return certList.TBSCertList
}

func requireSerialNumberInCRL(t *testing.T, revokeList pkix.TBSCertificateList, serialNum string) bool {
	if t != nil {
		t.Helper()
	}

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
	t.Helper()

	path := fmt.Sprintf("/v1/%s/crl", mountPoint)
	return getParsedCrlAtPath(t, client, path)
}

func getParsedCrlAtPath(t *testing.T, client *api.Client, path string) *pkix.CertificateList {
	t.Helper()

	req := client.NewRequest("GET", path)
	resp, err := client.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	crlBytes, err := io.ReadAll(resp.Body)
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
	t.Helper()

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

func CBHeader(b *backend, s logical.Storage, path string) (*logical.Response, error) {
	return CBReq(b, s, logical.HeaderOperation, path, make(map[string]interface{}))
}

func CBRead(b *backend, s logical.Storage, path string) (*logical.Response, error) {
	return CBReq(b, s, logical.ReadOperation, path, make(map[string]interface{}))
}

func CBWrite(b *backend, s logical.Storage, path string, data map[string]interface{}) (*logical.Response, error) {
	return CBReq(b, s, logical.UpdateOperation, path, data)
}

func CBPatch(b *backend, s logical.Storage, path string, data map[string]interface{}) (*logical.Response, error) {
	return CBReq(b, s, logical.PatchOperation, path, data)
}

func CBList(b *backend, s logical.Storage, path string) (*logical.Response, error) {
	return CBReq(b, s, logical.ListOperation, path, make(map[string]interface{}))
}

func CBDelete(b *backend, s logical.Storage, path string) (*logical.Response, error) {
	return CBReq(b, s, logical.DeleteOperation, path, make(map[string]interface{}))
}

func requireFieldsSetInResp(t *testing.T, resp *logical.Response, fields ...string) {
	t.Helper()

	var missingFields []string
	for _, field := range fields {
		value, ok := resp.Data[field]
		if !ok || value == nil {
			missingFields = append(missingFields, field)
		}
	}

	require.Empty(t, missingFields, "The following fields were required but missing from response:\n%v", resp.Data)
}

func requireSuccessNonNilResponse(t *testing.T, resp *logical.Response, err error, msgAndArgs ...interface{}) {
	t.Helper()

	require.NoError(t, err, msgAndArgs...)
	if resp.IsError() {
		errContext := fmt.Sprintf("Expected successful response but got error: %v", resp.Error())
		require.Falsef(t, resp.IsError(), errContext, msgAndArgs...)
	}
	require.NotNil(t, resp, msgAndArgs...)
}

func requireSuccessNilResponse(t *testing.T, resp *logical.Response, err error, msgAndArgs ...interface{}) {
	t.Helper()

	require.NoError(t, err, msgAndArgs...)
	if resp.IsError() {
		errContext := fmt.Sprintf("Expected successful response but got error: %v", resp.Error())
		require.Falsef(t, resp.IsError(), errContext, msgAndArgs...)
	}
	if resp != nil {
		msg := fmt.Sprintf("expected nil response but got: %v", resp)
		require.Nilf(t, resp, msg, msgAndArgs...)
	}
}

func getCRLNumber(t *testing.T, crl pkix.TBSCertificateList) int {
	t.Helper()

	for _, extension := range crl.Extensions {
		if extension.Id.Equal(certutil.CRLNumberOID) {
			bigInt := new(big.Int)
			leftOver, err := asn1.Unmarshal(extension.Value, &bigInt)
			require.NoError(t, err, "Failed unmarshalling crl number extension")
			require.Empty(t, leftOver, "leftover bytes from unmarshalling crl number extension")
			require.True(t, bigInt.IsInt64(), "parsed crl number integer is not an int64")
			require.False(t, math.MaxInt <= bigInt.Int64(), "parsed crl number integer can not fit in an int")
			return int(bigInt.Int64())
		}
	}

	t.Fatalf("failed to find crl number extension")
	return 0
}

func getCrlReferenceFromDelta(t *testing.T, crl pkix.TBSCertificateList) int {
	t.Helper()

	for _, extension := range crl.Extensions {
		if extension.Id.Equal(certutil.DeltaCRLIndicatorOID) {
			bigInt := new(big.Int)
			leftOver, err := asn1.Unmarshal(extension.Value, &bigInt)
			require.NoError(t, err, "Failed unmarshalling delta crl indicator extension")
			require.Empty(t, leftOver, "leftover bytes from unmarshalling delta crl indicator extension")
			require.True(t, bigInt.IsInt64(), "parsed delta crl integer is not an int64")
			require.False(t, math.MaxInt <= bigInt.Int64(), "parsed delta crl integer can not fit in an int")
			return int(bigInt.Int64())
		}
	}

	t.Fatalf("failed to find delta crl indicator extension")
	return 0
}

// waitForUpdatedCrl will wait until the CRL at the provided path has been reloaded
// up for a maxWait duration and gives up if the timeout has been reached. If a negative
// value for lastSeenCRLNumber is provided, the method will load the current CRL and wait
// for a newer CRL be generated.
func waitForUpdatedCrl(t *testing.T, client *api.Client, crlPath string, lastSeenCRLNumber int, maxWait time.Duration) pkix.TBSCertificateList {
	t.Helper()

	newCrl, didTimeOut := waitForUpdatedCrlUntil(t, client, crlPath, lastSeenCRLNumber, maxWait)
	if didTimeOut {
		t.Fatalf("Timed out waiting for new CRL rebuild on path %s", crlPath)
	}
	return newCrl.TBSCertList
}

// waitForUpdatedCrlUntil is a helper method that will wait for a CRL to be updated up until maxWait duration
// or give up and return the last CRL it loaded. It will not fail, if it does not see a new CRL within the
// max duration unlike waitForUpdatedCrl. Returns the last loaded CRL at the provided path and a boolean
// indicating if we hit maxWait duration or not.
func waitForUpdatedCrlUntil(t *testing.T, client *api.Client, crlPath string, lastSeenCrlNumber int, maxWait time.Duration) (*pkix.CertificateList, bool) {
	t.Helper()

	crl := getParsedCrlAtPath(t, client, crlPath)
	initialCrlRevision := getCRLNumber(t, crl.TBSCertList)
	newCrlRevision := initialCrlRevision

	// Short circuit the fetches if we have a version of the CRL we want
	if lastSeenCrlNumber > 0 && getCRLNumber(t, crl.TBSCertList) > lastSeenCrlNumber {
		return crl, false
	}

	start := time.Now()
	iteration := 0
	for {
		iteration++

		if time.Since(start) > maxWait {
			t.Logf("Timed out waiting for new CRL on path %s after iteration %d, delay: %v",
				crlPath, iteration, time.Now().Sub(start))
			return crl, true
		}

		crl = getParsedCrlAtPath(t, client, crlPath)
		newCrlRevision = getCRLNumber(t, crl.TBSCertList)
		if newCrlRevision > initialCrlRevision {
			t.Logf("Got new revision of CRL %s from %d to %d after iteration %d, delay %v",
				crlPath, initialCrlRevision, newCrlRevision, iteration, time.Now().Sub(start))
			return crl, false
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// A quick CRL to string to provide better test error messages
func summarizeCrl(t *testing.T, crl pkix.TBSCertificateList) string {
	version := getCRLNumber(t, crl)
	serials := []string{}
	for _, cert := range crl.RevokedCertificates {
		serials = append(serials, normalizeSerialFromBigInt(cert.SerialNumber))
	}
	return fmt.Sprintf("CRL Version: %d\n"+
		"This Update: %s\n"+
		"Next Update: %s\n"+
		"Revoked Serial Count: %d\n"+
		"Revoked Serials: %v", version, crl.ThisUpdate, crl.NextUpdate, len(serials), serials)
}

// OCSP helpers
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

func performOcspPost(t *testing.T, cert *x509.Certificate, issuerCert *x509.Certificate, client *api.Client, ocspPath string) *ocsp.Response {
	t.Helper()

	baseClient := client.WithNamespace("")

	ocspReq := generateRequest(t, crypto.SHA256, cert, issuerCert)
	ocspPostReq := baseClient.NewRequest(http2.MethodPost, ocspPath)
	ocspPostReq.Headers.Set("Content-Type", "application/ocsp-request")
	ocspPostReq.BodyBytes = ocspReq
	rawResp, err := baseClient.RawRequest(ocspPostReq)
	require.NoError(t, err, "failed sending unified-ocsp post request")

	require.Equal(t, 200, rawResp.StatusCode)
	require.Equal(t, ocspResponseContentType, rawResp.Header.Get("Content-Type"))
	bodyReader := rawResp.Body
	respDer, err := io.ReadAll(bodyReader)
	bodyReader.Close()
	require.NoError(t, err, "failed reading response body")

	ocspResp, err := ocsp.ParseResponse(respDer, issuerCert)
	require.NoError(t, err, "parsing ocsp get response")
	return ocspResp
}

func requireCertMissingFromStorage(t *testing.T, client *api.Client, cert *x509.Certificate) {
	serial := serialFromCert(cert)
	requireSerialMissingFromStorage(t, client, serial)
}

func requireSerialMissingFromStorage(t *testing.T, client *api.Client, serial string) {
	resp, err := client.Logical().ReadWithContext(context.Background(), "pki/cert/"+serial)
	require.NoErrorf(t, err, "failed reading certificate with serial %s", serial)
	require.Nilf(t, resp, "expected a nil response looking up serial %s got: %v", serial, resp)
}

func requireCertInStorage(t *testing.T, client *api.Client, cert *x509.Certificate) {
	serial := serialFromCert(cert)
	requireSerialInStorage(t, client, serial)
}

func requireSerialInStorage(t *testing.T, client *api.Client, serial string) {
	resp, err := client.Logical().ReadWithContext(context.Background(), "pki/cert/"+serial)
	require.NoErrorf(t, err, "failed reading certificate with serial %s", serial)
	require.NotNilf(t, resp, "reading certificate returned a nil response for serial: %s", serial)
	require.NotNilf(t, resp.Data, "reading certificate returned a nil data response for serial: %s", serial)
	require.NotEmpty(t, resp.Data["certificate"], "certificate field was empty for serial: %s", serial)
}

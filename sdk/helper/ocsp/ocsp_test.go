// Copyright (c) 2017-2022 Snowflake Computing Inc. All rights reserved.

package ocsp

import (
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/pki"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ocsp"
)

func TestOCSP(t *testing.T) {
	targetURL := []string{
		"https://sfcdev1.blob.core.windows.net/",
		"https://sfctest0.snowflakecomputing.com/",
		"https://s3-us-west-2.amazonaws.com/sfc-snowsql-updates/?prefix=1.1/windows_x86_64",
	}

	conf := VerifyConfig{
		OcspFailureMode: FailOpenFalse,
	}
	c := New(testLogFactory, 10)
	transports := []*http.Transport{
		newInsecureOcspTransport(nil),
		c.NewTransport(&conf),
	}

	for _, tgt := range targetURL {
		c.ocspResponseCache, _ = lru.New2Q(10)
		for _, tr := range transports {
			c := &http.Client{
				Transport: tr,
				Timeout:   30 * time.Second,
			}
			req, err := http.NewRequest("GET", tgt, bytes.NewReader(nil))
			if err != nil {
				t.Fatalf("fail to create a request. err: %v", err)
			}
			res, err := c.Do(req)
			if err != nil {
				t.Fatalf("failed to GET contents. err: %v", err)
			}
			defer res.Body.Close()
			_, err = ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("failed to read content body for %v", tgt)
			}

		}
	}
}

/**
// Used for development, requires an active Vault with PKI setup
func TestMultiOCSP(t *testing.T) {

	targetURL := []string{
		"https://localhost:8200/v1/pki/ocsp",
		"https://localhost:8200/v1/pki/ocsp",
		"https://localhost:8200/v1/pki/ocsp",
	}

	b, _ := pem.Decode([]byte(vaultCert))
	caCert, _ := x509.ParseCertificate(b.Bytes)
	conf := VerifyConfig{
		OcspFailureMode:     FailOpenFalse,
		QueryAllServers:     true,
		OcspServersOverride: targetURL,
		ExtraCas:            []*x509.Certificate{caCert},
	}
	c := New(testLogFactory, 10)
	transports := []*http.Transport{
		newInsecureOcspTransport(conf.ExtraCas),
		c.NewTransport(&conf),
	}

	tgt := "https://localhost:8200/v1/pki/ca/pem"
	c.ocspResponseCache, _ = lru.New2Q(10)
	for _, tr := range transports {
		c := &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		}
		req, err := http.NewRequest("GET", tgt, bytes.NewReader(nil))
		if err != nil {
			t.Fatalf("fail to create a request. err: %v", err)
		}
		res, err := c.Do(req)
		if err != nil {
			t.Fatalf("failed to GET contents. err: %v", err)
		}
		defer res.Body.Close()
		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("failed to read content body for %v", tgt)
		}
	}
}
*/

func TestUnitEncodeCertIDGood(t *testing.T) {
	targetURLs := []string{
		"faketestaccount.snowflakecomputing.com:443",
		"s3-us-west-2.amazonaws.com:443",
		"sfcdev1.blob.core.windows.net:443",
	}
	for _, tt := range targetURLs {
		chainedCerts := getCert(tt)
		for i := 0; i < len(chainedCerts)-1; i++ {
			subject := chainedCerts[i]
			issuer := chainedCerts[i+1]
			ocspServers := subject.OCSPServer
			if len(ocspServers) == 0 {
				t.Fatalf("no OCSP server is found. cert: %v", subject.Subject)
			}
			ocspReq, err := ocsp.CreateRequest(subject, issuer, &ocsp.RequestOptions{})
			if err != nil {
				t.Fatalf("failed to create OCSP request. err: %v", err)
			}
			var ost *ocspStatus
			_, ost = extractCertIDKeyFromRequest(ocspReq)
			if ost.err != nil {
				t.Fatalf("failed to extract cert ID from the OCSP request. err: %v", ost.err)
			}
			// better hash. Not sure if the actual OCSP server accepts this, though.
			ocspReq, err = ocsp.CreateRequest(subject, issuer, &ocsp.RequestOptions{Hash: crypto.SHA512})
			if err != nil {
				t.Fatalf("failed to create OCSP request. err: %v", err)
			}
			_, ost = extractCertIDKeyFromRequest(ocspReq)
			if ost.err != nil {
				t.Fatalf("failed to extract cert ID from the OCSP request. err: %v", ost.err)
			}
			// tweaked request binary
			ocspReq, err = ocsp.CreateRequest(subject, issuer, &ocsp.RequestOptions{Hash: crypto.SHA512})
			if err != nil {
				t.Fatalf("failed to create OCSP request. err: %v", err)
			}
			ocspReq[10] = 0 // random change
			_, ost = extractCertIDKeyFromRequest(ocspReq)
			if ost.err == nil {
				t.Fatal("should have failed")
			}
		}
	}
}

func TestUnitCheckOCSPResponseCache(t *testing.T) {
	c := New(testLogFactory, 10)
	dummyKey0 := certIDKey{
		NameHash:      "dummy0",
		IssuerKeyHash: "dummy0",
		SerialNumber:  "dummy0",
	}
	dummyKey := certIDKey{
		NameHash:      "dummy1",
		IssuerKeyHash: "dummy1",
		SerialNumber:  "dummy1",
	}
	currentTime := float64(time.Now().UTC().Unix())
	c.ocspResponseCache.Add(dummyKey0, &ocspCachedResponse{time: currentTime})
	subject := &x509.Certificate{}
	issuer := &x509.Certificate{}
	ost, err := c.checkOCSPResponseCache(&dummyKey, subject, issuer)
	if err != nil {
		t.Fatal(err)
	}
	if ost.code != ocspMissedCache {
		t.Fatalf("should have failed. expected: %v, got: %v", ocspMissedCache, ost.code)
	}
	// old timestamp
	c.ocspResponseCache.Add(dummyKey, &ocspCachedResponse{time: float64(1395054952)})
	ost, err = c.checkOCSPResponseCache(&dummyKey, subject, issuer)
	if err != nil {
		t.Fatal(err)
	}
	if ost.code != ocspCacheExpired {
		t.Fatalf("should have failed. expected: %v, got: %v", ocspCacheExpired, ost.code)
	}

	// invalid validity
	c.ocspResponseCache.Add(dummyKey, &ocspCachedResponse{time: float64(currentTime - 1000)})
	ost, err = c.checkOCSPResponseCache(&dummyKey, subject, nil)
	if err == nil && isValidOCSPStatus(ost.code) {
		t.Fatalf("should have failed.")
	}
}

func TestUnitValidateOCSP(t *testing.T) {
	ocspRes := &ocsp.Response{}
	ost, err := validateOCSP(ocspRes)
	if err == nil && isValidOCSPStatus(ost.code) {
		t.Fatalf("should have failed.")
	}

	currentTime := time.Now()
	ocspRes.ThisUpdate = currentTime.Add(-2 * time.Hour)
	ocspRes.NextUpdate = currentTime.Add(2 * time.Hour)
	ocspRes.Status = ocsp.Revoked
	ost, err = validateOCSP(ocspRes)
	if err != nil {
		t.Fatal(err)
	}

	if ost.code != ocspStatusRevoked {
		t.Fatalf("should have failed. expected: %v, got: %v", ocspStatusRevoked, ost.code)
	}
	ocspRes.Status = ocsp.Good
	ost, err = validateOCSP(ocspRes)
	if err != nil {
		t.Fatal(err)
	}

	if ost.code != ocspStatusGood {
		t.Fatalf("should have success. expected: %v, got: %v", ocspStatusGood, ost.code)
	}
	ocspRes.Status = ocsp.Unknown
	ost, err = validateOCSP(ocspRes)
	if err != nil {
		t.Fatal(err)
	}
	if ost.code != ocspStatusUnknown {
		t.Fatalf("should have failed. expected: %v, got: %v", ocspStatusUnknown, ost.code)
	}
	ocspRes.Status = ocsp.ServerFailed
	ost, err = validateOCSP(ocspRes)
	if err != nil {
		t.Fatal(err)
	}
	if ost.code != ocspStatusOthers {
		t.Fatalf("should have failed. expected: %v, got: %v", ocspStatusOthers, ost.code)
	}
}

func TestUnitEncodeCertID(t *testing.T) {
	var st *ocspStatus
	_, st = extractCertIDKeyFromRequest([]byte{0x1, 0x2})
	if st.code != ocspFailedDecomposeRequest {
		t.Fatalf("failed to get OCSP status. expected: %v, got: %v", ocspFailedDecomposeRequest, st.code)
	}
}

func getCert(addr string) []*x509.Certificate {
	tcpConn, err := net.DialTimeout("tcp", addr, 40*time.Second)
	if err != nil {
		panic(err)
	}
	defer tcpConn.Close()

	err = tcpConn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		panic(err)
	}
	config := tls.Config{InsecureSkipVerify: true, ServerName: addr}

	conn := tls.Client(tcpConn, &config)
	defer conn.Close()

	err = conn.Handshake()
	if err != nil {
		panic(err)
	}

	state := conn.ConnectionState()

	return state.PeerCertificates
}

func TestOCSPRetry(t *testing.T) {
	c := New(testLogFactory, 10)
	certs := getCert("s3-us-west-2.amazonaws.com:443")
	dummyOCSPHost := &url.URL{
		Scheme: "https",
		Host:   "dummyOCSPHost",
	}
	client := &fakeHTTPClient{
		cnt:     3,
		success: true,
		body:    []byte{1, 2, 3},
		logger:  hclog.New(hclog.DefaultOptions),
		t:       t,
	}
	res, b, st, err := c.retryOCSP(
		context.TODO(),
		client, fakeRequestFunc,
		dummyOCSPHost,
		make(map[string]string), []byte{0}, certs[len(certs)-1])
	if err == nil {
		fmt.Printf("should fail: %v, %v, %v\n", res, b, st)
	}
	client = &fakeHTTPClient{
		cnt:     30,
		success: true,
		body:    []byte{1, 2, 3},
		logger:  hclog.New(hclog.DefaultOptions),
		t:       t,
	}
	res, b, st, err = c.retryOCSP(
		context.TODO(),
		client, fakeRequestFunc,
		dummyOCSPHost,
		make(map[string]string), []byte{0}, certs[len(certs)-1])
	if err == nil {
		fmt.Printf("should fail: %v, %v, %v\n", res, b, st)
	}
}

type tcCanEarlyExit struct {
	results       []*ocspStatus
	resultLen     int
	retFailOpen   *ocspStatus
	retFailClosed *ocspStatus
}

func TestCanEarlyExitForOCSP(t *testing.T) {
	testcases := []tcCanEarlyExit{
		{ // 0
			results: []*ocspStatus{
				{
					code: ocspStatusGood,
				},
				{
					code: ocspStatusGood,
				},
				{
					code: ocspStatusGood,
				},
			},
			retFailOpen:   nil,
			retFailClosed: nil,
		},
		{ // 1
			results: []*ocspStatus{
				{
					code: ocspStatusRevoked,
					err:  errors.New("revoked"),
				},
				{
					code: ocspStatusGood,
				},
				{
					code: ocspStatusGood,
				},
			},
			retFailOpen:   &ocspStatus{ocspStatusRevoked, errors.New("revoked")},
			retFailClosed: &ocspStatus{ocspStatusRevoked, errors.New("revoked")},
		},
		{ // 2
			results: []*ocspStatus{
				{
					code: ocspStatusUnknown,
					err:  errors.New("unknown"),
				},
				{
					code: ocspStatusGood,
				},
				{
					code: ocspStatusGood,
				},
			},
			retFailOpen:   nil,
			retFailClosed: &ocspStatus{ocspStatusUnknown, errors.New("unknown")},
		},
		{ // 3: not taken as revoked if any invalid OCSP response (ocspInvalidValidity) is included.
			results: []*ocspStatus{
				{
					code: ocspStatusRevoked,
					err:  errors.New("revoked"),
				},
				{
					code: ocspInvalidValidity,
				},
				{
					code: ocspStatusGood,
				},
			},
			retFailOpen:   nil,
			retFailClosed: &ocspStatus{ocspStatusRevoked, errors.New("revoked")},
		},
		{ // 4: not taken as revoked if the number of results don't match the expected results.
			results: []*ocspStatus{
				{
					code: ocspStatusRevoked,
					err:  errors.New("revoked"),
				},
				{
					code: ocspStatusGood,
				},
			},
			resultLen:     3,
			retFailOpen:   nil,
			retFailClosed: &ocspStatus{ocspStatusRevoked, errors.New("revoked")},
		},
	}
	c := New(testLogFactory, 10)
	for idx, tt := range testcases {
		expectedLen := len(tt.results)
		if tt.resultLen > 0 {
			expectedLen = tt.resultLen
		}
		r := c.canEarlyExitForOCSP(tt.results, expectedLen, &VerifyConfig{OcspFailureMode: FailOpenTrue})
		if !(tt.retFailOpen == nil && r == nil) && !(tt.retFailOpen != nil && r != nil && tt.retFailOpen.code == r.code) {
			t.Fatalf("%d: failed to match return. expected: %v, got: %v", idx, tt.retFailOpen, r)
		}
		r = c.canEarlyExitForOCSP(tt.results, expectedLen, &VerifyConfig{OcspFailureMode: FailOpenFalse})
		if !(tt.retFailClosed == nil && r == nil) && !(tt.retFailClosed != nil && r != nil && tt.retFailClosed.code == r.code) {
			t.Fatalf("%d: failed to match return. expected: %v, got: %v", idx, tt.retFailClosed, r)
		}
	}
}

func TestWithVaultPKI(t *testing.T) {
	t.Parallel()

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": pki.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()
	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)
	client := cores[0].Client

	err := client.Sys().Mount("pki", &api.MountInput{
		Type: "pki",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "16h",
			MaxLeaseTTL:     "32h",
		},
	})
	require.NoError(t, err)

	resp, err := client.Logical().Write("pki/root/generate/internal", map[string]interface{}{
		"ttl":         "40h",
		"common_name": "Root R1",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["issuer_id"])
	rootIssuerId := resp.Data["issuer_id"].(string)

	// Set URLs pointing to the issuer.
	_, err = client.Logical().Write("pki/config/cluster", map[string]interface{}{
		"path":     client.Address() + "/v1/pki",
		"aia_path": client.Address() + "/v1/pki",
	})
	require.NoError(t, err)

	_, err = client.Logical().Write("pki/config/urls", map[string]interface{}{
		"enable_templating":       true,
		"crl_distribution_points": "{{cluster_aia_path}}/issuer/{{issuer_id}}/crl/der",
		"issuing_certificates":    "{{cluster_aia_path}}/issuer/{{issuer_id}}/der",
		"ocsp_servers":            "{{cluster_aia_path}}/ocsp",
	})
	require.NoError(t, err)

	// Build an intermediate CA
	resp, err = client.Logical().Write("pki/intermediate/generate/internal", map[string]interface{}{
		"common_name": "Int X1",
		"key_type":    "ec",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["csr"])
	intermediateCSR := resp.Data["csr"].(string)

	resp, err = client.Logical().Write("pki/root/sign-intermediate", map[string]interface{}{
		"csr": intermediateCSR,
		"ttl": "20h",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["certificate"])
	intermediateCert := resp.Data["certificate"]

	resp, err = client.Logical().Write("pki/intermediate/set-signed", map[string]interface{}{
		"certificate": intermediateCert,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.Data)
	require.NotEmpty(t, resp.Data["imported_issuers"])
	rawImportedIssuers := resp.Data["imported_issuers"].([]interface{})
	require.Equal(t, len(rawImportedIssuers), 1)
	importedIssuer := rawImportedIssuers[0].(string)
	require.NotEmpty(t, importedIssuer)

	// Set intermediate as default.
	_, err = client.Logical().Write("pki/config/issuers", map[string]interface{}{
		"default": importedIssuer,
	})
	require.NoError(t, err)

	// Setup roles for root, intermediate.
	_, err = client.Logical().Write("pki/roles/example-root", map[string]interface{}{
		"allowed_domains":  "example.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
		"key_type":         "ec",
		"issuer_ref":       rootIssuerId,
	})
	require.NoError(t, err)

	_, err = client.Logical().Write("pki/roles/example-int", map[string]interface{}{
		"allowed_domains":  "example.com",
		"allow_subdomains": "true",
		"max_ttl":          "1h",
		"key_type":         "ec",
	})
	require.NoError(t, err)

	// Issue certs and validate them against OCSP.
	for _, path := range []string{"pki/issue/example-int", "pki/issue/example-root"} {
		t.Logf("Validating against path: %v", path)
		resp, err = client.Logical().Write(path, map[string]interface{}{
			"common_name": "test.example.com",
			"ttl":         "5m",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotNil(t, resp.Data)
		require.NotEmpty(t, resp.Data["certificate"])
		require.NotEmpty(t, resp.Data["issuing_ca"])
		require.NotEmpty(t, resp.Data["serial_number"])

		certPEM := resp.Data["certificate"].(string)
		certBlock, _ := pem.Decode([]byte(certPEM))
		require.NotNil(t, certBlock)
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		require.NoError(t, err)
		require.NotNil(t, cert)

		issuerPEM := resp.Data["issuing_ca"].(string)
		issuerBlock, _ := pem.Decode([]byte(issuerPEM))
		require.NotNil(t, issuerBlock)
		issuer, err := x509.ParseCertificate(issuerBlock.Bytes)
		require.NoError(t, err)
		require.NotNil(t, issuer)

		serialNumber := resp.Data["serial_number"].(string)

		conf := &VerifyConfig{
			OcspFailureMode: FailOpenFalse,
			ExtraCas:        []*x509.Certificate{cluster.CACert},
		}
		ocspClient := New(testLogFactory, 10)

		err = ocspClient.VerifyLeafCertificate(context.Background(), cert, issuer, conf)
		require.NoError(t, err)

		_, err = client.Logical().Write("pki/revoke", map[string]interface{}{
			"serial_number": serialNumber,
		})
		require.NoError(t, err)

		err = ocspClient.VerifyLeafCertificate(context.Background(), cert, issuer, conf)
		require.Error(t, err)
	}
}

var testLogger = hclog.New(hclog.DefaultOptions)

func testLogFactory() hclog.Logger {
	return testLogger
}

type fakeHTTPClient struct {
	cnt        int    // number of retry
	success    bool   // return success after retry in cnt times
	timeout    bool   // timeout
	body       []byte // return body
	t          *testing.T
	logger     hclog.Logger
	redirected bool
}

func (c *fakeHTTPClient) Do(_ *retryablehttp.Request) (*http.Response, error) {
	c.cnt--
	if c.cnt < 0 {
		c.cnt = 0
	}
	c.t.Log("fakeHTTPClient.cnt", c.cnt)

	var retcode int
	if !c.redirected {
		c.redirected = true
		c.cnt++
		retcode = 405
	} else if c.success && c.cnt == 1 {
		retcode = 200
	} else {
		if c.timeout {
			// simulate timeout
			time.Sleep(time.Second * 1)
			return nil, &fakeHTTPError{
				err:     "Whatever reason (Client.Timeout exceeded while awaiting headers)",
				timeout: true,
			}
		}
		retcode = 0
	}

	ret := &http.Response{
		StatusCode: retcode,
		Body:       &fakeResponseBody{body: c.body},
	}
	return ret, nil
}

type fakeHTTPError struct {
	err     string
	timeout bool
}

func (e *fakeHTTPError) Error() string   { return e.err }
func (e *fakeHTTPError) Timeout() bool   { return e.timeout }
func (e *fakeHTTPError) Temporary() bool { return true }

type fakeResponseBody struct {
	body []byte
	cnt  int
}

func (b *fakeResponseBody) Read(p []byte) (n int, err error) {
	if b.cnt == 0 {
		copy(p, b.body)
		b.cnt = 1
		return len(b.body), nil
	}
	b.cnt = 0
	return 0, io.EOF
}

func (b *fakeResponseBody) Close() error {
	return nil
}

func fakeRequestFunc(_, _ string, _ interface{}) (*retryablehttp.Request, error) {
	return nil, nil
}

const vaultCert = `-----BEGIN CERTIFICATE-----
MIIDuTCCAqGgAwIBAgIUA6VeVD1IB5rXcCZRAqPO4zr/GAMwDQYJKoZIhvcNAQEL
BQAwcjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAlZBMREwDwYDVQQHDAhTb21lQ2l0
eTESMBAGA1UECgwJTXlDb21wYW55MRMwEQYDVQQLDApNeURpdmlzaW9uMRowGAYD
VQQDDBF3d3cuY29uaHVnZWNvLmNvbTAeFw0yMjA5MDcxOTA1MzdaFw0yNDA5MDYx
OTA1MzdaMHIxCzAJBgNVBAYTAlVTMQswCQYDVQQIDAJWQTERMA8GA1UEBwwIU29t
ZUNpdHkxEjAQBgNVBAoMCU15Q29tcGFueTETMBEGA1UECwwKTXlEaXZpc2lvbjEa
MBgGA1UEAwwRd3d3LmNvbmh1Z2Vjby5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQDL9qzEXi4PIafSAqfcwcmjujFvbG1QZbI8swxnD+w8i4ufAQU5
LDmvMrGo3ZbhJ0mCihYmFxpjhRdP2raJQ9TysHlPXHtDRpr9ckWTKBz2oIfqVtJ2
qzteQkWCkDAO7kPqzgCFsMeoMZeONRkeGib0lEzQAbW/Rqnphg8zVVkyQ71DZ7Pc
d5WkC2E28kKcSramhWfVFpxG3hSIrLOX2esEXteLRzKxFPf+gi413JZFKYIWrebP
u5t0++MLNpuX322geoki4BWMjQsd47XILmxZ4aj33ScZvdrZESCnwP76hKIxg9mO
lMxrqSWKVV5jHZrElSEj9LYJgDO1Y6eItn7hAgMBAAGjRzBFMAsGA1UdDwQEAwIE
MDATBgNVHSUEDDAKBggrBgEFBQcDATAhBgNVHREEGjAYggtleGFtcGxlLmNvbYIJ
bG9jYWxob3N0MA0GCSqGSIb3DQEBCwUAA4IBAQA5dPdf5SdtMwe2uSspO/EuWqbM
497vMQBW1Ey8KRKasJjhvOVYMbe7De5YsnW4bn8u5pl0zQGF4hEtpmifAtVvziH/
K+ritQj9VVNbLLCbFcg+b0kfjt4yrDZ64vWvIeCgPjG1Kme8gdUUWgu9dOud5gdx
qg/tIFv4TRS/eIIymMlfd9owOD3Ig6S5fy4NaAJFAwXf8+3Rzuc+e7JSAPgAufjh
tOTWinxvoiOLuYwo9CyGgq4qKBFsrY0aE0gdA7oTQkpbEbo2EbqiWUl/PTCl1Y4Z
nSZ0n+4q9QC9RLrWwYTwh838d5RVLUst2mBKSA+vn7YkqmBJbdBC6nkd7n7H
-----END CERTIFICATE-----
`

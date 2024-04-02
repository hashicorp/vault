// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"math/big"
	mathrand "math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/ocsp"
)

var ocspPort int

var source InMemorySource

type testLogger struct{}

func (t *testLogger) Log(args ...any) {
	fmt.Printf("%v", args)
}

func TestMain(m *testing.M) {
	source = make(InMemorySource)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return
	}

	ocspPort = listener.Addr().(*net.TCPAddr).Port
	srv := &http.Server{
		Addr:    "localhost:0",
		Handler: NewResponder(&testLogger{}, source, nil),
	}
	go func() {
		srv.Serve(listener)
	}()
	defer srv.Shutdown(context.Background())
	m.Run()
}

func TestCert_RoleResolve(t *testing.T) {
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:    []string{"example.com"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}

	tempDir, connState, err := generateTestCertAndConnState(t, certTemplate)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "example.com"}, false),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleWithName(t, connState, "web"),
			// Test with caching disabled
			testAccStepSetRoleCacheSize(t, -1),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleWithName(t, connState, "web"),
		},
	})
}

func testAccStepResolveRoleWithName(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.ResolveRoleOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		Check: func(resp *logical.Response) error {
			if resp.Data["role"] != certName {
				t.Fatalf("Role was not as expected. Expected %s, received %s", certName, resp.Data["role"])
			}
			return nil
		},
		Data: map[string]interface{}{
			"name": certName,
		},
	}
}

func TestCert_RoleResolveWithoutProvidingCertName(t *testing.T) {
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:    []string{"example.com"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}

	tempDir, connState, err := generateTestCertAndConnState(t, certTemplate)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "example.com"}, false),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleWithEmptyDataMap(t, connState, "web"),
			testAccStepSetRoleCacheSize(t, -1),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleWithEmptyDataMap(t, connState, "web"),
		},
	})
}

func testAccStepSetRoleCacheSize(t *testing.T, size int) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"role_cache_size": size,
		},
	}
}

func testAccStepResolveRoleWithEmptyDataMap(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.ResolveRoleOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		Check: func(resp *logical.Response) error {
			if resp.Data["role"] != certName {
				t.Fatalf("Role was not as expected. Expected %s, received %s", certName, resp.Data["role"])
			}
			return nil
		},
		Data: map[string]interface{}{},
	}
}

func testAccStepResolveRoleExpectRoleResolutionToFail(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.ResolveRoleOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		ErrorOk:         true,
		Check: func(resp *logical.Response) error {
			if resp == nil && !resp.IsError() {
				t.Fatalf("Response was not an error: resp:%#v", resp)
			}

			errString, ok := resp.Data["error"].(string)
			if !ok {
				t.Fatal("Error not part of response.")
			}

			if !strings.Contains(errString, "invalid certificate") {
				t.Fatalf("Error was not due to invalid role name. Error: %s", errString)
			}
			return nil
		},
		Data: map[string]interface{}{
			"name": certName,
		},
	}
}

func testAccStepResolveRoleOCSPFail(t *testing.T, connState tls.ConnectionState, certName string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation:       logical.ResolveRoleOperation,
		Path:            "login",
		Unauthenticated: true,
		ConnState:       &connState,
		ErrorOk:         true,
		Check: func(resp *logical.Response) error {
			if resp == nil || !resp.IsError() {
				t.Fatalf("Response was not an error: resp:%#v", resp)
			}

			errString, ok := resp.Data["error"].(string)
			if !ok {
				t.Fatal("Error not part of response.")
			}

			if !strings.Contains(errString, "no chain matching") {
				t.Fatalf("Error was not due to OCSP failure. Error: %s", errString)
			}
			return nil
		},
		Data: map[string]interface{}{
			"name": certName,
		},
	}
}

func TestCert_RoleResolve_RoleDoesNotExist(t *testing.T) {
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:    []string{"example.com"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
	}

	tempDir, connState, err := generateTestCertAndConnState(t, certTemplate)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	logicaltest.Test(t, logicaltest.TestCase{
		CredentialBackend: testFactory(t),
		Steps: []logicaltest.TestStep{
			testAccStepCert(t, "web", ca, "foo", allowed{dns: "example.com"}, false),
			testAccStepLoginWithName(t, connState, "web"),
			testAccStepResolveRoleExpectRoleResolutionToFail(t, connState, "notweb"),
		},
	})
}

func TestCert_RoleResolveOCSP(t *testing.T) {
	cases := []struct {
		name        string
		failOpen    bool
		certStatus  int
		errExpected bool
	}{
		{"failFalseGoodCert", false, ocsp.Good, false},
		{"failFalseRevokedCert", false, ocsp.Revoked, true},
		{"failFalseUnknownCert", false, ocsp.Unknown, true},
		{"failTrueGoodCert", true, ocsp.Good, false},
		{"failTrueRevokedCert", true, ocsp.Revoked, true},
		{"failTrueUnknownCert", true, ocsp.Unknown, false},
	}
	certTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:    []string{"example.com"},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
		SerialNumber: big.NewInt(mathrand.Int63()),
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(262980 * time.Hour),
		OCSPServer:   []string{fmt.Sprintf("http://localhost:%d", ocspPort)},
	}
	tempDir, connState, err := generateTestCertAndConnState(t, certTemplate)
	if tempDir != "" {
		defer os.RemoveAll(tempDir)
	}
	if err != nil {
		t.Fatalf("error testing connection state: %v", err)
	}
	ca, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_cert.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	issuer := parsePEM(ca)
	pkf, err := ioutil.ReadFile(filepath.Join(tempDir, "ca_key.pem"))
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	pk, err := certutil.ParsePEMBundle(string(pkf))
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp, err := ocsp.CreateResponse(issuer[0], issuer[0], ocsp.Response{
				Status:       c.certStatus,
				SerialNumber: certTemplate.SerialNumber,
				ProducedAt:   time.Now(),
				ThisUpdate:   time.Now(),
				NextUpdate:   time.Now().Add(time.Hour),
			}, pk.PrivateKey)
			if err != nil {
				t.Fatal(err)
			}
			source[certTemplate.SerialNumber.String()] = resp

			b := testFactory(t)
			b.(*backend).ocspClient.ClearCache()
			var resolveStep logicaltest.TestStep
			var loginStep logicaltest.TestStep
			if c.errExpected {
				loginStep = testAccStepLoginWithNameInvalid(t, connState, "web")
				resolveStep = testAccStepResolveRoleOCSPFail(t, connState, "web")
			} else {
				loginStep = testAccStepLoginWithName(t, connState, "web")
				resolveStep = testAccStepResolveRoleWithName(t, connState, "web")
			}
			logicaltest.Test(t, logicaltest.TestCase{
				CredentialBackend: b,
				Steps: []logicaltest.TestStep{
					testAccStepCertWithExtraParams(t, "web", ca, "foo", allowed{dns: "example.com"}, false,
						map[string]interface{}{"ocsp_enabled": true, "ocsp_fail_open": c.failOpen}),
					testAccStepReadCertPolicy(t, "web", false, map[string]interface{}{"ocsp_enabled": true, "ocsp_fail_open": c.failOpen}),
					loginStep,
					resolveStep,
				},
			})
		})
	}
}

func serialFromBigInt(serial *big.Int) string {
	return strings.TrimSpace(certutil.GetHexFormatted(serial.Bytes(), ":"))
}

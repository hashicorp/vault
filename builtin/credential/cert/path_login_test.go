package cert

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"math/big"
	mathrand "math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"

	"github.com/hashicorp/vault/sdk/logical"
)

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
		},
	})
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

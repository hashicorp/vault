package diagnose

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
)

// TestTLSValidCert is the positive test case to show that specifying a valid cert and key
// passes all checks.
func TestTLSValidCert(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./test-fixtures/goodcertwithroot.pem",
				TLSKeyFile:                    "./test-fixtures/goodkey.pem",
				TLSMinVersion:                 "tls10",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

// TestTLSFakeCert simply ensures that the certificate file must contain PEM data.
func TestTLSFakeCert(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./test-fixtures/fakecert.pem",
				TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
				TLSMinVersion:                 "tls10",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(err.Error(), "could not decode cert") {
		t.Fatalf("Bad error message: %s", err)
	}
}

// TestTLSTrailingData uses a certificate from:
// https://github.com/golang/go/issues/40545 that contains
// an extra DER sequence, and makes sure a trailing data error
// is returned.
func TestTLSTrailingData(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./test-fixtures/trailingdatacert.pem",
				TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
				TLSMinVersion:                 "tls10",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(err.Error(), "asn1: syntax error: trailing data") {
		t.Fatalf("Bad error message: %s", err)
	}
}

// TestTLSExpiredCert checks that an expired certificate fails TLS checks
// with an appropriate error.
func TestTLSExpiredCert(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./test-fixtures/expiredcert.pem",
				TLSKeyFile:                    "./test-fixtures/expiredprivatekey.pem",
				TLSMinVersion:                 "tls10",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(err.Error(), "certificate has expired or is not yet valid") {
		t.Fatalf("Bad error message: %s", err)
	}
}

// TestTLSMismatchedCryptographicInfo verifies that a cert and key of differing cryptographic
// types, when specified together, is met with a unique error message.
func TestTLSMismatchedCryptographicInfo(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
				TLSKeyFile:                    "./test-fixtures/ecdsa.key",
				TLSMinVersion:                 "tls10",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if err.Error() != "tls: private key type does not match public key type" {
		t.Fatalf("Bad error message: %s", err)
	}

	listeners = []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./test-fixtures/ecdsa.crt",
				TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
				TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
				TLSMinVersion:                 "tls10",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err = ListenerChecks(listeners)
	if err == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if err.Error() != "tls: private key type does not match public key type" {
		t.Fatalf("Bad error message: %s", err)
	}
}

// TestTLSMultiKeys verifies that a unique error message is thrown when a key is specified twice.
func TestTLSMultiKeys(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./../../api/test-fixtures/keys/key.pem",
				TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
				TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
				TLSMinVersion:                 "tls10",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(err.Error(), "pem block does not parse to a certificate") {
		t.Fatalf("Bad error message: %s", err)
	}
}

// TestTLSMultiCerts verifies that a unique error message is thrown when a cert is specified twice.
func TestTLSMultiCerts(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
				TLSKeyFile:                    "./../../api/test-fixtures/keys/cert.pem",
				TLSMinVersion:                 "tls10",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(err.Error(), "found a certificate rather than a key in the PEM for the private key") {
		t.Fatalf("Bad error message: %s", err)
	}
}

// TestTLSInvalidRoot makes sure that the Verify call in tls.go checks the authority of
// the root. The root certificate used in this test is the Baltimore Cyber Trust root
// certificate, downloaded from: https://www.digicert.com/kb/digicert-root-certificates.htm
func TestTLSInvalidRoot(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./test-fixtures/goodcertbadroot.pem",
				TLSKeyFile:                    "./test-fixtures/goodkey.pem",
				TLSMinVersion:                 "tls10",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if err.Error() != "failed to verify certificate: x509: certificate signed by unknown authority" {
		t.Fatalf("Bad error message: %s", err)
	}
}

// TestTLSNoRoot ensures that a server certificate that is passed in without a root
// is still accepted by diagnose as valid. This is an acceptable, though less secure,
// server configuration.
func TestTLSNoRoot(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
				TLSKeyFile:                    "./test-fixtures/goodkey.pem",
				TLSMinVersion:                 "tls10",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err != nil {
		t.Fatalf("Server certificate without root certificate is insecure, but still valid.")
	}
}

// TestTLSInvalidMinVersion checks that a listener with an invalid minimum configured
// version errors appropriately.
func TestTLSInvalidMinVersion(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
				TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
				TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
				TLSMinVersion:                 "0",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if err.Error() != fmt.Errorf(minVersionError, "0").Error() {
		t.Fatalf("Bad error message: %s", err)
	}
}

// TestTLSInvalidMaxVersion checks that a listener with an invalid maximum configured
// version errors appropriately.
func TestTLSInvalidMaxVersion(t *testing.T) {
	listeners := []listenerutil.Listener{
		{
			Config: &configutil.Listener{
				Type:                          "tcp",
				Address:                       "127.0.0.1:443",
				ClusterAddress:                "127.0.0.1:8201",
				TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
				TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
				TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
				TLSMaxVersion:                 "0",
				TLSRequireAndVerifyClientCert: true,
				TLSDisableClientCerts:         false,
			},
		},
	}
	err := ListenerChecks(listeners)
	if err == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if err.Error() != fmt.Errorf(maxVersionError, "0").Error() {
		t.Errorf("Bad error message: %w", err)
	}
}

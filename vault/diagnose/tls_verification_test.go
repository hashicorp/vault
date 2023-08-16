// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package diagnose

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
)

// TestTLSValidCert is the positive test case to show that specifying a valid cert and key
// passes all checks.
func TestTLSValidCert(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./test-fixtures/goodcertwithroot.pem",
			TLSKeyFile:            "./test-fixtures/goodkey.pem",
			TLSMinVersion:         "tls10",
			TLSDisableClientCerts: true,
		},
	}
	warnings, errs := ListenerChecks(context.Background(), listeners)
	if errs != nil {
		// The test failed -- we can just return one of the errors
		t.Fatalf(errs[0].Error())
	}
	if warnings != nil {
		t.Fatalf("warnings returned from good listener")
	}
}

// TestTLSFakeCert simply ensures that the certificate file must contain PEM data.
func TestTLSFakeCert(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./test-fixtures/fakecert.pem",
			TLSKeyFile:            "./../../api/test-fixtures/keys/key.pem",
			TLSMinVersion:         "tls10",
			TLSDisableClientCerts: true,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if len(errs) != 1 {
		t.Fatalf("more than one error returned: %+v", errs)
	}
	if !strings.Contains(errs[0].Error(), "Could not decode certificate") {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

// TestTLSTrailingData uses a certificate from:
// https://github.com/golang/go/issues/40545 that contains
// an extra DER sequence, and makes sure a trailing data error
// is returned.
func TestTLSTrailingData(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./test-fixtures/trailingdatacert.pem",
			TLSKeyFile:            "./../../api/test-fixtures/keys/key.pem",
			TLSMinVersion:         "tls10",
			TLSDisableClientCerts: true,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(errs[0].Error(), "x509: trailing data") {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

// TestTLSExpiredCert checks that an expired certificate fails TLS checks
// with an appropriate error.
func TestTLSExpiredCert(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./test-fixtures/expiredcert.pem",
			TLSKeyFile:            "./test-fixtures/expiredprivatekey.pem",
			TLSMinVersion:         "tls10",
			TLSDisableClientCerts: true,
		},
	}
	warnings, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(errs[0].Error(), "certificate has expired or is not yet valid") {
		t.Fatalf("Bad error message: %s", errs[0])
	}
	if warnings == nil || len(warnings) != 1 {
		t.Fatalf("TLS Config check on fake certificate should warn")
	}
	if !strings.Contains(warnings[0], "expired or near expiry") {
		t.Fatalf("Bad warning: %s", warnings[0])
	}
}

// TestTLSMismatchedCryptographicInfo verifies that a cert and key of differing cryptographic
// types, when specified together, is met with a unique error message.
func TestTLSMismatchedCryptographicInfo(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:            "./test-fixtures/ecdsa.key",
			TLSMinVersion:         "tls10",
			TLSDisableClientCerts: true,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(errs[0].Error(), "tls: private key type does not match public key type") {
		t.Fatalf("Bad error message: %s", errs[0])
	}

	listeners = []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./test-fixtures/ecdsa.crt",
			TLSKeyFile:            "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:       "./../../api/test-fixtures/root/rootcacert.pem",
			TLSMinVersion:         "tls10",
			TLSDisableClientCerts: true,
		},
	}
	_, errs = ListenerChecks(context.Background(), listeners)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(errs[0].Error(), "tls: private key type does not match public key type") {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

// TestTLSMultiKeys verifies that a unique error message is thrown when a key is specified twice.
func TestTLSMultiKeys(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./../../api/test-fixtures/keys/key.pem",
			TLSKeyFile:            "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:       "./../../api/test-fixtures/root/rootcacert.pem",
			TLSMinVersion:         "tls10",
			TLSDisableClientCerts: true,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(errs[0].Error(), "PEM block does not parse to a certificate") {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

// TestTLSCertAsKey verifies that a unique error message is thrown when a cert is specified twice.
func TestTLSCertAsKey(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:            "./../../api/test-fixtures/keys/cert.pem",
			TLSMinVersion:         "tls10",
			TLSDisableClientCerts: true,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(errs[0].Error(), "found a certificate rather than a key in the PEM for the private key") {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

// TestTLSInvalidRoot makes sure that the Verify call in tls.go checks the authority of
// the root. The root certificate used in this test is the Baltimore Cyber Trust root
// certificate, downloaded from: https://www.digicert.com/kb/digicert-root-certificates.htm
func TestTLSInvalidRoot(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./test-fixtures/goodcertbadroot.pem",
			TLSKeyFile:            "./test-fixtures/goodkey.pem",
			TLSMinVersion:         "tls10",
			TLSDisableClientCerts: true,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on fake certificate should fail")
	}
	if !strings.Contains(errs[0].Error(), "certificate signed by unknown authority") {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

// TestTLSNoRoot ensures that a server certificate that is passed in without a root
// is still accepted by diagnose as valid. This is an acceptable, though less secure,
// server configuration.
func TestTLSNoRoot(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:            "./test-fixtures/goodkey.pem",
			TLSMinVersion:         "tls10",
			TLSDisableClientCerts: true,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)

	if errs != nil {
		t.Fatalf("server certificate without root certificate is insecure, but still valid")
	}
}

// TestTLSInvalidMinVersion checks that a listener with an invalid minimum configured
// version errors appropriately.
func TestTLSInvalidMinVersion(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:            "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:       "./../../api/test-fixtures/root/rootcacert.pem",
			TLSMinVersion:         "0",
			TLSDisableClientCerts: true,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on invalid 'tls_min_version' should fail")
	}
	if !strings.Contains(errs[0].Error(), fmt.Errorf(minVersionError, "0").Error()) {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

// TestTLSInvalidMaxVersion checks that a listener with an invalid maximum configured
// version errors appropriately.
func TestTLSInvalidMaxVersion(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:            "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:       "./../../api/test-fixtures/root/rootcacert.pem",
			TLSMaxVersion:         "0",
			TLSDisableClientCerts: true,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on invalid 'tls_max_version' should fail")
	}
	if !strings.Contains(errs[0].Error(), fmt.Errorf(maxVersionError, "0").Error()) {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

// TestDisabledClientCertsAndDisabledTLSClientCAVerfiy checks that a listener works properly when both
// TLSRequireAndVerifyClientCert and TLSDisableClientCerts are false
func TestDisabledClientCertsAndDisabledTLSClientCAVerfiy(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
			TLSMaxVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: false,
			TLSDisableClientCerts:         false,
		},
	}
	status, _ := TLSMutualExclusionCertCheck(listeners[0])
	if status != 0 {
		t.Fatalf("TLS config failed when both TLSRequireAndVerifyClientCert and TLSDisableClientCerts are false")
	}
}

// TestTLSClientCAVerfiy checks that a listener which has TLS client certs checks enabled works as expected
func TestTLSClientCAVerfiy(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
			TLSMaxVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: true,
			TLSDisableClientCerts:         false,
		},
	}
	status, err := TLSMutualExclusionCertCheck(listeners[0])
	if status != 0 {
		t.Fatalf("TLS config check failed with %s", err)
	}
}

// TestTLSClientCAVerfiySkip checks that TLS client cert checks are skipped if TLSDisableClientCerts is true
// regardless of the value for TLSRequireAndVerifyClientCert
func TestTLSClientCAVerfiySkip(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
			TLSMaxVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: false,
			TLSDisableClientCerts:         true,
		},
	}
	status, err := TLSMutualExclusionCertCheck(listeners[0])
	if status != 0 {
		t.Fatalf("TLS config check did not skip verification and failed with %s", err)
	}
}

// TestTLSClientCAVerfiyMutualExclusion checks that TLS client cert checks are skipped if TLSDisableClientCerts is true
// regardless of the value for TLSRequireAndVerifyClientCert
func TestTLSClientCAVerfiyMutualExclusion(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
			TLSMaxVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: true,
			TLSDisableClientCerts:         true,
		},
	}
	status, err := TLSMutualExclusionCertCheck(listeners[0])
	if status == 0 {
		t.Fatalf("TLS config check should have failed when both 'tls_disable_client_certs' and 'tls_require_and_verify_client_cert' are true")
	}
	if !strings.Contains(err, "The tls_disable_client_certs and tls_require_and_verify_client_cert fields in the "+
		"listener stanza of the Vault server configuration are mutually exclusive fields") {
		t.Fatalf("Bad error message: %s", err)
	}
}

// TestTLSClientCAVerfiy checks that a listener which has TLS client certs checks enabled works as expected
func TestTLSClientCAFileCheck(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:               "./../../api/test-fixtures/root/rootcacert.pem",
			TLSMaxVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: true,
			TLSDisableClientCerts:         false,
		},
	}
	warnings, errs := ListenerChecks(context.Background(), listeners)
	if errs != nil {
		t.Fatalf("TLS config check failed while a good ClientCAFile was used")
	}
	if warnings != nil {
		t.Fatalf("TLS config check return warning while a good ClientCAFile was used")
	}
}

// TestTLSLeafCertInClientCAFile checks if a leafCert exist in TLSClientCAFile
func TestTLSLeafCertInClientCAFile(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:               "./test-fixtures/goodcertbadroot.pem",
			TLSMaxVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: true,
			TLSDisableClientCerts:         false,
		},
	}
	warnings, errs := ListenerChecks(context.Background(), listeners)
	fmt.Println(warnings)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on bad ClientCAFile certificate should fail once")
	}
	if warnings == nil || len(warnings) != 1 {
		t.Fatalf("TLS Config check on bad ClientCAFile certificate should warn once")
	}
	if !strings.Contains(warnings[0], "Found at least one leaf certificate in the CA certificate file.") {
		t.Fatalf("Bad error message: %s", warnings[0])
	}
	if !strings.Contains(errs[0].Error(), "signed by unknown authority") {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

// TestTLSNoRootInClientCAFile checks if no Root cert exist in TLSClientCAFile
func TestTLSNoRootInClientCAFile(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:               "./test-fixtures/intermediateCert.pem",
			TLSMaxVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: true,
			TLSDisableClientCerts:         false,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil {
		t.Fatalf("TLS Config check on bad ClientCAFile certificate should fail")
	}
	if !strings.Contains(errs[0].Error(), " No root certificate found") {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

// TestTLSIntermediateCertInClientCAFile checks if an intermediate cert is included in TLSClientCAFile
func TestTLSIntermediateCertInClientCAFile(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:               "./test-fixtures/chain.crt.pem",
			TLSMaxVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: true,
			TLSDisableClientCerts:         false,
		},
	}
	warnings, _ := ListenerChecks(context.Background(), listeners)
	if warnings == nil || len(warnings) != 1 {
		t.Fatalf("TLS Config check on bad ClientCAFile certificate should fail")
	}
	if !strings.Contains(warnings[0], "Found at least one intermediate certificate in the CA certificate file.") {
		t.Fatalf("Bad error message: %s", warnings[0])
	}
}

// TestTLSMultipleRootInClietCACert checks if multiple roots included in TLSClientCAFile
func TestTLSMultipleRootInClietCACert(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:               "./test-fixtures/twoRootCA.pem",
			TLSMinVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: true,
			TLSDisableClientCerts:         false,
		},
	}
	warnings, errs := ListenerChecks(context.Background(), listeners)
	if errs != nil {
		t.Fatalf("TLS Config check on valid certificate should not fail")
	}
	if warnings == nil {
		t.Fatalf("TLS Config check on valid but bad certificate should warn")
	}
	if !strings.Contains(warnings[0], "Found multiple root certificates in CA Certificate file instead of just one.") {
		t.Fatalf("Bad warning: %s", warnings[0])
	}
}

// TestTLSSelfSignedCerts tests invalid self-signed cert as TLSClientCAFile
func TestTLSSelfSignedCert(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   "./../../api/test-fixtures/keys/cert.pem",
			TLSKeyFile:                    "./../../api/test-fixtures/keys/key.pem",
			TLSClientCAFile:               "test-fixtures/selfSignedCert.pem",
			TLSMinVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: true,
			TLSDisableClientCerts:         false,
		},
	}
	_, errs := ListenerChecks(context.Background(), listeners)
	if errs == nil {
		t.Fatalf("Self-signed certificate is insecure")
	}
	if !strings.Contains(errs[0].Error(), "No root certificate found") {
		t.Fatalf("Bad error message: %s", errs[0])
	}
}

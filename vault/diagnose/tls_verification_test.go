// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package diagnose

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	mathrand2 "math/rand/v2"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/internalshared/configutil"
)

type testLeafWithRoot struct {
	testCa             generatedCert
	leafCertFile       string
	leafCertPem        *pem.Block
	leafKeyFile        string
	combinedLeafCaFile string
}

type generatedCert struct {
	keyFile  string
	certFile string
	certPem  *pem.Block
	cert     *x509.Certificate
	key      *ecdsa.PrivateKey
}

type testLeafWithIntermediary struct {
	rootCa         generatedCert
	intCa          generatedCert
	leaf           generatedCert
	combinedCaFile string
}

// generateCertWithIntermediaryRoot generates a leaf certificate signed by an intermediary root CA
func generateCertWithIntermediaryRoot(t testing.TB) testLeafWithIntermediary {
	t.Helper()
	tempDir := t.TempDir()
	template := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		SerialNumber: big.NewInt(mathrand2.Int64()),
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(60 * 24 * time.Hour),
	}

	ca := generateRootCa(t)
	caIntTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "Intermediary CA",
		},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		SerialNumber:          big.NewInt(mathrand2.Int64()),
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caInt := generateCertAndSign(t, caIntTemplate, ca, tempDir, "int_")
	leafCert := generateCertAndSign(t, template, caInt, tempDir, "leaf_")

	combinedCasFile := filepath.Join(tempDir, "cas.pem")
	err := os.WriteFile(combinedCasFile, append(pem.EncodeToMemory(caInt.certPem), pem.EncodeToMemory(ca.certPem)...), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	return testLeafWithIntermediary{
		rootCa:         ca,
		intCa:          caInt,
		leaf:           leafCert,
		combinedCaFile: combinedCasFile,
	}
}

// generateCertAndSign generates a certificate and associated key signed by a CA
func generateCertAndSign(t testing.TB, template *x509.Certificate, ca generatedCert, tempDir string, filePrefix string) generatedCert {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, template, ca.cert, key.Public(), ca.key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		t.Fatal(err)
	}
	certPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}
	certFile := filepath.Join(tempDir, filePrefix+"cert.pem")
	err = os.WriteFile(certFile, pem.EncodeToMemory(certPEMBlock), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	marshaledKey, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	keyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledKey,
	}
	keyFile := filepath.Join(tempDir, filePrefix+"key.pem")
	err = os.WriteFile(keyFile, pem.EncodeToMemory(keyPEMBlock), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	return generatedCert{
		keyFile:  keyFile,
		certFile: certFile,
		certPem:  certPEMBlock,
		cert:     cert,
		key:      key,
	}
}

// generateCertWithRoot generates a leaf certificate signed by a root CA
func generateCertWithRoot(t testing.TB) testLeafWithRoot {
	t.Helper()
	tempDir := t.TempDir()
	leafTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		SerialNumber: big.NewInt(mathrand2.Int64()),
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		NotBefore:    time.Now().Add(-30 * time.Second),
		NotAfter:     time.Now().Add(60 * 24 * time.Hour),
	}

	ca := generateRootCa(t)
	leafCert := generateCertAndSign(t, leafTemplate, ca, tempDir, "leaf_")

	combinedCaLeafFile := filepath.Join(tempDir, "leaf-ca.pem")
	err := os.WriteFile(combinedCaLeafFile, append(pem.EncodeToMemory(leafCert.certPem), pem.EncodeToMemory(ca.certPem)...), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	return testLeafWithRoot{
		testCa:             ca,
		leafCertPem:        leafCert.certPem,
		leafCertFile:       leafCert.certFile,
		leafKeyFile:        leafCert.keyFile,
		combinedLeafCaFile: combinedCaLeafFile,
	}
}

// generateRootCa generates a self-signed root CA certificate and key
func generateRootCa(t testing.TB) generatedCert {
	t.Helper()
	tempDir := t.TempDir()

	caCertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "Root CA",
		},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		SerialNumber:          big.NewInt(mathrand2.Int64()),
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, caKey.Public(), caKey)
	if err != nil {
		t.Fatal(err)
	}
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		t.Fatal(err)
	}
	caCertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	caFile := filepath.Join(tempDir, "ca_root_cert.pem")
	err = os.WriteFile(caFile, pem.EncodeToMemory(caCertPEMBlock), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	marshaledCAKey, err := x509.MarshalECPrivateKey(caKey)
	if err != nil {
		t.Fatal(err)
	}
	caKeyPEMBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: marshaledCAKey,
	}
	caKeyFile := filepath.Join(tempDir, "ca_root_key.pem")
	err = os.WriteFile(caKeyFile, pem.EncodeToMemory(caKeyPEMBlock), 0o644)
	if err != nil {
		t.Fatal(err)
	}
	return generatedCert{
		certPem:  caCertPEMBlock,
		certFile: caFile,
		keyFile:  caKeyFile,
		cert:     caCert,
		key:      caKey,
	}
}

// TestTLSValidCert is the positive test case to show that specifying a valid cert and key
// passes all checks.
func TestTLSValidCert(t *testing.T) {
	tlsFiles := generateCertWithRoot(t)

	listeners := []*configutil.Listener{
		{
			Type:                  "tcp",
			Address:               "127.0.0.1:443",
			ClusterAddress:        "127.0.0.1:8201",
			TLSCertFile:           tlsFiles.combinedLeafCaFile,
			TLSKeyFile:            tlsFiles.leafKeyFile,
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
	testCaFiles := generateCertWithRoot(t)
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   testCaFiles.leafCertFile,
			TLSKeyFile:                    testCaFiles.leafKeyFile,
			TLSClientCAFile:               testCaFiles.testCa.certFile,
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
	testCaFiles := generateCertWithRoot(t)

	tempDir := t.TempDir()

	otherRoot := generateRootCa(t)
	mixedLeafWithRoot := filepath.Join(tempDir, "goodcertbadroot.pem")
	err := os.WriteFile(mixedLeafWithRoot, append(pem.EncodeToMemory(testCaFiles.leafCertPem), pem.EncodeToMemory(otherRoot.certPem)...), 0o644)
	if err != nil {
		t.Fatalf("Failed to write file %s: %v", mixedLeafWithRoot, err)
	}

	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   testCaFiles.combinedLeafCaFile,
			TLSKeyFile:                    testCaFiles.leafKeyFile,
			TLSClientCAFile:               mixedLeafWithRoot,
			TLSMaxVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: true,
			TLSDisableClientCerts:         false,
		},
	}
	warnings, errs := ListenerChecks(context.Background(), listeners)
	fmt.Println(warnings)
	if errs == nil || len(errs) != 1 {
		t.Fatalf("TLS Config check on bad ClientCAFile certificate should fail once: got %v", errs)
	}
	if warnings == nil || len(warnings) != 1 {
		t.Fatalf("TLS Config check on bad ClientCAFile certificate should warn once: got %v", warnings)
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
	testCa := generateCertWithIntermediaryRoot(t)
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   testCa.leaf.certFile,
			TLSKeyFile:                    testCa.leaf.keyFile,
			TLSClientCAFile:               testCa.intCa.certFile,
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
	testCa := generateCertWithIntermediaryRoot(t)
	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   testCa.leaf.certFile,
			TLSKeyFile:                    testCa.leaf.keyFile,
			TLSClientCAFile:               testCa.combinedCaFile,
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

// TestTLSMultipleRootInClientCACert checks if multiple roots included in TLSClientCAFile
func TestTLSMultipleRootInClientCACert(t *testing.T) {
	testCa := generateCertWithRoot(t)
	otherRoot := generateRootCa(t)
	tempDir := t.TempDir()
	mixedRoots := filepath.Join(tempDir, "twoRootCA.pem")
	err := os.WriteFile(mixedRoots, append(pem.EncodeToMemory(testCa.testCa.certPem), pem.EncodeToMemory(otherRoot.certPem)...), 0o644)
	if err != nil {
		t.Fatalf("Failed to write file %s: %v", mixedRoots, err)
	}

	listeners := []*configutil.Listener{
		{
			Type:                          "tcp",
			Address:                       "127.0.0.1:443",
			ClusterAddress:                "127.0.0.1:8201",
			TLSCertFile:                   testCa.leafCertFile,
			TLSKeyFile:                    testCa.leafKeyFile,
			TLSClientCAFile:               mixedRoots,
			TLSMinVersion:                 "tls10",
			TLSRequireAndVerifyClientCert: true,
			TLSDisableClientCerts:         false,
		},
	}
	warnings, errs := ListenerChecks(context.Background(), listeners)
	if errs != nil {
		t.Fatalf("TLS Config check on valid certificate should not fail got: %v", errs)
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

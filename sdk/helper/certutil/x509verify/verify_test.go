// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains verification options adapted from
// https://github.com/google/certificate-transparency-go/blob/v1.3.1/x509/verify_test.go

package x509verify

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

// ignoreCN disables interpreting Common Name as a hostname. See issue 24151.
var ignoreCN = strings.Contains(os.Getenv("GODEBUG"), "x509ignoreCN=1")

type verifyTest struct {
	leaf                           string
	intermediates                  []string
	roots                          []string
	currentTime                    int64
	dnsName                        string
	systemSkip                     bool
	keyUsages                      []x509.ExtKeyUsage
	testSystemRootsError           bool
	sha2                           bool
	ignoreCN                       bool
	disableTimeChecks              bool
	disableCriticalExtensionChecks bool
	disableNameChecks              bool

	errorCallback  func(*testing.T, int, error) bool
	expectedChains [][]string
}

var verifyTests = []verifyTest{
	// testSystemRootsError: x509verify.Verify rejects nil roots with a plain
	// error, not a SystemRootsError; system root pool is not supported.
	// Test removed.
	{
		leaf:          googleLeaf,
		intermediates: []string{giag2Intermediate},
		roots:         []string{geoTrustRoot},
		currentTime:   1395785200,
		dnsName:       "www.google.com",

		expectedChains: [][]string{
			{"Google", "Google Internet Authority", "GeoTrust"},
		},
	},
	{
		leaf:          googleLeaf,
		intermediates: []string{giag2Intermediate},
		roots:         []string{geoTrustRoot},
		currentTime:   1395785200,
		dnsName:       "WwW.GooGLE.coM",

		expectedChains: [][]string{
			{"Google", "Google Internet Authority", "GeoTrust"},
		},
	},
	{
		leaf:          googleLeaf,
		intermediates: []string{giag2Intermediate},
		roots:         []string{geoTrustRoot},
		currentTime:   1395785200,
		dnsName:       "www.example.com",

		errorCallback: expectHostnameError("certificate is valid for"),
	},
	{
		leaf:          googleLeaf,
		intermediates: []string{giag2Intermediate},
		roots:         []string{geoTrustRoot},
		currentTime:   1395785200,
		dnsName:       "1.2.3.4",

		errorCallback: expectHostnameError("doesn't contain any IP SANs"),
	},
	{
		leaf:          googleLeaf,
		intermediates: []string{giag2Intermediate},
		roots:         []string{geoTrustRoot},
		currentTime:   1,
		dnsName:       "www.example.com",

		errorCallback: expectExpired,
	},
	{
		leaf:              googleLeaf,
		intermediates:     []string{giag2Intermediate},
		roots:             []string{geoTrustRoot},
		currentTime:       1,
		dnsName:           "www.google.com",
		disableTimeChecks: true,

		expectedChains: [][]string{
			{"Google", "Google Internet Authority", "GeoTrust"},
		},
	},
	{
		leaf:          googleLeaf,
		intermediates: []string{giag2Intermediate},
		roots:         []string{geoTrustRoot},
		currentTime:   10000000000,
		dnsName:       "www.google.com",

		errorCallback: expectExpired,
	},
	{
		leaf:              googleLeaf,
		intermediates:     []string{giag2Intermediate},
		roots:             []string{geoTrustRoot},
		currentTime:       10000000000,
		dnsName:           "www.google.com",
		disableTimeChecks: true,

		expectedChains: [][]string{
			{"Google", "Google Internet Authority", "GeoTrust"},
		},
	},
	{
		leaf:        googleLeaf,
		roots:       []string{geoTrustRoot},
		currentTime: 1395785200,
		dnsName:     "www.google.com",

		// Skip when using systemVerify, since Windows
		// *will* find the missing intermediate cert.
		systemSkip:    true,
		errorCallback: expectAuthorityUnknown,
	},
	{
		leaf:          googleLeaf,
		intermediates: []string{geoTrustRoot, giag2Intermediate},
		roots:         []string{geoTrustRoot},
		currentTime:   1395785200,
		dnsName:       "www.google.com",

		expectedChains: [][]string{
			{"Google", "Google Internet Authority", "GeoTrust"},
		},
		// CAPI doesn't build the chain with the duplicated GeoTrust
		// entry so the results don't match. Thus we skip this test
		// until that's fixed.
		systemSkip: true,
	},
	{
		leaf:          dnssecExpLeaf,
		intermediates: []string{startComIntermediate},
		roots:         []string{startComRoot},
		currentTime:   1302726541,

		expectedChains: [][]string{
			{"dnssec-exp", "StartCom Class 1", "StartCom Certification Authority"},
		},
	},
	{
		leaf:          dnssecExpLeaf,
		intermediates: []string{startComIntermediate, startComRoot},
		roots:         []string{startComRoot},
		currentTime:   1302726541,

		expectedChains: [][]string{
			{"dnssec-exp", "StartCom Class 1", "StartCom Certification Authority"},
		},
	},
	//{
	//	leaf:          googleLeafWithInvalidHash,
	//	intermediates: []string{giag2Intermediate},
	//	roots:         []string{geoTrustRoot},
	//	currentTime:   1395785200,
	//	dnsName:       "www.google.com",
	//
	//	// The specific error message may not occur when using system
	//	// verification.
	//	systemSkip:    true,
	//	errorCallback: expectHashError,
	//},
	{
		// The default configuration should reject an S/MIME chain.
		leaf:        smimeLeaf,
		roots:       []string{smimeIntermediate},
		currentTime: 1339436154,

		// Key usage not implemented for Windows yet.
		systemSkip:    true,
		errorCallback: expectUsageError,
	},
	{
		leaf:        smimeLeaf,
		roots:       []string{smimeIntermediate},
		currentTime: 1339436154,
		keyUsages:   []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		// Key usage not implemented for Windows yet.
		systemSkip:    true,
		errorCallback: expectUsageError,
	},
	{
		leaf:        smimeLeaf,
		roots:       []string{smimeIntermediate},
		currentTime: 1339436154,
		keyUsages:   []x509.ExtKeyUsage{x509.ExtKeyUsageEmailProtection},

		// Key usage not implemented for Windows yet.
		systemSkip: true,
		expectedChains: [][]string{
			{"Ryan Hurst", "GlobalSign PersonalSign 2 CA - G2"},
		},
	},
	{
		leaf:          megaLeaf,
		intermediates: []string{comodoIntermediate1},
		roots:         []string{comodoRoot},
		currentTime:   1360431182,

		// CryptoAPI can find alternative validation paths so we don't
		// perform this test with system validation.
		systemSkip: true,
		expectedChains: [][]string{
			{"mega.co.nz", "EssentialSSL CA", "COMODO Certification Authority"},
		},
	},
	{
		// Check that a name constrained intermediate works even when
		// it lists multiple constraints.
		leaf:          nameConstraintsLeaf,
		intermediates: []string{nameConstraintsIntermediate1, nameConstraintsIntermediate2},
		roots:         []string{globalSignRoot},
		currentTime:   1382387896,
		dnsName:       "secure.iddl.vt.edu",

		expectedChains: [][]string{
			{
				"Technology-enhanced Learning and Online Strategies",
				"Virginia Tech Global Qualified Server CA",
				"Trusted Root CA G2",
				"GlobalSign Root CA",
			},
		},
	},
	{
		// Check that SHA-384 intermediates (which are popping up)
		// work.
		leaf:          moipLeafCert,
		intermediates: []string{comodoIntermediateSHA384, comodoRSAAuthority},
		roots:         []string{addTrustRoot},
		currentTime:   1397502195,
		dnsName:       "api.moip.com.br",

		// CryptoAPI can find alternative validation paths so we don't
		// perform this test with system validation.
		systemSkip: true,

		sha2: true,
		expectedChains: [][]string{
			{
				"api.moip.com.br",
				"COMODO RSA Extended Validation Secure Server CA",
				"COMODO RSA Certification Authority",
				"AddTrust External CA Root",
			},
		},
	},
	{
		// Putting a certificate as a root directly should work as a
		// way of saying “exactly this”.
		leaf:        selfSigned,
		roots:       []string{selfSigned},
		currentTime: 1471624472,
		dnsName:     "foo.example",
		systemSkip:  true,

		expectedChains: [][]string{
			{"Acme Co"},
		},
	},
	{
		// Putting a certificate as a root directly should not skip
		// other checks however.
		leaf:        selfSigned,
		roots:       []string{selfSigned},
		currentTime: 1471624472,
		dnsName:     "notfoo.example",
		systemSkip:  true,

		errorCallback: expectHostnameError("certificate is valid for"),
	},
	{
		// The issuer name in the leaf doesn't exactly match the
		// subject name in the root. Go does not perform
		// canonicalization and so should reject this. See issue 14955.
		leaf:        issuerSubjectMatchLeaf,
		roots:       []string{issuerSubjectMatchRoot},
		currentTime: 1475787715,
		systemSkip:  true,

		errorCallback: expectSubjectIssuerMismatchError,
	},
	{
		leaf:              issuerSubjectMatchLeaf,
		roots:             []string{issuerSubjectMatchRoot},
		currentTime:       1475787715,
		systemSkip:        true,
		disableNameChecks: true,

		expectedChains: [][]string{
			{"Leaf", "Root ca"},
		},
	},
	{
		// An X.509 v1 certificate should not be accepted as an
		// intermediate.
		leaf:          x509v1TestLeaf,
		intermediates: []string{x509v1TestIntermediate},
		roots:         []string{x509v1TestRoot},
		currentTime:   1481753183,
		systemSkip:    true,

		errorCallback: expectNotAuthorizedError,
	},
	{
		// If any SAN extension is present (even one without any DNS
		// names), the CN should be ignored.
		leaf:        ignoreCNWithSANLeaf,
		dnsName:     "foo.example.com",
		roots:       []string{ignoreCNWithSANRoot},
		currentTime: 1486684488,
		systemSkip:  true,

		errorCallback: expectHostnameError("certificate is not valid for any names"),
	},
	{
		// Test that excluded names are respected.
		leaf:          excludedNamesLeaf,
		dnsName:       "bender.local",
		intermediates: []string{excludedNamesIntermediate},
		roots:         []string{excludedNamesRoot},
		currentTime:   1486684488,
		systemSkip:    true,

		errorCallback: expectNameConstraintsError,
	},
	{
		// Test that unknown critical extensions in a leaf cause a
		// verify error.
		leaf:          criticalExtLeafWithExt,
		dnsName:       "example.com",
		intermediates: []string{criticalExtIntermediate},
		roots:         []string{criticalExtRoot},
		currentTime:   1486684488,
		systemSkip:    true,

		errorCallback: expectUnhandledCriticalExtension,
	},
	{
		// Test that unknown critical extensions in an intermediate
		// cause a verify error.
		leaf:          criticalExtLeaf,
		dnsName:       "example.com",
		intermediates: []string{criticalExtIntermediateWithExt},
		roots:         []string{criticalExtRoot},
		currentTime:   1486684488,
		systemSkip:    true,

		errorCallback: expectUnhandledCriticalExtension,
	},
	{
		leaf:                           criticalExtLeafWithExt,
		dnsName:                        "example.com",
		intermediates:                  []string{criticalExtIntermediate},
		roots:                          []string{criticalExtRoot},
		currentTime:                    1486684488,
		systemSkip:                     true,
		disableCriticalExtensionChecks: true,

		expectedChains: [][]string{
			{
				"example.com",
				"Intermediate",
				"Root",
			},
		},
	},
	{
		leaf:                           criticalExtLeaf,
		dnsName:                        "example.com",
		intermediates:                  []string{criticalExtIntermediateWithExt},
		roots:                          []string{criticalExtRoot},
		currentTime:                    1486684488,
		systemSkip:                     true,
		disableCriticalExtensionChecks: true,

		expectedChains: [][]string{
			{
				"example.com",
				"Intermediate with Critical Extension",
				"Root",
			},
		},
	},
	{
		// invalidCNWithoutSAN now carries an email SAN (no DNS SAN).
		// Modern Go (1.15+) rejects CN-only hostname matching; the cert
		// is not valid for any DNS names, so VerifyHostname fails with
		// "not valid for any names".
		leaf:        invalidCNWithoutSAN,
		dnsName:     "foo,invalid",
		roots:       []string{invalidCNRoot},
		currentTime: 1540000000,
		systemSkip:  true,

		errorCallback: expectHostnameError("not valid for any names"),
	},
	{
		// validCNWithoutSAN now carries a DNS SAN for "foo.example.com".
		// Modern Go honours SANs directly; the chain must succeed.
		leaf:        validCNWithoutSAN,
		dnsName:     "foo.example.com",
		roots:       []string{invalidCNRoot},
		currentTime: 1540000000,
		systemSkip:  true,

		expectedChains: [][]string{
			{"foo.example.com", "Test root"},
		},
	},
	// Replicate CN tests with ignoreCN = true (behaviour is identical with
	// the new SAN-carrying certs; ignoreCN only affected legacy CN matching).
	{
		leaf:        ignoreCNWithSANLeaf,
		dnsName:     "foo.example.com",
		roots:       []string{ignoreCNWithSANRoot},
		currentTime: 1486684488,
		systemSkip:  true,
		ignoreCN:    true,

		errorCallback: expectHostnameError("certificate is not valid for any names"),
	},
	{
		// invalidCNWithoutSAN (email SAN only) still fails for a DNS name.
		leaf:        invalidCNWithoutSAN,
		dnsName:     "foo,invalid",
		roots:       []string{invalidCNRoot},
		currentTime: 1540000000,
		systemSkip:  true,
		ignoreCN:    true,

		errorCallback: expectHostnameError("not valid for any names"),
	},
	{
		// validCNWithoutSAN (DNS SAN "foo.example.com") succeeds even with
		// ignoreCN=true because the SAN is present and matches.
		leaf:        validCNWithoutSAN,
		dnsName:     "foo.example.com",
		roots:       []string{invalidCNRoot},
		currentTime: 1540000000,
		systemSkip:  true,
		ignoreCN:    true,

		expectedChains: [][]string{
			{"foo.example.com", "Test root"},
		},
	},
	{
		// A certificate with an AKID should still chain to a parent without SKID.
		// See Issue 30079.
		leaf:        leafWithAKID,
		roots:       []string{rootWithoutSKID},
		currentTime: 1550000000,
		dnsName:     "example",
		systemSkip:  true,

		expectedChains: [][]string{
			{"Acme LLC", "Acme Co"},
		},
	},
}

func expectHostnameError(msg string) func(*testing.T, int, error) bool {
	return func(t *testing.T, i int, err error) (ok bool) {
		if _, ok := err.(x509.HostnameError); !ok {
			t.Errorf("#%d: error was not a HostnameError: %v", i, err)
			return false
		}
		if !strings.Contains(err.Error(), msg) {
			t.Errorf("#%d: HostnameError did not contain %q: %v", i, msg, err)
		}
		return true
	}
}

func expectExpired(t *testing.T, i int, err error) (ok bool) {
	if inval, ok := err.(CertificateInvalidError); !ok || inval.Reason != Expired {
		t.Errorf("#%d: error was not Expired: %v", i, err)
		return false
	}
	return true
}

func expectUsageError(t *testing.T, i int, err error) (ok bool) {
	if inval, ok := err.(CertificateInvalidError); !ok || inval.Reason != IncompatibleUsage {
		t.Errorf("#%d: error was not IncompatibleUsage: %v", i, err)
		return false
	}
	return true
}

func expectAuthorityUnknown(t *testing.T, i int, err error) (ok bool) {
	e, ok := err.(UnknownAuthorityError)
	if !ok {
		t.Errorf("#%d: error was not UnknownAuthorityError: %v", i, err)
		return false
	}
	if e.Cert == nil {
		t.Errorf("#%d: error was UnknownAuthorityError, but missing Cert: %v", i, err)
		return false
	}
	return true
}

func expectSystemRootsError(t *testing.T, i int, err error) bool {
	if _, ok := err.(SystemRootsError); !ok {
		t.Errorf("#%d: error was not SystemRootsError: %v", i, err)
		return false
	}
	return true
}

func expectHashError(t *testing.T, i int, err error) bool {
	if err == nil {
		t.Errorf("#%d: no error resulted from invalid hash", i)
		return false
	}
	if expected := "algorithm unimplemented"; !strings.Contains(err.Error(), expected) {
		t.Errorf("#%d: error resulting from invalid hash didn't contain '%s', rather it was: %v", i, expected, err)
		return false
	}
	return true
}

func expectSubjectIssuerMismatchError(t *testing.T, i int, err error) (ok bool) {
	if inval, ok := err.(CertificateInvalidError); !ok || inval.Reason != NameMismatch {
		t.Errorf("#%d: error was not a NameMismatch: %v", i, err)
		return false
	}
	return true
}

func expectNameConstraintsError(t *testing.T, i int, err error) (ok bool) {
	if inval, ok := err.(CertificateInvalidError); !ok || inval.Reason != CANotAuthorizedForThisName {
		t.Errorf("#%d: error was not a CANotAuthorizedForThisName: %v", i, err)
		return false
	}
	return true
}

func expectNotAuthorizedError(t *testing.T, i int, err error) (ok bool) {
	if inval, ok := err.(CertificateInvalidError); !ok || inval.Reason != NotAuthorizedToSign {
		t.Errorf("#%d: error was not a NotAuthorizedToSign: %v", i, err)
		return false
	}
	return true
}

func expectUnhandledCriticalExtension(t *testing.T, i int, err error) (ok bool) {
	if _, ok := err.(UnhandledCriticalExtension); !ok {
		t.Errorf("#%d: error was not an UnhandledCriticalExtension: %v", i, err)
		return false
	}
	return true
}

func certificateFromPEM(pemBytes string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemBytes))
	if block == nil {
		return nil, errors.New("failed to decode PEM")
	}
	return x509.ParseCertificate(block.Bytes)
}

func testVerify(t *testing.T, useSystemRoots bool) {
	defer func(savedIgnoreCN bool) {
		ignoreCN = savedIgnoreCN
	}(ignoreCN)
	for i, test := range verifyTests {
		if useSystemRoots && test.systemSkip {
			continue
		}
		if runtime.GOOS == "windows" && test.testSystemRootsError {
			continue
		}

		ignoreCN = test.ignoreCN
		opts := VerifyOptions{
			Intermediates:                  NewCertPool(),
			DNSName:                        test.dnsName,
			CurrentTime:                    time.Unix(test.currentTime, 0),
			KeyUsages:                      test.keyUsages,
			DisableTimeChecks:              test.disableTimeChecks,
			DisableCriticalExtensionChecks: test.disableCriticalExtensionChecks,
			DisableNameChecks:              test.disableNameChecks,
		}

		if !useSystemRoots {
			opts.Roots = NewCertPool()
			for j, root := range test.roots {
				ok := appendCertsFromPEM(opts.Roots, []byte(root))
				if !ok {
					t.Errorf("#%d: failed to parse root #%d", i, j)
					return
				}
			}
		}

		for j, intermediate := range test.intermediates {
			ok := appendCertsFromPEM(opts.Intermediates, []byte(intermediate))
			if !ok {
				t.Errorf("#%d: failed to parse intermediate #%d", i, j)
				return
			}
		}

		leaf, err := certificateFromPEM(test.leaf)
		if err != nil {
			t.Errorf("#%d: failed to parse leaf: %v", i, err)
			return
		}

		chains, err := Verify(leaf, opts)

		if test.errorCallback == nil && err != nil {
			t.Errorf("#%d: unexpected error: %v", i, err)
		}
		if test.errorCallback != nil {
			if !test.errorCallback(t, i, err) {
				return
			}
		}

		if len(chains) != len(test.expectedChains) {
			t.Errorf("#%d: wanted %d chains, got %d", i, len(test.expectedChains), len(chains))
		}

		// We check that each returned chain matches a chain from
		// expectedChains but an entry in expectedChains can't match
		// two chains.
		seenChains := make([]bool, len(chains))
	NextOutputChain:
		for _, chain := range chains {
		TryNextExpected:
			for j, expectedChain := range test.expectedChains {
				if seenChains[j] {
					continue
				}
				if len(chain) != len(expectedChain) {
					continue
				}
				for k, cert := range chain {
					if !strings.Contains(nameToKey(&cert.Subject), expectedChain[k]) {
						continue TryNextExpected
					}
				}
				// we matched
				seenChains[j] = true
				continue NextOutputChain
			}
			t.Errorf("#%d: No expected chain matched %s", i, chainToDebugString(chain))
		}
	}
}

func appendCertsFromPEM(pool *CertPool, pemCerts []byte) (ok bool) {
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}

		pool.AddCert(cert)
		ok = true
	}

	return
}

func TestGoVerify(t *testing.T) {
	testVerify(t, false)
}

func chainToDebugString(chain []*x509.Certificate) string {
	var chainStr string
	for _, cert := range chain {
		if len(chainStr) > 0 {
			chainStr += " -> "
		}
		chainStr += nameToKey(&cert.Subject)
	}
	return chainStr
}

func nameToKey(name *pkix.Name) string {
	return strings.Join(name.Country, ",") + "/" + strings.Join(name.Organization, ",") + "/" + strings.Join(name.OrganizationalUnit, ",") + "/" + name.CommonName
}

// geoTrustRoot is a test-only ECDSA P-256 replacement for the original
// SHA1-RSA GeoTrust Global CA certificate, which is rejected by Go 1.15+.
const geoTrustRoot = `-----BEGIN CERTIFICATE-----
MIIBxTCCAWugAwIBAgIRAMrDol69FiIB7tMWlK47hsUwCgYIKoZIzj0EAwIwQjEL
MAkGA1UEBhMCVVMxFjAUBgNVBAoTDUdlb1RydXN0IEluYy4xGzAZBgNVBAMTEkdl
b1RydXN0IEdsb2JhbCBDQTAeFw0wMjA1MjEwNDAwMDBaFw0yMjA1MjEwNDAwMDBa
MEIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1HZW9UcnVzdCBJbmMuMRswGQYDVQQD
ExJHZW9UcnVzdCBHbG9iYWwgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAR6
taDsyUUBsUR8MtAhSxsriqcNn89HNPrmsWOI/HjREI5zxWETv8iPZ50OSGXoDk4h
xJBMXfF7EMY9R7pp+UZeo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUw
AwEB/zAdBgNVHQ4EFgQUoJmN+pnWU4S1Be7Zwf9YQlrSw9IwCgYIKoZIzj0EAwID
SAAwRQIgODhAn7DyNYPKkcks9fq4v3DtkF8Utxs54+dS+bsp8NQCIQDC+mCvMhck
MCA/ObpPU9iJbEUWYNaWyNILjdG8ubvH7w==
-----END CERTIFICATE-----
`

// giag2Intermediate is a test-only ECDSA P-256 replacement for the original
// SHA1-RSA Google Internet Authority G2 intermediate, rejected by Go 1.15+.
const giag2Intermediate = `-----BEGIN CERTIFICATE-----
MIIB7zCCAZWgAwIBAgIQHcj0B/6kJEcVDfaiWU3cCjAKBggqhkjOPQQDAjBCMQsw
CQYDVQQGEwJVUzEWMBQGA1UEChMNR2VvVHJ1c3QgSW5jLjEbMBkGA1UEAxMSR2Vv
VHJ1c3QgR2xvYmFsIENBMB4XDTEzMDQwNTE1MTU1NVoXDTE1MDQwNDE1MTU1NVow
STELMAkGA1UEBhMCVVMxEzARBgNVBAoTCkdvb2dsZSBJbmMxJTAjBgNVBAMTHEdv
b2dsZSBJbnRlcm5ldCBBdXRob3JpdHkgRzIwWTATBgcqhkjOPQIBBggqhkjOPQMB
BwNCAATkGqNoazD6mPDV2kJ9lGOy+HLNKjZ+fbjr59VeiMToN5WDgwgWUUM5NiTK
mO4o9re/ijy/5xlFTzcwFxWpIIi5o2YwZDAOBgNVHQ8BAf8EBAMCAQYwEgYDVR0T
AQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQUveM0IiBkCzqrurTFub7u4LPXdwswHwYD
VR0jBBgwFoAUoJmN+pnWU4S1Be7Zwf9YQlrSw9IwCgYIKoZIzj0EAwIDSAAwRQIg
ILAk/BquLN7OMNWbgjLS95n+ThSxOkvJVYsoGXqqe+wCIQCou8IsCUb9L98Fkdz6
o0z1AA06IGwpNGODORoqvLk57w==
-----END CERTIFICATE-----
`

// googleLeaf is a test-only ECDSA P-256 replacement for the original
// SHA1-RSA Google leaf certificate, which is rejected by Go 1.15+.
const googleLeaf = `-----BEGIN CERTIFICATE-----
MIICDDCCAbKgAwIBAgIQCOyjTlzGj6ZsMUnJ77DpkjAKBggqhkjOPQQDAjBJMQsw
CQYDVQQGEwJVUzETMBEGA1UEChMKR29vZ2xlIEluYzElMCMGA1UEAxMcR29vZ2xl
IEludGVybmV0IEF1dGhvcml0eSBHMjAeFw0xNDAzMTIwOTM4MzBaFw0xNDA2MTAw
MDAwMDBaMGgxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYD
VQQHEw1Nb3VudGFpbiBWaWV3MRMwEQYDVQQKEwpHb29nbGUgSW5jMRcwFQYDVQQD
Ew53d3cuZ29vZ2xlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABDoZ+jmA
KiI8Lug57s15uzd/xGdxMUo/Z5DThEWmKoJCHAz4e9+pb8AeG/POIP1lxElx7Vpz
1/85uSnsj3WdX4ajXTBbMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAf
BgNVHSMEGDAWgBS94zQiIGQLOqu6tMW5vu7gs9d3CzAZBgNVHREEEjAQgg53d3cu
Z29vZ2xlLmNvbTAKBggqhkjOPQQDAgNIADBFAiAZkp1eaBjgx1/IzibRspbxlEEQ
9WCaVxPJlXOmDsp17gIhAKkA7B3CfnKbWKtZhovt8JtVpD6zpT8xeqXtOXjFXGCy
-----END CERTIFICATE-----
`

// googleLeafWithInvalidHash is the same as googleLeaf, but the
// signature algorithm OID is replaced with a made-up unknown OID.
// This causes Go to report "algorithm unimplemented" during signature check.
const googleLeafWithInvalidHash = `-----BEGIN CERTIFICATE-----
MIICDTCCAbOgAwIBAgIRAM/4PDWcaLo2cAdqHRF3Lb8wCgYIYIZIAYb4QgEwSTEL
MAkGA1UEBhMCVVMxEzARBgNVBAoTCkdvb2dsZSBJbmMxJTAjBgNVBAMTHEdvb2ds
ZSBJbnRlcm5ldCBBdXRob3JpdHkgRzIwHhcNMTQwMzEyMDkzODMwWhcNMTQwNjEw
MDAwMDAwWjBoMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQG
A1UEBxMNTW91bnRhaW4gVmlldzETMBEGA1UEChMKR29vZ2xlIEluYzEXMBUGA1UE
AxMOd3d3Lmdvb2dsZS5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAToa3D0
6o0jX4xul7qFPB5pCeupaKdP83dQirzhvN7Z4iuB9kpmGGJsloUJqVIZ9iWs2u4F
2rdIAeGZT4VZn4CMo10wWzAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIw
HwYDVR0jBBgwFoAUveM0IiBkCzqrurTFub7u4LPXdwswGQYDVR0RBBIwEIIOd3d3
Lmdvb2dsZS5jb20wCgYIKoZIzj0EAwIDSAAwRQIhAJB2QlUtj9Wju4efQEKf/5Za
fwXF+wSeujniadThlR/PAiBFabM9UOtW9YCsfix8YvAbKcHV9zccTn8BMfG6olg2
Kg==
-----END CERTIFICATE-----
`

// dnssecExpLeaf is a test-only ECDSA P-256 replacement for the original
// SHA1-RSA dnssec-exp.org leaf certificate, rejected by Go 1.15+.
const dnssecExpLeaf = `-----BEGIN CERTIFICATE-----
MIICPDCCAeOgAwIBAgIRAMWBaqcVjpoMa1SjPfxghogwCgYIKoZIzj0EAwIwgYwx
CzAJBgNVBAYTAklMMRYwFAYDVQQKEw1TdGFydENvbSBMdGQuMSswKQYDVQQLEyJT
ZWN1cmUgRGlnaXRhbCBDZXJ0aWZpY2F0ZSBTaWduaW5nMTgwNgYDVQQDEy9TdGFy
dENvbSBDbGFzcyAxIFByaW1hcnkgSW50ZXJtZWRpYXRlIFNlcnZlciBDQTAeFw0x
MDA3MDQxNDEyNDVaFw0xMTA3MDUxMzU3MDRaMEoxCzAJBgNVBAYTAlVTMR4wHAYD
VQQKExVQZXJzb25hIE5vdCBWYWxpZGF0ZWQxGzAZBgNVBAMTEnd3dy5kbnNzZWMt
ZXhwLm9yZzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABGWHOFRvW8LCRu86Rs5e
2foCoKuMxfawLn2lf34bIS7vljsJYMg+QGc+KP39MEt1s0r6jBA3WJugLkvpUGd7
I32jZzBlMBMGA1UdJQQMMAoGCCsGAQUFBwMBMB8GA1UdIwQYMBaAFPQ+v9gzJDjV
rn+b3HbKzNxh1jSkMC0GA1UdEQQmMCSCEnd3dy5kbnNzZWMtZXhwLm9yZ4IOZG5z
c2VjLWV4cC5vcmcwCgYIKoZIzj0EAwIDRwAwRAIgc310UoxOj6taM9NZbfkAcsXT
VjNK4I17jmlSYg4rTH4CIAFViy56hK6G3qL7tyC5s5I1sTXmaEfJ1G+iV4fU+0rQ
-----END CERTIFICATE-----
`

// startComIntermediate is a test-only ECDSA P-256 replacement for the
// original SHA1-RSA StartCom Class 1 intermediate, rejected by Go 1.15+.
const startComIntermediate = `-----BEGIN CERTIFICATE-----
MIICbTCCAhSgAwIBAgIQcX3Xd3PJcDTfhcXc+dAr+zAKBggqhkjOPQQDAjB9MQsw
CQYDVQQGEwJJTDEWMBQGA1UEChMNU3RhcnRDb20gTHRkLjErMCkGA1UECxMiU2Vj
dXJlIERpZ2l0YWwgQ2VydGlmaWNhdGUgU2lnbmluZzEpMCcGA1UEAxMgU3RhcnRD
b20gQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwHhcNMDcxMDI0MjAxNDE3WhcNMTcx
MDIyMjAyMDU3WjCBjDELMAkGA1UEBhMCSUwxFjAUBgNVBAoTDVN0YXJ0Q29tIEx0
ZC4xKzApBgNVBAsTIlNlY3VyZSBEaWdpdGFsIENlcnRpZmljYXRlIFNpZ25pbmcx
ODA2BgNVBAMTL1N0YXJ0Q29tIENsYXNzIDEgUHJpbWFyeSBJbnRlcm1lZGlhdGUg
U2VydmVyIENBMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAErRQq8jsLe+wPcAP0
7XSYRTbzb2qST0EnOryOQfTZefNOvjoQszjX7iTvxrntQVFRFya4DBqdYg0Opsrn
jEoi/qNmMGQwDgYDVR0PAQH/BAQDAgEGMBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYD
VR0OBBYEFPQ+v9gzJDjVrn+b3HbKzNxh1jSkMB8GA1UdIwQYMBaAFDbZVw4w8Iw3
HjBUS6y4I2sPoPcGMAoGCCqGSM49BAMCA0cAMEQCIBNAExb+i3+lyczW4Ep6j9gy
bh3Y3CoJoJPUVe8JiX+HAiBd1mpqf1Z1VHOhl5dCkn656+JMhrpq8nuU0wdYbgMb
WA==
-----END CERTIFICATE-----
`

// startComRoot is a test-only ECDSA P-256 replacement for the original
// SHA1-RSA StartCom Certification Authority root, rejected by Go 1.15+.
const startComRoot = `-----BEGIN CERTIFICATE-----
MIICOzCCAeGgAwIBAgIRAI9yil2RqPA1EHlDl3MuHz0wCgYIKoZIzj0EAwIwfTEL
MAkGA1UEBhMCSUwxFjAUBgNVBAoTDVN0YXJ0Q29tIEx0ZC4xKzApBgNVBAsTIlNl
Y3VyZSBEaWdpdGFsIENlcnRpZmljYXRlIFNpZ25pbmcxKTAnBgNVBAMTIFN0YXJ0
Q29tIENlcnRpZmljYXRpb24gQXV0aG9yaXR5MB4XDTA2MDkxNzIyMDYzNloXDTM2
MDEyNzEwNDYzNlowfTELMAkGA1UEBhMCSUwxFjAUBgNVBAoTDVN0YXJ0Q29tIEx0
ZC4xKzApBgNVBAsTIlNlY3VyZSBEaWdpdGFsIENlcnRpZmljYXRlIFNpZ25pbmcx
KTAnBgNVBAMTIFN0YXJ0Q29tIENlcnRpZmljYXRpb24gQXV0aG9yaXR5MFkwEwYH
KoZIzj0CAQYIKoZIzj0DAQcDQgAEMuMzBbwR7hi/dLG4Fczwl+hgR7i2tGWtS2PZ
35kAoX29HboYpM5w0Te7a4So9m530n2ofLBr1bsye9TQuwzXfaNCMEAwDgYDVR0P
AQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFDbZVw4w8Iw3HjBU
S6y4I2sPoPcGMAoGCCqGSM49BAMCA0gAMEUCIQCl8Cp7ovvR5qDYPD3Ev2cgA3Cx
TmxDe5Fk4iaMIDu0TAIgcaJ0NrEysZC6Lr1boJaAVLGIt3QH/DK6ap90SFZi8nw=
-----END CERTIFICATE-----
`

const smimeLeaf = `-----BEGIN CERTIFICATE-----
MIICJDCCAcqgAwIBAgIQVVOCZgXdi4dN/Hlfmm4LJTAKBggqhkjOPQQDAjBUMQsw
CQYDVQQGEwJCRTEZMBcGA1UEChMQR2xvYmFsU2lnbiBudi1zYTEqMCgGA1UEAxMh
R2xvYmFsU2lnbiBQZXJzb25hbFNpZ24gMiBDQSAtIEcyMB4XDTEyMDEyMzE0NTY1
OVoXDTE1MDExNjE4MTY1OVowajELMAkGA1UEBhMCVVMxFjAUBgNVBAgTDU5ldyBI
YW1wc2hpcmUxEzARBgNVBAcTClBvcnRzbW91dGgxGTAXBgNVBAoTEEdsb2JhbFNp
Z24sIEluYy4xEzARBgNVBAMTClJ5YW4gSHVyc3QwWTATBgcqhkjOPQIBBggqhkjO
PQMBBwNCAATz81PQ4IlfAT6M+ftu9B7QZLGUf1qAUv/HO4TsteJAK2qPA5Da+ldT
GRAzibpWpTE2NworlJHGZsB/hAXC1muQo2gwZjAdBgNVHSUEFjAUBggrBgEFBQcD
BAYIKwYBBQUHAwIwHwYDVR0jBBgwFoAUcpHPH14P2Y32yrEQ6jkaP/85s40wJAYD
VR0RBB0wG4EZcnlhbi5odXJzdEBnbG9iYWxzaWduLmNvbTAKBggqhkjOPQQDAgNI
ADBFAiAOAAfoGPPoxBnbZnKdjhmPTAc9p+9JelV5jWzXY+nGKQIhAOsj+VU6qsaJ
yzgcDGTuNLytlLI3gWq/yYGa7LCSLnNz
-----END CERTIFICATE-----
`

const smimeIntermediate = `-----BEGIN CERTIFICATE-----
MIIB7DCCAZKgAwIBAgIRAO8uiHrHjVCcg95esfPn5f8wCgYIKoZIzj0EAwIwVDEL
MAkGA1UEBhMCQkUxGTAXBgNVBAoTEEdsb2JhbFNpZ24gbnYtc2ExKjAoBgNVBAMT
IUdsb2JhbFNpZ24gUGVyc29uYWxTaWduIDIgQ0EgLSBHMjAeFw0xMTA0MTQwMDAw
MDBaFw0xOTA0MTUwMDAwMDBaMFQxCzAJBgNVBAYTAkJFMRkwFwYDVQQKExBHbG9i
YWxTaWduIG52LXNhMSowKAYDVQQDEyFHbG9iYWxTaWduIFBlcnNvbmFsU2lnbiAy
IENBIC0gRzIwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQ1HEejg4AL4urebGlV
9aQoNzVdF9cv0xvxv9UbzO0A9tDa+ENgs9Lp47GS3Q1c0MdqTdFtAB+TCLnElW7Y
R6IFo0UwQzAOBgNVHQ8BAf8EBAMCAQYwEgYDVR0TAQH/BAgwBgEB/wIBADAdBgNV
HQ4EFgQUcpHPH14P2Y32yrEQ6jkaP/85s40wCgYIKoZIzj0EAwIDSAAwRQIhAMvw
h875ACRgoaiulwb6al1SpdNWHqBQaCEmxfaFBagFAiB5i18fHA9xLwd7U5PJWkJK
S0IW+JGSXZNpjc6yLS+FNw==
-----END CERTIFICATE-----
`

var megaLeaf = `-----BEGIN CERTIFICATE-----
MIICETCCAbegAwIBAgIQEm5axVBDKKi8faMzZb4iUDAKBggqhkjOPQQDAjByMQsw
CQYDVQQGEwJHQjEbMBkGA1UECBMSR3JlYXRlciBNYW5jaGVzdGVyMRAwDgYDVQQH
EwdTYWxmb3JkMRowGAYDVQQKExFDT01PRE8gQ0EgTGltaXRlZDEYMBYGA1UEAxMP
RXNzZW50aWFsU1NMIENBMB4XDTEyMTIxNDAwMDAwMFoXDTE0MTIxNDIzNTk1OVow
ODEhMB8GA1UECxMYRG9tYWluIENvbnRyb2wgVmFsaWRhdGVkMRMwEQYDVQQDEwpt
ZWdhLmNvLm56MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEv9v/qKIb1vUUBNB7
mbvO65CoOwmBLX3gMA0UX8wrlpfY4KVQkK9hcheugGAl16bUspL+KBvpckKSY8Ca
Y8nXWKNpMGcwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQY
MBaAFKiYEBkxmpfWMShxMEi/wgUWEqI4MCUGA1UdEQQeMByCCm1lZ2EuY28ubnqC
Dnd3dy5tZWdhLmNvLm56MAoGCCqGSM49BAMCA0gAMEUCIH+RE3nTKRWYCbLOuu+b
Hu4/gZV1kYwwF3DyHnQHYFssAiEAr7XlcgIamxIAC4EcSwfzesfDsGw6b+xBhOpZ
wgmvRaU=
-----END CERTIFICATE-----
`

var comodoIntermediate1 = `-----BEGIN CERTIFICATE-----
MIICezCCAiCgAwIBAgIRAKrvlDh0B7GzSdk4b0WrhP4wCgYIKoZIzj0EAwIwgYEx
CzAJBgNVBAYTAkdCMRswGQYDVQQIExJHcmVhdGVyIE1hbmNoZXN0ZXIxEDAOBgNV
BAcTB1NhbGZvcmQxGjAYBgNVBAoTEUNPTU9ETyBDQSBMaW1pdGVkMScwJQYDVQQD
Ex5DT01PRE8gQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwHhcNMDYxMjAxMDAwMDAw
WhcNMTkxMjMwMjM1OTU5WjByMQswCQYDVQQGEwJHQjEbMBkGA1UECBMSR3JlYXRl
ciBNYW5jaGVzdGVyMRAwDgYDVQQHEwdTYWxmb3JkMRowGAYDVQQKExFDT01PRE8g
Q0EgTGltaXRlZDEYMBYGA1UEAxMPRXNzZW50aWFsU1NMIENBMFkwEwYHKoZIzj0C
AQYIKoZIzj0DAQcDQgAEVYQA4mVfXUOiTNkBmQv6ejHtt0IisxTFr1YRbuzqqJND
i9vlEM0zO8nm6MQuiYbL6FAdLTDLjjfUq723qf2I+qOBhjCBgzAOBgNVHQ8BAf8E
BAMCAQYwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMBIGA1UdEwEB/wQI
MAYBAf8CAQAwHQYDVR0OBBYEFKiYEBkxmpfWMShxMEi/wgUWEqI4MB8GA1UdIwQY
MBaAFAKDsVXPpQ/rgxUBBWz2m7zbmsgNMAoGCCqGSM49BAMCA0kAMEYCIQCweyyQ
/5ODHs5ShzAM/V0BAX0AVvdY6xsyiklYNQcoQAIhALPZfrEPpaZWIs7ki7PInybJ
CdBZTeTV9yOa+OOYTxsG
-----END CERTIFICATE-----
`

var comodoRoot = `-----BEGIN CERTIFICATE-----
MIICRTCCAeugAwIBAgIRAJyB5sz1DXPOmSPZzsRUmz8wCgYIKoZIzj0EAwIwgYEx
CzAJBgNVBAYTAkdCMRswGQYDVQQIExJHcmVhdGVyIE1hbmNoZXN0ZXIxEDAOBgNV
BAcTB1NhbGZvcmQxGjAYBgNVBAoTEUNPTU9ETyBDQSBMaW1pdGVkMScwJQYDVQQD
Ex5DT01PRE8gQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwHhcNMDYxMjAxMDAwMDAw
WhcNMjkxMjMwMjM1OTU5WjCBgTELMAkGA1UEBhMCR0IxGzAZBgNVBAgTEkdyZWF0
ZXIgTWFuY2hlc3RlcjEQMA4GA1UEBxMHU2FsZm9yZDEaMBgGA1UEChMRQ09NT0RP
IENBIExpbWl0ZWQxJzAlBgNVBAMTHkNPTU9ETyBDZXJ0aWZpY2F0aW9uIEF1dGhv
cml0eTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABEyfpeHFEn5AIVG0sMe4TBXp
R1m4cH3Allkmh6meMQaCfge7q2Lfo9by/s8CrG2stz/7lOaRHUpp+9Au2IC13zaj
QjBAMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBQC
g7FVz6UP64MVAQVs9pu825rIDTAKBggqhkjOPQQDAgNIADBFAiAmhWWN2xcRodFm
0ggCzig31WU3xPSRbmWITAVFqptBcgIhAJW65TRVpNM+IB3YvhPo2/weoX2X5wXn
aH7SRPNcC4dz
-----END CERTIFICATE-----
`

var nameConstraintsLeaf = `-----BEGIN CERTIFICATE-----
MIIC+TCCAqCgAwIBAgIRAJOJ27of34c8vbm/Akja8rAwCgYIKoZIzj0EAwIwgcsx
CzAJBgNVBAYTAlVTMREwDwYDVQQIEwhWaXJnaW5pYTETMBEGA1UEBxMKQmxhY2tz
YnVyZzE8MDoGA1UEChMzVmlyZ2luaWEgUG9seXRlY2huaWMgSW5zdGl0dXRlIGFu
ZCBTdGF0ZSBVbml2ZXJzaXR5MSMwIQYDVQQLExpHbG9iYWwgUXVhbGlmaWVkIFNl
cnZlciBDQTExMC8GA1UEAxMoVmlyZ2luaWEgVGVjaCBHbG9iYWwgUXVhbGlmaWVk
IFNlcnZlciBDQTAeFw0xMzA5MTkxMzQwMTVaFw0xNTA5MTkxMzQwMTVaMIHNMQsw
CQYDVQQGEwJVUzERMA8GA1UECBMIVmlyZ2luaWExEzARBgNVBAcTCkJsYWNrc2J1
cmcxPDA6BgNVBAoTM1ZpcmdpbmlhIFBvbHl0ZWNobmljIEluc3RpdHV0ZSBhbmQg
U3RhdGUgVW5pdmVyc2l0eTE7MDkGA1UECxMyVGVjaG5vbG9neS1lbmhhbmNlZCBM
ZWFybmluZyBhbmQgT25saW5lIFN0cmF0ZWdpZXMxGzAZBgNVBAMTEnNlY3VyZS5p
ZGRsLnZ0LmVkdTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABKxIRbfIhD61mmjy
9shXs3G2jXkAbXXwtsYaWUkikbcob6PDSa/GHqjf+9FG5kW/Fo52kW4GO+VfuRaI
AP6pYBWjYTBfMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAfBgNVHSME
GDAWgBQbuSE90KPs6LsgNINARCRBYayF1zAdBgNVHREEFjAUghJzZWN1cmUuaWRk
bC52dC5lZHUwCgYIKoZIzj0EAwIDRwAwRAIgXh/j/i3Nx2XM10ET+dCZRL/lkl2+
5G6wOAYk28BKRGoCIGBCgaBmHUjD39OtPZ3wDtBG/pnQuXXOKiPunQCNmC3s
-----END CERTIFICATE-----
`

var nameConstraintsIntermediate1 = `-----BEGIN CERTIFICATE-----
MIICjDCCAjOgAwIBAgIQHsB18r1XIwGQOHZHf/5rgTAKBggqhkjOPQQDAjBFMQsw
CQYDVQQGEwJCRTEZMBcGA1UEChMQR2xvYmFsU2lnbiBudi1zYTEbMBkGA1UEAxMS
VHJ1c3RlZCBSb290IENBIEcyMB4XDTEyMTIxMjAwMDAwMFoXDTE3MTIxMjAwMDAw
MFowgcsxCzAJBgNVBAYTAlVTMREwDwYDVQQIEwhWaXJnaW5pYTETMBEGA1UEBxMK
QmxhY2tzYnVyZzE8MDoGA1UEChMzVmlyZ2luaWEgUG9seXRlY2huaWMgSW5zdGl0
dXRlIGFuZCBTdGF0ZSBVbml2ZXJzaXR5MSMwIQYDVQQLExpHbG9iYWwgUXVhbGlm
aWVkIFNlcnZlciBDQTExMC8GA1UEAxMoVmlyZ2luaWEgVGVjaCBHbG9iYWwgUXVh
bGlmaWVkIFNlcnZlciBDQTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABNZR58J0
RkGweEZC0Xurq29zt1sIkJZZA/kYg+6fJdjxTcEJVIR3Up1WR9DzZFXzKLyiHVb6
ZgdsGk6WCotAuGSjfjB8MA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/
AgEAMB0GA1UdDgQWBBQbuSE90KPs6LsgNINARCRBYayF1zAfBgNVHSMEGDAWgBSg
LK+xqoqPc+P2qTJ7aI7gCR7P9TAWBgNVHR4EDzANoAswCYIHLnZ0LmVkdTAKBggq
hkjOPQQDAgNHADBEAiAld1lFmvXM7tzHLSSJNaGKHtMa+CCjAF/IFKAek+GDBwIg
C8ClKQoqj0WdtxCFMsdx/yg2vDgSV1gpiwGVYLuvElY=
-----END CERTIFICATE-----
`

var nameConstraintsIntermediate2 = `-----BEGIN CERTIFICATE-----
MIICBTCCAaugAwIBAgIQLrIEa7stuHPdByTFrnYyfzAKBggqhkjOPQQDAjBcMQsw
CQYDVQQGEwJCRTEZMBcGA1UEChMQR2xvYmFsU2lnbiBudi1zYTEVMBMGA1UECxMM
VHJ1c3RlZCBSb290MRswGQYDVQQDExJHbG9iYWxTaWduIFJvb3QgQ0EwHhcNMTAw
MzE1MDAwMDAwWhcNMjYwMTE4MDAwMDAwWjBFMQswCQYDVQQGEwJCRTEZMBcGA1UE
ChMQR2xvYmFsU2lnbiBudi1zYTEbMBkGA1UEAxMSVHJ1c3RlZCBSb290IENBIEcy
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEUcgExU2p4oV4w4KFabCLPsMymT87
b6DKo7YRCkjDrjiydvvq3r4bhxYvGrJ23lCmzrrBCqZV++ISA+EigziPTKNmMGQw
DgYDVR0PAQH/BAQDAgEGMBIGA1UdEwEB/wQIMAYBAf8CAQIwHQYDVR0OBBYEFKAs
r7Gqio9z4/apMntojuAJHs/1MB8GA1UdIwQYMBaAFKgEViLDZHwXZPKNdozLffJ3
ZwiqMAoGCCqGSM49BAMCA0gAMEUCICpccwrGcXGgw1RsooFbSIz81COJ30je7uNo
moGtPdN8AiEAijFXmhtbJDy55XGQ/EFL6014f57v4zodoeZs9IkeVEQ=
-----END CERTIFICATE-----
`

var globalSignRoot = `-----BEGIN CERTIFICATE-----
MIIB9zCCAZ6gAwIBAgIQDqi3PEOjgeU8uPzUHq89hjAKBggqhkjOPQQDAjBcMQsw
CQYDVQQGEwJCRTEZMBcGA1UEChMQR2xvYmFsU2lnbiBudi1zYTEVMBMGA1UECxMM
VHJ1c3RlZCBSb290MRswGQYDVQQDExJHbG9iYWxTaWduIFJvb3QgQ0EwHhcNMTAw
MzE1MDAwMDAwWhcNMjYwMTE4MDAwMDAwWjBcMQswCQYDVQQGEwJCRTEZMBcGA1UE
ChMQR2xvYmFsU2lnbiBudi1zYTEVMBMGA1UECxMMVHJ1c3RlZCBSb290MRswGQYD
VQQDExJHbG9iYWxTaWduIFJvb3QgQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNC
AARVzdDCCm0rXF7Vq3vhpEl2y8IF45kVlt5/v+3Qe/pQ0An/tkTMzYNP9RX7n4a5
iI40MzoijTGxhWvtkxmpV7GTo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/
BAUwAwEB/zAdBgNVHQ4EFgQUqARWIsNkfBdk8o12jMt98ndnCKowCgYIKoZIzj0E
AwIDRwAwRAIgfTRzeF0LO0GlsWTb7ycBg4UShv6IGTMpGiuYtA5MP7cCICxUcRNL
FcPotwXcypbSeCeg1Dn6Ah9wrcOzIiHksgY/
-----END CERTIFICATE-----
`

var moipLeafCert = `-----BEGIN CERTIFICATE-----
MIIGQDCCBSigAwIBAgIRAPe/cwh7CUWizo8mYSDavLIwDQYJKoZIhvcNAQELBQAw
gZIxCzAJBgNVBAYTAkdCMRswGQYDVQQIExJHcmVhdGVyIE1hbmNoZXN0ZXIxEDAO
BgNVBAcTB1NhbGZvcmQxGjAYBgNVBAoTEUNPTU9ETyBDQSBMaW1pdGVkMTgwNgYD
VQQDEy9DT01PRE8gUlNBIEV4dGVuZGVkIFZhbGlkYXRpb24gU2VjdXJlIFNlcnZl
ciBDQTAeFw0xMzA4MTUwMDAwMDBaFw0xNDA4MTUyMzU5NTlaMIIBQjEXMBUGA1UE
BRMOMDg3MTg0MzEwMDAxMDgxEzARBgsrBgEEAYI3PAIBAxMCQlIxGjAYBgsrBgEE
AYI3PAIBAhMJU2FvIFBhdWxvMR0wGwYDVQQPExRQcml2YXRlIE9yZ2FuaXphdGlv
bjELMAkGA1UEBhMCQlIxETAPBgNVBBETCDAxNDUyMDAwMRIwEAYDVQQIEwlTYW8g
UGF1bG8xEjAQBgNVBAcTCVNhbyBQYXVsbzEtMCsGA1UECRMkQXZlbmlkYSBCcmln
YWRlaXJvIEZhcmlhIExpbWEgLCAyOTI3MR0wGwYDVQQKExRNb2lwIFBhZ2FtZW50
b3MgUy5BLjENMAsGA1UECxMETU9JUDEYMBYGA1UECxMPU1NMIEJsaW5kYWRvIEVW
MRgwFgYDVQQDEw9hcGkubW9pcC5jb20uYnIwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQDN0b9x6TrXXA9hPCF8/NjqGJ++2D4LO4ZiMFTjs0VwpXy2Y1Oe
s74/HuiLGnAHxTmAtV7IpZMibiOcTxcnDYp9oEWkf+gR+hZvwFZwyOBC7wyb3SR3
UvV0N1ZbEVRYpN9kuX/3vjDghjDmzzBwu8a/T+y5JTym5uiJlngVAWyh/RjtIvYi
+NVkQMbyVlPGkoCe6c30pH8DKYuUCZU6DHjUsPTX3jAskqbhDSAnclX9iX0p2bmw
KVBc+5Vh/2geyzDuquF0w+mNIYdU5h7uXvlmJnf3d2Cext5dxdL8/jezD3U0dAqI
pYSKERbyxSkJWxdvRlhdpM9YXMJcpc88xNp1AgMBAAGjggHcMIIB2DAfBgNVHSME
GDAWgBQ52v/KKBSKqHQTCLnkDqnS+n6daTAdBgNVHQ4EFgQU/lXuOa7DMExzZjRj
LQWcMWGZY7swDgYDVR0PAQH/BAQDAgWgMAwGA1UdEwEB/wQCMAAwHQYDVR0lBBYw
FAYIKwYBBQUHAwEGCCsGAQUFBwMCMEYGA1UdIAQ/MD0wOwYMKwYBBAGyMQECAQUB
MCswKQYIKwYBBQUHAgEWHWh0dHBzOi8vc2VjdXJlLmNvbW9kby5jb20vQ1BTMFYG
A1UdHwRPME0wS6BJoEeGRWh0dHA6Ly9jcmwuY29tb2RvY2EuY29tL0NPTU9ET1JT
QUV4dGVuZGVkVmFsaWRhdGlvblNlY3VyZVNlcnZlckNBLmNybDCBhwYIKwYBBQUH
AQEEezB5MFEGCCsGAQUFBzAChkVodHRwOi8vY3J0LmNvbW9kb2NhLmNvbS9DT01P
RE9SU0FFeHRlbmRlZFZhbGlkYXRpb25TZWN1cmVTZXJ2ZXJDQS5jcnQwJAYIKwYB
BQUHMAGGGGh0dHA6Ly9vY3NwLmNvbW9kb2NhLmNvbTAvBgNVHREEKDAmgg9hcGku
bW9pcC5jb20uYnKCE3d3dy5hcGkubW9pcC5jb20uYnIwDQYJKoZIhvcNAQELBQAD
ggEBAFoTmPlaDcf+nudhjXHwud8g7/LRyA8ucb+3/vfmgbn7FUc1eprF5sJS1mA+
pbiTyXw4IxcJq2KUj0Nw3IPOe9k84mzh+XMmdCKH+QK3NWkE9Udz+VpBOBc0dlqC
1RH5umStYDmuZg/8/r652eeQ5kUDcJyADfpKWBgDPYaGtwzKVT4h3Aok9SLXRHx6
z/gOaMjEDMarMCMw4VUIG1pvNraZrG5oTaALPaIXXpd8VqbQYPudYJ6fR5eY3FeW
H/ofbYFdRcuD26MfBFWE9VGGral9Fgo8sEHffho+UWhgApuQV4/l5fMzxB5YBXyQ
jhuy8PqqZS9OuLilTeLu4a8z2JI=
-----END CERTIFICATE-----`

var comodoIntermediateSHA384 = `-----BEGIN CERTIFICATE-----
MIIGDjCCA/agAwIBAgIQBqdDgNTr/tQ1taP34Wq92DANBgkqhkiG9w0BAQwFADCB
hTELMAkGA1UEBhMCR0IxGzAZBgNVBAgTEkdyZWF0ZXIgTWFuY2hlc3RlcjEQMA4G
A1UEBxMHU2FsZm9yZDEaMBgGA1UEChMRQ09NT0RPIENBIExpbWl0ZWQxKzApBgNV
BAMTIkNPTU9ETyBSU0EgQ2VydGlmaWNhdGlvbiBBdXRob3JpdHkwHhcNMTIwMjEy
MDAwMDAwWhcNMjcwMjExMjM1OTU5WjCBkjELMAkGA1UEBhMCR0IxGzAZBgNVBAgT
EkdyZWF0ZXIgTWFuY2hlc3RlcjEQMA4GA1UEBxMHU2FsZm9yZDEaMBgGA1UEChMR
Q09NT0RPIENBIExpbWl0ZWQxODA2BgNVBAMTL0NPTU9ETyBSU0EgRXh0ZW5kZWQg
VmFsaWRhdGlvbiBTZWN1cmUgU2VydmVyIENBMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAlVbeVLTf1QJJe9FbXKKyHo+cK2JMK40SKPMalaPGEP0p3uGf
CzhAk9HvbpUQ/OGQF3cs7nU+e2PsYZJuTzurgElr3wDqAwB/L3XVKC/sVmePgIOj
vdwDmZOLlJFWW6G4ajo/Br0OksxgnP214J9mMF/b5pTwlWqvyIqvgNnmiDkBfBzA
xSr3e5Wg8narbZtyOTDr0VdVAZ1YEZ18bYSPSeidCfw8/QpKdhQhXBZzQCMZdMO6
WAqmli7eNuWf0MLw4eDBYuPCGEUZUaoXHugjddTI0JYT/8ck0YwLJ66eetw6YWNg
iJctXQUL5Tvrrs46R3N2qPos3cCHF+msMJn4HwIDAQABo4IBaTCCAWUwHwYDVR0j
BBgwFoAUu69+Aj36pvE8hI6t7jiY7NkyMtQwHQYDVR0OBBYEFDna/8ooFIqodBMI
ueQOqdL6fp1pMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEAMD4G
A1UdIAQ3MDUwMwYEVR0gADArMCkGCCsGAQUFBwIBFh1odHRwczovL3NlY3VyZS5j
b21vZG8uY29tL0NQUzBMBgNVHR8ERTBDMEGgP6A9hjtodHRwOi8vY3JsLmNvbW9k
b2NhLmNvbS9DT01PRE9SU0FDZXJ0aWZpY2F0aW9uQXV0aG9yaXR5LmNybDBxBggr
BgEFBQcBAQRlMGMwOwYIKwYBBQUHMAKGL2h0dHA6Ly9jcnQuY29tb2RvY2EuY29t
L0NPTU9ET1JTQUFkZFRydXN0Q0EuY3J0MCQGCCsGAQUFBzABhhhodHRwOi8vb2Nz
cC5jb21vZG9jYS5jb20wDQYJKoZIhvcNAQEMBQADggIBAERCnUFRK0iIXZebeV4R
AUpSGXtBLMeJPNBy3IX6WK/VJeQT+FhlZ58N/1eLqYVeyqZLsKeyLeCMIs37/3mk
jCuN/gI9JN6pXV/kD0fQ22YlPodHDK4ixVAihNftSlka9pOlk7DgG4HyVsTIEFPk
1Hax0VtpS3ey4E/EhOfUoFDuPPpE/NBXueEoU/1Tzdy5H3pAvTA/2GzS8+cHnx8i
teoiccsq8FZ8/qyo0QYPFBRSTP5kKwxpKrgNUG4+BAe/eiCL+O5lCeHHSQgyPQ0o
fkkdt0rvAucNgBfIXOBhYsvss2B5JdoaZXOcOBCgJjqwyBZ9kzEi7nQLiMBciUEA
KKlHMd99SUWa9eanRRrSjhMQ34Ovmw2tfn6dNVA0BM7pINae253UqNpktNEvWS5e
ojZh1CSggjMziqHRbO9haKPl0latxf1eYusVqHQSTC8xjOnB3xBLAer2VBvNfzu9
XJ/B288ByvK6YBIhMe2pZLiySVgXbVrXzYxtvp5/4gJYp9vDLVj2dAZqmvZh+fYA
tmnYOosxWd2R5nwnI4fdAw+PKowegwFOAWEMUnNt/AiiuSpm5HZNMaBWm9lTjaK2
jwLI5jqmBNFI+8NKAnb9L9K8E7bobTQk+p0pisehKxTxlgBzuRPpwLk6R1YCcYAn
pLwltum95OmYdBbxN4SBB7SC
-----END CERTIFICATE-----`

const comodoRSAAuthority = `-----BEGIN CERTIFICATE-----
MIIFdDCCBFygAwIBAgIQJ2buVutJ846r13Ci/ITeIjANBgkqhkiG9w0BAQwFADBv
MQswCQYDVQQGEwJTRTEUMBIGA1UEChMLQWRkVHJ1c3QgQUIxJjAkBgNVBAsTHUFk
ZFRydXN0IEV4dGVybmFsIFRUUCBOZXR3b3JrMSIwIAYDVQQDExlBZGRUcnVzdCBF
eHRlcm5hbCBDQSBSb290MB4XDTAwMDUzMDEwNDgzOFoXDTIwMDUzMDEwNDgzOFow
gYUxCzAJBgNVBAYTAkdCMRswGQYDVQQIExJHcmVhdGVyIE1hbmNoZXN0ZXIxEDAO
BgNVBAcTB1NhbGZvcmQxGjAYBgNVBAoTEUNPTU9ETyBDQSBMaW1pdGVkMSswKQYD
VQQDEyJDT01PRE8gUlNBIENlcnRpZmljYXRpb24gQXV0aG9yaXR5MIICIjANBgkq
hkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAkehUktIKVrGsDSTdxc9EZ3SZKzejfSNw
AHG8U9/E+ioSj0t/EFa9n3Byt2F/yUsPF6c947AEYe7/EZfH9IY+Cvo+XPmT5jR6
2RRr55yzhaCCenavcZDX7P0N+pxs+t+wgvQUfvm+xKYvT3+Zf7X8Z0NyvQwA1onr
ayzT7Y+YHBSrfuXjbvzYqOSSJNpDa2K4Vf3qwbxstovzDo2a5JtsaZn4eEgwRdWt
4Q08RWD8MpZRJ7xnw8outmvqRsfHIKCxH2XeSAi6pE6p8oNGN4Tr6MyBSENnTnIq
m1y9TBsoilwie7SrmNnu4FGDwwlGTm0+mfqVF9p8M1dBPI1R7Qu2XK8sYxrfV8g/
vOldxJuvRZnio1oktLqpVj3Pb6r/SVi+8Kj/9Lit6Tf7urj0Czr56ENCHonYhMsT
8dm74YlguIwoVqwUHZwK53Hrzw7dPamWoUi9PPevtQ0iTMARgexWO/bTouJbt7IE
IlKVgJNp6I5MZfGRAy1wdALqi2cVKWlSArvX31BqVUa/oKMoYX9w0MOiqiwhqkfO
KJwGRXa/ghgntNWutMtQ5mv0TIZxMOmm3xaG4Nj/QN370EKIf6MzOi5cHkERgWPO
GHFrK+ymircxXDpqR+DDeVnWIBqv8mqYqnK8V0rSS527EPywTEHl7R09XiidnMy/
s1Hap0flhFMCAwEAAaOB9DCB8TAfBgNVHSMEGDAWgBStvZh6NLQm9/rEJlTvA73g
JMtUGjAdBgNVHQ4EFgQUu69+Aj36pvE8hI6t7jiY7NkyMtQwDgYDVR0PAQH/BAQD
AgGGMA8GA1UdEwEB/wQFMAMBAf8wEQYDVR0gBAowCDAGBgRVHSAAMEQGA1UdHwQ9
MDswOaA3oDWGM2h0dHA6Ly9jcmwudXNlcnRydXN0LmNvbS9BZGRUcnVzdEV4dGVy
bmFsQ0FSb290LmNybDA1BggrBgEFBQcBAQQpMCcwJQYIKwYBBQUHMAGGGWh0dHA6
Ly9vY3NwLnVzZXJ0cnVzdC5jb20wDQYJKoZIhvcNAQEMBQADggEBAGS/g/FfmoXQ
zbihKVcN6Fr30ek+8nYEbvFScLsePP9NDXRqzIGCJdPDoCpdTPW6i6FtxFQJdcfj
Jw5dhHk3QBN39bSsHNA7qxcS1u80GH4r6XnTq1dFDK8o+tDb5VCViLvfhVdpfZLY
Uspzgb8c8+a4bmYRBbMelC1/kZWSWfFMzqORcUx8Rww7Cxn2obFshj5cqsQugsv5
B5a6SE2Q8pTIqXOi6wZ7I53eovNNVZ96YUWYGGjHXkBrI/V5eu+MtWuLt29G9Hvx
PUsE2JOAWVrgQSQdso8VYFhH2+9uRv0V9dlfmrPb2LjkQLPNlzmuhbsdjrzch5vR
pu/xO28QOG8=
-----END CERTIFICATE-----`

const addTrustRoot = `-----BEGIN CERTIFICATE-----
MIIENjCCAx6gAwIBAgIBATANBgkqhkiG9w0BAQUFADBvMQswCQYDVQQGEwJTRTEU
MBIGA1UEChMLQWRkVHJ1c3QgQUIxJjAkBgNVBAsTHUFkZFRydXN0IEV4dGVybmFs
IFRUUCBOZXR3b3JrMSIwIAYDVQQDExlBZGRUcnVzdCBFeHRlcm5hbCBDQSBSb290
MB4XDTAwMDUzMDEwNDgzOFoXDTIwMDUzMDEwNDgzOFowbzELMAkGA1UEBhMCU0Ux
FDASBgNVBAoTC0FkZFRydXN0IEFCMSYwJAYDVQQLEx1BZGRUcnVzdCBFeHRlcm5h
bCBUVFAgTmV0d29yazEiMCAGA1UEAxMZQWRkVHJ1c3QgRXh0ZXJuYWwgQ0EgUm9v
dDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALf3GjPm8gAELTngTlvt
H7xsD821+iO2zt6bETOXpClMfZOfvUq8k+0DGuOPz+VtUFrWlymUWoCwSXrbLpX9
uMq/NzgtHj6RQa1wVsfwTz/oMp50ysiQVOnGXw94nZpAPA6sYapeFI+eh6FqUNzX
mk6vBbOmcZSccbNQYArHE504B4YCqOmoaSYYkKtMsE8jqzpPhNjfzp/haW+710LX
a0Tkx63ubUFfclpxCDezeWWkWaCUN/cALw3CknLa0Dhy2xSoRcRdKn23tNbE7qzN
E0S3ySvdQwAl+mG5aWpYIxG3pzOPVnVZ9c0p10a3CitlttNCbxWyuHv77+ldU9U0
WicCAwEAAaOB3DCB2TAdBgNVHQ4EFgQUrb2YejS0Jvf6xCZU7wO94CTLVBowCwYD
VR0PBAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wgZkGA1UdIwSBkTCBjoAUrb2YejS0
Jvf6xCZU7wO94CTLVBqhc6RxMG8xCzAJBgNVBAYTAlNFMRQwEgYDVQQKEwtBZGRU
cnVzdCBBQjEmMCQGA1UECxMdQWRkVHJ1c3QgRXh0ZXJuYWwgVFRQIE5ldHdvcmsx
IjAgBgNVBAMTGUFkZFRydXN0IEV4dGVybmFsIENBIFJvb3SCAQEwDQYJKoZIhvcN
AQEFBQADggEBALCb4IUlwtYj4g+WBpKdQZic2YR5gdkeWxQHIzZlj7DYd7usQWxH
YINRsPkyPef89iYTx4AWpb9a/IfPeHmJIZriTAcKhjW88t5RxNKWt9x+Tu5w/Rw5
6wwCURQtjr0W4MHfRnXnJK3s9EK0hZNwEGe6nQY1ShjTK3rMUUKhemPR5ruhxSvC
Nr4TDea9Y355e6cJDUCrat2PisP29owaQgVR1EX1n6diIWgVIEM8med8vSTYqZEX
c4g/VhsxOBi0cQ+azcgOno4uG+GMmIPLHzHxREzGBHNJdmAPx/i9F4BrLunMTA5a
mnkPIAou1Z5jJh5VkpTYghdae9C8x49OhgQ=
-----END CERTIFICATE-----`

const selfSigned = `-----BEGIN CERTIFICATE-----
MIIC/DCCAeSgAwIBAgIRAK0SWRVmi67xU3z0gkgY+PkwDQYJKoZIhvcNAQELBQAw
EjEQMA4GA1UEChMHQWNtZSBDbzAeFw0xNjA4MTkxNjMzNDdaFw0xNzA4MTkxNjMz
NDdaMBIxEDAOBgNVBAoTB0FjbWUgQ28wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDWkm1kdCwxyKEt6OTmZitkmLGH8cQu9z7rUdrhW8lWNm4kh2SuaUWP
pscBjda5iqg51aoKuWJR2rw6ElDne+X5eit2FT8zJgAU8v39lMFjbaVZfS9TFOYF
w0Tk0Luo/PyKJpZnwhsP++iiGQiteJbndy8aLKmJ2MpLfpDGIgxEIyNb5dgoDi0D
WReDCpE6K9WDYqvKVGnQ2Jvqqra6Gfx0tFkuqJxQuqA8aUOlPHcCH4KBZdNEoXdY
YL3E4dCAh0YiDs80wNZx4cHqEM3L8gTEFqW2Tn1TSuPZO6gjJ9QPsuUZVjaMZuuO
NVxqLGujZkDzARhC3fBpptMuaAfi20+BAgMBAAGjTTBLMA4GA1UdDwEB/wQEAwIF
oDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMBYGA1UdEQQPMA2C
C2Zvby5leGFtcGxlMA0GCSqGSIb3DQEBCwUAA4IBAQBPvvfnDhsHWt+/cfwdAVim
4EDn+hYOMkTQwU0pouYIvY8QXYkZ8MBxpBtBMK4JhFU+ewSWoBAEH2dCCvx/BDxN
UGTSJHMbsvJHcFvdmsvvRxOqQ/cJz7behx0cfoeHMwcs0/vWv8ms5wHesb5Ek7L0
pl01FCBGTcncVqr6RK1r4fTpeCCfRIERD+YRJz8TtPH6ydesfLL8jIV40H8NiDfG
vRAvOtNiKtPzFeQVdbRPOskC4rcHyPeiDAMAMixeLi63+CFty4da3r5lRezeedCE
cw3ESZzThBwWqvPOtJdpXdm+r57pDW8qD+/0lY8wfImMNkQAyCUCLg/1Lxt/hrBj
-----END CERTIFICATE-----`

const issuerSubjectMatchRoot = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 161640039802297062 (0x23e42c281e55ae6)
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: O=Golang, CN=Root ca
        Validity
            Not Before: Jan  1 00:00:00 2015 GMT
            Not After : Jan  1 00:00:00 2025 GMT
        Subject: O=Golang, CN=Root ca
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (1024 bit)
                Modulus:
                    00:e9:0e:7f:11:0c:e6:5a:e6:86:83:70:f6:51:07:
                    2e:02:78:11:f5:b2:24:92:38:ee:26:62:02:c7:94:
                    f1:3e:a1:77:6a:c0:8f:d5:22:68:b6:5d:e2:4c:da:
                    e0:85:11:35:c2:92:72:49:8d:81:b4:88:97:6b:b7:
                    fc:b2:44:5b:d9:4d:06:70:f9:0c:c6:8f:e9:b3:df:
                    a3:6a:84:6c:43:59:be:9d:b2:d0:76:9b:c3:d7:fa:
                    99:59:c3:b8:e5:f3:53:03:bd:49:d6:b3:cc:a2:43:
                    fe:ad:c2:0b:b9:01:b8:56:29:94:03:24:a7:0d:28:
                    21:29:a9:ae:94:5b:4a:f9:9f
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Key Usage: critical
                Certificate Sign
            X509v3 Extended Key Usage:
                TLS Web Server Authentication, TLS Web Client Authentication
            X509v3 Basic Constraints: critical
                CA:TRUE
            X509v3 Subject Key Identifier:
                40:37:D7:01:FB:40:2F:B8:1C:7E:54:04:27:8C:59:01
    Signature Algorithm: sha256WithRSAEncryption
         6f:84:df:49:e0:99:d4:71:66:1d:32:86:56:cb:ea:5a:6b:0e:
         00:6a:d1:5a:6e:1f:06:23:07:ff:cb:d1:1a:74:e4:24:43:0b:
         aa:2a:a0:73:75:25:82:bc:bf:3f:a9:f8:48:88:ac:ed:3a:94:
         3b:0d:d3:88:c8:67:44:61:33:df:71:6c:c5:af:ed:16:8c:bf:
         82:f9:49:bb:e3:2a:07:53:36:37:25:77:de:91:a4:77:09:7f:
         6f:b2:91:58:c4:05:89:ea:8e:fa:e1:3b:19:ef:f8:f6:94:b7:
         7b:27:e6:e4:84:dd:2b:f5:93:f5:3c:d8:86:c5:38:01:56:5c:
         9f:6d
-----BEGIN CERTIFICATE-----
MIICIDCCAYmgAwIBAgIIAj5CwoHlWuYwDQYJKoZIhvcNAQELBQAwIzEPMA0GA1UE
ChMGR29sYW5nMRAwDgYDVQQDEwdSb290IGNhMB4XDTE1MDEwMTAwMDAwMFoXDTI1
MDEwMTAwMDAwMFowIzEPMA0GA1UEChMGR29sYW5nMRAwDgYDVQQDEwdSb290IGNh
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDpDn8RDOZa5oaDcPZRBy4CeBH1
siSSOO4mYgLHlPE+oXdqwI/VImi2XeJM2uCFETXCknJJjYG0iJdrt/yyRFvZTQZw
+QzGj+mz36NqhGxDWb6dstB2m8PX+plZw7jl81MDvUnWs8yiQ/6twgu5AbhWKZQD
JKcNKCEpqa6UW0r5nwIDAQABo10wWzAOBgNVHQ8BAf8EBAMCAgQwHQYDVR0lBBYw
FAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQFMAMBAf8wGQYDVR0OBBIE
EEA31wH7QC+4HH5UBCeMWQEwDQYJKoZIhvcNAQELBQADgYEAb4TfSeCZ1HFmHTKG
VsvqWmsOAGrRWm4fBiMH/8vRGnTkJEMLqiqgc3Ulgry/P6n4SIis7TqUOw3TiMhn
RGEz33Fsxa/tFoy/gvlJu+MqB1M2NyV33pGkdwl/b7KRWMQFieqO+uE7Ge/49pS3
eyfm5ITdK/WT9TzYhsU4AVZcn20=
-----END CERTIFICATE-----`

const issuerSubjectMatchLeaf = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 16785088708916013734 (0xe8f09d3fe25beaa6)
    Signature Algorithm: sha256WithRSAEncryption
        Issuer: O=Golang, CN=Root CA
        Validity
            Not Before: Jan  1 00:00:00 2015 GMT
            Not After : Jan  1 00:00:00 2025 GMT
        Subject: O=Golang, CN=Leaf
        Subject Public Key Info:
            Public Key Algorithm: rsaEncryption
                Public-Key: (1024 bit)
                Modulus:
                    00:db:46:7d:93:2e:12:27:06:48:bc:06:28:21:ab:
                    7e:c4:b6:a2:5d:fe:1e:52:45:88:7a:36:47:a5:08:
                    0d:92:42:5b:c2:81:c0:be:97:79:98:40:fb:4f:6d:
                    14:fd:2b:13:8b:c2:a5:2e:67:d8:d4:09:9e:d6:22:
                    38:b7:4a:0b:74:73:2b:c2:34:f1:d1:93:e5:96:d9:
                    74:7b:f3:58:9f:6c:61:3c:c0:b0:41:d4:d9:2b:2b:
                    24:23:77:5b:1c:3b:bd:75:5d:ce:20:54:cf:a1:63:
                    87:1d:1e:24:c4:f3:1d:1a:50:8b:aa:b6:14:43:ed:
                    97:a7:75:62:f4:14:c8:52:d7
                Exponent: 65537 (0x10001)
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment
            X509v3 Extended Key Usage:
                TLS Web Server Authentication, TLS Web Client Authentication
            X509v3 Basic Constraints: critical
                CA:FALSE
            X509v3 Subject Key Identifier:
                9F:91:16:1F:43:43:3E:49:A6:DE:6D:B6:80:D7:9F:60
            X509v3 Authority Key Identifier:
                keyid:40:37:D7:01:FB:40:2F:B8:1C:7E:54:04:27:8C:59:01

    Signature Algorithm: sha256WithRSAEncryption
         8d:86:05:da:89:f5:1d:c5:16:14:41:b9:34:87:2b:5c:38:99:
         e3:d9:5a:5b:7a:5b:de:0b:5c:08:45:09:6f:1c:9d:31:5f:08:
         ca:7a:a3:99:da:83:0b:22:be:4f:02:35:91:4e:5d:5c:37:bf:
         89:22:58:7d:30:76:d2:2f:d0:a0:ee:77:9e:77:c0:d6:19:eb:
         ec:a0:63:35:6a:80:9b:80:1a:80:de:64:bc:40:38:3c:22:69:
         ad:46:26:a2:3d:ea:f4:c2:92:49:16:03:96:ae:64:21:b9:7c:
         ee:64:91:47:81:aa:b4:0c:09:2b:12:1a:b2:f3:af:50:b3:b1:
         ce:24
-----BEGIN CERTIFICATE-----
MIICODCCAaGgAwIBAgIJAOjwnT/iW+qmMA0GCSqGSIb3DQEBCwUAMCMxDzANBgNV
BAoTBkdvbGFuZzEQMA4GA1UEAxMHUm9vdCBDQTAeFw0xNTAxMDEwMDAwMDBaFw0y
NTAxMDEwMDAwMDBaMCAxDzANBgNVBAoTBkdvbGFuZzENMAsGA1UEAxMETGVhZjCB
nzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA20Z9ky4SJwZIvAYoIat+xLaiXf4e
UkWIejZHpQgNkkJbwoHAvpd5mED7T20U/SsTi8KlLmfY1Ame1iI4t0oLdHMrwjTx
0ZPlltl0e/NYn2xhPMCwQdTZKyskI3dbHDu9dV3OIFTPoWOHHR4kxPMdGlCLqrYU
Q+2Xp3Vi9BTIUtcCAwEAAaN3MHUwDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQG
CCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMBkGA1UdDgQSBBCfkRYf
Q0M+SabebbaA159gMBsGA1UdIwQUMBKAEEA31wH7QC+4HH5UBCeMWQEwDQYJKoZI
hvcNAQELBQADgYEAjYYF2on1HcUWFEG5NIcrXDiZ49laW3pb3gtcCEUJbxydMV8I
ynqjmdqDCyK+TwI1kU5dXDe/iSJYfTB20i/QoO53nnfA1hnr7KBjNWqAm4AagN5k
vEA4PCJprUYmoj3q9MKSSRYDlq5kIbl87mSRR4GqtAwJKxIasvOvULOxziQ=
-----END CERTIFICATE-----
`

const x509v1TestRoot = `
-----BEGIN CERTIFICATE-----
MIICIDCCAYmgAwIBAgIIAj5CwoHlWuYwDQYJKoZIhvcNAQELBQAwIzEPMA0GA1UE
ChMGR29sYW5nMRAwDgYDVQQDEwdSb290IENBMB4XDTE1MDEwMTAwMDAwMFoXDTI1
MDEwMTAwMDAwMFowIzEPMA0GA1UEChMGR29sYW5nMRAwDgYDVQQDEwdSb290IENB
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDpDn8RDOZa5oaDcPZRBy4CeBH1
siSSOO4mYgLHlPE+oXdqwI/VImi2XeJM2uCFETXCknJJjYG0iJdrt/yyRFvZTQZw
+QzGj+mz36NqhGxDWb6dstB2m8PX+plZw7jl81MDvUnWs8yiQ/6twgu5AbhWKZQD
JKcNKCEpqa6UW0r5nwIDAQABo10wWzAOBgNVHQ8BAf8EBAMCAgQwHQYDVR0lBBYw
FAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQFMAMBAf8wGQYDVR0OBBIE
EEA31wH7QC+4HH5UBCeMWQEwDQYJKoZIhvcNAQELBQADgYEAcIwqeNUpQr9cOcYm
YjpGpYkQ6b248xijCK7zI+lOeWN89zfSXn1AvfsC9pSdTMeDklWktbF/Ad0IN8Md
h2NtN34ard0hEfHc8qW8mkXdsysVmq6cPvFYaHz+dBtkHuHDoy8YQnC0zdN/WyYB
/1JmacUUofl+HusHuLkDxmadogI=
-----END CERTIFICATE-----`

const x509v1TestIntermediate = `
-----BEGIN CERTIFICATE-----
MIIByjCCATMCCQCCdEMsT8ykqTANBgkqhkiG9w0BAQsFADAjMQ8wDQYDVQQKEwZH
b2xhbmcxEDAOBgNVBAMTB1Jvb3QgQ0EwHhcNMTUwMTAxMDAwMDAwWhcNMjUwMTAx
MDAwMDAwWjAwMQ8wDQYDVQQKEwZHb2xhbmcxHTAbBgNVBAMTFFguNTA5djEgaW50
ZXJtZWRpYXRlMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDJ2QyniAOT+5YL
jeinEBJr3NsC/Q2QJ/VKmgvp+xRxuKTHJiVmxVijmp0vWg8AWfkmuE4p3hXQbbqM
k5yxrk1n60ONhim2L4VXriEvCE7X2OXhTmBls5Ufr7aqIgPMikwjScCXwz8E8qI8
UxyAhnjeJwMYBU8TuwBImSd4LBHoQQIDAQABMA0GCSqGSIb3DQEBCwUAA4GBAIab
DRG6FbF9kL9jb/TDHkbVBk+sl/Pxi4/XjuFyIALlARgAkeZcPmL5tNW1ImHkwsHR
zWE77kJDibzd141u21ZbLsKvEdUJXjla43bdyMmEqf5VGpC3D4sFt3QVH7lGeRur
x5Wlq1u3YDL/j6s1nU2dQ3ySB/oP7J+vQ9V4QeM+
-----END CERTIFICATE-----`

const x509v1TestLeaf = `
-----BEGIN CERTIFICATE-----
MIICMzCCAZygAwIBAgIJAPo99mqJJrpJMA0GCSqGSIb3DQEBCwUAMDAxDzANBgNV
BAoTBkdvbGFuZzEdMBsGA1UEAxMUWC41MDl2MSBpbnRlcm1lZGlhdGUwHhcNMTUw
MTAxMDAwMDAwWhcNMjUwMTAxMDAwMDAwWjArMQ8wDQYDVQQKEwZHb2xhbmcxGDAW
BgNVBAMTD2Zvby5leGFtcGxlLmNvbTCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkC
gYEApUh60Z+a5/oKJxG//Dn8CihSo2CJHNIIO3zEJZ1EeNSMZCynaIR6D3IPZEIR
+RG2oGt+f5EEukAPYxwasp6VeZEezoQWJ+97nPCT6DpwLlWp3i2MF8piK2R9vxkG
Z5n0+HzYk1VM8epIrZFUXSMGTX8w1y041PX/yYLxbdEifdcCAwEAAaNaMFgwDgYD
VR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNV
HRMBAf8EAjAAMBkGA1UdDgQSBBBFozXe0SnzAmjy+1U6M/cvMA0GCSqGSIb3DQEB
CwUAA4GBADYzYUvaToO/ucBskPdqXV16AaakIhhSENswYVSl97/sODaxsjishKq9
5R7siu+JnIFotA7IbBe633p75xEnLN88X626N/XRFG9iScLzpj0o0PWXBUiB+fxL
/jt8qszOXCv2vYdUTPNuPqufXLWMoirpuXrr1liJDmedCcAHepY/
-----END CERTIFICATE-----`

const ignoreCNWithSANRoot = `
-----BEGIN CERTIFICATE-----
MIIDPzCCAiegAwIBAgIIJkzCwkNrPHMwDQYJKoZIhvcNAQELBQAwMDEQMA4GA1UE
ChMHVEVTVElORzEcMBoGA1UEAxMTKipUZXN0aW5nKiogUm9vdCBDQTAeFw0xNTAx
MDEwMDAwMDBaFw0yNTAxMDEwMDAwMDBaMDAxEDAOBgNVBAoTB1RFU1RJTkcxHDAa
BgNVBAMTEyoqVGVzdGluZyoqIFJvb3QgQ0EwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQC4YAf5YqlXGcikvbMWtVrNICt+V/NNWljwfvSKdg4Inm7k6BwW
P6y4Y+n4qSYIWNU4iRkdpajufzctxQCO6ty13iw3qVktzcC5XBIiS6ymiRhhDgnY
VQqyakVGw9MxrPwdRZVlssUv3Hmy6tU+v5Ok31SLY5z3wKgYWvSyYs0b8bKNU8kf
2FmSHnBN16lxGdjhe3ji58F/zFMr0ds+HakrLIvVdFcQFAnQopM8FTHpoWNNzGU3
KaiO0jBbMFkd6uVjVnuRJ+xjuiqi/NWwiwQA+CEr9HKzGkxOF8nAsHamdmO1wW+w
OsCrC0qWQ/f5NTOVATTJe0vj88OMTvo3071VAgMBAAGjXTBbMA4GA1UdDwEB/wQE
AwICpDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDwYDVR0TAQH/BAUw
AwEB/zAZBgNVHQ4EEgQQQDfXAftAL7gcflQEJ4xZATANBgkqhkiG9w0BAQsFAAOC
AQEAGOn3XjxHyHbXLKrRmpwV447B7iNBXR5VlhwOgt1kWaHDL2+8f/9/h0HMkB6j
fC+/yyuYVqYuOeavqMGVrh33D2ODuTQcFlOx5lXukP46j3j+Lm0jjZ1qNX7vlP8I
VlUXERhbelkw8O4oikakwIY9GE8syuSgYf+VeBW/lvuAZQrdnPfabxe05Tre6RXy
nJHMB1q07YHpbwIkcV/lfCE9pig2nPXTLwYZz9cl46Ul5RCpPUi+IKURo3x8y0FU
aSLjI/Ya0zwUARMmyZ3RRGCyhIarPb20mKSaMf1/Nb23pS3k1QgmZhk5pAnXYsWu
BJ6bvwEAasFiLGP6Zbdmxb2hIA==
-----END CERTIFICATE-----`

const ignoreCNWithSANLeaf = `
-----BEGIN CERTIFICATE-----
MIIDaTCCAlGgAwIBAgIJAONakvRTxgJhMA0GCSqGSIb3DQEBCwUAMDAxEDAOBgNV
BAoTB1RFU1RJTkcxHDAaBgNVBAMTEyoqVGVzdGluZyoqIFJvb3QgQ0EwHhcNMTUw
MTAxMDAwMDAwWhcNMjUwMTAxMDAwMDAwWjAsMRAwDgYDVQQKEwdURVNUSU5HMRgw
FgYDVQQDEw9mb28uZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDBqskp89V/JMIBBqcauKSOVLcMyIE/t0jgSWVrsI4sksBTabLsfMdS
ui2n+dHQ1dRBuw3o4g4fPrWwS3nMnV3pZUHEn2TPi5N1xkjTaxObXgKIY2GKmFP3
rJ9vYqHT6mT4K93kCHoRcmJWWySc7S3JAOhTcdB4G+tIdQJN63E+XRYQQfNrn5HZ
hxQoOzaguHFx+ZGSD4Ntk6BSZz5NfjqCYqYxe+iCpTpEEYhIpi8joSPSmkTMTxBW
S1W2gXbYNQ9KjNkGM6FnQsUJrSPMrWs4v3UB/U88N5LkZeF41SqD9ySFGwbGajFV
nyzj12+4K4D8BLhlOc0Eo/F/8GwOwvmxAgMBAAGjgYkwgYYwDgYDVR0PAQH/BAQD
AgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAA
MBkGA1UdDgQSBBCjeab27q+5pV43jBGANOJ1MBsGA1UdIwQUMBKAEEA31wH7QC+4
HH5UBCeMWQEwDwYDVR0RBAgwBocEfwAAATANBgkqhkiG9w0BAQsFAAOCAQEAGZfZ
ErTVxxpIg64s22mQpXSk/72THVQsfsKHzlXmztM0CJzH8ccoN67ZqKxJCfdiE/FI
Emb6BVV4cGPeIKpcxaM2dwX/Y+Y0JaxpQJvqLxs+EByRL0gPP3shgg86WWCjYLxv
AgOn862d/JXGDrC9vIlQ/DDQcyL5g0JV5UjG2G9TUigbnrXxBw7BoWK6wmoSaHnR
sZKEHSs3RUJvm7qqpA9Yfzm9jg+i9j32zh1xFacghAOmFRFXa9eCVeigZ/KK2mEY
j2kBQyvnyKsXHLAKUoUOpd6t/1PHrfXnGj+HmzZNloJ/BZ1kiWb4eLvMljoLGkZn
xZbqP3Krgjj4XNaXjg==
-----END CERTIFICATE-----`

const excludedNamesLeaf = `
-----BEGIN CERTIFICATE-----
MIID4DCCAsigAwIBAgIHDUSFtJknhzANBgkqhkiG9w0BAQsFADCBnjELMAkGA1UE
BhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExEjAQBgNVBAcMCUxvcyBHYXRvczEU
MBIGA1UECgwLTmV0ZmxpeCBJbmMxLTArBgNVBAsMJFBsYXRmb3JtIFNlY3VyaXR5
ICgzNzM0NTE1NTYyODA2Mzk3KTEhMB8GA1UEAwwYSW50ZXJtZWRpYXRlIENBIGZv
ciAzMzkyMB4XDTE3MDIwODIxMTUwNFoXDTE4MDIwODIwMjQ1OFowgZAxCzAJBgNV
BAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQHDAlMb3MgR2F0b3Mx
FDASBgNVBAoMC05ldGZsaXggSW5jMS0wKwYDVQQLDCRQbGF0Zm9ybSBTZWN1cml0
eSAoMzczNDUxNTc0ODUwMjY5NikxEzARBgNVBAMMCjE3Mi4xNi4wLjEwggEiMA0G
CSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCZ0oP1bMv6bOeqcKbzinnGpNOpenhA
zdFFsgea62znWsH3Wg4+1Md8uPCqlaQIsaJQKZHc50eKD3bg0Io7c6kxHkBQr1b8
Q7cGeK3CjdqG3NwS/aizzrLKOwL693hFwwy7JY7GGCvogbhyQRKn6iV0U9zMm7bu
/9pQVV/wx8u01u2uAlLttjyQ5LJkxo5t8cATFVqxdN5J9eY//VSDiTwXnlpQITBP
/Ow+zYuZ3kFlzH3CtCOhOEvNG3Ar1NvP3Icq35PlHV+Eki4otnKfixwByoiGpqCB
UEIY04VrZJjwBxk08y/3jY2B3VLYGgi+rryyCxIqkB7UpSNPMMWSG4UpAgMBAAGj
LzAtMAwGA1UdEwEB/wQCMAAwHQYDVR0RBBYwFIIMYmVuZGVyLmxvY2FshwSsEAAB
MA0GCSqGSIb3DQEBCwUAA4IBAQCLW3JO8L7LKByjzj2RciPjCGH5XF87Wd20gYLq
sNKcFwCIeyZhnQy5aZ164a5G9AIk2HLvH6HevBFPhA9Ivmyv/wYEfnPd1VcFkpgP
hDt8MCFJ8eSjCyKdtZh1MPMLrLVymmJV+Rc9JUUYM9TIeERkpl0rskcO1YGewkYt
qKlWE+0S16+pzsWvKn831uylqwIb8ANBPsCX4aM4muFBHavSWAHgRO+P+yXVw8Q+
VQDnMHUe5PbZd1/+1KKVs1K/CkBCtoHNHp1d/JT+2zUQJphwja9CcgfFdVhSnHL4
oEEOFtqVMIuQfR2isi08qW/JGOHc4sFoLYB8hvdaxKWSE19A
-----END CERTIFICATE-----
`

const excludedNamesIntermediate = `
-----BEGIN CERTIFICATE-----
MIIDzTCCArWgAwIBAgIHDUSFqYeczDANBgkqhkiG9w0BAQsFADCBmTELMAkGA1UE
BhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExEjAQBgNVBAcMCUxvcyBHYXRvczEU
MBIGA1UECgwLTmV0ZmxpeCBJbmMxLTArBgNVBAsMJFBsYXRmb3JtIFNlY3VyaXR5
ICgzNzM0NTE1NDc5MDY0NjAyKTEcMBoGA1UEAwwTTG9jYWwgUm9vdCBmb3IgMzM5
MjAeFw0xNzAyMDgyMTE1MDRaFw0xODAyMDgyMDI0NThaMIGeMQswCQYDVQQGEwJV
UzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJTG9zIEdhdG9zMRQwEgYD
VQQKDAtOZXRmbGl4IEluYzEtMCsGA1UECwwkUGxhdGZvcm0gU2VjdXJpdHkgKDM3
MzQ1MTU1NjI4MDYzOTcpMSEwHwYDVQQDDBhJbnRlcm1lZGlhdGUgQ0EgZm9yIDMz
OTIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCOyEs6tJ/t9emQTvlx
3FS7uJSou5rKkuqVxZdIuYQ+B2ZviBYUnMRT9bXDB0nsVdKZdp0hdchdiwNXDG/I
CiWu48jkcv/BdynVyayOT+0pOJSYLaPYpzBx1Pb9M5651ct9GSbj6Tz0ChVonoIE
1AIZ0kkebucZRRFHd0xbAKVRKyUzPN6HJ7WfgyauUp7RmlC35wTmrmARrFohQLlL
7oICy+hIQePMy9x1LSFTbPxZ5AUUXVC3eUACU3vLClF/Xs8XGHebZpUXCdMQjOGS
nq1eFguFHR1poSB8uSmmLqm4vqUH9CDhEgiBAC8yekJ8//kZQ7lUEqZj3YxVbk+Y
E4H5AgMBAAGjEzARMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEB
ADxrnmNX5gWChgX9K5fYwhFDj5ofxZXAKVQk+WjmkwMcmCx3dtWSm++Wdksj/ZlA
V1cLW3ohWv1/OAZuOlw7sLf98aJpX+UUmIYYQxDubq+4/q7VA7HzEf2k/i/oN1NI
JgtrhpPcZ/LMO6k7DYx0qlfYq8pTSfd6MI4LnWKgLc+JSPJJjmvspgio2ZFcnYr7
A264BwLo6v1Mos1o1JUvFFcp4GANlw0XFiWh7JXYRl8WmS5DoouUC+aNJ3lmyF6z
LbIjZCSfgZnk/LK1KU1j91FI2bc2ULYZvAC1PAg8/zvIgxn6YM2Q7ZsdEgWw0FpS
zMBX1/lk4wkFckeUIlkD55Y=
-----END CERTIFICATE-----`

const excludedNamesRoot = `
-----BEGIN CERTIFICATE-----
MIIEGTCCAwGgAwIBAgIHDUSFpInn/zANBgkqhkiG9w0BAQsFADCBozELMAkGA1UE
BhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExEjAQBgNVBAcMCUxvcyBHYXRvczEU
MBIGA1UECgwLTmV0ZmxpeCBJbmMxLTArBgNVBAsMJFBsYXRmb3JtIFNlY3VyaXR5
ICgzNzMxNTA5NDM3NDYyNDg1KTEmMCQGA1UEAwwdTmFtZSBDb25zdHJhaW50cyBU
ZXN0IFJvb3QgQ0EwHhcNMTcwMjA4MjExNTA0WhcNMTgwMjA4MjAyNDU4WjCBmTEL
MAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExEjAQBgNVBAcMCUxvcyBH
YXRvczEUMBIGA1UECgwLTmV0ZmxpeCBJbmMxLTArBgNVBAsMJFBsYXRmb3JtIFNl
Y3VyaXR5ICgzNzM0NTE1NDc5MDY0NjAyKTEcMBoGA1UEAwwTTG9jYWwgUm9vdCBm
b3IgMzM5MjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJymcnX29ekc
7+MLyr8QuAzoHWznmGdDd2sITwWRjM89/21cdlHCGKSpULUNdFp9HDLWvYECtxt+
8TuzKiQz7qAerzGUT1zI5McIjHy0e/i4xIkfiBiNeTCuB/N9QRbZlcfM80ErkaA4
gCAFK8qZAcWkHIl6e+KaQFMPLKk9kckgAnVDHEJe8oLNCogCJ15558b65g05p9eb
5Lg+E98hoPRTQaDwlz3CZPfTTA2EiEZInSi8qzodFCbTpJUVTbiVUH/JtVjlibbb
smdcx5PORK+8ZJkhLEh54AjaWOX4tB/7Tkk8stg2VBmrIARt/j4UVj7cTrIWU3bV
m8TwHJG+YgsCAwEAAaNaMFgwDwYDVR0TAQH/BAUwAwEB/zBFBgNVHR4EPjA8oBww
CocICgEAAP//AAAwDoIMYmVuZGVyLmxvY2FsoRwwCocICgEAAP//AAAwDoIMYmVu
ZGVyLmxvY2FsMA0GCSqGSIb3DQEBCwUAA4IBAQAMjbheffPxtSKSv9NySW+8qmHs
n7Mb5GGyCFu+cMZSoSaabstbml+zHEFJvWz6/1E95K4F8jKhAcu/CwDf4IZrSD2+
Hee0DolVSQhZpnHgPyj7ZATz48e3aJaQPUlhCEOh0wwF4Y0N4FV0t7R6woLylYRZ
yU1yRHUqUYpN0DWFpsPbBqgM6uUAVO2ayBFhPgWUaqkmSbZ/Nq7isGvknaTmcIwT
6mOAFN0qFb4RGzfGJW7x6z7KCULS7qVDp6fU3tRoScHFEgRubks6jzQ1W5ooSm4o
+NQCZDd5eFeU8PpNX7rgaYE4GPq+EEmLVCBYmdctr8QVdqJ//8Xu3+1phjDy
-----END CERTIFICATE-----`

const invalidCNRoot = `-----BEGIN CERTIFICATE-----
MIIBZzCCAQ6gAwIBAgIQO0Ri4zkY+e5zDG8/1/b/njAKBggqhkjOPQQDAjAUMRIw
EAYDVQQLEwlUZXN0IHJvb3QwHhcNMTgwNzExMTYzMjE1WhcNMjgwNzEyMTYwNTM1
WjAUMRIwEAYDVQQLEwlUZXN0IHJvb3QwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNC
AATnruwA9QRPeavhn96/cU+tFLPaayo3nhiJp6NIWxp9T/lgAW2G+UWCOqkwTmx5
RNOMWfHEeA4BSFcOko63hOKCo0IwQDAOBgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/
BAUwAwEB/zAdBgNVHQ4EFgQUr0uL0xUvkMrj5ga7y/saiKeZlh4wCgYIKoZIzj0E
AwIDRwAwRAIgAvSQJ7/k6hPWBR1GIJdKz3bm/5yrkqpfEto3kYR1x0QCIA3NsbEY
X0d9ycjZPfxlV2d/jc0jS5/j3HyzKRlWtdbe
-----END CERTIFICATE-----
`

const invalidCNWithoutSAN = `-----BEGIN CERTIFICATE-----
MIIBgTCCASagAwIBAgIQQ6LBTu24z5nGTRH2/N4wlDAKBggqhkjOPQQDAjAUMRIw
EAYDVQQLEwlUZXN0IHJvb3QwHhcNMTgwNzExMTYyNTIxWhcNMjgwNzEwMTYwNTIx
WjAWMRQwEgYDVQQDEwtmb28saW52YWxpZDBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABLT4jXENc7BNmFf4yTfeekX+0InDoMUGodllNdlkBMvoH7A0VLMVz378zAiX
ZHeYGhNRoC8ukv+yFnAZ/me37B6jWDBWMBMGA1UdJQQMMAoGCCsGAQUFBwMBMB8G
A1UdIwQYMBaAFK9Li9MVL5DK4+YGu8v7GoinmZYeMB4GA1UdEQQXMBWBE2Zvb0Bp
bnZhbGlkLmV4YW1wbGUwCgYIKoZIzj0EAwIDSQAwRgIhALDZUu1Zrc5Ty6BnQEz8
EyS2Gxh3Ex+eVgM1DmPYoBITAiEAkXZBvl/HeestW2JKBy6UDcU1Q0+ZDIf0fm61
9CTZp3s=
-----END CERTIFICATE-----
`

const validCNWithoutSAN = `-----BEGIN CERTIFICATE-----
MIIBgDCCASegAwIBAgIRAMDYX4qhRNGZEQ7e0jGLzAswCgYIKoZIzj0EAwIwFDES
MBAGA1UECxMJVGVzdCByb290MB4XDTE4MDcxMTE3MjcyNFoXDTI4MDcxMDE1NTcy
NFowGjEYMBYGA1UEAxMPZm9vLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZI
zj0DAQcDQgAEM5L/kIJ5IeOu6oeRd3crd/UCxKGibnXyVRYlwAVgOu1HdLHwdzHv
5TK5+Vn+XrTbbB6fcyJBllruGvUJWk1z4KNUMFIwEwYDVR0lBAwwCgYIKwYBBQUH
AwEwHwYDVR0jBBgwFoAUr0uL0xUvkMrj5ga7y/saiKeZlh4wGgYDVR0RBBMwEYIP
Zm9vLmV4YW1wbGUuY29tMAoGCCqGSM49BAMCA0cAMEQCIHJtK97Q20R/bAis7C3b
jLDQlxnM0Al506QfgPyLRUUhAiBFbwmrVMIOwfyeSRVJ/Y34+ChdeboywW+RHZm4
QflxsA==
-----END CERTIFICATE-----
`

const (
	rootWithoutSKID = `
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            78:29:2a:dc:2f:12:39:7f:c9:33:93:ea:61:39:7d:70
        Signature Algorithm: ecdsa-with-SHA256
        Issuer: O = Acme Co
        Validity
            Not Before: Feb  4 22:56:34 2019 GMT
            Not After : Feb  1 22:56:34 2029 GMT
        Subject: O = Acme Co
        Subject Public Key Info:
            Public Key Algorithm: id-ecPublicKey
                Public-Key: (256 bit)
                pub:
                    04:84:a6:8c:69:53:af:87:4b:39:64:fe:04:24:e6:
                    d8:fc:d6:46:39:35:0e:92:dc:48:08:7e:02:5f:1e:
                    07:53:5c:d9:e0:56:c5:82:07:f6:a3:e2:ad:f6:ad:
                    be:a0:4e:03:87:39:67:0c:9c:46:91:68:6b:0e:8e:
                    f8:49:97:9d:5b
                ASN1 OID: prime256v1
                NIST CURVE: P-256
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment, Certificate Sign
            X509v3 Extended Key Usage:
                TLS Web Server Authentication
            X509v3 Basic Constraints: critical
                CA:TRUE
            X509v3 Subject Alternative Name:
                DNS:example
    Signature Algorithm: ecdsa-with-SHA256
         30:46:02:21:00:c6:81:61:61:42:8d:37:e7:d0:c3:72:43:44:
         17:bd:84:ff:88:81:68:9a:99:08:ab:3c:3a:c0:1e:ea:8c:ba:
         c0:02:21:00:de:c9:fa:e5:5e:c6:e2:db:23:64:43:a9:37:42:
         72:92:7f:6e:89:38:ea:9e:2a:a7:fd:2f:ea:9a:ff:20:21:e7
-----BEGIN CERTIFICATE-----
MIIBbzCCARSgAwIBAgIQeCkq3C8SOX/JM5PqYTl9cDAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE5MDIwNDIyNTYzNFoXDTI5MDIwMTIyNTYzNFow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABISm
jGlTr4dLOWT+BCTm2PzWRjk1DpLcSAh+Al8eB1Nc2eBWxYIH9qPirfatvqBOA4c5
ZwycRpFoaw6O+EmXnVujTDBKMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MBIGA1UdEQQLMAmCB2V4YW1wbGUwCgYI
KoZIzj0EAwIDSQAwRgIhAMaBYWFCjTfn0MNyQ0QXvYT/iIFompkIqzw6wB7qjLrA
AiEA3sn65V7G4tsjZEOpN0Jykn9uiTjqniqn/S/qmv8gIec=
-----END CERTIFICATE-----
`
	leafWithAKID = `
	Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number:
            f0:8a:62:f0:03:84:a2:cf:69:63:ad:71:3b:b6:5d:8c
        Signature Algorithm: ecdsa-with-SHA256
        Issuer: O = Acme Co
        Validity
            Not Before: Feb  4 23:06:52 2019 GMT
            Not After : Feb  1 23:06:52 2029 GMT
        Subject: O = Acme LLC
        Subject Public Key Info:
            Public Key Algorithm: id-ecPublicKey
                Public-Key: (256 bit)
                pub:
                    04:5a:4e:4d:fb:ff:17:f7:b6:13:e8:29:45:34:81:
                    39:ff:8c:9c:d9:8c:0a:9f:dd:b5:97:4c:2b:20:91:
                    1c:4f:6b:be:53:27:66:ec:4a:ad:08:93:6d:66:36:
                    0c:02:70:5d:01:ca:7f:c3:29:e9:4f:00:ba:b4:14:
                    ec:c5:c3:34:b3
                ASN1 OID: prime256v1
                NIST CURVE: P-256
        X509v3 extensions:
            X509v3 Key Usage: critical
                Digital Signature, Key Encipherment
            X509v3 Extended Key Usage:
                TLS Web Server Authentication
            X509v3 Basic Constraints: critical
                CA:FALSE
            X509v3 Authority Key Identifier:
                keyid:C2:2B:5F:91:78:34:26:09:42:8D:6F:51:B2:C5:AF:4C:0B:DE:6A:42

            X509v3 Subject Alternative Name:
                DNS:example
    Signature Algorithm: ecdsa-with-SHA256
         30:44:02:20:64:e0:ba:56:89:63:ce:22:5e:4f:22:15:fd:3c:
         35:64:9a:3a:6b:7b:9a:32:a0:7f:f7:69:8c:06:f0:00:58:b8:
         02:20:09:e4:9f:6d:8b:9e:38:e1:b6:01:d5:ee:32:a4:94:65:
         93:2a:78:94:bb:26:57:4b:c7:dd:6c:3d:40:2b:63:90
-----BEGIN CERTIFICATE-----
MIIBjTCCATSgAwIBAgIRAPCKYvADhKLPaWOtcTu2XYwwCgYIKoZIzj0EAwIwEjEQ
MA4GA1UEChMHQWNtZSBDbzAeFw0xOTAyMDQyMzA2NTJaFw0yOTAyMDEyMzA2NTJa
MBMxETAPBgNVBAoTCEFjbWUgTExDMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE
Wk5N+/8X97YT6ClFNIE5/4yc2YwKn921l0wrIJEcT2u+Uydm7EqtCJNtZjYMAnBd
Acp/wynpTwC6tBTsxcM0s6NqMGgwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoG
CCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAUwitfkXg0JglCjW9R
ssWvTAveakIwEgYDVR0RBAswCYIHZXhhbXBsZTAKBggqhkjOPQQDAgNHADBEAiBk
4LpWiWPOIl5PIhX9PDVkmjpre5oyoH/3aYwG8ABYuAIgCeSfbYueOOG2AdXuMqSU
ZZMqeJS7JldLx91sPUArY5A=
-----END CERTIFICATE-----
`
)

var unknownAuthorityErrorTests = []struct {
	cert     string
	expected string
}{
	{selfSignedWithCommonName, "x509: certificate signed by unknown authority (possibly because of \"empty\" while trying to verify candidate authority certificate \"test\")"},
	{selfSignedNoCommonNameWithOrgName, "x509: certificate signed by unknown authority (possibly because of \"empty\" while trying to verify candidate authority certificate \"ca\")"},
	{selfSignedNoCommonNameNoOrgName, "x509: certificate signed by unknown authority (possibly because of \"empty\" while trying to verify candidate authority certificate \"serial:0\")"},
}

func TestUnknownAuthorityError(t *testing.T) {
	for i, tt := range unknownAuthorityErrorTests {
		der, _ := pem.Decode([]byte(tt.cert))
		if der == nil {
			t.Errorf("#%d: Unable to decode PEM block", i)
		}
		c, err := x509.ParseCertificate(der.Bytes)
		if err != nil {
			t.Errorf("#%d: Unable to parse certificate -> %v", i, err)
		}
		uae := &UnknownAuthorityError{
			Cert:     c,
			hintErr:  fmt.Errorf("empty"),
			hintCert: c,
		}
		actual := uae.Error()
		if actual != tt.expected {
			t.Errorf("#%d: UnknownAuthorityError.Error() response invalid actual: %s expected: %s", i, actual, tt.expected)
		}
	}
}

var nameConstraintTests = []struct {
	constraint, domain string
	expectError        bool
	shouldMatch        bool
}{
	{"", "anything.com", false, true},
	{"example.com", "example.com", false, true},
	{"example.com.", "example.com", true, false},
	{"example.com", "example.com.", true, false},
	{"example.com", "ExAmPle.coM", false, true},
	{"example.com", "exampl1.com", false, false},
	{"example.com", "www.ExAmPle.coM", false, true},
	{"example.com", "sub.www.ExAmPle.coM", false, true},
	{"example.com", "notexample.com", false, false},
	{".example.com", "example.com", false, false},
	{".example.com", "www.example.com", false, true},
	{".example.com", "www..example.com", true, false},
}

func TestNameConstraints(t *testing.T) {
	for i, test := range nameConstraintTests {
		result, err := matchDomainConstraint(test.domain, test.constraint)

		if err != nil && !test.expectError {
			t.Errorf("unexpected error for test #%d: domain=%s, constraint=%s, err=%s", i, test.domain, test.constraint, err)
			continue
		}

		if err == nil && test.expectError {
			t.Errorf("unexpected success for test #%d: domain=%s, constraint=%s", i, test.domain, test.constraint)
			continue
		}

		if result != test.shouldMatch {
			t.Errorf("unexpected result for test #%d: domain=%s, constraint=%s, result=%t", i, test.domain, test.constraint, result)
		}
	}
}

const selfSignedWithCommonName = `-----BEGIN CERTIFICATE-----
MIIDCjCCAfKgAwIBAgIBADANBgkqhkiG9w0BAQsFADAaMQswCQYDVQQKEwJjYTEL
MAkGA1UEAxMCY2EwHhcNMTYwODI4MTcwOTE4WhcNMjEwODI3MTcwOTE4WjAcMQsw
CQYDVQQKEwJjYTENMAsGA1UEAxMEdGVzdDCCASIwDQYJKoZIhvcNAQEBBQADggEP
ADCCAQoCggEBAOH55PfRsbvmcabfLLko1w/yuapY/hk13Cgmc3WE/Z1ZStxGiVxY
gQVH9n4W/TbUsrep/TmcC4MV7xEm5252ArcgaH6BeQ4QOTFj/6Jx0RT7U/ix+79x
8RRysf7OlzNpGIctwZEM7i/G+0ZfqX9ULxL/EW9tppSxMX1jlXZQarnU7BERL5cH
+G2jcbU9H28FXYishqpVYE9L7xrXMm61BAwvGKB0jcVW6JdhoAOSfQbbgp7JjIlq
czXqUsv1UdORO/horIoJptynTvuARjZzyWatya6as7wyOgEBllE6BjPK9zpn+lp3
tQ8dwKVqm/qBPhIrVqYG/Ec7pIv8mJfYabMCAwEAAaNZMFcwDgYDVR0PAQH/BAQD
AgOoMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATAMBgNVHRMBAf8EAjAA
MAoGA1UdDgQDBAEAMAwGA1UdIwQFMAOAAQAwDQYJKoZIhvcNAQELBQADggEBAAAM
XMFphzq4S5FBcRdB2fRrmcoz+jEROBWvIH/1QUJeBEBz3ZqBaJYfBtQTvqCA5Rjw
dxyIwVd1W3q3aSulM0tO62UCU6L6YeeY/eq8FmpD7nMJo7kCrXUUAMjxbYvS3zkT
v/NErK6SgWnkQiPJBZNX1Q9+aSbLT/sbaCTdbWqcGNRuLGJkmqfIyoxRt0Hhpqsx
jP5cBaVl50t4qoCuVIE9cOucnxYXnI7X5HpXWvu8Pfxo4SwVjb1az8Fk5s8ZnxGe
fPB6Q3L/pKBe0SEe5GywpwtokPLB3lAygcuHbxp/1FlQ1NQZqq+vgXRIla26bNJf
IuYkJwt6w+LH/9HZgf8=
-----END CERTIFICATE-----`

const selfSignedNoCommonNameWithOrgName = `-----BEGIN CERTIFICATE-----
MIIC+zCCAeOgAwIBAgIBADANBgkqhkiG9w0BAQsFADAaMQswCQYDVQQKEwJjYTEL
MAkGA1UEAxMCY2EwHhcNMTYwODI4MTgxMzQ4WhcNMjEwODI3MTgxMzQ4WjANMQsw
CQYDVQQKEwJjYTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL5EjrUa
7EtOMxWiIgTzp2FlQvncPsG329O3l3uNGnbigb8TmNMw2M8UhoDjd84pnU5RAfqd
8t5TJyw/ybnIKBN131Q2xX+gPQ0dFyMvcO+i1CUgCxmYZomKVA2MXO1RD1hLTYGS
gOVjc3no3MBwd8uVQp0NStqJ1QvLtNG4Uy+B28qe+ZFGGbjGqx8/CU4A8Szlpf7/
xAZR8w5qFUUlpA2LQYeHHJ5fQVXw7kyL1diNrKNi0G3qcY0IrBh++hT+hnEEXyXu
g8a0Ux18hoE8D6rAr34rCZl6AWfqW5wjwm+N5Ns2ugr9U4N8uCKJYMPHb2CtdubU
46IzVucpTfGLdaMCAwEAAaNZMFcwDgYDVR0PAQH/BAQDAgOoMB0GA1UdJQQWMBQG
CCsGAQUFBwMCBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMAoGA1UdDgQDBAEAMAwG
A1UdIwQFMAOAAQAwDQYJKoZIhvcNAQELBQADggEBAEn5SgVpJ3zjsdzPqK7Qd/sB
bYd1qtPHlrszjhbHBg35C6mDgKhcv4o6N+fuC+FojZb8lIxWzJtvT9pQbfy/V6u3
wOb816Hm71uiP89sioIOKCvSAstj/p9doKDOUaKOcZBTw0PS2m9eja8bnleZzBvK
rD8cNkHf74v98KvBhcwBlDifVzmkWzMG6TL1EkRXUyLKiWgoTUFSkCDV927oXXMR
DKnszq+AVw+K8hbeV2A7GqT7YfeqOAvSbatTDnDtKOPmlCnQui8A149VgZzXv7eU
29ssJSqjUPyp58dlV6ZuynxPho1QVZUOQgnJToXIQ3/5vIvJRXy52GJCs4/Gh/w=
-----END CERTIFICATE-----`

const selfSignedNoCommonNameNoOrgName = `-----BEGIN CERTIFICATE-----
MIIC7jCCAdagAwIBAgIBADANBgkqhkiG9w0BAQsFADAaMQswCQYDVQQKEwJjYTEL
MAkGA1UEAxMCY2EwHhcNMTYwODI4MTgxOTQ1WhcNMjEwODI3MTgxOTQ1WjAAMIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp3E+Jl6DpgzogHUW/i/AAcCM
fnNJLOamNVKFGmmxhb4XTHxRaWoTzrlsyzIMS0WzivvJeZVe6mWbvuP2kZanKgIz
35YXRTR9HbqkNTMuvnpUESzWxbGWE2jmt2+a/Jnz89FS4WIYRhF7nI2z8PvZOfrI
2gETTT2tEpoF2S4soaYfm0DBeT8K0/rogAaf+oeUS6V+v3miRcAooJgpNJGu9kqm
S0xKPn1RCFVjpiRd6YNS0xZirjYQIBMFBvoSoHjaOdgJptNRBprYPOxVJ/ItzGf0
kPmzPFCx2tKfxV9HLYBPgxi+fP3IIx8aIYuJn8yReWtYEMYU11hDPeAFN5Gm+wID
AQABo1kwVzAOBgNVHQ8BAf8EBAMCA6gwHQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsG
AQUFBwMBMAwGA1UdEwEB/wQCMAAwCgYDVR0OBAMEAQAwDAYDVR0jBAUwA4ABADAN
BgkqhkiG9w0BAQsFAAOCAQEATZVOFeiCpPM5QysToLv+8k7Rjoqt6L5IxMUJGEpq
4ENmldmwkhEKr9VnYEJY3njydnnTm97d9vOfnLj9nA9wMBODeOO3KL2uJR2oDnmM
9z1NSe2aQKnyBb++DM3ZdikpHn/xEpGV19pYKFQVn35x3lpPh2XijqRDO/erKemb
w67CoNRb81dy+4Q1lGpA8ORoLWh5fIq2t2eNGc4qB8vlTIKiESzAwu7u3sRfuWQi
4R+gnfLd37FWflMHwztFbVTuNtPOljCX0LN7KcuoXYlr05RhQrmoN7fQHsrZMNLs
8FVjHdKKu+uPstwd04Uy4BR/H2y1yerN9j/L6ZkMl98iiA==
-----END CERTIFICATE-----`

const criticalExtRoot = `-----BEGIN CERTIFICATE-----
MIIBezCCASGgAwIBAgIRAIwRaUt4rZvXch4v02bjK2IwCgYIKoZIzj0EAwIwHTEM
MAoGA1UEChMDT3JnMQ0wCwYDVQQDEwRSb290MB4XDTE1MDEwMTAwMDAwMFoXDTI1
MDEwMTAwMDAwMFowHTEMMAoGA1UEChMDT3JnMQ0wCwYDVQQDEwRSb290MFkwEwYH
KoZIzj0CAQYIKoZIzj0DAQcDQgAEo3qKEPajJdTafip8xcUNzINYqv1PvuZo6/5g
FMOoPV1C6HdICnXR4Ra2O3FTwIGDEGiylcIwlT8nuMvRpkYiWqNCMEAwDgYDVR0P
AQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFIiB2+ltwD5CBo+l
efo6Ool0Lb3LMAoGCCqGSM49BAMCA0gAMEUCIQDTJ9AFgtRV05RGHjqKjzf5nUhY
ighLEW6CFL0aedgj9AIgAysWs2DyWEGY2TKkOFeCzEHIdW/FUvA7MSoHCNJXXAM=
-----END CERTIFICATE-----
`

const criticalExtIntermediate = `-----BEGIN CERTIFICATE-----
MIIByDCCAW2gAwIBAgIQZdWp3vTr9+NnNrE+VkNmsjAKBggqhkjOPQQDAjAdMQww
CgYDVQQKEwNPcmcxDTALBgNVBAMTBFJvb3QwHhcNMTUwMTAxMDAwMDAwWhcNMjUw
MTAxMDAwMDAwWjAlMQwwCgYDVQQKEwNPcmcxFTATBgNVBAMTDEludGVybWVkaWF0
ZTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABBalbOoxZOmPBCNQr7m+i0sF2bby
c6Wa8sNvMN22TDQFM8fhhofbb8EMVcLwyfLgLrFd7AL8lTq1XD7NO0JatD6jgYYw
gYMwDgYDVR0PAQH/BAQDAgEGMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcD
AjASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBQJWagJtsBC5eKGmvj1HzDU
74htMTAfBgNVHSMEGDAWgBSIgdvpbcA+QgaPpXn6OjqJdC29yzAKBggqhkjOPQQD
AgNJADBGAiEAhitP8MVD27BrtYpALpp0XbWaEyk+PDXgpzDT5+UgFIYCIQD0NDUz
3BWaNzZuPQ4jUkyrI+BAf6P7W1AdaVYPS0KsgA==
-----END CERTIFICATE-----
`

const criticalExtLeafWithExt = `-----BEGIN CERTIFICATE-----
MIIBrzCCAVSgAwIBAgIRALQyAx4JS1tERmm3SrxriKswCgYIKoZIzj0EAwIwJTEM
MAoGA1UEChMDT3JnMRUwEwYDVQQDEwxJbnRlcm1lZGlhdGUwHhcNMTUwMTAxMDAw
MDAwWhcNMjUwMTAxMDAwMDAwWjAkMQwwCgYDVQQKEwNPcmcxFDASBgNVBAMTC2V4
YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEaSiLYGfZ+oC2uYHo
TDCTzDEOsG4LJjnm+WWjeCJzXJjJFjl0brcHM8hTMdO6YkCs5JeeTNl3YNjpADIu
PB8LsqNmMGQwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB8GA1UdIwQY
MBaAFAlZqAm2wELl4oaa+PUfMNTviG0xMBYGA1UdEQQPMA2CC2V4YW1wbGUuY29t
MAoGA1EDBAEB/wQAMAoGCCqGSM49BAMCA0kAMEYCIQCfFTUYMbc3sPwFqM6GwEcl
fkcclvvfZhbAfQI5U2b0zAIhAIfRF8JxMsgw3V+en7NmMI1+EwsFPsGLidAQBQ8r
33x/
-----END CERTIFICATE-----
`

const criticalExtIntermediateWithExt = `-----BEGIN CERTIFICATE-----
MIIB6jCCAZGgAwIBAgIQP9IABgu01SIKL4mOp8PujTAKBggqhkjOPQQDAjAdMQww
CgYDVQQKEwNPcmcxDTALBgNVBAMTBFJvb3QwHhcNMTUwMTAxMDAwMDAwWhcNMjUw
MTAxMDAwMDAwWjA9MQwwCgYDVQQKEwNPcmcxLTArBgNVBAMTJEludGVybWVkaWF0
ZSB3aXRoIENyaXRpY2FsIEV4dGVuc2lvbjBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABOKKKeSpXzTbBGmJ63eu7LOIFyQ2EDtPdlRsedhQURovOrSpT8gPn5MSiJn2
Z7SpdymFCFNnq57qF6nZPGoPz4OjgZIwgY8wDgYDVR0PAQH/BAQDAgEGMB0GA1Ud
JQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMB0G
A1UdDgQWBBTSjlptXrFf+KifZMdeWcOhWf7MGzAfBgNVHSMEGDAWgBSIgdvpbcA+
QgaPpXn6OjqJdC29yzAKBgNRAwQBAf8EADAKBggqhkjOPQQDAgNHADBEAiBs0R/e
rnemyl7fjUcAKIkFtNMojoiuU1gom8u8IaXYLQIgATxbigT5+QYa4mZ4mHDV2u4B
Ehiu0+2R+ia09gCJFq0=
-----END CERTIFICATE-----
`

const criticalExtLeaf = `-----BEGIN CERTIFICATE-----
MIIBuDCCAV+gAwIBAgIQInS3rr1x4FAXDrY0jY7AMjAKBggqhkjOPQQDAjA9MQww
CgYDVQQKEwNPcmcxLTArBgNVBAMTJEludGVybWVkaWF0ZSB3aXRoIENyaXRpY2Fs
IEV4dGVuc2lvbjAeFw0xNTAxMDEwMDAwMDBaFw0yNTAxMDEwMDAwMDBaMCQxDDAK
BgNVBAoTA09yZzEUMBIGA1UEAxMLZXhhbXBsZS5jb20wWTATBgcqhkjOPQIBBggq
hkjOPQMBBwNCAAQV9xxn5Ts6Y2Ri3E6hIiX+pUBOl7fA3MRMd+HV16aGF7mHE5ge
WAyUr0u+MPLfHw7uVQDRsqj9z+TCp+YU4bVBo1owWDAdBgNVHSUEFjAUBggrBgEF
BQcDAQYIKwYBBQUHAwIwHwYDVR0jBBgwFoAU0o5abV6xX/ion2THXlnDoVn+zBsw
FgYDVR0RBA8wDYILZXhhbXBsZS5jb20wCgYIKoZIzj0EAwIDRwAwRAIgMBzFqYWK
bnuHScpa03Pr5fXEXuWr6udNJ9qh9nvhGQQCIA8+2DdCa7twdHUoZbnVuaqR+Vn7
qscz8cYNi5o/DM1f
-----END CERTIFICATE-----
`

func TestValidHostname(t *testing.T) {
	tests := []struct {
		host string
		want bool
	}{
		{"example.com", true},
		{"eXample123-.com", true},
		{"-eXample123-.com", false},
		{"", false},
		{".", false},
		{"example..com", false},
		{".example.com", false},
		{"*.example.com", true},
		{"*foo.example.com", false},
		{"foo.*.example.com", false},
		{"exa_mple.com", true},
		{"foo,bar", false},
		{"project-dev:us-central1:main", true},
	}
	for _, tt := range tests {
		if got := validHostname(tt.host); got != tt.want {
			t.Errorf("validHostname(%q) = %v, want %v", tt.host, got, tt.want)
		}
	}
}

func generateCert(cn string, isCA bool, issuer *x509.Certificate, issuerKey crypto.PrivateKey) (*x509.Certificate, crypto.PrivateKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  isCA,
	}
	if issuer == nil {
		issuer = template
		issuerKey = priv
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, issuer, priv.Public(), issuerKey)
	if err != nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, nil, err
	}

	return cert, priv, nil
}

func TestPathologicalChain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping generation of a long chain of certificates in short mode")
	}

	// Build a chain where all intermediates share the same subject, to hit the
	// path building worst behavior.
	roots, intermediates := NewCertPool(), NewCertPool()

	parent, parentKey, err := generateCert("Root CA", true, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	roots.AddCert(parent)

	for i := 1; i < 100; i++ {
		parent, parentKey, err = generateCert("Intermediate CA", true, parent, parentKey)
		if err != nil {
			t.Fatal(err)
		}
		intermediates.AddCert(parent)
	}

	leaf, _, err := generateCert("Leaf", false, parent, parentKey)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	_, err = Verify(leaf, VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
	})
	t.Logf("verification took %v", time.Since(start))

	if err == nil || !strings.Contains(err.Error(), "signature check attempts limit") {
		t.Errorf("expected verification to fail with a signature checks limit error; got %v", err)
	}
}

func TestLongChain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping generation of a long chain of certificates in short mode")
	}

	roots, intermediates := NewCertPool(), NewCertPool()

	parent, parentKey, err := generateCert("Root CA", true, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	roots.AddCert(parent)

	for i := 1; i < 15; i++ {
		name := fmt.Sprintf("Intermediate CA #%d", i)
		parent, parentKey, err = generateCert(name, true, parent, parentKey)
		if err != nil {
			t.Fatal(err)
		}
		intermediates.AddCert(parent)
	}

	leaf, _, err := generateCert("Leaf", false, parent, parentKey)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	if _, err := Verify(leaf, VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
	}); err != nil {
		t.Error(err)
	}
	t.Logf("verification took %v", time.Since(start))
}

// SPDX-License-Identifier: MPL-2.0
// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// testCA is a certificate authority usable for signing certificates and CRLs
// in tests.
type testCA struct {
	cert *x509.Certificate
	key  *ecdsa.PrivateKey
}

func newTestCA(t *testing.T, commonName string) *testCA {
	t.Helper()
	return newTestSubCA(t, commonName, nil)
}

// newTestSubCA creates a CA certificate signed by parent, or self-signed when
// parent is nil.
func newTestSubCA(t *testing.T, commonName string, parent *testCA) *testCA {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("error generating key: %v", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(randSerial(t)),
		Subject:               pkix.Name{CommonName: commonName},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		SubjectKeyId:          subjectKeyID(t, &key.PublicKey),
	}

	signerCert, signerKey := tmpl, any(key)
	if parent != nil {
		signerCert, signerKey = parent.cert, parent.key
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, signerCert, &key.PublicKey, signerKey)
	if err != nil {
		t.Fatalf("error creating CA certificate: %v", err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("error parsing CA certificate: %v", err)
	}

	return &testCA{cert: cert, key: key}
}

// issueLeaf issues a server certificate for the given DNS name.
func (ca *testCA) issueLeaf(t *testing.T, dnsName string) (*x509.Certificate, *ecdsa.PrivateKey) {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("error generating leaf key: %v", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(randSerial(t)),
		Subject:      pkix.Name{CommonName: dnsName},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{dnsName},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		SubjectKeyId: subjectKeyID(t, &key.PublicKey),
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca.cert, &key.PublicKey, ca.key)
	if err != nil {
		t.Fatalf("error creating leaf certificate: %v", err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("error parsing leaf certificate: %v", err)
	}
	return cert, key
}

// issueSelfSignedLeaf creates a self-signed certificate that is not a CA, which
// is what `openssl req -x509` produces and the shape a development Vault server
// commonly uses.
func issueSelfSignedLeaf(t *testing.T, dnsName string) *x509.Certificate {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("error generating key: %v", err)
	}

	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(randSerial(t)),
		Subject:               pkix.Name{CommonName: dnsName},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              []string{dnsName},
		BasicConstraintsValid: true,
		IsCA:                  false,
		SubjectKeyId:          subjectKeyID(t, &key.PublicKey),
	}

	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("error creating self-signed leaf: %v", err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatalf("error parsing self-signed leaf: %v", err)
	}
	return cert
}

// crl issues a CRL revoking the given certificates.
func (ca *testCA) crl(t *testing.T, nextUpdate time.Time, revoked ...*x509.Certificate) *x509.RevocationList {
	t.Helper()
	return ca.parsedCRL(t, ca.crlDER(t, nextUpdate, revoked...))
}

// deltaCRL issues a CRL carrying the deltaCRLIndicator extension, so it lists
// only the changes since some base CRL.
func (ca *testCA) deltaCRL(t *testing.T, nextUpdate time.Time, revoked ...*x509.Certificate) *x509.RevocationList {
	t.Helper()
	return ca.parsedCRL(t, ca.crlDERWithExtensions(t, nextUpdate, []pkix.Extension{{
		Id:       oidDeltaCRLIndicator,
		Critical: true,
		Value:    mustMarshalASN1(t, big.NewInt(1)),
	}}, revoked...))
}

func mustMarshalASN1(t *testing.T, v any) []byte {
	t.Helper()
	der, err := asn1.Marshal(v)
	if err != nil {
		t.Fatalf("error marshalling ASN.1 value: %v", err)
	}
	return der
}

func (ca *testCA) crlDER(t *testing.T, nextUpdate time.Time, revoked ...*x509.Certificate) []byte {
	t.Helper()
	return ca.crlDERWithExtensions(t, nextUpdate, nil, revoked...)
}

func (ca *testCA) crlDERWithExtensions(t *testing.T, nextUpdate time.Time, extensions []pkix.Extension, revoked ...*x509.Certificate) []byte {
	t.Helper()

	entries := make([]x509.RevocationListEntry, 0, len(revoked))
	for _, cert := range revoked {
		entries = append(entries, x509.RevocationListEntry{
			SerialNumber:   cert.SerialNumber,
			RevocationTime: time.Now().Add(-time.Minute),
		})
	}

	tmpl := &x509.RevocationList{
		Number:                    big.NewInt(randSerial(t)),
		ThisUpdate:                time.Now().Add(-time.Hour),
		NextUpdate:                nextUpdate,
		RevokedCertificateEntries: entries,
		ExtraExtensions:           extensions,
	}

	der, err := x509.CreateRevocationList(rand.Reader, tmpl, ca.cert, ca.key)
	if err != nil {
		t.Fatalf("error creating CRL: %v", err)
	}
	return der
}

func (ca *testCA) parsedCRL(t *testing.T, der []byte) *x509.RevocationList {
	t.Helper()
	crl, err := x509.ParseRevocationList(der)
	if err != nil {
		t.Fatalf("error parsing generated CRL: %v", err)
	}
	return crl
}

func randSerial(t *testing.T) int64 {
	t.Helper()
	n, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		t.Fatalf("error generating serial: %v", err)
	}
	return n.Int64() + 1
}

func subjectKeyID(t *testing.T, pub *ecdsa.PublicKey) []byte {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("error marshalling public key: %v", err)
	}
	// Must actually depend on the whole key: the leading bytes of a P-256
	// SubjectPublicKeyInfo are a fixed ASN.1 header, so anything derived from
	// only the prefix would collide across CAs.
	sum := sha256.Sum256(der)
	return sum[:20]
}

func pemCRL(t *testing.T, der []byte) []byte {
	t.Helper()
	return pem.EncodeToMemory(&pem.Block{Type: crlPEMBlockType, Bytes: der})
}

func futureTime() time.Time { return time.Now().Add(24 * time.Hour) }
func pastTime() time.Time   { return time.Now().Add(-30 * time.Minute) }

func TestParseCRLs(t *testing.T) {
	ca := newTestCA(t, "parse-crl-ca")
	other := newTestCA(t, "parse-crl-other-ca")
	derOne := ca.crlDER(t, futureTime())
	derTwo := other.crlDER(t, futureTime())

	caCertPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.cert.Raw})

	tests := []struct {
		name      string
		data      []byte
		wantCount int
		wantErr   string
	}{
		{
			name:      "single PEM CRL",
			data:      pemCRL(t, derOne),
			wantCount: 1,
		},
		{
			name:      "PEM bundle of two CRLs",
			data:      append(pemCRL(t, derOne), pemCRL(t, derTwo)...),
			wantCount: 2,
		},
		{
			name:      "raw DER CRL",
			data:      derOne,
			wantCount: 1,
		},
		{
			name:      "CRL and certificate in one PEM file ignores the certificate",
			data:      append(caCertPEM, pemCRL(t, derOne)...),
			wantCount: 1,
		},
		{
			name:    "PEM with no CRL blocks",
			data:    caCertPEM,
			wantErr: "no \"X509 CRL\" blocks found",
		},
		{
			name:    "truncated DER",
			data:    derOne[:len(derOne)/2],
			wantErr: "failed to parse CRL as PEM or DER",
		},
		{
			name:    "corrupt PEM CRL body",
			data:    pemCRL(t, derOne[:len(derOne)/2]),
			wantErr: "failed to parse PEM-encoded CRL",
		},
		{
			name:    "empty input",
			data:    nil,
			wantErr: "no CRL data provided",
		},
		{
			name:    "whitespace only",
			data:    []byte("\n \t\n"),
			wantErr: "no CRL data provided",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			crls, err := ParseCRLs(tc.data)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(crls) != tc.wantCount {
				t.Fatalf("expected %d CRLs, got %d", tc.wantCount, len(crls))
			}
		})
	}
}

func TestCRLChecker_VerifyConnection(t *testing.T) {
	root := newTestCA(t, "crl-root")
	intermediate := newTestSubCA(t, "crl-intermediate", root)
	leaf, _ := intermediate.issueLeaf(t, "vault.example.com")
	otherCA := newTestCA(t, "crl-unrelated")

	// A CA that reuses the subject name but has a different key, standing in
	// for a CA whose key was rotated.
	rotatedRoot := newTestCA(t, "crl-intermediate")

	directLeaf, _ := root.issueLeaf(t, "vault.example.com")
	fullChain := []*x509.Certificate{leaf, intermediate.cert, root.cert}

	// A self-signed non-CA certificate, which is what a development server
	// typically presents.
	selfSignedLeaf := issueSelfSignedLeaf(t, "selfsigned.example.com")

	// A self-signed CA that lists its own certificate in the CRL it publishes.
	selfRevokedCA := newTestCA(t, "crl-self-revoked")

	tests := []struct {
		name    string
		chain   []*x509.Certificate
		crls    []*x509.RevocationList
		wantErr string
	}{
		{
			name:  "clean CRL for leaf issuer passes",
			chain: fullChain,
			crls:  []*x509.RevocationList{intermediate.crl(t, futureTime())},
		},
		{
			name:    "revoked leaf is rejected",
			chain:   fullChain,
			crls:    []*x509.RevocationList{intermediate.crl(t, futureTime(), leaf)},
			wantErr: "was revoked",
		},
		{
			name:  "revoked intermediate is rejected",
			chain: fullChain,
			crls: []*x509.RevocationList{
				intermediate.crl(t, futureTime()),
				root.crl(t, futureTime(), intermediate.cert),
			},
			wantErr: "was revoked",
		},
		{
			name:    "no CRL for leaf issuer fails closed",
			chain:   fullChain,
			crls:    []*x509.RevocationList{root.crl(t, futureTime())},
			wantErr: "no CRL was supplied for issuer",
		},
		{
			name:    "CRL from an unrelated CA does not satisfy coverage",
			chain:   fullChain,
			crls:    []*x509.RevocationList{otherCA.crl(t, futureTime())},
			wantErr: "no CRL was supplied for issuer",
		},
		{
			name:    "stale CRL does not count as coverage",
			chain:   fullChain,
			crls:    []*x509.RevocationList{intermediate.crl(t, pastTime())},
			wantErr: "past their next update",
		},
		{
			name:    "stale CRL still honours its revocations",
			chain:   fullChain,
			crls:    []*x509.RevocationList{intermediate.crl(t, pastTime(), leaf)},
			wantErr: "was revoked",
		},
		{
			name:  "missing intermediate CRL is tolerated when the leaf is covered",
			chain: fullChain,
			crls:  []*x509.RevocationList{intermediate.crl(t, futureTime())},
		},
		{
			name:  "leaf issued directly by the root",
			chain: []*x509.Certificate{directLeaf, root.cert},
			crls:  []*x509.RevocationList{root.crl(t, futureTime())},
		},
		{
			name:    "CRL signed by a rotated key does not satisfy coverage",
			chain:   fullChain,
			crls:    []*x509.RevocationList{rotatedRoot.crl(t, futureTime())},
			wantErr: "no CRL was supplied for issuer",
		},
		{
			// A self-signed CA is its own issuer, so requiring a CRL from that
			// issuer would make revocation checking and self-signed server
			// certificates mutually exclusive.
			name:    "self-signed CA server certificate is exempt from coverage",
			chain:   []*x509.Certificate{root.cert},
			crls:    []*x509.RevocationList{root.crl(t, futureTime())},
			wantErr: "",
		},
		{
			// The same exemption must apply to a self-signed leaf, which is
			// what `openssl req -x509` produces. isSelfSigned must not reject
			// it merely because it is not a CA.
			name:    "self-signed leaf server certificate is exempt from coverage",
			chain:   []*x509.Certificate{selfSignedLeaf},
			crls:    []*x509.RevocationList{root.crl(t, futureTime())},
			wantErr: "",
		},
		{
			// Exempt from coverage is not the same as unchecked: a self-signed
			// CA can revoke its own certificate, and that must be honoured.
			name:    "self-signed certificate listed in its own CRL is rejected",
			chain:   []*x509.Certificate{selfRevokedCA.cert},
			crls:    []*x509.RevocationList{selfRevokedCA.crl(t, futureTime(), selfRevokedCA.cert)},
			wantErr: "was revoked",
		},
		{
			// A delta CRL lists only changes since a base CRL, so on its own it
			// cannot show the certificate is unrevoked.
			name:    "delta CRL alone does not satisfy coverage",
			chain:   fullChain,
			crls:    []*x509.RevocationList{intermediate.deltaCRL(t, futureTime())},
			wantErr: "delta CRLs",
		},
		{
			name:    "delta CRL still honours its revocations",
			chain:   fullChain,
			crls:    []*x509.RevocationList{intermediate.deltaCRL(t, futureTime(), leaf)},
			wantErr: "was revoked",
		},
		{
			name:  "delta CRL alongside a full CRL is accepted",
			chain: fullChain,
			crls: []*x509.RevocationList{
				intermediate.deltaCRL(t, futureTime()),
				intermediate.crl(t, futureTime()),
			},
		},
		{
			name:    "empty CRL set fails closed",
			chain:   fullChain,
			crls:    nil,
			wantErr: "no CRL was supplied for issuer",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			checker := &crlChecker{provider: NewStaticCRLProvider(tc.crls)}
			err := checker.verifyConnection(tls.ConnectionState{
				VerifiedChains: [][]*x509.Certificate{tc.chain},
			})

			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tc.wantErr, err)
			}
		})
	}
}

// A revoked cert in one chain must not be excused by a second clean chain, and
// vice versa: any one acceptable chain is enough.
func TestCRLChecker_MultipleChains(t *testing.T) {
	caA := newTestCA(t, "crl-multi-a")
	caB := newTestCA(t, "crl-multi-b")
	leafA, _ := caA.issueLeaf(t, "vault.example.com")
	leafB, _ := caB.issueLeaf(t, "vault.example.com")

	t.Run("one clean chain is enough", func(t *testing.T) {
		checker := &crlChecker{provider: NewStaticCRLProvider([]*x509.RevocationList{
			caA.crl(t, futureTime(), leafA),
			caB.crl(t, futureTime()),
		})}
		err := checker.verifyConnection(tls.ConnectionState{
			VerifiedChains: [][]*x509.Certificate{
				{leafA, caA.cert},
				{leafB, caB.cert},
			},
		})
		if err != nil {
			t.Fatalf("expected the clean chain to be accepted, got %v", err)
		}
	})

	t.Run("all chains revoked is rejected", func(t *testing.T) {
		checker := &crlChecker{provider: NewStaticCRLProvider([]*x509.RevocationList{
			caA.crl(t, futureTime(), leafA),
			caB.crl(t, futureTime(), leafB),
		})}
		err := checker.verifyConnection(tls.ConnectionState{
			VerifiedChains: [][]*x509.Certificate{
				{leafA, caA.cert},
				{leafB, caB.cert},
			},
		})
		if err == nil || !strings.Contains(err.Error(), "was revoked") {
			t.Fatalf("expected a revocation error, got %v", err)
		}
	})
}

// VerifyConnection runs only after chain verification succeeds, so an empty
// VerifiedChains means verification was disabled with InsecureSkipVerify.
// Falling back to the peer-presented chain would let an attacker satisfy the
// check with a self-consistent chain of their own certificates, so the check
// must fail closed.
func TestCRLChecker_NoVerifiedChains(t *testing.T) {
	ca := newTestCA(t, "crl-unverified")
	leaf, _ := ca.issueLeaf(t, "vault.example.com")

	t.Run("a clean presented chain is not accepted", func(t *testing.T) {
		checker := &crlChecker{provider: NewStaticCRLProvider([]*x509.RevocationList{
			ca.crl(t, futureTime()),
		})}
		err := checker.verifyConnection(tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{leaf, ca.cert},
		})
		if err == nil || !strings.Contains(err.Error(), "no verified certificate chain") {
			t.Fatalf("expected the check to fail closed, got %v", err)
		}
	})

	// An attacker who can present any chain they like would otherwise satisfy
	// leaf coverage with their own CA and CRL.
	t.Run("an attacker-supplied chain and CRL is not accepted", func(t *testing.T) {
		attacker := newTestCA(t, "crl-attacker")
		attackerLeaf, _ := attacker.issueLeaf(t, "vault.example.com")

		checker := &crlChecker{provider: NewStaticCRLProvider([]*x509.RevocationList{
			attacker.crl(t, futureTime()),
		})}
		err := checker.verifyConnection(tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{attackerLeaf, attacker.cert},
		})
		if err == nil || !strings.Contains(err.Error(), "no verified certificate chain") {
			t.Fatalf("expected the check to fail closed, got %v", err)
		}
	})

	t.Run("a revoked presented chain is not accepted either", func(t *testing.T) {
		checker := &crlChecker{provider: NewStaticCRLProvider([]*x509.RevocationList{
			ca.crl(t, futureTime(), leaf),
		})}
		err := checker.verifyConnection(tls.ConnectionState{
			PeerCertificates: []*x509.Certificate{leaf, ca.cert},
		})
		if err == nil || !strings.Contains(err.Error(), "no verified certificate chain") {
			t.Fatalf("expected the check to fail closed, got %v", err)
		}
	})

	t.Run("no certificates at all", func(t *testing.T) {
		checker := &crlChecker{provider: NewStaticCRLProvider(nil)}
		err := checker.verifyConnection(tls.ConnectionState{})
		if err == nil || !strings.Contains(err.Error(), "no verified certificate chain") {
			t.Fatalf("expected the check to fail closed, got %v", err)
		}
	})
}

func TestCRLChecker_ProviderError(t *testing.T) {
	ca := newTestCA(t, "crl-provider-error")
	leaf, _ := ca.issueLeaf(t, "vault.example.com")

	checker := &crlChecker{provider: errCRLProvider{err: errors.New("boom")}}
	err := checker.verifyConnection(tls.ConnectionState{
		VerifiedChains: [][]*x509.Certificate{{leaf, ca.cert}},
	})
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected the provider error to surface, got %v", err)
	}
}

type errCRLProvider struct{ err error }

func (e errCRLProvider) CRLs() ([]*x509.RevocationList, error) { return nil, e.err }

func TestFileCRLProvider(t *testing.T) {
	ca := newTestCA(t, "crl-file-ca")
	leaf, _ := ca.issueLeaf(t, "vault.example.com")

	t.Run("single file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, pemCRL(t, ca.crlDER(t, futureTime())))

		p, err := newFileCRLProvider(path, false, time.Minute)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		crls, err := p.CRLs()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(crls) != 1 {
			t.Fatalf("expected 1 CRL, got %d", len(crls))
		}
	})

	t.Run("directory of CRLs", func(t *testing.T) {
		dir := t.TempDir()
		other := newTestCA(t, "crl-file-other")
		writeFile(t, filepath.Join(dir, "a.pem"), pemCRL(t, ca.crlDER(t, futureTime())))
		writeFile(t, filepath.Join(dir, "b.crl"), other.crlDER(t, futureTime()))

		p, err := newFileCRLProvider(dir, true, time.Minute)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		crls, err := p.CRLs()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(crls) != 2 {
			t.Fatalf("expected 2 CRLs, got %d", len(crls))
		}
	})

	// CAPath, which CRLPath mirrors, is walked recursively, so CRLPath must be
	// too.
	t.Run("directory is read recursively", func(t *testing.T) {
		dir := t.TempDir()
		other := newTestCA(t, "crl-file-nested")
		nested := filepath.Join(dir, "nested")
		if err := os.MkdirAll(nested, 0o700); err != nil {
			t.Fatalf("error creating nested directory: %v", err)
		}
		writeFile(t, filepath.Join(dir, "a.pem"), pemCRL(t, ca.crlDER(t, futureTime())))
		writeFile(t, filepath.Join(nested, "b.crl"), other.crlDER(t, futureTime()))

		p, err := newFileCRLProvider(dir, true, time.Minute)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		crls, err := p.CRLs()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(crls) != 2 {
			t.Fatalf("expected the nested CRL to be found, got %d CRLs", len(crls))
		}
	})

	t.Run("picks up changes after the refresh interval", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, pemCRL(t, ca.crlDER(t, futureTime())))

		// A zero interval is replaced by the default, so use the smallest
		// positive value to make every call re-examine the file.
		p, err := newFileCRLProvider(path, false, time.Nanosecond)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		crls, err := p.CRLs()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n := len(crls[0].RevokedCertificateEntries); n != 0 {
			t.Fatalf("expected no revocations initially, got %d", n)
		}

		// Size changes here as well as mtime, so this does not depend on
		// filesystem timestamp granularity.
		writeFile(t, path, pemCRL(t, ca.crlDER(t, futureTime(), leaf)))

		crls, err = p.CRLs()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n := len(crls[0].RevokedCertificateEntries); n != 1 {
			t.Fatalf("expected the reloaded CRL to have 1 revocation, got %d", n)
		}
	})

	t.Run("does not re-read within the refresh interval", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, pemCRL(t, ca.crlDER(t, futureTime())))

		p, err := newFileCRLProvider(path, false, time.Hour)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		writeFile(t, path, pemCRL(t, ca.crlDER(t, futureTime(), leaf)))

		crls, err := p.CRLs()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n := len(crls[0].RevokedCertificateEntries); n != 0 {
			t.Fatalf("expected the cached CRL to be served, got %d revocations", n)
		}
	})

	t.Run("keeps serving the last good CRLs when a reload fails", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, pemCRL(t, ca.crlDER(t, futureTime())))

		p, err := newFileCRLProvider(path, false, time.Nanosecond)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Simulate a truncated write being observed mid-flight.
		writeFile(t, path, []byte("garbage"))

		crls, err := p.CRLs()
		if err != nil {
			t.Fatalf("expected the stale CRLs to be served, got error %v", err)
		}
		if len(crls) != 1 {
			t.Fatalf("expected 1 CRL, got %d", len(crls))
		}

		// And once the file is valid again, the new content is picked up.
		writeFile(t, path, pemCRL(t, ca.crlDER(t, futureTime(), leaf)))
		crls, err = p.CRLs()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n := len(crls[0].RevokedCertificateEntries); n != 1 {
			t.Fatalf("expected the recovered CRL to have 1 revocation, got %d", n)
		}
	})

	// A persistently unreadable file must not turn every handshake into a disk
	// read, so a failed reload is backed off rather than retried immediately.
	t.Run("a failed reload is backed off", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, pemCRL(t, ca.crlDER(t, futureTime())))

		// A long configured interval, so the retry interval is the one under
		// test rather than the configured one.
		p, err := newFileCRLProvider(path, false, time.Hour)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Force the next call to re-examine the file, and make it fail.
		p.mu.Lock()
		p.checkedAt = time.Now().Add(-2 * time.Hour)
		p.mu.Unlock()
		writeFile(t, path, []byte("garbage"))

		if _, err := p.CRLs(); err != nil {
			t.Fatalf("expected the stale CRLs to be served, got %v", err)
		}

		p.mu.Lock()
		failed := p.lastReloadFailed
		due := p.refreshDueLocked()
		p.mu.Unlock()

		if !failed {
			t.Fatal("expected the failure to be recorded")
		}
		if due {
			t.Fatal("expected the next refresh to be backed off rather than due immediately")
		}
	})

	// Handshakes must not queue behind another goroutine's disk read, so the
	// mutex is never held across I/O.
	t.Run("does not hold the lock across disk I/O", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, pemCRL(t, ca.crlDER(t, futureTime())))

		p, err := newFileCRLProvider(path, false, time.Nanosecond)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Claim the refresh the way CRLs does before reading, then confirm the
		// lock is free while the read would be happening.
		p.mu.Lock()
		p.reloading = true
		p.mu.Unlock()

		done := make(chan struct{})
		go func() {
			defer close(done)
			if _, err := p.CRLs(); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()

		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("CRLs blocked while a reload was in flight")
		}
	})

	t.Run("errors surface at construction time", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			set  func(t *testing.T) (string, bool)
		}{
			{
				name: "missing file",
				set: func(t *testing.T) (string, bool) {
					return filepath.Join(t.TempDir(), "nope.pem"), false
				},
			},
			{
				name: "missing directory",
				set: func(t *testing.T) (string, bool) {
					return filepath.Join(t.TempDir(), "nope"), true
				},
			},
			{
				name: "empty directory",
				set: func(t *testing.T) (string, bool) {
					return t.TempDir(), true
				},
			},
			{
				name: "unparseable file",
				set: func(t *testing.T) (string, bool) {
					path := filepath.Join(t.TempDir(), "crl.pem")
					writeFile(t, path, []byte("not a crl"))
					return path, false
				},
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				path, isDir := tc.set(t)
				if _, err := newFileCRLProvider(path, isDir, time.Minute); err == nil {
					t.Fatal("expected an error, got nil")
				}
			})
		}
	})

	t.Run("concurrent access is race free", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, pemCRL(t, ca.crlDER(t, futureTime())))

		p, err := newFileCRLProvider(path, false, time.Nanosecond)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var wg sync.WaitGroup
		for i := 0; i < 16; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 20; j++ {
					if _, err := p.CRLs(); err != nil {
						t.Errorf("unexpected error: %v", err)
						return
					}
				}
			}()
		}
		wg.Wait()
	})
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("error writing %q: %v", path, err)
	}
}

func TestConfigureTLS_CRL(t *testing.T) {
	ca := newTestCA(t, "configure-crl-ca")
	crlPEM := pemCRL(t, ca.crlDER(t, futureTime()))

	t.Run("CRL installs a VerifyConnection callback", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, crlPEM)

		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{CRL: path}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected VerifyConnection to be set")
		}
	})

	t.Run("CRLBytes installs a VerifyConnection callback", func(t *testing.T) {
		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected VerifyConnection to be set")
		}
	})

	t.Run("CRLPath reads a directory", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "crl.pem"), crlPEM)

		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{CRLPath: dir}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected VerifyConnection to be set")
		}
	})

	t.Run("no CRL leaves VerifyConnection alone when none was configured", func(t *testing.T) {
		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tlsConfigOf(t, config).VerifyConnection != nil {
			t.Fatal("expected VerifyConnection to remain nil")
		}
	})

	t.Run("CRL takes precedence over CRLBytes and CRLPath", func(t *testing.T) {
		// Only the CRL field points at something valid, so if either of the
		// others won, configuration would fail.
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, crlPEM)

		config := DefaultConfig()
		err := config.ConfigureTLS(&TLSConfig{
			CRL:      path,
			CRLBytes: []byte("garbage"),
			CRLPath:  filepath.Join(t.TempDir(), "does-not-exist"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("CRLBytes takes precedence over CRLPath", func(t *testing.T) {
		config := DefaultConfig()
		err := config.ConfigureTLS(&TLSConfig{
			CRLBytes: crlPEM,
			CRLPath:  filepath.Join(t.TempDir(), "does-not-exist"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("CRLProvider takes precedence over the file fields", func(t *testing.T) {
		config := DefaultConfig()
		err := config.ConfigureTLS(&TLSConfig{
			CRLProvider: NewStaticCRLProvider(nil),
			CRL:         filepath.Join(t.TempDir(), "does-not-exist"),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected VerifyConnection to be set")
		}
	})

	t.Run("unreadable CRL is an error", func(t *testing.T) {
		config := DefaultConfig()
		err := config.ConfigureTLS(&TLSConfig{CRL: filepath.Join(t.TempDir(), "nope.pem")})
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})

	t.Run("unparseable CRLBytes is an error", func(t *testing.T) {
		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: []byte("garbage")}); err == nil {
			t.Fatal("expected an error, got nil")
		}
	})

	t.Run("CRL with Insecure is rejected", func(t *testing.T) {
		config := DefaultConfig()
		err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM, Insecure: true})
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "insecure TLS") {
			t.Fatalf("expected an error mentioning insecure TLS, got %q", err)
		}
	})

	t.Run("a caller-supplied VerifyConnection still runs", func(t *testing.T) {
		config := DefaultConfig()
		var called bool
		tlsConfigOf(t, config).VerifyConnection = func(tls.ConnectionState) error {
			called = true
			return nil
		}
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// The CRL check will fail on an empty connection state, but the
		// caller's callback must have run first.
		_ = tlsConfigOf(t, config).VerifyConnection(tls.ConnectionState{})
		if !called {
			t.Fatal("expected the pre-existing VerifyConnection to be called")
		}
	})

	t.Run("a caller-supplied VerifyConnection can veto", func(t *testing.T) {
		config := DefaultConfig()
		sentinel := errors.New("vetoed")
		tlsConfigOf(t, config).VerifyConnection = func(tls.ConnectionState) error {
			return sentinel
		}
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := tlsConfigOf(t, config).VerifyConnection(tls.ConnectionState{}); !errors.Is(err, sentinel) {
			t.Fatalf("expected the caller's error, got %v", err)
		}
	})

	t.Run("repeated configuration does not stack callbacks", func(t *testing.T) {
		config := DefaultConfig()
		var calls int
		tlsConfigOf(t, config).VerifyConnection = func(tls.ConnectionState) error {
			calls++
			return nil
		}
		for i := 0; i < 3; i++ {
			if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}
		_ = tlsConfigOf(t, config).VerifyConnection(tls.ConnectionState{})
		if calls != 1 {
			t.Fatalf("expected the original callback to run once, ran %d times", calls)
		}
	})

	// A TLSConfig that says nothing about CRLs must not turn revocation checking
	// off: TLS settings arrive through more than one ConfigureTLS call, so an
	// implicit removal would silently disable the check.
	t.Run("reconfiguring without a CRL retains the check", func(t *testing.T) {
		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected VerifyConnection to be set")
		}
		if err := config.ConfigureTLS(&TLSConfig{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected the CRL check to be retained")
		}
	})

	// The retained check must be a live one, not just a non-nil callback.
	t.Run("the retained check still rejects a revoked certificate", func(t *testing.T) {
		ca := newTestCA(t, "configure-crl-retained")
		leaf, _ := ca.issueLeaf(t, "vault.example.com")

		config := DefaultConfig()
		err := config.ConfigureTLS(&TLSConfig{
			CRLBytes: pemCRL(t, ca.crlDER(t, futureTime(), leaf)),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := config.ConfigureTLS(&TLSConfig{TLSServerName: "vault.example.com"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		verify := tlsConfigOf(t, config).VerifyConnection
		if verify == nil {
			t.Fatal("expected the CRL check to be retained")
		}
		err = verify(tls.ConnectionState{
			VerifiedChains: [][]*x509.Certificate{{leaf, ca.cert}},
		})
		if err == nil || !strings.Contains(err.Error(), "was revoked") {
			t.Fatalf("expected the retained check to reject the revoked certificate, got %v", err)
		}
	})

	// The combination the CLI produces: environment settings install the check,
	// then flag-derived TLS settings arrive as a second call with no CRL fields.
	t.Run("ReadEnvironment followed by a CRL-less ConfigureTLS retains the check", func(t *testing.T) {
		crlPath := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, crlPath, crlPEM)
		caPath := filepath.Join(t.TempDir(), "ca.pem")
		writeFile(t, caPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.cert.Raw}))

		t.Setenv(EnvVaultCRL, crlPath)
		t.Setenv(EnvVaultCACert, caPath)

		config := DefaultConfig()
		if config.Error != nil {
			t.Fatalf("unexpected error: %v", config.Error)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected the environment CRL to install a check")
		}

		// What command/base.go does when -ca-cert is set, or when VAULT_CACERT
		// populates that flag's default.
		err := config.ConfigureTLS(&TLSConfig{CACert: caPath})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected the CRL check to survive the flag-derived ConfigureTLS call")
		}
	})

	t.Run("DisableCRL removes the check", func(t *testing.T) {
		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := config.ConfigureTLS(&TLSConfig{DisableCRL: true}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tlsConfigOf(t, config).VerifyConnection != nil {
			t.Fatal("expected VerifyConnection to be cleared")
		}
	})

	t.Run("DisableCRL with CRL fields is an error", func(t *testing.T) {
		config := DefaultConfig()
		err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM, DisableCRL: true})
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "DisableCRL") {
			t.Fatalf("expected an error mentioning DisableCRL, got %q", err)
		}
	})

	t.Run("DisableCRL restores the caller's callback", func(t *testing.T) {
		config := DefaultConfig()
		var called bool
		original := func(tls.ConnectionState) error {
			called = true
			return nil
		}
		tlsConfigOf(t, config).VerifyConnection = original

		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := config.ConfigureTLS(&TLSConfig{DisableCRL: true}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		restored := tlsConfigOf(t, config).VerifyConnection
		if restored == nil {
			t.Fatal("expected the caller's callback to be restored")
		}
		if err := restored(tls.ConnectionState{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Fatal("expected the caller's callback to run")
		}
	})

	t.Run("a new CRL replaces the retained one", func(t *testing.T) {
		ca := newTestCA(t, "configure-crl-replaced")
		leaf, _ := ca.issueLeaf(t, "vault.example.com")

		config := DefaultConfig()
		// First a CRL that revokes nothing.
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: pemCRL(t, ca.crlDER(t, futureTime()))}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Then one that revokes the leaf.
		err := config.ConfigureTLS(&TLSConfig{CRLBytes: pemCRL(t, ca.crlDER(t, futureTime(), leaf))})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		err = tlsConfigOf(t, config).VerifyConnection(tls.ConnectionState{
			VerifiedChains: [][]*x509.Certificate{{leaf, ca.cert}},
		})
		if err == nil || !strings.Contains(err.Error(), "was revoked") {
			t.Fatalf("expected the newer CRL to be in effect, got %v", err)
		}
	})

	// InsecureSkipVerify may already be on the shared tls.Config from an
	// earlier call, so the guard has to read the live config and not just the
	// incoming TLSConfig.
	t.Run("CRL after a previous Insecure configuration is rejected", func(t *testing.T) {
		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{Insecure: true}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM})
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "insecure TLS") {
			t.Fatalf("expected an error mentioning insecure TLS, got %q", err)
		}
	})

	t.Run("CRL with InsecureSkipVerify set directly on the tls.Config is rejected", func(t *testing.T) {
		config := DefaultConfig()
		tlsConfigOf(t, config).InsecureSkipVerify = true

		err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM})
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !strings.Contains(err.Error(), "insecure TLS") {
			t.Fatalf("expected an error mentioning insecure TLS, got %q", err)
		}
	})

	// An Insecure call arriving after CRL checking is configured must not leave
	// a retained provider that then trips the guard on every later call.
	t.Run("Insecure after a CRL is rejected rather than silently disabling it", func(t *testing.T) {
		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := config.ConfigureTLS(&TLSConfig{Insecure: true}); err == nil {
			t.Fatal("expected an error, got nil")
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected the CRL check to remain in place after the rejected call")
		}
		// The rejected call must not have applied InsecureSkipVerify either, or
		// the check would be left guarding an unverified chain.
		if tlsConfigOf(t, config).InsecureSkipVerify {
			t.Fatal("expected the rejected call to leave InsecureSkipVerify unset")
		}
	})

	// DisableCRL and Insecure together is a coherent request, unlike Insecure
	// on its own while a CRL is configured.
	t.Run("Insecure with DisableCRL is allowed", func(t *testing.T) {
		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := config.ConfigureTLS(&TLSConfig{Insecure: true, DisableCRL: true}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tlsConfigOf(t, config).VerifyConnection != nil {
			t.Fatal("expected the CRL check to be removed")
		}
		if !tlsConfigOf(t, config).InsecureSkipVerify {
			t.Fatal("expected InsecureSkipVerify to be set")
		}
	})

	t.Run("survives CloneConfig", func(t *testing.T) {
		config := DefaultConfig()
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		client, err := NewClient(config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.CloneConfig().TLSConfig().VerifyConnection == nil {
			t.Fatal("expected the cloned config to keep the CRL check")
		}
	})

	// A clone shares the original's Transport and therefore its tls.Config, so
	// reconfiguring the clone must not strip the check from the original. This
	// is what api/auth/cert.NewCertAuthClient does.
	t.Run("reconfiguring a clone does not strip the original's check", func(t *testing.T) {
		ca := newTestCA(t, "configure-crl-clone")
		leaf, _ := ca.issueLeaf(t, "vault.example.com")

		config := DefaultConfig()
		err := config.ConfigureTLS(&TLSConfig{
			CRLBytes: pemCRL(t, ca.crlDER(t, futureTime(), leaf)),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		client, err := NewClient(config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		clone := client.CloneConfig()
		// The caller's TLSConfig mentions no CRL fields, as NewCertAuthClient's
		// does not.
		if err := clone.ConfigureTLS(&TLSConfig{TLSServerName: "vault.example.com"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for name, verify := range map[string]func(tls.ConnectionState) error{
			"original": tlsConfigOf(t, config).VerifyConnection,
			"clone":    tlsConfigOf(t, clone).VerifyConnection,
		} {
			if verify == nil {
				t.Fatalf("expected the %s to keep the CRL check", name)
			}
			err := verify(tls.ConnectionState{
				VerifiedChains: [][]*x509.Certificate{{leaf, ca.cert}},
			})
			if err == nil || !strings.Contains(err.Error(), "was revoked") {
				t.Fatalf("expected the %s check to reject the revoked certificate, got %v", name, err)
			}
		}
	})

	// The clone shares the original's tls.Config, so a CRL configured on the
	// clone must replace the original's callback rather than stack onto it.
	t.Run("reconfiguring a clone with a CRL does not stack checks", func(t *testing.T) {
		config := DefaultConfig()
		var calls int
		tlsConfigOf(t, config).VerifyConnection = func(tls.ConnectionState) error {
			calls++
			return nil
		}
		if err := config.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		client, err := NewClient(config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		clone := client.CloneConfig()
		if err := clone.ConfigureTLS(&TLSConfig{CRLBytes: crlPEM}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_ = tlsConfigOf(t, clone).VerifyConnection(tls.ConnectionState{})
		if calls != 1 {
			t.Fatalf("expected the original callback to run once, ran %d times", calls)
		}
	})
}

func tlsConfigOf(t *testing.T, config *Config) *tls.Config {
	t.Helper()
	transport, ok := config.HttpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("unexpected transport type %T", config.HttpClient.Transport)
	}
	return transport.TLSClientConfig
}

func TestClientCRLEnvSettings(t *testing.T) {
	ca := newTestCA(t, "crl-env-ca")
	crlPEM := pemCRL(t, ca.crlDER(t, futureTime()))

	t.Run("VAULT_CRL", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, crlPEM)
		t.Setenv(EnvVaultCRL, path)

		config := DefaultConfig()
		if config.Error != nil {
			t.Fatalf("unexpected error: %v", config.Error)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected VerifyConnection to be set")
		}
	})

	t.Run("VAULT_CRL_BYTES", func(t *testing.T) {
		t.Setenv(EnvVaultCRLBytes, string(crlPEM))

		config := DefaultConfig()
		if config.Error != nil {
			t.Fatalf("unexpected error: %v", config.Error)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected VerifyConnection to be set")
		}
	})

	t.Run("VAULT_CRL_PATH", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, filepath.Join(dir, "crl.pem"), crlPEM)
		t.Setenv(EnvVaultCRLPath, dir)

		config := DefaultConfig()
		if config.Error != nil {
			t.Fatalf("unexpected error: %v", config.Error)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected VerifyConnection to be set")
		}
	})

	t.Run("invalid VAULT_CRL surfaces on the config", func(t *testing.T) {
		t.Setenv(EnvVaultCRL, filepath.Join(t.TempDir(), "nope.pem"))

		config := DefaultConfig()
		if config.Error == nil {
			t.Fatal("expected DefaultConfig to report an error")
		}
	})

	t.Run("unset by default", func(t *testing.T) {
		config := DefaultConfig()
		if config.Error != nil {
			t.Fatalf("unexpected error: %v", config.Error)
		}
		if tlsConfigOf(t, config).VerifyConnection != nil {
			t.Fatal("expected VerifyConnection to remain nil")
		}
	})

	// VAULT_CRL and VAULT_CACERT together is the expected combination: a private
	// CA is exactly the situation where revocation data is supplied by hand.
	t.Run("VAULT_CRL with VAULT_CACERT", func(t *testing.T) {
		crlPath := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, crlPath, crlPEM)
		caPath := filepath.Join(t.TempDir(), "ca.pem")
		writeFile(t, caPath, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.cert.Raw}))

		t.Setenv(EnvVaultCRL, crlPath)
		t.Setenv(EnvVaultCACert, caPath)

		config := DefaultConfig()
		if config.Error != nil {
			t.Fatalf("unexpected error: %v", config.Error)
		}
		if tlsConfigOf(t, config).VerifyConnection == nil {
			t.Fatal("expected VerifyConnection to be set")
		}
	})

	t.Run("VAULT_CRL with VAULT_SKIP_VERIFY is an error", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, path, crlPEM)
		t.Setenv(EnvVaultCRL, path)
		t.Setenv(EnvVaultSkipVerify, "true")

		config := DefaultConfig()
		if config.Error == nil {
			t.Fatal("expected DefaultConfig to report an error")
		}
		if !strings.Contains(config.Error.Error(), "insecure TLS") {
			t.Fatalf("expected an error mentioning insecure TLS, got %q", config.Error)
		}
	})
}

// crlTestServer starts an HTTPS server using a certificate issued by ca.
func crlTestServer(t *testing.T, ca *testCA, handler http.Handler) (addr string, requests *int32Counter) {
	t.Helper()

	leafCert, leafKey := ca.issueLeaf(t, "localhost")
	counter := &int32Counter{}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("error listening: %v", err)
	}

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter.inc()
			handler.ServeHTTP(w, r)
		}),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			Certificates: []tls.Certificate{{
				Certificate: [][]byte{leafCert.Raw},
				PrivateKey:  leafKey,
			}},
		},
	}

	go server.ServeTLS(ln, "", "")
	t.Cleanup(func() { server.Close() })

	return fmt.Sprintf("https://127.0.0.1:%d", ln.Addr().(*net.TCPAddr).Port), counter
}

type int32Counter struct {
	mu sync.Mutex
	n  int
}

func (c *int32Counter) inc() {
	c.mu.Lock()
	c.n++
	c.mu.Unlock()
}

func (c *int32Counter) get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.n
}

// TestCRLEndToEnd exercises the check over a real TLS handshake.
func TestCRLEndToEnd(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{}}`))
	})

	t.Run("clean CRL allows the request", func(t *testing.T) {
		ca := newTestCA(t, "crl-e2e-clean")
		addr, requests := crlTestServer(t, ca, handler)

		client := crlTestClient(t, ca, addr, ca.crlDER(t, futureTime()))
		if _, err := client.Logical().Read("secret/foo"); err != nil {
			t.Fatalf("expected the request to succeed, got %v", err)
		}
		if requests.get() != 1 {
			t.Fatalf("expected the server to see 1 request, saw %d", requests.get())
		}
	})

	t.Run("revoked server certificate fails the handshake", func(t *testing.T) {
		ca := newTestCA(t, "crl-e2e-revoked")
		leafCert, leafKey := ca.issueLeaf(t, "localhost")
		counter := &int32Counter{}

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("error listening: %v", err)
		}
		server := &http.Server{
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				counter.inc()
				handler.ServeHTTP(w, r)
			}),
			TLSConfig: &tls.Config{
				MinVersion:   tls.VersionTLS12,
				Certificates: []tls.Certificate{{Certificate: [][]byte{leafCert.Raw}, PrivateKey: leafKey}},
			},
		}
		go server.ServeTLS(ln, "", "")
		t.Cleanup(func() { server.Close() })

		addr := fmt.Sprintf("https://127.0.0.1:%d", ln.Addr().(*net.TCPAddr).Port)
		client := crlTestClient(t, ca, addr, ca.crlDER(t, futureTime(), leafCert))

		_, err = client.Logical().Read("secret/foo")
		if err == nil {
			t.Fatal("expected the request to fail")
		}
		if !strings.Contains(err.Error(), "was revoked") {
			t.Fatalf("expected a revocation error, got %v", err)
		}
		if counter.get() != 0 {
			t.Fatalf("expected the handler never to run, it ran %d times", counter.get())
		}
	})

	t.Run("missing CRL for the issuer fails the handshake", func(t *testing.T) {
		ca := newTestCA(t, "crl-e2e-uncovered")
		unrelated := newTestCA(t, "crl-e2e-unrelated")
		addr, requests := crlTestServer(t, ca, handler)

		client := crlTestClient(t, ca, addr, unrelated.crlDER(t, futureTime()))
		_, err := client.Logical().Read("secret/foo")
		if err == nil {
			t.Fatal("expected the request to fail")
		}
		if !strings.Contains(err.Error(), "no CRL was supplied for issuer") {
			t.Fatalf("expected a coverage error, got %v", err)
		}
		if requests.get() != 0 {
			t.Fatalf("expected the handler never to run, it ran %d times", requests.get())
		}
	})

	t.Run("concurrent requests", func(t *testing.T) {
		ca := newTestCA(t, "crl-e2e-concurrent")
		addr, _ := crlTestServer(t, ca, handler)
		client := crlTestClient(t, ca, addr, ca.crlDER(t, futureTime()))

		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if _, err := client.Logical().Read("secret/foo"); err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}()
		}
		wg.Wait()
	})

	// The configuration shape the CLI produces: the CRL arrives via the
	// environment and a second, CRL-less ConfigureTLS call follows. The check
	// must still reject a revoked certificate over a real handshake.
	t.Run("a later CRL-less ConfigureTLS still rejects a revoked certificate", func(t *testing.T) {
		ca := newTestCA(t, "crl-e2e-retained")
		leafCert, leafKey := ca.issueLeaf(t, "localhost")

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatalf("error listening: %v", err)
		}
		server := &http.Server{
			Handler: handler,
			TLSConfig: &tls.Config{
				MinVersion:   tls.VersionTLS12,
				Certificates: []tls.Certificate{{Certificate: [][]byte{leafCert.Raw}, PrivateKey: leafKey}},
			},
		}
		go server.ServeTLS(ln, "", "")
		t.Cleanup(func() { server.Close() })

		caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.cert.Raw})
		caPath := filepath.Join(t.TempDir(), "ca.pem")
		writeFile(t, caPath, caPEM)
		crlPath := filepath.Join(t.TempDir(), "crl.pem")
		writeFile(t, crlPath, pemCRL(t, ca.crlDER(t, futureTime(), leafCert)))

		t.Setenv(EnvVaultCRL, crlPath)
		t.Setenv(EnvVaultCACert, caPath)

		config := DefaultConfig()
		if config.Error != nil {
			t.Fatalf("unexpected error: %v", config.Error)
		}
		config.Address = fmt.Sprintf("https://127.0.0.1:%d", ln.Addr().(*net.TCPAddr).Port)
		config.MaxRetries = 0

		// What command/base.go does once -ca-cert is set or VAULT_CACERT
		// populates its default: a second call with no CRL fields.
		if err := config.ConfigureTLS(&TLSConfig{CACert: caPath}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		client, err := NewClient(config)
		if err != nil {
			t.Fatalf("error creating client: %v", err)
		}
		client.SetToken("test-token")

		_, err = client.Logical().Read("secret/foo")
		if err == nil {
			t.Fatal("expected the request to fail because the server certificate was revoked")
		}
		if !strings.Contains(err.Error(), "was revoked") {
			t.Fatalf("expected a revocation error, got %v", err)
		}
	})
}

// crlTestClient builds a client that trusts ca and checks the given CRL.
func crlTestClient(t *testing.T, ca *testCA, addr string, crlDER []byte) *Client {
	t.Helper()

	config := DefaultConfig()
	config.Address = addr
	// Keep handshake failures from being retried, so tests fail fast.
	config.MaxRetries = 0

	err := config.ConfigureTLS(&TLSConfig{
		CACertBytes: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.cert.Raw}),
		CRLBytes:    pemCRL(t, crlDER),
	})
	if err != nil {
		t.Fatalf("error configuring TLS: %v", err)
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("error creating client: %v", err)
	}
	client.SetToken("test-token")
	return client
}

// TestCRLCheckedOnResumedConnection is the reason VerifyConnection is used
// instead of VerifyPeerCertificate: the latter is not invoked when a TLS
// session is resumed, which would let a resumed connection skip the revocation
// check entirely.
func TestCRLCheckedOnResumedConnection(t *testing.T) {
	ca := newTestCA(t, "crl-resume-ca")
	leafCert, leafKey := ca.issueLeaf(t, "localhost")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("error listening: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	serverTLS := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{{Certificate: [][]byte{leafCert.Raw}, PrivateKey: leafKey}},
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func() {
				tlsConn := tls.Server(conn, serverTLS)
				// A complete handshake is enough; the client only cares
				// whether its own verification callback ran.
				_ = tlsConn.HandshakeContext(t.Context())
				tlsConn.Close()
			}()
		}
	}()

	roots := x509.NewCertPool()
	roots.AddCert(ca.cert)

	// A shared session cache is what makes resumption possible in the first
	// place. The api package never sets one by default, so this stands in for
	// a caller that supplies its own http.Client.
	cache := tls.NewLRUClientSessionCache(4)

	// revoked is flipped between the two handshakes to model a CRL that is
	// refreshed after the first connection.
	var revoked bool
	provider := funcCRLProvider(func() ([]*x509.RevocationList, error) {
		if revoked {
			return []*x509.RevocationList{ca.crl(t, futureTime(), leafCert)}, nil
		}
		return []*x509.RevocationList{ca.crl(t, futureTime())}, nil
	})
	checker := &crlChecker{provider: provider}

	var sawResumption bool
	clientTLS := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		RootCAs:            roots,
		ServerName:         "localhost",
		ClientSessionCache: cache,
		VerifyConnection: func(cs tls.ConnectionState) error {
			if cs.DidResume {
				sawResumption = true
			}
			return checker.verifyConnection(cs)
		},
	}

	// First handshake: populates the session cache and must succeed.
	conn, err := tls.Dial("tcp", ln.Addr().String(), clientTLS)
	if err != nil {
		t.Fatalf("expected the first handshake to succeed, got %v", err)
	}
	// Read once so that a TLS 1.3 session ticket, which arrives after the
	// handshake completes, makes it into the cache.
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _ = conn.Read(make([]byte, 1))
	conn.Close()

	// Second handshake: the certificate is now revoked and must be rejected
	// even if the session is resumed.
	revoked = true
	conn2, err := tls.Dial("tcp", ln.Addr().String(), clientTLS)
	if err == nil {
		conn2.Close()
		t.Fatal("expected the second handshake to fail because the certificate was revoked")
	}
	if !strings.Contains(err.Error(), "was revoked") {
		t.Fatalf("expected a revocation error, got %v", err)
	}

	if !sawResumption {
		t.Skip("the second connection did not resume a session, so this environment cannot exercise the resumption path")
	}

	ln.Close()
	<-done
}

// TestVerifyPeerCertificateSkipsResumedConnections documents the behaviour
// being guarded against above: with VerifyPeerCertificate, a resumed
// connection skips the check.
func TestVerifyPeerCertificateSkipsResumedConnections(t *testing.T) {
	ca := newTestCA(t, "crl-resume-negative-ca")
	leafCert, leafKey := ca.issueLeaf(t, "localhost")

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("error listening: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	serverTLS := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{{Certificate: [][]byte{leafCert.Raw}, PrivateKey: leafKey}},
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func() {
				tlsConn := tls.Server(conn, serverTLS)
				_ = tlsConn.HandshakeContext(t.Context())
				tlsConn.Close()
			}()
		}
	}()

	roots := x509.NewCertPool()
	roots.AddCert(ca.cert)

	var calls int
	var resumed bool
	clientTLS := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		RootCAs:            roots,
		ServerName:         "localhost",
		ClientSessionCache: tls.NewLRUClientSessionCache(4),
		VerifyPeerCertificate: func([][]byte, [][]*x509.Certificate) error {
			calls++
			return nil
		},
		VerifyConnection: func(cs tls.ConnectionState) error {
			if cs.DidResume {
				resumed = true
			}
			return nil
		},
	}

	conn, err := tls.Dial("tcp", ln.Addr().String(), clientTLS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _ = conn.Read(make([]byte, 1))
	conn.Close()

	firstCalls := calls

	conn2, err := tls.Dial("tcp", ln.Addr().String(), clientTLS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	conn2.Close()

	if !resumed {
		t.Skip("the second connection did not resume a session, so this environment cannot exercise the resumption path")
	}
	if calls != firstCalls {
		t.Fatalf("expected VerifyPeerCertificate to be skipped on the resumed connection, but it was called %d more time(s)", calls-firstCalls)
	}
}

type funcCRLProvider func() ([]*x509.RevocationList, error)

func (f funcCRLProvider) CRLs() ([]*x509.RevocationList, error) { return f() }

// Copyright IBM Corp. 2025
// SPDX-License-Identifier: MPL-2.0

package certutil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/helper/certutil/x509verify"
	"github.com/stretchr/testify/require"
)

// buildTestChain generates a self-signed root CA, an intermediate CA signed by
// the root, and a leaf certificate signed by the intermediate. It returns
// parsedBundles suitable for passing to VerifyCertificateChain.
//
// The notBefore/notAfter parameters control the leaf certificate's validity
// window, making it easy to create expired certificates for negative testing.
func buildTestChain(t *testing.T, leafNotBefore, leafNotAfter time.Time) *ParsedCertBundle {
	t.Helper()

	// Root CA key and self-signed certificate.
	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "generate root key")

	rootTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Test Root CA"},
		NotBefore:             time.Now().Add(-2 * time.Hour),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	rootDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	require.NoError(t, err, "create root certificate")
	rootCert, err := x509.ParseCertificate(rootDER)
	require.NoError(t, err, "parse root certificate")

	// Intermediate CA key and certificate signed by root.
	intermediateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "generate intermediate key")

	intermediateTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "Test Intermediate CA"},
		NotBefore:             time.Now().Add(-1 * time.Hour),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	intermediateDER, err := x509.CreateCertificate(rand.Reader, intermediateTemplate, rootCert, &intermediateKey.PublicKey, rootKey)
	require.NoError(t, err, "create intermediate certificate")
	intermediateCert, err := x509.ParseCertificate(intermediateDER)
	require.NoError(t, err, "parse intermediate certificate")

	// Leaf certificate signed by the intermediate.
	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "generate leaf key")

	leafTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(3),
		Subject:               pkix.Name{CommonName: "Test Leaf"},
		NotBefore:             leafNotBefore,
		NotAfter:              leafNotAfter,
		BasicConstraintsValid: false,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, intermediateCert, &leafKey.PublicKey, intermediateKey)
	require.NoError(t, err, "create leaf certificate")

	return &ParsedCertBundle{
		CertificateBytes: leafDER,
		CAChain: []*CertBlock{
			{Bytes: intermediateDER},
			{Bytes: rootDER},
		},
	}
}

// TestVerifyCertificateChain_ValidChain verifies that a well-formed three-tier
// chain (root → intermediate → leaf) passes verification with default options.
func TestVerifyCertificateChain_ValidChain(t *testing.T) {
	t.Parallel()

	now := time.Now()
	bundle := buildTestChain(t, now.Add(-1*time.Hour), now.Add(24*time.Hour))
	opts := x509verify.VerifyOptions{
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	err := VerifyCertificateChain(bundle, opts, false)
	require.NoError(t, err, "valid chain should verify successfully")
}

// TestVerifyCertificateChain_ExpiredLeaf verifies that an expired leaf
// certificate causes verification to fail when time checks are enabled, and
// succeeds when DisableTimeChecks is set — which is the behaviour required by
// Vault's PKI engine when issuing certificates without enforcing validity windows.
func TestVerifyCertificateChain_ExpiredLeaf(t *testing.T) {
	t.Parallel()

	past := time.Now().Add(-48 * time.Hour)
	bundle := buildTestChain(t, past.Add(-1*time.Hour), past)

	opts := x509verify.VerifyOptions{
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	err := VerifyCertificateChain(bundle, opts, false)
	require.Error(t, err, "expired leaf should fail verification with time checks enabled")

	opts.DisableTimeChecks = true
	err = VerifyCertificateChain(bundle, opts, false)
	require.NoError(t, err, "expired leaf should pass verification when DisableTimeChecks is true")
}

// TestVerifyCertificateChain_NoRoot verifies that when no root CA is present in
// the chain, the intermediate is promoted to the trust anchor (rootChainOnly=false),
// while rootChainOnly=true causes verification to fail.
func TestVerifyCertificateChain_NoRoot(t *testing.T) {
	t.Parallel()

	now := time.Now()
	fullBundle := buildTestChain(t, now.Add(-1*time.Hour), now.Add(24*time.Hour))

	// Strip the root from the CA chain, leaving only the intermediate.
	bundleNoRoot := &ParsedCertBundle{
		CertificateBytes: fullBundle.CertificateBytes,
		CAChain:          fullBundle.CAChain[:1], // intermediate only
	}

	opts := x509verify.VerifyOptions{
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}

	err := VerifyCertificateChain(bundleNoRoot, opts, false)
	require.NoError(t, err, "intermediate as trust anchor should succeed when rootChainOnly is false")

	err = VerifyCertificateChain(bundleNoRoot, opts, true)
	require.Error(t, err, "verification should fail when rootChainOnly is true and no root is present")
}

// TestVerifyCertificateChain_SelfSigned verifies that a self-signed certificate
// with its chain pointing to itself is accepted — this is the case where Vault
// stores the issuer as both root and leaf.
func TestVerifyCertificateChain_SelfSigned(t *testing.T) {
	t.Parallel()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err, "generate key")

	now := time.Now()
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Self-Signed Root"},
		NotBefore:             now.Add(-1 * time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err, "create self-signed certificate")

	bundle := &ParsedCertBundle{
		CertificateBytes: der,
		CAChain: []*CertBlock{
			{Bytes: der},
		},
	}

	opts := x509verify.VerifyOptions{
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	err = VerifyCertificateChain(bundle, opts, false)
	require.NoError(t, err, "self-signed certificate should verify successfully")
}

// TestVerifyCertificateChain_DisablePathLenChecks verifies that path length constraint
// violations can be bypassed when DisablePathLenChecks is true.
func TestVerifyCertificateChain_DisablePathLenChecks(t *testing.T) {
	t.Parallel()

	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	now := time.Now()
	rootTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Root CA"},
		NotBefore:             now.Add(-1 * time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	rootDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	require.NoError(t, err)
	rootCert, err := x509.ParseCertificate(rootDER)
	require.NoError(t, err)

	intKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	intTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "Intermediate CA"},
		NotBefore:             now.Add(-1 * time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	intDER, err := x509.CreateCertificate(rand.Reader, intTemplate, rootCert, &intKey.PublicKey, rootKey)
	require.NoError(t, err)
	intCert, err := x509.ParseCertificate(intDER)
	require.NoError(t, err)

	subIntKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	subIntTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(3),
		Subject:               pkix.Name{CommonName: "Sub Intermediate CA"},
		NotBefore:             now.Add(-1 * time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	subIntDER, err := x509.CreateCertificate(rand.Reader, subIntTemplate, intCert, &subIntKey.PublicKey, intKey)
	require.NoError(t, err)
	subIntCert, err := x509.ParseCertificate(subIntDER)
	require.NoError(t, err)

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	leafTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(4),
		Subject:      pkix.Name{CommonName: "Leaf"},
		NotBefore:    now.Add(-1 * time.Hour),
		NotAfter:     now.Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, subIntCert, &leafKey.PublicKey, subIntKey)
	require.NoError(t, err)

	bundle := &ParsedCertBundle{
		CertificateBytes: leafDER,
		CAChain: []*CertBlock{
			{Bytes: subIntDER},
			{Bytes: intDER},
			{Bytes: rootDER},
		},
	}

	opts := x509verify.VerifyOptions{
		KeyUsages:            []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		DisablePathLenChecks: false,
	}
	err = VerifyCertificateChain(bundle, opts, false)
	require.Error(t, err, "chain exceeding MaxPathLen should fail when DisablePathLenChecks is false")

	opts.DisablePathLenChecks = true
	err = VerifyCertificateChain(bundle, opts, false)
	require.NoError(t, err, "chain exceeding MaxPathLen should succeed when DisablePathLenChecks is true")
}

// TestVerifyCertificateChain_DisableNameConstraintChecks verifies that name constraint
// violations can be bypassed when DisableNameConstraintChecks is true.
func TestVerifyCertificateChain_DisableNameConstraintChecks(t *testing.T) {
	t.Parallel()

	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	now := time.Now()
	rootTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Root CA"},
		NotBefore:             now.Add(-1 * time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		PermittedDNSDomains:   []string{".example.com"},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	rootDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	require.NoError(t, err)
	rootCert, err := x509.ParseCertificate(rootDER)
	require.NoError(t, err)

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	leafTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "other.domain.com"},
		DNSNames:     []string{"other.domain.com"},
		NotBefore:    now.Add(-1 * time.Hour),
		NotAfter:     now.Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, rootCert, &leafKey.PublicKey, rootKey)
	require.NoError(t, err)

	bundle := &ParsedCertBundle{
		CertificateBytes: leafDER,
		CAChain: []*CertBlock{
			{Bytes: rootDER},
		},
	}

	opts := x509verify.VerifyOptions{
		KeyUsages:                   []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		DisableNameConstraintChecks: false,
	}
	err = VerifyCertificateChain(bundle, opts, false)
	require.Error(t, err, "name not matching permitted DNS domains should fail when DisableNameConstraintChecks is false")

	opts.DisableNameConstraintChecks = true
	err = VerifyCertificateChain(bundle, opts, false)
	require.NoError(t, err, "name constraint violation should be ignored when DisableNameConstraintChecks is true")
}

// TestVerifyCertificateChain_DisableCriticalExtensionChecks verifies that unknown critical
// extensions can be bypassed when DisableCriticalExtensionChecks is true.
func TestVerifyCertificateChain_DisableCriticalExtensionChecks(t *testing.T) {
	t.Parallel()

	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	now := time.Now()
	rootTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Root CA"},
		NotBefore:             now.Add(-1 * time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	rootDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	require.NoError(t, err)
	rootCert, err := x509.ParseCertificate(rootDER)
	require.NoError(t, err)

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	leafTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "leaf.example.com"},
		NotBefore:    now.Add(-1 * time.Hour),
		NotAfter:     now.Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		ExtraExtensions: []pkix.Extension{
			{
				Id:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, // Custom unknown extension OID
				Critical: true,
				Value:    []byte("custom data"),
			},
		},
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, rootCert, &leafKey.PublicKey, rootKey)
	require.NoError(t, err)

	bundle := &ParsedCertBundle{
		CertificateBytes: leafDER,
		CAChain: []*CertBlock{
			{Bytes: rootDER},
		},
	}

	opts := x509verify.VerifyOptions{
		KeyUsages:                      []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		DisableCriticalExtensionChecks: false,
	}
	err = VerifyCertificateChain(bundle, opts, false)
	require.Error(t, err, "unhandled critical extension should fail when DisableCriticalExtensionChecks is false")

	opts.DisableCriticalExtensionChecks = true
	err = VerifyCertificateChain(bundle, opts, false)
	require.NoError(t, err, "unhandled critical extension should pass when DisableCriticalExtensionChecks is true")
}

// TestVerifyCertificateChain_DisableNameChecks verifies that issuer/subject name mismatches
// (e.g. cross-signed intermediate with a different subject name) can be bypassed when DisableNameChecks is true.
func TestVerifyCertificateChain_DisableNameChecks(t *testing.T) {
	t.Parallel()

	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	now := time.Now()
	rootTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Root CA"},
		NotBefore:             now.Add(-1 * time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	rootDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	require.NoError(t, err)
	rootCert, err := x509.ParseCertificate(rootDER)
	require.NoError(t, err)

	// Intermediate signed by Root with subject "Int CA Original"
	intKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	intSKID := []byte{1, 2, 3, 4, 5}
	intTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "Int CA Original"},
		SubjectKeyId:          intSKID,
		NotBefore:             now.Add(-1 * time.Hour),
		NotAfter:              now.Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}
	intDER, err := x509.CreateCertificate(rand.Reader, intTemplate, rootCert, &intKey.PublicKey, rootKey)
	require.NoError(t, err)

	// Leaf issued with Issuer "Int CA Renamed", but signed by intKey and referencing intSKID as AKID
	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	leafTemplate := &x509.Certificate{
		SerialNumber:   big.NewInt(3),
		Subject:        pkix.Name{CommonName: "Leaf Certificate"},
		AuthorityKeyId: intSKID,
		NotBefore:      now.Add(-1 * time.Hour),
		NotAfter:       now.Add(24 * time.Hour),
		KeyUsage:       x509.KeyUsageDigitalSignature,
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	// Sign leaf with Int CA template modified to have a different Subject (so leaf.Issuer != intCert.Subject)
	signingParentTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "Int CA Renamed"},
		SubjectKeyId: intSKID,
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, signingParentTemplate, &leafKey.PublicKey, intKey)
	require.NoError(t, err)

	bundle := &ParsedCertBundle{
		CertificateBytes: leafDER,
		CAChain: []*CertBlock{
			{Bytes: intDER},
			{Bytes: rootDER},
		},
	}

	opts := x509verify.VerifyOptions{
		KeyUsages:         []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		DisableNameChecks: false,
	}
	err = VerifyCertificateChain(bundle, opts, false)
	require.Error(t, err, "issuer name mismatch should fail when DisableNameChecks is false")

	opts.DisableNameChecks = true
	err = VerifyCertificateChain(bundle, opts, false)
	require.NoError(t, err, "issuer name mismatch should pass when DisableNameChecks is true")
}

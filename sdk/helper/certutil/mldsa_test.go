// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package certutil

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/cloudflare/circl/sign/mldsa/mldsa65"
	"github.com/cloudflare/circl/sign/mldsa/mldsa87"
)

func TestMLDSA65KeyGeneration(t *testing.T) {
	var bundle ParsedCertBundle
	err := generatePrivateKey("ml-dsa-65", 0, &bundle, rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate ML-DSA-65 key: %v", err)
	}

	if bundle.PrivateKeyType != MLDSA65PrivateKey {
		t.Fatalf("expected key type %s, got %s", MLDSA65PrivateKey, bundle.PrivateKeyType)
	}

	if bundle.PrivateKey == nil {
		t.Fatal("private key is nil")
	}

	if len(bundle.PrivateKeyBytes) == 0 {
		t.Fatal("private key bytes are empty")
	}

	// Verify the key implements crypto.Signer
	signer, ok := bundle.PrivateKey.(crypto.Signer)
	if !ok {
		t.Fatal("private key does not implement crypto.Signer")
	}

	// Verify we can get the public key
	pub := signer.Public()
	if pub == nil {
		t.Fatal("public key is nil")
	}

	// Verify public key type detection
	keyType := GetPrivateKeyTypeFromSigner(signer)
	if keyType != MLDSA65PrivateKey {
		t.Fatalf("expected signer key type %s, got %s", MLDSA65PrivateKey, keyType)
	}

	pubKeyType := GetPrivateKeyTypeFromPublicKey(pub)
	if pubKeyType != MLDSA65PrivateKey {
		t.Fatalf("expected public key type %s, got %s", MLDSA65PrivateKey, pubKeyType)
	}
}

func TestMLDSA87KeyGeneration(t *testing.T) {
	var bundle ParsedCertBundle
	err := generatePrivateKey("ml-dsa-87", 0, &bundle, rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate ML-DSA-87 key: %v", err)
	}

	if bundle.PrivateKeyType != MLDSA87PrivateKey {
		t.Fatalf("expected key type %s, got %s", MLDSA87PrivateKey, bundle.PrivateKeyType)
	}

	if bundle.PrivateKey == nil {
		t.Fatal("private key is nil")
	}

	if len(bundle.PrivateKeyBytes) == 0 {
		t.Fatal("private key bytes are empty")
	}

	signer, ok := bundle.PrivateKey.(crypto.Signer)
	if !ok {
		t.Fatal("private key does not implement crypto.Signer")
	}

	pub := signer.Public()
	if pub == nil {
		t.Fatal("public key is nil")
	}

	keyType := GetPrivateKeyTypeFromSigner(signer)
	if keyType != MLDSA87PrivateKey {
		t.Fatalf("expected signer key type %s, got %s", MLDSA87PrivateKey, keyType)
	}
}

func TestMLDSA65KeyRoundTrip(t *testing.T) {
	// Generate a key
	var bundle ParsedCertBundle
	err := generatePrivateKey("ml-dsa-65", 0, &bundle, rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	// Round-trip through PEM encoding
	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bundle.PrivateKeyBytes,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)

	// Parse back
	decoded, _ := pem.Decode(pemBytes)
	if decoded == nil {
		t.Fatal("failed to decode PEM")
	}

	signer, blockType, err := ParseDERKey(decoded.Bytes)
	if err != nil {
		t.Fatalf("failed to parse DER key: %v", err)
	}

	if blockType != PKCS8Block {
		t.Fatalf("expected block type %s, got %s", PKCS8Block, blockType)
	}

	keyType := GetPrivateKeyTypeFromSigner(signer)
	if keyType != MLDSA65PrivateKey {
		t.Fatalf("expected key type %s after round-trip, got %s", MLDSA65PrivateKey, keyType)
	}

	// Verify the round-tripped key can sign
	msg := []byte("test message")
	sig, err := signer.Sign(nil, msg, crypto.Hash(0))
	if err != nil {
		t.Fatalf("failed to sign with round-tripped key: %v", err)
	}

	// Verify signature
	pk := signer.Public().(*mldsa65.PublicKey)
	if !mldsa65.Verify(pk, msg, nil, sig) {
		t.Fatal("signature verification failed for round-tripped key")
	}
}

func TestMLDSA87KeyRoundTrip(t *testing.T) {
	var bundle ParsedCertBundle
	err := generatePrivateKey("ml-dsa-87", 0, &bundle, rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bundle.PrivateKeyBytes,
	}
	pemBytes := pem.EncodeToMemory(pemBlock)

	decoded, _ := pem.Decode(pemBytes)
	if decoded == nil {
		t.Fatal("failed to decode PEM")
	}

	signer, _, err := ParseDERKey(decoded.Bytes)
	if err != nil {
		t.Fatalf("failed to parse DER key: %v", err)
	}

	keyType := GetPrivateKeyTypeFromSigner(signer)
	if keyType != MLDSA87PrivateKey {
		t.Fatalf("expected key type %s after round-trip, got %s", MLDSA87PrivateKey, keyType)
	}

	msg := []byte("test message for ML-DSA-87")
	sig, err := signer.Sign(nil, msg, crypto.Hash(0))
	if err != nil {
		t.Fatalf("failed to sign with round-tripped key: %v", err)
	}

	pk := signer.Public().(*mldsa87.PublicKey)
	if !mldsa87.Verify(pk, msg, nil, sig) {
		t.Fatal("signature verification failed for round-tripped key")
	}
}

func TestMLDSA65PublicKeyComparison(t *testing.T) {
	_, sk1, _ := mldsa65.GenerateKey(rand.Reader)
	_, sk2, _ := mldsa65.GenerateKey(rand.Reader)

	pk1 := sk1.Public()
	pk1Copy := sk1.Public()
	pk2 := sk2.Public()

	// Same key should be equal
	equal, err := ComparePublicKeys(pk1, pk1Copy)
	if err != nil {
		t.Fatalf("error comparing same keys: %v", err)
	}
	if !equal {
		t.Fatal("same ML-DSA-65 public keys should be equal")
	}

	// Different keys should not be equal
	equal, err = ComparePublicKeys(pk1, pk2)
	if err != nil {
		t.Fatalf("error comparing different keys: %v", err)
	}
	if equal {
		t.Fatal("different ML-DSA-65 public keys should not be equal")
	}
}

func TestMLDSA87PublicKeyComparison(t *testing.T) {
	_, sk1, _ := mldsa87.GenerateKey(rand.Reader)
	_, sk2, _ := mldsa87.GenerateKey(rand.Reader)

	pk1 := sk1.Public()
	pk1Copy := sk1.Public()
	pk2 := sk2.Public()

	equal, err := ComparePublicKeys(pk1, pk1Copy)
	if err != nil {
		t.Fatalf("error comparing same keys: %v", err)
	}
	if !equal {
		t.Fatal("same ML-DSA-87 public keys should be equal")
	}

	equal, err = ComparePublicKeys(pk1, pk2)
	if err != nil {
		t.Fatalf("error comparing different keys: %v", err)
	}
	if equal {
		t.Fatal("different ML-DSA-87 public keys should not be equal")
	}
}

func TestMLDSA65SubjectKeyID(t *testing.T) {
	_, sk, _ := mldsa65.GenerateKey(rand.Reader)

	skid, err := GetSubjKeyID(sk)
	if err != nil {
		t.Fatalf("failed to get subject key ID: %v", err)
	}

	if len(skid) == 0 {
		t.Fatal("subject key ID is empty")
	}

	// SHA-1 hash should be 20 bytes
	if len(skid) != 20 {
		t.Fatalf("expected 20 byte SKID, got %d bytes", len(skid))
	}
}

func TestMLDSA65GetPublicKeySize(t *testing.T) {
	_, sk, _ := mldsa65.GenerateKey(rand.Reader)
	pub := sk.Public()

	size := GetPublicKeySize(pub)
	expected := mldsa65.PublicKeySize * 8
	if size != expected {
		t.Fatalf("expected public key size %d bits, got %d", expected, size)
	}
}

func TestMLDSA87GetPublicKeySize(t *testing.T) {
	_, sk, _ := mldsa87.GenerateKey(rand.Reader)
	pub := sk.Public()

	size := GetPublicKeySize(pub)
	expected := mldsa87.PublicKeySize * 8
	if size != expected {
		t.Fatalf("expected public key size %d bits, got %d", expected, size)
	}
}

func TestMLDSA65ValidateKeyType(t *testing.T) {
	// ML-DSA-65 should accept keyBits=0 (no bit size concept)
	keyBits, err := ValidateDefaultOrValueKeyType("ml-dsa-65", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keyBits != 0 {
		t.Fatalf("expected keyBits 0, got %d", keyBits)
	}
}

func TestMLDSA87ValidateKeyType(t *testing.T) {
	keyBits, err := ValidateDefaultOrValueKeyType("ml-dsa-87", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keyBits != 0 {
		t.Fatalf("expected keyBits 0, got %d", keyBits)
	}
}

func TestMLDSA65ValidateSignatureLength(t *testing.T) {
	// ML-DSA uses built-in hashing, so any hash bits value should be accepted
	err := ValidateSignatureLength("ml-dsa-65", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMLDSA65HashBits(t *testing.T) {
	// ML-DSA should return 0 hash bits (uses built-in hashing)
	hashBits, err := DefaultOrValueHashBits("ml-dsa-65", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hashBits != 0 {
		t.Fatalf("expected 0 hash bits for ML-DSA-65, got %d", hashBits)
	}
}

func TestMLDSA65SelfSignedCertificate(t *testing.T) {
	// Generate ML-DSA-65 key pair
	_, sk, err := mldsa65.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "ML-DSA-65 Test Root CA",
			Organization: []string{"Test"},
		},
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Create self-signed certificate
	certDER, err := createMLDSACertificate(template, template, sk.Public(), sk)
	if err != nil {
		t.Fatalf("failed to create ML-DSA certificate: %v", err)
	}

	if len(certDER) == 0 {
		t.Fatal("certificate DER is empty")
	}

	// Parse the certificate back
	cert, err := parseMLDSACertificate(certDER)
	if err != nil {
		t.Fatalf("failed to parse ML-DSA certificate: %v", err)
	}

	if cert.Subject.CommonName != "ML-DSA-65 Test Root CA" {
		t.Fatalf("expected CN 'ML-DSA-65 Test Root CA', got '%s'", cert.Subject.CommonName)
	}

	if cert.SerialNumber.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("unexpected serial number: %v", cert.SerialNumber)
	}
}

func TestMLDSA87SelfSignedCertificate(t *testing.T) {
	_, sk, err := mldsa87.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName:   "ML-DSA-87 Test Root CA",
			Organization: []string{"Test"},
		},
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDER, err := createMLDSACertificate(template, template, sk.Public(), sk)
	if err != nil {
		t.Fatalf("failed to create ML-DSA-87 certificate: %v", err)
	}

	cert, err := parseMLDSACertificate(certDER)
	if err != nil {
		t.Fatalf("failed to parse ML-DSA-87 certificate: %v", err)
	}

	if cert.Subject.CommonName != "ML-DSA-87 Test Root CA" {
		t.Fatalf("expected CN 'ML-DSA-87 Test Root CA', got '%s'", cert.Subject.CommonName)
	}
}

func TestMLDSA65CertificateWithSANs(t *testing.T) {
	_, sk, err := mldsa65.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	uri, _ := url.Parse("https://example.com")
	template := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			CommonName: "ml-dsa-65-test.example.com",
		},
		NotBefore:      time.Now().Add(-30 * time.Second),
		NotAfter:       time.Now().Add(24 * time.Hour),
		KeyUsage:       x509.KeyUsageDigitalSignature,
		DNSNames:       []string{"ml-dsa-65-test.example.com", "*.example.com"},
		IPAddresses:    []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
		EmailAddresses: []string{"test@example.com"},
		URIs:           []*url.URL{uri},
	}

	certDER, err := createMLDSACertificate(template, template, sk.Public(), sk)
	if err != nil {
		t.Fatalf("failed to create ML-DSA certificate with SANs: %v", err)
	}

	if len(certDER) == 0 {
		t.Fatal("certificate DER is empty")
	}

	cert, err := parseMLDSACertificate(certDER)
	if err != nil {
		t.Fatalf("failed to parse ML-DSA certificate: %v", err)
	}

	if cert.Subject.CommonName != "ml-dsa-65-test.example.com" {
		t.Fatalf("unexpected CN: %s", cert.Subject.CommonName)
	}
}

func TestMLDSA65CSRCreation(t *testing.T) {
	_, sk, err := mldsa65.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   "ML-DSA-65 CSR Test",
			Organization: []string{"Test Org"},
		},
		DNSNames: []string{"csr-test.example.com"},
	}

	csrDER, err := createMLDSACSR(template, sk)
	if err != nil {
		t.Fatalf("failed to create ML-DSA CSR: %v", err)
	}

	if len(csrDER) == 0 {
		t.Fatal("CSR DER is empty")
	}
}

func TestMLDSA65CreateCertificateViaBundle(t *testing.T) {
	// Test the full createCertificate path with ML-DSA keys
	bundle := &CreationBundle{
		Params: &CreationParameters{
			Subject: pkix.Name{
				CommonName:   "ML-DSA-65 Bundle Test",
				Organization: []string{"Test Org"},
			},
			KeyType:  "ml-dsa-65",
			KeyBits:  0,
			NotAfter: time.Now().Add(365 * 24 * time.Hour),
			IsCA:     true,
			URLs: &URLEntries{
				IssuingCertificates:   []string{},
				CRLDistributionPoints: []string{},
				OCSPServers:           []string{},
			},
			KeyUsage:          x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
			NotBeforeDuration: 30 * time.Second,
		},
	}

	result, err := createCertificate(bundle, rand.Reader, generatePrivateKey)
	if err != nil {
		t.Fatalf("failed to create certificate via bundle: %v", err)
	}

	if result.PrivateKeyType != MLDSA65PrivateKey {
		t.Fatalf("expected key type %s, got %s", MLDSA65PrivateKey, result.PrivateKeyType)
	}

	if result.PrivateKey == nil {
		t.Fatal("private key is nil")
	}

	if len(result.CertificateBytes) == 0 {
		t.Fatal("certificate bytes are empty")
	}

	if result.Certificate == nil {
		t.Fatal("parsed certificate is nil")
	}

	if result.Certificate.Subject.CommonName != "ML-DSA-65 Bundle Test" {
		t.Fatalf("unexpected CN: %s", result.Certificate.Subject.CommonName)
	}
}

func TestMLDSA65CreateKeyBundle(t *testing.T) {
	kb, err := CreateKeyBundle("ml-dsa-65", 0, rand.Reader)
	if err != nil {
		t.Fatalf("failed to create key bundle: %v", err)
	}

	if kb.PrivateKeyType != MLDSA65PrivateKey {
		t.Fatalf("expected key type %s, got %s", MLDSA65PrivateKey, kb.PrivateKeyType)
	}

	if kb.PrivateKey == nil {
		t.Fatal("private key is nil")
	}

	if len(kb.PrivateKeyBytes) == 0 {
		t.Fatal("private key bytes are empty")
	}

	// Test PEM string generation
	pemStr, err := kb.ToPrivateKeyPemString()
	if err != nil {
		t.Fatalf("failed to get PEM string: %v", err)
	}

	if !strings.Contains(pemStr, "PRIVATE KEY") {
		t.Fatal("PEM string does not contain expected header")
	}
}

func TestMLDSA87CreateKeyBundle(t *testing.T) {
	kb, err := CreateKeyBundle("ml-dsa-87", 0, rand.Reader)
	if err != nil {
		t.Fatalf("failed to create key bundle: %v", err)
	}

	if kb.PrivateKeyType != MLDSA87PrivateKey {
		t.Fatalf("expected key type %s, got %s", MLDSA87PrivateKey, kb.PrivateKeyType)
	}

	if kb.PrivateKey == nil {
		t.Fatal("private key is nil")
	}

	pemStr, err := kb.ToPrivateKeyPemString()
	if err != nil {
		t.Fatalf("failed to get PEM string: %v", err)
	}

	if !strings.Contains(pemStr, "PRIVATE KEY") {
		t.Fatal("PEM string does not contain expected header")
	}
}

func TestMLDSA65CSRViaBundleCreation(t *testing.T) {
	bundle := &CreationBundle{
		Params: &CreationParameters{
			Subject: pkix.Name{
				CommonName:   "ML-DSA-65 CSR Bundle Test",
				Organization: []string{"Test Org"},
			},
			KeyType: "ml-dsa-65",
			KeyBits: 0,
			DNSNames: []string{"csr-bundle.example.com"},
			URLs: &URLEntries{
				IssuingCertificates:   []string{},
				CRLDistributionPoints: []string{},
				OCSPServers:           []string{},
			},
		},
	}

	result, err := createCSR(bundle, false, rand.Reader, generatePrivateKey)
	if err != nil {
		t.Fatalf("failed to create CSR via bundle: %v", err)
	}

	if result.PrivateKeyType != MLDSA65PrivateKey {
		t.Fatalf("expected key type %s, got %s", MLDSA65PrivateKey, result.PrivateKeyType)
	}

	if len(result.CSRBytes) == 0 {
		t.Fatal("CSR bytes are empty")
	}
}

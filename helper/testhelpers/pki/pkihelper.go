// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	mathrand2 "math/rand/v2"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// This file contains helper functions for generating CA hierarchies for testing

type LeafWithRoot struct {
	RootCa             GeneratedCert
	Leaf               GeneratedCert
	CombinedLeafCaFile string
}

type LeafWithIntermediary struct {
	RootCa         GeneratedCert
	IntCa          GeneratedCert
	Leaf           GeneratedCert
	CombinedCaFile string
}

type GeneratedCert struct {
	KeyFile  string
	CertFile string
	CertPem  *pem.Block
	Cert     *x509.Certificate
	Key      *ecdsa.PrivateKey
}

// GenerateCertWithIntermediaryRoot generates a leaf certificate signed by an intermediary root CA
func GenerateCertWithIntermediaryRoot(t testing.TB) LeafWithIntermediary {
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

	ca := GenerateRootCa(t)
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
	err := os.WriteFile(combinedCasFile, append(pem.EncodeToMemory(caInt.CertPem), pem.EncodeToMemory(ca.CertPem)...), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	return LeafWithIntermediary{
		RootCa:         ca,
		IntCa:          caInt,
		Leaf:           leafCert,
		CombinedCaFile: combinedCasFile,
	}
}

// generateCertAndSign generates a certificate and associated key signed by a CA
func generateCertAndSign(t testing.TB, template *x509.Certificate, ca GeneratedCert, tempDir string, filePrefix string) GeneratedCert {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, template, ca.Cert, key.Public(), ca.Key)
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
	return GeneratedCert{
		KeyFile:  keyFile,
		CertFile: certFile,
		CertPem:  certPEMBlock,
		Cert:     cert,
		Key:      key,
	}
}

// GenerateCertWithRoot generates a leaf certificate signed by a root CA
func GenerateCertWithRoot(t testing.TB) LeafWithRoot {
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

	ca := GenerateRootCa(t)
	leafCert := generateCertAndSign(t, leafTemplate, ca, tempDir, "leaf_")

	combinedCaLeafFile := filepath.Join(tempDir, "leaf-ca.pem")
	err := os.WriteFile(combinedCaLeafFile, append(pem.EncodeToMemory(leafCert.CertPem), pem.EncodeToMemory(ca.CertPem)...), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	return LeafWithRoot{
		RootCa:             ca,
		Leaf:               leafCert,
		CombinedLeafCaFile: combinedCaLeafFile,
	}
}

// GenerateRootCa generates a self-signed root CA certificate and key
func GenerateRootCa(t testing.TB) GeneratedCert {
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
	return GeneratedCert{
		CertPem:  caCertPEMBlock,
		CertFile: caFile,
		KeyFile:  caKeyFile,
		Cert:     caCert,
		Key:      caKey,
	}
}

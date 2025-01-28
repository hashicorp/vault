// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package certutil

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"net"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/go-test/deep"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
)

// Tests converting back and forth between a CertBundle and a ParsedCertBundle.
//
// Also tests the GetSubjKeyID, GetHexFormatted, ParseHexFormatted and
// ParsedCertBundle.getSigner functions.
func TestCertBundleConversion(t *testing.T) {
	cbuts := []*CertBundle{
		refreshRSACertBundle(),
		refreshRSACertBundleWithChain(),
		refreshRSA8CertBundle(),
		refreshRSA8CertBundleWithChain(),
		refreshECCertBundle(),
		refreshECCertBundleWithChain(),
		refreshEC8CertBundle(),
		refreshEC8CertBundleWithChain(),
		refreshEd255198CertBundle(),
		refreshEd255198CertBundleWithChain(),
	}

	for i, cbut := range cbuts {
		pcbut, err := cbut.ToParsedCertBundle()
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Errorf("Error converting to parsed cert bundle: %s", err)
			continue
		}

		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Error(err.Error())
		}

		cbut, err := pcbut.ToCertBundle()
		if err != nil {
			t.Fatalf("Error converting to cert bundle: %s", err)
		}

		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func BenchmarkCertBundleParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cbuts := []*CertBundle{
			refreshRSACertBundle(),
			refreshRSACertBundleWithChain(),
			refreshRSA8CertBundle(),
			refreshRSA8CertBundleWithChain(),
			refreshECCertBundle(),
			refreshECCertBundleWithChain(),
			refreshEC8CertBundle(),
			refreshEC8CertBundleWithChain(),
			refreshEd255198CertBundle(),
			refreshEd255198CertBundleWithChain(),
		}

		for i, cbut := range cbuts {
			pcbut, err := cbut.ToParsedCertBundle()
			if err != nil {
				b.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
				b.Errorf("Error converting to parsed cert bundle: %s", err)
				continue
			}

			cbut, err = pcbut.ToCertBundle()
			if err != nil {
				b.Fatalf("Error converting to cert bundle: %s", err)
			}
		}
	}
}

func TestCertBundleParsing(t *testing.T) {
	cbuts := []*CertBundle{
		refreshRSACertBundle(),
		refreshRSACertBundleWithChain(),
		refreshRSA8CertBundle(),
		refreshRSA8CertBundleWithChain(),
		refreshECCertBundle(),
		refreshECCertBundleWithChain(),
		refreshEC8CertBundle(),
		refreshEC8CertBundleWithChain(),
		refreshEd255198CertBundle(),
		refreshEd255198CertBundleWithChain(),
	}

	for i, cbut := range cbuts {
		jsonString, err := json.Marshal(cbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf("Error marshaling testing certbundle to JSON: %s", err)
		}
		pcbut, err := ParsePKIJSON(jsonString)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf("Error during JSON bundle handling: %s", err)
		}
		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatal(err.Error())
		}

		dataMap := structs.New(cbut).Map()
		pcbut, err = ParsePKIMap(dataMap)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf("Error during JSON bundle handling: %s", err)
		}
		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatal(err.Error())
		}

		pcbut, err = ParsePEMBundle(cbut.ToPEMBundle())
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf("Error during JSON bundle handling: %s", err)
		}
		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatal(err.Error())
		}
	}
}

func compareCertBundleToParsedCertBundle(cbut *CertBundle, pcbut *ParsedCertBundle) error {
	if cbut == nil {
		return errors.New("got nil bundle")
	}
	if pcbut == nil {
		return fmt.Errorf("got nil parsed bundle")
	}

	switch {
	case pcbut.Certificate == nil:
		return fmt.Errorf("parsed bundle has nil certificate")
	case pcbut.PrivateKey == nil:
		return fmt.Errorf("parsed bundle has nil private key")
	}

	switch cbut.PrivateKey {
	case privRSAKeyPem:
		if pcbut.PrivateKeyType != RSAPrivateKey {
			return fmt.Errorf("parsed bundle has wrong private key type: %v, should be 'rsa' (%v)", pcbut.PrivateKeyType, RSAPrivateKey)
		}
	case privRSA8KeyPem:
		if pcbut.PrivateKeyType != RSAPrivateKey {
			return fmt.Errorf("parsed bundle has wrong pkcs8 private key type: %v, should be 'rsa' (%v)", pcbut.PrivateKeyType, RSAPrivateKey)
		}
	case privECKeyPem:
		if pcbut.PrivateKeyType != ECPrivateKey {
			return fmt.Errorf("parsed bundle has wrong private key type: %v, should be 'ec' (%v)", pcbut.PrivateKeyType, ECPrivateKey)
		}
	case privEC8KeyPem:
		if pcbut.PrivateKeyType != ECPrivateKey {
			return fmt.Errorf("parsed bundle has wrong pkcs8 private key type: %v, should be 'ec' (%v)", pcbut.PrivateKeyType, ECPrivateKey)
		}
	case privEd255198KeyPem:
		if pcbut.PrivateKeyType != Ed25519PrivateKey {
			return fmt.Errorf("parsed bundle has wrong pkcs8 private key type: %v, should be 'ed25519' (%v)", pcbut.PrivateKeyType, ECPrivateKey)
		}
	default:
		return fmt.Errorf("parsed bundle has unknown private key type")
	}

	subjKeyID, err := GetSubjKeyID(pcbut.PrivateKey)
	if err != nil {
		return fmt.Errorf("error when getting subject key id: %s", err)
	}
	if bytes.Compare(subjKeyID, pcbut.Certificate.SubjectKeyId) != 0 {
		return fmt.Errorf("parsed bundle private key does not match subject key id\nGot\n%#v\nExpected\n%#v\nCert\n%#v", subjKeyID, pcbut.Certificate.SubjectKeyId, *pcbut.Certificate)
	}

	switch {
	case len(pcbut.CAChain) > 0 && len(cbut.CAChain) == 0:
		return fmt.Errorf("parsed bundle ca chain has certs when cert bundle does not")
	case len(pcbut.CAChain) == 0 && len(cbut.CAChain) > 0:
		return fmt.Errorf("cert bundle ca chain has certs when parsed cert bundle does not")
	}

	cb, err := pcbut.ToCertBundle()
	if err != nil {
		return fmt.Errorf("thrown error during parsed bundle conversion: %s\n\nInput was: %#v", err, *pcbut)
	}

	switch {
	case len(cb.Certificate) == 0:
		return fmt.Errorf("bundle has nil certificate")
	case len(cb.PrivateKey) == 0:
		return fmt.Errorf("bundle has nil private key")
	case len(cb.CAChain[0]) == 0:
		return fmt.Errorf("bundle has nil issuing CA")
	}

	switch pcbut.PrivateKeyType {
	case RSAPrivateKey:
		if cb.PrivateKey != privRSAKeyPem && cb.PrivateKey != privRSA8KeyPem {
			return fmt.Errorf("bundle private key does not match")
		}
	case ECPrivateKey:
		if cb.PrivateKey != privECKeyPem && cb.PrivateKey != privEC8KeyPem {
			return fmt.Errorf("bundle private key does not match")
		}
	case Ed25519PrivateKey:
		if cb.PrivateKey != privEd255198KeyPem {
			return fmt.Errorf("bundle private key does not match")
		}
	default:
		return fmt.Errorf("certBundle has unknown private key type")
	}

	if cb.SerialNumber != GetHexFormatted(pcbut.Certificate.SerialNumber.Bytes(), ":") {
		return fmt.Errorf("bundle serial number does not match")
	}

	if !bytes.Equal(pcbut.Certificate.SerialNumber.Bytes(), ParseHexFormatted(cb.SerialNumber, ":")) {
		return fmt.Errorf("failed re-parsing hex formatted number %s", cb.SerialNumber)
	}

	switch {
	case len(pcbut.CAChain) > 0 && len(cb.CAChain) == 0:
		return fmt.Errorf("parsed bundle ca chain has certs when cert bundle does not")
	case len(pcbut.CAChain) == 0 && len(cb.CAChain) > 0:
		return fmt.Errorf("cert bundle ca chain has certs when parsed cert bundle does not")
	case !reflect.DeepEqual(cbut.CAChain, cb.CAChain):
		return fmt.Errorf("cert bundle ca chain does not match: %#v\n\n%#v", cbut.CAChain, cb.CAChain)
	}

	return nil
}

func TestCSRBundleConversion(t *testing.T) {
	csrbuts := []*CSRBundle{
		refreshRSACSRBundle(),
		refreshECCSRBundle(),
		refreshEd25519CSRBundle(),
	}

	for _, csrbut := range csrbuts {
		pcsrbut, err := csrbut.ToParsedCSRBundle()
		if err != nil {
			t.Fatalf("Error converting to parsed CSR bundle: %v", err)
		}

		err = compareCSRBundleToParsedCSRBundle(csrbut, pcsrbut)
		if err != nil {
			t.Fatal(err.Error())
		}

		csrbut, err = pcsrbut.ToCSRBundle()
		if err != nil {
			t.Fatalf("Error converting to CSR bundle: %v", err)
		}

		err = compareCSRBundleToParsedCSRBundle(csrbut, pcsrbut)
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func compareCSRBundleToParsedCSRBundle(csrbut *CSRBundle, pcsrbut *ParsedCSRBundle) error {
	if csrbut == nil {
		return fmt.Errorf("got nil bundle")
	}
	if pcsrbut == nil {
		return fmt.Errorf("got nil parsed bundle")
	}

	switch {
	case pcsrbut.CSR == nil:
		return fmt.Errorf("parsed bundle has nil csr")
	case pcsrbut.PrivateKey == nil:
		return fmt.Errorf("parsed bundle has nil private key")
	}

	switch csrbut.PrivateKey {
	case privRSAKeyPem:
		if pcsrbut.PrivateKeyType != RSAPrivateKey {
			return fmt.Errorf("parsed bundle has wrong private key type")
		}
	case privECKeyPem:
		if pcsrbut.PrivateKeyType != ECPrivateKey {
			return fmt.Errorf("parsed bundle has wrong private key type")
		}
	case privEd255198KeyPem:
		if pcsrbut.PrivateKeyType != Ed25519PrivateKey {
			return fmt.Errorf("parsed bundle has wrong private key type")
		}
	default:
		return fmt.Errorf("parsed bundle has unknown private key type")
	}

	csrb, err := pcsrbut.ToCSRBundle()
	if err != nil {
		return fmt.Errorf("Thrown error during parsed bundle conversion: %s\n\nInput was: %#v", err, *pcsrbut)
	}

	switch {
	case len(csrb.CSR) == 0:
		return fmt.Errorf("bundle has nil certificate")
	case len(csrb.PrivateKey) == 0:
		return fmt.Errorf("bundle has nil private key")
	}

	switch csrb.PrivateKeyType {
	case "rsa":
		if pcsrbut.PrivateKeyType != RSAPrivateKey {
			return fmt.Errorf("bundle has wrong private key type")
		}
		if csrb.PrivateKey != privRSAKeyPem {
			return fmt.Errorf("bundle rsa private key does not match\nGot\n%#v\nExpected\n%#v", csrb.PrivateKey, privRSAKeyPem)
		}
	case "ec":
		if pcsrbut.PrivateKeyType != ECPrivateKey {
			return fmt.Errorf("bundle has wrong private key type")
		}
		if csrb.PrivateKey != privECKeyPem {
			return fmt.Errorf("bundle ec private key does not match")
		}
	case "ed25519":
		if pcsrbut.PrivateKeyType != Ed25519PrivateKey {
			return fmt.Errorf("bundle has wrong private key type")
		}
		if csrb.PrivateKey != privEd255198KeyPem {
			return fmt.Errorf("bundle ed25519 private key does not match")
		}
	default:
		return fmt.Errorf("bundle has unknown private key type")
	}

	return nil
}

func TestTLSConfig(t *testing.T) {
	cbut := refreshRSACertBundle()

	pcbut, err := cbut.ToParsedCertBundle()
	if err != nil {
		t.Fatalf("Error getting parsed cert bundle: %s", err)
	}

	usages := []TLSUsage{
		TLSUnknown,
		TLSClient,
		TLSServer,
		TLSClient | TLSServer,
	}

	for _, usage := range usages {
		tlsConfig, err := pcbut.GetTLSConfig(usage)
		if err != nil {
			t.Fatalf("Error getting tls config: %s", err)
		}
		if tlsConfig == nil {
			t.Fatalf("Got nil tls.Config")
		}

		if len(tlsConfig.Certificates) != 1 {
			t.Fatalf("Unexpected length in config.Certificates")
		}

		// Length should be 2, since we passed in a CA
		if len(tlsConfig.Certificates[0].Certificate) != 2 {
			t.Fatalf("Did not find both certificates in config.Certificates.Certificate")
		}

		if tlsConfig.Certificates[0].Leaf != pcbut.Certificate {
			t.Fatalf("Leaf certificate does not match parsed bundle's certificate")
		}

		if tlsConfig.Certificates[0].PrivateKey != pcbut.PrivateKey {
			t.Fatalf("Config's private key does not match parsed bundle's private key")
		}

		switch usage {
		case TLSServer | TLSClient:
			if len(tlsConfig.ClientCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.ClientCAs.Subjects()[0], pcbut.CAChain[0].Certificate.RawSubject) != 0 {
				t.Fatalf("CA certificate not in client cert pool as expected")
			}
			if len(tlsConfig.RootCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.RootCAs.Subjects()[0], pcbut.CAChain[0].Certificate.RawSubject) != 0 {
				t.Fatalf("CA certificate not in root cert pool as expected")
			}
		case TLSServer:
			if len(tlsConfig.ClientCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.ClientCAs.Subjects()[0], pcbut.CAChain[0].Certificate.RawSubject) != 0 {
				t.Fatalf("CA certificate not in client cert pool as expected")
			}
			if tlsConfig.RootCAs != nil {
				t.Fatalf("Found root pools in config object when not expected")
			}
		case TLSClient:
			if len(tlsConfig.RootCAs.Subjects()) != 1 || bytes.Compare(tlsConfig.RootCAs.Subjects()[0], pcbut.CAChain[0].Certificate.RawSubject) != 0 {
				t.Fatalf("CA certificate not in root cert pool as expected")
			}
			if tlsConfig.ClientCAs != nil {
				t.Fatalf("Found root pools in config object when not expected")
			}
		default:
			if tlsConfig.RootCAs != nil || tlsConfig.ClientCAs != nil {
				t.Fatalf("Found root pools in config object when not expected")
			}
		}
	}
}

func TestNewCertPool(t *testing.T) {
	caExample := `-----BEGIN CERTIFICATE-----
MIIC5zCCAc+gAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwptaW5p
a3ViZUNBMB4XDTE5MTIxMDIzMDUxOVoXDTI5MTIwODIzMDUxOVowFTETMBEGA1UE
AxMKbWluaWt1YmVDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANFi
/RIdMHd865X6JygTb9riX01DA3QnR+RoXDXNnj8D3LziLG2n8ItXMJvWbU3sxxyy
nX9HxJ0SIeexj1cYzdQBtJDjO1/PeuKc4CZ7zCukCAtHz8mC7BDPOU7F7pggpcQ0
/t/pa2m22hmCu8aDF9WlUYHtJpYATnI/A5vz/VFLR9daxmkl59Qo3oHITj7vAzSx
/75r9cibpQyJ+FhiHOZHQWYY2JYw2g4v5hm5hg5SFM9yFcZ75ISI9ebyFFIl9iBY
zAk9jqv1mXvLr0Q39AVwMTamvGuap1oocjM9NIhQvaFL/DNqF1ouDQjCf5u2imLc
TraO1/2KO8fqwOZCOrMCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgKkMB0GA1UdJQQW
MBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQBtVZCwCPqUUUpIClAlE9nc2fo2bTs9gsjXRmqdQ5oaSomSLE93
aJWYFuAhxPXtlApbLYZfW2m1sM3mTVQN60y0uE4e1jdSN1ErYQ9slJdYDAMaEmOh
iSexj+Nd1scUiMHV9lf3ps5J8sYeCpwZX3sPmw7lqZojTS12pANBDcigsaj5RRyN
9GyP3WkSQUsTpWlDb9Fd+KNdkCVw7nClIpBPA2KW4BQKw/rNSvOFD61mbzc89lo0
Q9IFGQFFF8jO18lbyWqnRBGXcS4/G7jQ3S7C121d14YLUeAYOM7pJykI1g4CLx9y
vitin0L6nprauWkKO38XgM4T75qKZpqtiOcT
-----END CERTIFICATE-----
`
	if _, err := NewCertPool(bytes.NewReader([]byte(caExample))); err != nil {
		t.Fatal(err)
	}
}

func TestGetPublicKeySize(t *testing.T) {
	rsa, err := cryptoutil.GenerateRSAKey(rand.Reader, 3072)
	if err != nil {
		t.Fatal(err)
	}
	if GetPublicKeySize(&rsa.PublicKey) != 3072 {
		t.Fatal("unexpected rsa key size")
	}
	ecdsa, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if GetPublicKeySize(&ecdsa.PublicKey) != 384 {
		t.Fatal("unexpected ecdsa key size")
	}
	ed25519, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if GetPublicKeySize(ed25519) != 256 {
		t.Fatal("unexpected ed25519 key size")
	}
	// Skipping DSA as too slow
}

func refreshRSA8CertBundle() *CertBundle {
	initTest.Do(setCerts)
	return &CertBundle{
		Certificate: certRSAPem,
		PrivateKey:  privRSA8KeyPem,
		CAChain:     []string{issuingCaChainPem[0]},
	}
}

func refreshRSA8CertBundleWithChain() *CertBundle {
	initTest.Do(setCerts)
	ret := refreshRSA8CertBundle()
	ret.CAChain = issuingCaChainPem
	return ret
}

func refreshRSACertBundle() *CertBundle {
	initTest.Do(setCerts)
	return &CertBundle{
		Certificate: certRSAPem,
		CAChain:     []string{issuingCaChainPem[0]},
		PrivateKey:  privRSAKeyPem,
	}
}

func refreshRSACertBundleWithChain() *CertBundle {
	initTest.Do(setCerts)
	ret := refreshRSACertBundle()
	ret.CAChain = issuingCaChainPem
	return ret
}

func refreshECCertBundle() *CertBundle {
	initTest.Do(setCerts)
	return &CertBundle{
		Certificate: certECPem,
		CAChain:     []string{issuingCaChainPem[0]},
		PrivateKey:  privECKeyPem,
	}
}

func refreshECCertBundleWithChain() *CertBundle {
	initTest.Do(setCerts)
	ret := refreshECCertBundle()
	ret.CAChain = issuingCaChainPem
	return ret
}

func refreshEd255198CertBundle() *CertBundle {
	initTest.Do(setCerts)
	return &CertBundle{
		Certificate: certEd25519Pem,
		PrivateKey:  privEd255198KeyPem,
		CAChain:     []string{issuingCaChainPem[0]},
	}
}

func refreshEd255198CertBundleWithChain() *CertBundle {
	initTest.Do(setCerts)
	ret := refreshEd255198CertBundle()
	ret.CAChain = issuingCaChainPem
	return ret
}

func refreshEd25519CSRBundle() *CSRBundle {
	initTest.Do(setCerts)
	return &CSRBundle{
		CSR:        csrEd25519Pem,
		PrivateKey: privEd255198KeyPem,
	}
}

func refreshRSACSRBundle() *CSRBundle {
	initTest.Do(setCerts)
	return &CSRBundle{
		CSR:        csrRSAPem,
		PrivateKey: privRSAKeyPem,
	}
}

func refreshECCSRBundle() *CSRBundle {
	initTest.Do(setCerts)
	return &CSRBundle{
		CSR:        csrECPem,
		PrivateKey: privECKeyPem,
	}
}

func refreshEC8CertBundle() *CertBundle {
	initTest.Do(setCerts)
	return &CertBundle{
		Certificate: certECPem,
		PrivateKey:  privEC8KeyPem,
		CAChain:     []string{issuingCaChainPem[0]},
	}
}

func refreshEC8CertBundleWithChain() *CertBundle {
	initTest.Do(setCerts)
	ret := refreshEC8CertBundle()
	ret.CAChain = issuingCaChainPem
	return ret
}

func setCerts() {
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	subjKeyID, err := GetSubjKeyID(caKey)
	if err != nil {
		panic(err)
	}
	caCertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "root.localhost",
		},
		SubjectKeyId:          subjKeyID,
		DNSNames:              []string{"root.localhost"},
		KeyUsage:              x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign),
		SerialNumber:          big.NewInt(mathrand.Int63()),
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, caKey.Public(), caKey)
	if err != nil {
		panic(err)
	}
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		panic(err)
	}
	caCertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	}
	caCertPEM := strings.TrimSpace(string(pem.EncodeToMemory(caCertPEMBlock)))

	intKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	subjKeyID, err = GetSubjKeyID(intKey)
	if err != nil {
		panic(err)
	}
	intCertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "int.localhost",
		},
		SubjectKeyId:          subjKeyID,
		DNSNames:              []string{"int.localhost"},
		KeyUsage:              x509.KeyUsage(x509.KeyUsageCertSign | x509.KeyUsageCRLSign),
		SerialNumber:          big.NewInt(mathrand.Int63()),
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	intBytes, err := x509.CreateCertificate(rand.Reader, intCertTemplate, caCert, intKey.Public(), caKey)
	if err != nil {
		panic(err)
	}
	intCert, err := x509.ParseCertificate(intBytes)
	if err != nil {
		panic(err)
	}
	intCertPEMBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: intBytes,
	}
	intCertPEM := strings.TrimSpace(string(pem.EncodeToMemory(intCertPEMBlock)))

	// EC generation
	{
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(err)
		}
		subjKeyID, err := GetSubjKeyID(key)
		if err != nil {
			panic(err)
		}
		certTemplate := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			SubjectKeyId: subjKeyID,
			DNSNames:     []string{"localhost"},
			ExtKeyUsage: []x509.ExtKeyUsage{
				x509.ExtKeyUsageServerAuth,
				x509.ExtKeyUsageClientAuth,
			},
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
			SerialNumber: big.NewInt(mathrand.Int63()),
			NotBefore:    time.Now().Add(-30 * time.Second),
			NotAfter:     time.Now().Add(262980 * time.Hour),
		}
		csrTemplate := &x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			DNSNames: []string{"localhost"},
		}
		csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, key)
		if err != nil {
			panic(err)
		}
		csrPEMBlock := &pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: csrBytes,
		}
		csrECPem = strings.TrimSpace(string(pem.EncodeToMemory(csrPEMBlock)))
		certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, intCert, key.Public(), intKey)
		if err != nil {
			panic(err)
		}
		certPEMBlock := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		}
		certECPem = strings.TrimSpace(string(pem.EncodeToMemory(certPEMBlock)))
		marshaledKey, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			panic(err)
		}
		keyPEMBlock := &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: marshaledKey,
		}
		privECKeyPem = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
		marshaledKey, err = x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			panic(err)
		}
		keyPEMBlock = &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: marshaledKey,
		}
		privEC8KeyPem = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
	}

	// RSA generation
	{
		key, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		subjKeyID, err := GetSubjKeyID(key)
		if err != nil {
			panic(err)
		}
		certTemplate := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			SubjectKeyId: subjKeyID,
			DNSNames:     []string{"localhost"},
			ExtKeyUsage: []x509.ExtKeyUsage{
				x509.ExtKeyUsageServerAuth,
				x509.ExtKeyUsageClientAuth,
			},
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
			SerialNumber: big.NewInt(mathrand.Int63()),
			NotBefore:    time.Now().Add(-30 * time.Second),
			NotAfter:     time.Now().Add(262980 * time.Hour),
		}
		csrTemplate := &x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			DNSNames: []string{"localhost"},
		}
		csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, key)
		if err != nil {
			panic(err)
		}
		csrPEMBlock := &pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: csrBytes,
		}
		csrRSAPem = strings.TrimSpace(string(pem.EncodeToMemory(csrPEMBlock)))
		certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, intCert, key.Public(), intKey)
		if err != nil {
			panic(err)
		}
		certPEMBlock := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		}
		certRSAPem = strings.TrimSpace(string(pem.EncodeToMemory(certPEMBlock)))
		marshaledKey := x509.MarshalPKCS1PrivateKey(key)
		keyPEMBlock := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: marshaledKey,
		}
		privRSAKeyPem = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
		marshaledKey, err = x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			panic(err)
		}
		keyPEMBlock = &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: marshaledKey,
		}
		privRSA8KeyPem = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
	}

	// Ed25519 generation
	{
		pubkey, privkey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
		subjKeyID, err := GetSubjKeyID(privkey)
		if err != nil {
			panic(err)
		}
		certTemplate := &x509.Certificate{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			SubjectKeyId: subjKeyID,
			DNSNames:     []string{"localhost"},
			ExtKeyUsage: []x509.ExtKeyUsage{
				x509.ExtKeyUsageServerAuth,
				x509.ExtKeyUsageClientAuth,
			},
			KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageKeyAgreement,
			SerialNumber: big.NewInt(mathrand.Int63()),
			NotBefore:    time.Now().Add(-30 * time.Second),
			NotAfter:     time.Now().Add(262980 * time.Hour),
		}
		csrTemplate := &x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName: "localhost",
			},
			DNSNames: []string{"localhost"},
		}
		csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, privkey)
		if err != nil {
			panic(err)
		}
		csrPEMBlock := &pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: csrBytes,
		}
		csrEd25519Pem = strings.TrimSpace(string(pem.EncodeToMemory(csrPEMBlock)))
		certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, intCert, pubkey, intKey)
		if err != nil {
			panic(err)
		}
		certPEMBlock := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		}
		certEd25519Pem = strings.TrimSpace(string(pem.EncodeToMemory(certPEMBlock)))
		marshaledKey, err := x509.MarshalPKCS8PrivateKey(privkey)
		if err != nil {
			panic(err)
		}
		keyPEMBlock := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: marshaledKey,
		}
		privEd255198KeyPem = strings.TrimSpace(string(pem.EncodeToMemory(keyPEMBlock)))
	}

	issuingCaChainPem = []string{intCertPEM, caCertPEM}
}

func TestComparePublicKeysAndType(t *testing.T) {
	rsa1 := genRsaKey(t).Public()
	rsa := genRsaKey(t).Public()
	eddsa1 := genEdDSA(t).Public()
	eddsa2 := genEdDSA(t).Public()
	ed25519_1, _ := genEd25519Key(t)
	ed25519_2, _ := genEd25519Key(t)

	type args struct {
		key1Iface crypto.PublicKey
		key2Iface crypto.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "RSA_Equal", args: args{key1Iface: rsa1, key2Iface: rsa1}, want: true, wantErr: false},
		{name: "RSA_NotEqual", args: args{key1Iface: rsa1, key2Iface: rsa}, want: false, wantErr: false},
		{name: "EDDSA_Equal", args: args{key1Iface: eddsa1, key2Iface: eddsa1}, want: true, wantErr: false},
		{name: "EDDSA_NotEqual", args: args{key1Iface: eddsa1, key2Iface: eddsa2}, want: false, wantErr: false},
		{name: "ED25519_Equal", args: args{key1Iface: ed25519_1, key2Iface: ed25519_1}, want: true, wantErr: false},
		{name: "ED25519_NotEqual", args: args{key1Iface: ed25519_1, key2Iface: ed25519_2}, want: false, wantErr: false},
		{name: "Mismatched_RSA", args: args{key1Iface: rsa1, key2Iface: ed25519_2}, want: false, wantErr: false},
		{name: "Mismatched_EDDSA", args: args{key1Iface: ed25519_1, key2Iface: rsa1}, want: false, wantErr: false},
		{name: "Mismatched_ED25519", args: args{key1Iface: ed25519_1, key2Iface: rsa1}, want: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ComparePublicKeysAndType(tt.args.key1Iface, tt.args.key2Iface)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComparePublicKeysAndType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ComparePublicKeysAndType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotAfterValues(t *testing.T) {
	if ErrNotAfterBehavior != 0 {
		t.Fatalf("Expected ErrNotAfterBehavior=%v to have value 0", ErrNotAfterBehavior)
	}

	if TruncateNotAfterBehavior != 1 {
		t.Fatalf("Expected TruncateNotAfterBehavior=%v to have value 1", TruncateNotAfterBehavior)
	}

	if PermitNotAfterBehavior != 2 {
		t.Fatalf("Expected PermitNotAfterBehavior=%v to have value 2", PermitNotAfterBehavior)
	}

	if AlwaysEnforceErr != 3 {
		t.Fatalf("Expected AlwaysEnforceErr=%v to have value 3", AlwaysEnforceErr)
	}
}

func TestSignatureAlgorithmRoundTripping(t *testing.T) {
	for leftName, value := range SignatureAlgorithmNames {
		if leftName == "pureed25519" && value == x509.PureEd25519 {
			continue
		}

		rightName, present := InvSignatureAlgorithmNames[value]
		if !present {
			t.Fatalf("%v=%v is present in SignatureAlgorithmNames but not in InvSignatureAlgorithmNames", leftName, value)
		}

		if strings.ToLower(rightName) != leftName {
			t.Fatalf("%v=%v is present in SignatureAlgorithmNames but inverse for %v has different name: %v", leftName, value, value, rightName)
		}
	}

	for leftValue, name := range InvSignatureAlgorithmNames {
		rightValue, present := SignatureAlgorithmNames[strings.ToLower(name)]
		if !present {
			t.Fatalf("%v=%v is present in InvSignatureAlgorithmNames but not in SignatureAlgorithmNames", leftValue, name)
		}

		if rightValue != leftValue {
			t.Fatalf("%v=%v is present in InvSignatureAlgorithmNames but forwards for %v has different value: %v", leftValue, name, name, rightValue)
		}
	}
}

// TestParseBasicConstraintExtension Verify extension generation and parsing of x509 basic constraint extensions
// works as expected.
func TestBasicConstraintExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		isCA       bool
		maxPathLen int
	}{
		{"empty-seq", false, -1},
		{"just-ca-true", true, -1},
		{"just-ca-with-maxpathlen", true, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext, err := CreateBasicConstraintExtension(tt.isCA, tt.maxPathLen)
			if err != nil {
				t.Fatalf("failed generating basic extension: %v", err)
			}

			gotIsCa, gotMaxPathLen, err := ParseBasicConstraintExtension(ext)
			if err != nil {
				t.Fatalf("failed parsing basic extension: %v", err)
			}

			if tt.isCA != gotIsCa {
				t.Fatalf("expected isCa (%v) got isCa (%v)", tt.isCA, gotIsCa)
			}

			if tt.maxPathLen != gotMaxPathLen {
				t.Fatalf("expected maxPathLen (%v) got maxPathLen (%v)", tt.maxPathLen, gotMaxPathLen)
			}
		})
	}

	t.Run("bad-extension-oid", func(t *testing.T) {
		// Test invalid type errors out
		_, _, err := ParseBasicConstraintExtension(pkix.Extension{})
		if err == nil {
			t.Fatalf("should have failed parsing non-basic constraint extension")
		}
	})

	t.Run("garbage-value", func(t *testing.T) {
		extraBytes, err := asn1.Marshal("a string")
		if err != nil {
			t.Fatalf("failed encoding the struct: %v", err)
		}
		ext := pkix.Extension{
			Id:    ExtensionBasicConstraintsOID,
			Value: extraBytes,
		}
		_, _, err = ParseBasicConstraintExtension(ext)
		if err == nil {
			t.Fatalf("should have failed parsing basic constraint with extra information")
		}
	})
}

// TestIgnoreCSRSigning Make sure we validate the CSR by default and that we can override
// the behavior disabling CSR signature checks
func TestIgnoreCSRSigning(t *testing.T) {
	t.Parallel()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed generating ca key: %v", err)
	}
	subjKeyID, err := GetSubjKeyID(caKey)
	if err != nil {
		t.Fatalf("failed generating ca subject key id: %v", err)
	}
	caCertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "root.localhost",
		},
		SubjectKeyId:          subjKeyID,
		DNSNames:              []string{"root.localhost"},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		SerialNumber:          big.NewInt(mathrand.Int63()),
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, caKey.Public(), caKey)
	if err != nil {
		t.Fatalf("failed creating ca certificate: %v", err)
	}
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		t.Fatalf("failed parsing ca certificate: %v", err)
	}

	signingBundle := &CAInfoBundle{
		ParsedCertBundle: ParsedCertBundle{
			PrivateKeyType:   ECPrivateKey,
			PrivateKey:       caKey,
			CertificateBytes: caBytes,
			Certificate:      caCert,
			CAChain:          nil,
		},
		URLs: &URLEntries{},
	}

	key := genEdDSA(t)
	csr := &x509.CertificateRequest{
		PublicKeyAlgorithm: x509.ECDSA,
		PublicKey:          key.Public(),
		Subject: pkix.Name{
			CommonName: "test.dadgarcorp.com",
		},
	}
	t.Run(fmt.Sprintf("ignore-csr-disabled"), func(t *testing.T) {
		params := &CreationParameters{
			URLs: &URLEntries{},
		}
		data := &CreationBundle{
			Params:        params,
			SigningBundle: signingBundle,
			CSR:           csr,
		}

		_, err := SignCertificate(data)
		if err == nil {
			t.Fatalf("should have failed signing csr with ignore csr signature disabled")
		}
		if !strings.Contains(err.Error(), "request signature invalid") {
			t.Fatalf("expected error to contain 'request signature invalid': got: %v", err)
		}
	})

	t.Run(fmt.Sprintf("ignore-csr-enabled"), func(t *testing.T) {
		params := &CreationParameters{
			IgnoreCSRSignature: true,
			URLs:               &URLEntries{},
		}
		data := &CreationBundle{
			Params:        params,
			SigningBundle: signingBundle,
			CSR:           csr,
		}

		cert, err := SignCertificate(data)
		if err != nil {
			t.Fatalf("failed to sign certificate: %v", err)
		}

		if err := cert.Verify(); err != nil {
			t.Fatalf("signature verification failed: %v", err)
		}
	})
}

// TestSignIntermediat_name_constraints verifies that all the name constraints extension fields are
// used when signing a certificate.
func TestSignCertificate_name_constraints(t *testing.T) {
	t.Parallel()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed generating ca key: %v", err)
	}
	subjKeyID, err := GetSubjKeyID(caKey)
	if err != nil {
		t.Fatalf("failed generating ca subject key id: %v", err)
	}
	caCertTemplate := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "root.localhost",
		},
		SubjectKeyId:          subjKeyID,
		DNSNames:              []string{"root.localhost"},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		SerialNumber:          big.NewInt(mathrand.Int63()),
		NotBefore:             time.Now().Add(-30 * time.Second),
		NotAfter:              time.Now().Add(262980 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, caCertTemplate, caCertTemplate, caKey.Public(), caKey)
	if err != nil {
		t.Fatalf("failed creating ca certificate: %v", err)
	}
	caCert, err := x509.ParseCertificate(caBytes)
	if err != nil {
		t.Fatalf("failed parsing ca certificate: %v", err)
	}

	signingBundle := &CAInfoBundle{
		ParsedCertBundle: ParsedCertBundle{
			PrivateKeyType:   ECPrivateKey,
			PrivateKey:       caKey,
			CertificateBytes: caBytes,
			Certificate:      caCert,
			CAChain:          nil,
		},
		URLs: &URLEntries{},
	}

	key := genEdDSA(t)
	csr := &x509.CertificateRequest{
		PublicKeyAlgorithm: x509.ECDSA,
		PublicKey:          key.Public(),
		Subject: pkix.Name{
			CommonName: "test.dadgarcorp.com",
		},
	}
	_, ipnet1, err := net.ParseCIDR("1.2.3.4/32")
	if err != nil {
		t.Fatal(err)
	}
	_, ipnet2, err := net.ParseCIDR("1.2.3.4/16")
	if err != nil {
		t.Fatal(err)
	}
	params := &CreationParameters{
		IgnoreCSRSignature:      true,
		URLs:                    &URLEntries{},
		NotAfter:                time.Now().Add(10000 * time.Hour),
		PermittedDNSDomains:     []string{"example.com", ".example.com"},
		ExcludedDNSDomains:      []string{"bad.example.com"},
		PermittedIPRanges:       []*net.IPNet{ipnet1},
		ExcludedIPRanges:        []*net.IPNet{ipnet2},
		PermittedEmailAddresses: []string{"one@example.com", "two@example.com"},
		ExcludedEmailAddresses:  []string{"un@example.com", "deux@example.com"},
		PermittedURIDomains:     []string{"domain1", "domain2"},
		ExcludedURIDomains:      []string{"domain3", "domain4"},
	}
	data := &CreationBundle{
		Params:        params,
		SigningBundle: signingBundle,
		CSR:           csr,
	}

	parsedBundle, err := SignCertificate(data)
	if err != nil {
		t.Fatal("should have failed signing csr with ignore csr signature disabled")
	}

	var failedChecks []error
	check := func(fieldName string, expected any, actual any) {
		diff := deep.Equal(expected, actual)
		if len(diff) > 0 {
			failedChecks = append(failedChecks, fmt.Errorf("error in field %q: %v", fieldName, diff))
		}
	}
	cert := parsedBundle.Certificate
	check("PermittedDNSDomains", params.PermittedDNSDomains, cert.PermittedDNSDomains)
	check("ExcludedDNSDomains", params.ExcludedDNSDomains, cert.ExcludedDNSDomains)
	check("PermittedIPRanges", params.PermittedIPRanges, cert.PermittedIPRanges)
	check("ExcludedIPRanges", params.ExcludedIPRanges, cert.ExcludedIPRanges)
	check("PermittedEmailAddresses", params.PermittedEmailAddresses, cert.PermittedEmailAddresses)
	check("ExcludedEmailAddresses", params.ExcludedEmailAddresses, cert.ExcludedEmailAddresses)
	check("PermittedURIDomains", params.PermittedURIDomains, cert.PermittedURIDomains)
	check("ExcludedURIDomains", params.ExcludedURIDomains, cert.ExcludedURIDomains)

	if err := errors.Join(failedChecks...); err != nil {
		t.Error(err)
	}
}

func genRsaKey(t *testing.T) *rsa.PrivateKey {
	key, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

func genEdDSA(t *testing.T) *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

func genEd25519Key(t *testing.T) (ed25519.PublicKey, ed25519.PrivateKey) {
	key, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key, priv
}

var (
	initTest           sync.Once
	privRSA8KeyPem     string
	privRSAKeyPem      string
	csrRSAPem          string
	certRSAPem         string
	privEd255198KeyPem string
	csrEd25519Pem      string
	certEd25519Pem     string
	privECKeyPem       string
	csrECPem           string
	privEC8KeyPem      string
	certECPem          string
	issuingCaChainPem  []string
)

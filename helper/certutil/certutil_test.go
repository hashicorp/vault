package certutil

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/api"
)

// Tests converting back and forth between a CertBundle and a ParsedCertBundle.
//
// Also tests the GetSubjKeyID, GetHexFormatted, and
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
			t.Errorf(err.Error())
		}

		cbut, err := pcbut.ToCertBundle()
		if err != nil {
			t.Fatalf("Error converting to cert bundle: %s", err)
		}

		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Fatalf(err.Error())
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
			t.Fatalf(err.Error())
		}

		secret := &api.Secret{
			Data: structs.New(cbut).Map(),
		}
		pcbut, err = ParsePKIMap(secret.Data)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf("Error during JSON bundle handling: %s", err)
		}
		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf(err.Error())
		}

		pcbut, err = ParsePEMBundle(cbut.ToPEMBundle())
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf("Error during JSON bundle handling: %s", err)
		}
		err = compareCertBundleToParsedCertBundle(cbut, pcbut)
		if err != nil {
			t.Logf("Error occurred with bundle %d in test array (index %d).\n", i+1, i)
			t.Fatalf(err.Error())
		}
	}
}

func compareCertBundleToParsedCertBundle(cbut *CertBundle, pcbut *ParsedCertBundle) error {
	if cbut == nil {
		return fmt.Errorf("got nil bundle")
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
	default:
		return fmt.Errorf("certBundle has unknown private key type")
	}

	if cb.SerialNumber != GetHexFormatted(pcbut.Certificate.SerialNumber.Bytes(), ":") {
		return fmt.Errorf("bundle serial number does not match")
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
	}

	for _, csrbut := range csrbuts {
		pcsrbut, err := csrbut.ToParsedCSRBundle()
		if err != nil {
			t.Fatalf("Error converting to parsed CSR bundle: %v", err)
		}

		err = compareCSRBundleToParsedCSRBundle(csrbut, pcsrbut)
		if err != nil {
			t.Fatalf(err.Error())
		}

		csrbut, err = pcsrbut.ToCSRBundle()
		if err != nil {
			t.Fatalf("Error converting to CSR bundle: %v", err)
		}

		err = compareCSRBundleToParsedCSRBundle(csrbut, pcsrbut)
		if err != nil {
			t.Fatalf(err.Error())
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
		key, err := rsa.GenerateKey(rand.Reader, 2048)
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

	issuingCaChainPem = []string{intCertPEM, caCertPEM}
}

var (
	initTest          sync.Once
	privRSA8KeyPem    string
	privRSAKeyPem     string
	csrRSAPem         string
	certRSAPem        string
	privECKeyPem      string
	csrECPem          string
	privEC8KeyPem     string
	certECPem         string
	issuingCaChainPem []string
)

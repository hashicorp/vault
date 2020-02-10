package mongodb

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

type Cert struct {
	Key    *rsa.PrivateKey
	KeyPem []byte

	Cert    *x509.Certificate
	CertPem []byte

	TLSCert tls.Certificate
}

func MakeCert(t *testing.T, parent *x509.Certificate) (bundle Cert) {
	key, keyPem := makeKey(t)

	now := time.Now()

	template := &x509.Certificate{
		IsCA:         false,
		SerialNumber: makeSerial(t),
		Subject: pkix.Name{
			CommonName: "unittest",
		},

		NotBefore: now,
		NotAfter:  now.Add(24 * time.Hour),

		KeyUsage: x509.KeyUsageDigitalSignature |
			x509.KeyUsageKeyEncipherment |
			x509.KeyUsageKeyAgreement,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SubjectKeyId: getSubjKeyID(t, key),
		DNSNames:     []string{"localhost"},
	}

	if parent == nil {
		parent = template
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, parent, key.Public(), key)
	if err != nil {
		t.Fatalf("Unable to generate cert: %s", err)
	}
	x509Cert, err := x509.ParseCertificate(cert)
	if err != nil {
		t.Fatalf("Unable to generate cert: %s", err)
	}

	certPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert,
		},
	)

	tlsCert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		t.Fatalf("Unable to parse X509 key pair: %s", err)
	}

	bundle = Cert{
		Key:     key,
		KeyPem:  keyPem,
		Cert:    x509Cert,
		CertPem: certPem,
		TLSCert: tlsCert,
	}
	return bundle
}

func makeKey(t *testing.T) (key *rsa.PrivateKey, pemBytes []byte) {
	t.Helper()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Unable to generate key for cert: %s", err)
	}

	privKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	)

	return privKey, privKeyPem
}

func makeSerial(t *testing.T) *big.Int {
	v := &big.Int{}
	serialNumberLimit := v.Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		t.Fatalf("Unable to generate serial number: %s", err)
	}
	return serialNumber
}

// Pulled from sdk/helper/certutil & slightly modified for test usage
func getSubjKeyID(t *testing.T, privateKey crypto.Signer) []byte {
	t.Helper()

	if privateKey == nil {
		t.Fatalf("passed-in private key is nil")
	}

	marshaledKey, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		t.Fatalf("error marshalling public key: %s", err)
	}

	subjKeyID := sha1.Sum(marshaledKey)

	return subjKeyID[:]
}

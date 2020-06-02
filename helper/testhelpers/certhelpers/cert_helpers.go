package certhelpers

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

type CertBuilder struct {
	tmpl       *x509.Certificate
	parentTmpl *x509.Certificate

	selfSign  bool
	parentKey *rsa.PrivateKey

	isCA bool
}

type CertOpt func(*CertBuilder) error

func CommonName(cn string) CertOpt {
	return func(builder *CertBuilder) error {
		builder.tmpl.Subject.CommonName = cn
		return nil
	}
}

func Parent(parent Certificate) CertOpt {
	return func(builder *CertBuilder) error {
		builder.parentKey = parent.privKey.privKey
		builder.parentTmpl = parent.Template
		return nil
	}
}

func IsCA(isCA bool) CertOpt {
	return func(builder *CertBuilder) error {
		builder.isCA = isCA
		return nil
	}
}

func SelfSign() CertOpt {
	return func(builder *CertBuilder) error {
		builder.selfSign = true
		return nil
	}
}

func Ip(ip ...string) CertOpt {
	return func(builder *CertBuilder) error {
		for _, addr := range ip {
			if ipAddr := net.ParseIP(addr); ipAddr != nil {
				builder.tmpl.IPAddresses = append(builder.tmpl.IPAddresses, ipAddr)
			}
		}
		return nil
	}
}

func Dns(dns ...string) CertOpt {
	return func(builder *CertBuilder) error {
		builder.tmpl.DNSNames = dns
		return nil
	}
}

func NewCert(t *testing.T, opts ...CertOpt) (cert Certificate) {
	t.Helper()

	builder := CertBuilder{
		tmpl: &x509.Certificate{
			SerialNumber: makeSerial(t),
			Subject: pkix.Name{
				CommonName: makeCommonName(),
			},
			NotBefore: time.Now().Add(-1 * time.Hour),
			NotAfter:  time.Now().Add(1 * time.Hour),
			IsCA:      false,
			KeyUsage: x509.KeyUsageDigitalSignature |
				x509.KeyUsageKeyEncipherment |
				x509.KeyUsageKeyAgreement,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		},
	}

	for _, opt := range opts {
		err := opt(&builder)
		if err != nil {
			t.Fatalf("Failed to set up certificate builder: %s", err)
		}
	}

	key := NewPrivateKey(t)

	builder.tmpl.SubjectKeyId = getSubjKeyID(t, key.privKey)

	tmpl := builder.tmpl
	parent := builder.parentTmpl
	publicKey := key.privKey.Public()
	signingKey := builder.parentKey

	if builder.selfSign {
		parent = tmpl
		signingKey = key.privKey
	}

	if builder.isCA {
		tmpl.IsCA = true
		tmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageCRLSign
		tmpl.ExtKeyUsage = nil
	} else {
		tmpl.KeyUsage = x509.KeyUsageDigitalSignature |
			x509.KeyUsageKeyEncipherment |
			x509.KeyUsageKeyAgreement
		tmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, tmpl, parent, publicKey, signingKey)
	if err != nil {
		t.Fatalf("Unable to generate certificate: %s", err)
	}
	certPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	tlsCert, err := tls.X509KeyPair(certPem, key.Pem)
	if err != nil {
		t.Fatalf("Unable to parse X509 key pair: %s", err)
	}

	return Certificate{
		Template: tmpl,
		privKey:  key,
		TLSCert:  tlsCert,
		rawCert:  certBytes,
		Pem:      certPem,
		isCA:     builder.isCA,
	}
}

// ////////////////////////////////////////////////////////////////////////////
// Private Key
// ////////////////////////////////////////////////////////////////////////////
type KeyWrapper struct {
	privKey *rsa.PrivateKey
	Pem     []byte
}

func NewPrivateKey(t *testing.T) (key KeyWrapper) {
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

	key = KeyWrapper{
		privKey: privKey,
		Pem:     privKeyPem,
	}

	return key
}

// ////////////////////////////////////////////////////////////////////////////
// Certificate
// ////////////////////////////////////////////////////////////////////////////
type Certificate struct {
	privKey  KeyWrapper
	Template *x509.Certificate
	TLSCert  tls.Certificate
	rawCert  []byte
	Pem      []byte
	isCA     bool
}

func (cert Certificate) CombinedPEM() []byte {
	if cert.isCA {
		return cert.Pem
	}
	return bytes.Join([][]byte{cert.privKey.Pem, cert.Pem}, []byte{'\n'})
}

func (cert Certificate) PrivateKeyPEM() []byte {
	return cert.privKey.Pem
}

// ////////////////////////////////////////////////////////////////////////////
// Writing to file
// ////////////////////////////////////////////////////////////////////////////
func WriteFile(t *testing.T, filename string, data []byte, perms os.FileMode) {
	t.Helper()

	err := ioutil.WriteFile(filename, data, perms)
	if err != nil {
		t.Fatalf("Unable to write to file [%s]: %s", filename, err)
	}
}

// ////////////////////////////////////////////////////////////////////////////
// Helpers
// ////////////////////////////////////////////////////////////////////////////
func makeSerial(t *testing.T) *big.Int {
	t.Helper()

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

func makeCommonName() (cn string) {
	return strings.ReplaceAll(time.Now().Format("20060102T150405.000"), ".", "")
}

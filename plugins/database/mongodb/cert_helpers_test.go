// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodb

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
	"math/big"
	"strings"
	"testing"
	"time"
)

type certBuilder struct {
	tmpl       *x509.Certificate
	parentTmpl *x509.Certificate

	selfSign  bool
	parentKey *rsa.PrivateKey

	isCA bool
}

type certOpt func(*certBuilder) error

func commonName(cn string) certOpt {
	return func(builder *certBuilder) error {
		builder.tmpl.Subject.CommonName = cn
		return nil
	}
}

func parent(parent certificate) certOpt {
	return func(builder *certBuilder) error {
		builder.parentKey = parent.privKey.privKey
		builder.parentTmpl = parent.template
		return nil
	}
}

func isCA(isCA bool) certOpt {
	return func(builder *certBuilder) error {
		builder.isCA = isCA
		return nil
	}
}

func selfSign() certOpt {
	return func(builder *certBuilder) error {
		builder.selfSign = true
		return nil
	}
}

func dns(dns ...string) certOpt {
	return func(builder *certBuilder) error {
		builder.tmpl.DNSNames = dns
		return nil
	}
}

func newCert(t *testing.T, opts ...certOpt) (cert certificate) {
	t.Helper()

	builder := certBuilder{
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

	key := newPrivateKey(t)

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

	tlsCert, err := tls.X509KeyPair(certPem, key.pem)
	if err != nil {
		t.Fatalf("Unable to parse X509 key pair: %s", err)
	}

	return certificate{
		template: tmpl,
		privKey:  key,
		tlsCert:  tlsCert,
		rawCert:  certBytes,
		pem:      certPem,
		isCA:     builder.isCA,
	}
}

// ////////////////////////////////////////////////////////////////////////////
// Private Key
// ////////////////////////////////////////////////////////////////////////////
type keyWrapper struct {
	privKey *rsa.PrivateKey
	pem     []byte
}

func newPrivateKey(t *testing.T) (key keyWrapper) {
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

	key = keyWrapper{
		privKey: privKey,
		pem:     privKeyPem,
	}

	return key
}

// ////////////////////////////////////////////////////////////////////////////
// Certificate
// ////////////////////////////////////////////////////////////////////////////
type certificate struct {
	privKey  keyWrapper
	template *x509.Certificate
	tlsCert  tls.Certificate
	rawCert  []byte
	pem      []byte
	isCA     bool
}

func (cert certificate) CombinedPEM() []byte {
	if cert.isCA {
		return cert.pem
	}
	return bytes.Join([][]byte{cert.privKey.pem, cert.pem}, []byte{'\n'})
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

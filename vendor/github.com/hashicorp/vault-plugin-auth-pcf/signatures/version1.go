package signatures

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hashicorp/go-multierror"
)

const TimeFormat = "2006-01-02T15:04:05Z"

type SignatureData struct {
	SigningTime time.Time
	Role        string
	Certificate string
}

func (s *SignatureData) hash() []byte {
	sum := sha256.Sum256([]byte(s.toSign()))
	return sum[:]
}

func (s *SignatureData) toSign() string {
	toHash := ""
	for _, field := range []string{s.SigningTime.UTC().Format(TimeFormat), s.Certificate, s.Role} {
		toHash += field
	}
	return toHash
}

func Sign(pathToPrivateKey string, signatureData *SignatureData) (string, error) {
	if signatureData == nil {
		return "", errors.New("signatureData must be provided")
	}

	keyBytes, err := ioutil.ReadFile(pathToPrivateKey)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return "", fmt.Errorf("unable to decode RSA private key from %s", keyBytes)
	}
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	signatureBytes, err := rsa.SignPSS(rand.Reader, rsaPrivateKey, crypto.SHA256, signatureData.hash(), nil)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(signatureBytes), nil
}

// Verify ensures that a given signature was created by one of the private keys
// matching one of the given client certificates. It is possible for a client
// certificate string given by PCF to contain multiple certificates within its
// body, hence the looping. The matching certificate is returned and should be
// further checked to ensure it contains the app, space, and org ID, and CN;
// otherwise it would be possible to match against an injected client certificate
// to gain authentication.
func Verify(signature string, signatureData *SignatureData) (*x509.Certificate, error) {
	if signatureData == nil {
		return nil, errors.New("signatureData must be provided")
	}

	// Use the CA certificate to verify the signature we've received.
	signatureBytes, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return nil, err
	}

	certBytes := []byte(signatureData.Certificate)
	var block *pem.Block
	var result error
	for {
		block, certBytes = pem.Decode(certBytes)
		if block == nil {
			break
		}
		clientCerts, err := x509.ParseCertificates(block.Bytes)
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}
		for _, clientCert := range clientCerts {
			publicKey, ok := clientCert.PublicKey.(*rsa.PublicKey)
			if !ok {
				result = multierror.Append(result, fmt.Errorf("not an rsa public key, it's a %t", clientCert.PublicKey))
				continue
			}

			if err := rsa.VerifyPSS(publicKey, crypto.SHA256, signatureData.hash(), signatureBytes, nil); err != nil {
				result = multierror.Append(result, err)
				continue
			}
			// Success
			return clientCert, nil
		}
	}
	if result == nil {
		return nil, fmt.Errorf("no matching client certificate found for %s in %s", signature, signatureData.Certificate)
	}
	return nil, result
}

func IsIssuer(pathToCACert string, clientCert *x509.Certificate) (bool, error) {
	caCertBytes, err := ioutil.ReadFile(pathToCACert)
	if err != nil {
		return false, err
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM(caCertBytes); !ok {
		return false, errors.New("couldn't append CA certificates")
	}

	verifyOpts := x509.VerifyOptions{
		Roots: pool,
	}

	if _, err := clientCert.Verify(verifyOpts); err != nil {
		return false, err
	}
	// Success
	return true, nil
}

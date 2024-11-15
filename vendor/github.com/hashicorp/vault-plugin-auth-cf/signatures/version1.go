// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
)

const TimeFormat = "2006-01-02T15:04:05Z"
const signatureVersion = "v1"

type SignatureData struct {
	SigningTime time.Time
	Role        string

	// CFInstanceCertContents are the full contents/body of the file
	// available at CF_INSTANCE_CERT. When viewed visually, this file
	// will contain two certificates. Generally, the first one is the
	// identity certificate itself, and the second one is the intermediate
	// certificate that issued it.
	CFInstanceCertContents string
}

func (s *SignatureData) hash() []byte {
	sum := sha256.Sum256([]byte(s.toSign()))
	return sum[:]
}

func (s *SignatureData) toSign() string {
	toHash := ""
	for _, field := range []string{s.SigningTime.UTC().Format(TimeFormat), s.CFInstanceCertContents, s.Role} {
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

	// This resolves to using a saltLength of 222.
	signatureBytes, err := rsa.SignPSS(rand.Reader, rsaPrivateKey, crypto.SHA256, signatureData.hash(), nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", signatureVersion, base64.StdEncoding.EncodeToString(signatureBytes)), nil
}

// Verify ensures that a given signature was created by a private key
// matching one of the given instance certificates. It returns the matching
// certificate, which should further be verified to be the identity certificate,
// and to be issued by a chain leading to the root CA certificate. There's a
// util function for this named Validate.
func Verify(signature string, signatureData *SignatureData) (*x509.Certificate, error) {
	var signatureBytes []byte
	var err error

	if signatureData == nil {
		return nil, errors.New("signatureData must be provided")
	}

	// Parse signature format
	parts := strings.Split(signature, ":")

	switch len(parts) {
	// Original release using URL-safe encoding and no embedded version
	case 1:
		signatureBytes, err = base64.URLEncoding.DecodeString(parts[0])
		if err != nil {
			return nil, err
		}
	case 2:
		if parts[0] != "v1" {
			return nil, fmt.Errorf("invalid signature version %q", parts[0])
		}
		signatureBytes, err = base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid signature format")
	}

	// Use the CA certificate to verify the signature we've received.
	cfInstanceCertContentsBytes := []byte(signatureData.CFInstanceCertContents)
	var block *pem.Block
	var result error
	for {
		block, cfInstanceCertContentsBytes = pem.Decode(cfInstanceCertContentsBytes)
		if block == nil {
			break
		}
		instanceCerts, err := x509.ParseCertificates(block.Bytes)
		if err != nil {
			result = multierror.Append(result, err)
			continue
		}
		for _, instanceCert := range instanceCerts {
			publicKey, ok := instanceCert.PublicKey.(*rsa.PublicKey)
			if !ok {
				result = multierror.Append(result, fmt.Errorf("not an rsa public key, it's a %t", instanceCert.PublicKey))
				continue
			}
			if err := rsa.VerifyPSS(publicKey, crypto.SHA256, signatureData.hash(), signatureBytes, nil); err != nil {
				result = multierror.Append(result, err)
				continue
			}
			// Success
			return instanceCert, nil
		}
	}
	if result == nil {
		return nil, fmt.Errorf("no matching certificate found for %s in %s", signature, signatureData.CFInstanceCertContents)
	}
	return nil, result
}

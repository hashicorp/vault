// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/mlkem"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/errutil"
)

// generateExportKeypair returns privPEM, pubPEM, and public key bytes for the given keyType.
// RSA and EC use x509.MarshalPKCS8PrivateKey / x509.MarshalPKIXPublicKey, emitting standard
// "PRIVATE KEY" / "PUBLIC KEY" PEM. ML-KEM uses custom PEM labels ("ML-KEM PRIVATE KEY" /
// "ML-KEM PUBLIC KEY") with raw bytes because x509.MarshalPKIXPublicKey and
// x509.MarshalPKCS8PrivateKey do not support *mlkem.EncapsulationKey* / *mlkem.DecapsulationKey*
// in Go 1.27.
func generateExportKeypair(keyType string) (privPEM string, pubPEM string, pubDER []byte, err error) {
	switch keyType {
	case "ml-kem-768":
		dk, genErr := mlkem.GenerateKey768()
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate ML-KEM-768 key: %v", genErr)}
		}
		pubBytes := dk.EncapsulationKey().Bytes()
		privPEM = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "ML-KEM PRIVATE KEY", Bytes: dk.Bytes()})))
		pubPEM = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "ML-KEM PUBLIC KEY", Bytes: pubBytes})))
		return privPEM, pubPEM, pubBytes, nil

	case "ml-kem-1024":
		dk, genErr := mlkem.GenerateKey1024()
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate ML-KEM-1024 key: %v", genErr)}
		}
		pubBytes := dk.EncapsulationKey().Bytes()
		privPEM = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "ML-KEM PRIVATE KEY", Bytes: dk.Bytes()})))
		pubPEM = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "ML-KEM PUBLIC KEY", Bytes: pubBytes})))
		return privPEM, pubPEM, pubBytes, nil
	}

	var privKey interface{}
	var pubKey interface{}

	switch keyType {
	case "rsa-2048":
		k, genErr := rsa.GenerateKey(rand.Reader, 2048)
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate RSA-2048 key: %v", genErr)}
		}
		privKey, pubKey = k, &k.PublicKey
	case "rsa-3072":
		k, genErr := rsa.GenerateKey(rand.Reader, 3072)
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate RSA-3072 key: %v", genErr)}
		}
		privKey, pubKey = k, &k.PublicKey
	case "rsa-4096":
		k, genErr := rsa.GenerateKey(rand.Reader, 4096)
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate RSA-4096 key: %v", genErr)}
		}
		privKey, pubKey = k, &k.PublicKey
	case "ec-p256":
		k, genErr := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate EC P-256 key: %v", genErr)}
		}
		privKey, pubKey = k, &k.PublicKey
	case "ec-p384":
		k, genErr := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate EC P-384 key: %v", genErr)}
		}
		privKey, pubKey = k, &k.PublicKey
	case "ec-p521":
		k, genErr := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate EC P-521 key: %v", genErr)}
		}
		privKey, pubKey = k, &k.PublicKey
	default:
		return "", "", nil, errutil.UserError{Err: fmt.Sprintf("unsupported key_type %q", keyType)}
	}

	privDER, marshalErr := x509.MarshalPKCS8PrivateKey(privKey)
	if marshalErr != nil {
		return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to marshal private key to PKCS8: %v", marshalErr)}
	}
	privPEM = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER})))

	pubDER, marshalErr = x509.MarshalPKIXPublicKey(pubKey)
	if marshalErr != nil {
		return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to marshal public key to PKIX: %v", marshalErr)}
	}
	pubPEM = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})))

	return privPEM, pubPEM, pubDER, nil
}

// computeExportKeyHMAC returns a deterministic "sha256:<hex>" fingerprint over the DER-encoded
// public key bytes. The fixed HMAC key is intentional — this is a lookup correlator.
func computeExportKeyHMAC(pubKeyDER []byte) string {
	const hmacKey = "vault-pki-export-key-hmac-v1"
	mac := hmac.New(sha256.New, []byte(hmacKey))
	mac.Write(pubKeyDER)
	return "sha256:" + hex.EncodeToString(mac.Sum(nil))
}

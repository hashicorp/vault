// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/mlkem"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"golang.org/x/crypto/hkdf"
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
		privPEM = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
			Type:    "ML-KEM PRIVATE KEY",
			Headers: map[string]string{"Algorithm": "ML-KEM-768"},
			Bytes:   dk.Bytes(),
		})))
		pubPEM = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "ML-KEM PUBLIC KEY", Bytes: pubBytes})))
		return privPEM, pubPEM, pubBytes, nil

	case "ml-kem-1024":
		dk, genErr := mlkem.GenerateKey1024()
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate ML-KEM-1024 key: %v", genErr)}
		}
		pubBytes := dk.EncapsulationKey().Bytes()
		privPEM = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{
			Type:    "ML-KEM PRIVATE KEY",
			Headers: map[string]string{"Algorithm": "ML-KEM-1024"},
			Bytes:   dk.Bytes(),
		})))
		pubPEM = strings.TrimSpace(string(pem.EncodeToMemory(&pem.Block{Type: "ML-KEM PUBLIC KEY", Bytes: pubBytes})))
		return privPEM, pubPEM, pubBytes, nil
	}

	var privKey interface{}
	var pubKey interface{}

	switch keyType {
	case "rsa-2048":
		k, genErr := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate RSA-2048 key: %v", genErr)}
		}
		privKey, pubKey = k, &k.PublicKey
	case "rsa-3072":
		k, genErr := cryptoutil.GenerateRSAKey(rand.Reader, 3072)
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate RSA-3072 key: %v", genErr)}
		}
		privKey, pubKey = k, &k.PublicKey
	case "rsa-4096":
		k, genErr := cryptoutil.GenerateRSAKey(rand.Reader, 4096)
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate RSA-4096 key: %v", genErr)}
		}
		privKey, pubKey = k, &k.PublicKey
	case "rsa-8192":
		k, genErr := cryptoutil.GenerateRSAKey(rand.Reader, 8192)
		if genErr != nil {
			return "", "", nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate RSA-8192 key: %v", genErr)}
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

// computeExportKeyHMACFromPEM decodes the first PEM block and calls
// computeExportKeyHMAC on its DER bytes. The caller must ensure pubKeyPEM is
// valid PEM; wrapCAPrivateKey already enforces this at the export endpoint so
// block will never be nil in practice.
func computeExportKeyHMACFromPEM(pubKeyPEM string) string {
	block, _ := pem.Decode([]byte(pubKeyPEM))
	if block == nil {
		// Should not be reachable: the export handler validates PEM before
		// calling this function.
		return ""
	}
	return computeExportKeyHMAC(block.Bytes)
}

// wrappedKeyBlob is the JSON envelope returned by wrapCAPrivateKey. All fields are
// base64url (standard base64 with padding, as produced by encoding/json).
//
// Wrapping algorithm:
//
//	RSA wrapping keys: RSA-OAEP-SHA256 encrypts a fresh 256-bit AES-GCM key
//	(the "content-encryption key", CEK). The CEK then AES-256-GCM-encrypts the
//	raw CA private-key PEM. The blob carries the OAEP-wrapped CEK, the GCM
//	nonce, and the GCM ciphertext+tag.
//
//	EC wrapping keys (P-256/384/521): ECDH-ES generates a shared secret with an
//	ephemeral sender key; HKDF-SHA256 (info="vault-pki-export-v1") derives a
//	256-bit AES-GCM CEK from that secret. The blob carries the ephemeral public
//	key (PKIX), GCM nonce, and GCM ciphertext+tag.
//
//	ML-KEM wrapping keys: ML-KEM Encapsulate produces a shared secret and a
//	KEM ciphertext. HKDF-SHA256 (info="vault-pki-export-v1") derives a 256-bit
//	AES-GCM CEK from the shared secret. The blob carries the KEM ciphertext,
//	GCM nonce, and GCM ciphertext+tag.
type wrappedKeyBlob struct {
	// Alg identifies the wrapping algorithm so the import side can select the
	// right unwrap path without guessing.
	Alg string `json:"alg"`

	// WrappedCEK holds the RSA-OAEP-wrapped content-encryption key (RSA only).
	WrappedCEK []byte `json:"wrapped_cek,omitempty"`

	// EphemeralPub is the PKIX-encoded ephemeral EC public key (EC only).
	EphemeralPub []byte `json:"ephemeral_pub,omitempty"`

	// KEMCiphertext is the ML-KEM encapsulation ciphertext (ML-KEM only).
	KEMCiphertext []byte `json:"kem_ciphertext,omitempty"`

	// Nonce is the 12-byte AES-GCM nonce.
	Nonce []byte `json:"nonce"`

	// Ciphertext is the AES-256-GCM ciphertext with the 16-byte tag appended.
	Ciphertext []byte `json:"ciphertext"`
}

// wrapCAPrivateKey encrypts caPrivKeyPEM with the public key in wrappingPubKeyPEM and
// returns the serialised wrappedKeyBlob as JSON bytes. The raw private key is
// never present in the return value or in any error message.
//
// Supported wrapping key types: RSA (OAEP-SHA256 + AES-256-GCM),
// EC P-256/P-384/P-521 (ECDH-ES + HKDF-SHA256 + AES-256-GCM),
// ML-KEM-768 / ML-KEM-1024 (KEM encapsulate + HKDF-SHA256 + AES-256-GCM).
func wrapCAPrivateKey(caPrivKeyPEM, wrappingPubKeyPEM string) ([]byte, error) {
	block, _ := pem.Decode([]byte(wrappingPubKeyPEM))
	if block == nil {
		return nil, errutil.UserError{Err: "public_key: no PEM block found"}
	}

	plaintext := []byte(caPrivKeyPEM)

	switch block.Type {
	case "PUBLIC KEY":
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("public_key: failed to parse PKIX public key: %v", err)}
		}
		switch k := pub.(type) {
		case *rsa.PublicKey:
			return wrapWithRSA(k, plaintext)
		case *ecdsa.PublicKey:
			return wrapWithEC(k, plaintext)
		default:
			return nil, errutil.UserError{Err: fmt.Sprintf("public_key: unsupported PKIX key type %T", pub)}
		}

	case "ML-KEM PUBLIC KEY":
		return wrapWithMLKEM(block.Bytes, plaintext)

	default:
		return nil, errutil.UserError{Err: fmt.Sprintf("public_key: unsupported PEM block type %q", block.Type)}
	}
}

// wrapWithRSA wraps plaintext using RSA-OAEP-SHA256 + AES-256-GCM.
func wrapWithRSA(pub *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	cek := make([]byte, 32)
	if _, err := rand.Read(cek); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate CEK: %v", err)}
	}

	wrappedCEK, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, cek, nil)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("RSA-OAEP encrypt failed: %v", err)}
	}

	ct, nonce, err := aesGCMSeal(cek, plaintext)
	if err != nil {
		return nil, err
	}

	return marshalBlob(wrappedKeyBlob{
		Alg:        "RSA-OAEP-SHA256+AES256GCM",
		WrappedCEK: wrappedCEK,
		Nonce:      nonce,
		Ciphertext: ct,
	})
}

// wrapWithEC wraps plaintext using ECDH-ES + HKDF-SHA256 + AES-256-GCM.
func wrapWithEC(pub *ecdsa.PublicKey, plaintext []byte) ([]byte, error) {
	// Convert the recipient public key to crypto/ecdh.
	recipientECDH, err := pub.ECDH()
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("failed to convert EC key to ecdh: %v", err)}
	}

	// Generate an ephemeral key on the same curve as the recipient key.
	var ephPriv *ecdh.PrivateKey
	switch pub.Curve {
	case elliptic.P384():
		ephPriv, err = ecdh.P384().GenerateKey(rand.Reader)
	case elliptic.P521():
		ephPriv, err = ecdh.P521().GenerateKey(rand.Reader)
	default: // P-256
		ephPriv, err = ecdh.P256().GenerateKey(rand.Reader)
	}
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate ephemeral EC key: %v", err)}
	}

	shared, err := ephPriv.ECDH(recipientECDH)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("ECDH failed: %v", err)}
	}

	cek := hkdfSHA256(shared, []byte("vault-pki-export-v1"), 32)

	ephPubDER, err := x509.MarshalPKIXPublicKey(ephPriv.PublicKey())
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("failed to marshal ephemeral public key: %v", err)}
	}

	ct, nonce, err := aesGCMSeal(cek, plaintext)
	if err != nil {
		return nil, err
	}

	return marshalBlob(wrappedKeyBlob{
		Alg:          "ECDH-ES+AES256GCM",
		EphemeralPub: ephPubDER,
		Nonce:        nonce,
		Ciphertext:   ct,
	})
}

// wrapWithMLKEM wraps plaintext using ML-KEM encapsulate + HKDF-SHA256 + AES-256-GCM.
// pubBytes is the raw encapsulation-key bytes (no ASN.1 wrapper).
func wrapWithMLKEM(pubBytes, plaintext []byte) ([]byte, error) {
	switch len(pubBytes) {
	case 1184: // ML-KEM-768
		ek, err := mlkem.NewEncapsulationKey768(pubBytes)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("public_key: invalid ML-KEM-768 encapsulation key: %v", err)}
		}
		shared, kemCT := ek.Encapsulate()
		cek := hkdfSHA256(shared, []byte("vault-pki-export-v1"), 32)
		ct, nonce, err := aesGCMSeal(cek, plaintext)
		if err != nil {
			return nil, err
		}
		return marshalBlob(wrappedKeyBlob{
			Alg:           "MLKEM768+AES256GCM",
			KEMCiphertext: kemCT,
			Nonce:         nonce,
			Ciphertext:    ct,
		})

	case 1568: // ML-KEM-1024
		ek, err := mlkem.NewEncapsulationKey1024(pubBytes)
		if err != nil {
			return nil, errutil.UserError{Err: fmt.Sprintf("public_key: invalid ML-KEM-1024 encapsulation key: %v", err)}
		}
		shared, kemCT := ek.Encapsulate()
		cek := hkdfSHA256(shared, []byte("vault-pki-export-v1"), 32)
		ct, nonce, err := aesGCMSeal(cek, plaintext)
		if err != nil {
			return nil, err
		}
		return marshalBlob(wrappedKeyBlob{
			Alg:           "MLKEM1024+AES256GCM",
			KEMCiphertext: kemCT,
			Nonce:         nonce,
			Ciphertext:    ct,
		})

	default:
		return nil, errutil.UserError{Err: fmt.Sprintf("public_key: unrecognised ML-KEM public key length %d (expected 1184 or 1568)", len(pubBytes))}
	}
}

// aesGCMSeal encrypts plaintext with AES-256-GCM using key and returns (ciphertext+tag, nonce, error).
func aesGCMSeal(key, plaintext []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, errutil.InternalError{Err: fmt.Sprintf("AES NewCipher failed: %v", err)}
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, errutil.InternalError{Err: fmt.Sprintf("AES-GCM init failed: %v", err)}
	}
	nonce = make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, nil, errutil.InternalError{Err: fmt.Sprintf("failed to generate GCM nonce: %v", err)}
	}
	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// hkdfSHA256 derives keyLen bytes from secret using RFC 5869 HKDF-SHA256 with
// no salt and the given info label.
func hkdfSHA256(secret, info []byte, keyLen int) []byte {
	r := hkdf.New(sha256.New, secret, nil, info)
	out := make([]byte, keyLen)
	if _, err := io.ReadFull(r, out); err != nil {
		// hkdf.New with a valid hash and keyLen ≤ 255*HashLen never errors.
		panic(fmt.Sprintf("hkdf read failed: %v", err))
	}
	return out
}

func marshalBlob(b wrappedKeyBlob) ([]byte, error) {
	out, err := json.Marshal(b)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("failed to marshal wrapped key blob: %v", err)}
	}
	return out, nil
}

// unwrapCAPrivateKey decrypts wrappedBlobJSON using the export private key in exportPrivKeyPEM
// and returns the unwrapped CA private key PEM.
func unwrapCAPrivateKey(wrappedBlobJSON []byte, exportPrivKeyPEM string) (string, error) {
	var blob wrappedKeyBlob
	if err := json.Unmarshal(wrappedBlobJSON, &blob); err != nil {
		return "", errutil.UserError{Err: fmt.Sprintf("wrapped_key: invalid JSON payload: %v", err)}
	}

	block, _ := pem.Decode([]byte(exportPrivKeyPEM))
	if block == nil {
		return "", errutil.InternalError{Err: "export private key: no PEM block found"}
	}

	var cek []byte
	switch blob.Alg {
	case "RSA-OAEP-SHA256+AES256GCM":
		if block.Type != "PRIVATE KEY" {
			return "", errutil.UserError{Err: fmt.Sprintf("wrapped_key algorithm %s requires RSA export private key", blob.Alg)}
		}
		privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", errutil.InternalError{Err: fmt.Sprintf("failed to parse export private key PKCS8: %v", err)}
		}
		rsaPriv, ok := privKey.(*rsa.PrivateKey)
		if !ok {
			return "", errutil.UserError{Err: "export private key is not an RSA key"}
		}
		decryptedCEK, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPriv, blob.WrappedCEK, nil)
		if err != nil {
			return "", errutil.UserError{Err: fmt.Sprintf("RSA-OAEP decryption failed: %v", err)}
		}
		cek = decryptedCEK

	case "ECDH-ES+AES256GCM":
		if block.Type != "PRIVATE KEY" {
			return "", errutil.UserError{Err: fmt.Sprintf("wrapped_key algorithm %s requires EC export private key", blob.Alg)}
		}
		privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", errutil.InternalError{Err: fmt.Sprintf("failed to parse export private key PKCS8: %v", err)}
		}
		ecPriv, ok := privKey.(*ecdsa.PrivateKey)
		if !ok {
			return "", errutil.UserError{Err: "export private key is not an EC key"}
		}
		privECDH, err := ecPriv.ECDH()
		if err != nil {
			return "", errutil.InternalError{Err: fmt.Sprintf("failed to convert EC private key to ECDH: %v", err)}
		}

		ephPubInterface, err := x509.ParsePKIXPublicKey(blob.EphemeralPub)
		if err != nil {
			return "", errutil.UserError{Err: fmt.Sprintf("failed to parse ephemeral public key: %v", err)}
		}
		// x509.ParsePKIXPublicKey always returns *ecdsa.PublicKey for EC keys,
		// even when the key was originally a *ecdh.PublicKey marshalled with
		// x509.MarshalPKIXPublicKey. The *ecdsa.PublicKey → ECDH() call below
		// converts it back to the crypto/ecdh representation needed for the
		// shared-secret computation.
		ephECDSA, ok := ephPubInterface.(*ecdsa.PublicKey)
		if !ok {
			return "", errutil.UserError{Err: "ephemeral public key is not an ECDSA key"}
		}
		ephECDH, err := ephECDSA.ECDH()
		if err != nil {
			return "", errutil.UserError{Err: fmt.Sprintf("failed to convert ephemeral EC public key to ECDH: %v", err)}
		}

		shared, err := privECDH.ECDH(ephECDH)
		if err != nil {
			return "", errutil.UserError{Err: fmt.Sprintf("ECDH shared secret generation failed: %v", err)}
		}
		cek = hkdfSHA256(shared, []byte("vault-pki-export-v1"), 32)

	case "MLKEM768+AES256GCM":
		if block.Type != "ML-KEM PRIVATE KEY" {
			return "", errutil.UserError{Err: fmt.Sprintf("wrapped_key algorithm %s requires ML-KEM private key", blob.Alg)}
		}
		dk, err := mlkem.NewDecapsulationKey768(block.Bytes)
		if err != nil {
			return "", errutil.InternalError{Err: fmt.Sprintf("failed to parse ML-KEM-768 decapsulation key: %v", err)}
		}
		shared, err := dk.Decapsulate(blob.KEMCiphertext)
		if err != nil {
			return "", errutil.UserError{Err: fmt.Sprintf("ML-KEM-768 decapsulation failed: %v", err)}
		}
		cek = hkdfSHA256(shared, []byte("vault-pki-export-v1"), 32)

	case "MLKEM1024+AES256GCM":
		if block.Type != "ML-KEM PRIVATE KEY" {
			return "", errutil.UserError{Err: fmt.Sprintf("wrapped_key algorithm %s requires ML-KEM private key", blob.Alg)}
		}
		dk, err := mlkem.NewDecapsulationKey1024(block.Bytes)
		if err != nil {
			return "", errutil.InternalError{Err: fmt.Sprintf("failed to parse ML-KEM-1024 decapsulation key: %v", err)}
		}
		shared, err := dk.Decapsulate(blob.KEMCiphertext)
		if err != nil {
			return "", errutil.UserError{Err: fmt.Sprintf("ML-KEM-1024 decapsulation failed: %v", err)}
		}
		cek = hkdfSHA256(shared, []byte("vault-pki-export-v1"), 32)

	default:
		return "", errutil.UserError{Err: fmt.Sprintf("unsupported wrapping algorithm %q in wrapped_key", blob.Alg)}
	}

	plaintext, err := aesGCMOpen(cek, blob.Nonce, blob.Ciphertext)
	if err != nil {
		return "", errutil.UserError{Err: fmt.Sprintf("AES-GCM decryption failed: %v", err)}
	}

	return string(plaintext), nil
}

// aesGCMOpen decrypts ciphertext with AES-256-GCM using key and nonce.
func aesGCMOpen(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("AES NewCipher failed: %v", err)}
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("AES-GCM init failed: %v", err)}
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errutil.UserError{Err: fmt.Sprintf("AES-GCM Open failed: %v", err)}
	}
	return plaintext, nil
}

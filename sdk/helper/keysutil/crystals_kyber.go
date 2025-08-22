package keysutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/kyber/kyber1024"
	"github.com/cloudflare/circl/kem/kyber/kyber512"
	"github.com/cloudflare/circl/kem/kyber/kyber768"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"golang.org/x/crypto/hkdf"
)

type kyberBox struct {
	s     kem.Scheme
	label string
}

func labelFor(kt KeyType) string {
	return kt.String() + "-aes256-gcm-v1" // "kyber512-aes256-gcm-v1", etc.
}

func newKyberBox(t KeyType) (kyberBox, error) {
	switch t {
	case KeyType_Kyber512:
		return kyberBox{s: kyber512.Scheme(), label: labelFor(t)}, nil
	case KeyType_Kyber768:
		return kyberBox{s: kyber768.Scheme(), label: labelFor(t)}, nil
	case KeyType_Kyber1024:
		return kyberBox{s: kyber1024.Scheme(), label: labelFor(t)}, nil
	default:
		return kyberBox{}, errutil.InternalError{Err: "failed to base64-decode Kyber public key"}
	}
}

func (k kyberBox) Encrypt(pk kem.PublicKey, plaintext, ad []byte) (capsule, nonce, ciphertext []byte, err error) {
	capsule, ss, err := k.s.Encapsulate(pk)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("encapsulate: %w", err)
	}

	key, err := k.deriveAES256Key(ss, capsule, ad)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("derive key: %w", err)
	}

	aead, n, err := k.newGCMWithNonce(key)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("gcm init: %w", err)
	}

	nonce = n
	ciphertext = aead.Seal(nil, nonce, plaintext, ad)
	return capsule, nonce, ciphertext, nil
}

func (k kyberBox) Decrypt(sk kem.PrivateKey, capsule, nonce, ciphertext, ad []byte) ([]byte, error) {
	if len(capsule) != k.s.CiphertextSize() {
		return nil, fmt.Errorf("invalid capsule length: got %d want %d", len(capsule), k.s.CiphertextSize())
	}

	ss, err := k.s.Decapsulate(sk, capsule)
	if err != nil {
		return nil, fmt.Errorf("decapsulate: %w", err)
	}

	key, err := k.deriveAES256Key(ss, capsule, ad)
	if err != nil {
		return nil, fmt.Errorf("derive key: %w", err)
	}

	aead, err := k.newGCM(key)
	if err != nil {
		return nil, fmt.Errorf("gcm init: %w", err)
	}

	if len(nonce) != aead.NonceSize() {
		return nil, fmt.Errorf("invalid nonce length: got %d want %d", len(nonce), aead.NonceSize())
	}

	plain, err := aead.Open(nil, nonce, ciphertext, ad)
	if err != nil {
		return nil, errors.New("decryption failed")
	}
	return plain, nil
}

func (k kyberBox) deriveAES256Key(secret, capsule, ad []byte) ([]byte, error) {
	// Bind KDF to transcript: label || H(capsule) || H(ad)
	hc := sha256.Sum256(capsule)
	ha := sha256.Sum256(ad)
	info := append(append([]byte(k.label), hc[:]...), ha[:]...)

	h := hkdf.New(sha256.New, secret, nil, info)
	key := make([]byte, 32)
	_, err := io.ReadFull(h, key)
	return key, err
}

func (k kyberBox) newGCMWithNonce(key []byte) (cipher.AEAD, []byte, error) {
	aead, err := k.newGCM(key)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("rand nonce: %w", err)
	}
	return aead, nonce, nil
}

func (k kyberBox) newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm: %w", err)
	}
	return aead, nil
}

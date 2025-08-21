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
	"github.com/cloudflare/circl/kem/kyber/kyber768"
	"golang.org/x/crypto/hkdf"
)

type kyberBox struct{ s kem.Scheme }

func newKyberBox() kyberBox { return kyberBox{s: kyber768.Scheme()} }

func (k kyberBox) Encrypt(pk kem.PublicKey, plaintext, ad []byte) (capsule, nonce, ciphertext []byte, err error) {
	capsule, ss, err := k.s.Encapsulate(pk)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("encapsulate: %w", err)
	}

	key, err := deriveAES256Key(ss, capsule, ad, "kyber768-aes256-gcm-v1")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("derive key: %w", err)
	}

	aead, n, err := newGCMWithNonce(key)
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

	key, err := deriveAES256Key(ss, capsule, ad, "kyber768-aes256-gcm-v1")
	if err != nil {
		return nil, fmt.Errorf("derive key: %w", err)
	}

	aead, err := newGCM(key)
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

func deriveAES256Key(secret, capsule, ad []byte, label string) ([]byte, error) {
	// Bind KDF to transcript: label || H(capsule) || H(ad)
	hc := sha256.Sum256(capsule)
	ha := sha256.Sum256(ad)
	info := append(append([]byte(label), hc[:]...), ha[:]...)

	h := hkdf.New(sha256.New, secret, nil, info)
	key := make([]byte, 32)
	_, err := io.ReadFull(h, key)
	return key, err
}

func newGCMWithNonce(key []byte) (cipher.AEAD, []byte, error) {
	aead, err := newGCM(key)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, fmt.Errorf("rand nonce: %w", err)
	}
	return aead, nonce, nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
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

package keysutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/kyber/kyber768"
	"golang.org/x/crypto/hkdf"
)

// kyberBox provides KEM+DEM encryption using Kyber768 and AES-256-GCM.
type kyberBox struct{ s kem.Scheme }

// newKyberBox returns a kyberBox using kyber768.
func newKyberBox() kyberBox {
	return kyberBox{s: kyber768.Scheme()}
}

// Encrypt encapsulates the shared secret, derives AES-256 key, and encrypts plaintext.
// Returns capsule, nonce, ciphertext.
func (k kyberBox) Encrypt(pk kem.PublicKey, plaintext, ad []byte) (capsule, nonce, ciphertext []byte, err error) {
	capsule, ss, err := k.s.Encapsulate(pk)
	if err != nil {
		return nil, nil, nil, err
	}
	key, err := deriveAES256Key(ss)
	if err != nil {
		return nil, nil, nil, err
	}
	aead, nonce, err := newGCMWithNonce(key)
	if err != nil {
		return nil, nil, nil, err
	}
	ciphertext = aead.Seal(nil, nonce, plaintext, ad)
	return capsule, nonce, ciphertext, nil
}

// Decrypt decapsulates to recover AES key and decrypts ciphertext.
func (k kyberBox) Decrypt(sk kem.PrivateKey, capsule, nonce, ciphertext, ad []byte) ([]byte, error) {
	ss, err := k.s.Decapsulate(sk, capsule)
	if err != nil {
		return nil, err
	}
	key, err := deriveAES256Key(ss)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(mustAES(key))
	if err != nil {
		return nil, err
	}
	return aead.Open(nil, nonce, ciphertext, ad)
}

func deriveAES256Key(secret []byte) ([]byte, error) {
	h := hkdf.New(sha256.New, secret, nil, []byte("Kyber→AES-256-GCM key"))
	key := make([]byte, 32)
	_, err := io.ReadFull(h, key)
	return key, err
}

func newGCMWithNonce(key []byte) (cipher.AEAD, []byte, error) {
	aead, err := cipher.NewGCM(mustAES(key))
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}
	return aead, nonce, nil
}

func mustAES(key []byte) cipher.Block {
	b, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return b
}

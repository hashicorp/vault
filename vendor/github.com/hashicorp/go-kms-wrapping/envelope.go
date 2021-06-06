package wrapping

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	fmt "fmt"

	uuid "github.com/hashicorp/go-uuid"
)

// NewEnvelope retuns an Envelope that is ready to use for use. It is valid to pass nil EnvelopeOptions.
func NewEnvelope(opts *EnvelopeOptions) *Envelope {
	return &Envelope{}
}

// Encrypt takes in plaintext and envelope encrypts it, generating an EnvelopeInfo value
func (e *Envelope) Encrypt(plaintext []byte, aad []byte) (*EnvelopeInfo, error) {
	// Generate DEK
	key, err := uuid.GenerateRandomBytes(32)
	if err != nil {
		return nil, err
	}
	iv, err := uuid.GenerateRandomBytes(12)
	if err != nil {
		return nil, err
	}
	aead, err := e.aeadEncrypter(key)
	if err != nil {
		return nil, err
	}

	return &EnvelopeInfo{
		Ciphertext: aead.Seal(nil, iv, plaintext, aad),
		Key:        key,
		IV:         iv,
	}, nil
}

// Decrypt takes in EnvelopeInfo and potentially additional data and decrypts. Additional data is separate from the encrypted blob info as it is expected that will be sourced from a separate location.
func (e *Envelope) Decrypt(data *EnvelopeInfo, aad []byte) ([]byte, error) {
	aead, err := e.aeadEncrypter(data.Key)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, data.IV, data.Ciphertext, aad)
}

func (e *Envelope) aeadEncrypter(key []byte) (cipher.AEAD, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create the GCM mode AEAD
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, errors.New("failed to initialize GCM mode")
	}

	return gcm, nil
}

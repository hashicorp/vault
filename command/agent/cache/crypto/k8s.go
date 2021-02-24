package crypto

import (
	"context"
	"crypto/rand"
	"fmt"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-kms-wrapping/wrappers/aead"
)

var _ KeyManager = (*KubeEncrypter)(nil)

type KubeEncrypter struct {
	renewable bool
	wrapper   *aead.Wrapper
}

// NewKubeEncrypter returns a new instance of the Kube encryption key. Kubernetes
// encryption keys aren't renewable.
func NewKubeEncrypter(existingKey []byte) (*KubeEncrypter, error) {
	k := &KubeEncrypter{
		renewable: false,
		wrapper:   aead.NewWrapper(nil),
	}

	k.wrapper.SetConfig(map[string]string{"key_id": KeyID})

	var rootKey []byte
	switch len(existingKey) {
	case 0:
		newKey := make([]byte, 32)
		_, err := rand.Read(newKey)
		if err != nil {
			return k, err
		}
		rootKey = newKey
	case 32:
		rootKey = existingKey
	default:
		return k, fmt.Errorf("unknown error")
	}

	if err := k.wrapper.SetAESGCMKeyBytes(rootKey); err != nil {
		return k, err
	}

	return k, nil
}

// GetKey returns the encryption key in a format optimized for storage.
// In k8s we store the key as is, so just return the key stored.
func (k *KubeEncrypter) GetKey() []byte {
	return k.wrapper.GetKeyBytes()
}

// GetPersistentKey returns the key which should be stored in the persisent
// cache. In k8s we store the key as is, so just return the key stored.
func (k *KubeEncrypter) GetPersistentKey() ([]byte, error) {
	return k.wrapper.GetKeyBytes(), nil
}

// Renewable lets the caller know if this encryption key type is
// renewable. In Kubernetes the key isn't renewable.
func (k *KubeEncrypter) Renewable() bool {
	return k.renewable
}

// Renewer is used when the encryption key type is renewable. Since Kubernetes
// keys aren't renewable, returning nothing.
func (k *KubeEncrypter) Renewer(ctx context.Context) error {
	return nil
}

// Encrypt takes plaintext values and encrypts them using the store key and additional
// data. For Kubernetes the AAD should be the service account JWT.
func (k *KubeEncrypter) Encrypt(ctx context.Context, plaintext, aad []byte) ([]byte, error) {
	blob, err := k.wrapper.Encrypt(ctx, plaintext, aad)
	if err != nil {
		return nil, err
	}
	return blob.Ciphertext, nil
}

// Decrypt takes ciphertext and AAD values and returns the decrypted value. For Kubernetes the AAD
// should be the service account JWT.
func (k *KubeEncrypter) Decrypt(ctx context.Context, ciphertext, aad []byte) ([]byte, error) {
	blob := &wrapping.EncryptedBlobInfo{
		Ciphertext: ciphertext,
		KeyInfo: &wrapping.KeyInfo{
			KeyID: KeyID,
		},
	}
	return k.wrapper.Decrypt(ctx, blob, aad)
}

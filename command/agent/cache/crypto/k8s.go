package crypto

import (
	"crypto/rand"
	"fmt"

	wrapping "github.com/hashicorp/go-kms-wrapping"
	"github.com/hashicorp/go-kms-wrapping/wrappers/aead"
)

var _ KeyManager = (*PassthroughKeyManager)(nil)

type PassthroughKeyManager struct {
	wrapper *aead.Wrapper
}

// NewPassthroughKeyManager returns a new instance of the Kube encryption key. Kubernetes
// encryption keys aren't renewable.
func NewPassthroughKeyManager(key []byte) (*PassthroughKeyManager, error) {

	var rootKey []byte = nil
	switch len(key) {
	case 0:
		newKey := make([]byte, 32)
		_, err := rand.Read(newKey)
		if err != nil {
			return nil, err
		}
		rootKey = newKey
	case 32:
		rootKey = key
	default:
		return nil, fmt.Errorf("invalid key size, should be 32, got %d", len(key))
	}

	wrapper := aead.NewWrapper(nil)

	if _, err := wrapper.SetConfig(map[string]string{"key_id": KeyID}); err != nil {
		return nil, err
	}

	if err := wrapper.SetAESGCMKeyBytes(rootKey); err != nil {
		return nil, err
	}

	k := &PassthroughKeyManager{
		wrapper: wrapper,
	}

	return k, nil
}

func (w *PassthroughKeyManager) Wrapper() wrapping.Wrapper {
	return w.wrapper
}

func (w *PassthroughKeyManager) RetrievalToken() ([]byte, error) {
	if w.wrapper == nil {
		return nil, fmt.Errorf("unable to get wrapper for token retrieval")
	}

	return w.wrapper.GetKeyBytes(), nil
}

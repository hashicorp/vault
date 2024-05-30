// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package keymanager

import (
	"context"
	"crypto/rand"
	"fmt"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/aead/v2"
)

var _ KeyManager = (*PassthroughKeyManager)(nil)

type PassthroughKeyManager struct {
	wrapper *aead.Wrapper
}

// NewPassthroughKeyManager returns a new instance of the Kube encryption key.
// If a key is provided, it will be used as the encryption key for the wrapper,
// otherwise one will be generated.
func NewPassthroughKeyManager(ctx context.Context, key []byte) (*PassthroughKeyManager, error) {
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

	wrapper := aead.NewWrapper()

	if _, err := wrapper.SetConfig(ctx, wrapping.WithConfigMap(map[string]string{"key_id": KeyID})); err != nil {
		return nil, err
	}

	if err := wrapper.SetAesGcmKeyBytes(rootKey); err != nil {
		return nil, err
	}

	k := &PassthroughKeyManager{
		wrapper: wrapper,
	}

	return k, nil
}

// Wrapper returns the manager's wrapper for key operations.
func (w *PassthroughKeyManager) Wrapper() wrapping.Wrapper {
	return w.wrapper
}

// RetrievalToken returns the key that was used on the wrapper since this key
// manager is simply a passthrough and does not provide a mechanism to abstract
// this key.
func (w *PassthroughKeyManager) RetrievalToken(ctx context.Context) ([]byte, error) {
	if w.wrapper == nil {
		return nil, fmt.Errorf("unable to get wrapper for token retrieval")
	}

	return w.wrapper.KeyBytes(ctx)
}

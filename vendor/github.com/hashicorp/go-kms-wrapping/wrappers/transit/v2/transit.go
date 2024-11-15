// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package transit

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

// Wrapper is a wrapper that leverages Vault's Transit secret
// engine
type Wrapper struct {
	logger       hclog.Logger
	client       transitClientEncryptor
	currentKeyId *atomic.Value
	keyIdPrefix  string
}

var _ wrapping.Wrapper = (*Wrapper)(nil)

// NewWrapper creates a new transit wrapper
func NewWrapper() *Wrapper {
	s := &Wrapper{
		currentKeyId: new(atomic.Value),
	}
	s.currentKeyId.Store("")
	return s
}

// SetConfig processes the config info from the server config
func (s *Wrapper) SetConfig(_ context.Context, opt ...wrapping.Option) (*wrapping.WrapperConfig, error) {
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, err
	}

	s.logger = opts.withLogger

	client, wrapConfig, err := newTransitClient(s.logger, opts)
	if err != nil {
		return nil, err
	}
	s.client = client
	s.keyIdPrefix = opts.withKeyIdPrefix

	// Send a value to test the wrapper and to set the current key id
	if _, err := s.Encrypt(context.Background(), []byte("a")); err != nil {
		client.Close()
		return nil, err
	}

	return wrapConfig, nil
}

// Init is called during core.Initialize
func (s *Wrapper) Init(_ context.Context) error {
	return nil
}

// Finalize is called during shutdown
func (s *Wrapper) Finalize(_ context.Context) error {
	s.client.Close()
	return nil
}

// Type returns the type for this particular Wrapper implementation
func (s *Wrapper) Type(_ context.Context) (wrapping.WrapperType, error) {
	return wrapping.WrapperTypeTransit, nil
}

// KeyId returns the last known key id
func (s *Wrapper) KeyId(_ context.Context) (string, error) {
	return s.currentKeyId.Load().(string), nil
}

// Encrypt is used to encrypt using Vault's Transit engine
func (s *Wrapper) Encrypt(ctx context.Context, plaintext []byte, _ ...wrapping.Option) (*wrapping.BlobInfo, error) {
	ciphertext, err := s.client.Encrypt(ctx, plaintext)
	if err != nil {
		return nil, err
	}

	splitKey := strings.Split(string(ciphertext), ":")
	if len(splitKey) != 3 {
		return nil, errors.New("invalid ciphertext returned")
	}
	keyId := s.keyIdPrefix + splitKey[1]
	s.currentKeyId.Store(keyId)

	ret := &wrapping.BlobInfo{
		Ciphertext: ciphertext,
		KeyInfo: &wrapping.KeyInfo{
			KeyId: keyId,
		},
	}
	return ret, nil
}

// Decrypt is used to decrypt the ciphertext
func (s *Wrapper) Decrypt(ctx context.Context, in *wrapping.BlobInfo, _ ...wrapping.Option) ([]byte, error) {
	plaintext, err := s.client.Decrypt(ctx, in.Ciphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// GetClient returns the transit Wrapper's transitClientEncryptor
func (s *Wrapper) GetClient() transitClientEncryptor {
	return s.client
}

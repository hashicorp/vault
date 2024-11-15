// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpkms

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	ErrKeyNotFound = errors.New("encryption key not found")
)

// Key represents a key from the storage backend.
type Key struct {
	// Name is the name of the key in Vault.
	Name string `json:"name"`

	// CryptoKeyID is the full resource ID of the key on GCP.
	CryptoKeyID string `json:"crypto_key_id"`

	// MinVersion is the minimum crypto key version to allow. If left unset or set
	// to a negative number, all versions are allowed.
	MinVersion int `json:"min_version"`

	// MaxVersion is the maximum crypto key version to allow. If left unset or set
	// to a negative number, all versions are allowed.
	MaxVersion int `json:"max_version"`
}

// Key retrieves the named key from the storage backend, or an error if one does
// not exist.
func (b *backend) Key(ctx context.Context, s logical.Storage, key string) (*Key, error) {
	entry, err := s.Get(ctx, "keys/"+key)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to retrieve key %q: {{err}}", key), err)
	}
	if entry == nil {
		return nil, ErrKeyNotFound
	}

	var result Key
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to decode entry for %q: {{err}}", key), err)
	}
	return &result, nil
}

// Keys returns the list of keys
func (b *backend) Keys(ctx context.Context, s logical.Storage) ([]string, error) {
	entries, err := s.List(ctx, "keys/")
	if err != nil {
		return nil, errwrap.Wrapf("failed to list keys: {{err}}", err)
	}
	return entries, nil
}

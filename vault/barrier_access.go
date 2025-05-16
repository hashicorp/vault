// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
)

// BarrierEncryptorAccess is a wrapper around BarrierEncryptor that allows Core
// to expose its barrier encrypt/decrypt operations through BarrierEncryptorAccess()
// while restricting the ability to modify Core.barrier itself.
type BarrierEncryptorAccess struct {
	barrierEncryptor BarrierEncryptor
}

var _ BarrierEncryptor = (*BarrierEncryptorAccess)(nil)

func NewBarrierEncryptorAccess(barrierEncryptor BarrierEncryptor) *BarrierEncryptorAccess {
	return &BarrierEncryptorAccess{barrierEncryptor: barrierEncryptor}
}

func (b *BarrierEncryptorAccess) Encrypt(ctx context.Context, key string, plaintext []byte) ([]byte, error) {
	return b.barrierEncryptor.Encrypt(ctx, key, plaintext)
}

func (b *BarrierEncryptorAccess) Decrypt(ctx context.Context, key string, ciphertext []byte) ([]byte, error) {
	return b.barrierEncryptor.Decrypt(ctx, key, ciphertext)
}

// NewBarrierDecryptingStorage returns a view of storage that will decrypt the
// storage values for Get operations. The returned storage is read-only, and
// will error on attempts to Put or Delete.
func NewBarrierDecryptingStorage(
	barrier BarrierEncryptor, underlying physical.Backend,
) logical.Storage {
	return &barrierDecryptingStorage{
		barrier:    barrier,
		underlying: underlying,
	}
}

type barrierDecryptingStorage struct {
	barrier    BarrierEncryptor
	underlying physical.Backend
}

func (b *barrierDecryptingStorage) List(ctx context.Context, prefix string) ([]string, error) {
	return b.underlying.List(ctx, prefix)
}

// ignore-nil-nil-function-check
func (b *barrierDecryptingStorage) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
	entry, err := b.underlying.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	decrypted, err := b.barrier.Decrypt(ctx, entry.Key, entry.Value)
	if err != nil {
		return nil, err
	}
	return &logical.StorageEntry{
		Key:      entry.Key,
		Value:    decrypted,
		SealWrap: entry.SealWrap,
	}, nil
}

var errReadOnly = errors.New("read-only storage")

func (b *barrierDecryptingStorage) Put(ctx context.Context, entry *logical.StorageEntry) error {
	return errReadOnly
}

func (b *barrierDecryptingStorage) Delete(ctx context.Context, key string) error {
	return errReadOnly
}

var _ logical.Storage = (*barrierDecryptingStorage)(nil)

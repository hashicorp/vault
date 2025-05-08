// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"

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

func (b *BarrierEncryptorAccess) EncryptUntracked(ctx context.Context, key string, plaintext []byte) ([]byte, error) {
	return b.barrierEncryptor.EncryptUntracked(ctx, key, plaintext)
}

func NewBarrierEncryptorStorage(
	barrier BarrierEncryptor, underlying physical.Backend,
) logical.Storage {
	return &barrierEncryptorStorage{
		barrier:    barrier,
		underlying: underlying,
	}
}

type barrierEncryptorStorage struct {
	barrier    BarrierEncryptor
	underlying physical.Backend
}

func (b *barrierEncryptorStorage) List(ctx context.Context, prefix string) ([]string, error) {
	return b.underlying.List(ctx, prefix)
}

// ignore-nil-nil-function-check
func (b *barrierEncryptorStorage) Get(ctx context.Context, key string) (*logical.StorageEntry, error) {
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

func (b *barrierEncryptorStorage) Put(ctx context.Context, entry *logical.StorageEntry) error {
	encrypted, err := b.barrier.EncryptUntracked(ctx, entry.Key, entry.Value)
	if err != nil {
		return err
	}
	pe := &physical.Entry{
		Key:      entry.Key,
		Value:    encrypted,
		SealWrap: entry.SealWrap,
	}
	return b.underlying.Put(ctx, pe)
}

func (b *barrierEncryptorStorage) Delete(ctx context.Context, key string) error {
	return b.underlying.Delete(ctx, key)
}

var _ logical.Storage = (*barrierEncryptorStorage)(nil)

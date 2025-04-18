// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package managed_key

import (
	"io"

	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/ssh"
)

type ManagedKeyInfo struct {
	publicKey ssh.PublicKey
	Name      NameKey
	Uuid      UUIDKey
}

type managedKeyId interface {
	String() string
}

type SSHManagedKeyView interface {
	BackendUUID() string
	GetManagedKeyView() (logical.ManagedKeySystemView, error)
	GetRandomReader() io.Reader
}

type (
	UUIDKey string
	NameKey string
)

func (u UUIDKey) String() string {
	return string(u)
}

func (n NameKey) String() string {
	return string(n)
}

func (m ManagedKeyInfo) PublicKey() ssh.PublicKey {
	return m.publicKey
}

func (m ManagedKeyInfo) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return nil, nil
}

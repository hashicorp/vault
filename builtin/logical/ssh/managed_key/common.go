// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package managed_key

import (
	"crypto"
	"io"

	"github.com/hashicorp/vault/sdk/logical"
)

type ManagedKeyInfo struct {
	PublicKey crypto.PublicKey
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

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

const operationPrefixSSH = "ssh"

type backend struct {
	*framework.Backend
	view        logical.Storage
	salt        *salt.Salt
	saltMutex   sync.RWMutex
	backendUUID string
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend(conf *logical.BackendConfig) (*backend, error) {
	var b backend
	b.view = conf.StorageView
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"verify",
				"public_key",
			},

			LocalStorage: []string{
				"otp/",
			},

			SealWrapStorage: []string{
				caPrivateKey,
				caPrivateKeyStoragePath,
				keysStoragePrefix,
			},
		},

		Paths: []*framework.Path{
			pathConfigZeroAddress(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
			pathLookup(&b),
			pathVerify(&b),
			pathConfigCA(&b),
			pathSign(&b),
			pathIssue(&b),
			pathFetchPublicKey(&b),
			pathCleanupKeys(&b),
		},

		Secrets: []*framework.Secret{
			secretOTP(&b),
		},

		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	b.backendUUID = conf.BackendUUID
	return &b, nil
}

func (b *backend) Salt(ctx context.Context) (*salt.Salt, error) {
	b.saltMutex.RLock()
	if b.salt != nil {
		defer b.saltMutex.RUnlock()
		return b.salt, nil
	}
	b.saltMutex.RUnlock()
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	if b.salt != nil {
		return b.salt, nil
	}
	salt, err := salt.NewSalt(ctx, b.view, &salt.Config{
		HashFunc: salt.SHA256Hash,
		Location: salt.DefaultLocation,
	})
	if err != nil {
		return nil, err
	}
	b.salt = salt
	return salt, nil
}

func (b *backend) invalidate(_ context.Context, key string) {
	switch key {
	case salt.DefaultLocation:
		b.saltMutex.Lock()
		defer b.saltMutex.Unlock()
		b.salt = nil
	}
}

func (b *backend) GetManagedKeyView() (logical.ManagedKeySystemView, error) {
	managedKeyView, ok := b.System().(logical.ManagedKeySystemView)
	if !ok {
		return nil, errutil.InternalError{Err: "unsupported system view"}
	}
	return managedKeyView, nil
}

func (b *backend) BackendUUID() string {
	return b.backendUUID
}

const backendHelp = `
The SSH backend generates credentials allowing clients to establish SSH
connections to remote hosts.

There are two variants of the backend, which generate different types of
credentials: One-Time Passwords (OTPs) and certificate authority. The desired behavior
is role-specific and chosen at role creation time with the 'key_type'
parameter.

Please see the backend documentation for a thorough description of both
types. The Vault team strongly recommends the OTP type.

After mounting this backend, before generating credentials, configure the
backend's lease behavior using the 'config/lease' endpoint and create roles
using the 'roles/' endpoint.
`

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki_backend

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/managed_key"
	"github.com/hashicorp/vault/sdk/logical"
)

type StorageContext interface {
	GetContext() context.Context
	GetStorage() logical.Storage

	UseLegacyBundleCaStorage() bool
	GetPkiManagedView() managed_key.PkiManagedKeyView
	CrlBuilder() CrlBuilderType
	GetCertificateCounter() issuing.CertificateCounter

	Logger() hclog.Logger
}

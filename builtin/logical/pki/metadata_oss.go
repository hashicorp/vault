// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pki

import (
	"context"
	"crypto/x509"
	"errors"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/logical"
)

var ErrMetadataIsEntOnly = errors.New("certificate metadata is only supported on Vault Enterprise")

func storeCertMetadata(ctx context.Context, storage logical.Storage, issuerId issuing.IssuerID, role string, certificate *x509.Certificate, certMetadata interface{}) error {
	return ErrMetadataIsEntOnly
}

func (b *backend) doTidyCertMetadata(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	return ErrMetadataIsEntOnly
}

func validateCertMetadataConfiguration(role *issuing.RoleEntry) error {
	return ErrMetadataIsEntOnly
}

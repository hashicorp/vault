// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pki

import (
	"context"
	"crypto/x509"
	"errors"
	"github.com/hashicorp/go-hclog"
	"math/big"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/logical"
)

var ErrMetadataIsEntOnly = errors.New("certificate metadata is only supported on Vault Enterprise")

func storeMetadata(ctx context.Context, storage logical.Storage, issuerId issuing.IssuerID, role string, certificate *x509.Certificate, metadata interface{}) error {
	return ErrMetadataIsEntOnly
}

func GetCertificateMetadata(ctx context.Context, storage logical.Storage, serialNumber *big.Int) (*CertificateMetadata, error) {
	return nil, ErrMetadataIsEntOnly
}

func (b *backend) doTidyCertMetadata(ctx context.Context, req *logical.Request, logger hclog.Logger, config *tidyConfig) error {
	return ErrMetadataIsEntOnly
}

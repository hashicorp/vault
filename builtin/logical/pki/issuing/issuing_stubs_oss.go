// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package issuing

import (
	"context"

	ctx509 "github.com/google/certificate-transparency-go/x509"
	"github.com/hashicorp/vault/sdk/logical"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func entSetCertVerifyOptions(ctx context.Context, storage logical.Storage, issuerId IssuerID, options *ctx509.VerifyOptions) error {
	return nil
}

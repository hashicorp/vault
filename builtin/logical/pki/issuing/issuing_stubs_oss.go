// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package issuing

import (
	ctx509 "github.com/google/certificate-transparency-go/x509"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

func entSetCertVerifyOptions(issuer *IssuerEntry, view logical.SystemView, options *ctx509.VerifyOptions) (bool, error) {
	return false, nil
}

func EntAdjustCreationBundle(view logical.SystemView, bundle *certutil.CreationBundle) {
	return
}

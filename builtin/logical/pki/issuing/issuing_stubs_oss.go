// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package issuing

import (
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/certutil/x509verify"
	"github.com/hashicorp/vault/sdk/logical"
)

func entSetCertVerifyOptions(issuer *IssuerEntry, view logical.SystemView, options *x509verify.VerifyOptions) (bool, error) {
	return false, nil
}

func EntAdjustCreationBundle(view logical.SystemView, bundle *certutil.CreationBundle) {
	return
}

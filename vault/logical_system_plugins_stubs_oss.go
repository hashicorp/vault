// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"github.com/hashicorp/vault/sdk/logical"
)

func validateSHA256(sha256 string) *logical.Response {
	if sha256 == "" {
		return logical.ErrorResponse("missing SHA-256 value")
	}
	return nil
}

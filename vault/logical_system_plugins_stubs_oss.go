// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"github.com/hashicorp/vault/sdk/logical"
)

func validateSha256IsEmptyForEntPluginVersion(pluginVersion string, sha256 string) *logical.Response {
	return nil
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// go:build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

func entExtractVerifyPlugin(context.Context, *pluginutil.PluginRunner) error {
	// Do nothing in OSS
	return nil
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package plugincatalog

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

func (c *PluginCatalog) entPrepareDownloadedPlugin(ctx context.Context, plugin pluginutil.SetPluginInput) (string, string, error) {
	return "", "", nil
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package plugincatalog

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

func (c *PluginCatalog) entPrepareDownloadedPlugin(ctx context.Context, plugin pluginutil.SetPluginInput) (string, string, consts.PluginTier, error) {
	return "", "", consts.PluginTierUnknown, nil
}

func (c *PluginCatalog) entDownloadExtractVerifyPlugin(ctx context.Context, pr *pluginutil.PluginRunner) error {
	return nil
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// DownloadExtractVerifyPlugin returns an error as this is an enterprise only feature
func (d dynamicSystemView) DownloadExtractVerifyPlugin(_ context.Context, _ *pluginutil.PluginRunner) error {
	return fmt.Errorf("enterprise only feature")
}

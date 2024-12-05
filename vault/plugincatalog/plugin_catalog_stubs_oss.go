// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package plugincatalog

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

func (c *PluginCatalog) unpackPluginArtifact(context.Context, pluginutil.SetPluginInput) (bool, string, error) {
	return false, "", fmt.Errorf("enterprise-only feature: plugin artifact unpacking")
}

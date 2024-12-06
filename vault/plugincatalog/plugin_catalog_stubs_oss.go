// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package plugincatalog

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

func (c *PluginCatalog) unpackPluginArtifact(plugin pluginutil.SetPluginInput) (bool, string, error) {
	return false, plugin.Command, fmt.Errorf("enterprise-only feature: plugin artifact unpacking")
}

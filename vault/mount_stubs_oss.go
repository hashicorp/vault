// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

type EntMountConfig struct{}

type EntAPIMountConfig struct{}

func entExtractVerifyPlugin(context.Context, *pluginutil.PluginRunner) error {
	// Do nothing in OSS
	return nil
}

// resolveMountEntryVersion allows entry.Version to be overridden if there is a
// corresponding pinned version.
func (c *Core) resolveMountEntryVersion(ctx context.Context, pluginType consts.PluginType, entry *MountEntry) (string, error) {
	pluginName := entry.Type
	if alias, ok := mountAliases[pluginName]; ok {
		pluginName = alias
	}
	pinnedVersion, err := c.pluginCatalog.GetPinnedVersion(ctx, pluginType, pluginName)
	if err != nil && !errors.Is(err, pluginutil.ErrPinnedVersionNotFound) {
		return "", err
	}
	if pinnedVersion != nil {
		return pinnedVersion.Version, nil
	}
	return entry.Version, nil
}

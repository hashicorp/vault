// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugincatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pinnedVersionStoragePrefix = "pinned"
)

func pinnedVersionStorageKey(pluginType consts.PluginType, pluginName string) string {
	return path.Join(pinnedVersionStoragePrefix, pluginType.String(), pluginName)
}

// SetPinnedVersion creates a pinned version for the given plugin name and type.
func (c *PluginCatalog) SetPinnedVersion(ctx context.Context, pin *pluginutil.PinnedVersion) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	plugin, err := c.get(ctx, pin.Name, pin.Type, pin.Version)
	if err != nil {
		return err
	}
	if plugin == nil {
		return fmt.Errorf("%w; %s plugin %q version %s does not exist", ErrPluginNotFound, pin.Type.String(), pin.Name, pin.Version)
	}

	bytes, err := json.Marshal(pin)
	if err != nil {
		return fmt.Errorf("failed to encode pinned version entry: %w", err)
	}

	logicalEntry := logical.StorageEntry{
		Key:   path.Join(pinnedVersionStoragePrefix, pin.Type.String(), pin.Name),
		Value: bytes,
	}

	if err := c.catalogView.Put(ctx, &logicalEntry); err != nil {
		return fmt.Errorf("failed to persist pinned version entry: %w", err)
	}

	return nil
}

// GetPinnedVersion returns the pinned version for the given plugin name and type.
func (c *PluginCatalog) GetPinnedVersion(ctx context.Context, pluginType consts.PluginType, pluginName string) (*pluginutil.PinnedVersion, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.getPinnedVersionInternal(ctx, pinnedVersionStorageKey(pluginType, pluginName))
}

func (c *PluginCatalog) getPinnedVersionInternal(ctx context.Context, key string) (*pluginutil.PinnedVersion, error) {
	logicalEntry, err := c.catalogView.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve pinned version entry: %w", err)
	}

	if logicalEntry == nil {
		return nil, pluginutil.ErrPinnedVersionNotFound
	}

	var pin pluginutil.PinnedVersion
	if err := json.Unmarshal(logicalEntry.Value, &pin); err != nil {
		return nil, fmt.Errorf("failed to decode pinned version entry: %w", err)
	}

	return &pin, nil
}

// DeletePinnedVersion deletes the pinned version for the given plugin name and type.
func (c *PluginCatalog) DeletePinnedVersion(ctx context.Context, pluginType consts.PluginType, pluginName string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if err := c.catalogView.Delete(ctx, path.Join(pinnedVersionStoragePrefix, pluginType.String(), pluginName)); err != nil {
		return fmt.Errorf("failed to delete pinned version entry: %w", err)
	}

	return nil
}

// ListPinnedVersions returns a list of pinned versions for the given plugin type.
func (c *PluginCatalog) ListPinnedVersions(ctx context.Context) ([]*pluginutil.PinnedVersion, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keys, err := logical.CollectKeys(ctx, c.catalogView)
	if err != nil {
		return nil, err
	}

	var pinnedVersions []*pluginutil.PinnedVersion
	for _, key := range keys {
		// Skip: plugin entry.
		if !strings.HasPrefix(key, pinnedVersionStoragePrefix) {
			continue
		}

		pin, err := c.getPinnedVersionInternal(ctx, key)
		if err != nil {
			return nil, err
		}

		pinnedVersions = append(pinnedVersions, pin)
	}

	return pinnedVersions, nil
}

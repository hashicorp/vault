// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package plugincatalog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// setInternal creates a new plugin entry in the catalog and persists it to storage
func (c *PluginCatalog) setInternal(ctx context.Context, plugin pluginutil.SetPluginInput) (*pluginutil.PluginRunner, error) {
	command := plugin.Command

	if plugin.OCIImage == "" {
		// Best effort check to make sure the command isn't breaking out of the
		// configured plugin directory.
		command = filepath.Join(c.directory, plugin.Command)
		sym, err := filepath.EvalSymlinks(command)
		if err != nil {
			return nil, fmt.Errorf("error while validating the command path: %w", err)
		}
		symAbs, err := filepath.Abs(filepath.Dir(sym))
		if err != nil {
			return nil, fmt.Errorf("error while validating the command path: %w", err)
		}

		if symAbs != c.directory {
			return nil, errors.New("cannot execute files outside of configured plugin directory")
		}
	}

	// entryTmp should only be used for the below type and version checks. It uses the
	// full command instead of the relative command because get() normally prepends
	// the plugin directory to the command, but we can't use get() here.
	entryTmp := &pluginutil.PluginRunner{
		Name:     plugin.Name,
		Command:  command,
		OCIImage: plugin.OCIImage,
		Runtime:  plugin.Runtime,
		Args:     plugin.Args,
		Env:      plugin.Env,
		Sha256:   plugin.Sha256,
		Builtin:  false,
	}
	if entryTmp.OCIImage != "" && entryTmp.Runtime != "" {
		var err error
		entryTmp.RuntimeConfig, err = c.runtimeCatalog.Get(ctx, entryTmp.Runtime, consts.PluginRuntimeTypeContainer)
		if err != nil {
			return nil, fmt.Errorf("failed to get configured runtime for plugin %q: %w", plugin.Name, err)
		}
	}

	entry := &pluginutil.PluginRunner{
		Name:     plugin.Name,
		Type:     plugin.Type,
		Version:  plugin.Version,
		Command:  plugin.Command,
		OCIImage: plugin.OCIImage,
		Runtime:  plugin.Runtime,
		Args:     plugin.Args,
		Env:      plugin.Env,
		Sha256:   plugin.Sha256,
		Builtin:  false,
	}

	buf, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to encode plugin entry: %w", err)
	}
	storageKey := path.Join(plugin.Type.String(), plugin.Name)
	if plugin.Version != "" {
		storageKey = path.Join(storageKey, plugin.Version)
	}
	logicalEntry := logical.StorageEntry{
		Key:   storageKey,
		Value: buf,
	}
	if err := c.catalogView.Put(ctx, &logicalEntry); err != nil {
		return nil, fmt.Errorf("failed to persist plugin entry: %w", err)
	}
	return entry, nil
}

func (c *PluginCatalog) entSanitize(context.Context) error {
	return nil
}

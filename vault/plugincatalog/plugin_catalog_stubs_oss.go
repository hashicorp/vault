// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package plugincatalog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"

	semver "github.com/hashicorp/go-version"
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
	// If the plugin type is unknown, we want to attempt to determine the type
	if plugin.Type == consts.PluginTypeUnknown {
		var err error
		plugin.Type, err = c.getPluginTypeFromUnknown(ctx, entryTmp)
		if err != nil {
			return nil, err
		}
		if plugin.Type == consts.PluginTypeUnknown {
			return nil, ErrPluginBadType
		}
	}

	// getting the plugin version is best-effort, so errors are not fatal
	runningVersion := logical.EmptyPluginVersion
	var versionErr error
	switch plugin.Type {
	case consts.PluginTypeSecrets, consts.PluginTypeCredential:
		runningVersion, versionErr = c.getBackendRunningVersion(ctx, entryTmp)
	case consts.PluginTypeDatabase:
		runningVersion, versionErr = c.getDatabaseRunningVersion(ctx, entryTmp)
	default:
		return nil, fmt.Errorf("unknown plugin type: %v", plugin.Type)
	}
	if versionErr != nil {
		c.logger.Warn("Error determining plugin version", "error", versionErr)
		if errors.Is(versionErr, ErrPluginUnableToRun) {
			return nil, versionErr
		}
	} else if plugin.Version != "" && runningVersion.Version != "" && plugin.Version != runningVersion.Version {
		c.logger.Error("Plugin self-reported version did not match requested version",
			"plugin", plugin.Name, "requestedVersion", plugin.Version, "reportedVersion", runningVersion.Version)
		return nil, fmt.Errorf("%w: %s reported version (%s) did not match requested version (%s)",
			ErrPluginVersionMismatch, plugin.Name, runningVersion.Version, plugin.Version)
	} else if plugin.Version == "" && runningVersion.Version != "" {
		plugin.Version = runningVersion.Version
		_, err := semver.NewVersion(plugin.Version)
		if err != nil {
			return nil, fmt.Errorf("plugin self-reported version %q is not a valid semantic version: %w", plugin.Version, err)
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

func (c *PluginCatalog) entValidate(context.Context) error {
	return nil
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugin

import (
	"context"
	"net/rpc"
	"sync"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin"
	bplugin "github.com/hashicorp/vault/sdk/plugin"
)

// Backend returns an instance of the backend, either as a plugin if external
// or as a concrete implementation if builtin, casted as logical.Backend.
func Backend(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	var b backend
	name := conf.Config["plugin_name"]
	pluginType, err := consts.ParsePluginType(conf.Config["plugin_type"])
	if err != nil {
		return nil, err
	}
	pluginVersion := conf.Config["plugin_version"]

	sys := conf.System

	raw, err := plugin.NewBackendV5(ctx, name, pluginType, pluginVersion, sys, conf)
	if err != nil {
		return nil, err
	}
	b.Backend = raw
	b.config = conf

	return &b, nil
}

// backend is a thin wrapper around a builtin plugin or a plugin.BackendPluginClientV5
type backend struct {
	logical.Backend
	mu sync.RWMutex

	config *logical.BackendConfig

	// Used to detect if we already reloaded
	canary string
}

func (b *backend) reloadBackend(ctx context.Context, storage logical.Storage) error {
	pluginName := b.config.Config["plugin_name"]
	pluginType, err := consts.ParsePluginType(b.config.Config["plugin_type"])
	if err != nil {
		return err
	}
	pluginVersion := b.config.Config["plugin_version"]

	b.Logger().Debug("plugin: reloading plugin backend", "plugin", pluginName)

	// Ensure proper cleanup of the backend
	// Pass a context value so that the plugin client will call the appropriate
	// cleanup method for reloading
	reloadCtx := context.WithValue(ctx, plugin.ContextKeyPluginReload, "reload")
	b.Backend.Cleanup(reloadCtx)

	nb, err := plugin.NewBackendV5(ctx, pluginName, pluginType, pluginVersion, b.config.System, b.config)
	if err != nil {
		return err
	}
	err = nb.Setup(ctx, b.config)
	if err != nil {
		return err
	}
	b.Backend = nb

	// Re-initialize the backend in case plugin was reloaded
	// after it crashed
	err = b.Backend.Initialize(ctx, &logical.InitializationRequest{
		Storage: storage,
	})
	if err != nil {
		return err
	}

	return nil
}

// HandleRequest is a thin wrapper implementation of HandleRequest that includes automatic plugin reload.
func (b *backend) HandleRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	b.mu.RLock()
	canary := b.canary
	resp, err := b.Backend.HandleRequest(ctx, req)
	b.mu.RUnlock()
	// Need to compare string value for case were err comes from plugin RPC
	// and is returned as plugin.BasicError type.
	if err != nil &&
		(err.Error() == rpc.ErrShutdown.Error() || err == bplugin.ErrPluginShutdown) {
		// Reload plugin if it's an rpc.ErrShutdown
		b.mu.Lock()
		if b.canary == canary {
			err := b.reloadBackend(ctx, req.Storage)
			if err != nil {
				b.mu.Unlock()
				return nil, err
			}
			b.canary, err = uuid.GenerateUUID()
			if err != nil {
				b.mu.Unlock()
				return nil, err
			}
		}
		b.mu.Unlock()

		// Try request once more
		b.mu.RLock()
		defer b.mu.RUnlock()
		return b.Backend.HandleRequest(ctx, req)
	}
	return resp, err
}

// HandleExistenceCheck is a thin wrapper implementation of HandleRequest that includes automatic plugin reload.
func (b *backend) HandleExistenceCheck(ctx context.Context, req *logical.Request) (bool, bool, error) {
	b.mu.RLock()
	canary := b.canary
	checkFound, exists, err := b.Backend.HandleExistenceCheck(ctx, req)
	b.mu.RUnlock()
	if err != nil &&
		(err.Error() == rpc.ErrShutdown.Error() || err == bplugin.ErrPluginShutdown) {
		// Reload plugin if it's an rpc.ErrShutdown
		b.mu.Lock()
		if b.canary == canary {
			err := b.reloadBackend(ctx, req.Storage)
			if err != nil {
				b.mu.Unlock()
				return false, false, err
			}
			b.canary, err = uuid.GenerateUUID()
			if err != nil {
				b.mu.Unlock()
				return false, false, err
			}
		}
		b.mu.Unlock()

		// Try request once more
		b.mu.RLock()
		defer b.mu.RUnlock()
		return b.Backend.HandleExistenceCheck(ctx, req)
	}
	return checkFound, exists, err
}

// InvalidateKey is a thin wrapper used to ensure we grab the lock for race purposes
func (b *backend) InvalidateKey(ctx context.Context, key string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	b.Backend.InvalidateKey(ctx, key)
}

func (b *backend) IsExternal() bool {
	switch b.Backend.(type) {
	case *plugin.BackendPluginClientV5:
		return true
	}
	return false
}

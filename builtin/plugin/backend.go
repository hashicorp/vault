// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package plugin

import (
	"context"
	"fmt"
	"net/rpc"
	"reflect"
	"sync"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	v5 "github.com/hashicorp/vault/builtin/plugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	bplugin "github.com/hashicorp/vault/sdk/plugin"
)

var (
	ErrMismatchType  = fmt.Errorf("mismatch on mounted backend and plugin backend type")
	ErrMismatchPaths = fmt.Errorf("mismatch on mounted backend and plugin backend special paths")
)

// Factory returns a configured plugin logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	merr := &multierror.Error{}
	_, ok := conf.Config["plugin_name"]
	if !ok {
		return nil, fmt.Errorf("plugin_name not provided")
	}
	b, err := v5.Backend(ctx, conf)
	if err == nil {
		if err := b.Setup(ctx, conf); err != nil {
			return nil, err
		}
		return b, nil
	}
	merr = multierror.Append(merr, err)

	b, err = Backend(ctx, conf)
	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, fmt.Errorf("invalid backend version: %s", merr)
	}

	if err := b.Setup(ctx, conf); err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr.ErrorOrNil()
	}
	return b, nil
}

// Backend returns an instance of the backend, either as a plugin if external
// or as a concrete implementation if builtin, casted as logical.Backend.
func Backend(ctx context.Context, conf *logical.BackendConfig) (*PluginBackend, error) {
	var b PluginBackend

	name := conf.Config["plugin_name"]
	pluginType, err := consts.ParsePluginType(conf.Config["plugin_type"])
	if err != nil {
		return nil, err
	}
	version := conf.Config["plugin_version"]

	sys := conf.System

	// NewBackendWithVersion with isMetadataMode set to true
	raw, err := bplugin.NewBackendWithVersion(ctx, name, pluginType, sys, conf, true, version)
	if err != nil {
		return nil, err
	}
	err = raw.Setup(ctx, conf)
	if err != nil {
		raw.Cleanup(ctx)
		return nil, err
	}
	// Get SpecialPaths and BackendType
	paths := raw.SpecialPaths()
	btype := raw.Type()
	runningVersion := ""
	if versioner, ok := raw.(logical.PluginVersioner); ok {
		runningVersion = versioner.PluginVersion().Version
	}

	// Cleanup meta plugin backend
	raw.Cleanup(ctx)

	// Initialize b.Backend with placeholder backend since plugin
	// backends will need to be lazy loaded.
	b.Backend = &framework.Backend{
		PathsSpecial:   paths,
		BackendType:    btype,
		RunningVersion: runningVersion,
	}

	b.config = conf

	return &b, nil
}

// PluginBackend is a thin wrapper around plugin.BackendPluginClient
type PluginBackend struct {
	Backend logical.Backend
	sync.RWMutex

	config *logical.BackendConfig

	// Used to detect if we already reloaded
	canary string

	// Used to detect if plugin is set
	loaded bool
}

// startBackend starts a plugin backend
func (b *PluginBackend) startBackend(ctx context.Context, storage logical.Storage) error {
	pluginName := b.config.Config["plugin_name"]
	pluginType, err := consts.ParsePluginType(b.config.Config["plugin_type"])
	if err != nil {
		return err
	}

	// Ensure proper cleanup of the backend (i.e. call client.Kill())
	b.Backend.Cleanup(ctx)

	nb, err := bplugin.NewBackendWithVersion(ctx, pluginName, pluginType, b.config.System, b.config, false, b.config.Config["plugin_version"])
	if err != nil {
		return err
	}
	err = nb.Setup(ctx, b.config)
	if err != nil {
		nb.Cleanup(ctx)
		return err
	}

	// If the backend has not been loaded (i.e. still in metadata mode),
	// check if type and special paths still matches
	if !b.loaded {
		if b.Backend.Type() != nb.Type() {
			nb.Cleanup(ctx)
			b.Backend.Logger().Warn("failed to start plugin process", "plugin", pluginName, "error", ErrMismatchType)
			return ErrMismatchType
		}
		if !reflect.DeepEqual(b.Backend.SpecialPaths(), nb.SpecialPaths()) {
			nb.Cleanup(ctx)
			b.Backend.Logger().Warn("failed to start plugin process", "plugin", pluginName, "error", ErrMismatchPaths)
			return ErrMismatchPaths
		}
	}

	b.Backend = nb
	b.loaded = true

	// call Initialize() explicitly here.
	return b.Backend.Initialize(ctx, &logical.InitializationRequest{
		Storage: storage,
	})
}

// lazyLoad lazy-loads the backend before running a method
func (b *PluginBackend) lazyLoadBackend(ctx context.Context, storage logical.Storage, methodWrapper func() error) error {
	b.RLock()
	canary := b.canary

	// Lazy-load backend
	if !b.loaded {
		// Upgrade lock
		b.RUnlock()
		b.Lock()
		// Check once more after lock swap
		if !b.loaded {
			err := b.startBackend(ctx, storage)
			if err != nil {
				b.Unlock()
				return err
			}
		}
		b.Unlock()
		b.RLock()
	}

	err := methodWrapper()
	b.RUnlock()

	// Need to compare string value for case were err comes from plugin RPC
	// and is returned as plugin.BasicError type.
	if err != nil &&
		(err.Error() == rpc.ErrShutdown.Error() || err == bplugin.ErrPluginShutdown) {
		// Reload plugin if it's an rpc.ErrShutdown
		b.Lock()
		if b.canary == canary {
			b.Backend.Logger().Debug("reloading plugin backend", "plugin", b.config.Config["plugin_name"])
			err := b.startBackend(ctx, storage)
			if err != nil {
				b.Unlock()
				return err
			}
			b.canary, err = uuid.GenerateUUID()
			if err != nil {
				b.Unlock()
				return err
			}
		}
		b.Unlock()

		// Try once more
		b.RLock()
		defer b.RUnlock()
		return methodWrapper()
	}
	return err
}

// HandleRequest is a thin wrapper implementation of HandleRequest that includes
// automatic plugin reload.
func (b *PluginBackend) HandleRequest(ctx context.Context, req *logical.Request) (resp *logical.Response, err error) {
	err = b.lazyLoadBackend(ctx, req.Storage, func() error {
		var merr error
		resp, merr = b.Backend.HandleRequest(ctx, req)
		return merr
	})

	return
}

// HandleExistenceCheck is a thin wrapper implementation of HandleExistenceCheck
// that includes automatic plugin reload.
func (b *PluginBackend) HandleExistenceCheck(ctx context.Context, req *logical.Request) (checkFound bool, exists bool, err error) {
	err = b.lazyLoadBackend(ctx, req.Storage, func() error {
		var merr error
		checkFound, exists, merr = b.Backend.HandleExistenceCheck(ctx, req)
		return merr
	})

	return
}

// Initialize is intentionally a no-op here, the backend will instead be
// initialized when it is lazily loaded.
func (b *PluginBackend) Initialize(ctx context.Context, req *logical.InitializationRequest) error {
	return nil
}

// SpecialPaths is a thin wrapper used to ensure we grab the lock for race purposes
func (b *PluginBackend) SpecialPaths() *logical.Paths {
	b.RLock()
	defer b.RUnlock()
	return b.Backend.SpecialPaths()
}

// System is a thin wrapper used to ensure we grab the lock for race purposes
func (b *PluginBackend) System() logical.SystemView {
	b.RLock()
	defer b.RUnlock()
	return b.Backend.System()
}

// Logger is a thin wrapper used to ensure we grab the lock for race purposes
func (b *PluginBackend) Logger() log.Logger {
	b.RLock()
	defer b.RUnlock()
	return b.Backend.Logger()
}

// Cleanup is a thin wrapper used to ensure we grab the lock for race purposes
func (b *PluginBackend) Cleanup(ctx context.Context) {
	b.RLock()
	defer b.RUnlock()
	b.Backend.Cleanup(ctx)
}

// InvalidateKey is a thin wrapper used to ensure we grab the lock for race purposes
func (b *PluginBackend) InvalidateKey(ctx context.Context, key string) {
	b.RLock()
	defer b.RUnlock()
	b.Backend.InvalidateKey(ctx, key)
}

// Setup is a thin wrapper used to ensure we grab the lock for race purposes
func (b *PluginBackend) Setup(ctx context.Context, config *logical.BackendConfig) error {
	b.RLock()
	defer b.RUnlock()
	return b.Backend.Setup(ctx, config)
}

// Type is a thin wrapper used to ensure we grab the lock for race purposes
func (b *PluginBackend) Type() logical.BackendType {
	b.RLock()
	defer b.RUnlock()
	return b.Backend.Type()
}

func (b *PluginBackend) PluginVersion() logical.PluginVersion {
	if versioner, ok := b.Backend.(logical.PluginVersioner); ok {
		return versioner.PluginVersion()
	}
	return logical.EmptyPluginVersion
}

var _ logical.PluginVersioner = (*PluginBackend)(nil)

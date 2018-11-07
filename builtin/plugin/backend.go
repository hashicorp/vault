package plugin

import (
	"context"
	"fmt"
	"net/rpc"
	"reflect"
	"sync"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	bplugin "github.com/hashicorp/vault/logical/plugin"
)

var (
	ErrMismatchType  = fmt.Errorf("mismatch on mounted backend and plugin backend type")
	ErrMismatchPaths = fmt.Errorf("mismatch on mounted backend and plugin backend special paths")
)

// Factory returns a configured plugin logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	_, ok := conf.Config["plugin_name"]
	if !ok {
		return nil, fmt.Errorf("plugin_name not provided")
	}
	b, err := Backend(ctx, conf)
	if err != nil {
		return nil, err
	}

	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend returns an instance of the backend, either as a plugin if external
// or as a concrete implementation if builtin, casted as logical.Backend.
func Backend(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	var b PluginBackend

	name := conf.Config["plugin_name"]
	pluginType, err := consts.ParsePluginType(conf.Config["plugin_type"])
	if err != nil {
		return nil, err
	}

	sys := conf.System

	// NewBackend with isMetadataMode set to true
	raw, err := bplugin.NewBackend(ctx, name, pluginType, sys, conf, true)
	if err != nil {
		return nil, err
	}
	err = raw.Setup(ctx, conf)
	if err != nil {
		return nil, err
	}
	// Get SpecialPaths and BackendType
	paths := raw.SpecialPaths()
	btype := raw.Type()

	// Cleanup meta plugin backend
	raw.Cleanup(ctx)

	// Initialize b.Backend with dummy backend since plugin
	// backends will need to be lazy loaded.
	b.Backend = &framework.Backend{
		PathsSpecial: paths,
		BackendType:  btype,
	}

	b.config = conf

	return &b, nil
}

// PluginBackend is a thin wrapper around plugin.BackendPluginClient
type PluginBackend struct {
	logical.Backend
	sync.RWMutex

	config *logical.BackendConfig

	// Used to detect if we already reloaded
	canary string

	// Used to detect if plugin is set
	loaded bool
}

func (b *PluginBackend) reloadBackend(ctx context.Context) error {
	b.Logger().Debug("reloading plugin backend", "plugin", b.config.Config["plugin_name"])
	return b.startBackend(ctx)
}

// startBackend starts a plugin backend
func (b *PluginBackend) startBackend(ctx context.Context) error {
	pluginName := b.config.Config["plugin_name"]
	pluginType, err := consts.ParsePluginType(b.config.Config["plugin_type"])
	if err != nil {
		return err
	}

	// Ensure proper cleanup of the backend (i.e. call client.Kill())
	b.Backend.Cleanup(ctx)

	nb, err := bplugin.NewBackend(ctx, pluginName, pluginType, b.config.System, b.config, false)
	if err != nil {
		return err
	}
	err = nb.Setup(ctx, b.config)
	if err != nil {
		return err
	}

	// If the backend has not been loaded (i.e. still in metadata mode),
	// check if type and special paths still matches
	if !b.loaded {
		if b.Backend.Type() != nb.Type() {
			nb.Cleanup(ctx)
			b.Logger().Warn("failed to start plugin process", "plugin", b.config.Config["plugin_name"], "error", ErrMismatchType)
			return ErrMismatchType
		}
		if !reflect.DeepEqual(b.Backend.SpecialPaths(), nb.SpecialPaths()) {
			nb.Cleanup(ctx)
			b.Logger().Warn("failed to start plugin process", "plugin", b.config.Config["plugin_name"], "error", ErrMismatchPaths)
			return ErrMismatchPaths
		}
	}

	b.Backend = nb
	b.loaded = true

	return nil
}

// HandleRequest is a thin wrapper implementation of HandleRequest that includes automatic plugin reload.
func (b *PluginBackend) HandleRequest(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	b.RLock()
	canary := b.canary

	// Lazy-load backend
	if !b.loaded {
		// Upgrade lock
		b.RUnlock()
		b.Lock()
		// Check once more after lock swap
		if !b.loaded {
			err := b.startBackend(ctx)
			if err != nil {
				b.Unlock()
				return nil, err
			}
		}
		b.Unlock()
		b.RLock()
	}
	resp, err := b.Backend.HandleRequest(ctx, req)
	b.RUnlock()
	// Need to compare string value for case were err comes from plugin RPC
	// and is returned as plugin.BasicError type.
	if err != nil &&
		(err.Error() == rpc.ErrShutdown.Error() || err == bplugin.ErrPluginShutdown) {
		// Reload plugin if it's an rpc.ErrShutdown
		b.Lock()
		if b.canary == canary {
			err := b.reloadBackend(ctx)
			if err != nil {
				b.Unlock()
				return nil, err
			}
			b.canary, err = uuid.GenerateUUID()
			if err != nil {
				b.Unlock()
				return nil, err
			}
		}
		b.Unlock()

		// Try request once more
		b.RLock()
		defer b.RUnlock()
		return b.Backend.HandleRequest(ctx, req)
	}
	return resp, err
}

// HandleExistenceCheck is a thin wrapper implementation of HandleRequest that includes automatic plugin reload.
func (b *PluginBackend) HandleExistenceCheck(ctx context.Context, req *logical.Request) (bool, bool, error) {
	b.RLock()
	canary := b.canary

	// Lazy-load backend
	if !b.loaded {
		// Upgrade lock
		b.RUnlock()
		b.Lock()
		// Check once more after lock swap
		if !b.loaded {
			err := b.startBackend(ctx)
			if err != nil {
				b.Unlock()
				return false, false, err
			}
		}
		b.Unlock()
		b.RLock()
	}

	checkFound, exists, err := b.Backend.HandleExistenceCheck(ctx, req)
	b.RUnlock()
	if err != nil &&
		(err.Error() == rpc.ErrShutdown.Error() || err == bplugin.ErrPluginShutdown) {
		// Reload plugin if it's an rpc.ErrShutdown
		b.Lock()
		if b.canary == canary {
			err := b.reloadBackend(ctx)
			if err != nil {
				b.Unlock()
				return false, false, err
			}
			b.canary, err = uuid.GenerateUUID()
			if err != nil {
				b.Unlock()
				return false, false, err
			}
		}
		b.Unlock()

		// Try request once more
		b.RLock()
		defer b.RUnlock()
		return b.Backend.HandleExistenceCheck(ctx, req)
	}
	return checkFound, exists, err
}

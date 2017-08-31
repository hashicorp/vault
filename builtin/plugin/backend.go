package plugin

import (
	"fmt"
	"net/rpc"
	"reflect"
	"sync"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	bplugin "github.com/hashicorp/vault/logical/plugin"
)

var (
	ErrMismatchType  = fmt.Errorf("mismatch on mounted backend and plugin backend type")
	ErrMismatchPaths = fmt.Errorf("mismatch on mounted backend and plugin backend special paths")
)

// Factory returns a configured plugin logical.Backend.
func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	_, ok := conf.Config["plugin_name"]
	if !ok {
		return nil, fmt.Errorf("plugin_name not provided")
	}
	b, err := Backend(conf)
	if err != nil {
		return nil, err
	}

	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

// Backend returns an instance of the backend, either as a plugin if external
// or as a concrete implementation if builtin, casted as logical.Backend.
func Backend(conf *logical.BackendConfig) (logical.Backend, error) {
	var b backend

	name := conf.Config["plugin_name"]
	sys := conf.System

	// NewBackend with isMetadataMode set to true
	raw, err := bplugin.NewBackend(name, sys, conf.Logger, true)
	if err != nil {
		return nil, err
	}
	err = raw.Setup(conf)
	if err != nil {
		return nil, err
	}
	// Get SpecialPaths and BackendType
	paths := raw.SpecialPaths()
	btype := raw.Type()

	// Cleanup meta plugin backend
	raw.Cleanup()

	// Initialize b.Backend with dummy backend since plugin
	// backends will need to be lazy loaded.
	b.Backend = &framework.Backend{
		PathsSpecial: paths,
		BackendType:  btype,
	}

	b.config = conf

	return &b, nil
}

// backend is a thin wrapper around plugin.BackendPluginClient
type backend struct {
	logical.Backend
	sync.RWMutex

	config *logical.BackendConfig

	// Used to detect if we already reloaded
	canary string

	// Used to detect if plugin is set
	loaded bool
}

func (b *backend) reloadBackend() error {
	b.Logger().Trace("plugin: reloading plugin backend", "plugin", b.config.Config["plugin_name"])
	return b.startBackend()
}

// startBackend starts a plugin backend
func (b *backend) startBackend() error {
	pluginName := b.config.Config["plugin_name"]

	// Ensure proper cleanup of the backend (i.e. call client.Kill())
	b.Backend.Cleanup()

	nb, err := bplugin.NewBackend(pluginName, b.config.System, b.config.Logger, false)
	if err != nil {
		return err
	}
	err = nb.Setup(b.config)
	if err != nil {
		return err
	}

	// If the backend has not been loaded (i.e. still in metadata mode),
	// check if type and special paths still matches
	if !b.loaded {
		if b.Backend.Type() != nb.Type() {
			nb.Cleanup()
			b.Logger().Warn("plugin: failed to start plugin process", "plugin", b.config.Config["plugin_name"], "error", ErrMismatchType)
			return ErrMismatchType
		}
		if !reflect.DeepEqual(b.Backend.SpecialPaths(), nb.SpecialPaths()) {
			nb.Cleanup()
			b.Logger().Warn("plugin: failed to start plugin process", "plugin", b.config.Config["plugin_name"], "error", ErrMismatchPaths)
			return ErrMismatchPaths
		}
	}

	b.Backend = nb
	b.loaded = true

	// Call initialize
	if err := b.Backend.Initialize(); err != nil {
		return err
	}

	return nil
}

// HandleRequest is a thin wrapper implementation of HandleRequest that includes automatic plugin reload.
func (b *backend) HandleRequest(req *logical.Request) (*logical.Response, error) {
	b.RLock()
	canary := b.canary

	// Lazy-load backend
	if !b.loaded {
		// Upgrade lock
		b.RUnlock()
		b.Lock()
		// Check once more after lock swap
		if !b.loaded {
			err := b.startBackend()
			if err != nil {
				b.Unlock()
				return nil, err
			}
		}
		b.Unlock()
		b.RLock()
	}
	resp, err := b.Backend.HandleRequest(req)
	b.RUnlock()
	// Need to compare string value for case were err comes from plugin RPC
	// and is returned as plugin.BasicError type.
	if err != nil && err.Error() == rpc.ErrShutdown.Error() {
		// Reload plugin if it's an rpc.ErrShutdown
		b.Lock()
		if b.canary == canary {
			err := b.reloadBackend()
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
		return b.Backend.HandleRequest(req)
	}
	return resp, err
}

// HandleExistenceCheck is a thin wrapper implementation of HandleRequest that includes automatic plugin reload.
func (b *backend) HandleExistenceCheck(req *logical.Request) (bool, bool, error) {
	b.RLock()
	canary := b.canary

	// Lazy-load backend
	if !b.loaded {
		// Upgrade lock
		b.RUnlock()
		b.Lock()
		// Check once more after lock swap
		if !b.loaded {
			err := b.startBackend()
			if err != nil {
				b.Unlock()
				return false, false, err
			}
		}
		b.Unlock()
		b.RLock()
	}

	checkFound, exists, err := b.Backend.HandleExistenceCheck(req)
	b.RUnlock()
	if err != nil && err.Error() == rpc.ErrShutdown.Error() {
		// Reload plugin if it's an rpc.ErrShutdown
		b.Lock()
		if b.canary == canary {
			err := b.reloadBackend()
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
		return b.Backend.HandleExistenceCheck(req)
	}
	return checkFound, exists, err
}

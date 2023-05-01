// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
)

// BackendPluginClientV5 is a wrapper around backendPluginClient
// that also contains its plugin.Client instance. It's primarily
// used to cleanly kill the client on Cleanup()
type BackendPluginClientV5 struct {
	client pluginutil.PluginClient

	logical.Backend
}

type ContextKey string

func (c ContextKey) String() string {
	return "plugin" + string(c)
}

const ContextKeyPluginReload = ContextKey("plugin-reload")

// Cleanup cleans up the go-plugin client and the plugin catalog
func (b *BackendPluginClientV5) Cleanup(ctx context.Context) {
	_, ok := ctx.Value(ContextKeyPluginReload).(string)
	if !ok {
		b.Backend.Cleanup(ctx)
		b.client.Close()
		return
	}
	b.Backend.Cleanup(ctx)
	b.client.Reload()
}

func (b *BackendPluginClientV5) IsExternal() bool {
	return true
}

func (b *BackendPluginClientV5) PluginVersion() logical.PluginVersion {
	if versioner, ok := b.Backend.(logical.PluginVersioner); ok {
		return versioner.PluginVersion()
	}
	return logical.EmptyPluginVersion
}

var _ logical.PluginVersioner = (*BackendPluginClientV5)(nil)

// NewBackendV5 will return an instance of an RPC-based client implementation of
// the backend for external plugins, or a concrete implementation of the
// backend if it is a builtin backend. The backend is returned as a
// logical.Backend interface.
func NewBackendV5(ctx context.Context, pluginName string, pluginType consts.PluginType, pluginVersion string, sys pluginutil.LookRunnerUtil, conf *logical.BackendConfig) (logical.Backend, error) {
	// Look for plugin in the plugin catalog
	pluginRunner, err := sys.LookupPluginVersion(ctx, pluginName, pluginType, pluginVersion)
	if err != nil {
		return nil, err
	}

	var backend logical.Backend
	if pluginRunner.Builtin {
		// Plugin is builtin so we can retrieve an instance of the interface
		// from the pluginRunner. Then cast it to logical.Factory.
		rawFactory, err := pluginRunner.BuiltinFactory()
		if err != nil {
			return nil, fmt.Errorf("error getting plugin type: %q", err)
		}

		if factory, ok := rawFactory.(logical.Factory); !ok {
			return nil, fmt.Errorf("unsupported backend type: %q", pluginName)
		} else {
			if backend, err = factory(ctx, conf); err != nil {
				return nil, err
			}
		}
	} else {
		// create a backendPluginClient instance
		config := pluginutil.PluginClientConfig{
			Name:            pluginName,
			PluginSets:      PluginSet,
			PluginType:      pluginType,
			Version:         pluginVersion,
			HandshakeConfig: HandshakeConfig,
			Logger:          conf.Logger.Named(pluginName),
			AutoMTLS:        true,
			Wrapper:         sys,
		}
		backend, err = NewPluginClientV5(ctx, sys, config)
		if err != nil {
			return nil, err
		}
	}

	return backend, nil
}

// PluginSet is the map of plugins we can dispense.
var PluginSet = map[int]plugin.PluginSet{
	5: {
		"backend": &GRPCBackendPlugin{},
	},
}

func Dispense(rpcClient plugin.ClientProtocol, pluginClient pluginutil.PluginClient) (logical.Backend, error) {
	// Request the plugin
	raw, err := rpcClient.Dispense("backend")
	if err != nil {
		return nil, err
	}

	var backend logical.Backend
	// We should have a logical backend type now. This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	switch c := raw.(type) {
	case *backendGRPCPluginClient:
		// This is an abstraction leak from go-plugin but it is necessary in
		// order to enable multiplexing on multiplexed plugins
		c.client = pb.NewBackendClient(pluginClient.Conn())
		c.versionClient = logical.NewPluginVersionClient(pluginClient.Conn())

		backend = c
	default:
		return nil, errors.New("unsupported plugin client type")
	}

	return &BackendPluginClientV5{
		client:  pluginClient,
		Backend: backend,
	}, nil
}

func NewPluginClientV5(ctx context.Context, sys pluginutil.RunnerUtil, config pluginutil.PluginClientConfig) (logical.Backend, error) {
	pluginClient, err := sys.NewPluginClient(ctx, config)
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := pluginClient.Dispense("backend")
	if err != nil {
		return nil, err
	}

	var backend logical.Backend
	var transport string
	// We should have a logical backend type now. This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	switch c := raw.(type) {
	case *backendGRPCPluginClient:
		// This is an abstraction leak from go-plugin but it is necessary in
		// order to enable multiplexing on multiplexed plugins
		c.client = pb.NewBackendClient(pluginClient.Conn())
		c.versionClient = logical.NewPluginVersionClient(pluginClient.Conn())

		backend = c
		transport = "gRPC"
	default:
		return nil, errors.New("unsupported plugin client type")
	}

	// Wrap the backend in a tracing middleware
	if config.Logger.IsTrace() {
		backend = &BackendTracingMiddleware{
			logger: config.Logger.With("transport", transport),
			next:   backend,
		}
	}

	return &BackendPluginClientV5{
		client:  pluginClient,
		Backend: backend,
	}, nil
}

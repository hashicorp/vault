// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"errors"
	"fmt"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// BackendPluginClient is a wrapper around backendPluginClient
// that also contains its plugin.Client instance. It's primarily
// used to cleanly kill the client on Cleanup()
type BackendPluginClient struct {
	client *plugin.Client

	logical.Backend
}

// Cleanup calls the RPC client's Cleanup() func and also calls
// the go-plugin's client Kill() func
func (b *BackendPluginClient) Cleanup(ctx context.Context) {
	b.Backend.Cleanup(ctx)
	b.client.Kill()
}

// NewBackendWithVersion will return an instance of an RPC-based client implementation of the backend for
// external plugins, or a concrete implementation of the backend if it is a builtin backend.
// The backend is returned as a logical.Backend interface. The isMetadataMode param determines whether
// the plugin should run in metadata mode.
func NewBackendWithVersion(ctx context.Context, pluginName string, pluginType consts.PluginType, sys pluginutil.LookRunnerUtil, conf *logical.BackendConfig, isMetadataMode bool, version string) (logical.Backend, error) {
	// Look for plugin in the plugin catalog
	pluginRunner, err := sys.LookupPluginVersion(ctx, pluginName, pluginType, version)
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
		backend, err = NewPluginClient(ctx, sys, pluginRunner, conf.Logger, isMetadataMode)
		if err != nil {
			return nil, err
		}
	}

	return backend, nil
}

// NewBackend will return an instance of an RPC-based client implementation of the backend for
// external plugins, or a concrete implementation of the backend if it is a builtin backend.
// The backend is returned as a logical.Backend interface. The isMetadataMode param determines whether
// the plugin should run in metadata mode.
func NewBackend(ctx context.Context, pluginName string, pluginType consts.PluginType, sys pluginutil.LookRunnerUtil, conf *logical.BackendConfig, isMetadataMode bool) (logical.Backend, error) {
	return NewBackendWithVersion(ctx, pluginName, pluginType, sys, conf, isMetadataMode, "")
}

func NewPluginClient(ctx context.Context, sys pluginutil.RunnerUtil, pluginRunner *pluginutil.PluginRunner, logger log.Logger, isMetadataMode bool) (logical.Backend, error) {
	// pluginMap is the map of plugins we can dispense.
	pluginSet := map[int]plugin.PluginSet{
		// Version 3 used to supports both protocols. We want to keep it around
		// since it's possible old plugins built against this version will still
		// work with gRPC. There is currently no difference between version 3
		// and version 4.
		3: {
			"backend": &GRPCBackendPlugin{
				MetadataMode: isMetadataMode,
			},
		},
		4: {
			"backend": &GRPCBackendPlugin{
				MetadataMode: isMetadataMode,
			},
		},
	}

	namedLogger := logger.Named(pluginRunner.Name)

	var client *plugin.Client
	var err error
	if isMetadataMode {
		client, err = pluginRunner.RunMetadataMode(ctx, sys, pluginSet, HandshakeConfig, []string{}, namedLogger)
	} else {
		client, err = pluginRunner.Run(ctx, sys, pluginSet, HandshakeConfig, []string{}, namedLogger)
	}
	if err != nil {
		return nil, err
	}

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("backend")
	if err != nil {
		return nil, err
	}

	var backend logical.Backend
	var transport string
	// We should have a logical backend type now. This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	switch b := raw.(type) {
	case *backendGRPCPluginClient:
		backend = b
		transport = "gRPC"
	default:
		return nil, errors.New("unsupported plugin client type")
	}

	// Wrap the backend in a tracing middleware
	if namedLogger.IsTrace() {
		backend = &BackendTracingMiddleware{
			logger: namedLogger.With("transport", transport),
			next:   backend,
		}
	}

	return &BackendPluginClient{
		client:  client,
		Backend: backend,
	}, nil
}

func (b *BackendPluginClient) PluginVersion() logical.PluginVersion {
	if versioner, ok := b.Backend.(logical.PluginVersioner); ok {
		return versioner.PluginVersion()
	}
	return logical.EmptyPluginVersion
}

var _ logical.PluginVersioner = (*BackendPluginClient)(nil)

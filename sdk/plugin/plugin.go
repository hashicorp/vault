package plugin

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/errwrap"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// BackendPluginClient is a wrapper around backendPluginClient
// that also contains its plugin.Client instance. It's primarily
// used to cleanly kill the client on Cleanup()
type BackendPluginClient struct {
	client *plugin.Client
	sync.Mutex

	logical.Backend
}

// AutoMTLSSupported inspects the go-plugin client's negotiated version. This
// allows us to determine if a given client supports AutoMTLS
func (b *BackendPluginClient) AutoMTLSSupported() bool {
	return b.client.NegotiatedVersion() >= pluginutil.BackendPluginAutoMTLSMinVersion
}

// Cleanup calls the RPC client's Cleanup() func and also calls
// the go-plugin's client Kill() func
func (b *BackendPluginClient) Cleanup(ctx context.Context) {
	b.Backend.Cleanup(ctx)
	b.client.Kill()
}

// NewBackend will return an instance of an RPC-based client implementation of the backend for
// external plugins, or a concrete implementation of the backend if it is a builtin backend.
// The backend is returned as a logical.Backend interface. The isMetadataMode param determines whether
// the plugin should run in metadata mode.
func NewBackend(ctx context.Context, pluginName string, pluginType consts.PluginType, sys pluginutil.LookRunnerUtil, conf *logical.BackendConfig, isMetadataMode bool, autoMTLS bool) (logical.Backend, error) {
	// Look for plugin in the plugin catalog
	pluginRunner, err := sys.LookupPlugin(ctx, pluginName, pluginType)
	if err != nil {
		return nil, err
	}

	var backend logical.Backend
	if pluginRunner.Builtin {
		// Plugin is builtin so we can retrieve an instance of the interface
		// from the pluginRunner. Then cast it to logical.Factory.
		rawFactory, err := pluginRunner.BuiltinFactory()
		if err != nil {
			return nil, errwrap.Wrapf("error getting plugin type: {{err}}", err)
		}

		if factory, ok := rawFactory.(logical.Factory); !ok {
			return nil, fmt.Errorf("unsupported backend type: %q", pluginName)
		} else {
			if backend, err = factory(ctx, conf); err != nil {
				return nil, err
			}
		}
	} else {
		config := pluginutil.PluginClientConfig{
			Name:           pluginName,
			PluginType:     pluginType,
			Logger:         conf.Logger.Named(pluginName),
			IsMetadataMode: isMetadataMode,
			AutoMTLS:       autoMTLS,
			Wrapper:        sys,
		}
		// create a backendPluginClient instance
		backend, err = NewPluginClient(ctx, pluginRunner, config)
		if err != nil {
			return nil, err
		}
	}

	return backend, nil
}

// pluginSet returns the go-plugin PluginSet that we can dispense. This ensures
// that plugins that don't support AutoMTLS are run on the appropriate version.
func pluginSet(autoMTLS, metadataMode bool) map[int]plugin.PluginSet {
	if autoMTLS {
		return map[int]plugin.PluginSet{
			5: {
				"backend": &GRPCBackendPlugin{
					MetadataMode:      false,
					AutoMTLSSupported: true,
				},
			},
		}
	}
	return map[int]plugin.PluginSet{
		// Version 3 used to supports both protocols. We want to keep it around
		// since it's possible old plugins built against this version will still
		// work with gRPC. There is currently no difference between version 3
		// and version 4.
		3: {
			"backend": &GRPCBackendPlugin{
				MetadataMode: metadataMode,
			},
		},
		4: {
			"backend": &GRPCBackendPlugin{
				MetadataMode: metadataMode,
			},
		},
	}
}

func NewPluginClient(ctx context.Context, pluginRunner *pluginutil.PluginRunner, config pluginutil.PluginClientConfig) (logical.Backend, error) {
	var client *plugin.Client
	var err error
	pluginSet := pluginSet(config.AutoMTLS, config.IsMetadataMode)

	if config.AutoMTLS {
		client, err = pluginRunner.RunAutoMTLS(ctx, config.Wrapper, pluginSet, handshakeConfig, []string{}, config.Logger)
	} else if config.IsMetadataMode {
		client, err = pluginRunner.RunMetadataMode(ctx, config.Wrapper, pluginSet, handshakeConfig, []string{}, config.Logger)
	} else {
		client, err = pluginRunner.Run(ctx, config.Wrapper, pluginSet, handshakeConfig, []string{}, config.Logger)
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
	switch raw.(type) {
	case *backendGRPCPluginClient:
		backend = raw.(*backendGRPCPluginClient)
		transport = "gRPC"
	default:
		return nil, errors.New("unsupported plugin client type")
	}

	// Wrap the backend in a tracing middleware
	if config.Logger.IsTrace() {
		backend = &backendTracingMiddleware{
			logger: config.Logger.With("transport", transport),
			next:   backend,
		}
	}

	return &BackendPluginClient{
		client:  client,
		Backend: backend,
	}, nil
}

// wrapError takes a generic error type and makes it usable with the plugin
// interface. Only errors which have exported fields and have been registered
// with gob can be unwrapped and transported. This checks error types and, if
// none match, wrap the error in a plugin.BasicError.
func wrapError(err error) error {
	if err == nil {
		return nil
	}

	switch err.(type) {
	case *plugin.BasicError,
		logical.HTTPCodedError,
		*logical.StatusBadRequest:
		return err
	}

	return plugin.NewBasicError(err)
}

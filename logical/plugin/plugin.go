package plugin

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"sync"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
)

// init registers basic structs with gob which will be used to transport complex
// types through the plugin server and client.
func init() {
	// Common basic structs
	gob.Register([]interface{}{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[string]string{})
	gob.Register(map[string]int{})

	// Register these types since we have to serialize and de-serialize
	// tls.ConnectionState over the wire as part of logical.Request.Connection.
	gob.Register(rsa.PublicKey{})
	gob.Register(ecdsa.PublicKey{})
	gob.Register(time.Duration(0))

	// Custom common error types for requests. If you add something here, you must
	// also add it to the switch statement in `wrapError`!
	gob.Register(&plugin.BasicError{})
	gob.Register(logical.CodedError(0, ""))
	gob.Register(&logical.StatusBadRequest{})
}

// BackendPluginClient is a wrapper around backendPluginClient
// that also contains its plugin.Client instance. It's primarily
// used to cleanly kill the client on Cleanup()
type BackendPluginClient struct {
	client *plugin.Client
	sync.Mutex

	logical.Backend
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
func NewBackend(ctx context.Context, pluginName string, sys pluginutil.LookRunnerUtil, logger log.Logger, isMetadataMode bool) (logical.Backend, error) {
	// Look for plugin in the plugin catalog
	pluginRunner, err := sys.LookupPlugin(ctx, pluginName)
	if err != nil {
		return nil, err
	}

	var backend logical.Backend
	if pluginRunner.Builtin {
		// Plugin is builtin so we can retrieve an instance of the interface
		// from the pluginRunner. Then cast it to logical.Backend.
		backendRaw, err := pluginRunner.BuiltinFactory()
		if err != nil {
			return nil, fmt.Errorf("error getting plugin type: %s", err)
		}

		var ok bool
		backend, ok = backendRaw.(logical.Backend)
		if !ok {
			return nil, fmt.Errorf("unsupported backend type: %s", pluginName)
		}

	} else {
		// create a backendPluginClient instance
		backend, err = newPluginClient(ctx, sys, pluginRunner, logger, isMetadataMode)
		if err != nil {
			return nil, err
		}
	}

	return backend, nil
}

func newPluginClient(ctx context.Context, sys pluginutil.RunnerUtil, pluginRunner *pluginutil.PluginRunner, logger log.Logger, isMetadataMode bool) (logical.Backend, error) {
	// pluginMap is the map of plugins we can dispense.
	pluginMap := map[string]plugin.Plugin{
		"backend": &BackendPlugin{
			metadataMode: isMetadataMode,
		},
	}

	namedLogger := logger.Named(pluginRunner.Name)

	var client *plugin.Client
	var err error
	if isMetadataMode {
		client, err = pluginRunner.RunMetadataMode(ctx, sys, pluginMap, handshakeConfig, []string{}, namedLogger)
	} else {
		client, err = pluginRunner.Run(ctx, sys, pluginMap, handshakeConfig, []string{}, namedLogger)
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
	case *backendPluginClient:
		backend = raw.(*backendPluginClient)
		transport = "netRPC"
	case *backendGRPCPluginClient:
		backend = raw.(*backendGRPCPluginClient)
		transport = "gRPC"
	default:
		return nil, errors.New("Unsupported plugin client type")
	}

	// Wrap the backend in a tracing middleware
	if namedLogger.IsTrace() {
		backend = &backendTracingMiddleware{
			logger: namedLogger.With("transport", transport),
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

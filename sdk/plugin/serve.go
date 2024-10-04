// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"crypto/tls"
	"math"
	"os"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/grpc"
)

// BackendPluginName is the name of the plugin that can be
// dispensed from the plugin server.
const BackendPluginName = "backend"

type TLSProviderFunc func() (*tls.Config, error)

type ServeOpts struct {
	BackendFactoryFunc logical.Factory
	TLSProviderFunc    TLSProviderFunc
	Logger             log.Logger
}

// Serve is a helper function used to serve a backend plugin. This
// should be ran on the plugin's main process.
func Serve(opts *ServeOpts) error {
	logger := opts.Logger
	if logger == nil {
		logger = log.New(&log.LoggerOptions{
			Level:      log.Trace,
			Output:     os.Stderr,
			JSONFormat: true,
		})
	}

	// pluginMap is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		// Version 3 used to supports both protocols. We want to keep it around
		// since it's possible old plugins built against this version will still
		// work with gRPC. There is currently no difference between version 3
		// and version 4.
		3: {
			"backend": &GRPCBackendPlugin{
				Factory: opts.BackendFactoryFunc,
				Logger:  logger,
			},
		},
		4: {
			"backend": &GRPCBackendPlugin{
				Factory: opts.BackendFactoryFunc,
				Logger:  logger,
			},
		},
		5: {
			"backend": &GRPCBackendPlugin{
				Factory:             opts.BackendFactoryFunc,
				MultiplexingSupport: false,
				Logger:              logger,
			},
		},
	}

	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		return err
	}

	serveOpts := &plugin.ServeConfig{
		HandshakeConfig:  HandshakeConfig,
		VersionedPlugins: pluginSets,
		TLSProvider:      opts.TLSProviderFunc,
		Logger:           logger,

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			opts = append(opts, grpc.MaxRecvMsgSize(math.MaxInt32))
			opts = append(opts, grpc.MaxSendMsgSize(math.MaxInt32))
			return plugin.DefaultGRPCServer(opts)
		},
	}

	plugin.Serve(serveOpts)

	return nil
}

// ServeMultiplex is a helper function used to serve a backend plugin. This
// should be ran on the plugin's main process.
func ServeMultiplex(opts *ServeOpts) error {
	logger := opts.Logger
	if logger == nil {
		logger = log.New(&log.LoggerOptions{
			Level:      log.Trace,
			Output:     os.Stderr,
			JSONFormat: true,
		})
	}

	// pluginMap is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		// Version 3 used to supports both protocols. We want to keep it around
		// since it's possible old plugins built against this version will still
		// work with gRPC. There is currently no difference between version 3
		// and version 4.
		3: {
			"backend": &GRPCBackendPlugin{
				Factory: opts.BackendFactoryFunc,
				Logger:  logger,
			},
		},
		4: {
			"backend": &GRPCBackendPlugin{
				Factory: opts.BackendFactoryFunc,
				Logger:  logger,
			},
		},
		5: {
			"backend": &GRPCBackendPlugin{
				Factory:             opts.BackendFactoryFunc,
				MultiplexingSupport: true,
				Logger:              logger,
			},
		},
	}

	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		return err
	}

	serveOpts := &plugin.ServeConfig{
		HandshakeConfig:  HandshakeConfig,
		VersionedPlugins: pluginSets,
		Logger:           logger,

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			opts = append(opts, grpc.MaxRecvMsgSize(math.MaxInt32))
			opts = append(opts, grpc.MaxSendMsgSize(math.MaxInt32))
			return plugin.DefaultGRPCServer(opts)
		},

		// TLSProvider is required to support v3 and v4 plugins.
		// It will be ignored for v5 which uses AutoMTLS
		TLSProvider: opts.TLSProviderFunc,
	}

	plugin.Serve(serveOpts)

	return nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var HandshakeConfig = plugin.HandshakeConfig{
	MagicCookieKey:   "VAULT_BACKEND_PLUGIN",
	MagicCookieValue: "6669da05-b1c8-4f49-97d9-c8e5bed98e20",
}

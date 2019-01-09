package plugin

import (
	"crypto/tls"
	"math"
	"os"

	"google.golang.org/grpc"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
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
		3: plugin.PluginSet{
			"backend": &BackendPlugin{
				GRPCBackendPlugin: &GRPCBackendPlugin{
					Factory: opts.BackendFactoryFunc,
					Logger:  logger,
				},
			},
		},
		4: plugin.PluginSet{
			"backend": &GRPCBackendPlugin{
				Factory: opts.BackendFactoryFunc,
				Logger:  logger,
			},
		},
	}

	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		return err
	}

	serveOpts := &plugin.ServeConfig{
		HandshakeConfig:  handshakeConfig,
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

	// If we do not have gRPC support fallback to version 3
	// Remove this block in 0.13
	if !pluginutil.GRPCSupport() {
		serveOpts.GRPCServer = nil
		delete(pluginSets, 4)
	}

	plugin.Serve(serveOpts)

	return nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  4,
	MagicCookieKey:   "VAULT_BACKEND_PLUGIN",
	MagicCookieValue: "6669da05-b1c8-4f49-97d9-c8e5bed98e20",
}

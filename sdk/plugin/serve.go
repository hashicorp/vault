package plugin

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"google.golang.org/grpc"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// BackendPluginName is the name of the plugin that can be
// dispensed from the plugin server.
const BackendPluginName = "backend"

const (
	// envVaultReattachPlugins is the environment variable used by Vault CLI
	// to directly connect to already running plugin processes, such as those
	// being inspected by debugging processes. When connecting to plugins in
	// this manner, Vault CLI disables certain plugin handshake checks and
	// will not stop the plugin process.
	envVaultReattachPlugins = "VAULT_REATTACH_PLUGINS"
)

type TLSProviderFunc func() (*tls.Config, error)

type ServeOpts struct {
	BackendFactoryFunc logical.Factory
	TLSProviderFunc    TLSProviderFunc
	Logger             log.Logger

	// Debug starts a debug server and controls its lifecycle, printing the
	// information needed for Vault to connect to the plugin to stdout.
	// os.Interrupt will be captured and used to stop the server.
	Debug bool
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

// ServeConfig contains the configured options for how a plugin should be
// served.
type ServeConfig struct {
	debugCtx     context.Context
	debugCh      chan *plugin.ReattachConfig
	debugCloseCh chan struct{}

	managedDebugReattachConfigTimeout time.Duration
	managedDebugStopSignals           []os.Signal
}

// ServeMultiplex is a helper function used to serve a backend plugin. This
// should be ran on the plugin's main process.
func ServeMultiplex(opts *ServeOpts) error {
	// Defaults
	conf := ServeConfig{
		managedDebugReattachConfigTimeout: 2 * time.Second,
		managedDebugStopSignals:           []os.Signal{os.Interrupt},
	}

	logger := opts.Logger
	if logger == nil {
		logger = log.New(&log.LoggerOptions{
			Level:      log.Info,
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

	serveConfig := &plugin.ServeConfig{
		HandshakeConfig:  HandshakeConfig,
		VersionedPlugins: pluginSets,
		Logger:           logger,

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			opts = append(opts, grpc.MaxRecvMsgSize(math.MaxInt32))
			opts = append(opts, grpc.MaxSendMsgSize(math.MaxInt32))
			return plugin.DefaultGRPCServer(opts)
		},
	}

	if opts.Debug {
		ctx, cancel := context.WithCancel(context.Background())
		signalCh := make(chan os.Signal, len(conf.managedDebugStopSignals))

		signal.Notify(signalCh, conf.managedDebugStopSignals...)

		defer func() {
			signal.Stop(signalCh)
			cancel()
		}()

		go func() {
			select {
			case <-signalCh:
				cancel()
			case <-ctx.Done():
			}
		}()

		conf.debugCh = make(chan *plugin.ReattachConfig)
		conf.debugCloseCh = make(chan struct{})
		conf.debugCtx = ctx
	}

	if conf.debugCh != nil {
		serveConfig.Test = &plugin.ServeTestConfig{
			Context:          conf.debugCtx,
			ReattachConfigCh: conf.debugCh,
			CloseCh:          conf.debugCloseCh,
		}
	}

	if !opts.Debug {
		plugin.Serve(serveConfig)
		return nil
	}

	go plugin.Serve(serveConfig)

	var pluginReattachConfig *plugin.ReattachConfig

	select {
	case pluginReattachConfig = <-conf.debugCh:
	case <-time.After(conf.managedDebugReattachConfigTimeout):
		return errors.New("timeout waiting on reattach configuration")
	}

	if pluginReattachConfig == nil {
		return errors.New("nil reattach configuration received")
	}
	fmt.Printf("\n\n%#+v\n\n", pluginReattachConfig)

	// Duplicate implementation is required because the go-plugin
	// ReattachConfig.Addr implementation is not friendly for JSON encoding
	// and to avoid importing terraform-exec.
	type reattachConfigAddr struct {
		Network string
		String  string
	}

	type reattachConfig struct {
		Protocol        string
		ProtocolVersion int
		Pid             int
		Test            bool
		Addr            reattachConfigAddr
	}

	reattachBytes, err := json.Marshal(map[string]reattachConfig{
		"backend": {
			Protocol:        string(pluginReattachConfig.Protocol),
			ProtocolVersion: pluginReattachConfig.ProtocolVersion,
			Pid:             pluginReattachConfig.Pid,
			Test:            pluginReattachConfig.Test,
			Addr: reattachConfigAddr{
				Network: pluginReattachConfig.Addr.Network(),
				String:  pluginReattachConfig.Addr.String(),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Error building reattach string: %w", err)
	}

	reattachStr := string(reattachBytes)

	// This is currently intended to be executed via plugin main function and
	// human friendly, so output directly to stdout.
	fmt.Printf("Plugin started. To attach Vault CLI, set the %s environment variable with the following:\n\n", envVaultReattachPlugins)

	switch runtime.GOOS {
	case "windows":
		fmt.Printf("\tCommand Prompt:\tset \"%s=%s\"\n", envVaultReattachPlugins, reattachStr)
		fmt.Printf("\tPowerShell:\t$env:%s='%s'\n", envVaultReattachPlugins, strings.ReplaceAll(reattachStr, `'`, `''`))
	default:
		fmt.Printf("\t%s='%s'\n", envVaultReattachPlugins, strings.ReplaceAll(reattachStr, `'`, `'"'"'`))
	}

	fmt.Println("")

	// Wait for the server to be done.
	<-conf.debugCloseCh

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

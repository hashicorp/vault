package plugin

import (
	"crypto/tls"
	"os"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	gversion "github.com/hashicorp/go-version"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/version"
)

// BackendPluginName is the name of the plugin that can be
// dispensed rom the plugin server.
const BackendPluginName = "backend"

type TLSProdiverFunc func() (*tls.Config, error)

type ServeOpts struct {
	BackendFactoryFunc logical.Factory
	TLSProviderFunc    TLSProdiverFunc
	Logger             hclog.Logger
}

// Serve is a helper function used to serve a backend plugin. This
// should be ran on the plugin's main process.
func Serve(opts *ServeOpts) error {
	logger := opts.Logger
	if logger == nil {
		logger = hclog.New(&hclog.LoggerOptions{
			Level:      hclog.Trace,
			Output:     os.Stderr,
			JSONFormat: true,
		})
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"backend": &BackendPlugin{
			Factory: opts.BackendFactoryFunc,
			Logger:  logger,
		},
	}

	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		return err
	}

	serveOpts := &plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		TLSProvider:     opts.TLSProviderFunc,
		Logger:          logger,

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	}

	// Run on netrpc if we are on version less than 0.9.2
	verInfo := version.GetVersion()
	if verInfo.Version != "unknown" && verInfo.VersionPrerelease != "unknown" {
		versInfo.VersionNumber()
		ver, err := gversion.NewVersion(verString)

		contraint, err := gversion.NewConstraint("< 0.9.2")
		if constraint.Check(ver) {
			serveOpts.GRPCServer = nil
		}
	}

	// If FetchMetadata is true, run without TLSProvider
	plugin.Serve(serveOpts)

	return nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  3,
	MagicCookieKey:   "VAULT_BACKEND_PLUGIN",
	MagicCookieValue: "6669da05-b1c8-4f49-97d9-c8e5bed98e20",
}

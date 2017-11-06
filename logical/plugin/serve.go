package plugin

import (
	"crypto/tls"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
)

// BackendPluginName is the name of the plugin that can be
// dispensed rom the plugin server.
const BackendPluginName = "backend"

type BackendFactoryFunc func(*logical.BackendConfig) (logical.Backend, error)
type TLSProdiverFunc func() (*tls.Config, error)

type ServeOpts struct {
	BackendFactoryFunc BackendFactoryFunc
	TLSProviderFunc    TLSProdiverFunc
}

// Serve is a helper function used to serve a backend plugin. This
// should be ran on the plugin's main process.
func Serve(opts *ServeOpts) error {
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"backend": &BackendPlugin{
			Factory: opts.BackendFactoryFunc,
		},
	}

	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		return err
	}

	// If FetchMetadata is true, run without TLSProvider
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		TLSProvider:     opts.TLSProviderFunc,
	})

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

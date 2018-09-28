package dbplugin

import (
	"crypto/tls"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
)

// Serve is called from within a plugin and wraps the provided
// Database implementation in a databasePluginRPCServer object and starts a
// RPC server.
func Serve(db Database, tlsProvider func() (*tls.Config, error)) {
	plugin.Serve(ServeConfig(db, tlsProvider))
}

func ServeConfig(db Database, tlsProvider func() (*tls.Config, error)) *plugin.ServeConfig {
	// pluginSets is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		3: plugin.PluginSet{
			"database": &DatabasePlugin{
				GRPCDatabasePlugin: &GRPCDatabasePlugin{
					Impl: db,
				},
			},
		},
		4: plugin.PluginSet{
			"database": &GRPCDatabasePlugin{
				Impl: db,
			},
		},
	}

	conf := &plugin.ServeConfig{
		HandshakeConfig:  handshakeConfig,
		VersionedPlugins: pluginSets,
		TLSProvider:      tlsProvider,
		GRPCServer:       plugin.DefaultGRPCServer,
	}

	// If we do not have gRPC support fallback to version 3
	// Remove this block in 0.13
	if !pluginutil.GRPCSupport() {
		conf.GRPCServer = nil
		delete(conf.VersionedPlugins, 4)
	}

	return conf
}

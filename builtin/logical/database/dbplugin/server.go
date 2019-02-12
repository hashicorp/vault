package dbplugin

import (
	"crypto/tls"

	plugin "github.com/hashicorp/go-plugin"
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
		// Version 3 used to supports both protocols. We want to keep it around
		// since it's possible old plugins built against this version will still
		// work with gRPC. There is currently no difference between version 3
		// and version 4.
		3: plugin.PluginSet{
			"database": &GRPCDatabasePlugin{
				Impl: db,
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

	return conf
}

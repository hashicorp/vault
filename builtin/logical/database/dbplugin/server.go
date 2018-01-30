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
	dbPlugin := &DatabasePlugin{
		impl: db,
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"database": dbPlugin,
	}

	conf := &plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		TLSProvider:     tlsProvider,
		GRPCServer:      plugin.DefaultGRPCServer,
	}

	if !pluginutil.GRPCSupport() {
		conf.GRPCServer = nil
	}

	return conf
}

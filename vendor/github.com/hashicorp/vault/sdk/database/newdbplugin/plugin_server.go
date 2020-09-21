package newdbplugin

import (
	"crypto/tls"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// Serve is called from within a plugin and wraps the provided
// Database implementation in a databasePluginRPCServer object and starts a
// RPC server.
func Serve(db Database, tlsProvider func() (*tls.Config, error)) {
	plugin.Serve(ServeConfig(db, tlsProvider))
}

func ServeConfig(db Database, tlsProvider func() (*tls.Config, error)) *plugin.ServeConfig {
	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// pluginSets is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		5: plugin.PluginSet{
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

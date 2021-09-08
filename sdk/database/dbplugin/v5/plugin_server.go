package dbplugin

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// Serve is called from within a plugin and wraps the provided
// Database implementation in a databasePluginRPCServer object and starts a
// RPC server.
func Serve(dbFactory func() (Database, error)) {
	plugin.Serve(ServeConfig(dbFactory))
}

func ServeConfig(dbFactory func() (Database, error)) *plugin.ServeConfig {
	err := pluginutil.OptionallyEnableMlock()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// dbPluginSet is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		6: {
			"database": &GRPCDatabasePlugin{
				DBFactory: dbFactory,
			},
		},
	}

	conf := &plugin.ServeConfig{
		HandshakeConfig:  handshakeConfig,
		VersionedPlugins: pluginSets,
		GRPCServer:       plugin.DefaultGRPCServer,
	}

	return conf
}

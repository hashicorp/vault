package dbplugin

import (
	"errors"
	"sync"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	log "github.com/mgutz/logxi/v1"
)

// DatabasePluginClient embeds a databasePluginRPCClient and wraps it's Close
// method to also call Kill() on the plugin.Client.
type DatabasePluginClient struct {
	client *plugin.Client
	sync.Mutex

	Database
}

func (dc *DatabasePluginClient) Close() error {
	err := dc.Database.Close()
	dc.client.Kill()

	return err
}

// newPluginClient returns a databaseRPCClient with a connection to a running
// plugin. The client is wrapped in a DatabasePluginClient object to ensure the
// plugin is killed on call of Close().
func newPluginClient(sys pluginutil.RunnerUtil, pluginRunner *pluginutil.PluginRunner, logger log.Logger) (Database, error) {
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"database": new(DatabasePlugin),
	}

	client, err := pluginRunner.Run(sys, pluginMap, handshakeConfig, []string{}, logger)
	if err != nil {
		return nil, err
	}

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("database")
	if err != nil {
		return nil, err
	}

	// We should have a database type now. This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	var db Database
	switch raw.(type) {
	case *gRPCClient:
		db = raw.(*gRPCClient)
	case *databasePluginRPCClient:
		logger.Warn("database: plugin is using deprecated net RPC transport, recompile plugin to upgrade to gRPC", "plugin", pluginRunner.Name)
		db = raw.(*databasePluginRPCClient)
	default:
		return nil, errors.New("unsupported client type")
	}

	// Wrap RPC implimentation in DatabasePluginClient
	return &DatabasePluginClient{
		client:   client,
		Database: db,
	}, nil
}

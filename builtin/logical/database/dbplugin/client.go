package dbplugin

import (
	"context"
	"errors"
	"sync"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
)

// DatabasePluginClient embeds a databasePluginRPCClient and wraps it's Close
// method to also call Kill() on the plugin.Client.
type DatabasePluginClient struct {
	client *plugin.Client
	sync.Mutex

	Database
}

// This wraps the Close call and ensures we both close the database connection
// and kill the plugin.
func (dc *DatabasePluginClient) Close() error {
	err := dc.Database.Close()
	dc.client.Kill()

	return err
}

// NewPluginClient returns a databaseRPCClient with a connection to a running
// plugin. The client is wrapped in a DatabasePluginClient object to ensure the
// plugin is killed on call of Close().
func NewPluginClient(ctx context.Context, sys pluginutil.RunnerUtil, pluginRunner *pluginutil.PluginRunner, logger log.Logger, isMetadataMode bool) (Database, error) {

	// pluginSets is the map of plugins we can dispense.
	pluginSets := map[int]plugin.PluginSet{
		// Version 3 supports both protocols
		3: plugin.PluginSet{
			"database": &DatabasePlugin{
				GRPCDatabasePlugin: new(GRPCDatabasePlugin),
			},
		},
		// Version 4 only supports gRPC
		4: plugin.PluginSet{
			"database": new(GRPCDatabasePlugin),
		},
	}

	var client *plugin.Client
	var err error
	if isMetadataMode {
		client, err = pluginRunner.RunMetadataMode(ctx, sys, pluginSets, handshakeConfig, []string{}, logger)
	} else {
		client, err = pluginRunner.Run(ctx, sys, pluginSets, handshakeConfig, []string{}, logger)
	}
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
		logger.Warn("plugin is using deprecated netRPC transport, recompile plugin to upgrade to gRPC", "plugin", pluginRunner.Name)
		db = raw.(*databasePluginRPCClient)
	default:
		return nil, errors.New("unsupported client type")
	}

	// Wrap RPC implementation in DatabasePluginClient
	return &DatabasePluginClient{
		client:   client,
		Database: db,
	}, nil
}

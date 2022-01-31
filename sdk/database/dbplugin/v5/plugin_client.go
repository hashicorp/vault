package dbplugin

import (
	"context"
	"errors"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

type DatabasePluginClient struct {
	client pluginutil.Multiplexer
	Database
}

// This wraps the Close call and ensures we both close the database connection
// and kill the plugin.
func (dc *DatabasePluginClient) Close() error {
	err := dc.Database.Close()
	dc.client.Close()

	return err
}

// pluginSets is the map of plugins we can dispense.
var PluginSets = map[int]plugin.PluginSet{
	5: {
		"database": &GRPCDatabasePlugin{multiplexingSupport: false},
	},
	6: {
		"database": &GRPCDatabasePlugin{multiplexingSupport: true},
	},
}

// NewPluginClient returns a databaseRPCClient with a connection to a running
// plugin.
func NewPluginClient(ctx context.Context, sys pluginutil.RunnerUtil, pluginRunner *pluginutil.PluginRunner, logger log.Logger, isMetadataMode bool) (Database, error) {
	pluginClient, err := sys.NewPluginClient(ctx, pluginRunner, logger, isMetadataMode)
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := pluginClient.Dispense("database")
	if err != nil {
		return nil, err
	}

	// We should have a database type now. This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	var db Database
	switch raw.(type) {
	case gRPCClient:

		gRPCClient := raw.(gRPCClient)

		// Wrap clientConn with our implementation so that we can inject the
		// ID into the context
		cc := &databaseClientConn{
			ClientConn: pluginClient.Conn(),
			id:         pluginClient.ID(),
		}
		gRPCClient.client = proto.NewDatabaseClient(cc)
		db = gRPCClient
	default:
		return nil, errors.New("unsupported client type")
	}

	return &DatabasePluginClient{
		client:   pluginClient,
		Database: db,
	}, nil
}

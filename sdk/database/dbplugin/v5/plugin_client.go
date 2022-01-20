package dbplugin

import (
	"context"
	"errors"
	"sync"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

type DatabasePluginClient struct {
	sync.Mutex

	Database
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
	rpcClient, id, err := sys.NewPluginClient(ctx, pluginRunner, logger, isMetadataMode)
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
	case gRPCClient:

		gRPCClient := raw.(gRPCClient)

		gc := rpcClient.(*plugin.GRPCClient)
		// Wrap clientConn with our implementation and get rid of middleware
		// and then cast it back
		cc := &databaseClientConn{
			ClientConn: gc.Conn,
			id:         id,
		}
		gRPCClient.client = proto.NewDatabaseClient(cc)
		db = gRPCClient
	default:
		return nil, errors.New("unsupported client type")
	}

	// Wrap RPC implementation in DatabasePluginClient
	return &DatabasePluginClient{
		Database: db,
	}, nil
}

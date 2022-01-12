package dbplugin

import (
	"context"
	"errors"
	"sync"

	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

type DatabasePluginClient struct {
	id string
	sync.Mutex

	Database
}

// pluginSets is the map of plugins we can dispense.
// TODO(JM): add multiplexingSupport
var PluginSets = map[int]plugin.PluginSet{
	5: {
		"database": new(GRPCDatabasePlugin),
	},
}

// NewPluginClient returns a databaseRPCClient with a connection to a running
// plugin.
func NewPluginClient(ctx context.Context, sys pluginutil.RunnerUtil, pluginRunner *pluginutil.PluginRunner, logger log.Logger, isMetadataMode bool) (Database, error) {
	rpcClient, err := sys.NewPluginClient(ctx, pluginRunner, logger, isMetadataMode)
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
		db = raw.(gRPCClient)
	default:
		return nil, errors.New("unsupported client type")
	}

	// Wrap RPC implementation in DatabasePluginClient
	return &DatabasePluginClient{
		Database: db,
		id:       "",
	}, nil
}

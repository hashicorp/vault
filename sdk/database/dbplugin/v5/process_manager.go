package dbplugin

import (
	"context"
	"fmt"
	"sync"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

type DatabaseManager struct {
	mu     *sync.RWMutex
	logger log.Logger

	dbs                map[string]Database       // config name -> Database
	multiplexedPlugins map[string]*pluginProcess // plugin name -> plugin process
}

// TODO: Flesh out constructor
func NewDatabaseManager(logger log.Logger) *DatabaseManager {
	dbm := &DatabaseManager{
		mu:     new(sync.RWMutex),
		logger: logger,
	}

	return dbm
}

var (
	dbPluginSet = map[int]plugin.PluginSet{
		5: {
			"database": new(GRPCDatabasePlugin),
		},
		6: {
			"database": new(GRPCDatabasePlugin),
		},
	}
)

func (dbm *DatabaseManager) Get(ctx context.Context, id string, pluginName string, sys pluginutil.LookRunnerUtil) (Database, error) {
	// Check if a v5 instance already exists
	dbm.mu.RLock()
	db, exists := dbm.dbs[id]
	if exists {
		dbm.mu.RUnlock()
		return db, nil
	}

	// If a v5 instance doesn't exist, check to see if it's multiplexed
	multiplexedProcess, exists := dbm.multiplexedPlugins[pluginName]
	if exists {
		db, err := multiplexedProcess.getDatabase(id)
		if err != nil {
			dbm.mu.RUnlock()
			return nil, fmt.Errorf("failed to get multiplexed database: %w", err)
		}
		dbm.dbs[id] = db
		dbm.mu.RUnlock()
		return db, nil
	}

	// Database isn't in the v5 cache, nor is there a multiplexed plugin instance running
	// Upgrade to a write lock
	dbm.mu.RUnlock()
	dbm.mu.Lock()
	defer dbm.mu.Unlock()

	// Check if the Database was added by a different goroutine
	db, exists = dbm.dbs[id]
	if exists {
		return db, nil
	}

	// Don't need to check the multiplexedPlugins map since the DB will be added to both

	// No plugin instance found, start a new one
	pluginRunner, err := sys.LookupPlugin(ctx, pluginName, consts.PluginTypeDatabase)
	if err != nil {
		return nil, fmt.Errorf("failed to find plugin: %w", err)
	}

	// TODO: Check pluginRunner.Builtin - if true, create the DB via BuiltinFactory(), then add it to the dbs cache & return

	// TODO: Pull plugin startup logic into its own function so it can be referenced in plugin_catalog.go isDatabasePlugin()
	process, err := pluginRunner.RunConfig(ctx,
		pluginutil.Runner(sys),
		pluginutil.PluginSets(dbPluginSet),
		pluginutil.HandshakeConfig(handshakeConfig),
		pluginutil.Logger(dbm.logger),
		pluginutil.MetadataMode(false),
		pluginutil.AutoMTLS(true),
	)
	if err != nil {
		return nil, err
	}

	namedLogger := dbm.logger.Named(pluginName)

	// Switch on the version of the plugin
	// If v5, it's a 1:1 between the Database instance and the plugin process
	//   - Add the Database to the dbs map and return
	// If v6, it's a *:1 between the Database instances and the plugin process
	//   - Create a multiplexed pluginProcess object that will manage
	//     the plugin process itself
	//   - Get a new Database instance from the pluginProcess
	switch process.NegotiatedVersion() {
	case 5:
		db, err := dbm.newV5Database(process, namedLogger)
		if err != nil {
			// Kill the process because something went wrong, and we don't want dangling processes
			process.Kill()
			return nil, fmt.Errorf("failed to initialize database process: %w", err)
		}
		// TODO: Add wrapper to remove DB from cache
		dbm.dbs[id] = db
		return db, nil
	case 6:
		// TODO: Should the logger include the ID as well?
		multiplexedPlugin, err := newPluginProcess(dbm.logger, process)
		if err != nil {
			// Kill the process because something went wrong, and we don't want dangling processes
			process.Kill()
			return nil, fmt.Errorf("failed to initialize multiplexed plugin process: %w", err)
		}

		db, err := multiplexedPlugin.getDatabase(id)
		if err != nil {
			// Kill the process because something went wrong, and we don't want dangling processes
			process.Kill()
			return nil, fmt.Errorf("failed to initialize multiplexed plugin process: %w", err)
		}

		dbm.dbs[id] = db
		dbm.multiplexedPlugins[id] = multiplexedPlugin
		// TODO: Add wrapper to remove DB from cache
		return db, nil
	default:
		return nil, fmt.Errorf("unsupported plugin version %d", process.NegotiatedVersion())
	}
}

func (dbm *DatabaseManager) newV5Database(process *plugin.Client, logger log.Logger) (Database, error) {
	rpcClient, err := process.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	rawDB, err := rpcClient.Dispense("database")
	if err != nil {
		return nil, err
	}

	var db Database
	db, ok := rawDB.(*gRPCClient)
	if !ok {
		return nil, fmt.Errorf("unsupported client type")
	}

	pluginType, err := db.Type()
	if err != nil {
		return nil, err
	}

	db = addMiddleware(db, pluginType, logger)

	db = &DatabasePluginClient{
		client:   process,
		Database: db,
	}

	return db, nil
}

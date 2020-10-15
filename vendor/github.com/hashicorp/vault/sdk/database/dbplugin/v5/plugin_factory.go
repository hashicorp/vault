package dbplugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// PluginFactory is used to build plugin database types. It wraps the database
// object in a logging and metrics middleware.
func PluginFactory(ctx context.Context, pluginName string, sys pluginutil.LookRunnerUtil, logger log.Logger) (Database, error) {
	// Look for plugin in the plugin catalog
	pluginRunner, err := sys.LookupPlugin(ctx, pluginName, consts.PluginTypeDatabase)
	if err != nil {
		return nil, err
	}

	namedLogger := logger.Named(pluginName)

	var transport string
	var db Database
	if pluginRunner.Builtin {
		// Plugin is builtin so we can retrieve an instance of the interface
		// from the pluginRunner. Then cast it to a Database.
		dbRaw, err := pluginRunner.BuiltinFactory()
		if err != nil {
			return nil, errwrap.Wrapf("error initializing plugin: {{err}}", err)
		}

		var ok bool
		db, ok = dbRaw.(Database)
		if !ok {
			return nil, fmt.Errorf("unsupported database type: %q", pluginName)
		}

		transport = "builtin"

	} else {
		// create a DatabasePluginClient instance
		db, err = NewPluginClient(ctx, sys, pluginRunner, namedLogger, false)
		if err != nil {
			return nil, err
		}

		// Switch on the underlying database client type to get the transport
		// method.
		switch db.(*DatabasePluginClient).Database.(type) {
		case *gRPCClient:
			transport = "gRPC"
		}

	}

	typeStr, err := db.Type()
	if err != nil {
		return nil, errwrap.Wrapf("error getting plugin type: {{err}}", err)
	}

	// Wrap with metrics middleware
	db = &databaseMetricsMiddleware{
		next:    db,
		typeStr: typeStr,
	}

	// Wrap with tracing middleware
	if namedLogger.IsTrace() {
		db = &databaseTracingMiddleware{
			next:   db,
			logger: namedLogger.With("transport", transport),
		}
	}

	return db, nil
}

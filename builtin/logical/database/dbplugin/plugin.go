package dbplugin

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/pluginutil"
)

// Database is the interface that all database objects must implement.
type Database interface {
	Type() (string, error)
	CreateUser(ctx context.Context, statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error)
	RenewUser(ctx context.Context, statements Statements, username string, expiration time.Time) error
	RevokeUser(ctx context.Context, statements Statements, username string) error

	RotateRootCredentials(ctx context.Context, statements []string) (config map[string]interface{}, err error)

	Init(ctx context.Context, config map[string]interface{}, verifyConnection bool) (saveConfig map[string]interface{}, err error)
	Close() error

	// DEPRECATED, will be removed in a future plugin version bump.
	Initialize(ctx context.Context, config map[string]interface{}, verifyConnection bool) (err error)
}

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

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  4,
	MagicCookieKey:   "VAULT_DATABASE_PLUGIN",
	MagicCookieValue: "926a0820-aea2-be28-51d6-83cdf00e8edb",
}

var _ plugin.Plugin = &GRPCDatabasePlugin{}
var _ plugin.GRPCPlugin = &GRPCDatabasePlugin{}

// GRPCDatabasePlugin is the plugin.Plugin implementation that only supports GRPC
// transport
type GRPCDatabasePlugin struct {
	Impl Database

	// Embeding this will disable the netRPC protocol
	plugin.NetRPCUnsupportedPlugin
}

func (d GRPCDatabasePlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	impl := &DatabaseErrorSanitizerMiddleware{
		next: d.Impl,
	}

	RegisterDatabaseServer(s, &gRPCServer{impl: impl})
	return nil
}

func (GRPCDatabasePlugin) GRPCClient(doneCtx context.Context, _ *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &gRPCClient{
		client:     NewDatabaseClient(c),
		clientConn: c,
		doneCtx:    doneCtx,
	}, nil
}

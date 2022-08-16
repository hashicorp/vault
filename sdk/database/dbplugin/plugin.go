package dbplugin

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// Database is the interface that all database objects must implement.
type Database interface {
	// Type returns the TypeName for the particular database backend
	// implementation. This type name is usually set as a constant within the
	// database backend implementation, e.g. "mysql" for the MySQL database
	// backend.
	Type() (string, error)

	// CreateUser is called on `$ vault read database/creds/:role-name` and it's
	// also the first time anything is touched from `$ vault write
	// database/roles/:role-name`. This is likely to be the highest-throughput
	// method for most plugins.
	CreateUser(ctx context.Context, statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error)

	// RenewUser is triggered by a renewal call to the API. In many database
	// backends, this triggers a call on the underlying database that extends a
	// VALID UNTIL clause on a user. However, if no such need exists, setting
	// this as a NO-OP means that when renewal is called, the lease renewal time
	// is pushed further out as appropriate, thus pushing out the time until the
	// RevokeUser method is called.
	RenewUser(ctx context.Context, statements Statements, username string, expiration time.Time) error

	// RevokeUser is triggered either automatically by a lease expiration, or by
	// a revocation call to the API.
	RevokeUser(ctx context.Context, statements Statements, username string) error

	// RotateRootCredentials is triggered by a root credential rotation call to
	// the API.
	RotateRootCredentials(ctx context.Context, statements []string) (config map[string]interface{}, err error)

	// GenerateCredentials returns a generated password for the plugin. This is
	// used in combination with SetCredentials to set a specific password for a
	// database user and preserve the password in WAL entries.
	GenerateCredentials(ctx context.Context) (string, error)

	// SetCredentials uses provided information to create or set the credentials
	// for a database user. Unlike CreateUser, this method requires both a
	// username and a password given instead of generating them. This is used for
	// creating and setting the password of static accounts, as well as rolling
	// back passwords in the database in the event an updated database fails to
	// save in Vault's storage.
	SetCredentials(ctx context.Context, statements Statements, staticConfig StaticUserConfig) (username string, password string, err error)

	// Init is called on `$ vault write database/config/:db-name`, or when you
	// do a creds call after Vault's been restarted. The config provided won't
	// hold all the keys and values provided in the API call, some will be
	// stripped by the database engine before the config is provided. The config
	// returned will be stored, which will persist it across shutdowns.
	Init(ctx context.Context, config map[string]interface{}, verifyConnection bool) (saveConfig map[string]interface{}, err error)

	// Close attempts to close the underlying database connection that was
	// established by the backend.
	Close() error
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

var (
	_ plugin.Plugin     = &GRPCDatabasePlugin{}
	_ plugin.GRPCPlugin = &GRPCDatabasePlugin{}
)

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

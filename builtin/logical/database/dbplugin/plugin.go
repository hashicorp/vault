package dbplugin

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	log "github.com/mgutz/logxi/v1"
)

// Database is the interface that all database objects must implement.
type Database interface {
	Type() (string, error)
	CreateUser(statements Statements, usernameConfig UsernameConfig, expiration time.Time) (username string, password string, err error)
	RenewUser(statements Statements, username string, expiration time.Time) error
	RevokeUser(statements Statements, username string) error

	Initialize(config map[string]interface{}, verifyConnection bool) error
	Close() error
}

// Statements set in role creation and passed into the database type's functions.
type Statements struct {
	CreationStatements   string `json:"creation_statments" mapstructure:"creation_statements" structs:"creation_statments"`
	RevocationStatements string `json:"revocation_statements" mapstructure:"revocation_statements" structs:"revocation_statements"`
	RollbackStatements   string `json:"rollback_statements" mapstructure:"rollback_statements" structs:"rollback_statements"`
	RenewStatements      string `json:"renew_statements" mapstructure:"renew_statements" structs:"renew_statements"`
}

// UsernameConfig is used to configure prefixes for the username to be
// generated.
type UsernameConfig struct {
	DisplayName string
	RoleName    string
}

// PluginFactory is used to build plugin database types. It wraps the database
// object in a logging and metrics middleware.
func PluginFactory(pluginName string, sys pluginutil.LookRunnerUtil, logger log.Logger) (Database, error) {
	// Look for plugin in the plugin catalog
	pluginRunner, err := sys.LookupPlugin(pluginName)
	if err != nil {
		return nil, err
	}

	var db Database
	if pluginRunner.Builtin {
		// Plugin is builtin so we can retrieve an instance of the interface
		// from the pluginRunner. Then cast it to a Database.
		dbRaw, err := pluginRunner.BuiltinFactory()
		if err != nil {
			return nil, fmt.Errorf("error getting plugin type: %s", err)
		}

		var ok bool
		db, ok = dbRaw.(Database)
		if !ok {
			return nil, fmt.Errorf("unsuported database type: %s", pluginName)
		}

	} else {
		// create a DatabasePluginClient instance
		db, err = newPluginClient(sys, pluginRunner, logger)
		if err != nil {
			return nil, err
		}
	}

	typeStr, err := db.Type()
	if err != nil {
		return nil, fmt.Errorf("error getting plugin type: %s", err)
	}

	// Wrap with metrics middleware
	db = &databaseMetricsMiddleware{
		next:    db,
		typeStr: typeStr,
	}

	// Wrap with tracing middleware
	if logger.IsTrace() {
		db = &databaseTracingMiddleware{
			next:    db,
			typeStr: typeStr,
			logger:  logger,
		}
	}

	return db, nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  3,
	MagicCookieKey:   "VAULT_DATABASE_PLUGIN",
	MagicCookieValue: "926a0820-aea2-be28-51d6-83cdf00e8edb",
}

// DatabasePlugin implements go-plugin's Plugin interface. It has methods for
// retrieving a server and a client instance of the plugin.
type DatabasePlugin struct {
	impl Database
}

func (d DatabasePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &databasePluginRPCServer{impl: d.impl}, nil
}

func (DatabasePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &databasePluginRPCClient{client: c}, nil
}

// ---- RPC Request Args Domain ----

type InitializeRequest struct {
	Config           map[string]interface{}
	VerifyConnection bool
}

type CreateUserRequest struct {
	Statements     Statements
	UsernameConfig UsernameConfig
	Expiration     time.Time
}

type RenewUserRequest struct {
	Statements Statements
	Username   string
	Expiration time.Time
}

type RevokeUserRequest struct {
	Statements Statements
	Username   string
}

// ---- RPC Response Args Domain ----

type CreateUserResponse struct {
	Username string
	Password string
}

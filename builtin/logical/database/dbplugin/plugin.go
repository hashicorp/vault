package dbplugin

import (
	"errors"
	"net/rpc"
	"time"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/vault/helper/pluginutil"
	log "github.com/mgutz/logxi/v1"
)

var (
	ErrEmptyPluginName = errors.New("empty plugin name")
)

// DatabaseType is the interface that all database objects must implement.
type DatabaseType interface {
	Type() string
	CreateUser(statements Statements, usernamePrefix string, expiration time.Time) (username string, password string, err error)
	RenewUser(statements Statements, username string, expiration time.Time) error
	RevokeUser(statements Statements, username string) error

	Initialize(map[string]interface{}) error
	Close() error
}

// Statements set in role creation and passed into the database type's functions.
// TODO: Add a way of setting defaults here.
type Statements struct {
	CreationStatements   string `json:"creation_statments" mapstructure:"creation_statements" structs:"creation_statments"`
	RevocationStatements string `json:"revocation_statements" mapstructure:"revocation_statements" structs:"revocation_statements"`
	RollbackStatements   string `json:"rollback_statements" mapstructure:"rollback_statements" structs:"rollback_statements"`
	RenewStatements      string `json:"renew_statements" mapstructure:"renew_statements" structs:"renew_statements"`
}

// PluginFactory is used to build plugin database types. It wraps the database
// object in a logging and metrics middleware.
func PluginFactory(pluginName string, sys pluginutil.LookWrapper, logger log.Logger) (DatabaseType, error) {
	if pluginName == "" {
		return nil, ErrEmptyPluginName
	}

	pluginMeta, err := sys.LookupPlugin(pluginName)
	if err != nil {
		return nil, err
	}

	db, err := newPluginClient(sys, pluginMeta)
	if err != nil {
		return nil, err
	}

	// Wrap with metrics middleware
	db = &databaseMetricsMiddleware{
		next:    db,
		typeStr: db.Type(),
	}

	// Wrap with tracing middleware
	db = &databaseTracingMiddleware{
		next:    db,
		typeStr: db.Type(),
		logger:  logger,
	}

	return db, nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "VAULT_DATABASE_PLUGIN",
	MagicCookieValue: "926a0820-aea2-be28-51d6-83cdf00e8edb",
}

type DatabasePlugin struct {
	impl DatabaseType
}

func (d DatabasePlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &databasePluginRPCServer{impl: d.impl}, nil
}

func (DatabasePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &databasePluginRPCClient{client: c}, nil
}

// ---- RPC Request Args Domain ----

type CreateUserRequest struct {
	Statements     Statements
	UsernamePrefix string
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

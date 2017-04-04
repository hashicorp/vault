package dbs

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	log "github.com/mgutz/logxi/v1"
)

const (
	postgreSQLTypeName = "postgres"
	mySQLTypeName      = "mysql"
	cassandraTypeName  = "cassandra"
	pluginTypeName     = "plugin"
)

var (
	ErrUnsupportedDatabaseType = errors.New("unsupported database type")
	ErrEmptyCreationStatement  = errors.New("empty creation statements")
	ErrEmptyPluginName         = errors.New("empty plugin name")
)

// Factory function definition
type Factory func(*DatabaseConfig, logical.SystemView, log.Logger) (DatabaseType, error)

// BuiltinFactory is used to build builtin database types. It wraps the database
// object in a logging and metrics middleware.
func BuiltinFactory(conf *DatabaseConfig, sys logical.SystemView, logger log.Logger) (DatabaseType, error) {
	var dbType DatabaseType

	switch conf.DatabaseType {
	case postgreSQLTypeName:
		connProducer := &sqlConnectionProducer{}
		connProducer.config = conf

		credsProducer := &sqlCredentialsProducer{
			displayNameLen: 23,
			usernameLen:    63,
		}

		dbType = &PostgreSQL{
			ConnectionProducer:  connProducer,
			CredentialsProducer: credsProducer,
		}

	case mySQLTypeName:
		connProducer := &sqlConnectionProducer{}
		connProducer.config = conf

		credsProducer := &sqlCredentialsProducer{
			displayNameLen: 4,
			usernameLen:    16,
		}

		dbType = &MySQL{
			ConnectionProducer:  connProducer,
			CredentialsProducer: credsProducer,
		}

	case cassandraTypeName:
		connProducer := &cassandraConnectionProducer{}
		connProducer.config = conf

		credsProducer := &cassandraCredentialsProducer{}

		dbType = &Cassandra{
			ConnectionProducer:  connProducer,
			CredentialsProducer: credsProducer,
		}

	default:
		return nil, ErrUnsupportedDatabaseType
	}

	// Wrap with metrics middleware
	dbType = &databaseMetricsMiddleware{
		next:    dbType,
		typeStr: dbType.Type(),
	}

	// Wrap with tracing middleware
	dbType = &databaseTracingMiddleware{
		next:    dbType,
		typeStr: dbType.Type(),
		logger:  logger,
	}

	return dbType, nil
}

// PluginFactory is used to build plugin database types. It wraps the database
// object in a logging and metrics middleware.
func PluginFactory(conf *DatabaseConfig, sys logical.SystemView, logger log.Logger) (DatabaseType, error) {
	if conf.PluginName == "" {
		return nil, ErrEmptyPluginName
	}

	pluginMeta, err := sys.LookupPlugin(conf.PluginName)
	if err != nil {
		return nil, err
	}

	// Make sure the database type is set to plugin
	conf.DatabaseType = pluginTypeName

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

// DatabaseType is the interface that all database objects must implement.
type DatabaseType interface {
	Type() string
	CreateUser(statements Statements, username, password, expiration string) error
	RenewUser(statements Statements, username, expiration string) error
	RevokeUser(statements Statements, username string) error

	Initialize(map[string]interface{}) error
	Close() error
	CredentialsProducer
}

// DatabaseConfig is used by the Factory function to configure a DatabaseType
// object.
type DatabaseConfig struct {
	DatabaseType string `json:"type" structs:"type" mapstructure:"type"`
	// ConnectionDetails stores the database specific connection settings needed
	// by each database type.
	ConnectionDetails     map[string]interface{} `json:"connection_details" structs:"connection_details" mapstructure:"connection_details"`
	MaxOpenConnections    int                    `json:"max_open_connections" structs:"max_open_connections" mapstructure:"max_open_connections"`
	MaxIdleConnections    int                    `json:"max_idle_connections" structs:"max_idle_connections" mapstructure:"max_idle_connections"`
	MaxConnectionLifetime time.Duration          `json:"max_connection_lifetime" structs:"max_connection_lifetime" mapstructure:"max_connection_lifetime"`
	PluginName            string                 `json:"plugin_name" structs:"plugin_name" mapstructure:"plugin_name"`
}

// GetFactory returns the appropriate factory method for the given database
// type.
func (dc *DatabaseConfig) GetFactory() Factory {
	if dc.DatabaseType == pluginTypeName {
		return PluginFactory
	}

	return BuiltinFactory
}

// Statments set in role creation and passed into the database type's functions.
// TODO: Add a way of setting defaults here.
type Statements struct {
	CreationStatements   string `json:"creation_statments" mapstructure:"creation_statements" structs:"creation_statments"`
	RevocationStatements string `json:"revocation_statements" mapstructure:"revocation_statements" structs:"revocation_statements"`
	RollbackStatements   string `json:"rollback_statements" mapstructure:"rollback_statements" structs:"rollback_statements"`
	RenewStatements      string `json:"renew_statements" mapstructure:"renew_statements" structs:"renew_statements"`
}

// Query templates a query for us.
func queryHelper(tpl string, data map[string]string) string {
	for k, v := range data {
		tpl = strings.Replace(tpl, fmt.Sprintf("{{%s}}", k), v, -1)
	}

	return tpl
}

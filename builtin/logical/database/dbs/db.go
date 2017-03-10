package dbs

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

const (
	postgreSQLTypeName = "postgres"
	mySQLTypeName      = "mysql"
	cassandraTypeName  = "cassandra"
	pluginTypeName     = "plugin"
)

var (
	ErrUnsupportedDatabaseType = errors.New("Unsupported database type")
)

type Factory func(*DatabaseConfig) (DatabaseType, error)

func BuiltinFactory(conf *DatabaseConfig) (DatabaseType, error) {
	switch conf.DatabaseType {
	case postgreSQLTypeName:
		var connProducer *sqlConnectionProducer
		err := mapstructure.Decode(conf.ConnectionDetails, &connProducer)
		if err != nil {
			return nil, err
		}
		connProducer.config = conf

		credsProducer := &sqlCredentialsProducer{
			displayNameLen: 23,
			usernameLen:    63,
		}

		return &PostgreSQL{
			ConnectionProducer:  connProducer,
			CredentialsProducer: credsProducer,
		}, nil

	case mySQLTypeName:
		var connProducer *sqlConnectionProducer
		err := mapstructure.Decode(conf.ConnectionDetails, &connProducer)
		if err != nil {
			return nil, err
		}
		connProducer.config = conf

		credsProducer := &sqlCredentialsProducer{
			displayNameLen: 4,
			usernameLen:    16,
		}

		return &MySQL{
			ConnectionProducer:  connProducer,
			CredentialsProducer: credsProducer,
		}, nil

	case cassandraTypeName:
		var connProducer *cassandraConnectionProducer
		err := mapstructure.Decode(conf.ConnectionDetails, &connProducer)
		if err != nil {
			return nil, err
		}
		connProducer.config = conf

		credsProducer := &cassandraCredentialsProducer{}

		return &Cassandra{
			ConnectionProducer:  connProducer,
			CredentialsProducer: credsProducer,
		}, nil
	}

	return nil, ErrUnsupportedDatabaseType
}

func PluginFactory(conf *DatabaseConfig) (DatabaseType, error) {
	if conf.PluginCommand == "" {
		return nil, errors.New("ERROR")
	}

	db, err := newPluginClient(conf.PluginCommand)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type DatabaseType interface {
	Type() string
	CreateUser(statements Statements, username, password, expiration string) error
	RenewUser(statements Statements, username, expiration string) error
	RevokeUser(statements Statements, username string) error

	Close()
	CredentialsProducer
}

type DatabaseConfig struct {
	DatabaseType          string                 `json:"type" structs:"type" mapstructure:"type"`
	ConnectionDetails     map[string]interface{} `json:"connection_details" structs:"connection_details" mapstructure:"connection_details"`
	MaxOpenConnections    int                    `json:"max_open_connections" structs:"max_open_connections" mapstructure:"max_open_connections"`
	MaxIdleConnections    int                    `json:"max_idle_connections" structs:"max_idle_connections" mapstructure:"max_idle_connections"`
	MaxConnectionLifetime time.Duration          `json:"max_connection_lifetime" structs:"max_connection_lifetime" mapstructure:"max_connection_lifetime"`
	PluginCommand         string                 `json:"plugin_command" structs:"plugin_command" mapstructure:"plugin_command"`
}

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

package dbs

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

const (
	postgreSQLTypeName = "postgres"
	cassandraTypeName  = "cassandra"
)

var (
	ErrUnsupportedDatabaseType = errors.New("Unsupported database type")
)

func Factory(conf ConnectionConfig) (DatabaseType, error) {
	switch conf.ConnectionType {
	case postgreSQLTypeName:
		return &PostgreSQL{
			config: conf,
		}, nil
	}

	return nil, ErrUnsupportedDatabaseType
}

type DatabaseType interface {
	Type() string
	Connection() (*sql.DB, error)
	Close()
	Reset(ConnectionConfig) (*sql.DB, error)
	CreateUser(createStmt, username, password, expiration string) error
	RenewUser(username, expiration string) error
	CustomRevokeUser(username, revocationSQL string) error
	DefaultRevokeUser(username string) error
}

type ConnectionConfig struct {
	ConnectionType     string            `json:"type" structs:"type" mapstructure:"type"`
	ConnectionURL      string            `json:"connection_url" structs:"connection_url" mapstructure:"connection_url"`
	ConnectionDetails  map[string]string `json:"connection_details" structs:"connection_details" mapstructure:"connection_details"`
	MaxOpenConnections int               `json:"max_open_connections" structs:"max_open_connections" mapstructure:"max_open_connections"`
	MaxIdleConnections int               `json:"max_idle_connections" structs:"max_idle_connections" mapstructure:"max_idle_connections"`
}

// Query templates a query for us.
func queryHelper(tpl string, data map[string]string) string {
	for k, v := range data {
		tpl = strings.Replace(tpl, fmt.Sprintf("{{%s}}", k), v, -1)
	}

	return tpl
}

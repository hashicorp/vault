package database

import (
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const databaseConfigPath = "database/dbs/"

// DatabaseType is the interface that all database objects must implement.
type DatabaseType interface {
	Type() string
	CreateUser(statements Statements, username, password, expiration string) error
	RenewUser(statements Statements, username, expiration string) error
	RevokeUser(statements Statements, username string) error

	Initialize(map[string]interface{}) error
	Close() error

	GenerateUsername(displayName string) (string, error)
	GeneratePassword() (string, error)
	GenerateExpiration(ttl time.Duration) (string, error)
}

// DatabaseConfig is used by the Factory function to configure a DatabaseType
// object.
type DatabaseConfig struct {
	PluginName string `json:"plugin_name" structs:"plugin_name" mapstructure:"plugin_name"`
	// ConnectionDetails stores the database specific connection settings needed
	// by each database type.
	ConnectionDetails     map[string]interface{} `json:"connection_details" structs:"connection_details" mapstructure:"connection_details"`
	MaxOpenConnections    int                    `json:"max_open_connections" structs:"max_open_connections" mapstructure:"max_open_connections"`
	MaxIdleConnections    int                    `json:"max_idle_connections" structs:"max_idle_connections" mapstructure:"max_idle_connections"`
	MaxConnectionLifetime time.Duration          `json:"max_connection_lifetime" structs:"max_connection_lifetime" mapstructure:"max_connection_lifetime"`
}

// Statements set in role creation and passed into the database type's functions.
// TODO: Add a way of setting defaults here.
type Statements struct {
	CreationStatements   string `json:"creation_statments" mapstructure:"creation_statements" structs:"creation_statments"`
	RevocationStatements string `json:"revocation_statements" mapstructure:"revocation_statements" structs:"revocation_statements"`
	RollbackStatements   string `json:"rollback_statements" mapstructure:"rollback_statements" structs:"rollback_statements"`
	RenewStatements      string `json:"renew_statements" mapstructure:"renew_statements" structs:"renew_statements"`
}

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend(conf).Setup(conf)
}

func Backend(conf *logical.BackendConfig) *databaseBackend {
	var b databaseBackend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			Root: []string{
				"dbs/plugin/*",
			},
		},

		Paths: []*framework.Path{
			pathConfigurePluginConnection(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathRoleCreate(&b),
			pathResetConnection(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},

		Clean: b.closeAllDBs,

		Invalidate: b.invalidate,
	}

	b.logger = conf.Logger
	b.connections = make(map[string]DatabaseType)
	return &b
}

type databaseBackend struct {
	connections map[string]DatabaseType
	logger      log.Logger

	*framework.Backend
	sync.Mutex
}

// resetAllDBs closes all connections from all database types
func (b *databaseBackend) closeAllDBs() {
	b.Lock()
	defer b.Unlock()

	for _, db := range b.connections {
		db.Close()
	}
}

// This function is used to retrieve a database object either from the cached
// connection map or by using the database config in storage. The caller of this
// function needs to hold the backend's lock.
func (b *databaseBackend) getOrCreateDBObj(s logical.Storage, name string) (DatabaseType, error) {
	// if the object already is built and cached, return it
	db, ok := b.connections[name]
	if ok {
		return db, nil
	}

	entry, err := s.Get(fmt.Sprintf("dbs/%s", name))
	if err != nil {
		return nil, fmt.Errorf("failed to read connection configuration with name: %s", name)
	}
	if entry == nil {
		return nil, fmt.Errorf("failed to find entry for connection with name: %s", name)
	}

	var config DatabaseConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, err
	}

	db, err = PluginFactory(&config, b.System(), b.logger)
	if err != nil {
		return nil, err
	}

	err = db.Initialize(config.ConnectionDetails)
	if err != nil {
		return nil, err
	}

	b.connections[name] = db

	return db, nil
}

func (b *databaseBackend) Role(s logical.Storage, n string) (*roleEntry, error) {
	entry, err := s.Get("role/" + n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *databaseBackend) invalidate(key string) {
	b.Lock()
	defer b.Unlock()

	switch {
	case strings.HasPrefix(key, databaseConfigPath):
		name := strings.TrimPrefix(key, databaseConfigPath)
		b.clearConnection(name)
	}
}

// clearConnection closes the database connection and
// removes it from the b.connections map.
func (b *databaseBackend) clearConnection(name string) {
	db, ok := b.connections[name]
	if ok {
		db.Close()
		delete(b.connections, name)
	}
}

const backendHelp = `
The database backend supports using many different databases
as secret backends, including but not limited to:
cassandra, msslq, mysql, postgres

After mounting this backend, configure it using the endpoints within
the "database/dbs/" path.
`

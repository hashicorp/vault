package database

import (
	"context"
	"fmt"
	"net/rpc"
	"strings"
	"sync"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const databaseConfigPath = "database/config/"

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend(conf *logical.BackendConfig) *databaseBackend {
	var b databaseBackend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config/*",
			},
		},

		Paths: []*framework.Path{
			pathListPluginConnection(&b),
			pathConfigurePluginConnection(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
			pathResetConnection(&b),
		},

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},
		Clean:       b.closeAllDBs,
		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	b.logger = conf.Logger
	b.connections = make(map[string]dbplugin.Database)
	return &b
}

type databaseBackend struct {
	connections map[string]dbplugin.Database
	logger      log.Logger

	*framework.Backend
	sync.RWMutex
}

// closeAllDBs closes all connections from all database types
func (b *databaseBackend) closeAllDBs() {
	b.Lock()
	defer b.Unlock()

	for _, db := range b.connections {
		db.Close()
	}

	b.connections = make(map[string]dbplugin.Database)
}

// This function is used to retrieve a database object either from the cached
// connection map. The caller of this function needs to hold the backend's read
// lock.
func (b *databaseBackend) getDBObj(name string) (dbplugin.Database, bool) {
	db, ok := b.connections[name]
	return db, ok
}

// This function creates a new db object from the stored configuration and
// caches it in the connections map. The caller of this function needs to hold
// the backend's write lock
func (b *databaseBackend) createDBObj(ctx context.Context, s logical.Storage, name string) (dbplugin.Database, error) {
	db, ok := b.connections[name]
	if ok {
		return db, nil
	}

	config, err := b.DatabaseConfig(s, name)
	if err != nil {
		return nil, err
	}

	db, err = dbplugin.PluginFactory(config.PluginName, b.System(), b.logger)
	if err != nil {
		return nil, err
	}

	err = db.Initialize(ctx, config.ConnectionDetails, true)
	if err != nil {
		return nil, err
	}

	b.connections[name] = db

	return db, nil
}

func (b *databaseBackend) DatabaseConfig(s logical.Storage, name string) (*DatabaseConfig, error) {
	entry, err := s.Get(fmt.Sprintf("config/%s", name))
	if err != nil {
		return nil, fmt.Errorf("failed to read connection configuration: %s", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("failed to find entry for connection with name: %s", name)
	}

	var config DatabaseConfig
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

type upgradeStatements struct {
	// This json tag has a typo in it, the new version does not. This
	// necessitates this upgrade logic.
	CreationStatements   string `json:"creation_statments"`
	RevocationStatements string `json:"revocation_statements"`
	RollbackStatements   string `json:"rollback_statements"`
	RenewStatements      string `json:"renew_statements"`
}

type upgradeCheck struct {
	// This json tag has a typo in it, the new version does not. This
	// necessitates this upgrade logic.
	Statements upgradeStatements `json:"statments"`
}

func (b *databaseBackend) Role(s logical.Storage, roleName string) (*roleEntry, error) {
	entry, err := s.Get("role/" + roleName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var upgradeCh upgradeCheck
	if err := entry.DecodeJSON(&upgradeCh); err != nil {
		return nil, err
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	empty := upgradeCheck{}
	if upgradeCh != empty {
		result.Statements.CreationStatements = upgradeCh.Statements.CreationStatements
		result.Statements.RevocationStatements = upgradeCh.Statements.RevocationStatements
		result.Statements.RollbackStatements = upgradeCh.Statements.RollbackStatements
		result.Statements.RenewStatements = upgradeCh.Statements.RenewStatements
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

func (b *databaseBackend) closeIfShutdown(name string, err error) {
	// Plugin has shutdown, close it so next call can reconnect.
	switch err {
	case rpc.ErrShutdown, dbplugin.ErrPluginShutdown:
		b.Lock()
		b.clearConnection(name)
		b.Unlock()
	}
}

const backendHelp = `
The database backend supports using many different databases
as secret backends, including but not limited to:
cassandra, mssql, mysql, postgres

After mounting this backend, configure it using the endpoints within
the "database/config/" path.
`

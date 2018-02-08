package database

import (
	"context"
	"fmt"
	"net/rpc"
	"strings"
	"sync"

	log "github.com/mgutz/logxi/v1"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const databaseConfigPath = "database/config/"

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	if err := b.Setup(ctx, conf); err != nil {
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

func (b *databaseBackend) DatabaseConfig(ctx context.Context, s logical.Storage, name string) (*DatabaseConfig, error) {
	entry, err := s.Get(ctx, fmt.Sprintf("config/%s", name))
	if err != nil {
		return nil, errwrap.Wrapf("failed to read connection configuration: {{err}}", err)
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

type oldStatements struct {
	CreationStatements   string `json:"creation_statements"`
	RevocationStatements string `json:"revocation_statements"`
	RollbackStatements   string `json:"rollback_statements"`
	RenewStatements      string `json:"renew_statements"`
}

type upgradeCheck struct {
	// This json tag has a typo in it, the new version does not. This
	// necessitates this upgrade logic.
	Statements    *upgradeStatements `json:"statments,omitempty"`
	OldStatements *oldStatements     `json:"statements,omitempty"`
}

func (b *databaseBackend) Role(ctx context.Context, s logical.Storage, roleName string) (*roleEntry, error) {
	entry, err := s.Get(ctx, "role/"+roleName)
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

	switch {
	case upgradeCh.Statements != nil:
		var stmts dbplugin.Statements
		if upgradeCh.Statements.CreationStatements != "" {
			stmts.CreationStatements = []string{upgradeCh.Statements.CreationStatements}
		}
		if upgradeCh.Statements.RevocationStatements != "" {
			stmts.RevocationStatements = []string{upgradeCh.Statements.RevocationStatements}
		}
		if upgradeCh.Statements.RollbackStatements != "" {
			stmts.RollbackStatements = []string{upgradeCh.Statements.RollbackStatements}
		}
		if upgradeCh.Statements.RenewStatements != "" {
			stmts.RenewStatements = []string{upgradeCh.Statements.RenewStatements}
		}
		result.Statements = stmts
	case upgradeCh.OldStatements != nil:
		var stmts dbplugin.Statements
		if upgradeCh.OldStatements.CreationStatements != "" {
			stmts.CreationStatements = []string{upgradeCh.OldStatements.CreationStatements}
		}
		if upgradeCh.OldStatements.RevocationStatements != "" {
			stmts.RevocationStatements = []string{upgradeCh.OldStatements.RevocationStatements}
		}
		if upgradeCh.OldStatements.RollbackStatements != "" {
			stmts.RollbackStatements = []string{upgradeCh.OldStatements.RollbackStatements}
		}
		if upgradeCh.OldStatements.RenewStatements != "" {
			stmts.RenewStatements = []string{upgradeCh.OldStatements.RenewStatements}
		}
		result.Statements = stmts
	}

	return &result, nil
}

func (b *databaseBackend) invalidate(ctx context.Context, key string) {
	switch {
	case strings.HasPrefix(key, databaseConfigPath):
		name := strings.TrimPrefix(key, databaseConfigPath)
		b.ClearConnection(name)
	}
}

func (b *databaseBackend) GetConnection(ctx context.Context, s logical.Storage, name string) (dbplugin.Database, error) {
	b.RLock()
	unlockFunc := b.RUnlock
	defer func() { unlockFunc() }()

	db, ok := b.connections[name]
	if ok {
		return db, nil
	}

	// Upgrade lock
	b.RUnlock()
	b.Lock()
	unlockFunc = b.Unlock

	db, ok = b.connections[name]
	if ok {
		return db, nil
	}

	config, err := b.DatabaseConfig(ctx, s, name)
	if err != nil {
		return nil, err
	}

	db, err = dbplugin.PluginFactory(ctx, config.PluginName, b.System(), b.logger)
	if err != nil {
		return nil, err
	}

	_, err = db.Initialize(ctx, config.ConnectionDetails, true)
	if err != nil {
		db.Close()
		return nil, err
	}
	b.connections[name] = db

	return db, nil
}

// clearConnection closes the database connection and
// removes it from the b.connections map.
func (b *databaseBackend) ClearConnection(name string) error {
	b.Lock()
	defer b.Unlock()

	db, ok := b.connections[name]
	if ok {
		if err := db.Close(); err != nil {
			return err
		}
		delete(b.connections, name)
	}
	return nil
}

func (b *databaseBackend) CloseIfShutdown(name string, err error) {
	// Plugin has shutdown, close it so next call can reconnect.
	switch err {
	case rpc.ErrShutdown, dbplugin.ErrPluginShutdown:
		b.ClearConnection(name)
	}
}

// closeAllDBs closes all connections from all database types
func (b *databaseBackend) closeAllDBs(ctx context.Context) {
	b.Lock()
	defer b.Unlock()

	for _, db := range b.connections {
		db.Close()
	}
	b.connections = make(map[string]dbplugin.Database)
}

const backendHelp = `
The database backend supports using many different databases
as secret backends, including but not limited to:
cassandra, mssql, mysql, postgres

After mounting this backend, configure it using the endpoints within
the "database/config/" path.
`

package database

import (
	"context"
	"fmt"
	"net/rpc"
	"strings"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-uuid"
	v4 "github.com/hashicorp/vault/sdk/database/dbplugin"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	databaseConfigPath     = "config/"
	databaseRolePath       = "role/"
	databaseStaticRolePath = "static-role/"
	minRootCredRollbackAge = 1 * time.Minute
)

type dbPluginInstance struct {
	sync.RWMutex
	database databaseVersionWrapper

	id     string
	name   string
	closed bool
}

func (dbi *dbPluginInstance) Close() error {
	dbi.Lock()
	defer dbi.Unlock()

	if dbi.closed {
		return nil
	}
	dbi.closed = true

	return dbi.database.Close()
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	b.credRotationQueue = queue.New()
	// Create a context with a cancel method for processing any WAL entries and
	// populating the queue
	initCtx := context.Background()
	b.ctx, b.cancelQueue = context.WithCancel(initCtx)
	// Load queue and kickoff new periodic ticker
	go b.initQueue(b.ctx, conf, conf.System.ReplicationState())
	return b, nil
}

func Backend(conf *logical.BackendConfig) *databaseBackend {
	var b databaseBackend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				"config/*",
				"static-role/*",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				pathListPluginConnection(&b),
				pathConfigurePluginConnection(&b),
				pathResetConnection(&b),
			},
			pathListRoles(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
			pathRotateRootCredentials(&b),
		),

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},
		Clean:             b.clean,
		Invalidate:        b.invalidate,
		WALRollback:       b.walRollback,
		WALRollbackMinAge: minRootCredRollbackAge,
		BackendType:       logical.TypeLogical,
	}

	b.logger = conf.Logger
	b.connections = make(map[string]*dbPluginInstance)

	b.roleLocks = locksutil.CreateLocks()

	return &b
}

type databaseBackend struct {
	// used to synchronize access to the connections map
	connLock sync.RWMutex
	// connections holds configured database connections by config name
	connections map[string]*dbPluginInstance
	logger      log.Logger

	*framework.Backend
	// credRotationQueue is an in-memory priority queue used to track Static Roles
	// that require periodic rotation. Backends will have a PriorityQueue
	// initialized on setup, but only backends that are mounted by a primary
	// server or mounted as a local mount will perform the rotations.
	//
	// cancelQueue is used to remove the priority queue and terminate the
	// background ticker.
	credRotationQueue *queue.PriorityQueue
	// context used for canceling operations
	ctx         context.Context
	cancelQueue context.CancelFunc

	// roleLocks is used to lock modifications to roles in the queue, to ensure
	// concurrent requests are not modifying the same role and possibly causing
	// issues with the priority queue.
	roleLocks []*locksutil.LockEntry
}

func (b *databaseBackend) connGet(name string) *dbPluginInstance {
	b.connLock.RLock()
	defer b.connLock.RUnlock()
	return b.connections[name]
}

func (b *databaseBackend) connPop(name string) *dbPluginInstance {
	b.connLock.Lock()
	defer b.connLock.Unlock()
	dbi := b.connections[name]
	delete(b.connections, name)
	return dbi
}

func (b *databaseBackend) connPut(name string, newDbi *dbPluginInstance) *dbPluginInstance {
	b.connLock.Lock()
	defer b.connLock.Unlock()
	dbi := b.connections[name]
	b.connections[name] = newDbi
	return dbi
}

func (b *databaseBackend) connClear() map[string]*dbPluginInstance {
	b.connLock.Lock()
	defer b.connLock.Unlock()
	old := b.connections
	b.connections = make(map[string]*dbPluginInstance)
	return old
}

func (b *databaseBackend) DatabaseConfig(ctx context.Context, s logical.Storage, name string) (*DatabaseConfig, error) {
	entry, err := s.Get(ctx, fmt.Sprintf("config/%s", name))
	if err != nil {
		return nil, fmt.Errorf("failed to read connection configuration: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("failed to find entry for connection with name: %q", name)
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
	Statements *upgradeStatements `json:"statments,omitempty"`
}

func (b *databaseBackend) Role(ctx context.Context, s logical.Storage, roleName string) (*roleEntry, error) {
	return b.roleAtPath(ctx, s, roleName, databaseRolePath)
}

func (b *databaseBackend) StaticRole(ctx context.Context, s logical.Storage, roleName string) (*roleEntry, error) {
	return b.roleAtPath(ctx, s, roleName, databaseStaticRolePath)
}

func (b *databaseBackend) roleAtPath(ctx context.Context, s logical.Storage, roleName string, pathPrefix string) (*roleEntry, error) {
	entry, err := s.Get(ctx, pathPrefix+roleName)
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
		var stmts v4.Statements
		if upgradeCh.Statements.CreationStatements != "" {
			stmts.Creation = []string{upgradeCh.Statements.CreationStatements}
		}
		if upgradeCh.Statements.RevocationStatements != "" {
			stmts.Revocation = []string{upgradeCh.Statements.RevocationStatements}
		}
		if upgradeCh.Statements.RollbackStatements != "" {
			stmts.Rollback = []string{upgradeCh.Statements.RollbackStatements}
		}
		if upgradeCh.Statements.RenewStatements != "" {
			stmts.Renewal = []string{upgradeCh.Statements.RenewStatements}
		}
		result.Statements = stmts
	}

	result.Statements.Revocation = strutil.RemoveEmpty(result.Statements.Revocation)

	// For backwards compatibility, copy the values back into the string form
	// of the fields
	result.Statements = dbutil.StatementCompatibilityHelper(result.Statements)

	return &result, nil
}

func (b *databaseBackend) invalidate(ctx context.Context, key string) {
	switch {
	case strings.HasPrefix(key, databaseConfigPath):
		name := strings.TrimPrefix(key, databaseConfigPath)
		b.ClearConnection(name)
	}
}

func (b *databaseBackend) GetConnection(ctx context.Context, s logical.Storage, name string) (*dbPluginInstance, error) {
	config, err := b.DatabaseConfig(ctx, s, name)
	if err != nil {
		return nil, err
	}

	return b.GetConnectionWithConfig(ctx, name, config)
}

func (b *databaseBackend) GetConnectionWithConfig(ctx context.Context, name string, config *DatabaseConfig) (*dbPluginInstance, error) {
	dbi := b.connGet(name)
	if dbi != nil {
		return dbi, nil
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	dbw, err := newDatabaseWrapper(ctx, config.PluginName, b.System(), b.logger)
	if err != nil {
		return nil, fmt.Errorf("unable to create database instance: %w", err)
	}

	initReq := v5.InitializeRequest{
		Config:           config.ConnectionDetails,
		VerifyConnection: true,
	}
	_, err = dbw.Initialize(ctx, initReq)
	if err != nil {
		dbw.Close()
		return nil, err
	}

	dbi = &dbPluginInstance{
		database: dbw,
		id:       id,
		name:     name,
	}
	oldConn := b.connPut(name, dbi)
	if oldConn != nil {
		err := oldConn.Close()
		if err != nil {
			b.Logger().Warn("Error closing database connection", "error", err)
		}
	}
	return dbi, nil
}

// ClearConnection closes the database connection and
// removes it from the b.connections map.
func (b *databaseBackend) ClearConnection(name string) error {
	db := b.connPop(name)
	if db != nil {
		// Ignore error here since the database client is always killed
		db.Close()
	}
	return nil
}

func (b *databaseBackend) CloseIfShutdown(db *dbPluginInstance, err error) {
	// Plugin has shutdown, close it so next call can reconnect.
	switch err {
	case rpc.ErrShutdown, v4.ErrPluginShutdown, v5.ErrPluginShutdown:
		// Put this in a goroutine so that requests can run with the read or write lock
		// and simply defer the unlock.  Since we are attaching the instance and matching
		// the id in the connection map, we can safely do this.
		go func() {
			db.Close()

			// Ensure we are deleting the correct connection
			mapDB := b.connPop(db.name)
			if mapDB != nil && db.id != mapDB.id {
				// oops, put it back
				oldDbi := b.connPut(db.name, mapDB)
				if oldDbi != nil {
					// there is a small chance that something else was inserted in that slot during that time
					// if so, clean it up
					oldDbi.Close()
				}
			}
		}()
	}
}

// clean closes all connections from all database types
// and cancels any rotation queue loading operation.
func (b *databaseBackend) clean(_ context.Context) {
	// kill the queue and terminate the background ticker
	b.cancelQueue()

	connections := b.connClear()
	for _, db := range connections {
		go db.Close()
	}
}

const backendHelp = `
The database backend supports using many different databases
as secret backends, including but not limited to:
cassandra, mssql, mysql, postgres

After mounting this backend, configure it using the endpoints within
the "database/config/" path.
`

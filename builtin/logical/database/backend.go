package database

import (
	"context"
	"fmt"
	"net/rpc"
	"strings"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/errwrap"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	databaseConfigPath     = "database/config/"
	databaseRolePath       = "role/"
	databaseStaticRolePath = "static-role/"

	// interval to check the queue for items needing rotation
	queueTickSeconds  = 5
	queueTickInterval = queueTickSeconds * time.Second

	// wal storage key used for static account rotations
	staticWALKey = "staticRotationKey"
)

type dbPluginInstance struct {
	sync.RWMutex
	dbplugin.Database

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

	return dbi.Database.Close()
}

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend(conf)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}

	b.credRotationQueue = queue.New()
	// create a context with a cancel method for processing any WAL entries and
	// populating the queue
	ctx, cancel := context.WithCancel(ctx)
	b.cancelQueue = cancel
	// load queue and kickoff new periodic ticker
	go b.initQueue(ctx, conf)
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
			pathRotateCredentials(&b),
		),

		Secrets: []*framework.Secret{
			secretCreds(&b),
		},
		Clean:       b.clean,
		Invalidate:  b.invalidate,
		BackendType: logical.TypeLogical,
	}

	b.logger = conf.Logger
	b.connections = make(map[string]*dbPluginInstance)

	return &b
}

type databaseBackend struct {
	connections map[string]*dbPluginInstance
	logger      log.Logger

	*framework.Backend
	sync.RWMutex
	// credRotationQueue is an in-memory priority queue used to track Roles that
	// are associated with static accounts and require periodic rotation. Only
	// backends that are mounted by a primary server, or mounted as a local mount,
	// will have a priority queue and perform the rotations.
	//
	// cancelQueue is used to remove the priority queue and terminate the
	// background ticker.
	credRotationQueue *queue.PriorityQueue
	cancelQueue       context.CancelFunc
}

func (b *databaseBackend) DatabaseConfig(ctx context.Context, s logical.Storage, name string) (*DatabaseConfig, error) {
	entry, err := s.Get(ctx, fmt.Sprintf("config/%s", name))
	if err != nil {
		return nil, errwrap.Wrapf("failed to read connection configuration: {{err}}", err)
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
		var stmts dbplugin.Statements
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

	dbp, err := dbplugin.PluginFactory(ctx, config.PluginName, b.System(), b.logger)
	if err != nil {
		return nil, err
	}

	_, err = dbp.Init(ctx, config.ConnectionDetails, true)
	if err != nil {
		dbp.Close()
		return nil, err
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	db = &dbPluginInstance{
		Database: dbp,
		name:     name,
		id:       id,
	}

	b.connections[name] = db
	return db, nil
}

// initQueue preforms the necessary checks and initializations needed to preform
// automatic credential rotation for roles associated with static accounts. This
// method verifies if a queue is needed (primary server or local mount), and if
// so initializes the queue and launches a go-routine to periodically invoke a
// method to preform the rotations.
//
// initQueue is invoked by the Factory method in a go-routine. The Factory does
// not wait for success or failure of it's tasks before continuing. This is to
// avoid blocking the mount process while loading and evaluating existing roles,
// etc.
func (b *databaseBackend) initQueue(ctx context.Context, conf *logical.BackendConfig) {
	// verify this mount is on the primary server, or is a local mount. If not, do
	// not create a queue or launch a ticker. Both processing the WAL list and
	// populating the queue are done sequentially and before launching a
	// go-routine to run the periodic ticker.
	replicationState := conf.System.ReplicationState()
	if (conf.System.LocalMount() || !replicationState.HasState(consts.ReplicationPerformanceSecondary)) &&
		!replicationState.HasState(consts.ReplicationDRSecondary) &&
		!replicationState.HasState(consts.ReplicationPerformanceStandby) {
		b.Logger().Info("initializing database rotation queue")

		// Sleep a few seconds to allow Vault to mount and complete setup, so
		// that we can write to storage
		time.Sleep(3 * time.Second)

		// load roles and populate queue with static accounts
		b.populateQueue(ctx, conf.StorageView)

		// launch ticker
		go b.runTicker(ctx, conf.StorageView)
	}
}

// loadStaticWALs reads WAL entries and returns a map of roles and their
// setCredentialsWAL, if found.
func (b *databaseBackend) loadStaticWALs(ctx context.Context, s logical.Storage) (map[string]*setCredentialsWAL, error) {
	keys, err := framework.ListWAL(ctx, s)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		b.Logger().Debug("no WAL entries found")
		return nil, nil
	}

	walMap := make(map[string]*setCredentialsWAL)
	// loop through WAL keys and process any rotation ones
	for _, walID := range keys {
		walEntry := b.findStaticWAL(ctx, s, walID)
		if walEntry == nil {
			continue
		}

		// verify the static role still exists
		roleName := walEntry.RoleName
		role, err := b.StaticRole(ctx, s, roleName)
		if err != nil {
			b.Logger().Warn("unable to read static role", "error", err, "role", roleName)
			continue
		}
		if role == nil || role.StaticAccount == nil {
			if err := framework.DeleteWAL(ctx, s, walEntry.walID); err != nil {
				b.Logger().Warn("unable to delete WAL", "error", err, "WAL ID", walEntry.walID)
			}
			continue
		}

		walEntry.walID = walID
		walMap[walEntry.RoleName] = walEntry
	}
	return walMap, nil
}

// invalidateQueue cancels any background queue loading and destroys the queue.
func (b *databaseBackend) invalidateQueue() {
	b.Lock()
	defer b.Unlock()

	// cancelQueue
	if b.cancelQueue != nil {
		b.cancelQueue()
	}
	b.credRotationQueue = nil
}

// ClearConnection closes the database connection and
// removes it from the b.connections map.
func (b *databaseBackend) ClearConnection(name string) error {
	b.Lock()
	defer b.Unlock()
	return b.clearConnection(name)
}

func (b *databaseBackend) clearConnection(name string) error {
	db, ok := b.connections[name]
	if ok {
		// Ignore error here since the database client is always killed
		db.Close()
		delete(b.connections, name)
	}
	return nil
}

func (b *databaseBackend) CloseIfShutdown(db *dbPluginInstance, err error) {
	// Plugin has shutdown, close it so next call can reconnect.
	switch err {
	case rpc.ErrShutdown, dbplugin.ErrPluginShutdown:
		// Put this in a goroutine so that requests can run with the read or write lock
		// and simply defer the unlock.  Since we are attaching the instance and matching
		// the id in the connection map, we can safely do this.
		go func() {
			b.Lock()
			defer b.Unlock()
			db.Close()

			// Ensure we are deleting the correct connection
			mapDB, ok := b.connections[db.name]
			if ok && db.id == mapDB.id {
				delete(b.connections, db.name)
			}
		}()
	}
}

// clean closes all connections from all database types
// and cancels any rotation queue loading operation.
func (b *databaseBackend) clean(ctx context.Context) {
	// invalidateQueue acquires it's own lock on the backend, removes queue, and
	// terminates the background ticker
	b.invalidateQueue()

	b.Lock()
	defer b.Unlock()

	for _, db := range b.connections {
		db.Close()
	}
	b.connections = make(map[string]*dbPluginInstance)
}

// pushItem wraps the internal queue's Push call, to make sure a queue is
// actually available. This is needed because both runTicker and initQueue
// operate in go-routines, and could be accessing the queue concurrently
func (b *databaseBackend) pushItem(item *queue.Item) error {
	b.RLock()
	unlockFunc := b.RUnlock
	defer func() { unlockFunc() }()

	if b.credRotationQueue != nil {
		return b.credRotationQueue.Push(item)
	}

	b.Logger().Warn("no queue found during push item")
	return nil
}

// popItem wraps the internal queue's Pop call, to make sure a queue is
// actually available. This is needed because both runTicker and initQueue
// operate in go-routines, and could be accessing the queue concurrently
func (b *databaseBackend) popItem() (*queue.Item, error) {
	b.RLock()
	defer b.RUnlock()
	if b.credRotationQueue != nil {
		return b.credRotationQueue.Pop()
	}
	return nil, queue.ErrEmpty
}

// popItemByKey wraps the internal queue's PopByKey call, to make sure a queue is
// actually available. This is needed because both runTicker and initQueue
// operate in go-routines, and could be accessing the queue concurrently
func (b *databaseBackend) popItemByKey(name string) (*queue.Item, error) {
	b.RLock()
	defer b.RUnlock()
	if b.credRotationQueue != nil {
		return b.credRotationQueue.PopByKey(name)
	}
	return nil, queue.ErrEmpty
}

const backendHelp = `
The database backend supports using many different databases
as secret backends, including but not limited to:
cassandra, mssql, mysql, postgres

After mounting this backend, configure it using the endpoints within
the "database/config/" path.
`

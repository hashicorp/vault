package database

import (
	"context"
	"fmt"
	"net/rpc"
	"strings"
	"sync"
        "time"

        "github.com/hashicorp/errwrap"
        log "github.com/hashicorp/go-hclog"
        "github.com/hashicorp/go-uuid"
        "github.com/hashicorp/vault/builtin/logical/database/dbplugin"
        "github.com/hashicorp/vault/helper/consts"
        "github.com/hashicorp/vault/helper/queue"
        "github.com/hashicorp/vault/helper/strutil"
        "github.com/hashicorp/vault/logical"
        "github.com/hashicorp/vault/logical/framework"
        "github.com/hashicorp/vault/plugins/helper/database/dbutil"
)

const databaseConfigPath = "database/config/"
const databaseRolePath = "role/"

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
                                // TODO: will want to encrypt static accounts / roles with password info
                                // in them
			},
		},
		Paths: []*framework.Path{
			pathListPluginConnection(&b),
			pathConfigurePluginConnection(&b),
			pathListRoles(&b),
			pathRoles(&b),
			pathCredsCreate(&b),
			pathResetConnection(&b),
			pathRotateCredentials(&b),
                        pathRotateRoleCredentials(&b),
		},

                Secrets: []*framework.Secret{
                        secretCreds(&b),
                },
                Clean:             b.clean,
                Invalidate:        b.invalidate,
                BackendType:       logical.TypeLogical,
                WALRollback:       b.walRollback,
                WALRollbackMinAge: 2 * time.Minute,
        }

        b.logger = conf.Logger
        b.connections = make(map[string]*dbPluginInstance)

        replicationState := conf.System.ReplicationState()
        if (conf.System.LocalMount() || !replicationState.HasState(consts.ReplicationPerformanceSecondary)) &&
                !replicationState.HasState(consts.ReplicationDRSecondary) &&
                !replicationState.HasState(consts.ReplicationPerformanceStandby) {

                b.initQueue(conf.StorageView)

                b.Backend.PeriodicFunc = func(ctx context.Context, req *logical.Request) error {
                        // Re-init queue if it was invalidated
                        b.Lock()
                        if b.credRotationQueue == nil {
                                b.initQueue(req.Storage)
                        }
                        b.Unlock()

                        // This will loop until either:
                        // - The queue of passwords needing rotation is completely empty.
                        // - It encounters the first password not yet needing rotation.

                        // reQueue is a collection of any queue items that failed to rotate, to be
                        // re-added at the end
                        var reQueue []*queue.Item
                        for {
                                item, err := b.credRotationQueue.PopItem()
                                if err != nil {
                                        if err == queue.ErrEmpty {
                                                return nil
                                        }
                                        return err
                                }

                                // TODO: load role and don't just use the value
                                role := item.Value.(*roleEntry)

                                if time.Now().Unix() > item.Priority {
                                        // We've found our first item not in need of rotation
                                        if err := b.createUpdateStaticAccount(ctx, req.Storage, item.Key, role, false); err != nil {
                                                b.logger.Warn("unable rotate credentials in periodic function", "error", err)
                                                // add the item to the re-queue slice
                                                reQueue = append(reQueue, item)
                                        }
                                        b.credRotationQueue.PushItem(&queue.Item{
                                                Key:      item.Key,
                                                Value:    role,
                                                Priority: time.Now().Add(role.StaticAccount.RotationPeriod).Unix(),
                                        })
                                } else {
                                        b.credRotationQueue.PushItem(item)
                                        return nil
                                }
                        }

                        // re-enqueue items that failed to rotate
                        for _, item := range reQueue {
                                b.credRotationQueue.PushItem(item)
                        }
                        return nil
                }
        }

        return &b
}

type databaseBackend struct {
	connections map[string]*dbPluginInstance
	logger      log.Logger

        *framework.Backend
        sync.RWMutex
        credRotationQueue queue.PriorityQueue
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
        case strings.HasPrefix(key, databaseRolePath):
                b.invalidateQueue()
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

// initQueue creates a new priority queue and starts a cancellable goroutine
// to load existing static roles.
//
// This method is should not be called concurrently.
func (b *databaseBackend) initQueue(s logical.Storage) {
        b.credRotationQueue = queue.NewTimeQueue()
        ctx, cancel := context.WithCancel(context.Background())
        b.cancelQueue = cancel
        go b.populateQueue(ctx, s)
}

// invalidateQueue cancels any background queue loading and destroys the queue.
func (b *databaseBackend) invalidateQueue() {
        b.Lock()
        defer b.Unlock()

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
        // invalidateQueue acquires it's own lock
        b.invalidateQueue()

        b.Lock()
        defer b.Unlock()

        for _, db := range b.connections {
                db.Close()
        }
	b.connections = make(map[string]*dbPluginInstance)
}

const backendHelp = `
The database backend supports using many different databases
as secret backends, including but not limited to:
cassandra, mssql, mysql, postgres

After mounting this backend, configure it using the endpoints within
the "database/config/" path.
`

// populate queue loads the priority queue with existing static accounts.
func (b *databaseBackend) populateQueue(ctx context.Context, s logical.Storage) {
        log := b.Logger()

        log.Info("restoring role rotation queue")

        roles, err := s.List(ctx, "role/")
        if err != nil {
                log.Warn("unable to list role for enqueueing", "error", err)
                return
        }

        for _, roleName := range roles {
                select {
                case <-ctx.Done():
                        log.Info("rotation queue restore cancelled")
                        return
                default:
                }

                role, err := b.Role(ctx, s, roleName)
                if err != nil {
                        log.Warn("unable to read role", "error", err, "role", roleName)
                        continue
                }
                if role == nil || role.StaticAccount == nil {
                        continue
                }

                if err := b.credRotationQueue.PushItem(&queue.Item{
                        Key:      roleName,
                        Value:    role,
                        Priority: role.StaticAccount.LastVaultRotation.Add(role.StaticAccount.RotationPeriod).Unix(),
                }); err != nil {
                        log.Warn("unable to enqueue item", "error", err, "role", roleName)
                }
        }
        log.Info("successfully restored role rotation queue")
}

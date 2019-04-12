package database

import (
	"context"
        "errors"
	"fmt"
	"net/rpc"
	"strings"
	"sync"
        "time"

        "github.com/hashicorp/errwrap"
        log "github.com/hashicorp/go-hclog"
        "github.com/hashicorp/go-multierror"
        "github.com/hashicorp/go-uuid"
        "github.com/hashicorp/vault/builtin/logical/database/dbplugin"
        "github.com/hashicorp/vault/helper/consts"
        "github.com/hashicorp/vault/helper/queue"
        "github.com/hashicorp/vault/helper/strutil"
        "github.com/hashicorp/vault/logical"
        "github.com/hashicorp/vault/logical/framework"
        "github.com/hashicorp/vault/plugins/helper/database/dbutil"
        "github.com/mitchellh/mapstructure"
)

const (
        databaseConfigPath = "database/config/"
        databaseRolePath   = "role/"

        // interval to check the queue for items needing rotation
        QueueTickInterval = 5 * time.Second

        // key used for WAL entry kind information
        walRotationKey = "staticRotationKey"
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
        credRotationQueue queue.PriorityQueue

        // cancelQueue is used to remove the priority queue and terminate the
        // background ticker
        cancelQueue context.CancelFunc
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

// initQueue preforms the necessary checks and initializations needed to preform
// automatic credential rotation for roles associated with static accounts. This
// method verifies if a queue is needed (primary, non-local mount), and if so
// initializes the queue and launches a go-routine to periodically invoke a
// method to preform the rotations.
//
// initQueue is invoked by the Factory method in a go-routine. The Factory does
// not wait for success or failure of it's tasks before continuing. This is to
// avoid blocking the mount process while loading and evaluating existing roles,
// etc.
func (b *databaseBackend) initQueue(ctx context.Context, conf *logical.BackendConfig) {
        // verify this mount is on the primary server, or is a local mount. If not, do
        // not create a queue or launch a ticker
        replicationState := conf.System.ReplicationState()
        if (conf.System.LocalMount() || !replicationState.HasState(consts.ReplicationPerformanceSecondary)) &&
                !replicationState.HasState(consts.ReplicationDRSecondary) &&
                !replicationState.HasState(consts.ReplicationPerformanceStandby) {
                b.Logger().Info("initializing database rotation queue")

                b.Lock()
                if b.credRotationQueue == nil {
                        b.credRotationQueue = queue.NewTimeQueue()
                }
                b.Unlock()

                // search through WAL for any rotations that were interrupted
                if err := b.checkQueueWAL(ctx, conf); err != nil {
                        b.Logger().Warn("error(s) loading WAL entries into queue: ", err)
                }

                // load roles and populate queue with static accounts
                ctx, cancel := context.WithCancel(context.Background())
                b.cancelQueue = cancel
                b.populateQueue(ctx, conf.StorageView)
                // launch ticker
                go b.runTicker(ctx, conf.StorageView)
        }
}

// checkQueueWAL reads WAL entries at backend initialization. If WAL entries for
// static account rotation are round, attempt to re-set the password for the
// role given the NewPassword stored in the WAL. If the matching Role does not
// exist, the Role's LastVaultRotation is newer than the WAL, or the Role does
// not have a static account, delete the WAL.
func (b *databaseBackend) checkQueueWAL(ctx context.Context, conf *logical.BackendConfig) error {
        keys, err := framework.ListWAL(ctx, conf.StorageView)
        if err != nil {
                return err
        }
        if len(keys) == 0 {
                b.Logger().Info("no WAL entries found")
                return nil
        }

        // loop through WAL keys and process any rotation ones
        var merr error
        for _, walID := range keys {
                select {
                case <-ctx.Done():
                        b.Logger().Info("checkQueueWAL cancelled")
                        return merr
                default:
                }
                walEntry := b.walForItemValue(ctx, conf.StorageView, walID)
                if walEntry == nil {
                        continue
                }

                // TODO: just use createUpdateStaticAccount
                // load matching role and verify
                role, err := b.Role(ctx, conf.StorageView, walEntry.RoleName)
                if err != nil {
                        b.Logger().Warn("error loading role", err)
                        merr = multierror.Append(merr, err)
                        continue
                }

                if role == nil || role.StaticAccount == nil {
                        b.Logger().Warn("role or static account not found")
                        if err = framework.DeleteWAL(ctx, conf.StorageView, walID); err != nil {
                                b.Logger().Warn("error deleting WAL for role with no static account", err)
                                merr = multierror.Append(merr, err)
                        }
                        continue
                }

                if role.StaticAccount.LastVaultRotation.After(walEntry.LastVaultRotation) {
                        // role password has been rotated since the WAL was created, so let's
                        // delete this
                        if err = framework.DeleteWAL(ctx, conf.StorageView, walID); err != nil {
                                b.Logger().Warn("error deleting WAL for role with newer rotation date", err)
                                merr = multierror.Append(merr, err)
                        }
                        continue
                }

                // createUpdateStaticAccount which will attempt to set the password and
                // delete the WAL if successful
                resp, err := b.createUpdateStaticAccount(ctx, conf.StorageView, &setPasswordInput{
                        RoleName: walEntry.RoleName,
                        Role:     role,
                        WALID:    walID,
                        Password: walEntry.NewPassword,
                })
                if err != nil {
                        merr = multierror.Append(merr, err)
                        // if response contains a WALID, create an item to push to the queue with
                        // a backoff time and include the WAL ID
                        merr = multierror.Append(merr, err)
                        if resp.WALID != "" {
                                // Add their rotation to the queue
                                if err := b.credRotationQueue.PushItem(&queue.Item{
                                        Key:      walEntry.RoleName,
                                        Value:    walID,
                                        Priority: walEntry.LastVaultRotation.Add(time.Second * 60).Unix(),
                                }); err != nil {
                                        b.Logger().Warn("error pushing item on to queue after failed WAL restore", err)
                                        merr = multierror.Append(merr, err)
                                }
                        }
                }
        } // end range keys
        return merr
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

const backendHelp = `
The database backend supports using many different databases
as secret backends, including but not limited to:
cassandra, mssql, mysql, postgres

After mounting this backend, configure it using the endpoints within
the "database/config/" path.
`

// populateQueue loads the priority queue with existing static accounts. This
// occurs at initialization, after any WAL entries of failed or interrupted
// rotations have been processed. It lists the roles from storage and searches
// for any that have an associated static account, then adds them to the
// priority queue for rotations.
func (b *databaseBackend) populateQueue(ctx context.Context, s logical.Storage) {
        log := b.Logger()
        log.Info("populating role rotation queue")

        roles, err := s.List(ctx, "role/")
        if err != nil {
                log.Warn("unable to list role for enqueueing", "error", err)
                return
        }

        // guard against a nil queue
        b.Lock()
        if b.credRotationQueue == nil {
                b.credRotationQueue = queue.NewTimeQueue()
        }
        b.Unlock()

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
                        Priority: role.StaticAccount.LastVaultRotation.Add(role.StaticAccount.RotationPeriod).Unix(),
                }); err != nil {
                        log.Warn("unable to enqueue item", "error", err, "role", roleName)
                }
        }
}

// runTicker kicks off a periodic ticker that invoke the automatic credential
// rotation method at a determined interval
func (b *databaseBackend) runTicker(ctx context.Context, s logical.Storage) {
        b.logger.Info("starting periodic ticker")
        tick := time.NewTicker(QueueTickInterval)
        defer tick.Stop()
        for {
                select {
                case <-tick.C:
                        b.rotateCredentials(ctx, s)

                case <-ctx.Done():
                        b.logger.Info("stopping periodic ticker")
                        return
                }
        }
        return
}

// walSetCredentials is used to store information in a WAL that can retry a
// credential setting or rotation in the event of partial failure.
type walSetCredentials struct {
        // TODO: need this?
        ID       string
        Attempts int
        //

        Username          string
        NewPassword       string
        OldPassword       string
        RoleName          string
        Statements        []string
        LastVaultRotation time.Time
}

// rotateCredentials sets a new password for a static account. This method is
// invoked by a go-routine launched the runTicker method, and invoked
// periodically (approximately every 5 seconds).
// This will loop until either:
// - The queue of passwords needing rotation is completely empty.
// - It encounters the first password not yet needing rotation.
func (b *databaseBackend) rotateCredentials(ctx context.Context, s logical.Storage) error {
        for {
                item, err := b.credRotationQueue.PopItem()
                if err != nil {
                        if err == queue.ErrEmpty {
                                return nil
                        }
                        return err
                }

                role, err := b.Role(ctx, s, item.Key)
                if err != nil {
                        b.logger.Warn(fmt.Sprintf("unable load role (%s)", item.Key), "error", err)
                        continue
                }
                if role == nil {
                        b.logger.Warn(fmt.Sprintf("role (%s) not found", item.Key), "error", err)
                        continue
                }

                if time.Now().Unix() > item.Priority {
                        // We've found our first item not in need of rotation
                        input := &setPasswordInput{
                                RoleName: item.Key,
                                Role:     role,
                        }

                        // check for existing WAL entry with a Password
                        if walID, ok := item.Value.(string); ok {
                                walEntry := b.walForItemValue(ctx, s, walID)
                                if walEntry != nil && walEntry.NewPassword != "" {
                                        input.Password = walEntry.NewPassword
                                        input.WALID = walID
                                }
                        }

                        // lvr is the roles' last vault rotation
                        resp, err := b.createUpdateStaticAccount(ctx, s, input)
                        if err != nil {
                                b.logger.Warn("unable rotate credentials in periodic function", "error", err)
                                // add the item to the re-queue slice
                                newItem := queue.Item{
                                        Key:      item.Key,
                                        Priority: item.Priority + 10,
                                }

                                // preserve the WALID if it was returned
                                if resp.WALID != "" {
                                        newItem.Value = resp.WALID
                                }

                                if err := b.credRotationQueue.PushItem(&newItem); err != nil {
                                        b.logger.Warn("unable to push item on to queue", "error", err)
                                }
                                // go to next item
                                continue
                        }

                        // guard against RotationTime not being set or zero-value
                        lvr := resp.RotationTime
                        if lvr.IsZero() {
                                lvr = time.Now()
                        }

                        nextRotation := lvr.Add(role.StaticAccount.RotationPeriod)
                        newItem := queue.Item{
                                Key:      item.Key,
                                Priority: nextRotation.Unix(),
                        }
                        if err := b.credRotationQueue.PushItem(&newItem); err != nil {
                                b.logger.Warn("unable to push item on to queue", "error", err)
                        }
                } else {
                        // highest priority item does not need rotation, so we push it back on
                        // the queue and break the loop
                        b.credRotationQueue.PushItem(item)
                        break
                }
        }
        return nil
}

// walForItemValue looks in the Backend's WAL entries for an entry with a
// specific ID. If found and of type walRotationKey, return the parsed
// walSetCredentials struct, otherwise return nil
func (b *databaseBackend) walForItemValue(ctx context.Context, s logical.Storage, id string) *walSetCredentials {
        wal, err := framework.GetWAL(ctx, s, id)
        if err != nil {
                b.Logger().Warn(fmt.Sprintf("error reading WAL for ID (%s):", id), err)
                return nil
        }

        if wal == nil || wal.Kind != walRotationKey {
                return nil
        }

        var walEntry walSetCredentials
        if mapErr := mapstructure.Decode(wal.Data, &walEntry); err != nil {
                b.Logger().Warn("error decoding walEntry", mapErr.Error())
                return nil
        }

        return &walEntry
}

// TODO: rename to match the method these go with
type setPasswordInput struct {
        RoleName   string
        Role       *roleEntry
        Password   string
        CreateUser bool
        WALID      string
}

type setPasswordResponse struct {
        RotationTime time.Time
        // Optional return field, in the event WAL was created and not destroyed
        // during the operation
        WALID string
}

func (b *databaseBackend) createUpdateStaticAccount(ctx context.Context, s logical.Storage, input *setPasswordInput) (*setPasswordResponse, error) {
        var lvr time.Time
        var merr error
        // re-use WAL ID if present, otherwise PUT a new WAL
        setResponse := &setPasswordResponse{WALID: input.WALID}

        dbConfig, err := b.DatabaseConfig(ctx, s, input.Role.DBName)
        if err != nil {
                return setResponse, err
        }

        // If role name isn't in the database's allowed roles, send back a
        // permission denied.
        if !strutil.StrListContains(dbConfig.AllowedRoles, "*") && !strutil.StrListContainsGlob(dbConfig.AllowedRoles, input.RoleName) {
                return setResponse, fmt.Errorf("%q is not an allowed role", input.RoleName)
        }

        // Get the Database object
        db, err := b.GetConnection(ctx, s, input.Role.DBName)
        if err != nil {
                return setResponse, err
        }

        // Use password from input if available. This happens if we're restoring from
        // a WAL item or processing the rotation queue with an item that has a WAL
        // associated with it
        newPassword := input.Password
        if newPassword == "" {
                // Generate a new password
                newPassword, err = db.GenerateCredentials(ctx)
                if err != nil {
                        return setResponse, err
                }
        }

        db.RLock()
        defer db.RUnlock()

        config := dbplugin.StaticUserConfig{
                Username: input.Role.StaticAccount.Username,
                Password: newPassword,
        }

        // Create or rotate the user
        stmts := input.Role.Statements.Creation
        if !input.CreateUser {
                stmts = input.Role.Statements.Rotation
        }

        if setResponse.WALID == "" {
                setResponse.WALID, err = framework.PutWAL(ctx, s, walRotationKey, &walSetCredentials{
                        RoleName:          input.RoleName,
                        Username:          config.Username,
                        NewPassword:       config.Password,
                        OldPassword:       input.Role.StaticAccount.Password,
                        Statements:        stmts,
                        LastVaultRotation: input.Role.StaticAccount.LastVaultRotation,
                })
                if err != nil {
                        // TODO: error wrap here?
                        return setResponse, errwrap.Wrapf("error writing WAL entry: {{err}}", err)
                }
        }

        var sterr error
        _, password, _, sterr := db.SetCredentials(ctx, config, stmts)
        if sterr != nil {
                b.CloseIfShutdown(db, sterr)
                return setResponse, sterr
        }

        // TODO set credentials doesn't need to return all these things
        if newPassword != password {
                return setResponse, errors.New("mismatch password returned")
        }

        // Store updated role information
        lvr = time.Now()
        input.Role.StaticAccount.LastVaultRotation = lvr
        input.Role.StaticAccount.Password = password
        setResponse.RotationTime = lvr

        entry, err := logical.StorageEntryJSON("role/"+input.RoleName, input.Role)
        if err != nil {
                return setResponse, err
        }
        if err := s.Put(ctx, entry); err != nil {
                return setResponse, err
        }

        // cleanup WAL after successfully rotating and pushing new item on to queue
        if err := framework.DeleteWAL(ctx, s, setResponse.WALID); err != nil {
                merr = multierror.Append(merr, err)
        }

        // return a new setPasswordResponse without the WALID, since we deleted it
        return &setPasswordResponse{RotationTime: lvr}, merr
}

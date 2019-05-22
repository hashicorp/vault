package database

import (
	"context"
	"errors"
	"fmt"
	"net/rpc"
	"strings"
	"sync"
	"time"

	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
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
	queueTickInterval = 5 * time.Second

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

	// load queue and kickoff new periodic ticker
	go b.initQueue(conf)
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
func (b *databaseBackend) initQueue(conf *logical.BackendConfig) {
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

		b.Lock()
		if b.credRotationQueue == nil {
			b.credRotationQueue = queue.New()
		}

		// create a context with a cancel method for processing any WAL entries and
		// populating the queue
		ctx, cancel := context.WithCancel(context.Background())
		b.cancelQueue = cancel
		b.Unlock()

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

	// build map of role name / wal entries
	walMap, err := b.loadStaticWALs(ctx, s)
	if err != nil {
		log.Warn("unable to load rotation WALs", "error", err)
	}

	roles, err := s.List(ctx, databaseStaticRolePath)
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

		role, err := b.StaticRole(ctx, s, roleName)
		if err != nil {
			log.Warn("unable to read static role", "error", err, "role", roleName)
			continue
		}
		walEntry := walMap[roleName]
		if role == nil || role.StaticAccount == nil {
			if walEntry != nil {
				if err := framework.DeleteWAL(ctx, s, walEntry.walID); err != nil {
					log.Warn("unable to delete WAL", "error", err, "WAL ID", walEntry.walID)
				}
			}
			continue
		}
		item := queue.Item{
			Key:      roleName,
			Priority: role.StaticAccount.LastVaultRotation.Add(role.StaticAccount.RotationPeriod).Unix(),
		}

		// check if role name is in map
		if walEntry != nil {
			// check walEntry last vault time
			if !walEntry.LastVaultRotation.IsZero() && walEntry.LastVaultRotation.Before(role.StaticAccount.LastVaultRotation) {
				// WAL's last vault rotation record is older than the role's data, so
				// delete and move on
				if err := framework.DeleteWAL(ctx, s, walEntry.walID); err != nil {
					log.Warn("unable to delete WAL", "error", err, "WAL ID", walEntry.walID)
				}
			} else {
				log.Info("adjusting priority for Role")
				item.Value = walEntry.walID
				item.Priority = time.Now().Unix()
			}
		}

		if err := b.pushItem(&item); err != nil {
			log.Warn("unable to enqueue item", "error", err, "role", roleName)
		}
	}
}

// runTicker kicks off a periodic ticker that invoke the automatic credential
// rotation method at a determined interval. The default interval is 5 seconds.
func (b *databaseBackend) runTicker(ctx context.Context, s logical.Storage) {
	b.logger.Info("starting periodic ticker")
	tick := time.NewTicker(queueTickInterval)
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
}

// setCredentialsWAL is used to store information in a WAL that can retry a
// credential setting or rotation in the event of partial failure.
type setCredentialsWAL struct {
	NewPassword string
	OldPassword string
	RoleName    string
	Username    string

	LastVaultRotation time.Time

	walID string
}

// rotateCredentials sets a new password for a static account. This method is
// invoked in the runTicker method, which is in it's own go-routine, and invoked
// periodically (approximately every 5 seconds).
//
// This method loops through the priority queue, popping the highest priority
// item until it encounters the first item that does not yet need rotation,
// based on the current time.
func (b *databaseBackend) rotateCredentials(ctx context.Context, s logical.Storage) error {
	for {
		// quit rotating credentials if shutdown has started
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		item, err := b.popItem()
		if err != nil {
			if err == queue.ErrEmpty {
				return nil
			}
			return err
		}

		// guard against possible nil item
		if item == nil {
			return nil
		}

		// validate the role still exists
		role, err := b.StaticRole(ctx, s, item.Key)
		if err != nil {
			b.logger.Warn("unable load role", "role", item.Key, "error", err)
			continue
		}
		if role == nil {
			b.logger.Warn("role not found", "role", item.Key, "error", err)
			continue
		}

		// if "now" is less than the Item priority, then this item does not need to
		// be rotated
		if time.Now().Unix() < item.Priority {
			if err := b.pushItem(item); err != nil {
				b.logger.Warn("unable to push item on to queue", "error", err)
			}
			// break out of the for loop
			break
		}

		input := &setStaticAccountInput{
			RoleName: item.Key,
			Role:     role,
		}

		// If there is a WAL entry related to this Role, the corresponding WAL ID
		// should be stored in the Item's Value field.
		if walID, ok := item.Value.(string); ok {
			walEntry := b.findStaticWAL(ctx, s, walID)
			if walEntry != nil && walEntry.NewPassword != "" {
				input.Password = walEntry.NewPassword
				input.WALID = walID
			}
		}

		resp, err := b.setStaticAccount(ctx, s, input)
		if err != nil {
			b.logger.Warn("unable to rotate credentials in periodic function", "error", err)
			// Increment the priority enough so that the next call to this method
			// likely will not attempt to rotate it, as a back-off of sorts
			item.Priority = time.Now().Add(10 * time.Second).Unix()

			// preserve the WALID if it was returned
			if resp.WALID != "" {
				item.Value = resp.WALID
			}

			if err := b.pushItem(item); err != nil {
				b.logger.Warn("unable to push item on to queue", "error", err)
			}
			// go to next item
			continue
		}

		lvr := resp.RotationTime
		if lvr.IsZero() {
			lvr = time.Now()
		}

		// update priority and push updated Item to the queue
		nextRotation := lvr.Add(role.StaticAccount.RotationPeriod)
		item.Priority = nextRotation.Unix()
		if err := b.pushItem(item); err != nil {
			b.logger.Warn("unable to push item on to queue", "error", err)
		}
	}
	return nil
}

// findStaticWAL loads a WAL entry by ID. If found, only return the WAL if it
// is of type staticWALKey, otherwise return nil
func (b *databaseBackend) findStaticWAL(ctx context.Context, s logical.Storage, id string) *setCredentialsWAL {
	wal, err := framework.GetWAL(ctx, s, id)
	if err != nil {
		b.Logger().Warn("error reading WAL", "id", id, "error", err)
		return nil
	}

	if wal == nil || wal.Kind != staticWALKey {
		return nil
	}

	data := wal.Data.(map[string]interface{})
	walEntry := setCredentialsWAL{
		walID:       id,
		NewPassword: data["NewPassword"].(string),
		OldPassword: data["OldPassword"].(string),
		RoleName:    data["RoleName"].(string),
		Username:    data["Username"].(string),
	}
	lvr, err := time.Parse(time.RFC3339, data["LastVaultRotation"].(string))
	if err != nil {
		b.Logger().Warn("error decoding walEntry", err.Error())
		return nil
	}
	walEntry.LastVaultRotation = lvr

	return &walEntry
}

type setStaticAccountInput struct {
	RoleName   string
	Role       *roleEntry
	Password   string
	CreateUser bool
	WALID      string
}

type setStaticAccountOutput struct {
	RotationTime time.Time
	Password     string
	// Optional return field, in the event WAL was created and not destroyed
	// during the operation
	WALID string
}

// setStaticAccount sets the password for a static account associated with a
// Role. This method does many things:
// - verifies role exists and is in the allowed roles list
// - loads an existing WAL entry if WALID input is given, otherwise creates a
// new WAL entry
// - gets a database connection
// - accepts an input password, otherwise generates a new one via gRPC to the
// database plugin
// - sets new password for the static account
// - uses WAL for ensuring passwords are not lost if storage to Vault fails
//
// This method does not perform any operations on the priority queue. Those
// tasks must be handled outside of this method.
func (b *databaseBackend) setStaticAccount(ctx context.Context, s logical.Storage, input *setStaticAccountInput) (*setStaticAccountOutput, error) {
	// lvr is the known LastVaultRotation
	var lvr time.Time
	var merr error
	// re-use WAL ID if present, otherwise PUT a new WAL
	if input == nil {
		input = &setStaticAccountInput{}
	}
	output := &setStaticAccountOutput{WALID: input.WALID}

	dbConfig, err := b.DatabaseConfig(ctx, s, input.Role.DBName)
	if err != nil {
		return output, err
	}

	// If role name isn't in the database's allowed roles, send back a
	// permission denied.
	if !strutil.StrListContains(dbConfig.AllowedRoles, "*") && !strutil.StrListContainsGlob(dbConfig.AllowedRoles, input.RoleName) {
		return output, fmt.Errorf("%q is not an allowed role", input.RoleName)
	}

	// Get the Database object
	db, err := b.GetConnection(ctx, s, input.Role.DBName)
	if err != nil {
		return output, err
	}

	db.RLock()
	defer db.RUnlock()

	// Use password from input if available. This happens if we're restoring from
	// a WAL item or processing the rotation queue with an item that has a WAL
	// associated with it
	newPassword := input.Password
	if newPassword == "" {
		// Generate a new password
		newPassword, err = db.GenerateCredentials(ctx)
		if err != nil {
			return output, err
		}
	}
	output.Password = newPassword

	config := dbplugin.StaticUserConfig{
		Username: input.Role.StaticAccount.Username,
		Password: newPassword,
	}

	if output.WALID == "" {
		output.WALID, err = framework.PutWAL(ctx, s, staticWALKey, &setCredentialsWAL{
			RoleName:          input.RoleName,
			Username:          config.Username,
			NewPassword:       config.Password,
			OldPassword:       input.Role.StaticAccount.Password,
			LastVaultRotation: input.Role.StaticAccount.LastVaultRotation,
		})
		if err != nil {
			return output, errwrap.Wrapf("error writing WAL entry: {{err}}", err)
		}
	}

	_, password, err := db.SetCredentials(ctx, input.Role.Statements, config)
	if err != nil {
		b.CloseIfShutdown(db, err)
		return output, errwrap.Wrapf("error setting credentials: {{err}}", err)
	}

	if newPassword != password {
		return output, errors.New("mismatch passwords returned")
	}

	// Store updated role information
	lvr = time.Now()
	input.Role.StaticAccount.LastVaultRotation = lvr
	input.Role.StaticAccount.Password = password
	output.RotationTime = lvr

	entry, err := logical.StorageEntryJSON(databaseStaticRolePath+input.RoleName, input.Role)
	if err != nil {
		return output, err
	}
	if err := s.Put(ctx, entry); err != nil {
		return output, err
	}

	// cleanup WAL after successfully rotating and pushing new item on to queue
	if err := framework.DeleteWAL(ctx, s, output.WALID); err != nil {
		merr = multierror.Append(merr, err)
		return output, merr
	}

	// the WAL has been deleted, return new setStaticAccountOutput without it
	return &setStaticAccountOutput{RotationTime: lvr}, merr
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
	// Upgrade lock
	b.RUnlock()
	b.Lock()
	unlockFunc = b.Unlock

	// check again
	if b.credRotationQueue != nil {
		return b.credRotationQueue.Push(item)
	}
	b.credRotationQueue = queue.New()

	return b.credRotationQueue.Push(item)
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

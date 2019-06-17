package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	// Interval to check the queue for items needing rotation
	queueTickSeconds  = 5
	queueTickInterval = queueTickSeconds * time.Second

	// WAL storage key used for static account rotations
	staticWALKey = "staticRotationKey"
)

// populateQueue loads the priority queue with existing static accounts. This
// occurs at initialization, after any WAL entries of failed or interrupted
// rotations have been processed. It lists the roles from storage and searches
// for any that have an associated static account, then adds them to the
// priority queue for rotations.
func (b *databaseBackend) populateQueue(ctx context.Context, s logical.Storage) {
	log := b.Logger()
	log.Info("populating role rotation queue")

	// Build map of role name / wal entries
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

		item := queue.Item{
			Key:      roleName,
			Priority: role.StaticAccount.LastVaultRotation.Add(role.StaticAccount.RotationPeriod).Unix(),
		}

		// Check if role name is in map
		walEntry := walMap[roleName]
		if walEntry != nil {
			// Check walEntry last vault time
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
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
	RoleName    string `json:"role_name"`
	Username    string `json:"username"`

	LastVaultRotation time.Time `json:"last_vault_rotation"`

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
		// Quit rotating credentials if shutdown has started
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		item, err := b.popFromRotationQueue()
		if err != nil {
			if err == queue.ErrEmpty {
				return nil
			}
			return err
		}

		// Guard against possible nil item
		if item == nil {
			return nil
		}

		// Grab the exclusive lock for this Role, to make sure we don't incur and
		// writes during the rotation process
		lock := locksutil.LockForKey(b.roleLocks, item.Key)
		lock.Lock()
		defer lock.Unlock()

		// Validate the role still exists
		role, err := b.StaticRole(ctx, s, item.Key)
		if err != nil {
			b.logger.Error("unable to load role", "role", item.Key, "error", err)
			item.Priority = time.Now().Add(10 * time.Second).Unix()
			if err := b.pushItem(item); err != nil {
				b.logger.Error("unable to push item on to queue", "error", err)
			}
			continue
		}
		if role == nil {
			b.logger.Warn("role not found", "role", item.Key, "error", err)
			continue
		}

		// If "now" is less than the Item priority, then this item does not need to
		// be rotated
		if time.Now().Unix() < item.Priority {
			if err := b.pushItem(item); err != nil {
				b.logger.Error("unable to push item on to queue", "error", err)
			}
			// Break out of the for loop
			break
		}

		input := &setStaticAccountInput{
			RoleName: item.Key,
			Role:     role,
		}

		// If there is a WAL entry related to this Role, the corresponding WAL ID
		// should be stored in the Item's Value field.
		if walID, ok := item.Value.(string); ok {
			walEntry, err := b.findStaticWAL(ctx, s, walID)
			if err != nil {
				b.logger.Error("error finding static WAL", "error", err)
				item.Priority = time.Now().Add(10 * time.Second).Unix()
				if err := b.pushItem(item); err != nil {
					b.logger.Error("unable to push item on to queue", "error", err)
				}
			}
			if walEntry != nil && walEntry.NewPassword != "" {
				input.Password = walEntry.NewPassword
				input.WALID = walID
			}
		}

		resp, err := b.setStaticAccount(ctx, s, input)
		if err != nil {
			b.logger.Error("unable to rotate credentials in periodic function", "error", err)
			// Increment the priority enough so that the next call to this method
			// likely will not attempt to rotate it, as a back-off of sorts
			item.Priority = time.Now().Add(10 * time.Second).Unix()

			// Preserve the WALID if it was returned
			if resp != nil && resp.WALID != "" {
				item.Value = resp.WALID
			}

			if err := b.pushItem(item); err != nil {
				b.logger.Error("unable to push item on to queue", "error", err)
			}
			// Go to next item
			continue
		}

		lvr := resp.RotationTime
		if lvr.IsZero() {
			lvr = time.Now()
		}

		// Update priority and push updated Item to the queue
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
func (b *databaseBackend) findStaticWAL(ctx context.Context, s logical.Storage, id string) (*setCredentialsWAL, error) {
	wal, err := framework.GetWAL(ctx, s, id)
	if err != nil {
		return nil, err
	}

	if wal == nil || wal.Kind != staticWALKey {
		return nil, nil
	}

	data := wal.Data.(map[string]interface{})
	walEntry := setCredentialsWAL{
		walID:       id,
		NewPassword: data["new_password"].(string),
		OldPassword: data["old_password"].(string),
		RoleName:    data["role_name"].(string),
		Username:    data["username"].(string),
	}
	lvr, err := time.Parse(time.RFC3339, data["last_vault_rotation"].(string))
	if err != nil {
		return nil, err
	}
	walEntry.LastVaultRotation = lvr

	return &walEntry, nil
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
	var merr error
	if input == nil || input.Role == nil || input.RoleName == "" {
		return nil, errors.New("input was empty when attempting to set credentials for static account")
	}
	// Re-use WAL ID if present, otherwise PUT a new WAL
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
	// lvr is the known LastVaultRotation
	lvr := time.Now()
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

	// Cleanup WAL after successfully rotating and pushing new item on to queue
	if err := framework.DeleteWAL(ctx, s, output.WALID); err != nil {
		merr = multierror.Append(merr, err)
		return output, merr
	}

	// The WAL has been deleted, return new setStaticAccountOutput without it
	return &setStaticAccountOutput{RotationTime: lvr}, merr
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
	// Verify this mount is on the primary server, or is a local mount. If not, do
	// not create a queue or launch a ticker. Both processing the WAL list and
	// populating the queue are done sequentially and before launching a
	// go-routine to run the periodic ticker.
	replicationState := conf.System.ReplicationState()
	if (conf.System.LocalMount() || !replicationState.HasState(consts.ReplicationPerformanceSecondary)) &&
		!replicationState.HasState(consts.ReplicationDRSecondary) &&
		!replicationState.HasState(consts.ReplicationPerformanceStandby) {
		b.Logger().Info("initializing database rotation queue")

		// Poll for a PutWAL call that does not return a "read-only storage" error.
		// This ensures the startup phases of loading WAL entries from any possible
		// failed rotations can complete without error when deleting from storage.
	READONLY_LOOP:
		for {
			select {
			case <-ctx.Done():
				b.Logger().Info("queue initialization canceled")
				return
			default:
			}

			walID, err := framework.PutWAL(ctx, conf.StorageView, staticWALKey, &setCredentialsWAL{RoleName: "vault-readonlytest"})
			if walID != "" {
				defer framework.DeleteWAL(ctx, conf.StorageView, walID)
			}
			switch {
			case err == nil:
				break READONLY_LOOP
			case err.Error() == logical.ErrSetupReadOnly.Error():
				time.Sleep(10 * time.Millisecond)
			default:
				b.Logger().Error("deleting nil key resulted in error", "error", err)
				return
			}
		}

		// Load roles and populate queue with static accounts
		b.populateQueue(ctx, conf.StorageView)

		// Launch ticker
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
	// Loop through WAL keys and process any rotation ones
	for _, walID := range keys {
		walEntry, err := b.findStaticWAL(ctx, s, walID)
		if err != nil {
			b.Logger().Error("error loading static WAL", "id", walID, "error", err)
			continue
		}
		if walEntry == nil {
			continue
		}

		// Verify the static role still exists
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

// popFromRotationQueue wraps the internal queue's Pop call, to make sure a queue is
// actually available. This is needed because both runTicker and initQueue
// operate in go-routines, and could be accessing the queue concurrently
func (b *databaseBackend) popFromRotationQueue() (*queue.Item, error) {
	b.RLock()
	defer b.RUnlock()
	if b.credRotationQueue != nil {
		return b.credRotationQueue.Pop()
	}
	return nil, queue.ErrEmpty
}

// popFromRotationQueueByKey wraps the internal queue's PopByKey call, to make sure a queue is
// actually available. This is needed because both runTicker and initQueue
// operate in go-routines, and could be accessing the queue concurrently
func (b *databaseBackend) popFromRotationQueueByKey(name string) (*queue.Item, error) {
	b.RLock()
	defer b.RUnlock()
	if b.credRotationQueue != nil {
		return b.credRotationQueue.PopByKey(name)
	}
	return nil, queue.ErrEmpty
}

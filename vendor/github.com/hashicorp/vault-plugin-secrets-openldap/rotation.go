package openldap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
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
func (b *backend) populateQueue(ctx context.Context, s logical.Storage) {
	log := b.Logger()
	log.Info("populating role rotation queue")

	// Build map of role name / wal entries
	walMap, err := b.loadStaticWALs(ctx, s)
	if err != nil {
		log.Warn("unable to load rotation WALs", "error", err)
	}

	roles, err := s.List(ctx, staticRolePath)
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

		role, err := b.staticRole(ctx, s, roleName)
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
func (b *backend) runTicker(ctx context.Context, s logical.Storage) {
	b.Logger().Info("starting periodic ticker")
	tick := time.NewTicker(queueTickInterval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			b.rotateCredentials(ctx, s)

		case <-ctx.Done():
			b.Logger().Info("stopping periodic ticker")
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
	DN          string `json:"dn"`

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
func (b *backend) rotateCredentials(ctx context.Context, s logical.Storage) {
	for b.rotateCredential(ctx, s) {
	}
}

func (b *backend) rotateCredential(ctx context.Context, s logical.Storage) bool {
	// Quit rotating credentials if shutdown has started
	select {
	case <-ctx.Done():
		return false
	default:
	}
	item, err := b.popFromRotationQueue()
	if err != nil {
		if err != queue.ErrEmpty {
			b.Logger().Error("error popping item from queue", "err", err)
		}
		return false
	}

	// Guard against possible nil item
	if item == nil {
		return false
	}

	// Grab the exclusive lock for this Role, to make sure we don't incur and
	// writes during the rotation process
	lock := locksutil.LockForKey(b.roleLocks, item.Key)
	lock.Lock()
	defer lock.Unlock()

	// Validate the role still exists
	role, err := b.staticRole(ctx, s, item.Key)
	if err != nil {
		b.Logger().Error("unable to load role", "role", item.Key, "error", err)
		item.Priority = time.Now().Add(10 * time.Second).Unix()
		if err := b.pushItem(item); err != nil {
			b.Logger().Error("unable to push item on to queue", "error", err)
		}
		return true
	}
	if role == nil {
		b.Logger().Warn("role not found", "role", item.Key, "error", err)
		return true
	}

	// If "now" is less than the Item priority, then this item does not need to
	// be rotated
	if time.Now().Unix() < item.Priority {
		if err := b.pushItem(item); err != nil {
			b.Logger().Error("unable to push item on to queue", "error", err)
		}
		// Break out of the for loop
		return false
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
			b.Logger().Error("error finding static WAL", "error", err)
			item.Priority = time.Now().Add(10 * time.Second).Unix()
			if err := b.pushItem(item); err != nil {
				b.Logger().Error("unable to push item on to queue", "error", err)
			}
		}
		if walEntry != nil && walEntry.NewPassword != "" {
			input.Password = walEntry.NewPassword
			input.WALID = walID
		}
	}

	resp, err := b.setStaticAccountPassword(ctx, s, input)
	if err != nil {
		b.Logger().Error("unable to rotate credentials in periodic function", "error", err)
		// Increment the priority enough so that the next call to this method
		// likely will not attempt to rotate it, as a back-off of sorts
		item.Priority = time.Now().Add(10 * time.Second).Unix()

		// Preserve the WALID if it was returned
		if resp != nil && resp.WALID != "" {
			item.Value = resp.WALID
		}

		if err := b.pushItem(item); err != nil {
			b.Logger().Error("unable to push item on to queue", "error", err)
		}
		// Go to next item
		return true
	}

	lvr := resp.RotationTime
	if lvr.IsZero() {
		lvr = time.Now()
	}

	// Update priority and push updated Item to the queue
	nextRotation := lvr.Add(role.StaticAccount.RotationPeriod)
	item.Priority = nextRotation.Unix()
	if err := b.pushItem(item); err != nil {
		b.Logger().Warn("unable to push item on to queue", "error", err)
	}
	return true
}

// findStaticWAL loads a WAL entry by ID. If found, only return the WAL if it
// is of type staticWALKey, otherwise return nil
func (b *backend) findStaticWAL(ctx context.Context, s logical.Storage, id string) (*setCredentialsWAL, error) {
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
		DN:          data["dn"].(string),
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

// setStaticAccountPassword sets the password for a static account associated with a
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
func (b *backend) setStaticAccountPassword(ctx context.Context, s logical.Storage, input *setStaticAccountInput) (*setStaticAccountOutput, error) {
	var merr error
	if input == nil || input.Role == nil || input.RoleName == "" {
		return nil, errors.New("input was empty when attempting to set credentials for static account")
	}

	if _, hasTimeout := ctx.Deadline(); !hasTimeout {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, defaultCtxTimeout)
		defer cancel()
	}

	// Re-use WAL ID if present, otherwise PUT a new WAL
	output := &setStaticAccountOutput{WALID: input.WALID}

	config, err := readConfig(ctx, s)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("the config is currently unset")
	}

	newPassword, err := b.GeneratePassword(ctx, config)
	if err != nil {
		return nil, err
	}

	oldPassword := input.Role.StaticAccount.Password

	// Take out the backend lock since we are swapping out the connection
	b.Lock()
	defer b.Unlock()

	if output.WALID == "" {
		output.WALID, err = framework.PutWAL(ctx, s, staticWALKey, &setCredentialsWAL{
			RoleName:          input.RoleName,
			Username:          input.Role.StaticAccount.Username,
			DN:                input.Role.StaticAccount.DN,
			NewPassword:       newPassword,
			OldPassword:       oldPassword,
			LastVaultRotation: input.Role.StaticAccount.LastVaultRotation,
		})
		if err != nil {
			return output, errwrap.Wrapf("error writing WAL entry: {{err}}", err)
		}
	}

	// Update the password remotely.
	if err := b.client.UpdatePassword(config.LDAP, input.Role.StaticAccount.DN, newPassword); err != nil {
		return nil, err
	}

	// Update the password locally.
	if pwdStoringErr := storePassword(ctx, s, config); pwdStoringErr != nil {
		// We were unable to store the new password locally. We can't continue in this state because we won't be able
		// to roll any passwords, including our own to get back into a state of working. So, we need to roll back to
		// the last password we successfully got into storage.
		if rollbackErr := b.rollBackPassword(ctx, config, oldPassword); rollbackErr != nil {
			return nil, fmt.Errorf(`unable to store new password due to %s and unable to return to previous password due
to %s, configure a new binddn and bindpass to restore openldap function`, pwdStoringErr, rollbackErr)
		}
		return nil, fmt.Errorf("unable to update password due to storage err: %s", pwdStoringErr)
	}

	// Store updated role information
	// lvr is the known LastVaultRotation
	lvr := time.Now()
	input.Role.StaticAccount.LastVaultRotation = lvr
	input.Role.StaticAccount.Password = newPassword
	output.RotationTime = lvr

	entry, err := logical.StorageEntryJSON(staticRolePath+input.RoleName, input.Role)
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

func (b *backend) GeneratePassword(ctx context.Context, cfg *config) (string, error) {
	if cfg.PasswordPolicy == "" {
		if cfg.PasswordLength == 0 {
			return base62.Random(defaultPasswordLength)
		}
		return base62.Random(cfg.PasswordLength)
	}

	password, err := b.System().GeneratePasswordFromPolicy(ctx, cfg.PasswordPolicy)
	if err != nil {
		return "", fmt.Errorf("unable to generate password: %w", err)
	}
	return password, nil
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
func (b *backend) initQueue(ctx context.Context, conf *logical.InitializationRequest) {
	// Verify this mount is on the primary server, or is a local mount. If not, do
	// not create a queue or launch a ticker. Both processing the WAL list and
	// populating the queue are done sequentially and before launching a
	// go-routine to run the periodic ticker.
	replicationState := b.System().ReplicationState()
	if (b.System().LocalMount() || !replicationState.HasState(consts.ReplicationPerformanceSecondary)) &&
		!replicationState.HasState(consts.ReplicationDRSecondary) &&
		!replicationState.HasState(consts.ReplicationPerformanceStandby) {
		b.Logger().Info("initializing database rotation queue")

		// Load roles and populate queue with static accounts
		b.populateQueue(ctx, conf.Storage)

		// Launch ticker
		go b.runTicker(ctx, conf.Storage)
	}
}

// loadStaticWALs reads WAL entries and returns a map of roles and their
// setCredentialsWAL, if found.
func (b *backend) loadStaticWALs(ctx context.Context, s logical.Storage) (map[string]*setCredentialsWAL, error) {
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
		role, err := b.staticRole(ctx, s, roleName)
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
func (b *backend) pushItem(item *queue.Item) error {
	b.RLock()
	defer b.RUnlock()

	if b.credRotationQueue != nil {
		return b.credRotationQueue.Push(item)
	}

	b.Logger().Warn("no queue found during push item")
	return nil
}

// popFromRotationQueue wraps the internal queue's Pop call, to make sure a queue is
// actually available. This is needed because both runTicker and initQueue
// operate in go-routines, and could be accessing the queue concurrently
func (b *backend) popFromRotationQueue() (*queue.Item, error) {
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
func (b *backend) popFromRotationQueueByKey(name string) (*queue.Item, error) {
	b.RLock()
	defer b.RUnlock()
	if b.credRotationQueue != nil {
		item, err := b.credRotationQueue.PopByKey(name)
		if err != nil {
			return nil, err
		}
		if item != nil {
			return item, nil
		}
	}
	return nil, queue.ErrEmpty
}

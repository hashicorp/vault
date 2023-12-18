// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	// Default interval to check the queue for items needing rotation
	defaultQueueTickSeconds = 5

	// Config key to set an alternate interval
	queueTickIntervalKey = "rotation_queue_tick_interval"

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
			Priority: role.StaticAccount.NextRotationTime().Unix(),
		}

		// Check if role name is in map
		walEntry := walMap[roleName]
		if walEntry != nil {
			// Check walEntry last vault time
			if walEntry.LastVaultRotation.IsZero() {
				// A WAL's last Vault rotation can only ever be 0 for a role that
				// was never successfully created. So we know this WAL couldn't
				// have been created for this role we just retrieved from storage.
				// i.e. it must be a hangover from a previous attempt at creating
				// a role with the same name
				log.Debug("deleting WAL with zero last rotation time", "WAL ID", walEntry.walID, "created", walEntry.walCreatedAt)
				if err := framework.DeleteWAL(ctx, s, walEntry.walID); err != nil {
					log.Warn("unable to delete zero-time WAL", "error", err, "WAL ID", walEntry.walID)
				}
			} else if walEntry.LastVaultRotation.Before(role.StaticAccount.LastVaultRotation) {
				// WAL's last vault rotation record is older than the role's data, so
				// delete and move on
				log.Debug("deleting outdated WAL", "WAL ID", walEntry.walID, "created", walEntry.walCreatedAt)
				if err := framework.DeleteWAL(ctx, s, walEntry.walID); err != nil {
					log.Warn("unable to delete WAL", "error", err, "WAL ID", walEntry.walID)
				}
			} else {
				// previous rotation attempt was interrupted, so we set the
				// Priority as highest to be processed immediately
				log.Info("found WAL for role", "role", item.Key, "WAL ID", walEntry.walID)
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
func (b *databaseBackend) runTicker(ctx context.Context, queueTickInterval time.Duration, s logical.Storage) {
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
	CredentialType v5.CredentialType `json:"credential_type"`
	NewPassword    string            `json:"new_password"`
	NewPublicKey   []byte            `json:"new_public_key"`
	NewPrivateKey  []byte            `json:"new_private_key"`
	RoleName       string            `json:"role_name"`
	Username       string            `json:"username"`

	LastVaultRotation time.Time `json:"last_vault_rotation"`

	// Private fields which will not be included in json.Marshal/Unmarshal.
	walID        string
	walCreatedAt int64 // Unix time at which the WAL was created.
}

// credentialIsSet returns true if the credential associated with the
// CredentialType field is properly set. See field comments to for a
// mapping of CredentialType values to respective credential fields.
func (w *setCredentialsWAL) credentialIsSet() bool {
	if w == nil {
		return false
	}

	switch w.CredentialType {
	case v5.CredentialTypePassword:
		return w.NewPassword != ""
	case v5.CredentialTypeRSAPrivateKey:
		return len(w.NewPublicKey) > 0
	default:
		return false
	}
}

// rotateCredentials sets a new password for a static account. This method is
// invoked in the runTicker method, which is in it's own go-routine, and invoked
// periodically (approximately every 5 seconds).
//
// This method loops through the priority queue, popping the highest priority
// item until it encounters the first item that does not yet need rotation,
// based on the current time.
func (b *databaseBackend) rotateCredentials(ctx context.Context, s logical.Storage) {
	for b.rotateCredential(ctx, s) {
	}
}

func (b *databaseBackend) rotateCredential(ctx context.Context, s logical.Storage) bool {
	// Quit rotating credentials if shutdown has started
	select {
	case <-ctx.Done():
		return false
	default:
	}
	item, err := b.popFromRotationQueue()
	if err != nil {
		if err != queue.ErrEmpty {
			b.logger.Error("error popping item from queue", "err", err)
		}
		return false
	}

	// Guard against possible nil item
	if item == nil {
		return false
	}

	roleName := item.Key
	logger := b.Logger().With("role", roleName)

	// Grab the exclusive lock for this Role, to make sure we don't incur and
	// writes during the rotation process
	lock := locksutil.LockForKey(b.roleLocks, roleName)
	lock.Lock()
	defer lock.Unlock()

	// Validate the role still exists
	role, err := b.StaticRole(ctx, s, roleName)
	if err != nil {
		logger.Error("unable to load role", "error", err)

		item.Priority = time.Now().Add(10 * time.Second).Unix()
		if err := b.pushItem(item); err != nil {
			logger.Error("unable to push item on to queue", "error", err)
		}
		return true
	}
	if role == nil {
		logger.Warn("role not found", "error", err)
		return true
	}

	logger = logger.With("database", role.DBName)

	input := &setStaticAccountInput{
		RoleName: roleName,
		Role:     role,
	}

	now := time.Now()
	if !role.StaticAccount.ShouldRotate(item.Priority, now) {
		if !role.StaticAccount.IsInsideRotationWindow(now) {
			// We are a schedule-based rotation and we are outside a rotation
			// window so we update priority and NextVaultRotation
			item.Priority = role.StaticAccount.NextRotationTimeFromInput(now).Unix()
			role.StaticAccount.SetNextVaultRotation(now)
			b.logger.Trace("outside schedule-based rotation window, update priority", "next", role.StaticAccount.NextRotationTime())

			// write to storage after updating NextVaultRotation so the next
			// time this item is checked for rotation our role that we retrieve
			// from storage reflects that change
			entry, err := logical.StorageEntryJSON(databaseStaticRolePath+input.RoleName, input.Role)
			if err != nil {
				logger.Error("unable to encode entry for storage", "error", err)
				return false
			}
			if err := s.Put(ctx, entry); err != nil {
				logger.Error("unable to write to storage", "error", err)
				return false
			}
		}
		// do not rotate now, push item back onto queue to be rotated later
		if err := b.pushItem(item); err != nil {
			logger.Error("unable to push item on to queue", "error", err)
		}
		// Break out of the for loop
		return false
	}

	// If there is a WAL entry related to this Role, the corresponding WAL ID
	// should be stored in the Item's Value field.
	if walID, ok := item.Value.(string); ok {
		input.WALID = walID
	}

	resp, err := b.setStaticAccount(ctx, s, input)
	if err != nil {
		logger.Error("unable to rotate credentials in periodic function", "error", err)

		// Increment the priority enough so that the next call to this method
		// likely will not attempt to rotate it, as a back-off of sorts
		item.Priority = time.Now().Add(10 * time.Second).Unix()

		// Preserve the WALID if it was returned
		if resp != nil && resp.WALID != "" {
			item.Value = resp.WALID
		}

		if err := b.pushItem(item); err != nil {
			logger.Error("unable to push item on to queue", "error", err)
		}
		// Go to next item
		return true
	}
	// Clear any stored WAL ID as we must have successfully deleted our WAL to get here.
	item.Value = ""

	lvr := resp.RotationTime
	if lvr.IsZero() {
		lvr = time.Now()
	}

	// Update priority and push updated Item to the queue
	item.Priority = role.StaticAccount.NextRotationTimeFromInput(lvr).Unix()

	if err := b.pushItem(item); err != nil {
		logger.Warn("unable to push item on to queue", "error", err)
	}
	return true
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
		walID:        id,
		walCreatedAt: wal.CreatedAt,
		NewPassword:  data["new_password"].(string),
		RoleName:     data["role_name"].(string),
		Username:     data["username"].(string),
	}
	lvr, err := time.Parse(time.RFC3339, data["last_vault_rotation"].(string))
	if err != nil {
		return nil, err
	}
	walEntry.LastVaultRotation = lvr

	return &walEntry, nil
}

type setStaticAccountInput struct {
	RoleName string
	Role     *roleEntry
	WALID    string
}

type setStaticAccountOutput struct {
	RotationTime time.Time
	// Optional return field, in the event WAL was created and not destroyed
	// during the operation
	WALID string
}

// setStaticAccount sets the credential for a static account associated with a
// Role. This method does many things:
// - verifies role exists and is in the allowed roles list
// - loads an existing WAL entry if WALID input is given, otherwise creates a
// new WAL entry
//   - gets a database connection
//   - accepts an input credential, otherwise generates a new one for
//     the role's credential type
//   - sets new credential for the static account
//   - uses WAL for ensuring new credentials are not lost if storage to Vault fails,
//     resulting in a partial failure.
//
// This method does not perform any operations on the priority queue. Those
// tasks must be handled outside of this method.
func (b *databaseBackend) setStaticAccount(ctx context.Context, s logical.Storage, input *setStaticAccountInput) (*setStaticAccountOutput, error) {
	if input == nil || input.Role == nil || input.RoleName == "" {
		return nil, errors.New("input was empty when attempting to set credentials for static account")
	}
	// Re-use WAL ID if present, otherwise PUT a new WAL
	output := &setStaticAccountOutput{WALID: input.WALID}

	dbConfig, err := b.DatabaseConfig(ctx, s, input.Role.DBName)
	if err != nil {
		return output, err
	}
	if dbConfig == nil {
		return output, errors.New("the config is currently unset")
	}

	// If role name isn't in the database's allowed roles, send back a
	// permission denied.
	if !strutil.StrListContains(dbConfig.AllowedRoles, "*") && !strutil.StrListContainsGlob(dbConfig.AllowedRoles, input.RoleName) {
		return output, fmt.Errorf("%q is not an allowed role", input.RoleName)
	}

	// If the plugin doesn't support the credential type, return an error
	if !dbConfig.SupportsCredentialType(input.Role.CredentialType) {
		return output, fmt.Errorf("unsupported credential_type: %q",
			input.Role.CredentialType.String())
	}

	// Get the Database object
	dbi, err := b.GetConnection(ctx, s, input.Role.DBName)
	if err != nil {
		return output, err
	}

	dbi.RLock()
	defer dbi.RUnlock()

	updateReq := v5.UpdateUserRequest{
		Username: input.Role.StaticAccount.Username,
	}
	statements := v5.Statements{
		Commands: input.Role.Statements.Rotation,
	}

	// Use credential from input if available. This happens if we're restoring from
	// a WAL item or processing the rotation queue with an item that has a WAL
	// associated with it
	if output.WALID != "" {
		wal, err := b.findStaticWAL(ctx, s, output.WALID)
		if err != nil {
			return output, fmt.Errorf("error retrieving WAL entry: %w", err)
		}
		switch {
		case wal == nil:
			b.Logger().Error("expected role to have WAL, but WAL not found in storage", "role", input.RoleName, "WAL ID", output.WALID)

			// Generate a new WAL entry and credential
			output.WALID = ""
		case !wal.credentialIsSet():
			b.Logger().Error("expected WAL to have a new credential set, but empty", "role", input.RoleName, "WAL ID", output.WALID)
			if err := framework.DeleteWAL(ctx, s, output.WALID); err != nil {
				b.Logger().Warn("failed to delete WAL with no new credential", "error", err, "WAL ID", output.WALID)
			}

			// Generate a new WAL entry and credential
			output.WALID = ""
		case wal.CredentialType == v5.CredentialTypePassword:
			// Roll forward by using the credential in the existing WAL entry
			updateReq.CredentialType = v5.CredentialTypePassword
			updateReq.Password = &v5.ChangePassword{
				NewPassword: wal.NewPassword,
				Statements:  statements,
			}
			input.Role.StaticAccount.Password = wal.NewPassword
		case wal.CredentialType == v5.CredentialTypeRSAPrivateKey:
			// Roll forward by using the credential in the existing WAL entry
			updateReq.CredentialType = v5.CredentialTypeRSAPrivateKey
			updateReq.PublicKey = &v5.ChangePublicKey{
				NewPublicKey: wal.NewPublicKey,
				Statements:   statements,
			}
			input.Role.StaticAccount.PrivateKey = wal.NewPrivateKey
		}
	}

	// Generate a new credential
	if output.WALID == "" {
		walEntry := &setCredentialsWAL{
			RoleName:          input.RoleName,
			Username:          input.Role.StaticAccount.Username,
			LastVaultRotation: input.Role.StaticAccount.LastVaultRotation,
		}

		switch input.Role.CredentialType {
		case v5.CredentialTypePassword:
			generator, err := newPasswordGenerator(input.Role.CredentialConfig)
			if err != nil {
				return output, fmt.Errorf("failed to construct credential generator: %s", err)
			}

			// Fall back to database config-level password policy if not set on role
			if generator.PasswordPolicy == "" {
				generator.PasswordPolicy = dbConfig.PasswordPolicy
			}

			// Generate the password
			newPassword, err := generator.generate(ctx, b, dbi.database)
			if err != nil {
				b.CloseIfShutdown(dbi, err)
				return output, fmt.Errorf("failed to generate password: %s", err)
			}

			// Set new credential in WAL entry and update user request
			walEntry.NewPassword = newPassword
			updateReq.CredentialType = v5.CredentialTypePassword
			updateReq.Password = &v5.ChangePassword{
				NewPassword: newPassword,
				Statements:  statements,
			}

			// Set new credential in static account
			input.Role.StaticAccount.Password = newPassword
		case v5.CredentialTypeRSAPrivateKey:
			generator, err := newRSAKeyGenerator(input.Role.CredentialConfig)
			if err != nil {
				return output, fmt.Errorf("failed to construct credential generator: %s", err)
			}

			// Generate the RSA key pair
			public, private, err := generator.generate(b.GetRandomReader())
			if err != nil {
				return output, fmt.Errorf("failed to generate RSA key pair: %s", err)
			}

			// Set new credential in WAL entry and update user request
			walEntry.NewPublicKey = public
			updateReq.CredentialType = v5.CredentialTypeRSAPrivateKey
			updateReq.PublicKey = &v5.ChangePublicKey{
				NewPublicKey: public,
				Statements:   statements,
			}

			// Set new credential in static account
			input.Role.StaticAccount.PrivateKey = private
		}

		output.WALID, err = framework.PutWAL(ctx, s, staticWALKey, walEntry)
		if err != nil {
			return output, fmt.Errorf("error writing WAL entry: %w", err)
		}
		b.Logger().Debug("writing WAL", "role", input.RoleName, "WAL ID", output.WALID)
	}

	_, err = dbi.database.UpdateUser(ctx, updateReq, false)
	if err != nil {
		b.CloseIfShutdown(dbi, err)
		return output, fmt.Errorf("error setting credentials: %w", err)
	}

	// Store updated role information
	// lvr is the known LastVaultRotation
	lvr := time.Now()
	input.Role.StaticAccount.LastVaultRotation = lvr
	input.Role.StaticAccount.SetNextVaultRotation(lvr)
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
		b.Logger().Warn("error deleting WAL", "WAL ID", output.WALID, "error", err)
		return output, err
	}
	b.Logger().Debug("deleted WAL", "WAL ID", output.WALID)

	// The WAL has been deleted, return new setStaticAccountOutput without it
	return &setStaticAccountOutput{RotationTime: lvr}, nil
}

// initQueue preforms the necessary checks and initializations needed to perform
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
	if b.WriteSafeReplicationState() {
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
			if walID != "" && err == nil {
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
		queueTickerInterval := defaultQueueTickSeconds * time.Second
		if strVal, ok := conf.Config[queueTickIntervalKey]; ok {
			newVal, err := strconv.Atoi(strVal)
			if err == nil {
				queueTickerInterval = time.Duration(newVal) * time.Second
			} else {
				b.Logger().Error("bad value for %q option: %q", queueTickIntervalKey, strVal)
			}
		}
		go b.runTicker(ctx, queueTickerInterval, conf.StorageView)
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
			b.Logger().Debug("deleting WAL with nil role or static account", "WAL ID", walEntry.walID)
			if err := framework.DeleteWAL(ctx, s, walEntry.walID); err != nil {
				b.Logger().Warn("unable to delete WAL", "error", err, "WAL ID", walEntry.walID)
			}
			continue
		}

		if existingWALEntry, exists := walMap[walEntry.RoleName]; exists {
			b.Logger().Debug("multiple WALs detected for role", "role", walEntry.RoleName,
				"loaded WAL ID", existingWALEntry.walID, "created at", existingWALEntry.walCreatedAt, "last vault rotation", existingWALEntry.LastVaultRotation,
				"candidate WAL ID", walEntry.walID, "created at", walEntry.walCreatedAt, "last vault rotation", walEntry.LastVaultRotation)

			if walEntry.walCreatedAt > existingWALEntry.walCreatedAt {
				// If the existing WAL is older, delete it from storage and fall
				// through to inserting our current WAL into the map.
				b.Logger().Debug("deleting stale loaded WAL", "WAL ID", existingWALEntry.walID)
				err = framework.DeleteWAL(ctx, s, existingWALEntry.walID)
				if err != nil {
					b.Logger().Warn("unable to delete loaded WAL", "error", err, "WAL ID", existingWALEntry.walID)
				}
			} else {
				// If we already have a more recent WAL entry in the map, delete
				// this one and continue onto the next WAL.
				b.Logger().Debug("deleting stale candidate WAL", "WAL ID", walEntry.walID)
				err = framework.DeleteWAL(ctx, s, walID)
				if err != nil {
					b.Logger().Warn("unable to delete candidate WAL", "error", err, "WAL ID", walEntry.walID)
				}
				continue
			}
		}

		b.Logger().Debug("loaded WAL", "WAL ID", walID)
		walMap[walEntry.RoleName] = walEntry
	}
	return walMap, nil
}

// pushItem wraps the internal queue's Push call, to make sure a queue is
// actually available. This is needed because both runTicker and initQueue
// operate in go-routines, and could be accessing the queue concurrently
func (b *databaseBackend) pushItem(item *queue.Item) error {
	select {
	case <-b.queueCtx.Done():
	default:
		return b.credRotationQueue.Push(item)
	}
	b.Logger().Warn("no queue found during push item")
	return nil
}

// popFromRotationQueue wraps the internal queue's Pop call, to make sure a queue is
// actually available. This is needed because both runTicker and initQueue
// operate in go-routines, and could be accessing the queue concurrently
func (b *databaseBackend) popFromRotationQueue() (*queue.Item, error) {
	select {
	case <-b.queueCtx.Done():
	default:
		return b.credRotationQueue.Pop()
	}
	return nil, queue.ErrEmpty
}

// popFromRotationQueueByKey wraps the internal queue's PopByKey call, to make sure a queue is
// actually available. This is needed because both runTicker and initQueue
// operate in go-routines, and could be accessing the queue concurrently
func (b *databaseBackend) popFromRotationQueueByKey(name string) (*queue.Item, error) {
	select {
	case <-b.queueCtx.Done():
	default:
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

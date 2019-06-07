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
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	// interval to check the queue for items needing rotation
	queueTickSeconds  = 5
	queueTickInterval = queueTickSeconds * time.Second

	// wal storage key used for static account rotations
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
				// no static role exists for this wal entry, delete it
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
			b.logger.Error("unable to rotate credentials in periodic function", "error", err)
			// Increment the priority enough so that the next call to this method
			// likely will not attempt to rotate it, as a back-off of sorts
			item.Priority = time.Now().Add(10 * time.Second).Unix()

			// preserve the WALID if it was returned
			if resp != nil && resp.WALID != "" {
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
	if input == nil || input.Role == nil || input.RoleName == "" {
		return nil, errors.New("input was empty when attempting to set credentials for static account")
	}
	// re-use WAL ID if present, otherwise PUT a new WAL
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

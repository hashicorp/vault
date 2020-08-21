package database

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WAL storage key used for the rollback of root database credentials
const rotateRootWALKey = "rotateRootWALKey"

// WAL entry used for the rollback of root database credentials
type rotateRootCredentialsWAL struct {
	ConnectionName string
	UserName       string
	NewPassword    string
	OldPassword    string
}

// walRollback handles WAL entries that result from partial failures
// to rotate the root credentials of a database. It is responsible
// for rolling back root database credentials when doing so would
// reconcile the credentials with Vault storage.
func (b *databaseBackend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
	if kind != rotateRootWALKey {
		return errors.New("unknown type to rollback")
	}

	// Decode the WAL data
	var entry rotateRootCredentialsWAL
	if err := mapstructure.Decode(data, &entry); err != nil {
		return err
	}

	// Get the current database configuration from storage
	config, err := b.DatabaseConfig(ctx, req.Storage, entry.ConnectionName)
	if err != nil {
		return err
	}

	// The password in storage doesn't match the new password
	// in the WAL entry. This means there was a partial failure
	// to update either the database or storage.
	if config.ConnectionDetails["password"] != entry.NewPassword {
		// Clear any cached connection to inform the rollback decision
		err := b.ClearConnection(entry.ConnectionName)
		if err != nil {
			return err
		}

		// Attempt to get a connection with the current configuration.
		// If successful, the WAL entry can be deleted. This means
		// the root credentials are the same according to the database
		// and storage.
		_, err = b.GetConnection(ctx, req.Storage, entry.ConnectionName)
		if err == nil {
			return nil
		}

		return b.rollbackDatabaseCredentials(ctx, config, entry)
	}

	// The password in storage matches the new password in
	// the WAL entry, so there is nothing to roll back. This
	// means the new password was successfully updated in the
	// database and storage, but the WAL wasn't deleted.
	return nil
}

// rollbackDatabaseCredentials rolls back root database credentials for
// the connection associated with the passed WAL entry. It will create
// a connection to the database using the WAL entry new password in
// order to alter the password to be the WAL entry old password.
func (b *databaseBackend) rollbackDatabaseCredentials(ctx context.Context, config *DatabaseConfig, entry rotateRootCredentialsWAL) error {
	// Attempt to get a connection with the WAL entry new password.
	config.ConnectionDetails["password"] = entry.NewPassword
	db, err := b.GetConnectionWithConfig(ctx, entry.ConnectionName, config)
	if err != nil {
		return err
	}

	// Ensure the connection used to roll back the database password is not cached
	defer func() {
		if err := b.ClearConnection(entry.ConnectionName); err != nil {
			b.Logger().Error("error closing database plugin connection", "err", err)
		}
	}()

	err = changeUserPassword(ctx, db.database, entry.UserName, entry.OldPassword, config.RootCredentialsRotateStatements)
	if status.Code(err) == codes.Unimplemented {
		return nil
	}
	return err
}

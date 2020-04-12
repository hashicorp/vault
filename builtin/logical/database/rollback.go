package database

import (
	"context"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

// WAL storage key used for root credential rotations on database plugins
const rootWALKey = "rootRotationKey"

type rotateRootCredentialsWAL struct {
	ConnectionName string
	UserName       string
	NewPassword    string
	OldPassword    string
}

// TODO: Considerations for HA and Replication?
func (b *databaseBackend) walRollback(ctx context.Context, req *logical.Request, kind string,
	data interface{}) error {
	if kind != rootWALKey {
		return nil
	}

	// Decode the WAL data
	var entry rotateRootCredentialsWAL
	if err := mapstructure.Decode(data, &entry); err != nil {
		b.Logger().Info("error decoding WAL data", "data", data)
		return err
	}

	// Get the current database configuration from Vault storage
	config, err := b.DatabaseConfig(ctx, req.Storage, entry.ConnectionName)
	if err != nil {
		return err
	}

	// The password in Vault storage does not match the new password
	// in the WAL entry. This means there was a partial failure where
	// the database password was updated but Vault storage was not.
	// To reconcile the password between Vault and the database, roll
	// back the database password to the old password.
	if config.ConnectionDetails["password"] != entry.NewPassword {
		return b.rollbackDatabasePassword(ctx, config, entry)
	}

	// The password in Vault storage matches the new password
	// in the WAL entry, so there is nothing to roll back. This
	// means the new password was successfully updated in the
	// database and Vault storage, but the WAL was not deleted.
	return nil
}

func (b *databaseBackend) rollbackDatabasePassword(ctx context.Context, config *DatabaseConfig, entry rotateRootCredentialsWAL) error {
	// Get a connection using the new password
	config.ConnectionDetails["password"] = entry.NewPassword
	dbc, err := b.GetConnectionWithConfig(ctx, entry.ConnectionName, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := b.ClearConnection(entry.ConnectionName); err != nil {
			b.Logger().Error("error closing database plugin connection", "err", err)
		}
	}()

	// Roll back the database password to the WAL old password
	// in order to reconcile the database and Vault storage.
	rotationStatements := dbplugin.Statements{
		Rotation: config.RootCredentialsRotateStatements,
	}
	userConfig := dbplugin.StaticUserConfig{
		Username: entry.UserName,
		Password: entry.OldPassword,
	}
	_, _, err = dbc.SetCredentials(ctx, rotationStatements, userConfig)

	// If SetCredentials is unimplemented in the plugin, this means that
	// the root credential rotation happened via the RotateRootCredentials
	// RPC. Delete the WAL by returning nil.
	if err != nil && status.Code(err) == codes.Unimplemented {
		return nil
	}

	return err
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"bytes"
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WAL storage key used for the rollback of root database credentials
const (
	rotateRootWALKey       = "rotateRootWALKey"
	rollbackTypePassword   = "password"
	rollbackTypePrivateKey = "private_key"
)

// WAL entry used for the rollback of root database credentials
type rotateRootCredentialsWAL struct {
	ConnectionName string
	UserName       string
	NewPassword    string
	OldPassword    string

	// NewPrivateKey is used to update the config value
	NewPrivateKey []byte
	// OldPublicKey is used in the UpdateUserRequest
	OldPublicKey []byte
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

	// default rollback type
	isPrivatekeyRollback := config.ConnectionDetails["private_key"] != nil && !bytes.Equal(config.ConnectionDetails["private_key"].([]byte), entry.NewPrivateKey)
	isPasswordRollback := config.ConnectionDetails["password"] != "" && config.ConnectionDetails["password"] != entry.NewPassword

	// The credential in storage doesn't match the new credential
	// in the WAL entry. This means there was a partial failure
	// to update either the database or storage.
	if isPasswordRollback || isPrivatekeyRollback {
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

		switch {
		case isPrivatekeyRollback:
			return b.rollbackDatabaseCredentials(ctx, config, entry, rollbackTypePrivateKey)
		case isPasswordRollback:
			return b.rollbackDatabaseCredentials(ctx, config, entry, rollbackTypePassword)
		default:
		}
	}

	// The credential in storage matches the new password in
	// the WAL entry, so there is nothing to roll back. This
	// means the new credential was successfully updated in the
	// database and storage, but the WAL wasn't deleted.
	return nil
}

// rollbackDatabaseCredentials rolls back root database credentials for
// the connection associated with the passed WAL entry. It will create
// a connection to the database using the WAL entry new password in
// order to alter the password to be the WAL entry old password.
func (b *databaseBackend) rollbackDatabaseCredentials(ctx context.Context, config *DatabaseConfig, entry rotateRootCredentialsWAL, rollbackType string) error {
	// Attempt to get a connection with the WAL entry new password.

	var dbi *dbPluginInstance
	var updateReq v5.UpdateUserRequest
	var err error
	switch rollbackType {
	case rollbackTypePassword:
		config.ConnectionDetails["password"] = entry.NewPassword
		dbi, err = b.GetConnectionWithConfig(ctx, entry.ConnectionName, config)
		if err != nil {
			return err
		}

		// Ensure the connection used to roll back the database password is not cached
		defer func() {
			if err := b.ClearConnection(entry.ConnectionName); err != nil {
				b.Logger().Error("error closing database plugin connection", "err", err)
			}
		}()

		updateReq = v5.UpdateUserRequest{
			Username:       entry.UserName,
			CredentialType: v5.CredentialTypePassword,
			Password: &v5.ChangePassword{
				NewPassword: entry.OldPassword,
				Statements: v5.Statements{
					Commands: config.RootCredentialsRotateStatements,
				},
			},
		}

	case rollbackTypePrivateKey:
		config.ConnectionDetails["private_key"] = entry.NewPrivateKey
		dbi, err = b.GetConnectionWithConfig(ctx, entry.ConnectionName, config)
		if err != nil {
			return err
		}

		// Ensure the connection used to roll back the database password is not cached
		defer func() {
			if err := b.ClearConnection(entry.ConnectionName); err != nil {
				b.Logger().Error("error closing database plugin connection", "err", err)
			}
		}()

		updateReq = v5.UpdateUserRequest{
			Username:       entry.UserName,
			CredentialType: v5.CredentialTypeRSAPrivateKey,
			PublicKey: &v5.ChangePublicKey{
				NewPublicKey: entry.OldPublicKey,
				Statements: v5.Statements{
					Commands: config.RootCredentialsRotateStatements,
				},
			},
		}
	default:
		return errors.New("unknown rollback type")

	}

	// It actually is the root user here, but we only want to use SetCredentials since
	// RotateRootCredentials doesn't give any control over what password is used
	_, err = dbi.database.UpdateUser(ctx, updateReq, false)
	if status.Code(err) == codes.Unimplemented || err == dbplugin.ErrPluginStaticUnsupported {
		return nil
	}
	return err
}

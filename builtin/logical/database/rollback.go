// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WAL storage key used for the rollback of root database credentials
const rotateRootWALKey = "rotateRootWALKey"

// snowflakeErrJWTTokenInvalid is the Snowflake server-side error code for JWT
// authentication failure.
const snowflakeErrJWTTokenInvalid = "390144"

// WAL entry used for the rollback of root database credentials
type rotateRootCredentialsWAL struct {
	ConnectionName string
	UserName       string

	NewPassword string
	OldPassword string

	NewPublicKey  string
	NewPrivateKey string
	OldPrivateKey string
}

func NewRotateRootCredentialsWALPasswordEntry(connectionName, userName, newPassword, oldPassword string) *rotateRootCredentialsWAL {
	return &rotateRootCredentialsWAL{
		ConnectionName: connectionName,
		UserName:       userName,
		NewPassword:    newPassword,
		OldPassword:    oldPassword,
	}
}

func NewRotateRootCredentialsWALPrivateKeyEntry(connectionName, userName, newPublicKey, newPrivateKey, oldPrivateKey string) *rotateRootCredentialsWAL {
	return &rotateRootCredentialsWAL{
		ConnectionName: connectionName,
		UserName:       userName,
		NewPublicKey:   newPublicKey,
		NewPrivateKey:  newPrivateKey,
		OldPrivateKey:  oldPrivateKey,
	}
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

	// Route based on credential type in the WAL entry.
	if entry.NewPrivateKey != "" {
		// Stored key matches WAL new key: rotation completed, WAL not yet deleted.
		if config.ConnectionDetails["private_key"] == entry.NewPrivateKey {
			b.Logger().Info("WAL rollback: private key already rotated, nothing to roll back",
				"connection", entry.ConnectionName)
			return nil
		}

		b.Logger().Warn("WAL rollback: private key out of sync, starting rollback",
			"connection", entry.ConnectionName, "username", entry.UserName)

		if err := b.ClearConnection(entry.ConnectionName); err != nil {
			return err
		}

		return b.rollbackDatabasePrivateKey(ctx, config, entry)
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

		// An initialization timeout means the database was unreachable within
		// Vault's deadline, not that the stored credentials are wrong. A timeout
		// is not a reliable signal of credential state: the rotation may have
		// already applied the new password before the database became slow.
		// Returning the error here lets the WAL framework retry later rather
		// than risking a rollback that reverts a successfully rotated credential.
		if errors.Is(err, errDatabaseInitializeTimeout) {
			return err
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
	dbi, err := b.GetConnectionWithConfig(ctx, entry.ConnectionName, config)
	if err != nil {
		return err
	}

	// Ensure the connection used to roll back the database password is not cached
	defer func() {
		if err := b.ClearConnection(entry.ConnectionName); err != nil {
			b.Logger().Error("error closing database plugin connection", "err", err)
		}
	}()

	updateReq := v5.UpdateUserRequest{
		Username:       entry.UserName,
		CredentialType: v5.CredentialTypePassword,
		Password: &v5.ChangePassword{
			NewPassword: entry.OldPassword,
			Statements: v5.Statements{
				Commands: config.RootCredentialsRotateStatements,
			},
		},
	}

	// It actually is the root user here, but we only want to use SetCredentials since
	// RotateRootCredentials doesn't give any control over what password is used
	_, err = dbi.database.UpdateUser(ctx, updateReq, false)
	if status.Code(err) == codes.Unimplemented || err == dbplugin.ErrPluginStaticUnsupported {
		return nil
	}
	return err
}

// rollbackDatabasePrivateKey restores the old public key on Snowflake for key-pair
// auth connections by connecting with the new private key and issuing ALTER USER.
func (b *databaseBackend) rollbackDatabasePrivateKey(ctx context.Context, config *DatabaseConfig, entry rotateRootCredentialsWAL) error {
	oldPublicKey, err := derivePublicKeyFromPrivateKeyPEM(entry.OldPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to derive old public key for rollback: %w", err)
	}

	config.ConnectionDetails["private_key"] = entry.NewPrivateKey
	dbi, err := b.GetConnectionWithConfig(ctx, entry.ConnectionName, config)
	if err != nil {
		b.Logger().Error("WAL rollback: failed to connect using new private key", "connection", entry.ConnectionName, "error", err.Error())
		return err
	}

	defer func() {
		if err := b.ClearConnection(entry.ConnectionName); err != nil {
			b.Logger().Error("error closing database plugin connection", "error", err)
		}
	}()

	b.Logger().Info("WAL rollback: restoring old public key on Snowflake", "connection", entry.ConnectionName, "username", entry.UserName)

	updateReq := v5.UpdateUserRequest{
		Username:       entry.UserName,
		CredentialType: v5.CredentialTypeRSAPrivateKey,
		PublicKey: &v5.ChangePublicKey{
			NewPublicKey: oldPublicKey,
			Statements: v5.Statements{
				Commands: config.RootCredentialsRotateStatements,
			},
		},
	}

	_, err = dbi.database.UpdateUser(ctx, updateReq, false)
	if status.Code(err) == codes.Unimplemented || err == dbplugin.ErrPluginStaticUnsupported {
		return nil
	}
	if err != nil {
		// Snowflake error 390144 means JWT authentication failed. This occurs when
		// the new private key was never registered with Snowflake (crash before UpdateUser),
		// so the system is already consistent with the old key — delete the WAL cleanly.
		if strings.Contains(err.Error(), snowflakeErrJWTTokenInvalid) {
			b.Logger().Info("WAL rollback: new private key rejected by Snowflake (crash before UpdateUser), system already consistent",
				"connection", entry.ConnectionName)
			return nil
		}
		b.Logger().Error("WAL rollback: failed to restore old public key", "connection", entry.ConnectionName, "error", err.Error())
		return err
	}
	b.Logger().Info("WAL rollback: successfully restored old public key", "connection", entry.ConnectionName, "username", entry.UserName)
	return nil
}

func derivePublicKeyFromPrivateKeyPEM(privateKeyPEM string) ([]byte, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an RSA key")
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}), nil
}

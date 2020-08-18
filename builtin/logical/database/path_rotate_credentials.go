package database

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/vault/sdk/database/newdbplugin"

	"github.com/hashicorp/vault/sdk/database/dbplugin"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

func pathRotateCredentials(b *databaseBackend) []*framework.Path {
	return []*framework.Path{
		&framework.Path{
			Pattern: "rotate-root/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of this database connection",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateCredentialsUpdate(),
					ForwardPerformanceSecondary: true,
					ForwardPerformanceStandby:   true,
				},
			},

			HelpSynopsis:    pathCredsCreateReadHelpSyn,
			HelpDescription: pathCredsCreateReadHelpDesc,
		},
		&framework.Path{
			Pattern: "rotate-role/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the static role",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateRoleCredentialsUpdate(),
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},

			HelpSynopsis:    pathCredsCreateReadHelpSyn,
			HelpDescription: pathCredsCreateReadHelpDesc,
		},
	}
}

func (b *databaseBackend) pathRotateCredentialsUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		config, err := b.DatabaseConfig(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}

		db, err := b.GetConnection(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}

		defer func() {
			// Close the plugin
			db.closed = true
			if err := db.Close(); err != nil {
				b.Logger().Error("error closing the database plugin connection", "err", err)
			}
			// Even on error, still remove the connection
			delete(b.connections, name)
		}()

		// Take out the backend lock since we are swapping out the connection
		b.Lock()
		defer b.Unlock()

		// Take the write lock on the instance
		db.Lock()
		defer db.Unlock()

		// Generate new credentials
		username := config.ConnectionDetails["username"].(string)
		oldPassword := config.ConnectionDetails["password"].(string)
		newPassword, err := generatePassword(ctx, b.System(), config.PasswordPolicy)
		if err != nil {
			return nil, err
		}
		config.ConnectionDetails["password"] = newPassword

		// Write a WAL entry
		walID, err := framework.PutWAL(ctx, req.Storage, rotateRootWALKey, &rotateRootCredentialsWAL{
			ConnectionName: name,
			UserName:       username,
			OldPassword:    oldPassword,
			NewPassword:    newPassword,
		})
		if err != nil {
			return nil, err
		}

		// Database v5
		if db.database.database != nil {
			err := b.changeUserPasswordNew(ctx, db.database, username, newPassword, config.RootCredentialsRotateStatements)
			if err != nil {
				return nil, fmt.Errorf("failed to change root user password: %w", err)
			}
			b.deleteWal(ctx, req.Storage, walID)
			return nil, nil
		}

		// Database v4
		err = b.changeUserPasswordLegacy(ctx, db.database, config.RootCredentialsRotateStatements, username, newPassword)
		if status.Code(err) == codes.Unimplemented {
			// Fall back to using RotateRootCredentials if unimplemented
			newConfigDetails, err := db.database.legacyDatabase.RotateRootCredentials(ctx, config.RootCredentialsRotateStatements)
			if err != nil {
				return nil, err
			}
			config.ConnectionDetails = newConfigDetails
			err = store(ctx, req.Storage, name, config)
			if err != nil {
				return nil, err
			}
			b.deleteWal(ctx, req.Storage, walID)
			return nil, nil
		}
		if err != nil {
			return nil, err
		}

		// v4: SetCredentials worked
		b.deleteWal(ctx, req.Storage, walID)
		return nil, nil
	}
}

func (b *databaseBackend) deleteWal(ctx context.Context, storage logical.Storage, walID string) {
	err := framework.DeleteWAL(ctx, storage, walID)
	if err != nil {
		b.Logger().Warn("unable to delete WAL", "error", err, "WAL ID", walID)
	}
}

func (b *databaseBackend) changeUserPasswordNew(ctx context.Context, dbw databaseVersionWrapper, username, newPassword string, rotationCmds []string) error {
	req := newdbplugin.UpdateUserRequest{
		Username: username,
		Password: &newdbplugin.ChangePassword{
			NewPassword: newPassword,
			Statements: newdbplugin.Statements{
				Commands: rotationCmds,
			},
		},
	}
	_, err := dbw.database.UpdateUser(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (b *databaseBackend) changeUserPasswordLegacy(
	ctx context.Context,
	dbw databaseVersionWrapper,
	rotationCmds []string,
	username, newPassword string) (err error) {

	// Attempt to use SetCredentials for the root credential rotation
	statements := dbplugin.Statements{Rotation: rotationCmds}
	userConfig := dbplugin.StaticUserConfig{
		Username: username,
		Password: newPassword,
	}

	_, _, err = dbw.legacyDatabase.SetCredentials(ctx, statements, userConfig)
	return err
}

func store(ctx context.Context, storage logical.Storage, name string, value interface{}) error {
	entry, err := logical.StorageEntryJSON(fmt.Sprintf("config/%s", name), value)
	if err != nil {
		return fmt.Errorf("unable to marshal object to JSON: %w", err)
	}

	err = storage.Put(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to save object: %w", err)
	}
	return nil
}

// func (b *databaseBackend) rotateRootCredentials(
// 	ctx context.Context,
// 	dbw databaseVersionWrapper,
// 	storage putter,
// 	name string,
// 	config *DatabaseConfig,
// 	rotationCmds []string) (newConfig map[string]interface{}, err error) {
//
// 	newConfig, err := dbw.legacyDatabase.RotateRootCredentials(ctx, rotationCmds)
// 	if err != nil {
// 		return err
// 	}
// 	config.ConnectionDetails = newConfig
//
// 	// Update storage with the new root credentials
// 	entry, err := logical.StorageEntryJSON(fmt.Sprintf("config/%s", name), config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := req.Storage.Put(ctx, entry); err != nil {
// 		return nil, err
// 	}
//
// 	// config.ConnectionDetails, err := dbw.legacyDatabase.RotateRootCredentials(ctx, rotationCmds)
// 	// if err != nil {
// 	// 	return err
// 	// }
// }

func (b *databaseBackend) pathRotateRoleCredentialsUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse("empty role name attribute given"), nil
		}

		role, err := b.StaticRole(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return logical.ErrorResponse("no static role found for role name"), nil
		}

		// In create/update of static accounts, we only care if the operation
		// err'd , and this call does not return credentials
		item, err := b.popFromRotationQueueByKey(name)
		if err != nil {
			item = &queue.Item{
				Key: name,
			}
		}

		resp, err := b.setStaticAccount(ctx, req.Storage, &setStaticAccountInput{
			RoleName: name,
			Role:     role,
		})
		// if err is not nil, we need to attempt to update the priority and place
		// this item back on the queue. The err should still be returned at the end
		// of this method.
		if err != nil {
			b.logger.Warn("unable to rotate credentials in rotate-role", "error", err)
			// Update the priority to re-try this rotation and re-add the item to
			// the queue
			item.Priority = time.Now().Add(10 * time.Second).Unix()

			// Preserve the WALID if it was returned
			if resp != nil && resp.WALID != "" {
				item.Value = resp.WALID
			}
		} else {
			item.Priority = resp.RotationTime.Add(role.StaticAccount.RotationPeriod).Unix()
		}

		// Add their rotation to the queue
		if err := b.pushItem(item); err != nil {
			return nil, err
		}

		// return any err from the setStaticAccount call
		return nil, err
	}
}

const pathRotateCredentialsUpdateHelpSyn = `
Request to rotate the root credentials for a certain database connection.
`

const pathRotateCredentialsUpdateHelpDesc = `
This path attempts to rotate the root credentials for the given database. 
`

const pathRotateRoleCredentialsUpdateHelpSyn = `
Request to rotate the credentials for a static user account.
`
const pathRotateRoleCredentialsUpdateHelpDesc = `
This path attempts to rotate the credentials for the given static user account.
`

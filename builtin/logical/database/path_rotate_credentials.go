package database

import (
	"context"
	"fmt"
	"time"

	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

func pathRotateRootCredentials(b *databaseBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "rotate-root/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of this database connection",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateRootCredentialsUpdate(),
					ForwardPerformanceSecondary: true,
					ForwardPerformanceStandby:   true,
				},
			},

			HelpSynopsis:    pathRotateCredentialsUpdateHelpSyn,
			HelpDescription: pathRotateCredentialsUpdateHelpDesc,
		},
		{
			Pattern: "rotate-role/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
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

			HelpSynopsis:    pathRotateRoleCredentialsUpdateHelpSyn,
			HelpDescription: pathRotateRoleCredentialsUpdateHelpDesc,
		},
	}
}

func (b *databaseBackend) pathRotateRootCredentialsUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse(respErrEmptyName), nil
		}

		config, err := b.DatabaseConfig(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}

		rootUsername, ok := config.ConnectionDetails["username"].(string)
		if !ok || rootUsername == "" {
			return nil, fmt.Errorf("unable to rotate root credentials: no username in configuration")
		}

		dbi, err := b.GetConnection(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}

		// Take out the backend lock since we are swapping out the connection
		b.Lock()
		defer b.Unlock()

		// Take the write lock on the instance
		dbi.Lock()
		defer dbi.Unlock()

		defer func() {
			// Close the plugin
			dbi.closed = true
			if err := dbi.database.Close(); err != nil {
				b.Logger().Error("error closing the database plugin connection", "err", err)
			}
			// Even on error, still remove the connection
			delete(b.connections, name)
		}()

		generator, err := newPasswordGenerator(nil)
		if err != nil {
			return nil, fmt.Errorf("failed to construct credential generator: %s", err)
		}
		generator.PasswordPolicy = config.PasswordPolicy

		// Generate new credentials
		oldPassword := config.ConnectionDetails["password"].(string)
		newPassword, err := generator.generate(ctx, b, dbi.database)
		if err != nil {
			b.CloseIfShutdown(dbi, err)
			return nil, fmt.Errorf("failed to generate password: %s", err)
		}
		config.ConnectionDetails["password"] = newPassword

		// Write a WAL entry
		walID, err := framework.PutWAL(ctx, req.Storage, rotateRootWALKey, &rotateRootCredentialsWAL{
			ConnectionName: name,
			UserName:       rootUsername,
			OldPassword:    oldPassword,
			NewPassword:    newPassword,
		})
		if err != nil {
			return nil, err
		}

		updateReq := v5.UpdateUserRequest{
			Username:       rootUsername,
			CredentialType: v5.CredentialTypePassword,
			Password: &v5.ChangePassword{
				NewPassword: newPassword,
				Statements: v5.Statements{
					Commands: config.RootCredentialsRotateStatements,
				},
			},
		}
		newConfigDetails, err := dbi.database.UpdateUser(ctx, updateReq, true)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
		if newConfigDetails != nil {
			config.ConnectionDetails = newConfigDetails
		}

		err = storeConfig(ctx, req.Storage, name, config)
		if err != nil {
			return nil, err
		}

		err = framework.DeleteWAL(ctx, req.Storage, walID)
		if err != nil {
			b.Logger().Warn("unable to delete WAL", "error", err, "WAL ID", walID)
		}
		return nil, nil
	}
}

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

		input := &setStaticAccountInput{
			RoleName: name,
			Role:     role,
		}
		if walID, ok := item.Value.(string); ok {
			input.WALID = walID
		}
		resp, err := b.setStaticAccount(ctx, req.Storage, input)
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
			// Clear any stored WAL ID as we must have successfully deleted our WAL to get here.
			item.Value = ""
		}

		// Add their rotation to the queue
		if err := b.pushItem(item); err != nil {
			return nil, err
		}

		if err != nil {
			return nil, fmt.Errorf("unable to finish rotating credentials; retries will "+
				"continue in the background but it is also safe to retry manually: %w", err)
		}

		// return any err from the setStaticAccount call
		return nil, nil
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

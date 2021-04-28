package openldap

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	rotateRootPath = "rotate-root"
	rotateRolePath = "rotate-role/"
)

func (b *backend) pathRotateCredentials() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: rotateRootPath,
			Fields:  map[string]*framework.FieldSchema{},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateCredentialsUpdate,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateCredentialsUpdate,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},
			HelpSynopsis: "Request to rotate the root credentials Vault uses for the OpenLDAP administrator account.",
			HelpDescription: "This path attempts to rotate the root credentials of the administrator account " +
				"(binddn) used by Vault to manage OpenLDAP.",
		},
		{
			Pattern: rotateRolePath + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the static role",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateRoleCredentialsUpdate,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateRoleCredentialsUpdate,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},
			HelpSynopsis:    "Request to rotate the credentials for a static user account.",
			HelpDescription: "This path attempts to rotate the credentials for the given OpenLDAP static user account.",
		},
	}
}

func (b *backend) pathRotateCredentialsUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if _, hasTimeout := ctx.Deadline(); !hasTimeout {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, defaultCtxTimeout)
		defer cancel()
	}

	config, err := readConfig(ctx, req.Storage)
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
	oldPassword := config.LDAP.BindPassword

	// Take out the backend lock since we are swapping out the connection
	b.Lock()
	defer b.Unlock()

	// Update the password remotely.
	if err := b.client.UpdateRootPassword(config.LDAP, newPassword); err != nil {
		return nil, err
	}
	config.LDAP.BindPassword = newPassword

	// Update the password locally.
	if pwdStoringErr := storePassword(ctx, req.Storage, config); pwdStoringErr != nil {
		// We were unable to store the new password locally. We can't continue in this state because we won't be able
		// to roll any passwords, including our own to get back into a state of working. So, we need to roll back to
		// the last password we successfully got into storage.
		if rollbackErr := b.rollBackPassword(ctx, config, oldPassword); rollbackErr != nil {
			return nil, fmt.Errorf(`unable to store new password due to %s and unable to return to previous password
due to %s, configure a new binddn and bindpass to restore openldap function`, pwdStoringErr, rollbackErr)
		}
		return nil, fmt.Errorf("unable to update password due to storage err: %s", pwdStoringErr)
	}

	// Respond with a 204.
	return nil, nil
}
func (b *backend) pathRotateRoleCredentialsUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("empty role name attribute given"), nil
	}

	role, err := b.staticRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("role doesn't exist: %s", name), nil
	}

	// In create/update of static accounts, we only care if the operation
	// err'd , and this call does not return credentials
	item, err := b.popFromRotationQueueByKey(name)
	if err != nil {
		item = &queue.Item{
			Key: name,
		}
	}

	resp, err := b.setStaticAccountPassword(ctx, req.Storage, &setStaticAccountInput{
		RoleName: name,
		Role:     role,
	})
	if err != nil {
		b.Logger().Warn("unable to rotate credentials in rotate-role", "error", err)
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

	// Add their rotation to the queue. We use pushErr here to distinguish between
	// the error returned from setStaticAccount. They are scoped differently but
	// it's more clear to developers that err above can still be non nil, and not
	// overwritten or reused here.
	if pushErr := b.pushItem(item); pushErr != nil {
		return nil, pushErr
	}

	// We're not returning creds here because we do not know if its been processed
	// by the queue.
	return nil, err
}

// rollBackPassword uses naive exponential backoff to retry updating to an old password,
// because LDAP may still be propagating the previous password change.
func (b *backend) rollBackPassword(ctx context.Context, config *config, oldPassword string) error {
	var err error
	for i := 0; i < 10; i++ {
		waitSeconds := math.Pow(float64(i), 2)
		timer := time.NewTimer(time.Duration(waitSeconds) * time.Second)
		select {
		case <-timer.C:
		case <-ctx.Done():
			// Outer environment is closing.
			return fmt.Errorf("unable to roll back password because enclosing environment is shutting down")
		}
		if err = b.client.UpdateRootPassword(config.LDAP, oldPassword); err == nil {
			// Success.
			return nil
		}
	}
	// Failure after looping.
	return err
}

func storePassword(ctx context.Context, s logical.Storage, config *config) error {
	entry, err := logical.StorageEntryJSON(configPath, config)
	if err != nil {
		return err
	}
	return s.Put(ctx, entry)
}

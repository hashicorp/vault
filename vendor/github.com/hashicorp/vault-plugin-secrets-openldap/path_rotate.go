// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openldap

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/backoff"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

var (
	rollbackAttempts    = 10
	minRollbackDuration = 1 * time.Second
	maxRollbackDuration = 100 * time.Second
)

const (
	rotateRootPath = "rotate-root"
	rotateRolePath = "rotate-role/"
)

func (b *backend) pathRotateCredentials() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: rotateRootPath,
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAP,
				OperationVerb:   "rotate",
				OperationSuffix: "root-credentials",
			},
			Fields: map[string]*framework.FieldSchema{},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateRootCredentialsUpdate,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},
			HelpSynopsis: "Request to rotate the root credentials Vault uses for the LDAP administrator account.",
			HelpDescription: "This path attempts to rotate the root credentials of the administrator account " +
				"(binddn) used by Vault to manage LDAP.",
		},
		{
			Pattern: strings.TrimSuffix(rotateRolePath, "/") + genericNameWithForwardSlashRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAP,
				OperationVerb:   "rotate",
				OperationSuffix: "static-role",
			},
			Fields: fieldsForType(rotateRolePath),
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathRotateRoleCredentialsUpdate,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},
			HelpSynopsis:    "Request to rotate the credentials for a static user account.",
			HelpDescription: "This path attempts to rotate the credentials for the given LDAP static user account.",
		},
	}
}

func (b *backend) pathRotateRootCredentialsUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
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
	if err := b.client.UpdateDNPassword(config.LDAP, config.LDAP.BindDN, newPassword); err != nil {
		return nil, err
	}
	config.LDAP.BindPassword = newPassword
	config.LDAP.LastBindPassword = oldPassword
	config.LDAP.LastBindPasswordRotation = time.Now()

	// Update the password locally.
	if pwdStoringErr := storePassword(ctx, req.Storage, config); pwdStoringErr != nil {
		// We were unable to store the new password locally. We can't continue in this state because we won't be able
		// to roll any passwords, including our own to get back into a state of working. So, we need to roll back to
		// the last password we successfully got into storage.
		if rollbackErr := b.rollbackPassword(ctx, config, oldPassword); rollbackErr != nil {
			return nil, fmt.Errorf(`unable to store new password due to %s and unable to return to previous password
due to %s, configure a new binddn and bindpass to restore ldap function`, pwdStoringErr, rollbackErr)
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

	input := &setStaticAccountInput{
		RoleName: name,
		Role:     role,
	}
	if walID, ok := item.Value.(string); ok {
		input.WALID = walID
	}
	resp, err := b.setStaticAccountPassword(ctx, req.Storage, input)
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
		// Clear any stored WAL ID as we must have successfully deleted our WAL to get here.
		item.Value = ""
	}

	// Add their rotation to the queue. We use pushErr here to distinguish between
	// the error returned from setStaticAccount. They are scoped differently but
	// it's more clear to developers that err above can still be non nil, and not
	// overwritten or reused here.
	if pushErr := b.pushItem(item); pushErr != nil {
		return nil, pushErr
	}

	if err != nil {
		return nil, fmt.Errorf("unable to finish rotating credentials; retries will "+
			"continue in the background but it is also safe to retry manually: %w", err)
	}

	// We're not returning creds here because we do not know if its been processed
	// by the queue.
	return nil, nil
}

// rollbackPassword uses exponential backoff to retry updating to an old password,
// because LDAP may still be propagating the previous password change.
func (b *backend) rollbackPassword(ctx context.Context, config *config, oldPassword string) error {
	expbackoff := backoff.NewBackoff(rollbackAttempts, minRollbackDuration, maxRollbackDuration)
	var err error
	for {
		nextsleep, terr := expbackoff.Next()
		if terr != nil {
			// exponential backoff has failed every attempt; return last error
			return err
		}
		timer := time.NewTimer(nextsleep)
		select {
		case <-timer.C:
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C // drain the channel so that it will be garbage collected
			}
			// Outer environment is closing.
			return fmt.Errorf("unable to rollback password because enclosing environment is shutting down")
		}
		err = b.client.UpdateDNPassword(config.LDAP, config.LDAP.BindDN, oldPassword)
		if err == nil {
			return nil
		}
	}
}

func storePassword(ctx context.Context, s logical.Storage, config *config) error {
	entry, err := logical.StorageEntryJSON(configPath, config)
	if err != nil {
		return err
	}
	return s.Put(ctx, entry)
}

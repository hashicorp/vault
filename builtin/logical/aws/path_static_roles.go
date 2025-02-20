// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	pathStaticRole = "static-roles"

	paramRoleName        = "name"
	paramUsername        = "username"
	paramRotationPeriod  = "rotation_period"
	paramAssumeRoleARN   = "assume_role_arn"
	paramRoleSessionName = "assume_role_session_name"
	paramExternalID      = "external_id"
)

type staticRoleEntry struct {
	Name                  string        `json:"name" structs:"name" mapstructure:"name"`
	ID                    string        `json:"id" structs:"id" mapstructure:"id"`
	Username              string        `json:"username" structs:"username" mapstructure:"username"`
	RotationPeriod        time.Duration `json:"rotation_period" structs:"rotation_period" mapstructure:"rotation_period"`
	AssumeRoleARN         string        `json:"assume_role_arn" structs:"assume_role_arn" mapstructure:"assume_role_arn"`
	AssumeRoleSessionName string        `json:"assume_role_session_name" structs:"assume_role_session_name" mapstructure:"assume_role_session_name"`
	ExternalID            string        `json:"external_id" structs:"external_id" mapstructure:"external_id"`
}

func pathStaticRoles(b *backend) *framework.Path {
	roleResponse := map[int][]framework.Response{
		http.StatusOK: {{
			Description: http.StatusText(http.StatusOK),
			Fields: map[string]*framework.FieldSchema{
				paramRoleName: {
					Type:        framework.TypeString,
					Description: descRoleName,
				},
				paramUsername: {
					Type:        framework.TypeString,
					Description: descUsername,
				},
				paramRotationPeriod: {
					Type:        framework.TypeDurationSecond,
					Description: descRotationPeriod,
				},
			},
		}},
	}
	fields := roleResponse[http.StatusOK][0].Fields
	AddStaticAssumeRoleFieldsEnt(fields)

	return &framework.Path{
		Pattern: fmt.Sprintf("%s/%s", pathStaticRole, framework.GenericNameWithAtRegex(paramRoleName)),
		Fields:  fields,

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:  b.pathStaticRolesRead,
				Responses: roleResponse,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathStaticRolesWrite,
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
				Responses:                   roleResponse,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback:                    b.pathStaticRolesDelete,
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
				Responses: map[int][]framework.Response{
					http.StatusNoContent: {{
						Description: http.StatusText(http.StatusNoContent),
					}},
				},
			},
		},

		HelpSynopsis:    pathStaticRolesHelpSyn,
		HelpDescription: pathStaticRolesHelpDesc,
	}
}

func (b *backend) pathStaticRolesRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName, ok := data.GetOk(paramRoleName)
	if !ok {
		return nil, fmt.Errorf("missing %q parameter", paramRoleName)
	}

	b.roleMutex.RLock()
	defer b.roleMutex.RUnlock()

	entry, err := req.Storage.Get(ctx, formatRoleStoragePath(roleName.(string)))
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration for static role %q: %w", roleName, err)
	}
	if entry == nil {
		return nil, nil
	}

	var config staticRoleEntry
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("failed to decode configuration for static role %q: %w", roleName, err)
	}

	return &logical.Response{
		Data: formatResponse(config),
	}, nil
}

func (b *backend) pathStaticRolesWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Create & validate config from request parameters
	config := staticRoleEntry{}
	isCreate := req.Operation == logical.CreateOperation

	if rawRoleName, ok := data.GetOk(paramRoleName); ok {
		config.Name = rawRoleName.(string)

		if err := b.validateRoleName(config.Name); err != nil {
			return nil, err
		}
	} else {
		return logical.ErrorResponse("missing %q parameter", paramRoleName), nil
	}

	// retrieve old role value
	entry, err := req.Storage.Get(ctx, formatRoleStoragePath(config.Name))
	if err != nil {
		return nil, fmt.Errorf("couldn't check storage for pre-existing role: %w", err)
	}

	if entry != nil {
		err = entry.DecodeJSON(&config)
		if err != nil {
			return nil, fmt.Errorf("couldn't convert existing role into config struct: %w", err)
		}
	} else {
		// if we couldn't find an entry, this is a create event
		isCreate = true
	}

	// other params are optional if we're not Creating

	err = validateAssumeRoleFields(data, &config)
	if err != nil {
		return nil, err
	}

	if rawUsername, ok := data.GetOk(paramUsername); ok {
		config.Username = rawUsername.(string)

		if err := b.validateIAMUserExists(ctx, req.Storage, &config, isCreate); err != nil {
			return nil, err
		}
	} else if isCreate {
		return logical.ErrorResponse("missing %q parameter", paramUsername), nil
	}

	if rawRotationPeriod, ok := data.GetOk(paramRotationPeriod); ok {
		config.RotationPeriod = time.Duration(rawRotationPeriod.(int)) * time.Second

		if err := b.validateRotationPeriod(config.RotationPeriod); err != nil {
			return nil, err
		}
	} else if isCreate {
		return logical.ErrorResponse("missing %q parameter", paramRotationPeriod), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	// Upsert role config
	newRole, err := logical.StorageEntryJSON(formatRoleStoragePath(config.Name), config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object to JSON: %w", err)
	}
	err = req.Storage.Put(ctx, newRole)
	if err != nil {
		return nil, fmt.Errorf("failed to save object in storage: %w", err)
	}

	// Bootstrap initial set of keys if they did not exist before. AWS Secret Access Keys can only be obtained on creation,
	// so we need to boostrap new roles with a new initial set of keys to be able to serve valid credentials to Vault clients.
	credsPath := formatCredsStoragePath(config.Name)
	existingCredsEntry, err := req.Storage.Get(ctx, credsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to verify if credentials already exist for role %q: %w", config.Name, err)
	}
	if existingCredsEntry == nil {
		creds, err := b.createCredential(ctx, req.Storage, config, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create new credentials for role %q: %w", config.Name, err)
		}

		err = b.credRotationQueue.Push(&queue.Item{
			Key:      config.Name,
			Value:    config,
			Priority: creds.priority(config),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add item into the rotation queue for role %q: %w", config.Name, err)
		}
	} else {
		var existingCreds awsCredentials
		err := existingCredsEntry.DecodeJSON(&existingCreds)
		if err != nil {
			return nil, fmt.Errorf("unable to decode existing credentials for role %s: %w", config.Name, err)
		}
		// creds already exist, so all we need to do is update the rotation
		// what here stays the same and what changes? Can we change the name?
		i, err := b.credRotationQueue.PopByKey(config.Name)
		if err != nil {
			return nil, fmt.Errorf("expected an item with name %q, but got an error: %w", config.Name, err)
		}
		// check if i is nil to prevent panic because
		// 1. PopByKey returns nil if the key does not exist; and
		// 2. the static cred queue is not repopulated on reload (see VAULT-30877)
		if i == nil {
			return nil, fmt.Errorf("expected an item with name %q, but got nil", config.Name)
		}
		i.Value = config
		// update the next rotation to occur at now + the new rotation period
		newExpiration := time.Now().Add(config.RotationPeriod)
		existingCreds.Expiration = &newExpiration
		_, err = logical.StorageEntryJSON(credsPath, &existingCreds)
		if err != nil {
			return nil, fmt.Errorf("error updating credentials for role %s: %w", config.Name, err)
		}
		i.Priority = existingCreds.priority(config)

		err = b.credRotationQueue.Push(i)
		if err != nil {
			return nil, fmt.Errorf("failed to add updated item into the rotation queue for role %q: %w", config.Name, err)
		}
	}

	return &logical.Response{
		Data: formatResponse(config),
	}, nil
}

func (b *backend) pathStaticRolesDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName, ok := data.GetOk(paramRoleName)
	if !ok {
		return nil, fmt.Errorf("missing %q parameter", paramRoleName)
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	entry, err := req.Storage.Get(ctx, formatRoleStoragePath(roleName.(string)))
	if err != nil {
		return nil, fmt.Errorf("couldn't locate role in storage due to error: %w", err)
	}
	// no entry in storage, but no error either, congrats, it's deleted!
	if entry == nil {
		return nil, nil
	}
	var cfg staticRoleEntry
	err = entry.DecodeJSON(&cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't convert storage entry to role config")
	}

	err = b.deleteCredential(ctx, req.Storage, cfg, false)
	if err != nil {
		return nil, fmt.Errorf("failed to clean credentials while deleting role %q: %w", roleName.(string), err)
	}

	// delete from the queue
	_, err = b.credRotationQueue.PopByKey(cfg.Name)
	if err != nil {
		return nil, fmt.Errorf("couldn't delete key from queue: %w", err)
	}

	return nil, req.Storage.Delete(ctx, formatRoleStoragePath(roleName.(string)))
}

func (b *backend) validateRoleName(name string) error {
	if name == "" {
		return errors.New("empty role name attribute given")
	}
	return nil
}

// validateIAMUser checks the user information we have for the role against the information on AWS. On a create, it uses the username
// to retrieve the user information and _sets_ the userID. On update, it validates the userID and username.
func (b *backend) validateIAMUserExists(ctx context.Context, storage logical.Storage, entry *staticRoleEntry, isCreate bool) error {
	c, err := b.getNonCachedIAMClient(ctx, storage, *entry)
	if err != nil {
		return fmt.Errorf("unable to get client to validate username %q: %w", entry.Username, err)
	}
	b.iamClient = c

	// we don't really care about the content of the result, just that it's not an error
	out, err := c.GetUser(&iam.GetUserInput{
		UserName: aws.String(entry.Username),
	})
	if err != nil || out.User == nil {
		return fmt.Errorf("unable to validate username %q: %w", entry.Username, err)
	}
	if *out.User.UserName != entry.Username {
		return fmt.Errorf("AWS GetUser returned a username, but it didn't match: %q was requested, but %q was returned", entry.Username, *out.User.UserName)
	}

	if !isCreate && *out.User.UserId != entry.ID {
		return fmt.Errorf("AWS GetUser returned a user, but the ID did not match: %q was requested, but %q was returned", entry.ID, *out.User.UserId)
	} else {
		// if this is an insert, store the userID. This is the immutable part of an IAM user, but it's not exactly user-friendly.
		// So, we allow users to specify usernames, but on updates we'll use the ID as a verification cross-check.
		entry.ID = *out.User.UserId
	}

	return nil
}

const (
	minAllowableRotationPeriod = 1 * time.Minute
)

func (b *backend) validateRotationPeriod(period time.Duration) error {
	if period < b.minAllowableRotationPeriod {
		return fmt.Errorf("role rotation period out of range: must be greater than %.2f seconds", b.minAllowableRotationPeriod.Seconds())
	}
	return nil
}

func formatResponse(cfg staticRoleEntry) map[string]interface{} {
	response := structs.New(cfg).Map()
	response[paramRotationPeriod] = int64(cfg.RotationPeriod.Seconds())

	return response
}

func formatRoleStoragePath(roleName string) string {
	return fmt.Sprintf("%s/%s", pathStaticRole, roleName)
}

const pathStaticRolesHelpSyn = `
Manage static roles for AWS.
`

const pathStaticRolesHelpDesc = `
This path lets you manage static roles (users) for the AWS secret backend.
A static role is associated with a single IAM user, and manages the access
keys based on a rotation period, automatically rotating the credential. If
the IAM user has multiple access keys, the oldest key will be rotated.
`

const (
	descRoleName       = "The name of this role."
	descUsername       = "The IAM user to adopt as a static role."
	descRotationPeriod = `Period by which to rotate the backing credential of the adopted user. 
This can be a Go duration (e.g, '1m', 24h'), or an integer number of seconds.`
	descAssumeRoleARN   = `The AWS ARN for the role to be assumed when interacting with the account specified.`
	descRoleSessionName = `An identifier for the assumed role session.`
	descExternalID      = `An external ID to be passed to the assumed role session.`
)

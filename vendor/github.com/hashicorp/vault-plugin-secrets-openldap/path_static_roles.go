// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package openldap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	staticRolePath = "static-role/"
)

// genericNameWithForwardSlashRegex is a regex which requires a role name. The
// role name can include any number of alphanumeric characters separated by
// forward slashes.
func genericNameWithForwardSlashRegex(name string) string {
	return fmt.Sprintf(`(/(?P<%s>\w(([\w-./]+)?\w)?))`, name)
}

// optionalGenericNameWithForwardSlashListRegex is a regex for optionally
// including a role path in list options. The role path can be used to list
// nested roles at arbitrary depth.
func optionalGenericNameWithForwardSlashListRegex(name string) string {
	return fmt.Sprintf("/?(?P<%s>.+)?", name)
}

func (b *backend) pathListStaticRoles() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: strings.TrimSuffix(staticRolePath, "/") + optionalGenericNameWithForwardSlashListRegex("path"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAP,
				OperationVerb:   "list",
				OperationSuffix: "static-roles",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathStaticRoleList,
				},
			},
			Fields: map[string]*framework.FieldSchema{
				"path": {
					Type:        framework.TypeLowerCaseString,
					Description: "Path of roles to list",
				},
			},
			HelpSynopsis:    staticRolesListHelpSynopsis,
			HelpDescription: staticRolesListHelpDescription,
		},
	}
}

func (b *backend) pathStaticRoles() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: strings.TrimSuffix(staticRolePath, "/") + genericNameWithForwardSlashRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixLDAP,
				OperationSuffix: "static-role",
			},
			Fields:         fieldsForType(staticRolePath),
			ExistenceCheck: b.pathStaticRoleExistenceCheck,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    b.pathStaticRoleCreateUpdate,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback:                    b.pathStaticRoleCreateUpdate,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathStaticRoleRead,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback:                    b.pathStaticRoleDelete,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},
			HelpSynopsis:    staticRoleHelpSynopsis,
			HelpDescription: staticRoleHelpDescription,
		},
	}
}

// fieldsForType returns a map of string/FieldSchema items for the given role
// type. The purpose is to keep the shared fields between dynamic and static
// roles consistent, and allow for each type to override or provide their own
// specific fields
func fieldsForType(roleType string) map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"name": {
			Type:        framework.TypeLowerCaseString,
			Description: "Name of the role",
		},
		"username": {
			Type:        framework.TypeString,
			Description: "The username/logon name for the entry with which this role will be associated.",
		},
		"dn": {
			Type:        framework.TypeString,
			Description: "The distinguished name of the entry to manage.",
		},
	}

	// Get the fields that are specific to the type of role, and add them to the
	// common fields. In the future we can add additional for dynamic roles.
	var typeFields map[string]*framework.FieldSchema
	switch roleType {
	case staticRolePath:
		typeFields = staticFields()
	}

	for k, v := range typeFields {
		fields[k] = v
	}

	return fields
}

// staticFields returns a map of key and field schema items that are specific
// only to static roles
func staticFields() map[string]*framework.FieldSchema {
	fields := map[string]*framework.FieldSchema{
		"rotation_period": {
			Type:        framework.TypeDurationSecond,
			Description: "Period for automatic credential rotation of the given entry.",
		},
		"skip_import_rotation": {
			Type:        framework.TypeBool,
			Description: "Skip the initial pasword rotation on import (has no effect on updates)",
		},
	}
	return fields
}

func (b *backend) pathStaticRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	role, err := b.staticRole(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}
	return role != nil, nil
}

func (b *backend) pathStaticRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Grab the exclusive lock
	lock := locksutil.LockForKey(b.roleLocks, name)
	lock.Lock()
	defer lock.Unlock()

	// TODO: Add retry logic

	// Remove the item from the queue
	_, err := b.popFromRotationQueueByKey(name)
	if err != nil {
		return nil, err
	}

	role, err := b.staticRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	err = req.Storage.Delete(ctx, staticRolePath+name)
	if err != nil {
		return nil, err
	}

	b.managedUserLock.Lock()
	defer b.managedUserLock.Unlock()
	delete(b.managedUsers, role.StaticAccount.Username)

	walIDs, err := framework.ListWAL(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	var merr *multierror.Error
	for _, walID := range walIDs {
		wal, err := b.findStaticWAL(ctx, req.Storage, walID)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}
		if wal != nil && name == wal.RoleName {
			b.Logger().Debug("deleting WAL for deleted role", "WAL ID", walID, "role", name)
			err = framework.DeleteWAL(ctx, req.Storage, walID)
			if err != nil {
				b.Logger().Debug("failed to delete WAL for deleted role", "WAL ID", walID, "error", err)
				merr = multierror.Append(merr, err)
			}
		}
	}

	return nil, merr.ErrorOrNil()
}

func (b *backend) pathStaticRoleRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role, err := b.staticRole(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	data := map[string]interface{}{
		"dn":       role.StaticAccount.DN,
		"username": role.StaticAccount.Username,
	}

	data["rotation_period"] = role.StaticAccount.RotationPeriod.Seconds()
	if !role.StaticAccount.LastVaultRotation.IsZero() {
		data["last_vault_rotation"] = role.StaticAccount.LastVaultRotation
	}

	return &logical.Response{Data: data}, nil
}

func (b *backend) pathStaticRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)

	// Grab the exclusive lock as well potentially pop and re-push the queue item
	// for this role
	lock := locksutil.LockForKey(b.roleLocks, name)
	lock.Lock()
	defer lock.Unlock()

	role, err := b.staticRole(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}

	if role == nil {
		role = &roleEntry{
			StaticAccount: &staticAccount{},
		}
	}

	isCreate := req.Operation == logical.CreateOperation

	b.managedUserLock.Lock()
	defer b.managedUserLock.Unlock()

	usernameRaw, ok := data.GetOk("username")
	if !ok && isCreate {
		return logical.ErrorResponse("username is a required field to manage a static account"), nil
	}
	if ok {
		username := usernameRaw.(string)
		if username == "" {
			return logical.ErrorResponse("username must not be empty"), nil
		}
		if _, exists := b.managedUsers[username]; exists && isCreate {
			return logical.ErrorResponse("%q is already managed by the secrets engine", username), nil
		}
		if !isCreate && username != role.StaticAccount.Username {
			return logical.ErrorResponse("cannot update static account username"), nil
		}

		role.StaticAccount.Username = username
	}

	// DN is optional. Unless it is unset via providing the empty string, it
	// cannot be modified after creation. If given, it will take precedence
	// over username for LDAP search during password rotation.
	if dnRaw, ok := data.GetOk("dn"); ok {
		dn := dnRaw.(string)
		if !isCreate && dn != "" && dn != role.StaticAccount.DN {
			return logical.ErrorResponse("cannot update static account distinguished name (dn)"), nil
		}

		role.StaticAccount.DN = dn
	}

	rotationPeriodSecondsRaw, ok := data.GetOk("rotation_period")
	if !ok && isCreate {
		return logical.ErrorResponse("rotation_period is required to create static accounts"), nil
	}
	if ok {
		rotationPeriodSeconds := rotationPeriodSecondsRaw.(int)
		if rotationPeriodSeconds < queueTickSeconds {
			// If rotation frequency is specified the value must be at least
			// that of the constant queueTickSeconds (5 seconds at time of writing),
			// otherwise we won't be able to rotate in time
			return logical.ErrorResponse("rotation_period must be %d seconds or more", queueTickSeconds), nil
		}
		role.StaticAccount.RotationPeriod = time.Duration(rotationPeriodSeconds) * time.Second
	}

	skipRotation := false
	skipRotationRaw, ok := data.GetOk("skip_import_rotation")
	if ok {
		// if skip rotation was set, use it (or validation error on an update)
		if !isCreate {
			return logical.ErrorResponse("skip_import_rotation has no effect on updates"), nil
		}
		skipRotation = skipRotationRaw.(bool)
	} else if isCreate {
		// otherwise, go get it if this is a create request.
		c, err := readConfig(ctx, req.Storage)
		if err != nil {
			return nil, err
		}
		if c == nil {
			return logical.ErrorResponse("missing LDAP configuration"), nil
		}

		skipRotation = c.SkipStaticRoleImportRotation
	}

	// lvr represents the role's LastVaultRotation
	lvr := role.StaticAccount.LastVaultRotation

	// Only call setStaticAccountPassword if we're creating the role for the first time
	var item *queue.Item
	switch req.Operation {
	case logical.CreateOperation:
		// if we were asked to not rotate, just add the entry - this essentially becomes an update operation, except
		// the item is new
		if skipRotation {
			entry, err := logical.StorageEntryJSON(staticRolePath+name, role)
			if err != nil {
				return nil, err
			}
			if err := req.Storage.Put(ctx, entry); err != nil {
				return nil, err
			}

			// set the item
			item = &queue.Item{
				Key: name,
			}
			// synthetically set lvr to now, so that it gets queued correctly
			lvr = time.Now()
			break
		} else {
			// setStaticAccountPassword calls Storage.Put and saves the role to storage
			resp, err := b.setStaticAccountPassword(ctx, req.Storage, &setStaticAccountInput{
				RoleName: name,
				Role:     role,
			})
			if err != nil {
				if resp != nil && resp.WALID != "" {
					b.Logger().Debug("deleting WAL for failed role creation", "WAL ID", resp.WALID, "role", name)
					walDeleteErr := framework.DeleteWAL(ctx, req.Storage, resp.WALID)
					if walDeleteErr != nil {
						b.Logger().Debug("failed to delete WAL for failed role creation", "WAL ID", resp.WALID, "error", walDeleteErr)
						var merr *multierror.Error
						merr = multierror.Append(merr, err)
						merr = multierror.Append(merr, fmt.Errorf("failed to clean up WAL from failed role creation: %w", walDeleteErr))
						err = merr.ErrorOrNil()
					}
				}
				return nil, err
			}
			// guard against RotationTime not being set or zero-value
			lvr = resp.RotationTime
			item = &queue.Item{
				Key: name,
			}
		}
	case logical.UpdateOperation:
		// store updated Role
		entry, err := logical.StorageEntryJSON(staticRolePath+name, role)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		// In case this is an update, remove any previous version of the item from
		// the queue. The existing item could be tracking a WAL ID for this role,
		// so it's important to keep the existing item rather than recreate it.
		// TODO: Add retry logic
		item, err = b.popFromRotationQueueByKey(name)
		if err != nil {
			return nil, err
		}
	}

	item.Priority = lvr.Add(role.StaticAccount.RotationPeriod).Unix()

	// Add their rotation to the queue
	if err := b.pushItem(item); err != nil {
		return nil, err
	}

	b.managedUsers[role.StaticAccount.Username] = struct{}{}

	return nil, nil
}

type roleEntry struct {
	StaticAccount *staticAccount `json:"static_account" mapstructure:"static_account"`
}

type staticAccount struct {
	// DN to create or assume management for static accounts
	DN string `json:"dn"`

	// Username to create or assume management for static accounts
	Username string `json:"username"`

	// Password is the current password for static accounts. As an input, this is
	// used/required when trying to assume management of an existing static
	// account. This is returned on credential requests if it exists.
	Password string `json:"password"`

	// LastPassword is the prior password after a rotation for static accounts.
	// This is returned on credential requests if it exists.
	LastPassword string `json:"last_password"`

	// LastVaultRotation represents the last time Vault rotated the password
	LastVaultRotation time.Time `json:"last_vault_rotation"`

	// RotationPeriod is number in seconds between each rotation, effectively a
	// "time to live". This value is compared to the LastVaultRotation to
	// determine if a password needs to be rotated
	RotationPeriod time.Duration `json:"rotation_period"`
}

// NextRotationTime calculates the next rotation by adding the Rotation Period
// to the last known vault rotation
func (s *staticAccount) NextRotationTime() time.Time {
	return s.LastVaultRotation.Add(s.RotationPeriod)
}

// PasswordTTL calculates the approximate time remaining until the password is
// no longer valid. This is approximate because the periodic rotation is only
// checked approximately every 5 seconds, and each rotation can take a small
// amount of time to process. This can result in a negative TTL time while the
// rotation function processes the Static Role and performs the rotation. If the
// TTL is negative, zero is returned. Users should not trust passwords with a
// Zero TTL, as they are likely in the process of being rotated and will quickly
// be invalidated.
func (s *staticAccount) PasswordTTL() time.Duration {
	next := s.NextRotationTime()
	ttl := next.Sub(time.Now()).Round(time.Second)
	if ttl < 0 {
		ttl = time.Duration(0)
	}
	return ttl
}

func (b *backend) pathStaticRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	rolePath := data.Get("path").(string)
	roles, err := req.Storage.List(ctx, staticRolePath+rolePath)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	return logical.ListResponse(roles), nil
}

func (b *backend) staticRole(ctx context.Context, s logical.Storage, roleName string) (*roleEntry, error) {
	completeRole := staticRolePath + roleName
	entry, err := s.Get(ctx, completeRole)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

const staticRoleHelpSynopsis = `
Manage the static roles that can be created with this backend.
`

const staticRoleHelpDescription = `
This path lets you manage the static roles that can be created with this
backend. Static Roles are associated with a single LDAP entry, and manage the
password based on a rotation period, automatically rotating the password.

The "username" parameter is required and configures the username for the LDAP entry. 
This is helpful to provide a usable name when distinguished name (DN) isn't used 
directly for authentication. If DN not provided, "username" will be used for LDAP 
subtree search, rooted at the "userdn" configuration value. The name attribute to use 
when searching for the user can be configured with the "userattr" configuration value.

The "dn" parameter is optional and configures the distinguished name to use 
when managing the existing entry. If the "dn" parameter is set, it will take 
precedence over the "username" when LDAP searches are performed.

The "rotation_period' parameter is required and configures how often, in seconds, 
the credentials should be automatically rotated by Vault.  The minimum is 5 seconds (5s).
`

const staticRolesListHelpDescription = `
List all the static roles being managed by Vault.
`

const staticRolesListHelpSynopsis = `
This path lists all the static roles Vault is currently managing within the LDAP system.
`

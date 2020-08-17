package openldap

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/locksutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/queue"
)

const (
	staticRolePath = "static-role/"
)

func (b *backend) pathListRoles() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: staticRolePath + "?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathRoleList,
				},
			},
			HelpSynopsis:    staticRolesListHelpSynopsis,
			HelpDescription: staticRolesListHelpDescription,
		},
	}
}

func (b *backend) pathRoles() []*framework.Path {
	return []*framework.Path{
		{
			Pattern:        staticRolePath + framework.GenericNameRegex("name"),
			Fields:         fieldsForType(staticRolePath),
			ExistenceCheck: b.pathStaticRoleExistenceCheck,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathStaticRoleCreateUpdate,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathStaticRoleCreateUpdate,
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
		"ttl": {
			Type:        framework.TypeDurationSecond,
			Description: "The time-to-live for the password.",
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
	}
	return fields
}

func (b *backend) pathStaticRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	role, err := b.StaticRole(ctx, req.Storage, data.Get("name").(string))
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

	//TODO: Add retry logic

	// Remove the item from the queue
	_, err := b.popFromRotationQueueByKey(name)
	if err != nil {
		return nil, err
	}

	role, err := b.StaticRole(ctx, req.Storage, name)
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

	return nil, nil
}

func (b *backend) pathStaticRoleRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role, err := b.StaticRole(ctx, req.Storage, d.Get("name").(string))
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

	role, err := b.StaticRole(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}

	if role == nil {
		role = &roleEntry{
			StaticAccount: &staticAccount{},
		}
	}

	dn := data.Get("dn").(string)
	if dn == "" {
		return logical.ErrorResponse("dn is a required field to manage a static account"), nil
	}
	role.StaticAccount.DN = dn

	username := data.Get("username").(string)
	if username == "" {
		return logical.ErrorResponse("username is a required field to manage a static account"), nil
	}
	role.StaticAccount.Username = username

	rotationPeriodSecondsRaw, ok := data.GetOk("rotation_period")
	if !ok {
		return logical.ErrorResponse("rotation_period is required for static accounts"), nil
	}
	rotationPeriodSeconds := rotationPeriodSecondsRaw.(int)
	if rotationPeriodSeconds < queueTickSeconds {
		// If rotation frequency is specified the value
		// must be at least that of the constant queueTickSeconds (5 seconds at
		// time of writing), otherwise we wont be able to rotate in time
		return logical.ErrorResponse("rotation_period must be %d seconds or more", queueTickSeconds), nil
	}
	role.StaticAccount.RotationPeriod = time.Duration(rotationPeriodSeconds) * time.Second

	// lvr represents the role's LastVaultRotation
	lvr := role.StaticAccount.LastVaultRotation

	// Only call setStaticAccountPassword if we're creating the role for the
	// first time
	switch req.Operation {
	case logical.CreateOperation:
		// setStaticAccountPassword calls Storage.Put and saves the role to storage
		resp, err := b.setStaticAccountPassword(ctx, req.Storage, &setStaticAccountInput{
			RoleName: name,
			Role:     role,
		})
		if err != nil {
			return nil, err
		}
		// guard against RotationTime not being set or zero-value
		lvr = resp.RotationTime
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
		// the queue

		//TODO: Add retry logic
		_, err = b.popFromRotationQueueByKey(name)
		if err != nil {
			return nil, err
		}
	}

	// Add their rotation to the queue
	if err := b.pushItem(&queue.Item{
		Key:      name,
		Priority: lvr.Add(role.StaticAccount.RotationPeriod).Unix(),
	}); err != nil {
		return nil, err
	}

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
	// account. Return this on credential request if it exists.
	Password string `json:"password"`

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

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	path := staticRolePath
	entries, err := req.Storage.List(ctx, path)
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) StaticRole(ctx context.Context, s logical.Storage, roleName string) (*roleEntry, error) {
	return b.roleAtPath(ctx, s, roleName, staticRolePath)
}

func (b *backend) roleAtPath(ctx context.Context, s logical.Storage, roleName string, pathPrefix string) (*roleEntry, error) {
	entry, err := s.Get(ctx, pathPrefix+roleName)
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

The "dn" parameter is required and configures the domain name to use when managing 
the existing entry.

The "username" parameter is required and configures the username for the LDAP entry. 
This is helpful to provide a usable name when domain name (DN) isn't used directly for 
authentication.


The "rotation_period' parameter is required and configures how often, in seconds, the credentials should be 
automatically rotated by Vault.  The minimum is 5 seconds (5s).
`

const staticRolesListHelpDescription = `
List all the static roles being managed by Vault.
`

const staticRolesListHelpSynopsis = `
This path lists all the static roles Vault is currently managing in OpenLDAP.
`

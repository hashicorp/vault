package database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/queue"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func pathRoles(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},

			"db_name": {
				Type:        framework.TypeString,
				Description: "Name of the database this role acts on.",
			},
			"creation_statements": {
				Type: framework.TypeStringSlice,
				Description: `Specifies the database statements executed to
				create and configure a user. See the plugin's API page for more
				information on support and formatting for this parameter.`,
			},
			"revocation_statements": {
				Type: framework.TypeStringSlice,
				Description: `Specifies the database statements to be executed
				to revoke a user. See the plugin's API page for more information
				on support and formatting for this parameter.`,
			},
			"renew_statements": {
				Type: framework.TypeStringSlice,
				Description: `Specifies the database statements to be executed
				to renew a user. Not every plugin type will support this
				functionality. See the plugin's API page for more information on
				support and formatting for this parameter. `,
			},
			"rollback_statements": {
				Type: framework.TypeStringSlice,
				Description: `Specifies the database statements to be executed
				rollback a create operation in the event of an error. Not every
				plugin type will support this functionality. See the plugin's
				API page for more information on support and formatting for this
				parameter.`,
			},

			"default_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Default ttl for role.",
			},

			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Maximum time a credential is valid for",
			},
			// TODO: consider renaming to "static_username" for clarity
			"username": {
				Type: framework.TypeString,
				Description: `Name of the static user account for Vault to manage.
				Requires "rotation_period" to be specified`,
			},
			"rotation_period": {
				Type: framework.TypeDurationSecond,
				Description: `Period for automatic credential rotation of the given
				username. Not valid unless used with "username".`,
			},
			"rotation_statements": {
				Type: framework.TypeStringSlice,
				Description: `Specifies the database statements to be executed to rotate
				the accounts credentials. Not every plugin type will support this
				functionality. See the plugin's API page for more information on support
				and formatting for this parameter.`,
			},
		},

		ExistenceCheck: b.pathRoleExistenceCheck,
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead,
			logical.CreateOperation: b.pathRoleCreateUpdate,
			logical.UpdateOperation: b.pathRoleCreateUpdate,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func (b *databaseBackend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	role, err := b.Role(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}
	return role != nil, nil
}

func (b *databaseBackend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	// if this role is a static account, we need to revoke the user from the
	// database
	// TODO: wrap this in a WAL
	role, err := b.Role(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// clean up the static useraccount, if it exists
	if role.StaticAccount != nil {
		db, err := b.GetConnection(ctx, req.Storage, role.DBName)
		if err != nil {
			return nil, err
		}

		db.RLock()
		defer db.RUnlock()

		if err := db.RevokeUser(ctx, role.Statements, role.StaticAccount.Username); err != nil {
			b.CloseIfShutdown(db, err)
			return nil, err
		}
	}

	err = req.Storage.Delete(ctx, "role/"+name)
	if err != nil {
		return nil, err
	}

	if b.credRotationQueue != nil {
		if _, err := b.credRotationQueue.PopItemByKey(name); err != nil {
			if _, ok := err.(*queue.ErrItemNotFound); !ok {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (b *databaseBackend) pathRoleRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role, err := b.Role(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	data := map[string]interface{}{
		"db_name":               role.DBName,
		"creation_statements":   role.Statements.Creation,
		"revocation_statements": role.Statements.Revocation,
		"rollback_statements":   role.Statements.Rollback,
		"renew_statements":      role.Statements.Renewal,
		"rotation_statements":   role.Statements.Rotation,
		"default_ttl":           role.DefaultTTL.Seconds(),
		"max_ttl":               role.MaxTTL.Seconds(),
	}
	if len(role.Statements.Creation) == 0 {
		data["creation_statements"] = []string{}
	}
	if len(role.Statements.Revocation) == 0 {
		data["revocation_statements"] = []string{}
	}
	if len(role.Statements.Rollback) == 0 {
		data["rollback_statements"] = []string{}
	}
	if len(role.Statements.Renewal) == 0 {
		data["renew_statements"] = []string{}
	}
	if len(role.Statements.Rotation) == 0 {
		data["rotation_statements"] = []string{}
	}

	if role.StaticAccount != nil {
		data["username"] = role.StaticAccount.Username
		data["rotation_period"] = role.StaticAccount.RotationPeriod.Seconds()
		if !role.StaticAccount.LastVaultRotation.IsZero() {
			// TODO: formatting
			data["last_vault_rotation"] = role.StaticAccount.LastVaultRotation
		}
	}

	return &logical.Response{
		Data: data,
	}, nil
}

func (b *databaseBackend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *databaseBackend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name := data.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("empty role name attribute given"), nil
	}

	role, err := b.Role(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}

	// createRole is a boolean to indicate if this is a new role creation. This
	// is used to ensure we do not allow an existing role to be "migrated" to
	// role with a static account. If createRole is false and static_account
	// data is given, return an error
	createRole := req.Operation == logical.CreateOperation
	if role == nil {
		role = &roleEntry{}
		createRole = true
	}

	// Static Account information
	if username, ok := data.Get("username").(string); ok && username != "" {
		// If the role exists and there is no StaticAccount, return error
		if role.StaticAccount == nil {
			if !createRole {
				return logical.ErrorResponse("cannot change an existing role to a static account"), nil
			}
			role.StaticAccount = &staticAccount{}
		}

		// If it's a Create operation, both username and rotation_period must be included
		rotationPeriodSecondsRaw, ok := data.GetOk("rotation_period")
		if !ok && createRole {
			return logical.ErrorResponse("rotation_period is required to create static accounts"), nil
		}
		if ok {
			rotationPeriodSeconds := rotationPeriodSecondsRaw.(int)
			if rotationPeriodSeconds < 5 {
				// If rotation frequency is specified, and this is an update, the value
				// must be at least 5 seconds because our periodic func runs about once a
				// minute.
				return logical.ErrorResponse("rotation_period must be 5 seconds or more"), nil
			}
			role.StaticAccount.RotationPeriod = time.Duration(rotationPeriodSeconds) * time.Second
		}

		if role.StaticAccount.Username != "" && role.StaticAccount.Username != username {
			return logical.ErrorResponse("cannot update static account username"), nil
		}
		role.StaticAccount.Username = username

		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid rotation_period: %s", err)), nil
		}

		// TODO: not sure why the check on logical.CreateOperation here
		if rotationStmtsRaw, ok := data.GetOk("rotation_statements"); ok {
			role.Statements.Rotation = rotationStmtsRaw.([]string)
		} else if req.Operation == logical.CreateOperation {
			role.Statements.Rotation = data.Get("rotation_statements").([]string)
		}

		if len(role.Statements.Rotation) == 0 {
			return logical.ErrorResponse("rotation_statements is a required field for static accounts"), nil
		}
	}

	// DB Attributes
	{
		if dbNameRaw, ok := data.GetOk("db_name"); ok {
			role.DBName = dbNameRaw.(string)
		} else if req.Operation == logical.CreateOperation {
			role.DBName = data.Get("db_name").(string)
		}
		if role.DBName == "" {
			return logical.ErrorResponse("empty database name attribute"), nil
		}
	}

	// TTLs
	{
		if defaultTTLRaw, ok := data.GetOk("default_ttl"); ok {
			role.DefaultTTL = time.Duration(defaultTTLRaw.(int)) * time.Second
		} else if req.Operation == logical.CreateOperation {
			role.DefaultTTL = time.Duration(data.Get("default_ttl").(int)) * time.Second
		}
		if maxTTLRaw, ok := data.GetOk("max_ttl"); ok {
			role.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
		} else if req.Operation == logical.CreateOperation {
			role.MaxTTL = time.Duration(data.Get("max_ttl").(int)) * time.Second
		}
	}

	// Statements
	{
		if creationStmtsRaw, ok := data.GetOk("creation_statements"); ok {
			role.Statements.Creation = creationStmtsRaw.([]string)
		} else if req.Operation == logical.CreateOperation {
			role.Statements.Creation = data.Get("creation_statements").([]string)
		}

		if revocationStmtsRaw, ok := data.GetOk("revocation_statements"); ok {
			role.Statements.Revocation = revocationStmtsRaw.([]string)
		} else if req.Operation == logical.CreateOperation {
			role.Statements.Revocation = data.Get("revocation_statements").([]string)
		}

		if rollbackStmtsRaw, ok := data.GetOk("rollback_statements"); ok {
			role.Statements.Rollback = rollbackStmtsRaw.([]string)
		} else if req.Operation == logical.CreateOperation {
			role.Statements.Rollback = data.Get("rollback_statements").([]string)
		}

		if renewStmtsRaw, ok := data.GetOk("renew_statements"); ok {
			role.Statements.Renewal = renewStmtsRaw.([]string)
		} else if req.Operation == logical.CreateOperation {
			role.Statements.Renewal = data.Get("renew_statements").([]string)
		}

		// Do not persist deprecated statements that are populated on role read
		role.Statements.CreationStatements = ""
		role.Statements.RevocationStatements = ""
		role.Statements.RenewStatements = ""
		role.Statements.RollbackStatements = ""
	}

	role.Statements.Revocation = strutil.RemoveEmpty(role.Statements.Revocation)

	if role.StaticAccount != nil {
		// in create/update of static accounts, we only care if the operation
		// err'd , and this call does not return credentials

		// lvr represents the roles' LastVaultRotation
		lvr := role.StaticAccount.LastVaultRotation

		// only call createUpdateStaticAccount if we're creating the role for the
		// first time
		switch req.Operation {
		case logical.CreateOperation:
			resp, err := b.createUpdateStaticAccount(ctx, req.Storage, &setPasswordInput{
				RoleName:   name,
				Role:       role,
				CreateUser: createRole,
			})
			if err != nil {
				return nil, err
			}
			// guard against RotationTime not being set or zero-value
			lvr = resp.RotationTime
			if lvr.IsZero() {
				lvr = time.Now()
			}
		case logical.UpdateOperation:
			// In case this is an update, remove any previous version of the item from the queue
			if _, err := b.credRotationQueue.PopItemByKey(name); err != nil {
				if _, ok := err.(*queue.ErrItemNotFound); !ok {
					return nil, err
				}
			}
		}

		// Add their rotation to the queue
		if err := b.credRotationQueue.PushItem(&queue.Item{
			Key:      name,
			Priority: lvr.Add(role.StaticAccount.RotationPeriod).Unix(),
		}); err != nil {
			return nil, err
		}
	}
	// END create/update static account

	// Store it
	entry, err := logical.StorageEntryJSON("role/"+name, role)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type roleEntry struct {
	DBName        string              `json:"db_name"`
	Statements    dbplugin.Statements `json:"statements"`
	DefaultTTL    time.Duration       `json:"default_ttl"`
	MaxTTL        time.Duration       `json:"max_ttl"`
	StaticAccount *staticAccount      `json:"static_account" mapstructure:"static_account"`
}

type staticAccount struct {
	// Fields used in Static Accounts and automatic password rotation

	// Username to create or assume management for static accounts
	Username string `json:"username"`

	// Password is the current password for static accounts. As an input, this is
	// used/required when trying to assume management of an existing static
	// account. Return this on credential request if it exists.
	Password string `json:"password"`

	// LastVaultRotation represents the last time Vault rotated the password
	// PasswordLastSet represents the last time a manual password rotation was
	// preformed, using the Vault endpoint
	LastVaultRotation time.Time `json:"last_vault_rotation"`
	PasswordLastSet   time.Time `json:"password_last_set"`

	// RotationPeriod is number in seconds between each rotation, effectively a
	// "time to live". This value is compared to the LastVaultRotation to
	// determine if a password needs to be rotated
	RotationPeriod time.Duration `json:"rotation_period"`

	// previousPassword is used to preserve the previous password during a
	// rotation. If any step in the process fails, we have record of the previous
	// password and can attempt to roll back.
	previousPassword string
}

const pathRoleHelpSyn = `
Manage the roles that can be created with this backend.
`

const pathRoleHelpDesc = `
This path lets you manage the roles that can be created with this backend.

The "db_name" parameter is required and configures the name of the database
connection to use.

The "creation_statements" parameter customizes the string used to create the
credentials. This can be a sequence of SQL queries, or other statement formats
for a particular database type. Some substitution will be done to the statement
strings for certain keys. The names of the variables must be surrounded by "{{"
and "}}" to be replaced.

  * "name" - The random username generated for the DB user.

  * "password" - The random password generated for the DB user.

  * "expiration" - The timestamp when this user will expire.

Example of a decent creation_statements for a postgresql database plugin:

	CREATE ROLE "{{name}}" WITH
	  LOGIN
	  PASSWORD '{{password}}'
	  VALID UNTIL '{{expiration}}';
	GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO "{{name}}";

The "revocation_statements" parameter customizes the statement string used to
revoke a user. Example of a decent revocation_statements for a postgresql
database plugin:

	REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM {{name}};
	REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM {{name}};
	REVOKE USAGE ON SCHEMA public FROM {{name}};
	DROP ROLE IF EXISTS {{name}};

The "renew_statements" parameter customizes the statement string used to renew a
user.
The "rollback_statements' parameter customizes the statement string used to
rollback a change if needed.
`

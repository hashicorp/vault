package database

import (
	"context"
	"time"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList(),
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
		},

		ExistenceCheck: b.pathRoleExistenceCheck(),
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead(),
			logical.CreateOperation: b.pathRoleCreateUpdate(),
			logical.UpdateOperation: b.pathRoleCreateUpdate(),
			logical.DeleteOperation: b.pathRoleDelete(),
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func (b *databaseBackend) pathRoleExistenceCheck() framework.ExistenceFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
		role, err := b.Role(ctx, req.Storage, data.Get("name").(string))
		if err != nil {
			return false, err
		}

		return role != nil, nil
	}
}

func (b *databaseBackend) pathRoleDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		err := req.Storage.Delete(ctx, "role/"+data.Get("name").(string))
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (b *databaseBackend) pathRoleRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

		return &logical.Response{
			Data: data,
		}, nil
	}
}

func (b *databaseBackend) pathRoleList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		entries, err := req.Storage.List(ctx, "role/")
		if err != nil {
			return nil, err
		}

		return logical.ListResponse(entries), nil
	}
}

func (b *databaseBackend) pathRoleCreateUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse("empty role name attribute given"), nil
		}

		role, err := b.Role(ctx, req.Storage, data.Get("name").(string))
		if err != nil {
			return nil, err
		}
		if role == nil {
			role = &roleEntry{}
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
}

type roleEntry struct {
	DBName     string              `json:"db_name"`
	Statements dbplugin.Statements `json:"statements"`
	DefaultTTL time.Duration       `json:"default_ttl"`
	MaxTTL     time.Duration       `json:"max_ttl"`
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

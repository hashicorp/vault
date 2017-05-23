package database

import (
	"time"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
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
				Type: framework.TypeString,
				Description: `Specifies the database statements executed to
				create and configure a user. See the plugin's API page for more
				information on support and formatting for this parameter.`,
			},
			"revocation_statements": {
				Type: framework.TypeString,
				Description: `Specifies the database statements to be executed
				to revoke a user. See the plugin's API page for more information
				on support and formatting for this parameter.`,
			},
			"renew_statements": {
				Type: framework.TypeString,
				Description: `Specifies the database statements to be executed
				to renew a user. Not every plugin type will support this
				functionality. See the plugin's API page for more information on
				support and formatting for this parameter. `,
			},
			"rollback_statements": {
				Type: framework.TypeString,
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

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRoleRead(),
			logical.UpdateOperation: b.pathRoleCreate(),
			logical.DeleteOperation: b.pathRoleDelete(),
		},

		HelpSynopsis:    pathRoleHelpSyn,
		HelpDescription: pathRoleHelpDesc,
	}
}

func (b *databaseBackend) pathRoleDelete() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		err := req.Storage.Delete("role/" + data.Get("name").(string))
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func (b *databaseBackend) pathRoleRead() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		role, err := b.Role(req.Storage, data.Get("name").(string))
		if err != nil {
			return nil, err
		}
		if role == nil {
			return nil, nil
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"db_name":               role.DBName,
				"creation_statements":   role.Statements.CreationStatements,
				"revocation_statements": role.Statements.RevocationStatements,
				"rollback_statements":   role.Statements.RollbackStatements,
				"renew_statements":      role.Statements.RenewStatements,
				"default_ttl":           role.DefaultTTL.Seconds(),
				"max_ttl":               role.MaxTTL.Seconds(),
			},
		}, nil
	}
}

func (b *databaseBackend) pathRoleList() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		entries, err := req.Storage.List("role/")
		if err != nil {
			return nil, err
		}

		return logical.ListResponse(entries), nil
	}
}

func (b *databaseBackend) pathRoleCreate() framework.OperationFunc {
	return func(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)
		if name == "" {
			return logical.ErrorResponse("empty role name attribute given"), nil
		}

		dbName := data.Get("db_name").(string)
		if dbName == "" {
			return logical.ErrorResponse("empty database name attribute given"), nil
		}

		// Get statements
		creationStmts := data.Get("creation_statements").(string)
		revocationStmts := data.Get("revocation_statements").(string)
		rollbackStmts := data.Get("rollback_statements").(string)
		renewStmts := data.Get("renew_statements").(string)

		// Get TTLs
		defaultTTLRaw := data.Get("default_ttl").(int)
		maxTTLRaw := data.Get("max_ttl").(int)
		defaultTTL := time.Duration(defaultTTLRaw) * time.Second
		maxTTL := time.Duration(maxTTLRaw) * time.Second

		statements := dbplugin.Statements{
			CreationStatements:   creationStmts,
			RevocationStatements: revocationStmts,
			RollbackStatements:   rollbackStmts,
			RenewStatements:      renewStmts,
		}

		// Store it
		entry, err := logical.StorageEntryJSON("role/"+name, &roleEntry{
			DBName:     dbName,
			Statements: statements,
			DefaultTTL: defaultTTL,
			MaxTTL:     maxTTL,
		})
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(entry); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

type roleEntry struct {
	DBName     string              `json:"db_name" mapstructure:"db_name" structs:"db_name"`
	Statements dbplugin.Statements `json:"statments" mapstructure:"statements" structs:"statments"`
	DefaultTTL time.Duration       `json:"default_ttl" mapstructure:"default_ttl" structs:"default_ttl"`
	MaxTTL     time.Duration       `json:"max_ttl" mapstructure:"max_ttl" structs:"max_ttl"`
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

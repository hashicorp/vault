package database

import (
        "context"
        "fmt"
        "time"

        "github.com/hashicorp/vault/builtin/logical/database/dbplugin"
        "github.com/hashicorp/vault/helper/parseutil"
        "github.com/hashicorp/vault/helper/strutil"
        "github.com/hashicorp/vault/logical"
        "github.com/hashicorp/vault/logical/framework"
        "github.com/y0ssar1an/q"
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

                        "static_account": &framework.FieldSchema{
                                Type:        framework.TypeMap,
                                Description: `Static account thing. only accpets a few things`,
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
                // rn := data.Get("name").(string)
                // q.Q("pathRoleExistenceCheck:=", rn)
                role, err := b.Role(ctx, req.Storage, data.Get("name").(string))
                if err != nil {
                        // q.Q("role does not exist")
                        return false, err
                }

                // q.Q("role exists:", role != nil)
                return role != nil, nil
        }
}

func (b *databaseBackend) pathRoleDelete() framework.OperationFunc {
        return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
                q.Q("name to delete:", data.Get("name"))
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

                if role.StaticAccount != nil {
                        sa := make(map[string]interface{})
                        sa["username"] = role.StaticAccount.Username
                        sa["rotation_frequency"] = role.StaticAccount.RotationFrequency.Nanoseconds()
                        if role.StaticAccount.Password != "" {
                                sa["password"] = role.StaticAccount.Password
                        }
                        data["static_account"] = sa
                        q.Q("static account=", sa)
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
                        // q.Q("did not find role for name:", data.Get("name").(string))
                        role = &roleEntry{}
                } else {
                        // q.Q("found role for name:", data.Get("name").(string))
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

                // Static Account information
                staticRaw := data.Get("static_account").(map[string]interface{})
                var sa *staticAccount
                if len(staticRaw) > 0 {
                        sa = &staticAccount{}

                        if v, ok := staticRaw["username"].(string); ok {
                                sa.Username = v
                        } else {
                                return logical.ErrorResponse("username is a required field for static accounts"), nil
                        }

                        if v, ok := staticRaw["rotation_frequency"]; ok {
                                sa.RotationFrequency, err = parseutil.ParseDurationSecond(v)
                                if err != nil {
                                        return logical.ErrorResponse(fmt.Sprintf("invalid rotation_frequency: %s", err)), nil
                                }
                        } else {
                                return logical.ErrorResponse("rotation_frequency is a required field for static accounts"), nil
                        }

                        if p, ok := staticRaw["password"].(string); ok {
                                sa.Password = p
                        }
                }

                if sa != nil {
                        role.StaticAccount = sa
                }

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

        // RotationFrequency is numer in seconds between each rotation, effectively a
        // "time to live". This value is compared to the LastVaultRotation to
        // determine if a password needs to be rotated
        RotationFrequency time.Duration `json:"rotation_frequency"`

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

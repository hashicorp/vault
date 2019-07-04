package database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCredsCreate(b *databaseBackend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			"ttl": &framework.FieldSchema{
				Type:    framework.TypeDurationSecond,
				Default: 0,
				Description: `Duration in seconds after which the issued credentials should expire.
Defaults to 0, in which the value will fallback to default_ttl on the role`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathCredsCreateRead(),
		},

		HelpSynopsis:    pathCredsCreateReadHelpSyn,
		HelpDescription: pathCredsCreateReadHelpDesc,
	}
}

func (b *databaseBackend) pathCredsCreateRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)

		// Get the role
		role, err := b.Role(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
		}

		dbConfig, err := b.DatabaseConfig(ctx, req.Storage, role.DBName)
		if err != nil {
			return nil, err
		}

		// If role name isn't in the database's allowed roles, send back a
		// permission denied.
		if !strutil.StrListContains(dbConfig.AllowedRoles, "*") && !strutil.StrListContainsGlob(dbConfig.AllowedRoles, name) {
			return nil, fmt.Errorf("%q is not an allowed role", name)
		}

		// Get the Database object
		db, err := b.GetConnection(ctx, req.Storage, role.DBName)
		if err != nil {
			return nil, err
		}

		db.RLock()
		defer db.RUnlock()

		var ttl time.Duration
		if rawTTL, ok := data.GetOk("ttl"); ok {
			ttl = time.Duration(rawTTL.(int)) * time.Second
		}

		calcTTL, warnings, err := framework.CalculateTTL(b.System(), 0, role.DefaultTTL, ttl, role.MaxTTL, 0, time.Time{})
		if err != nil {
			return nil, err
		}
		expiration := time.Now().Add(calcTTL)
		// Adding a small buffer since the TTL will be calculated again after this call
		// to ensure the database credential does not expire before the lease
		expiration = expiration.Add(5 * time.Second)

		usernameConfig := dbplugin.UsernameConfig{
			DisplayName: req.DisplayName,
			RoleName:    name,
		}

		// Create the user
		username, password, err := db.CreateUser(ctx, role.Statements, usernameConfig, expiration)
		if err != nil {
			b.CloseIfShutdown(db, err)
			return nil, err
		}

		resp := b.Secret(SecretCredsType).Response(map[string]interface{}{
			"username": username,
			"password": password,
		}, map[string]interface{}{
			"username":              username,
			"role":                  name,
			"db_name":               role.DBName,
			"revocation_statements": role.Statements.Revocation,
		})
		resp.Secret.TTL = calcTTL
		resp.Secret.MaxTTL = role.MaxTTL
		for _, warning := range warnings {
			resp.AddWarning(warning)
		}
		return resp, nil
	}
}

const pathCredsCreateReadHelpSyn = `
Request database credentials for a certain role.
`

const pathCredsCreateReadHelpDesc = `
This path reads database credentials for a certain role. The
database credentials will be generated on demand and will be automatically
revoked when the lease is up.
`

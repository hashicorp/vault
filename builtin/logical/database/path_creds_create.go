// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCredsCreate(b *databaseBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "creds/" + framework.GenericNameRegex("name"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixDatabase,
				OperationVerb:   "generate",
				OperationSuffix: "credentials",
			},

			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathCredsCreateRead(),
			},

			HelpSynopsis:    pathCredsCreateReadHelpSyn,
			HelpDescription: pathCredsCreateReadHelpDesc,
		},
		{
			Pattern: "static-creds/" + framework.GenericNameRegex("name"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixDatabase,
				OperationVerb:   "read",
				OperationSuffix: "static-role-credentials",
			},

			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the static role.",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathStaticCredsRead(),
			},

			HelpSynopsis:    pathStaticCredsReadHelpSyn,
			HelpDescription: pathStaticCredsReadHelpDesc,
		},
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

		// If the plugin doesn't support the credential type, return an error
		if !dbConfig.SupportsCredentialType(role.CredentialType) {
			return logical.ErrorResponse("unsupported credential_type: %q",
				role.CredentialType.String()), nil
		}

		// Get the Database object
		dbi, err := b.GetConnection(ctx, req.Storage, role.DBName)
		if err != nil {
			return nil, err
		}

		dbi.RLock()
		defer dbi.RUnlock()

		ttl, _, err := framework.CalculateTTL(b.System(), 0, role.DefaultTTL, 0, role.MaxTTL, 0, time.Time{})
		if err != nil {
			return nil, err
		}
		expiration := time.Now().Add(ttl)
		// Adding a small buffer since the TTL will be calculated again after this call
		// to ensure the database credential does not expire before the lease
		expiration = expiration.Add(5 * time.Second)

		newUserReq := v5.NewUserRequest{
			UsernameConfig: v5.UsernameMetadata{
				DisplayName: req.DisplayName,
				RoleName:    name,
			},
			Statements: v5.Statements{
				Commands: role.Statements.Creation,
			},
			RollbackStatements: v5.Statements{
				Commands: role.Statements.Rollback,
			},
			Expiration: expiration,
		}

		respData := make(map[string]interface{})

		// Generate the credential based on the role's credential type
		switch role.CredentialType {
		case v5.CredentialTypePassword:
			generator, err := newPasswordGenerator(role.CredentialConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to construct credential generator: %s", err)
			}

			// Fall back to database config-level password policy if not set on role
			if generator.PasswordPolicy == "" {
				generator.PasswordPolicy = dbConfig.PasswordPolicy
			}

			// Generate the password
			password, err := generator.generate(ctx, b, dbi.database)
			if err != nil {
				b.CloseIfShutdown(dbi, err)
				return nil, fmt.Errorf("failed to generate password: %s", err)
			}

			// Set input credential
			newUserReq.CredentialType = v5.CredentialTypePassword
			newUserReq.Password = password

		case v5.CredentialTypeRSAPrivateKey:
			generator, err := newRSAKeyGenerator(role.CredentialConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to construct credential generator: %s", err)
			}

			// Generate the RSA key pair
			public, private, err := generator.generate(b.GetRandomReader())
			if err != nil {
				return nil, fmt.Errorf("failed to generate RSA key pair: %s", err)
			}

			// Set input credential
			newUserReq.CredentialType = v5.CredentialTypeRSAPrivateKey
			newUserReq.PublicKey = public

			// Set output credential
			respData["rsa_private_key"] = string(private)
		case v5.CredentialTypeClientCertificate:
			generator, err := newClientCertificateGenerator(role.CredentialConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to construct credential generator: %s", err)
			}

			// Generate the client certificate
			cb, subject, err := generator.generate(b.GetRandomReader(), expiration,
				newUserReq.UsernameConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to generate client certificate: %w", err)
			}

			// Set input credential
			newUserReq.CredentialType = dbplugin.CredentialTypeClientCertificate
			newUserReq.Subject = subject

			// Set output credential
			respData["client_certificate"] = cb.Certificate
			respData["private_key"] = cb.PrivateKey
			respData["private_key_type"] = cb.PrivateKeyType
		}

		// Overwriting the password in the event this is a legacy database
		// plugin and the provided password is ignored
		newUserResp, password, err := dbi.database.NewUser(ctx, newUserReq)
		if err != nil {
			b.CloseIfShutdown(dbi, err)
			return nil, err
		}
		respData["username"] = newUserResp.Username

		// Database plugins using the v4 interface generate and return the password.
		// Set the password response to what is returned by the NewUser request.
		if role.CredentialType == v5.CredentialTypePassword {
			respData["password"] = password
		}

		internal := map[string]interface{}{
			"username":              newUserResp.Username,
			"role":                  name,
			"db_name":               role.DBName,
			"revocation_statements": role.Statements.Revocation,
		}
		resp := b.Secret(SecretCredsType).Response(respData, internal)
		resp.Secret.TTL = role.DefaultTTL
		resp.Secret.MaxTTL = role.MaxTTL
		return resp, nil
	}
}

func (b *databaseBackend) pathStaticCredsRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)

		role, err := b.StaticRole(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return logical.ErrorResponse("unknown role: %s", name), nil
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

		respData := map[string]interface{}{
			"username":            role.StaticAccount.Username,
			"ttl":                 role.StaticAccount.CredentialTTL().Seconds(),
			"last_vault_rotation": role.StaticAccount.LastVaultRotation,
		}

		if role.StaticAccount.UsesRotationPeriod() {
			respData["rotation_period"] = role.StaticAccount.RotationPeriod.Seconds()
		} else if role.StaticAccount.UsesRotationSchedule() {
			respData["rotation_schedule"] = role.StaticAccount.RotationSchedule
			if role.StaticAccount.RotationWindow.Seconds() != 0 {
				respData["rotation_window"] = role.StaticAccount.RotationWindow.Seconds()
			}
		}

		switch role.CredentialType {
		case v5.CredentialTypePassword:
			respData["password"] = role.StaticAccount.Password
		case v5.CredentialTypeRSAPrivateKey:
			respData["rsa_private_key"] = string(role.StaticAccount.PrivateKey)
		}

		return &logical.Response{
			Data: respData,
		}, nil
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

const pathStaticCredsReadHelpSyn = `
Request database credentials for a certain static role. These credentials are
rotated periodically.
`

const pathStaticCredsReadHelpDesc = `
This path reads database credentials for a certain static role. The database
credentials are rotated periodically according to their configuration, and will
return the same password until they are rotated.
`

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-2.0

package ldap

import (
	"context"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathConfigRotateRoot(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "config/rotate-root",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixLDAP,
			OperationVerb:   "rotate",
			OperationSuffix: "root-credentials",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathConfigRotateRootUpdate,
				ForwardPerformanceSecondary: true,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathConfigRotateRootHelpSyn,
		HelpDescription: pathConfigRotateRootHelpDesc,
	}
}

func (b *backend) pathConfigRotateRootUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// lock the backend's state - really just the config state - for mutating
	b.mu.Lock()
	defer b.mu.Unlock()

	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("attempted to rotate root on an undefined config"), nil
	}

	u, p := cfg.BindDN, cfg.BindPassword
	if u == "" || p == "" {
		return logical.ErrorResponse("auth is not using authenticated search, no root to rotate"), nil
	}

	// grab our ldap client
	client := ldaputil.Client{
		Logger: b.Logger(),
		LDAP:   ldaputil.NewLDAP(),
	}

	conn, err := client.DialLDAP(cfg.ConfigEntry)
	if err != nil {
		return nil, err
	}

	err = conn.Bind(u, p)
	if err != nil {
		return nil, err
	}

	lreq := &ldap.ModifyRequest{
		DN: cfg.BindDN,
	}

	var newPassword string
	if cfg.PasswordPolicy != "" {
		newPassword, err = b.System().GeneratePasswordFromPolicy(ctx, cfg.PasswordPolicy)
	} else {
		newPassword, err = base62.Random(defaultPasswordLength)
	}
	if err != nil {
		return nil, err
	}

	lreq.Replace("userPassword", []string{newPassword})

	err = conn.Modify(lreq)
	if err != nil {
		return nil, err
	}
	// update config with new password
	cfg.BindPassword = newPassword
	entry, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		// we might have to roll-back the password here?
		return nil, err
	}

	return nil, nil
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the LDAP credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the LDAP bindpass used by Vault for this mount.
`

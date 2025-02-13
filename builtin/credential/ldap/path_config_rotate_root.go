// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-2.0

package ldap

import (
	"context"
	"errors"

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

func (b *backend) pathConfigRotateRootUpdate(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	err := b.rotateRootCredential(ctx, req)
	var responseError responseError
	if errors.As(err, &responseError) {
		return logical.ErrorResponse(responseError.Error()), nil
	}

	// naturally this is `nil, nil` if the err is nil
	return nil, err
}

// responseError exists to capture the cases in the old rotate call that returned specific error responses
type responseError struct {
	error
}

func (b *backend) rotateRootCredential(ctx context.Context, req *logical.Request) error {
	// lock the backend's state - really just the config state - for mutating
	b.mu.Lock()
	defer b.mu.Unlock()

	cfg, err := b.Config(ctx, req)
	if err != nil {
		return err
	}
	if cfg == nil {
		return responseError{errors.New("attempted to rotate root on an undefined config")}
	}

	u, p := cfg.BindDN, cfg.BindPassword
	if u == "" || p == "" {
		// Logging this is as it may be useful to know that the binddn/bindpass is not set.
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth is not using authenticated search, no root to rotate")
		}
		return responseError{errors.New("auth is not using authenticated search, no root to rotate")}
	}

	// grab our ldap client
	client := ldaputil.Client{
		Logger: b.Logger(),
		LDAP:   ldaputil.NewLDAP(),
	}

	conn, err := client.DialLDAP(cfg.ConfigEntry)
	if err != nil {
		return err
	}

	err = conn.Bind(u, p)
	if err != nil {
		return err
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
		return err
	}

	lreq.Replace("userPassword", []string{newPassword})

	err = conn.Modify(lreq)
	if err != nil {
		return err
	}
	// update config with new password
	cfg.BindPassword = newPassword
	entry, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		// we might have to roll-back the password here?
		return err
	}

	return nil
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the LDAP credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the LDAP bindpass used by Vault for this mount.
`

// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-2.0

package ldap

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/text/encoding/unicode"
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
	if err != nil {
		// log here instead of inside the actual rotate call because the rotation manager also logs, so this is
		// the "equivalent" place for manual rotations.
		b.Logger().Error("failed to rotate root credential on user request", "path", req.Path, "error", err.Error())
	} else {
		// err is nil in this case
		b.Logger().Info("succesfully rotated root credential on user request", "path", req.Path)
	}
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

	// Create a copy of the config to modify for rotation
	rotateConfig := *cfg.ConfigEntry
	if cfg.RotationUrl != "" {
		rotateConfig.Url = cfg.RotationUrl
	}
	conn, err := client.DialLDAP(&rotateConfig)
	if err != nil {
		return err
	}

	err = conn.Bind(u, p)
	if err != nil {
		return err
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
	lreq, err := b.getModifyRequest(cfg, newPassword)
	if err != nil {
		return err
	}
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

func (b *backend) getModifyRequest(cfg *ldapConfigEntry, newPassword string) (*ldap.ModifyRequest, error) {
	lreq := &ldap.ModifyRequest{
		DN: cfg.BindDN,
	}
	switch cfg.RotationSchema {
	case SchemaOpenLDAP:
		lreq.Replace("userPassword", []string{newPassword})
	case SchemaRACF:
		// Password and password phrase management are mutually exclusive
		// operations. When the system is configured to manage one, it will not
		// modify the other.
		if cfg.RotationCredentialType == CredentialTypePhrase {
			lreq.Replace("racfPassPhrase", []string{newPassword})
		} else {
			lreq.Replace("racfPassword", []string{newPassword})
		}
		lreq.Replace("racfAttributes", []string{"noexpired"})
	case SchemaAD:
		utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		pwdEncoded, err := utf16.NewEncoder().String("\"" + newPassword + "\"")
		if err != nil {
			return nil, err
		}
		lreq.Replace("unicodePwd", []string{pwdEncoded})
	default:
		return nil, fmt.Errorf("configured schema %s not valid", cfg.RotationSchema)
	}
	return lreq, nil
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the LDAP credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the LDAP bindpass used by Vault for this mount.
`

// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-2.0

package ldap

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"strings"
	"unicode/utf16"

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
		Fields:          map[string]*framework.FieldSchema{},
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
		b.Logger().Info("successfully rotated root credential on user request", "path", req.Path)
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

	// Support both OpenLDAP and AD root rotation
	if cfg.Schema == "ad" {
		// AD requires unicodePwd, UTF-16LE encoding, and quoted string
		// AD password changes must be done over a secure connection (LDAPS)
		quotedPwd := "\"" + newPassword + "\""
		utf16Pwd := utf16leEncode(quotedPwd)
		lreq.Replace("unicodePwd", []string{utf16Pwd})

		// Log a warning if not using LDAPS for AD
		if !strings.HasPrefix(strings.ToLower(rotateConfig.Url), "ldaps://") {
			b.Logger().Warn("Active Directory password rotation should use LDAPS (ldaps://) for security")
		}
	} else {
		// OpenLDAP uses userPassword attribute with plain text
		lreq.Replace("userPassword", []string{newPassword})
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

func utf16leEncode(s string) string {
	utf := utf16.Encode([]rune(s))
	buf := new(bytes.Buffer)
	for _, v := range utf {
		binary.Write(buf, binary.LittleEndian, v)
	}
	return buf.String()
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the LDAP credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the LDAP bindpass used by Vault for this mount.
`

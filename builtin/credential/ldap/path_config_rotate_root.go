// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-2.0

package ldap

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
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

	// Validate TLS requirements for AD password rotation
	schema := ldaputil.NormalizedSchema(cfg.Schema)
	if schema == ldaputil.SchemaAD {
		// Validate URL(s) which will actually be used for rotation
		urlToValidate := cfg.Url
		if cfg.RotationUrl != "" {
			urlToValidate = cfg.RotationUrl
		}
		if err := validateADRotationURLs(urlToValidate, cfg.StartTLS); err != nil {
			return responseError{err}
		}
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
	// Close the connection when done to avoid leaking connections, especially during repeated rotation attempts.defer conn.Close()
	defer conn.Close()

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

	switch schema {
	case ldaputil.SchemaAD:
		// AD root rotation requires:
		// 1) Quoted password
		// 2) UTF-16LE encoding
		// 3) Encrypted connection (LDAPS/StartTLS)
		// Without these, AD rejects password updates.
		b.Logger().Debug("rotating root password using AD schema")
		quotedPwd := fmt.Sprintf("\"%s\"", newPassword)
		utf16Bytes := encodeUTF16LEBytes(quotedPwd)

		modReq := ldap.NewModifyRequest(cfg.BindDN, nil)
		modReq.Replace("unicodePwd", []string{string(utf16Bytes)})

		if err := conn.Modify(modReq); err != nil {
			return fmt.Errorf("failed to modify AD password for %q: %w", cfg.BindDN, err)
		}

	case ldaputil.SchemaOpenLDAP:
		b.Logger().Debug("rotating root password using openldap schema")
		lreq := &ldap.ModifyRequest{
			DN: cfg.BindDN,
		}
		lreq.Replace("userPassword", []string{newPassword})

		if err := conn.Modify(lreq); err != nil {
			return fmt.Errorf("failed to modify OpenLDAP password: %w", err)
		}
	default:
		return responseError{fmt.Errorf("unsupported schema type for password rotation: %s", schema)}
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

// encodeUTF16LEBytes encodes a string as UTF-16LE bytes for AD password changes.
// This encoding is required for Active Directory password changes via the unicodePwd attribute.
func encodeUTF16LEBytes(s string) []byte {
	utf16Runes := utf16.Encode([]rune(s))
	buf := make([]byte, len(utf16Runes)*2)
	for i, r := range utf16Runes {
		binary.LittleEndian.PutUint16(buf[i*2:], r)
	}
	return buf
}

// validateADRotationURLs validates that all URLs in the provided URL string
// meet the security requirements for AD password rotation.
func validateADRotationURLs(urlString string, startTLS bool) error {
	// AD password rotation requires encrypted connections (LDAPS or StartTLS)
	// Supported configurations:
	// 1. ldaps:// with proper certificate validation (recommended)
	// 2. ldaps:// with insecure_tls=true (skips certificate validation)
	// 3. ldap:// with starttls=true (with or without explicit certificates)
	if strings.TrimSpace(urlString) == "" {
		return errors.New("AD password rotation requires a configured URL")
	}
	// Split on commas to handle multiple URLs
	rawURLs := strings.Split(urlString, ",")
	hasNonTLSURL := false
	for _, rawURL := range rawURLs {
		rawURL := strings.TrimSpace(rawURL)
		if rawURL == "" {
			continue
		}
		urlLower := strings.ToLower(rawURL)
		isLDAPS := strings.HasPrefix(urlLower, "ldaps://")
		isLDAP := strings.HasPrefix(urlLower, "ldap://")

		// Validate that URL uses a supported protocol
		if !isLDAPS && !isLDAP {
			return fmt.Errorf("AD password rotation requires ldap:// or ldaps:// protocol, got: %s", rawURL)
		}
		// Track if any URL uses non-TLS ldap://
		if isLDAP {
			hasNonTLSURL = true
		}
	}
	// If any URL uses ldap:// (non-TLS), require StartTLS to ensure encryption
	if hasNonTLSURL && !startTLS {
		return errors.New("AD password rotation with ldap:// requires starttls=true for encrypted connection")
	}
	return nil
}

const pathConfigRotateRootHelpSyn = `
Request to rotate the LDAP credentials used by Vault
`

const pathConfigRotateRootHelpDesc = `
This path attempts to rotate the LDAP bindpass used by Vault for this mount.
`

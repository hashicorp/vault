// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/keysutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathKeysConfig() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/config",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "configure",
			OperationSuffix: "key",
		},

		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"min_decryption_version": {
				Type: framework.TypeInt,
				Description: `If set, the minimum version of the key allowed
to be decrypted. For signing keys, the minimum
version allowed to be used for verification.`,
			},

			"min_encryption_version": {
				Type: framework.TypeInt,
				Description: `If set, the minimum version of the key allowed
to be used for encryption; or for signing keys,
to be used for signing. If set to zero, only
the latest version of the key is allowed.`,
			},

			"deletion_allowed": {
				Type:        framework.TypeBool,
				Description: "Whether to allow deletion of the key",
			},

			"exportable": {
				Type:        framework.TypeBool,
				Description: `Enables export of the key. Once set, this cannot be disabled.`,
			},

			"allow_plaintext_backup": {
				Type:        framework.TypeBool,
				Description: `Enables taking a backup of the named key in plaintext format. Once set, this cannot be disabled.`,
			},

			"auto_rotate_period": {
				Type: framework.TypeDurationSecond,
				Description: `Amount of time the key should live before
being automatically rotated. A value of 0
disables automatic rotation for the key.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathKeysConfigWrite,
		},

		HelpSynopsis:    pathKeysConfigHelpSyn,
		HelpDescription: pathKeysConfigHelpDesc,
	}
}

func (b *backend) pathKeysConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (resp *logical.Response, retErr error) {
	name := d.Get("name").(string)

	// Check if the policy already exists before we lock everything
	p, _, err := b.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
	}, b.GetRandomReader())
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse(
				fmt.Sprintf("no existing key named %s could be found", name)),
			logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(true)
	}
	defer p.Unlock()

	var warning string

	originalMinDecryptionVersion := p.MinDecryptionVersion
	originalMinEncryptionVersion := p.MinEncryptionVersion
	originalDeletionAllowed := p.DeletionAllowed
	originalExportable := p.Exportable
	originalAllowPlaintextBackup := p.AllowPlaintextBackup

	defer func() {
		if retErr != nil || (resp != nil && resp.IsError()) {
			p.MinDecryptionVersion = originalMinDecryptionVersion
			p.MinEncryptionVersion = originalMinEncryptionVersion
			p.DeletionAllowed = originalDeletionAllowed
			p.Exportable = originalExportable
			p.AllowPlaintextBackup = originalAllowPlaintextBackup
		}
	}()

	persistNeeded := false

	minDecryptionVersionRaw, ok := d.GetOk("min_decryption_version")
	if ok {
		minDecryptionVersion := minDecryptionVersionRaw.(int)

		if minDecryptionVersion < 0 {
			return logical.ErrorResponse("min decryption version cannot be negative"), nil
		}

		if minDecryptionVersion == 0 {
			minDecryptionVersion = 1
			warning = "since Vault 0.3, transit key numbering starts at 1; forcing minimum to 1"
		}

		if minDecryptionVersion != p.MinDecryptionVersion {
			if minDecryptionVersion > p.LatestVersion {
				return logical.ErrorResponse(
					fmt.Sprintf("cannot set min decryption version of %d, latest key version is %d", minDecryptionVersion, p.LatestVersion)), nil
			}
			p.MinDecryptionVersion = minDecryptionVersion
			persistNeeded = true
		}
	}

	minEncryptionVersionRaw, ok := d.GetOk("min_encryption_version")
	if ok {
		minEncryptionVersion := minEncryptionVersionRaw.(int)

		if minEncryptionVersion < 0 {
			return logical.ErrorResponse("min encryption version cannot be negative"), nil
		}

		if minEncryptionVersion != p.MinEncryptionVersion {
			if minEncryptionVersion > p.LatestVersion {
				return logical.ErrorResponse(
					fmt.Sprintf("cannot set min encryption version of %d, latest key version is %d", minEncryptionVersion, p.LatestVersion)), nil
			}
			p.MinEncryptionVersion = minEncryptionVersion
			persistNeeded = true
		}
	}

	// Check here to get the final picture after the logic on each
	// individually. MinDecryptionVersion will always be 1 or above.
	if p.MinEncryptionVersion > 0 &&
		p.MinEncryptionVersion < p.MinDecryptionVersion {
		return logical.ErrorResponse(
			fmt.Sprintf("cannot set min encryption/decryption values; min encryption version of %d must be greater than or equal to min decryption version of %d", p.MinEncryptionVersion, p.MinDecryptionVersion)), nil
	}

	allowDeletionInt, ok := d.GetOk("deletion_allowed")
	if ok {
		allowDeletion := allowDeletionInt.(bool)
		if allowDeletion != p.DeletionAllowed {
			p.DeletionAllowed = allowDeletion
			persistNeeded = true
		}
	}

	// Add this as a guard here before persisting since we now require the min
	// decryption version to start at 1; even if it's not explicitly set here,
	// force the upgrade
	if p.MinDecryptionVersion == 0 {
		p.MinDecryptionVersion = 1
		persistNeeded = true
	}

	exportableRaw, ok := d.GetOk("exportable")
	if ok {
		exportable := exportableRaw.(bool)
		// Don't unset the already set value
		if exportable && !p.Exportable {
			p.Exportable = exportable
			persistNeeded = true
		}
	}

	allowPlaintextBackupRaw, ok := d.GetOk("allow_plaintext_backup")
	if ok {
		allowPlaintextBackup := allowPlaintextBackupRaw.(bool)
		// Don't unset the already set value
		if allowPlaintextBackup && !p.AllowPlaintextBackup {
			p.AllowPlaintextBackup = allowPlaintextBackup
			persistNeeded = true
		}
	}

	autoRotatePeriodRaw, ok, err := d.GetOkErr("auto_rotate_period")
	if err != nil {
		return nil, err
	}
	if ok {
		autoRotatePeriod := time.Second * time.Duration(autoRotatePeriodRaw.(int))
		// Provided value must be 0 to disable or at least an hour
		if autoRotatePeriod != 0 && autoRotatePeriod < time.Hour {
			return logical.ErrorResponse("auto rotate period must be 0 to disable or at least an hour"), nil
		}

		if autoRotatePeriod != p.AutoRotatePeriod {
			p.AutoRotatePeriod = autoRotatePeriod
			persistNeeded = true
		}

		if p.Type == keysutil.KeyType_MANAGED_KEY && autoRotatePeriod != 0 {
			return logical.ErrorResponse("Auto rotation can not be set for managed keys"), nil
		}
	}

	if !persistNeeded {
		resp, err := b.formatKeyPolicy(p, nil)
		if err != nil {
			return nil, err
		}
		if warning != "" {
			resp.AddWarning(warning)
		}
		return resp, nil
	}

	switch {
	case p.MinAvailableVersion > p.MinEncryptionVersion:
		return logical.ErrorResponse("min encryption version should not be less than min available version"), nil
	case p.MinAvailableVersion > p.MinDecryptionVersion:
		return logical.ErrorResponse("min decryption version should not be less then min available version"), nil
	}

	if err := p.Persist(ctx, req.Storage); err != nil {
		return nil, err
	}

	resp, err = b.formatKeyPolicy(p, nil)
	if err != nil {
		return nil, err
	}
	if warning != "" {
		resp.AddWarning(warning)
	}
	return resp, nil
}

const pathKeysConfigHelpSyn = `Configure a named encryption key`

const pathKeysConfigHelpDesc = `
This path is used to configure the named key. Currently, this
supports adjusting the minimum version of the key allowed to
be used for decryption via the min_decryption_version parameter.
`

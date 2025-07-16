// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathRestore() *framework.Path {
	return &framework.Path{
		Pattern: "restore" + framework.OptionalParamRegex("name"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixTransit,
			OperationVerb:   "restore",
			OperationSuffix: "key|and-rename-key",
		},

		Fields: map[string]*framework.FieldSchema{
			"backup": {
				Type:        framework.TypeString,
				Description: "Backed up key data to be restored. This should be the output from the 'backup/' endpoint.",
			},
			"name": {
				Type:        framework.TypeString,
				Description: "If set, this will be the name of the restored key.",
			},
			"force": {
				Type:        framework.TypeBool,
				Description: "If set and a key by the given name exists, force the restore operation and override the key.",
				Default:     false,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRestoreUpdate,
		},

		HelpSynopsis:    pathRestoreHelpSyn,
		HelpDescription: pathRestoreHelpDesc,
	}
}

func (b *backend) pathRestoreUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	backupB64 := d.Get("backup").(string)
	force := d.Get("force").(bool)
	if backupB64 == "" {
		return logical.ErrorResponse("'backup' must be supplied"), nil
	}

	// If a name is given, make sure it does not contain any slashes. The Transit
	// secret engine does not allow sub-paths in key names
	keyName := d.Get("name").(string)
	if strings.Contains(keyName, "/") {
		return nil, ErrInvalidKeyName
	}

	return nil, b.lm.RestorePolicy(ctx, req.Storage, keyName, backupB64, force)
}

const (
	pathRestoreHelpSyn  = `Restore the named key`
	pathRestoreHelpDesc = `This path is used to restore the named key.`
)

var ErrInvalidKeyName = errors.New("key names cannot be paths")

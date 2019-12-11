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
		Fields: map[string]*framework.FieldSchema{
			"backup": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Backed up key data to be restored. This should be the output from the 'backup/' endpoint.",
			},
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "If set, this will be the name of the restored key.",
			},
			"force": &framework.FieldSchema{
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

	keyName := d.Get("name").(string)
	// if a name is given, make sure it does not contain any slashes and look like
	// a path
	if keyName != "" {
		if strings.Contains(keyName, "/") {
			return nil, ErrInvalidKeyName
		}
	}

	return nil, b.lm.RestorePolicy(ctx, req.Storage, keyName, backupB64, force)
}

const pathRestoreHelpSyn = `Restore the named key`
const pathRestoreHelpDesc = `This path is used to restore the named key.`

var ErrInvalidKeyName = errors.New("key names cannot be paths")

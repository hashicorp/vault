package transit

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRestoreUpdate,
		},

		HelpSynopsis:    pathRestoreHelpSyn,
		HelpDescription: pathRestoreHelpDesc,
	}
}

func (b *backend) pathRestoreUpdate(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	backupB64 := d.Get("backup").(string)
	if backupB64 == "" {
		return logical.ErrorResponse("'backup' must be supplied"), nil
	}

	return nil, b.lm.RestorePolicy(req.Storage, d.Get("name").(string), backupB64)
}

const pathRestoreHelpSyn = `Restore the named key`
const pathRestoreHelpDesc = `This path is used to restore the named key.`

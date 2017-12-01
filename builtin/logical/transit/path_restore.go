package transit

import (
	"encoding/base64"

	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathRestore() *framework.Path {
	return &framework.Path{
		Pattern: "restore" + framework.OptionalParamRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"backup": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Backed up key data to be restored. This should be the output from the 'backup/' endpoint",
			},
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name to be assigned to the restored key",
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

	backupBytes, err := base64.StdEncoding.DecodeString(backupB64)
	if err != nil {
		return nil, err
	}

	var keyData keysutil.KeyData
	keyData.Policy = &keysutil.Policy{
		Keys: keysutil.KeyEntryMap{},
	}
	err = jsonutil.DecodeJSON(backupBytes, &keyData)
	if err != nil {
		return nil, err
	}

	name := d.Get("name").(string)
	if name != "" {
		keyData.Policy.Name = name
	}

	err = b.lm.RestorePolicy(req.Storage, keyData)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

const pathRestoreHelpSyn = `Restore the named key`
const pathRestoreHelpDesc = `This path is used to restore the named key.`

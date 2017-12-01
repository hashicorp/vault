package transit

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathBackup() *framework.Path {
	return &framework.Path{
		Pattern: "backup/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathBackupRead,
		},

		HelpSynopsis:    pathBackupHelpSyn,
		HelpDescription: pathBackupHelpDesc,
	}
}

func (b *backend) pathBackupRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("empty 'name'"), nil
	}

	p, lock, err := b.lm.GetPolicyExclusive(req.Storage, name)
	if lock != nil {
		defer lock.Unlock()
	}
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid key %q", name)), nil
	}

	backup, err := p.Backup(req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"backup": backup,
		},
	}, nil
}

const pathBackupHelpSyn = `Backup the named key`
const pathBackupHelpDesc = `This path is used to backup the named key.`

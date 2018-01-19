package transit

import (
	"context"

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

func (b *backend) pathBackupRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	backup, err := b.lm.BackupPolicy(ctx, req.Storage, d.Get("name").(string))
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

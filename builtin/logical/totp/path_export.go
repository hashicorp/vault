package totp

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathExport(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "export/" + framework.GenericNameWithAtRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the key.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathExportRead,
		},

		HelpSynopsis:    pathExportHelpSyn,
		HelpDescription: pathExportHelpDesc,
	}
}

func (b *backend) pathExportRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	key, err := b.Key(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, nil
	}

	// Translate algorithm back to string
	algorithm := key.Algorithm.String()

	// Return values of key
	return &logical.Response{
		Data: map[string]interface{}{
			"key":          key.Key,
			"issuer":       key.Issuer,
			"account_name": key.AccountName,
			"period":       key.Period,
			"algorithm":    algorithm,
			"digits":       key.Digits,
		},
	}, nil
}

const pathExportHelpSyn = `
Export the keys data and secret.
`

const pathExportHelpDesc = `
This path lets you export the keys.
`

package mongodbatlas

import (
	"context"
	"errors"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *Backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: "config",
		Fields: map[string]*framework.FieldSchema{
			"public_key": {
				Type:        framework.TypeString,
				Description: "MongoDB Atlas Programmatic Public Key",
				Required:    true,
			},
			"private_key": {
				Type:        framework.TypeString,
				Description: "MongoDB Atlas Programmatic Private Key",
				Required:    true,
				DisplayAttrs: &framework.DisplayAttributes{
					Sensitive: true,
				},
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
			logical.ReadOperation:   b.pathConfigRead,
		},
		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
}

func (b *Backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	publicKey := data.Get("public_key").(string)
	if publicKey == "" {
		return nil, errors.New("public_key is empty")
	}

	privateKey := data.Get("private_key").(string)
	if privateKey == "" {
		return nil, errors.New("private_key is empty")
	}

	entry, err := logical.StorageEntryJSON("config", config{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// Clean cached client (if any)
	b.client = nil

	return nil, nil
}

func (b *Backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cfg, err := getRootConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"public_key": cfg.PublicKey,
		},
	}, nil
}

type config struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

const pathConfigHelpSyn = `
Configure the  credentials that are used to manage Database Users.
`

const pathConfigHelpDesc = `
Before doing anything, the Atlas backend needs credentials that are able
to manage databaseusers, access keys, etc. This endpoint is used to 
configure those credentials.
`

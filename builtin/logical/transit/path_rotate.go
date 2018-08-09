package transit

import (
	"context"

	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathRotate() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/rotate",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathRotateWrite,
		},

		HelpSynopsis:    pathRotateHelpSyn,
		HelpDescription: pathRotateHelpDesc,
	}
}

func (b *backend) pathRotateWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Get the policy
	p, _, err := b.lm.GetPolicy(ctx, keysutil.PolicyRequest{
		Storage: req.Storage,
		Name:    name,
	})
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse("key not found"), logical.ErrInvalidRequest
	}
	if !b.System().CachingDisabled() {
		p.Lock(true)
	}

	// Rotate the policy
	err = p.Rotate(ctx, req.Storage)

	p.Unlock()
	return nil, err
}

const pathRotateHelpSyn = `Rotate named encryption key`

const pathRotateHelpDesc = `
This path is used to rotate the named key. After rotation,
new encryption requests using this name will use the new key,
but decryption will still be supported for older versions.
`

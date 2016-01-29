package transit

import (
	"fmt"

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

func (b *backend) pathRotateWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Get the policy
	lp, err := b.policies.getPolicy(req, name)
	if err != nil {
		return nil, err
	}

	// Error if invalid policy
	if lp == nil {
		return logical.ErrorResponse("policy not found"), logical.ErrInvalidRequest
	}

	lp.Lock()
	defer lp.Unlock()

	// Verify if wasn't deleted before we grabbed the lock
	if lp.policy == nil {
		return nil, fmt.Errorf("no existing policy named %s could be found", name)
	}

	// Generate the policy
	err = lp.policy.rotate(req.Storage)

	return nil, err
}

const pathRotateHelpSyn = `Rotate named encryption key`

const pathRotateHelpDesc = `
This path is used to rotate the named key. After rotation,
new encryption requests using this name will use the new key,
but decryption will still be supported for older versions.
`

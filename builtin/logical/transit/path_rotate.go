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
	lp, err := b.policies.getPolicy(req.Storage, name)
	if err != nil {
		return nil, err
	}

	// Error if invalid policy
	if lp == nil {
		return logical.ErrorResponse("key not found"), logical.ErrInvalidRequest
	}

	// Store so we can detect later if this has changed out from under us
	keyVersion := lp.Policy().LatestVersion

	// lock the policies object so we can refresh
	b.policies.Lock()
	defer b.policies.Unlock()
	lp.Lock()
	defer lp.Unlock()

	// Refresh in case it's changed since before we grabbed the lock
	lp, err = b.policies.refreshPolicy(req.Storage, name)
	if err != nil {
		return nil, err
	}
	if lp == nil {
		return nil, fmt.Errorf("error finding key %s after locking for changes", name)
	}

	// Verify if wasn't deleted before we grabbed the lock
	if lp.Policy() == nil {
		return nil, fmt.Errorf("no existing key named %s could be found", name)
	}

	// Make sure that the policy hasn't been rotated simultaneously
	if keyVersion != lp.Policy().LatestVersion {
		resp := &logical.Response{}
		resp.AddWarning("key has been rotated since this endpoint was called; did not perform rotation")
		return resp, nil
	}

	// Rotate the policy
	err = lp.Policy().rotate(req.Storage)

	return nil, err
}

const pathRotateHelpSyn = `Rotate named encryption key`

const pathRotateHelpDesc = `
This path is used to rotate the named key. After rotation,
new encryption requests using this name will use the new key,
but decryption will still be supported for older versions.
`

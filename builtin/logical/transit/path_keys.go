package transit

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathKeys() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},

			"derived": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Enables key derivation mode. This allows for per-transaction unique keys",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathPolicyWrite,
			logical.DeleteOperation: b.pathPolicyDelete,
			logical.ReadOperation:   b.pathPolicyRead,
		},

		HelpSynopsis:    pathPolicyHelpSyn,
		HelpDescription: pathPolicyHelpDesc,
	}
}

func (b *backend) pathPolicyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	derived := d.Get("derived").(bool)

	// Check if the policy already exists
	existing, err := b.policies.getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, nil
	}

	// Generate the policy
	_, err = b.policies.generatePolicy(req.Storage, name, derived)
	return nil, err
}

func (b *backend) pathPolicyRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	lp, err := b.policies.getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if lp == nil {
		return nil, nil
	}

	lp.RLock()
	defer lp.RUnlock()

	// Verify if wasn't deleted before we grabbed the lock
	if lp.policy == nil {
		return nil, fmt.Errorf("no existing policy named %s could be found", name)
	}

	// Return the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"name":                   lp.policy.Name,
			"cipher_mode":            lp.policy.CipherMode,
			"derived":                lp.policy.Derived,
			"deletion_allowed":       lp.policy.DeletionAllowed,
			"min_decryption_version": lp.policy.MinDecryptionVersion,
			"latest_version":         lp.policy.LatestVersion,
		},
	}
	if lp.policy.Derived {
		resp.Data["kdf_mode"] = lp.policy.KDFMode
	}

	retKeys := map[string]int64{}
	for k, v := range lp.policy.Keys {
		retKeys[strconv.Itoa(k)] = v.CreationTime
	}
	resp.Data["keys"] = retKeys

	return resp, nil
}

func (b *backend) pathPolicyDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	lp, err := b.policies.getPolicy(req, name)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error looking up policy %s, error is %s", name, err)), err
	}
	if lp == nil {
		return logical.ErrorResponse(fmt.Sprintf("no such key %s", name)), logical.ErrInvalidRequest
	}

	err = b.policies.deletePolicy(req.Storage, name)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("error deleting policy %s: %s", name, err)), err
	}

	return nil, nil
}

const pathPolicyHelpSyn = `Managed named encryption keys`

const pathPolicyHelpDesc = `
This path is used to manage the named keys that are available.
Doing a write with no value against a new named key will create
it using a randomly generated key.
`

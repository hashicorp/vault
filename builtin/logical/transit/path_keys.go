package transit

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathKeys() *framework.Path {
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
			logical.WriteOperation:  pathPolicyWrite,
			logical.DeleteOperation: pathPolicyDelete,
			logical.ReadOperation:   pathPolicyRead,
		},

		HelpSynopsis:    pathPolicyHelpSyn,
		HelpDescription: pathPolicyHelpDesc,
	}
}

func pathPolicyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	derived := d.Get("derived").(bool)

	// Check if the policy already exists
	existing, err := getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, nil
	}

	// Generate the policy
	_, err = generatePolicy(req.Storage, name, derived)
	return nil, err
}

func pathPolicyRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	p, err := getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	// Return the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"name":        p.Name,
			"cipher_mode": p.CipherMode,
			"derived":     p.Derived,
			"disabled":    p.Disabled,
		},
	}
	if p.Derived {
		resp.Data["kdf_mode"] = p.KDFMode
	}
	return resp, nil
}

func pathPolicyDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	p, err := getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return logical.ErrorResponse(fmt.Sprintf("no such key %s", name)), logical.ErrInvalidRequest
	}

	if !p.Disabled {
		return logical.ErrorResponse(fmt.Sprintf("key must be disabled before deletion")), logical.ErrInvalidRequest
	}

	err = req.Storage.Delete("policy/" + name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

const pathPolicyHelpSyn = `Managed named encryption keys`

const pathPolicyHelpDesc = `
This path is used to manage the named keys that are available.
Doing a write with no value against a new named key will create
it using a randomly generated key.
`

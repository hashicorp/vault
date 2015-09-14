package transit

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathEnable() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/enable",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathEnableWrite,
		},

		HelpSynopsis:    pathEnableHelpSyn,
		HelpDescription: pathEnableHelpDesc,
	}
}

func pathDisable() *framework.Path {
	return &framework.Path{
		Pattern: "keys/" + framework.GenericNameRegex("name") + "/disable",
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the key",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathDisableWrite,
		},

		HelpSynopsis:    pathDisableHelpSyn,
		HelpDescription: pathDisableHelpDesc,
	}
}

func pathEnableWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Check if the policy already exists
	policy, err := getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return logical.ErrorResponse(
				fmt.Sprintf("no existing role named %s could be found", name)),
			logical.ErrInvalidRequest
	}

	if !policy.Disabled {
		return nil, nil
	}

	policy.Disabled = false

	return nil, policy.Persist(req.Storage, name)
}

func pathDisableWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Check if the policy already exists
	policy, err := getPolicy(req, name)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return logical.ErrorResponse(
				fmt.Sprintf("no existing role named %s could be found", name)),
			logical.ErrInvalidRequest
	}

	if policy.Disabled {
		return nil, nil
	}

	policy.Disabled = true

	return nil, policy.Persist(req.Storage, name)
}

const pathEnableHelpSyn = `Enable a named encryption key`

const pathEnableHelpDesc = `
This path is used to enable the named key. After enabling,
the key will be available for use for encryption.
`

const pathDisableHelpSyn = `Disable a named encryption key`

const pathDisableHelpDesc = `
This path is used to disable the named key. After disabling,
the key cannot be used to encrypt values. This is useful when
when switching to a new named key, but wanting to be able to
decrypt against old keys while guarding against additional
data being encrypted with the old key.
`

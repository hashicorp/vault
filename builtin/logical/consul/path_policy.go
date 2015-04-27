package consul

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathPolicy() *framework.Path {
	return &framework.Path{
		Pattern: `policy/(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the policy",
			},

			"policy": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Policy document, base64 encoded.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   pathPolicyRead,
			logical.WriteOperation:  pathPolicyWrite,
			logical.DeleteOperation: pathPolicyDelete,
		},
	}
}

func pathPolicyRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Read the policy
	policy, err := req.Storage.Get("policy/" + name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving policy: %s", err)
	}
	if policy == nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Policy '%s' not found", name)), nil
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"policy": base64.StdEncoding.EncodeToString(policy.Value),
		},
	}
	return resp, nil
}

func pathPolicyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	policyRaw, err := base64.StdEncoding.DecodeString(d.Get("policy").(string))
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error decoding policy base64: %s", err)), nil
	}

	// Write the policy into storage
	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "policy/" + d.Get("name").(string),
		Value: policyRaw,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func pathPolicyDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if err := req.Storage.Delete("policy/" + name); err != nil {
		return nil, err
	}
	return nil, nil
}

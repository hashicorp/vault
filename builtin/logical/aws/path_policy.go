package aws

import (
	"bytes"
	"encoding/json"
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
				Description: "IAM policy document",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: pathPolicyDelete,
			logical.ReadOperation:   pathPolicyRead,
			logical.WriteOperation:  pathPolicyWrite,
		},

		HelpSynopsis:    pathPolicyHelpSyn,
		HelpDescription: pathPolicyHelpDesc,
	}
}

func pathPolicyDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete("policy/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func pathPolicyRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := req.Storage.Get("policy/" + d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"policy": string(entry.Value),
		},
	}, nil
}

func pathPolicyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(d.Get("policy").(string))); err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error compacting policy: %s", err)), nil
	}

	// Write the policy into storage
	err := req.Storage.Put(&logical.StorageEntry{
		Key:   "policy/" + d.Get("name").(string),
		Value: buf.Bytes(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

const pathPolicyHelpSyn = `
Read and write IAM policies that access keys can be made for.
`

const pathPolicyHelpDesc = `
This path allows you to read and write policies that are used to
create access keys. These policies map directly to the route to read the
access keys. For example, if the backend is mounted at "aws" and you
wrote a policy to "aws/policy/deploy" then a user could request access
credentials at "aws/deploy".

The policies written are normal IAM policies. Vault will not attempt to
parse these except to validate that they're basic JSON. To validate the
keys, attempt to read an access key after writing the policy.
`

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
				Description: "Policy document",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathPolicyWrite,
		},
	}
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

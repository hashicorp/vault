package aws

import (
	"bytes"
	"encoding/base64"
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
				Description: "Policy document, base64 encoded.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: pathPolicyWrite,
		},
	}
}

func pathPolicyWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Decode and compact the policy. AWS requires a JSON-compacted policy
	// because it mustn't contain newlines.
	var policyBuf bytes.Buffer
	policyRaw, err := base64.StdEncoding.DecodeString(d.Get("policy").(string))
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error decoding policy base64: %s", err)), nil
	}
	if err := json.Compact(&policyBuf, []byte(policyRaw)); err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error compacting policy: %s", err)), nil
	}

	// Write the policy into storage
	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "policy/" + d.Get("name").(string),
		Value: policyBuf.Bytes(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

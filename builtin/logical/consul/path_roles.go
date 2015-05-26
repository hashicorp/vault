package consul

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathRoles() *framework.Path {
	return &framework.Path{
		Pattern: `roles/(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},

			"policy": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Policy document, base64 encoded.",
			},

			"lease": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Lease time of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   pathRolesRead,
			logical.WriteOperation:  pathRolesWrite,
			logical.DeleteOperation: pathRolesDelete,
		},
	}
}

func pathRolesRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	// Read the policy
	policy, err := req.Storage.Get("policy/" + name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %s", err)
	}
	if policy == nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Role '%s' not found", name)), nil
	}

	leaseRaw, err := req.Storage.Get("policy/" + name + "/lease")
	if err != nil {
		return nil, fmt.Errorf("error retrieving lease: %s", err)
	}
	lease, err := time.ParseDuration(string(leaseRaw.Value))
	if err != nil {
		return nil, fmt.Errorf("error retrieving lease: %s", err)
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"policy": base64.StdEncoding.EncodeToString(policy.Value),
			"lease":  lease.String(),
		},
	}
	return resp, nil
}

func pathRolesWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	policyRaw, err := base64.StdEncoding.DecodeString(d.Get("policy").(string))
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error decoding policy base64: %s", err)), nil
	}

	// Write the policy into storage
	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "policy/" + name,
		Value: policyRaw,
	})
	if err != nil {
		return nil, err
	}

	// Write the policy lease into storage
	lease, err := time.ParseDuration(d.Get("lease").(string))
	if err != nil || lease == time.Duration(0) {
		lease = DefaultLeaseDuration
	}
	err = req.Storage.Put(&logical.StorageEntry{
		Key:   "policy/" + name + "/lease",
		Value: []byte(lease.String()),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func pathRolesDelete(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if err := req.Storage.Delete("policy/" + name); err != nil {
		return nil, err
	}
	return nil, nil
}

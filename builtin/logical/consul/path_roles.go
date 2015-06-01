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

	entry, err := req.Storage.Get("policy/" + name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result roleConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"policy": base64.StdEncoding.EncodeToString([]byte(result.Policy)),
			"lease":  result.Lease.String(),
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
	lease, err := time.ParseDuration(d.Get("lease").(string))
	if err != nil || lease == time.Duration(0) {
		lease = DefaultLeaseDuration
	}

	entry, err := logical.StorageEntryJSON("policy/"+name, roleConfig{
		Policy: string(policyRaw),
		Lease:  lease,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(entry); err != nil {
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

type roleConfig struct {
	Policy string        `json:"policy"`
	Lease  time.Duration `json:"lease"`
}

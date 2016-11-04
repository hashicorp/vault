package consul

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},
	}
}

func pathRoles() *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},

			"policy": &framework.FieldSchema{
				Type: framework.TypeString,
				Description: `Policy document, base64 encoded. Required
for 'client' tokens.`,
			},

			"token_type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "client",
				Description: `Which type of token to create: 'client'
or 'management'. If a 'management' token,
the "policy" parameter is not required.
Defaults to 'client'.`,
			},

			"lease": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Lease time of the role.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   pathRolesRead,
			logical.UpdateOperation: pathRolesWrite,
			logical.DeleteOperation: pathRolesDelete,
		},
	}
}

func (b *backend) pathRoleList(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List("policy/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
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

	if result.TokenType == "" {
		result.TokenType = "client"
	}

	// Generate the response
	resp := &logical.Response{
		Data: map[string]interface{}{
			"lease":      result.Lease.String(),
			"token_type": result.TokenType,
		},
	}
	if result.Policy != "" {
		resp.Data["policy"] = base64.StdEncoding.EncodeToString([]byte(result.Policy))
	}
	return resp, nil
}

func pathRolesWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	tokenType := d.Get("token_type").(string)

	switch tokenType {
	case "client":
	case "management":
	default:
		return logical.ErrorResponse(
			"token_type must be \"client\" or \"management\""), nil
	}

	name := d.Get("name").(string)
	policy := d.Get("policy").(string)
	var policyRaw []byte
	var err error
	if tokenType != "management" {
		if policy == "" {
			return logical.ErrorResponse(
				"policy cannot be empty when not using management tokens"), nil
		}
		policyRaw, err = base64.StdEncoding.DecodeString(d.Get("policy").(string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"Error decoding policy base64: %s", err)), nil
		}
	}

	var lease time.Duration
	leaseParam := d.Get("lease").(string)
	if leaseParam != "" {
		lease, err = time.ParseDuration(leaseParam)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf(
				"error parsing given lease of %s: %s", leaseParam, err)), nil
		}
	}

	entry, err := logical.StorageEntryJSON("policy/"+name, roleConfig{
		Policy:    string(policyRaw),
		Lease:     lease,
		TokenType: tokenType,
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
	Policy    string        `json:"policy"`
	Lease     time.Duration `json:"lease"`
	TokenType string        `json:"token_type"`
}

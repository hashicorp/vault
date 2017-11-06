package nomad

import (
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},
	}
}

func pathRoles() *framework.Path {
	return &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},

			"policy": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Policy name as previously created in Nomad. Required",
			},

			"global": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Policy name as previously created in Nomad. Required",
			},

			"token_type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "client",
				Description: `Which type of token to create: 'client'
or 'management'. If a 'management' token,
the "policy" parameter is not required.
Defaults to 'client'.`,
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
	entries, err := req.Storage.List("role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func pathRolesRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get("role/" + name)
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
			"token_type": result.TokenType,
			"global":     result.Global,
		},
	}
	if len(result.Policy) != 0 {
		resp.Data["policy"] = result.Policy
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
	global := d.Get("global").(bool)
	policy := d.Get("policy").([]string)
	var err error
	if tokenType != "management" {
		if len(policy) == 0 {
			return logical.ErrorResponse(
				"policy cannot be empty when not using management tokens"), nil
		}
	} else {
		if len(policy) != 0 {
			return logical.ErrorResponse(
				"policy should be empty when using management tokens"), nil
		}
	}

	entry, err := logical.StorageEntryJSON("role/"+name, roleConfig{
		Policy:    policy,
		TokenType: tokenType,
		Global:    global,
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
	if err := req.Storage.Delete("role/" + name); err != nil {
		return nil, err
	}
	return nil, nil
}

type roleConfig struct {
	Policy    []string `json:"policy"`
	TokenType string   `json:"token_type"`
	Global    bool     `json:"global"`
}

package consul

import (
	"context"
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

func pathRoles(b *backend) *framework.Path {
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
for 'client' tokens. Required for Consul pre-1.4.`,
			},

			"policies": &framework.FieldSchema{
				Type: framework.TypeCommaStringSlice,
				Description: `List of policies to attach to the token. Required
for Consul 1.4 or above.`,
			},

			"local": &framework.FieldSchema{
				Type: framework.TypeBool,
				Description: `Indicates that the token should not be replicated globally 
and instead be local to the current datacenter.  Available in Consul 1.4 and above.`,
			},

			"token_type": &framework.FieldSchema{
				Type:    framework.TypeString,
				Default: "client",
				Description: `Which type of token to create: 'client'
or 'management'. If a 'management' token,
the "policy" parameter is not required.
Defaults to 'client'.`,
			},

			"ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: "TTL for the Consul token created from the role.",
			},

			"max_ttl": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: "Max TTL for the Consul token created from the role.",
			},

			"lease": &framework.FieldSchema{
				Type:        framework.TypeDurationSecond,
				Description: "DEPRECATED: Use ttl.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathRolesRead,
			logical.UpdateOperation: b.pathRolesWrite,
			logical.DeleteOperation: b.pathRolesDelete,
		},
	}
}

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "policy/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *backend) pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, "policy/"+name)
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
			"lease":      int64(result.TTL.Seconds()),
			"ttl":        int64(result.TTL.Seconds()),
			"max_ttl":    int64(result.MaxTTL.Seconds()),
			"token_type": result.TokenType,
			"local":      result.Local,
		},
	}
	if result.Policy != "" {
		resp.Data["policy"] = base64.StdEncoding.EncodeToString([]byte(result.Policy))
	}
	if len(result.Policies) > 0 {
		resp.Data["policies"] = result.Policies
	}
	return resp, nil
}

func (b *backend) pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	tokenType := d.Get("token_type").(string)
	policy := d.Get("policy").(string)
	name := d.Get("name").(string)
	policies := d.Get("policies").([]string)
	local := d.Get("local").(bool)

	if len(policies) == 0 {
		switch tokenType {
		case "client":
			if policy == "" {
				return logical.ErrorResponse(
					"Use either a policy document, or a list of policies, depending on your Consul version"), nil
			}
		case "management":
		default:
			return logical.ErrorResponse(
				"token_type must be \"client\" or \"management\""), nil
		}
	}

	policyRaw, err := base64.StdEncoding.DecodeString(policy)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error decoding policy base64: %s", err)), nil
	}

	var ttl time.Duration
	ttlRaw, ok := d.GetOk("ttl")
	if ok {
		ttl = time.Second * time.Duration(ttlRaw.(int))
	} else {
		leaseParamRaw, ok := d.GetOk("lease")
		if ok {
			ttl = time.Second * time.Duration(leaseParamRaw.(int))
		}
	}

	var maxTTL time.Duration
	maxTTLRaw, ok := d.GetOk("max_ttl")
	if ok {
		maxTTL = time.Second * time.Duration(maxTTLRaw.(int))
	}

	entry, err := logical.StorageEntryJSON("policy/"+name, roleConfig{
		Policy:    string(policyRaw),
		Policies:  policies,
		TokenType: tokenType,
		TTL:       ttl,
		MaxTTL:    maxTTL,
		Local:     local,
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *backend) pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if err := req.Storage.Delete(ctx, "policy/"+name); err != nil {
		return nil, err
	}
	return nil, nil
}

type roleConfig struct {
	Policy    string        `json:"policy"`
	Policies  []string      `json:"policies"`
	TTL       time.Duration `json:"lease"`
	MaxTTL    time.Duration `json:"max_ttl"`
	TokenType string        `json:"token_type"`
	Local     bool          `json:"local"`
}

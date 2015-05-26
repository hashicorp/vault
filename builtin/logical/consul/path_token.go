package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathToken(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `creds/(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathTokenRead,
		},
	}
}

func (b *backend) pathTokenRead(
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

	// Get the consul client
	c, err := client(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Generate a random name for the token
	tokenName := fmt.Sprintf("Vault %s %d", req.DisplayName, time.Now().Unix())
	// Create it
	token, _, err := c.ACL().Create(&api.ACLEntry{
		Name:  tokenName,
		Type:  "client",
		Rules: string(policy.Value),
	}, nil)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Use the helper to create the secret
	s := b.Secret(SecretTokenType)
	if leaseRaw != nil {
		lease, err := time.ParseDuration(string(leaseRaw.Value))
		if err == nil {
			s.DefaultDuration = lease
		}
	}
	return s.Response(map[string]interface{}{
		"token": token,
	}, nil), nil
}

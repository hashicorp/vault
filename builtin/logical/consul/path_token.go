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
		Pattern: `(?P<name>\w+)`,
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the policy",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathTokenRead,
		},
	}
}

func (b *backend) pathTokenRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	policyName := d.Get("name").(string)

	// Generate a random name for the token
	name := fmt.Sprintf("Vault %s %d", req.DisplayName, time.Now().Unix())

	// Read the policy
	policy, err := req.Storage.Get("policy/" + policyName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving policy: %s", err)
	}
	if policy == nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Policy '%s' not found", policyName)), nil
	}

	// Get the consul client
	c, err := client(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Create it
	token, _, err := c.ACL().Create(&api.ACLEntry{
		Name:  name,
		Type:  "client",
		Rules: string(policy.Value),
	}, nil)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Use the helper to create the secret
	return b.Secret(SecretTokenType).Response(map[string]interface{}{
		"token": token,
	}, nil), nil
}

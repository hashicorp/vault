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
		Pattern: "creds/" + framework.GenericNameRegex("name"),
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

	entry, err := req.Storage.Get("policy/" + name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %s", err)
	}
	if entry == nil {
		return logical.ErrorResponse(fmt.Sprintf("Role '%s' not found", name)), nil
	}

	var result roleConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	if result.TokenType == "" {
		result.TokenType = "client"
	}

	// Get the consul client
	c, userErr, intErr := client(req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		return logical.ErrorResponse(userErr.Error()), nil
	}

	// Generate a name for the token
	tokenName := fmt.Sprintf("Vault %s %s %d", name, req.DisplayName, time.Now().UnixNano())

	// Create it
	token, _, err := c.ACL().Create(&api.ACLEntry{
		Name:  tokenName,
		Type:  result.TokenType,
		Rules: result.Policy,
	}, nil)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Use the helper to create the secret
	s := b.Secret(SecretTokenType).Response(map[string]interface{}{
		"token": token,
	}, map[string]interface{}{
		"token": token,
	})
	s.Secret.TTL = result.Lease

	return s, nil
}

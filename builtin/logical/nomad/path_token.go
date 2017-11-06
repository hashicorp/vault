package nomad

import (
	"fmt"
	"time"

	"github.com/hashicorp/nomad/api"
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

	entry, err := req.Storage.Get("role/" + name)
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

	// Get the nomad client
	c, err := b.client(req.Storage)
	if err != nil {
		return nil, err
	}

	// Generate a name for the token
	tokenName := fmt.Sprintf("Vault %s %s %d", name, req.DisplayName, time.Now().UnixNano())

	// Create it
	token, _, err := c.ACLTokens().Create(&api.ACLToken{
		Name:     tokenName,
		Type:     result.TokenType,
		Policies: result.Policy,
		Global:   result.Global,
	}, nil)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Use the helper to create the secret
	s := b.Secret(SecretTokenType).Response(map[string]interface{}{
		"secret_id":   token.SecretID,
		"accessor_id": token.AccessorID,
	}, map[string]interface{}{
		"accessor_id": token.AccessorID,
	})

	return s, nil
}

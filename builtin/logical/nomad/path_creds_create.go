package nomad

import (
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathCredsCreate(b *backend) *framework.Path {
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

	role, err := b.Role(req.Storage, name)
	if err != nil {
		return nil, errwrap.Wrapf("error retrieving role: {{err}}", err)
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("role %q not found", name)), nil
	}

	// Determine if we have a lease configuration
	leaseConfig, err := b.LeaseConfig(req.Storage)
	if err != nil {
		return nil, err
	}
	if leaseConfig == nil {
		leaseConfig = &configLease{}
	}

	// Get the nomad client
	c, err := b.client(req.Storage)
	if err != nil {
		return nil, err
	}

	// Generate a name for the token
	tokenName := fmt.Sprintf("vault-%s-%s-%d", name, req.DisplayName, time.Now().UnixNano())

	// Create it
	token, _, err := c.ACLTokens().Create(&api.ACLToken{
		Name:     tokenName,
		Type:     role.TokenType,
		Policies: role.Policies,
		Global:   role.Global,
	}, nil)
	if err != nil {
		return nil, err
	}

	// Use the helper to create the secret
	resp := b.Secret(SecretTokenType).Response(map[string]interface{}{
		"secret_id":   token.SecretID,
		"accessor_id": token.AccessorID,
	}, map[string]interface{}{
		"accessor_id": token.AccessorID,
	})
	resp.Secret.TTL = leaseConfig.TTL

	return resp, nil
}

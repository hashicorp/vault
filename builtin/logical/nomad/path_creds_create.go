package nomad

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// maxTokenNameLength is the maximum length for the name of a Nomad access
// token
const maxTokenNameLength = 256

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

func (b *backend) pathTokenRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	conf, _ := b.readConfigAccess(ctx, req.Storage)
	// establish a default
	tokenNameLength := maxTokenNameLength
	if conf != nil && conf.MaxTokenNameLength > 0 {
		tokenNameLength = conf.MaxTokenNameLength
	}

	role, err := b.Role(ctx, req.Storage, name)
	if err != nil {
		return nil, errwrap.Wrapf("error retrieving role: {{err}}", err)
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("role %q not found", name)), nil
	}

	// Determine if we have a lease configuration
	leaseConfig, err := b.LeaseConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if leaseConfig == nil {
		leaseConfig = &configLease{}
	}

	// Get the nomad client
	c, err := b.client(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	// Generate a name for the token
	tokenName := fmt.Sprintf("vault-%s-%s-%d", name, req.DisplayName, time.Now().UnixNano())

	// Note: if the given role name is sufficiently long, the UnixNano() portion
	// of the pseudo randomized token name is the part that gets trimmed off,
	// weakening it's randomness.
	if len(tokenName) > tokenNameLength {
		tokenName = tokenName[:tokenNameLength]
	}

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
	resp.Secret.MaxTTL = leaseConfig.MaxTTL

	return resp, nil
}

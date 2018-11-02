package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	SecretTokenType = "token"
)

func secretToken(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTokenType,
		Fields: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Request token",
			},
		},

		Renew:  b.secretTokenRenew,
		Revoke: b.secretTokenRevoke,
	}
}

func (b *backend) secretTokenRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp := &logical.Response{Secret: req.Secret}
	roleRaw, ok := req.Secret.InternalData["role"]
	if !ok {
		return resp, nil
	}

	role, ok := roleRaw.(string)
	if !ok {
		return resp, nil
	}

	entry, err := req.Storage.Get(ctx, "policy/"+role)
	if err != nil {
		return nil, errwrap.Wrapf("error retrieving role: {{err}}", err)
	}
	if entry == nil {
		return logical.ErrorResponse(fmt.Sprintf("issuing role %q not found", role)), nil
	}

	var result roleConfig
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	resp.Secret.TTL = result.TTL
	resp.Secret.MaxTTL = result.MaxTTL
	return resp, nil
}

func (b *backend) secretTokenRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, userErr, intErr := b.client(ctx, req.Storage)
	if intErr != nil {
		return nil, intErr
	}
	if userErr != nil {
		// Returning logical.ErrorResponse from revocation function is risky
		return nil, userErr
	}

	tokenRaw, ok := req.Secret.InternalData["token"]
	if !ok {
		// We return nil here because this is a pre-0.5.3 problem and there is
		// nothing we can do about it. We already can't revoke the lease
		// properly if it has been renewed and this is documented pre-0.5.3
		// behavior with a security bulletin about it.
		return nil, nil
	}

	var version string
	versionRaw, ok := req.Secret.InternalData["version"]
	if ok {
		version = versionRaw.(string)
	}

	switch version {
	case "":
		// Pre 1.4 tokens
		_, err := c.ACL().Destroy(tokenRaw.(string), nil)
		if err != nil {
			return nil, err
		}
	case tokenPolicyType:
		_, err := c.ACL().TokenDelete(tokenRaw.(string), nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Invalid version string in data: %s", version)
	}

	return nil, nil
}

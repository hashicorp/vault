// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	SecretTokenType = "token"
)

func secretToken(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTokenType,
		Fields: map[string]*framework.FieldSchema{
			"token": {
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
		return nil, fmt.Errorf("error retrieving role: %w", err)
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

	b.TryRecordObservationWithRequest(ctx, req, ObservationTypeConsulCredentialRenew, map[string]interface{}{
		"role_name":  role,
		"ttl":        result.TTL.String(),
		"max_ttl":    result.MaxTTL.String(),
		"token_type": result.TokenType,
		"local":      result.Local,
	})

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

	// Extract Consul Namespace and Partition info from secret
	var revokeWriteOptions *api.WriteOptions
	var namespace, partition string

	namespaceRaw, ok := req.Data["consul_namespace"]
	if ok {
		namespace = namespaceRaw.(string)
	}
	partitionRaw, ok := req.Data["partition"]
	if ok {
		partition = partitionRaw.(string)
	}

	revokeWriteOptions = &api.WriteOptions{
		Namespace: namespace,
		Partition: partition,
	}

	switch version {
	case "":
		// Pre 1.4 tokens
		token := tokenRaw.(string)
		_, err := c.ACL().Destroy(token, nil)
		if err != nil {
			return nil, err
		}
		b.TryRecordObservationWithRequest(ctx, req, ObservationTypeConsulCredentialRevoke, map[string]interface{}{
			"role_name": req.Secret.InternalData["role"],
		})
	case tokenPolicyType:
		token := tokenRaw.(string)
		_, err := c.ACL().TokenDelete(token, revokeWriteOptions)
		if err != nil {
			return nil, err
		}
		b.TryRecordObservationWithRequest(ctx, req, ObservationTypeConsulCredentialRevoke, map[string]interface{}{
			"role_name": req.Secret.InternalData["role"],
		})
	default:
		return nil, fmt.Errorf("Invalid version string in data: %s", version)
	}

	return nil, nil
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package rabbitmq

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// SecretCredsType is the key for this backend's secrets.
const SecretCredsType = "creds"

func secretCreds(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretCredsType,
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "RabbitMQ username",
			},
			"password": {
				Type:        framework.TypeString,
				Description: "Password for the RabbitMQ username",
			},
		},
		Renew:  b.secretCredsRenew,
		Revoke: b.secretCredsRevoke,
	}
}

// Renew the previously issued secret
func (b *backend) secretCredsRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the lease information
	lease, err := b.Lease(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}

	resp := &logical.Response{Secret: req.Secret}
	resp.Secret.TTL = lease.TTL
	resp.Secret.MaxTTL = lease.MaxTTL
	return resp, nil
}

// Revoke the previously issued secret
func (b *backend) secretCredsRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the username from the internal data
	usernameRaw, ok := req.Secret.InternalData["username"]
	if !ok {
		return nil, fmt.Errorf("secret is missing username internal data")
	}
	username := usernameRaw.(string)

	// Get our connection
	client, err := b.Client(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if _, err = client.DeleteUser(username); err != nil {
		return nil, fmt.Errorf("could not delete user: %w", err)
	}

	return nil, nil
}

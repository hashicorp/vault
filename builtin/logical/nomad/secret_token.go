// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package nomad

import (
	"context"
	"errors"
	"fmt"

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
	lease, err := b.LeaseConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if lease == nil {
		lease = &configLease{}
	}
	resp := &logical.Response{Secret: req.Secret}
	resp.Secret.TTL = lease.TTL
	resp.Secret.MaxTTL = lease.MaxTTL
	accessorID := ""
	if req.Secret.InternalData != nil {
		accessorIDRaw, ok := req.Secret.InternalData["accessor_id"]
		if ok {
			accessorID, _ = accessorIDRaw.(string)
		}
	}
	b.TryRecordObservationWithRequest(ctx, req, ObservationTypeNomadCredentialRenew, map[string]interface{}{
		"ttl":         lease.TTL.String(),
		"max_ttl":     lease.MaxTTL.String(),
		"accessor_id": accessorID,
	})
	return resp, nil
}

func (b *backend) secretTokenRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, err := b.client(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, fmt.Errorf("error getting Nomad client")
	}

	accessorIDRaw, ok := req.Secret.InternalData["accessor_id"]
	if !ok {
		return nil, fmt.Errorf("accessor_id is missing on the lease")
	}
	accessorID, ok := accessorIDRaw.(string)
	if !ok {
		return nil, errors.New("unable to convert accessor_id")
	}
	_, err = c.ACLTokens().Delete(accessorID, nil)
	if err != nil {
		return nil, err
	}

	b.TryRecordObservationWithRequest(ctx, req, ObservationTypeNomadCredentialRevoke, map[string]interface{}{
		"accessor_id": accessorID,
	})

	return nil, nil
}

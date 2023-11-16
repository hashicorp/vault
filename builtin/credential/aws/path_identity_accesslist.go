// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const identityAccessListStorage = "whitelist/identity/"

func (b *backend) pathIdentityAccessList() *framework.Path {
	return &framework.Path{
		Pattern: "identity-accesslist/" + framework.GenericNameRegex("instance_id"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationSuffix: "identity-access-list",
		},

		Fields: map[string]*framework.FieldSchema{
			"instance_id": {
				Type: framework.TypeString,
				Description: `EC2 instance ID. A successful login operation from an EC2 instance
gets cached in this accesslist, keyed off of instance ID.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathIdentityAccesslistRead,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathIdentityAccesslistDelete,
			},
		},

		HelpSynopsis:    pathIdentityAccessListSyn,
		HelpDescription: pathIdentityAccessListDesc,
	}
}

func (b *backend) pathListIdentityAccessList() *framework.Path {
	return &framework.Path{
		Pattern: "identity-accesslist/?",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAWS,
			OperationSuffix: "identity-access-list",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathAccessListIdentitiesList,
			},
		},

		HelpSynopsis:    pathListIdentityAccessListHelpSyn,
		HelpDescription: pathListIdentityAccessListHelpDesc,
	}
}

// pathAccessListIdentitiesList is used to list all the instance IDs that are present
// in the identity access list. This will list both valid and expired entries.
func (b *backend) pathAccessListIdentitiesList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	identities, err := req.Storage.List(ctx, identityAccessListStorage)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(identities), nil
}

// Fetch an item from the access list given an instance ID.
func accessListIdentityEntry(ctx context.Context, s logical.Storage, instanceID string) (*accessListIdentity, error) {
	entry, err := s.Get(ctx, identityAccessListStorage+instanceID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result accessListIdentity
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Stores an instance ID and the information required to validate further login/renewal attempts from
// the same instance ID.
func setAccessListIdentityEntry(ctx context.Context, s logical.Storage, instanceID string, identity *accessListIdentity) error {
	entry, err := logical.StorageEntryJSON(identityAccessListStorage+instanceID, identity)
	if err != nil {
		return err
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}
	return nil
}

// pathIdentityAccesslistDelete is used to delete an entry from the identity access list given an instance ID.
func (b *backend) pathIdentityAccesslistDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	instanceID := data.Get("instance_id").(string)
	if instanceID == "" {
		return logical.ErrorResponse("missing instance_id"), nil
	}

	return nil, req.Storage.Delete(ctx, identityAccessListStorage+instanceID)
}

// pathIdentityAccesslistRead is used to view an entry in the identity access list given an instance ID.
func (b *backend) pathIdentityAccesslistRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	instanceID := data.Get("instance_id").(string)
	if instanceID == "" {
		return logical.ErrorResponse("missing instance_id"), nil
	}

	entry, err := accessListIdentityEntry(ctx, req.Storage, instanceID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"role":                      entry.Role,
			"client_nonce":              entry.ClientNonce,
			"creation_time":             entry.CreationTime.Format(time.RFC3339Nano),
			"disallow_reauthentication": entry.DisallowReauthentication,
			"pending_time":              entry.PendingTime,
			"expiration_time":           entry.ExpirationTime.Format(time.RFC3339Nano),
			"last_updated_time":         entry.LastUpdatedTime.Format(time.RFC3339Nano),
		},
	}, nil
}

// Struct to represent each item in the identity access list.
type accessListIdentity struct {
	Role                     string    `json:"role"`
	ClientNonce              string    `json:"client_nonce"`
	CreationTime             time.Time `json:"creation_time"`
	DisallowReauthentication bool      `json:"disallow_reauthentication"`
	PendingTime              string    `json:"pending_time"`
	ExpirationTime           time.Time `json:"expiration_time"`
	LastUpdatedTime          time.Time `json:"last_updated_time"`
}

const pathIdentityAccessListSyn = `
Read or delete entries in the identity access list.
`

const pathIdentityAccessListDesc = `
Each login from an EC2 instance creates/updates an entry in the identity access list.

Entries in this list can be viewed or deleted using this endpoint.

By default, a cron task will periodically look for expired entries in the access list
and deletes them. The duration to periodically run this, is one hour by default.
However, this can be configured using the 'config/tidy/identities' endpoint. This tidy
action can be triggered via the API as well, using the 'tidy/identities' endpoint.
`

const pathListIdentityAccessListHelpSyn = `
Lists the items present in the identity access list.
`

const pathListIdentityAccessListHelpDesc = `
The entries in the identity access list is keyed off of the EC2 instance IDs.
This endpoint lists all the entries present in the identity access list, both
expired and un-expired entries. Use 'tidy/identities' endpoint to clean-up
the access list of identities.
`

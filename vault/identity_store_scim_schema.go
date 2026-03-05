// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"

	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/identity"
)

// SCIM client storage prefix
const scimClientStoragePrefix = "scim/client/"

func scimClientSchema(_ bool) *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: scimClientsTable,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ClientID",
				},
			},
			"client_id": {
				Name: "client_id",
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "NamespaceID",
						},
						&memdb.StringFieldIndex{
							Field: "ClientID",
						},
					},
				},
			},
			"client_name": {
				Name:   "client_name",
				Unique: true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "NamespaceID",
						},
						&memdb.StringFieldIndex{
							Field: "ClientName",
						},
					},
				},
			},
			"namespace_id": {
				Name: "namespace_id",
				Indexer: &memdb.StringFieldIndex{
					Field: "NamespaceID",
				},
			},
			"access_grant_principal": {
				Name:   "access_grant_principal",
				Unique: true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "NamespaceID",
						},
						&memdb.StringFieldIndex{
							Field: "AccessGrantPrincipal",
						},
					},
				},
			},
		},
	}
}

type scimClientRequest struct{}

func addSCIMClientIDToContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, scimClientRequest{}, id)
}

func scimClientIDFromContext(ctx context.Context) string {
	val := ctx.Value(scimClientRequest{})
	if val == nil {
		return ""
	}
	return val.(string)
}

func (i *IdentityStore) scimResourceCheck(ctx context.Context, resource scimManaged, originalSCIMID string, isCreate bool) error {
	reqSCIMClientID := scimClientIDFromContext(ctx)
	resourceSCIMClientID := resource.SCIMClientID()

	switch isCreate {
	case true:
		// The request must have come via a SCIM API in order to set
		// the SCIM client ID
		if reqSCIMClientID == "" && resourceSCIMClientID != "" {
			return errors.New("cannot set scim_client_id")
		}
		if reqSCIMClientID != "" && resourceSCIMClientID == "" {
			// this shouldn't ever happen
			return errors.New("cannot create a resource via SCIM without a SCIM client ID")
		}
		if reqSCIMClientID != resourceSCIMClientID {
			// this also shouldn't ever happen
			return errors.New("cannot create resource via SCIM with a different SCIM client ID")
		}

	default:
		// if the resource is being updated, the SCIM client ID
		// cannot be modified
		if originalSCIMID != resourceSCIMClientID {
			return errors.New("cannot update scim_client_id")
		}
		// if the resource is being updated, this must be via SCIM
		if originalSCIMID != reqSCIMClientID {
			return errors.New("SCIM-managed resources must be modified through SCIM")
		}
	}
	return nil
}

type scimManaged interface {
	SCIMClientID() string
}

var (
	_ scimManaged = (*identity.Entity)(nil)
	_ scimManaged = (*identity.Group)(nil)
	_ scimManaged = (*identity.Alias)(nil)
)

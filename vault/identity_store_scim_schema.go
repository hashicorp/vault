// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import "github.com/hashicorp/go-memdb"

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

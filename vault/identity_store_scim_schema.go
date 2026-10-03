// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import "github.com/hashicorp/go-memdb"

const (
	scimConfigTable = "scim_config"
)

func scimConfigSchema(_ bool) *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: scimConfigTable,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ID",
				},
			},
			"client_id": {
				Name:   "client_id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ClientID",
				},
			},
			"client_role": {
				Name: "client_role",
				Indexer: &memdb.StringFieldIndex{
					Field: "ClientRole",
				},
			},
			"access_grant_method": {
				Name: "access_grant_method",
				Indexer: &memdb.StringFieldIndex{
					Field: "AccessGrantMethod",
				},
			},
			"access_grant_principal": {
				Name:   "access_grant_principal",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "AccessGrantPrincipal",
				},
			},
			"alias_mount_accessor": {
				Name:   "alias_mount_accessor",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "AliasMountAccessor",
				},
				AllowMissing: true,
			},
		},
	}
}

// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	memdb "github.com/hashicorp/go-memdb"
)

const (
	tpmgroupsTable = "tpmgroups"
)

func tpmgroupsTableSchema(lowerCaseName bool) *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: tpmgroupsTable,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ID",
				},
			},
			"name": {
				Name:   "name",
				Unique: true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "NamespaceID",
						},
						&memdb.StringFieldIndex{
							Field:     "Name",
							Lowercase: lowerCaseName,
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
			"member_tpm_ids": {
				Name:         "member_tpm_ids",
				AllowMissing: true,
				Indexer: &memdb.StringSliceFieldIndex{
					Field: "MemberTpmIDs",
				},
			},
			"bucket_key": {
				Name: "bucket_key",
				Indexer: &memdb.StringFieldIndex{
					Field: "BucketKey",
				},
			},
		},
	}
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"fmt"

	memdb "github.com/hashicorp/go-memdb"
)

const (
	entitiesTable      = "entities"
	entityAliasesTable = "entity_aliases"
	groupsTable        = "groups"
	groupAliasesTable  = "group_aliases"
	oidcClientsTable   = "oidc_clients"
)

func identityStoreSchema(lowerCaseName bool) *memdb.DBSchema {
	iStoreSchema := &memdb.DBSchema{
		Tables: make(map[string]*memdb.TableSchema),
	}

	schemas := []func(bool) *memdb.TableSchema{
		entitiesTableSchema,
		aliasesTableSchema,
		groupsTableSchema,
		groupAliasesTableSchema,
		oidcClientsTableSchema,
	}

	for _, schemaFunc := range schemas {
		schema := schemaFunc(lowerCaseName)
		if _, ok := iStoreSchema.Tables[schema.Name]; ok {
			panic(fmt.Sprintf("duplicate table name: %s", schema.Name))
		}
		iStoreSchema.Tables[schema.Name] = schema
	}

	return iStoreSchema
}

func aliasesTableSchema(lowerCaseName bool) *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: entityAliasesTable,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ID",
				},
			},
			"factors": {
				Name:   "factors",
				Unique: true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "MountAccessor",
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
			"local_bucket_key": {
				Name:         "local_bucket_key",
				AllowMissing: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "LocalBucketKey",
				},
			},
		},
	}
}

func entitiesTableSchema(lowerCaseName bool) *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: entitiesTable,
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
			"merged_entity_ids": {
				Name:         "merged_entity_ids",
				Unique:       true,
				AllowMissing: true,
				Indexer: &memdb.StringSliceFieldIndex{
					Field: "MergedEntityIDs",
				},
			},
			"bucket_key": {
				Name: "bucket_key",
				Indexer: &memdb.StringFieldIndex{
					Field: "BucketKey",
				},
			},
			"namespace_id": {
				Name: "namespace_id",
				Indexer: &memdb.StringFieldIndex{
					Field: "NamespaceID",
				},
			},
		},
	}
}

func groupsTableSchema(lowerCaseName bool) *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: groupsTable,
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
			"member_entity_ids": {
				Name:         "member_entity_ids",
				AllowMissing: true,
				Indexer: &memdb.StringSliceFieldIndex{
					Field: "MemberEntityIDs",
				},
			},
			"parent_group_ids": {
				Name:         "parent_group_ids",
				AllowMissing: true,
				Indexer: &memdb.StringSliceFieldIndex{
					Field: "ParentGroupIDs",
				},
			},
			"bucket_key": {
				Name: "bucket_key",
				Indexer: &memdb.StringFieldIndex{
					Field: "BucketKey",
				},
			},
			"namespace_id": {
				Name: "namespace_id",
				Indexer: &memdb.StringFieldIndex{
					Field: "NamespaceID",
				},
			},
		},
	}
}

func groupAliasesTableSchema(lowerCaseName bool) *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: groupAliasesTable,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ID",
				},
			},
			"factors": {
				Name:   "factors",
				Unique: true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "MountAccessor",
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
		},
	}
}

func oidcClientsTableSchema(_ bool) *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: oidcClientsTable,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ClientID",
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
							Field: "Name",
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
		},
	}
}

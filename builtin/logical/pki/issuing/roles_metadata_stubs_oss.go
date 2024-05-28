// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package issuing

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

import "github.com/hashicorp/vault/sdk/framework"

func AddNoStoreMetadata(roleData map[string]interface{}, r *RoleEntry) {
	return
}

func WithNoStoreMetadata(noStoreMetadata bool) RoleModifier {
	return func(r *RoleEntry) {
		r.NoStoreMetadata = true
	}
}

const MetadataPermitted = false

func AddNoStoreMetadataRoleField(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	return fields
}

func GetNoStoreMetadata(data *framework.FieldData) bool {
	return true
}

func NoStoreMetadataValue(value bool) bool {
	return true
}

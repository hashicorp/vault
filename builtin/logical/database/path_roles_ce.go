// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package database

import "github.com/hashicorp/vault/sdk/framework"

// AddStaticFieldsEnt is a no-op for comminuty edition
func AddStaticFieldsEnt(fields map[string]*framework.FieldSchema) {
	// no-op
}

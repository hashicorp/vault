// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !enterprise

package database

import "github.com/hashicorp/vault/sdk/framework"

// AddConnectionFieldsEnt is a no-op for community edition
func AddConnectionFieldsEnt(fields map[string]*framework.FieldSchema) {
	// no-op
}

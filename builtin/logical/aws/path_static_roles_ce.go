// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package aws

import "github.com/hashicorp/vault/sdk/framework"

// AddStaticAssumeRoleFieldsEnt is a no-op for community edition
func AddStaticAssumeRoleFieldsEnt(fields map[string]*framework.FieldSchema) {
	// no-op
}

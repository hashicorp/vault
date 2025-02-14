// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package aws

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
)

// AddStaticAssumeRoleFieldsEnt is a no-op for community edition
func AddStaticAssumeRoleFieldsEnt(fields map[string]*framework.FieldSchema) {
	// no-op
}

func validateAssumeRoleFields(data *framework.FieldData, config *staticRoleEntry) error {
	_, hasAssumeRoleARN := data.GetOk(paramAssumeRoleARN)
	_, hasRoleSessionName := data.GetOk(paramRoleSessionName)
	_, hasExternalID := data.GetOk(paramExternalID)

	if hasAssumeRoleARN || hasRoleSessionName || hasExternalID {
		return fmt.Errorf("cross-account static roles are only supported in Vault Enterprise")
	}

	return nil
}

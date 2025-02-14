// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/vault/sdk/logical"
)

// assumeRoleStatic assumes an AWS role for cross-account static role management.
// It uses the role ARN and session name provided in the staticRoleEntry configuration
// to generate credentials for the assumed role.
func (b *backend) assumeRoleStatic(ctx context.Context, s logical.Storage, entry *staticRoleEntry) (*aws.Config, error) {
	return nil, fmt.Errorf("assumeRoleStatic is not supported in this version of Vault")
}

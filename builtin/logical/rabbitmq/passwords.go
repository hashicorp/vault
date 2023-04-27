// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rabbitmq

import (
	"context"

	"github.com/hashicorp/go-secure-stdlib/base62"
)

func (b *backend) generatePassword(ctx context.Context, policyName string) (password string, err error) {
	if policyName != "" {
		return b.System().GeneratePasswordFromPolicy(ctx, policyName)
	}
	return base62.Random(36)
}

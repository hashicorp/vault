package rabbitmq

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/base62"
)

func (b *backend) generatePassword(ctx context.Context, policyName string) (password string, err error) {
	if policyName != "" {
		return b.System().GeneratePasswordFromPolicy(ctx, policyName)
	}
	return base62.Random(36)
}

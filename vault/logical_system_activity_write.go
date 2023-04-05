//go:build !testonly

package vault

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *SystemBackend) activityWritePath() *framework.Path {
	return &framework.Path{
		Pattern:    "internal/counters/activity/write$",
		Operations: map[logical.Operation]framework.OperationHandler{},
	}
}

package azuresecrets

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/base62"
)

const (
	passwordLength = 36
)

type passwordGenerator interface {
	GeneratePasswordFromPolicy(ctx context.Context, policyName string) (password string, err error)
}

type passwords struct {
	policyGenerator passwordGenerator
	policyName      string
}

func (p passwords) generate(ctx context.Context) (password string, err error) {
	if p.policyName == "" {
		return base62.Random(passwordLength)
	}
	if p.policyGenerator == nil {
		return "", fmt.Errorf("policy set, but no policy generator specified")
	}
	return p.policyGenerator.GeneratePasswordFromPolicy(ctx, p.policyName)
}

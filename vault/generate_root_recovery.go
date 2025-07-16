// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"go.uber.org/atomic"
)

// GenerateRecoveryTokenStrategy is the strategy used to generate a
// recovery token
func GenerateRecoveryTokenStrategy(token *atomic.String) GenerateRootStrategy {
	return &generateRecoveryToken{token: token}
}

// generateRecoveryToken implements the GenerateRootStrategy and is in
// charge of creating recovery tokens.
type generateRecoveryToken struct {
	token *atomic.String
}

func (g *generateRecoveryToken) authenticate(ctx context.Context, c *Core, combinedKey []byte) error {
	key, err := c.unsealKeyToRootKeyPostUnseal(ctx, combinedKey)
	if err != nil {
		return fmt.Errorf("unable to authenticate: %w", err)
	}

	// Use the retrieved root key to unseal the barrier
	if err := c.barrier.Unseal(ctx, key); err != nil {
		return fmt.Errorf("recovery token generation failed, cannot unseal barrier: %w", err)
	}

	for _, v := range c.postRecoveryUnsealFuncs {
		if err := v(); err != nil {
			return fmt.Errorf("failed to run post unseal func: %w", err)
		}
	}
	return nil
}

func (g *generateRecoveryToken) generate(ctx context.Context, c *Core) (string, func(), error) {
	var id string
	var err error
	id, err = base62.Random(TokenLength)
	if err != nil {
		return "", nil, err
	}
	var token string
	if c.DisableSSCTokens() {
		token = consts.LegacyRecoveryTokenPrefix + id
	} else {
		token = consts.RecoveryTokenPrefix + id
	}
	g.token.Store(token)

	return token, func() { g.token.Store("") }, nil
}

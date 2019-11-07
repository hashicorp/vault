package vault

import (
	"context"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/base62"
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
	key, err := c.unsealKeyToMasterKey(ctx, combinedKey)
	if err != nil {
		return errwrap.Wrapf("unable to authenticate: {{err}}", err)
	}

	// Use the retrieved master key to unseal the barrier
	if err := c.barrier.Unseal(ctx, key); err != nil {
		return errwrap.Wrapf("recovery operation token generation failed, cannot unseal barrier: {{err}}", err)
	}

	for _, v := range c.postRecoveryUnsealFuncs {
		if err := v(); err != nil {
			return errwrap.Wrapf("failed to run post unseal func: {{err}}", err)
		}
	}
	return nil
}

func (g *generateRecoveryToken) generate(ctx context.Context, c *Core) (string, func(), error) {
	id, err := base62.Random(TokenLength)
	if err != nil {
		return "", nil, err
	}
	token := "r." + id
	g.token.Store(token)

	return token, func() { g.token.Store("") }, nil
}

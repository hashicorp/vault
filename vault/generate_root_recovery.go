package vault

import (
	"context"
	"errors"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/shamir"

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

func (g *generateRecoveryToken) authenticate(ctx context.Context, c *Core, key []byte) error {
	// If recovery keys are supported, verify the recovery key and fetch the stored keys
	if c.seal.RecoveryKeySupported() {
		if err := c.seal.VerifyRecoveryKey(ctx, key); err != nil {
			return errwrap.Wrapf("recovery key verification failed: {{err}}", err)
		}

		if !c.seal.StoredKeysSupported() {
			return errors.New("recovery key verified but stored keys unsupported")
		}

		masterKeyShares, err := c.seal.GetStoredKeys(ctx)
		if err != nil {
			return errwrap.Wrapf("unable to retrieve stored keys: {{err}}", err)
		}

		switch len(masterKeyShares) {
		case 0:
			return errors.New("seal returned no master key shares")
		case 1:
			key = masterKeyShares[0]
		default:
			key, err = shamir.Combine(masterKeyShares)
			if err != nil {
				return errwrap.Wrapf("failed to compute master key: {{err}}", err)
			}
		}
	}

	if err := c.barrier.Unseal(ctx, key); err != nil {
		return errwrap.Wrapf("recovery operation token verification failed: {{err}}", err)
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

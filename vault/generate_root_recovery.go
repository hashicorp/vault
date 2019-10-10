package vault

import (
	"context"

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

func (g *generateRecoveryToken) generate(ctx context.Context, c *Core) (string, func(), error) {
	id, err := base62.Random(TokenLength)
	if err != nil {
		return "", nil, err
	}
	token := "r." + id
	g.token.Store(token)

	return token, func() { g.token.Store("") }, nil
}

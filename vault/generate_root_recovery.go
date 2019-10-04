package vault

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/errwrap"
)

const coreRecoveryTokenPath = "core/recovery-token"

var (
	// GenerateRecoveryTokenStrategy is the strategy used to generate a
	// recovery token
	GenerateRecoveryTokenStrategy GenerateRootStrategy = generateRecoveryToken{}
)

// generateRecoveryToken implements the GenerateRootStrategy and is in
// charge of creating recovery tokens.
type generateRecoveryToken struct{}

func (g generateRecoveryToken) generate(ctx context.Context, c *Core) (string, func(), error) {
	tokenUUID, err := c.CreateDROperationToken(ctx)
	if err != nil {
		return "", nil, err
	}

	cleanupFunc := func() {
		c.DeleteRecoveryToken(ctx)
	}

	return tokenUUID, cleanupFunc, nil
}

func (c *Core) DeleteRecoveryToken(ctx context.Context) error {
	return c.barrier.Delete(ctx, coreDROperationTokenPath)
}

func (c *Core) CreateDROperationToken(ctx context.Context) (string, error) {
	id, err := base62.Random(TokenLength)
	if err != nil {
		return "", err
	}
	id = "s." + id

	root := &RecoveryToken{
		ID:           id,
		CreationTime: time.Now(),
	}

	buf, err := json.Marshal(root)
	if err != nil {
		return "", errwrap.Wrapf("failed to encode DR Root Token: {{err}}", err)
	}

	if err := c.barrier.Put(ctx, &logical.StorageEntry{
		Key:   coreRecoveryTokenPath,
		Value: buf,
	}); err != nil {
		return "", err
	}

	return id, nil
}

func (c *Core) GetRecoveryToken(ctx context.Context) (string, error) {
	raw, err := c.barrier.Get(ctx, coreDROperationTokenPath)
	if err != nil {
		return "", err
	}
	if raw == nil {
		return "", nil
	}

	drRootToken := &RecoveryToken{}
	err = json.Unmarshal(raw.Value, drRootToken)
	if err != nil {
		return "", err
	}

	return drRootToken.ID, nil
}

type RecoveryToken struct {
	// ID of this entry, generally a random UUID
	ID string `json:"id" mapstructure:"id" structs:"id"`
	// Time of token creation
	CreationTime time.Time `json:"creation_time" mapstructure:"creation_time" structs:"creation_time"`
}

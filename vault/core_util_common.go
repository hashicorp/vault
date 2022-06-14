package vault

import (
	"context"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *Core) loadHeaderHMACKey(ctx context.Context) error {
	ent, err := c.barrier.Get(ctx, indexHeaderHMACKeyPath)
	if err != nil {
		return err
	}

	if ent != nil {
		c.IndexHeaderHMACKey.Store(ent.Value)
	}
	return nil
}

func (c *Core) headerHMACKey() []byte {
	key := c.IndexHeaderHMACKey.Load()
	if key == nil {
		return nil
	}
	return key.([]byte)
}

func (c *Core) setupHeaderHMACKey(ctx context.Context, isPerfStandby bool) error {
	if c.IsPerfSecondary() || c.IsDRSecondary() || isPerfStandby {
		return c.loadHeaderHMACKey(ctx)
	}
	ent, err := c.barrier.Get(ctx, indexHeaderHMACKeyPath)
	if err != nil {
		return err
	}

	if ent != nil {
		c.IndexHeaderHMACKey.Store(ent.Value)
		return nil
	}

	key, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	err = c.barrier.Put(ctx, &logical.StorageEntry{
		Key:   indexHeaderHMACKeyPath,
		Value: []byte(key),
	})
	if err != nil {
		return err
	}
	c.IndexHeaderHMACKey.Store([]byte(key))
	return nil
}

// +build !enterprise

package vault

import "context"

func (c *Core) postSealMigration(ctx context.Context) error { return nil }

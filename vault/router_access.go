package vault

import "context"

// RouterAccess provides access into some things necessary for testing
type RouterAccess struct {
	c *Core
}

func NewRouterAccess(c *Core) *RouterAccess {
	return &RouterAccess{c: c}
}

func (r *RouterAccess) StoragePrefixByAPIPath(ctx context.Context, path string) (string, bool) {
	return r.c.router.MatchingStoragePrefixByAPIPath(ctx, path)
}

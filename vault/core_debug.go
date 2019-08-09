package vault

import (
	"context"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

// SysDebugTokenCheck is a wrapper that simply checks for the token permission
// on the debug/pprof route. The call to the pprof handler should be done at the
// http layer.
func (c *Core) SysDebugTokenCheck(httpCtx context.Context, req *logical.Request) (retErr error) {
	ctx, cancel := context.WithCancel(namespace.RootContext(nil))
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-httpCtx.Done():
			cancel()
		}
	}()

	_, _, retErr = c.checkToken(ctx, req, false)

	return retErr
}

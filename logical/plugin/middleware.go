package plugin

import (
	"context"
	"time"

	"github.com/hashicorp/vault/logical"
	log "github.com/mgutz/logxi/v1"
)

// backendPluginClient implements logical.Backend and is the
// go-plugin client.
type backendTracingMiddleware struct {
	logger    log.Logger
	transport string
	typeStr   string

	next logical.Backend
}

// Validate the backendTracingMiddle object satisfies the backend interface
var _ logical.Backend = &backendTracingMiddleware{}

func (b *backendTracingMiddleware) HandleRequest(ctx context.Context, req *logical.Request) (resp *logical.Response, err error) {
	defer func(then time.Time) {
		b.logger.Trace("plugin.HandleRequest", "path", req.Path, "status", "finished", "type", b.typeStr, "transport", b.transport, "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("plugin.HandleRequest", "path", req.Path, "status", "started", "type", b.typeStr, "transport", b.transport)
	return b.next.HandleRequest(ctx, req)
}

func (b *backendTracingMiddleware) SpecialPaths() *logical.Paths {
	defer func(then time.Time) {
		b.logger.Trace("plugin.SpecialPaths", "status", "finished", "type", b.typeStr, "transport", b.transport, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("plugin.SpecialPaths", "status", "started", "type", b.typeStr, "transport", b.transport)
	return b.next.SpecialPaths()
}

func (b *backendTracingMiddleware) System() logical.SystemView {
	return b.next.System()
}

func (b *backendTracingMiddleware) Logger() log.Logger {
	return b.next.Logger()
}

func (b *backendTracingMiddleware) HandleExistenceCheck(ctx context.Context, req *logical.Request) (found bool, exists bool, err error) {
	defer func(then time.Time) {
		b.logger.Trace("plugin.HandleExistenceCheck", "path", req.Path, "status", "finished", "type", b.typeStr, "transport", b.transport, "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("plugin.HandleExistenceCheck", "path", req.Path, "status", "started", "type", b.typeStr, "transport", b.transport)
	return b.next.HandleExistenceCheck(ctx, req)
}

func (b *backendTracingMiddleware) Cleanup(ctx context.Context) {
	defer func(then time.Time) {
		b.logger.Trace("plugin.Cleanup", "status", "finished", "type", b.typeStr, "transport", b.transport, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("plugin.Cleanup", "status", "started", "type", b.typeStr, "transport", b.transport)
	b.next.Cleanup(ctx)
}

func (b *backendTracingMiddleware) InvalidateKey(ctx context.Context, key string) {
	defer func(then time.Time) {
		b.logger.Trace("plugin.InvalidateKey", "key", key, "status", "finished", "type", b.typeStr, "transport", b.transport, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("plugin.InvalidateKey", "key", key, "status", "started", "type", b.typeStr, "transport", b.transport)
	b.next.InvalidateKey(ctx, key)
}

func (b *backendTracingMiddleware) Setup(ctx context.Context, config *logical.BackendConfig) (err error) {
	defer func(then time.Time) {
		b.logger.Trace("plugin.Setup", "status", "finished", "type", b.typeStr, "transport", b.transport, "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("plugin.Setup", "status", "started", "type", b.typeStr, "transport", b.transport)
	return b.next.Setup(ctx, config)
}

func (b *backendTracingMiddleware) Type() logical.BackendType {
	defer func(then time.Time) {
		b.logger.Trace("plugin.Type", "status", "finished", "type", b.typeStr, "transport", b.transport, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("plugin.Type", "status", "started", "type", b.typeStr, "transport", b.transport)
	return b.next.Type()
}

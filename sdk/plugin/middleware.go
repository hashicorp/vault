package plugin

import (
	"context"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

// backendPluginClient implements logical.Backend and is the
// go-plugin client.
type BackendTracingMiddleware struct {
	BLogger log.Logger

	Next logical.Backend
}

// Validate the backendTracingMiddle object satisfies the backend interface
var _ logical.Backend = &BackendTracingMiddleware{}

func (b *BackendTracingMiddleware) Initialize(ctx context.Context, req *logical.InitializationRequest) (err error) {
	defer func(then time.Time) {
		b.BLogger.Trace("initialize", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.BLogger.Trace("initialize", "status", "started")
	return b.Next.Initialize(ctx, req)
}

func (b *BackendTracingMiddleware) HandleRequest(ctx context.Context, req *logical.Request) (resp *logical.Response, err error) {
	defer func(then time.Time) {
		b.BLogger.Trace("handle request", "path", req.Path, "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.BLogger.Trace("handle request", "path", req.Path, "status", "started")
	return b.Next.HandleRequest(ctx, req)
}

func (b *BackendTracingMiddleware) SpecialPaths() *logical.Paths {
	defer func(then time.Time) {
		b.BLogger.Trace("special paths", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.BLogger.Trace("special paths", "status", "started")
	return b.Next.SpecialPaths()
}

func (b *BackendTracingMiddleware) System() logical.SystemView {
	return b.Next.System()
}

func (b *BackendTracingMiddleware) Logger() log.Logger {
	return b.Next.Logger()
}

func (b *BackendTracingMiddleware) HandleExistenceCheck(ctx context.Context, req *logical.Request) (found bool, exists bool, err error) {
	defer func(then time.Time) {
		b.BLogger.Trace("handle existence check", "path", req.Path, "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.BLogger.Trace("handle existence check", "path", req.Path, "status", "started")
	return b.Next.HandleExistenceCheck(ctx, req)
}

func (b *BackendTracingMiddleware) Cleanup(ctx context.Context) {
	defer func(then time.Time) {
		b.BLogger.Trace("cleanup", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.BLogger.Trace("cleanup", "status", "started")
	b.Next.Cleanup(ctx)
}

func (b *BackendTracingMiddleware) InvalidateKey(ctx context.Context, key string) {
	defer func(then time.Time) {
		b.BLogger.Trace("invalidate key", "key", key, "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.BLogger.Trace("invalidate key", "key", key, "status", "started")
	b.Next.InvalidateKey(ctx, key)
}

func (b *BackendTracingMiddleware) Setup(ctx context.Context, config *logical.BackendConfig) (err error) {
	defer func(then time.Time) {
		b.BLogger.Trace("setup", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.BLogger.Trace("setup", "status", "started")
	return b.Next.Setup(ctx, config)
}

func (b *BackendTracingMiddleware) Type() logical.BackendType {
	defer func(then time.Time) {
		b.BLogger.Trace("type", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.BLogger.Trace("type", "status", "started")
	return b.Next.Type()
}

package plugin

import (
	"context"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/logical"
)

// backendPluginClient implements logical.Backend and is the
// go-plugin client.
type backendTracingMiddleware struct {
	logger log.Logger

	next logical.Backend
}

// Validate the backendTracingMiddle object satisfies the backend interface
var _ logical.Backend = &backendTracingMiddleware{}

func (b *backendTracingMiddleware) HandleRequest(ctx context.Context, req *logical.Request) (resp *logical.Response, err error) {
	defer func(then time.Time) {
		b.logger.Debug("HandleRequest", "path", req.Path, "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Debug("HandleRequest", "path", req.Path, "status", "started")
	return b.next.HandleRequest(ctx, req)
}

func (b *backendTracingMiddleware) SpecialPaths() *logical.Paths {
	defer func(then time.Time) {
		b.logger.Debug("SpecialPaths", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Debug("SpecialPaths", "status", "started")
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
		b.logger.Debug("HandleExistenceCheck", "path", req.Path, "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Debug("HandleExistenceCheck", "path", req.Path, "status", "started")
	return b.next.HandleExistenceCheck(ctx, req)
}

func (b *backendTracingMiddleware) Cleanup(ctx context.Context) {
	defer func(then time.Time) {
		b.logger.Debug("Cleanup", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Debug("Cleanup", "status", "started")
	b.next.Cleanup(ctx)
}

func (b *backendTracingMiddleware) InvalidateKey(ctx context.Context, key string) {
	defer func(then time.Time) {
		b.logger.Debug("InvalidateKey", "key", key, "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Debug("InvalidateKey", "key", key, "status", "started")
	b.next.InvalidateKey(ctx, key)
}

func (b *backendTracingMiddleware) Setup(ctx context.Context, config *logical.BackendConfig) (err error) {
	defer func(then time.Time) {
		b.logger.Debug("Setup", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Debug("Setup", "status", "started")
	return b.next.Setup(ctx, config)
}

func (b *backendTracingMiddleware) Type() logical.BackendType {
	defer func(then time.Time) {
		b.logger.Debug("Type", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Debug("Type", "status", "started")
	return b.next.Type()
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	logger log.Logger

	next logical.Backend
}

// Validate the backendTracingMiddle object satisfies the backend interface
var _ logical.Backend = &BackendTracingMiddleware{}

func (b *BackendTracingMiddleware) Initialize(ctx context.Context, req *logical.InitializationRequest) (err error) {
	defer func(then time.Time) {
		b.logger.Trace("initialize", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("initialize", "status", "started")
	return b.next.Initialize(ctx, req)
}

func (b *BackendTracingMiddleware) HandleRequest(ctx context.Context, req *logical.Request) (resp *logical.Response, err error) {
	defer func(then time.Time) {
		b.logger.Trace("handle request", "path", req.Path, "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("handle request", "path", req.Path, "status", "started")
	return b.next.HandleRequest(ctx, req)
}

func (b *BackendTracingMiddleware) SpecialPaths() *logical.Paths {
	defer func(then time.Time) {
		b.logger.Trace("special paths", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("special paths", "status", "started")
	return b.next.SpecialPaths()
}

func (b *BackendTracingMiddleware) System() logical.SystemView {
	return b.next.System()
}

func (b *BackendTracingMiddleware) Logger() log.Logger {
	return b.next.Logger()
}

func (b *BackendTracingMiddleware) HandleExistenceCheck(ctx context.Context, req *logical.Request) (found bool, exists bool, err error) {
	defer func(then time.Time) {
		b.logger.Trace("handle existence check", "path", req.Path, "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("handle existence check", "path", req.Path, "status", "started")
	return b.next.HandleExistenceCheck(ctx, req)
}

func (b *BackendTracingMiddleware) Cleanup(ctx context.Context) {
	defer func(then time.Time) {
		b.logger.Trace("cleanup", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("cleanup", "status", "started")
	b.next.Cleanup(ctx)
}

func (b *BackendTracingMiddleware) InvalidateKey(ctx context.Context, key string) {
	defer func(then time.Time) {
		b.logger.Trace("invalidate key", "key", key, "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("invalidate key", "key", key, "status", "started")
	b.next.InvalidateKey(ctx, key)
}

func (b *BackendTracingMiddleware) Setup(ctx context.Context, config *logical.BackendConfig) (err error) {
	defer func(then time.Time) {
		b.logger.Trace("setup", "status", "finished", "err", err, "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("setup", "status", "started")
	return b.next.Setup(ctx, config)
}

func (b *BackendTracingMiddleware) Type() logical.BackendType {
	defer func(then time.Time) {
		b.logger.Trace("type", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("type", "status", "started")
	return b.next.Type()
}

func (b *BackendTracingMiddleware) PluginVersion() logical.PluginVersion {
	defer func(then time.Time) {
		b.logger.Trace("version", "status", "finished", "took", time.Since(then))
	}(time.Now())

	b.logger.Trace("version", "status", "started")
	if versioner, ok := b.next.(logical.PluginVersioner); ok {
		return versioner.PluginVersion()
	}
	return logical.EmptyPluginVersion
}

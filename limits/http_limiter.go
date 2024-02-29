// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package limits

import (
	"context"
	"errors"
	"net/http"
)

//lint:ignore ST1005 Vault is the product name
var ErrCapacity = errors.New("Vault server temporarily overloaded")

const (
	WriteLimiter       = "write"
	SpecialPathLimiter = "special-path"
)

// HTTPLimiter is a convenience struct that we use to wrap some logical request
// context and prevent dependence on Core.
type HTTPLimiter struct {
	Method      string
	PathLimited bool
	LookupFunc  func(key string) *RequestLimiter
}

// CtxKeyDisableRequestLimiter holds the HTTP Listener's disable config if set.
type CtxKeyDisableRequestLimiter struct{}

func (c CtxKeyDisableRequestLimiter) String() string {
	return "disable_request_limiter"
}

// Acquire checks the HTTPLimiter metadata to determine if an HTTP request
// should be limited, or simply passed through as a no-op.
func (h *HTTPLimiter) Acquire(ctx context.Context) (*RequestListener, bool) {
	// If the limiter is disabled, return an empty wrapper so the limiter is a
	// no-op and indicate that the request can proceed.
	if disable := ctx.Value(CtxKeyDisableRequestLimiter{}); disable != nil && disable.(bool) {
		return &RequestListener{}, true
	}

	lim := &RequestLimiter{}
	if h.PathLimited {
		lim = h.LookupFunc(SpecialPathLimiter)
	} else {
		switch h.Method {
		case http.MethodGet, http.MethodHead, http.MethodTrace, http.MethodOptions:
			// We're only interested in the inverse, so do nothing here.
		default:
			lim = h.LookupFunc(WriteLimiter)
		}
	}
	return lim.Acquire(ctx)
}

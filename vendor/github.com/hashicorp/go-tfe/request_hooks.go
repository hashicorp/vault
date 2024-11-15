// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

import (
	"context"
	"fmt"
	"net/http"
)

// ContextWithResponseHeaderHook returns a context that will, if passed to
// [ClientRequest.Do] or to any of the wrapper methods that call it, arrange
// for the given callback to be called with the headers from the raw HTTP
// response.
//
// This is intended for allowing callers to respond to out-of-band metadata
// such as cache-control-related headers, rate limiting headers, etc. Hooks
// must not modify the given [http.Header] or otherwise attempt to change how
// the response is handled by [ClientRequest.Do].
//
// If the given context already has a response header hook then the returned
// context will call both the existing hook and the newly-provided one, with
// the newer being called first.
func ContextWithResponseHeaderHook(parentCtx context.Context, cb func(status int, header http.Header)) context.Context {
	// If the given context already has a notification callback then we'll
	// arrange to notify both the previous and the new one. This is not
	// a super efficient way to achieve that but we expect it to be rare
	// for there to be more than one or two hooks associated with a particular
	// request, so it's not warranted to optimize this further.
	existingI := parentCtx.Value(contextResponseHeaderHookKey)
	finalCb := cb
	if existingI != nil {
		existing, ok := existingI.(func(int, http.Header))
		// This explicit check-and-panic is redundant but required by our linter.
		if !ok {
			panic(fmt.Sprintf("context has response header hook of invalid type %T", existingI))
		}
		finalCb = func(status int, header http.Header) {
			cb(status, header)
			existing(status, header)
		}
	}
	return context.WithValue(parentCtx, contextResponseHeaderHookKey, finalCb)
}

func contextResponseHeaderHook(ctx context.Context) func(int, http.Header) {
	cbI := ctx.Value(contextResponseHeaderHookKey)
	if cbI == nil {
		// Stub callback that does absolutely nothing, then.
		return func(int, http.Header) {}
	}
	return cbI.(func(int, http.Header))
}

// contextResponseHeaderHookKey is the type of the internal key used to store
// the callback for [ContextWithResponseHeaderHook] inside a [context.Context]
// object.
type contextResponseHeaderHookKeyType struct{}

// contextResponseHeaderHookKey is the internal key used to store the callback
// for [ContextWithResponseHeaderHook] inside a [context.Context] object.
var contextResponseHeaderHookKey contextResponseHeaderHookKeyType

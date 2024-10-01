// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package priority

import (
	"context"
	"net/http"
	"strconv"

	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// VaultAOPForceRejectHeaderName is the name of an HTTP header that is used primarily
	// for testing (it's not documented publicly). If set to "true" in a request
	// that is subject to any form of Adaptive Overload Protection, the request
	// will be rejected as if there is an overload. This is useful for
	// deterministically testing the error handling plumbing as there are many
	// possible code paths that need to be tested.
	VaultAOPForceRejectHeaderName = "X-Vault-AOP-Force-Reject"
)

// Priorities are limited to 256 levels to keep the state space small making
// enforcement data structures much more efficient.
type AOPWritePriority uint8

const (
	// AlwaysDrop is intended for testing only and will cause the request to be
	// rejected with a 503 even if the server is not overloaded.
	AlwaysDrop AOPWritePriority = 0

	// StandardHTTP is the default AOPWritePriority for HTTP requests.
	StandardHTTP AOPWritePriority = 128

	// NeverDrop is used to mark a request such that it will never be rejected.
	// This is currently used as an administrative priority used for requests on
	// paths which require sudo capabilities.
	NeverDrop AOPWritePriority = 255
)

// String returns the string representation of the AOPWritePriority.
func (p AOPWritePriority) String() string {
	return strconv.FormatUint(uint64(p), 8)
}

// StringToAOPWritePriority converts a string to an AOPWritePriority.
func StringToAOPWritePriority(s string) AOPWritePriority {
	// Just swallow the error and fall back to the standard priority
	p, err := strconv.ParseUint(s, 8, 8)
	if err != nil {
		return StandardHTTP
	}
	return AOPWritePriority(p)
}

// WrapRequestPriorityHandler provides special handling for headers with
// X-Vault-AOP-Force-Reject set to `true`. This is useful for testing status
// codes and return values related to Adaptive Overload Protection without
// overloading Vault.
func WrapRequestPriorityHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if raw := req.Header.Get(VaultAOPForceRejectHeaderName); raw != "" {
			if fail, _ := parseutil.ParseBool(raw); fail {
				// Make the request fail as if Vault was overloaded. We don't
				// explicitly error out here, but rather attach some context
				// indicating that the PID controller should perform a
				// rejection.  This allows us to test errors propagated from the
				// WAL backend.
				req = req.WithContext(ContextWithRequestPriority(req.Context(), AlwaysDrop))
			}
		}
		handler.ServeHTTP(w, req)
	})
}

// ContextWithRequestPriority returns a new context derived from ctx with the
// given priority set.
func ContextWithRequestPriority(ctx context.Context, priority AOPWritePriority) context.Context {
	if _, ok := ctx.Value(logical.CtxKeyInFlightRequestPriority{}).(AOPWritePriority); ok {
		return ctx
	}

	return context.WithValue(ctx, logical.CtxKeyInFlightRequestPriority{}, priority)
}

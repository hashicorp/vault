// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build testonly

package vault

import (
	"context"
	"net/http"

	"github.com/hashicorp/vault/limits"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// RequestLimiterResponse is a struct for marshalling Request Limiter status responses.
type RequestLimiterResponse struct {
	GlobalDisabled   bool                      `json:"global_disabled" mapstructure:"global_disabled"`
	ListenerDisabled bool                      `json:"listener_disabled" mapstructure:"listener_disabled"`
	Limiters         map[string]*LimiterStatus `json:"types" mapstructure:"types"`
}

// LimiterStatus holds the per-limiter status and flags for testing.
type LimiterStatus struct {
	Enabled bool                `json:"enabled" mapstructure:"enabled"`
	Flags   limits.LimiterFlags `json:"flags,omitempty" mapstructure:"flags,omitempty"`
}

const readRequestLimiterHelpText = `
Read the current status of the request limiter.
`

func (b *SystemBackend) requestLimiterReadPath() *framework.Path {
	return &framework.Path{
		Pattern:         "internal/request-limiter/status$",
		HelpDescription: readRequestLimiterHelpText,
		HelpSynopsis:    readRequestLimiterHelpText,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.handleReadRequestLimiter,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb:   "read",
					OperationSuffix: "verbosity-level-for",
				},
				Responses: map[int][]framework.Response{
					http.StatusOK: {{
						Description: "OK",
					}},
				},
				Summary: "Read the current status of the request limiter.",
			},
		},
	}
}

// handleReadRequestLimiter returns the enabled Request Limiter status for this node.
func (b *SystemBackend) handleReadRequestLimiter(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp := &RequestLimiterResponse{
		Limiters: make(map[string]*LimiterStatus),
	}

	b.Core.limiterRegistryLock.Lock()
	registry := b.Core.limiterRegistry
	b.Core.limiterRegistryLock.Unlock()

	resp.GlobalDisabled = !registry.Enabled
	resp.ListenerDisabled = req.RequestLimiterDisabled
	enabled := !(resp.GlobalDisabled || resp.ListenerDisabled)

	for name := range limits.DefaultLimiterFlags {
		var flags limits.LimiterFlags
		if requestLimiter := b.Core.GetRequestLimiter(name); requestLimiter != nil && enabled {
			flags = requestLimiter.Flags
		}

		resp.Limiters[name] = &LimiterStatus{
			Enabled: enabled,
			Flags:   flags,
		}
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"request_limiter": resp,
		},
	}, nil
}

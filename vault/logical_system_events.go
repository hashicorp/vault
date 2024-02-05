// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"net/http"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// handleEventsSubscribe
func (b *SystemBackend) handleEventsSubscribe(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// TODO
	return logical.RespondWithStatusCode(nil, req, http.StatusNoContent)
}

// handleEventsUnsubscribe
func (b *SystemBackend) handleEventsUnsubscribe(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// TODO
	return logical.RespondWithStatusCode(nil, req, http.StatusNoContent)
}

// handleEventsListSubscriptions
func (b *SystemBackend) handleEventsListSubscriptions(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// TODO
	return logical.RespondWithStatusCode(nil, req, http.StatusNoContent)
}

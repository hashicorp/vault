// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package http

import (
	"net/http"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/quotas"
)

func entWrapGenericHandler(core *vault.Core, in http.Handler, props *vault.HandlerProperties) http.Handler {
	// Wrap the help wrapped handler with another layer with a generic
	// handler
	return wrapGenericHandler(core, in, props)
}

func entAdditionalRoutes(mux *http.ServeMux, core *vault.Core) {}

func entAdjustResponse(core *vault.Core, w http.ResponseWriter, req *logical.Request) {
}

func entRlqRequestFields(core *vault.Core, r *http.Request, quotaReq *quotas.Request) {}

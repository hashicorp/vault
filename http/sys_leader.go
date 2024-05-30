// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

// This endpoint is needed to answer queries before Vault unseals
// or becomes the leader.
func handleSysLeader(core *vault.Core, opt ...ListenerConfigOption) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysLeaderGet(core, w, r, opt...)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysLeaderGet(core *vault.Core, w http.ResponseWriter, r *http.Request, opt ...ListenerConfigOption) {
	var tokenPresent bool
	token := r.Header.Get(consts.AuthHeaderName)
	ctx := r.Context()

	if token != "" {
		// We don't care about the error, we just want to know if token exists
		lock := core.HALock()
		lock.Lock()
		tokenEntry, err := core.LookupToken(ctx, token)
		lock.Unlock()
		tokenPresent = err == nil && tokenEntry != nil
	}

	if tokenPresent {
		ctx = logical.CreateContextRedactionSettings(r.Context(), false, false, false)
	}

	resp, err := core.GetLeaderStatus(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, resp)
}

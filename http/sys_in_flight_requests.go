// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func handleUnAuthenticatedInFlightRequest(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		currentInFlightReqMap := core.LoadInFlightReqData()

		respondOk(w, currentInFlightReqMap)
	})
}

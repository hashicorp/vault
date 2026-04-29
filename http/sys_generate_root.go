// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"io"
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func handleSysGenerateRootAttempt(core *vault.Core, generateStrategy vault.GenerateRootStrategy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysGenerateRootAttemptGet(core, w, r, "")
		case "POST", "PUT":
			handleSysGenerateRootAttemptPut(core, w, r, generateStrategy)
		case "DELETE":
			handleSysGenerateRootAttemptDelete(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysGenerateRootAttemptGet(core *vault.Core, w http.ResponseWriter, r *http.Request, otp string) {
	ctx, cancel := core.GetContext()
	defer cancel()

	status, code, err := vault.HandleSysGenerateRootAttemptGet(ctx, core, otp, true)
	if err != nil {
		respondError(w, code, err)
		return
	}

	respondOk(w, status)
}

func handleSysGenerateRootAttemptPut(core *vault.Core, w http.ResponseWriter, r *http.Request, generateStrategy vault.GenerateRootStrategy) {
	// Parse the request
	var req vault.GenerateRootInitRequest
	if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil && err != io.EOF {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	otp, code, err := vault.HandleSysGenerateRootAttemptPut(core, generateStrategy, &req, true)
	if err != nil {
		respondError(w, code, err)
		return
	}

	handleSysGenerateRootAttemptGet(core, w, r, otp)
}

func handleSysGenerateRootAttemptDelete(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	code, err := vault.HandleSysGenerateRootAttemptDelete(core, true)
	if err != nil {
		respondError(w, code, err)
		return
	}
	respondOk(w, nil)
}

func handleSysGenerateRootUpdate(core *vault.Core, generateStrategy vault.GenerateRootStrategy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the request
		var req vault.GenerateRootUpdateRequest
		if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		ctx, cancel := core.GetContext()
		defer cancel()

		resp, code, err := vault.HandleSysGenerateRootUpdate(ctx, core, generateStrategy, &req, true)
		if err != nil {
			respondError(w, code, err)
			return
		}

		respondOk(w, resp)
	})
}

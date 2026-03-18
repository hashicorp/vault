// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault"
)

func handleSysRekeyInit(core *vault.Core, recovery bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		standby, _ := core.Standby()
		if standby {
			respondStandby(core, w, r)
			return
		}

		repState := core.ReplicationState()
		if repState.HasState(consts.ReplicationPerformanceSecondary) {
			respondError(w, http.StatusBadRequest,
				fmt.Errorf("rekeying can only be performed on the primary cluster when replication is activated"))
			return
		}

		ctx, cancel := core.GetContext()
		defer cancel()

		switch {
		case recovery && !core.SealAccess().RecoveryKeySupported():
			respondError(w, http.StatusBadRequest, fmt.Errorf("recovery rekeying not supported"))
		case r.Method == "GET":
			handleSysRekeyInitGet(ctx, core, recovery, w, r)
		case r.Method == "POST" || r.Method == "PUT":
			handleSysRekeyInitPut(ctx, core, recovery, w, r)
		case r.Method == "DELETE":
			handleSysRekeyInitDelete(ctx, core, recovery, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysRekeyInitGet(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	status, code, err := vault.HandleSysRekeyInitGet(ctx, core, recovery, true)
	if err != nil {
		respondError(w, code, err)
		return
	}

	respondOk(w, status)
}

func handleSysRekeyInitPut(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req *vault.RekeyRequest
	if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	code, err := vault.HandleSysRekeyInitPut(core, recovery, req, true)
	if err != nil {
		respondError(w, code, err)
		return
	}

	handleSysRekeyInitGet(ctx, core, recovery, w, r)
}

func handleSysRekeyInitDelete(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	var req vault.RekeyDeleteRequest
	if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if err := core.RekeyCancel(recovery, req.Nonce, 10*time.Minute, true); err != nil {
		respondError(w, err.Code(), err)
		return
	}
	respondOk(w, nil)
}

func handleSysRekeyUpdate(core *vault.Core, recovery bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		standby, _ := core.Standby()
		if standby {
			respondStandby(core, w, r)
			return
		}

		// Parse the request
		var req vault.RekeyUpdateRequest
		if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		ctx, cancel := core.GetContext()
		defer cancel()

		result, code, err := vault.HandleSysRekeyUpdatePut(ctx, core, recovery, &req, true)
		if err != nil {
			respondError(w, code, err)
			return
		}
		if result != nil {
			respondOk(w, result)
			return
		}

		handleSysRekeyInitGet(r.Context(), core, recovery, w, r)
	})
}

func handleSysRekeyVerify(core *vault.Core, recovery bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		standby, _ := core.Standby()
		if standby {
			respondStandby(core, w, r)
			return
		}

		repState := core.ReplicationState()
		if repState.HasState(consts.ReplicationPerformanceSecondary) {
			respondError(w, http.StatusBadRequest,
				fmt.Errorf("rekeying can only be performed on the primary cluster when replication is activated"))
			return
		}

		ctx, cancel := core.GetContext()
		defer cancel()

		switch {
		case recovery && !core.SealAccess().RecoveryKeySupported():
			respondError(w, http.StatusBadRequest, fmt.Errorf("recovery rekeying not supported"))
		case r.Method == "GET":
			handleSysRekeyVerifyGet(ctx, core, recovery, w, r)
		case r.Method == "POST" || r.Method == "PUT":
			handleSysRekeyVerifyPut(ctx, core, recovery, w, r)
		case r.Method == "DELETE":
			handleSysRekeyVerifyDelete(ctx, core, recovery, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysRekeyVerifyGet(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	status, code, err := vault.HandleSysRekeyVerifyGet(ctx, core, recovery, true)
	if err != nil {
		respondError(w, code, err)
		return
	}

	respondOk(w, status)
}

func handleSysRekeyVerifyDelete(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	if err := core.RekeyVerifyRestart(recovery, true); err != nil {
		respondError(w, err.Code(), err)
		return
	}

	handleSysRekeyVerifyGet(ctx, core, recovery, w, r)
}

func handleSysRekeyVerifyPut(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req vault.RekeyVerificationUpdateRequest
	if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := core.GetContext()
	defer cancel()

	resp, code, err := vault.HandleSysRekeyVerifyPut(ctx, core, recovery, true, &req)
	if err != nil {
		respondError(w, code, err)
		return
	}

	// Format the response
	if resp != nil {
		respondOk(w, resp)
	} else {
		handleSysRekeyVerifyGet(ctx, core, recovery, w, r)
	}
}

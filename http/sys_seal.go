// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func handleSysSeal(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, _, statusCode, err := buildLogicalRequest(core, w, r, "")
		if err != nil || statusCode != 0 {
			respondError(w, statusCode, err)
			return
		}

		switch req.Operation {
		case logical.UpdateOperation:
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Seal with the token above
		// We use context.Background since there won't be a request context if the node isn't active
		if err := core.SealWithRequest(r.Context(), req); err != nil {
			if errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
				respondError(w, http.StatusForbidden, err)
				return
			}
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		respondOk(w, nil)
	})
}

func handleSysStepDown(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, _, statusCode, err := buildLogicalRequest(core, w, r, "")
		if err != nil || statusCode != 0 {
			respondError(w, statusCode, err)
			return
		}

		switch req.Operation {
		case logical.UpdateOperation:
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Seal with the token above
		if err := core.StepDown(r.Context(), req); err != nil {
			if errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
				respondError(w, http.StatusForbidden, err)
				return
			}
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		respondOk(w, nil)
	})
}

func handleSysUnseal(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
		case "POST":
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Check if this node was removed from the cluster. If so, respond with an error and return,
		// since we don't want a removed node to be able to unseal.
		removed, ok := core.IsRemovedFromCluster()
		if ok && removed {
			respondError(w, http.StatusInternalServerError, errors.New("node was removed from a HA cluster"))
			return
		}

		// Parse the request
		var req UnsealRequest
		if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		if req.Reset {
			if !core.Sealed() {
				respondError(w, http.StatusBadRequest, errors.New("vault is unsealed"))
				return
			}
			core.ResetUnsealProcess()
			handleSysSealStatusRaw(core, w, r)
			return
		}

		if req.Key == "" {
			respondError(
				w, http.StatusBadRequest,
				errors.New("'key' must be specified in request body as JSON, or 'reset' set to true"))
			return
		}

		// Decode the key, which is base64 or hex encoded
		min, max := core.BarrierKeyLength()
		key, err := hex.DecodeString(req.Key)
		// We check min and max here to ensure that a string that is base64
		// encoded but also valid hex will not be valid and we instead base64
		// decode it
		if err != nil || len(key) < min || len(key) > max {
			key, err = base64.StdEncoding.DecodeString(req.Key)
			if err != nil {
				respondError(
					w, http.StatusBadRequest,
					errors.New("'key' must be a valid hex or base64 string"))
				return
			}
		}

		// Attempt the unseal.  If migrate was specified, the key should correspond
		// to the old seal.
		if req.Migrate {
			_, err = core.UnsealMigrate(key)
		} else {
			_, err = core.Unseal(key)
		}
		if err != nil {
			switch {
			case errwrap.ContainsType(err, new(vault.ErrInvalidKey)):
			case errwrap.Contains(err, vault.ErrBarrierInvalidKey.Error()):
			case errwrap.Contains(err, vault.ErrBarrierNotInit.Error()):
			case errwrap.Contains(err, vault.ErrBarrierSealed.Error()):
			case errwrap.Contains(err, consts.ErrStandby.Error()):
			default:
				respondError(w, http.StatusInternalServerError, err)
				return
			}
			respondError(w, http.StatusBadRequest, err)
			return
		}

		// Return the seal status
		handleSysSealStatusRaw(core, w, r)
	})
}

func handleSysSealStatus(core *vault.Core, opt ...ListenerConfigOption) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		handleSysSealStatusRaw(core, w, r, opt...)
	})
}

func handleSysSealBackendStatus(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		handleSysSealBackendStatusRaw(core, w, r)
	})
}

func handleSysSealStatusRaw(core *vault.Core, w http.ResponseWriter, r *http.Request, opt ...ListenerConfigOption) {
	ctx := r.Context()

	var tokenPresent bool
	token := r.Header.Get(consts.AuthHeaderName)
	if token != "" {
		// We don't care about the error, we just want to know if the token exists
		lock := core.HALock()
		lock.Lock()
		tokenEntry, err := core.LookupToken(ctx, token)
		lock.Unlock()
		tokenPresent = err == nil && tokenEntry != nil
	}

	// If there are is no valid token then we will redact the specified values
	if tokenPresent {
		ctx = logical.CreateContextRedactionSettings(ctx, false, false, false)
	}

	status, err := core.GetSealStatus(ctx, true)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, status)
}

func handleSysSealBackendStatusRaw(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	status, err := core.GetSealBackendStatus(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, status)
}

// Note: because we didn't provide explicit tagging in the past we can't do it
// now because if it then no longer accepts capitalized versions it could break
// clients
type UnsealRequest struct {
	Key     string
	Reset   bool
	Migrate bool
}

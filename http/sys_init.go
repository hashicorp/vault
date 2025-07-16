// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/vault"
)

func handleSysInit(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysInitGet(core, w, r)
		case "PUT", "POST":
			handleSysInitPut(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysInitGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	init, err := core.Initialized(context.Background())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, &InitStatusResponse{
		Initialized: init,
	})
}

func handleSysInitPut(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Parse the request
	var req InitRequest
	if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	// Validate init request parameters
	if err := validateInitParameters(core, req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	// Initialize
	barrierConfig := &vault.SealConfig{
		SecretShares:    req.SecretShares,
		SecretThreshold: req.SecretThreshold,
		StoredShares:    req.StoredShares,
		PGPKeys:         req.PGPKeys,
	}

	recoveryConfig := &vault.SealConfig{
		SecretShares:    req.RecoveryShares,
		SecretThreshold: req.RecoveryThreshold,
		PGPKeys:         req.RecoveryPGPKeys,
	}

	initParams := &vault.InitParams{
		BarrierConfig:   barrierConfig,
		RecoveryConfig:  recoveryConfig,
		RootTokenPGPKey: req.RootTokenPGPKey,
	}

	result, initErr := core.Initialize(ctx, initParams)
	if initErr != nil {
		if vault.IsFatalError(initErr) {
			respondError(w, http.StatusBadRequest, initErr)
			return
		} else {
			// Add a warnings field? The error will be logged in the vault log
			// already.
		}
	}

	// Encode the keys
	keys := make([]string, 0, len(result.SecretShares))
	keysB64 := make([]string, 0, len(result.SecretShares))
	for _, k := range result.SecretShares {
		keys = append(keys, hex.EncodeToString(k))
		keysB64 = append(keysB64, base64.StdEncoding.EncodeToString(k))
	}

	resp := &InitResponse{
		Keys:      keys,
		KeysB64:   keysB64,
		RootToken: result.RootToken,
	}

	if len(result.RecoveryShares) > 0 {
		resp.RecoveryKeys = make([]string, 0, len(result.RecoveryShares))
		resp.RecoveryKeysB64 = make([]string, 0, len(result.RecoveryShares))
		for _, k := range result.RecoveryShares {
			resp.RecoveryKeys = append(resp.RecoveryKeys, hex.EncodeToString(k))
			resp.RecoveryKeysB64 = append(resp.RecoveryKeysB64, base64.StdEncoding.EncodeToString(k))
		}
	}

	if err := core.UnsealWithStoredKeys(ctx); err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, resp)
}

type InitRequest struct {
	SecretShares      int      `json:"secret_shares"`
	SecretThreshold   int      `json:"secret_threshold"`
	StoredShares      int      `json:"stored_shares"`
	PGPKeys           []string `json:"pgp_keys"`
	RecoveryShares    int      `json:"recovery_shares"`
	RecoveryThreshold int      `json:"recovery_threshold"`
	RecoveryPGPKeys   []string `json:"recovery_pgp_keys"`
	RootTokenPGPKey   string   `json:"root_token_pgp_key"`
}

type InitResponse struct {
	Keys            []string `json:"keys"`
	KeysB64         []string `json:"keys_base64"`
	RecoveryKeys    []string `json:"recovery_keys,omitempty"`
	RecoveryKeysB64 []string `json:"recovery_keys_base64,omitempty"`
	RootToken       string   `json:"root_token"`
}

type InitStatusResponse struct {
	Initialized bool `json:"initialized"`
}

// Validates if the right parameters are used based on AutoUnseal
func validateInitParameters(core *vault.Core, req InitRequest) error {
	recoveryFlags := make([]string, 0)
	barrierFlags := make([]string, 0)

	if req.SecretShares != 0 {
		barrierFlags = append(barrierFlags, "secret_shares")
	}
	if req.SecretThreshold != 0 {
		barrierFlags = append(barrierFlags, "secret_threshold")
	}
	if len(req.PGPKeys) != 0 {
		barrierFlags = append(barrierFlags, "pgp_keys")
	}
	if req.RecoveryShares != 0 {
		recoveryFlags = append(recoveryFlags, "recovery_shares")
	}
	if req.RecoveryThreshold != 0 {
		recoveryFlags = append(recoveryFlags, "recovery_threshold")
	}
	if len(req.RecoveryPGPKeys) != 0 {
		recoveryFlags = append(recoveryFlags, "recovery_pgp_keys")
	}

	switch core.SealAccess().RecoveryKeySupported() {
	case true:
		if len(barrierFlags) > 0 {
			return fmt.Errorf("parameters %s not applicable to seal type %s", strings.Join(barrierFlags, ","), core.SealAccess().BarrierSealConfigType())
		}
	default:
		if len(recoveryFlags) > 0 {
			return fmt.Errorf("parameters %s not applicable to seal type %s", strings.Join(recoveryFlags, ","), core.SealAccess().BarrierSealConfigType())
		}

	}
	return nil
}

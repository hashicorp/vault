// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-secure-stdlib/base62"
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

	// Get the current seal configuration
	barrierConfig, err := core.SealAccess().BarrierConfig(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if barrierConfig == nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("server is not yet initialized"))
		return
	}

	sealConfig := barrierConfig
	if core.SealAccess().RecoveryKeySupported() {
		sealConfig, err = core.SealAccess().RecoveryConfig(ctx)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
	}

	// Get the generation configuration
	generationConfig, err := core.GenerateRootConfiguration()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Get the progress
	progress, err := core.GenerateRootProgress()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	var otpLength int
	if core.DisableSSCTokens() {
		otpLength = vault.TokenLength + vault.OldTokenPrefixLength
	} else {
		otpLength = vault.TokenLength + vault.TokenPrefixLength
	}

	// Format the status
	status := &GenerateRootStatusResponse{
		Started:   false,
		Progress:  progress,
		Required:  sealConfig.SecretThreshold,
		Complete:  false,
		OTPLength: otpLength,
		OTP:       otp,
	}
	if generationConfig != nil {
		status.Nonce = generationConfig.Nonce
		status.Started = true
		status.PGPFingerprint = generationConfig.PGPFingerprint
	}

	respondOk(w, status)
}

func handleSysGenerateRootAttemptPut(core *vault.Core, w http.ResponseWriter, r *http.Request, generateStrategy vault.GenerateRootStrategy) {
	// Parse the request
	var req GenerateRootInitRequest
	if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil && err != io.EOF {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	var err error
	var genned bool

	switch {
	case len(req.PGPKey) > 0, len(req.OTP) > 0:
	default:
		genned = true
		if core.DisableSSCTokens() {
			req.OTP, err = base62.Random(vault.TokenLength + vault.OldTokenPrefixLength)
		} else {
			req.OTP, err = base62.Random(vault.TokenLength + vault.TokenPrefixLength)
		}
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
	}

	// Attemptialize the generation
	if err := core.GenerateRootInit(req.OTP, req.PGPKey, generateStrategy); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if genned {
		handleSysGenerateRootAttemptGet(core, w, r, req.OTP)
		return
	}

	handleSysGenerateRootAttemptGet(core, w, r, "")
}

func handleSysGenerateRootAttemptDelete(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	err := core.GenerateRootCancel()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, nil)
}

func handleSysGenerateRootUpdate(core *vault.Core, generateStrategy vault.GenerateRootStrategy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the request
		var req GenerateRootUpdateRequest
		if _, err := parseJSONRequest(core.PerfStandby(), r, w, &req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}
		if req.Key == "" {
			respondError(
				w, http.StatusBadRequest,
				errors.New("'key' must be specified in request body as JSON"))
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

		ctx, cancel := core.GetContext()
		defer cancel()

		// Use the key to make progress on root generation
		result, err := core.GenerateRootUpdate(ctx, key, req.Nonce, generateStrategy)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		resp := &GenerateRootStatusResponse{
			Complete:       result.Progress == result.Required,
			Nonce:          req.Nonce,
			Progress:       result.Progress,
			Required:       result.Required,
			Started:        true,
			EncodedToken:   result.EncodedToken,
			PGPFingerprint: result.PGPFingerprint,
		}

		if generateStrategy == vault.GenerateStandardRootTokenStrategy {
			resp.EncodedRootToken = result.EncodedToken
		}

		respondOk(w, resp)
	})
}

type GenerateRootInitRequest struct {
	OTP    string `json:"otp"`
	PGPKey string `json:"pgp_key"`
}

type GenerateRootStatusResponse struct {
	Nonce            string `json:"nonce"`
	Started          bool   `json:"started"`
	Progress         int    `json:"progress"`
	Required         int    `json:"required"`
	Complete         bool   `json:"complete"`
	EncodedToken     string `json:"encoded_token"`
	EncodedRootToken string `json:"encoded_root_token"`
	PGPFingerprint   string `json:"pgp_fingerprint"`
	OTP              string `json:"otp"`
	OTPLength        int    `json:"otp_length"`
}

type GenerateRootUpdateRequest struct {
	Nonce string
	Key   string
}

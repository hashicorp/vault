package http

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func handleSysGenerateRootAttempt(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysGenerateRootAttemptGet(core, w, r)
		case "POST", "PUT":
			handleSysGenerateRootAttemptPut(core, w, r)
		case "DELETE":
			handleSysGenerateRootAttemptDelete(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysGenerateRootAttemptGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Get the current seal configuration
	barrierConfig, err := core.SealAccess().BarrierConfig()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if barrierConfig == nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf(
			"server is not yet initialized"))
		return
	}

	sealConfig := barrierConfig
	if core.SealAccess().RecoveryKeySupported() {
		sealConfig, err = core.SealAccess().RecoveryConfig()
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

	// Format the status
	status := &GenerateRootStatusResponse{
		Started:  false,
		Progress: progress,
		Required: sealConfig.SecretThreshold,
		Complete: false,
	}
	if generationConfig != nil {
		status.Nonce = generationConfig.Nonce
		status.Started = true
		status.PGPFingerprint = generationConfig.PGPFingerprint
	}

	respondOk(w, status)
}

func handleSysGenerateRootAttemptPut(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req GenerateRootInitRequest
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if len(req.OTP) > 0 && len(req.PGPKey) > 0 {
		respondError(w, http.StatusBadRequest, fmt.Errorf("only one of \"otp\" and \"pgp_key\" must be specified"))
		return
	}

	// Attemptialize the generation
	err := core.GenerateRootInit(req.OTP, req.PGPKey)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	handleSysGenerateRootAttemptGet(core, w, r)
}

func handleSysGenerateRootAttemptDelete(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	err := core.GenerateRootCancel()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, nil)
}

func handleSysGenerateRootUpdate(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the request
		var req GenerateRootUpdateRequest
		if err := parseRequest(r, &req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}
		if req.Key == "" {
			respondError(
				w, http.StatusBadRequest,
				errors.New("'key' must specified in request body as JSON"))
			return
		}

		// Decode the key, which is hex encoded
		key, err := hex.DecodeString(req.Key)
		if err != nil {
			respondError(
				w, http.StatusBadRequest,
				errors.New("'key' must be a valid hex-string"))
			return
		}

		// Use the key to make progress on root generation
		result, err := core.GenerateRootUpdate(key, req.Nonce)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		resp := &GenerateRootStatusResponse{
			Complete:         result.Progress == result.Required,
			Nonce:            req.Nonce,
			Progress:         result.Progress,
			Required:         result.Required,
			Started:          true,
			EncodedRootToken: result.EncodedRootToken,
			PGPFingerprint:   result.PGPFingerprint,
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
	EncodedRootToken string `json:"encoded_root_token"`
	PGPFingerprint   string `json:"pgp_fingerprint"`
}

type GenerateRootUpdateRequest struct {
	Nonce string
	Key   string
}

package http

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func handleSysGenerateShareAttempt(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysGenerateShareAttemptGet(core, w, r)
		case "POST", "PUT":
			handleSysGenerateShareAttemptPut(core, w, r)
		case "DELETE":
			handleSysGenerateShareAttemptDelete(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysGenerateShareAttemptGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {
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
	generationConfig, err := core.GenerateShareConfiguration()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Get the progress
	progress, err := core.GenerateShareProgress()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Format the status
	status := &GenerateShareStatusResponse{
		Started:  false,
		Progress: progress,
		Required: sealConfig.SecretThreshold,
		Complete: false,
	}
	if generationConfig != nil {
		status.Started = true
		status.PGPFingerprint = generationConfig.PGPFingerprint
	}

	respondOk(w, status)
}

func handleSysGenerateShareAttemptPut(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req GenerateShareInitRequest
	if err := parseRequest(r, w, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	// Attemptialize the generation
	err := core.GenerateShareInit(req.PGPKey)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	handleSysGenerateShareAttemptGet(core, w, r)
}

func handleSysGenerateShareAttemptDelete(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	err := core.GenerateShareCancel()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, nil)
}

func handleSysGenerateShareUpdate(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the request
		var req GenerateShareUpdateRequest
		if err := parseRequest(r, w, &req); err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}
		if req.Key == "" {
			respondError(
				w, http.StatusBadRequest,
				errors.New("'key' must specified in request body as JSON"))
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

		// Use the key to make progress on root generation
		result, err := core.GenerateShareUpdate(key)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		resp := &GenerateShareStatusResponse{
			Complete:       result.Progress == result.Required,
			Progress:       result.Progress,
			Required:       result.Required,
			Started:        true,
			Key:            result.Key,
			PGPFingerprint: result.PGPFingerprint,
		}

		respondOk(w, resp)
	})
}

type GenerateShareInitRequest struct {
	PGPKey string `json:"pgp_key"`
}

type GenerateShareStatusResponse struct {
	Started        bool   `json:"started"`
	Progress       int    `json:"progress"`
	Required       int    `json:"required"`
	Complete       bool   `json:"complete"`
	Key            string `json:"key"`
	PGPFingerprint string `json:"pgp_fingerprint"`
}

type GenerateShareUpdateRequest struct {
	Key string
}

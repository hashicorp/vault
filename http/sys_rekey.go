package http

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func handleSysRekeyInit(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysRekeyInitGet(core, w, r)
		case "POST", "PUT":
			handleSysRekeyInitPut(core, w, r)
		case "DELETE":
			handleSysRekeyInitDelete(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysRekeyInitGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Get the current configuration
	sealConfig, err := core.SealConfig()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if sealConfig == nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf(
			"server is not yet initialized"))
		return
	}

	// Get the rekey configuration
	rekeyConf, err := core.RekeyConfig()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Get the progress
	progress, err := core.RekeyProgress()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Format the status
	status := &RekeyStatusResponse{
		Started:  false,
		T:        0,
		N:        0,
		Progress: progress,
		Required: sealConfig.SecretThreshold,
	}
	if rekeyConf != nil {
		status.Started = true
		status.T = rekeyConf.SecretThreshold
		status.N = rekeyConf.SecretShares
	}
	respondOk(w, status)
}

func handleSysRekeyInitPut(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req RekeyRequest
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	// Initialize the rekey
	err := core.RekeyInit(&vault.SealConfig{
		SecretShares:    req.SecretShares,
		SecretThreshold: req.SecretThreshold,
	})
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	respondOk(w, nil)
}

func handleSysRekeyInitDelete(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	err := core.RekeyCancel()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, nil)
}

func handleSysRekeyUpdate(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}

		// Parse the request
		var req RekeyUpdateRequest
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

		// Use the key to make progress on rekey
		result, err := core.RekeyUpdate(key)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		// Format the response
		resp := &RekeyUpdateResponse{}
		if result != nil {
			resp.Complete = true

			// Encode the keys
			keys := make([]string, 0, len(result.SecretShares))
			for _, k := range result.SecretShares {
				keys = append(keys, hex.EncodeToString(k))
			}
			resp.Keys = keys
		}
		respondOk(w, resp)
	})
}

type RekeyRequest struct {
	SecretShares    int `json:"secret_shares"`
	SecretThreshold int `json:"secret_threshold"`
}

type RekeyStatusResponse struct {
	Started  bool `json:"started"`
	T        int  `json:"t"`
	N        int  `json:"n"`
	Progress int  `json:"progress"`
	Required int  `json:"required"`
}

type RekeyUpdateRequest struct {
	Key string
}

type RekeyUpdateResponse struct {
	Complete bool     `json:"complete"`
	Keys     []string `json:"keys"`
}

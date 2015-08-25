package http

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/vault"
)

func handleSysInit(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleSysInitGet(core, w, r)
		case "PUT":
			handleSysInitPut(core, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysInitGet(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	init, err := core.Initialized()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, &InitStatusResponse{
		Initialized: init,
	})
}

func handleSysInitPut(core *vault.Core, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req InitRequest
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	switch {
	case req.SecretShares > 0 && len(req.SecretPGPKeys) > 0 && len(req.SecretPGPKeys) != req.SecretShares:
		respondError(w, http.StatusBadRequest, fmt.Errorf("Mismatch between key-shares and length of pgp-keys (you can specify pgp-keys alone)"))
		return
	case req.SecretShares == 0 && len(req.SecretPGPKeys) > 0:
		req.SecretShares = len(req.SecretPGPKeys)
	}

	// Initialize
	result, err := core.Initialize(&vault.SealConfig{
		SecretShares:    req.SecretShares,
		SecretThreshold: req.SecretThreshold,
		SecretPGPKeys:   req.SecretPGPKeys,
	})
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	// Encode the keys
	keys := make([]string, 0, len(result.SecretShares))
	for _, k := range result.SecretShares {
		keys = append(keys, hex.EncodeToString(k))
	}

	respondOk(w, &InitResponse{
		Keys:      keys,
		RootToken: result.RootToken,
	})
}

type InitRequest struct {
	SecretShares    int      `json:"secret_shares"`
	SecretThreshold int      `json:"secret_threshold"`
	SecretPGPKeys   []string `json:"secret_pgp_keys"`
}

type InitResponse struct {
	Keys      []string `json:"keys"`
	RootToken string   `json:"root_token"`
}

type InitStatusResponse struct {
	Initialized bool `json:"initialized"`
}

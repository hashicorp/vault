package http

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/hashicorp/errwrap"
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

	if core.SealAccess().StoredKeysSupported() {
		if barrierConfig.SecretShares != 1 {
			respondError(w, http.StatusBadRequest, fmt.Errorf("secret shares must be 1"))
			return
		}
		if barrierConfig.SecretThreshold != barrierConfig.SecretShares {
			respondError(w, http.StatusBadRequest, fmt.Errorf("secret threshold must be same as secret shares"))
			return
		}
		if barrierConfig.StoredShares != barrierConfig.SecretShares {
			respondError(w, http.StatusBadRequest, fmt.Errorf("stored shares must be same as secret shares"))
			return
		}
		if barrierConfig.PGPKeys != nil && len(barrierConfig.PGPKeys) > 0 {
			respondError(w, http.StatusBadRequest, fmt.Errorf("PGP keys not supported when storing shares"))
			return
		}
	} else {
		if barrierConfig.StoredShares > 0 {
			respondError(w, http.StatusBadRequest, fmt.Errorf("stored keys are not supported"))
			return
		}
	}

	result, initErr := core.Initialize(barrierConfig, recoveryConfig)
	if initErr != nil {
		if !errwrap.ContainsType(initErr, new(vault.NonFatalError)) {
			respondError(w, http.StatusBadRequest, initErr)
			return
		} else {
			// Add a warnings field? The error will be logged in the vault log
			// already.
		}
	}

	// Encode the keys
	keys := make([]string, 0, len(result.SecretShares))
	for _, k := range result.SecretShares {
		keys = append(keys, hex.EncodeToString(k))
	}

	resp := &InitResponse{
		Keys:      keys,
		RootToken: result.RootToken,
	}

	if len(result.RecoveryShares) > 0 {
		resp.RecoveryKeys = make([]string, 0, len(result.RecoveryShares))
		for _, k := range result.RecoveryShares {
			resp.RecoveryKeys = append(resp.RecoveryKeys, hex.EncodeToString(k))
		}
	}

	core.UnsealWithStoredKeys()

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
}

type InitResponse struct {
	Keys         []string `json:"keys"`
	RecoveryKeys []string `json:"recovery_keys,omitempty"`
	RootToken    string   `json:"root_token"`
}

type InitStatusResponse struct {
	Initialized bool `json:"initialized"`
}

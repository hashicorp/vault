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

	// If idempotentcy is enabled, check if we are initialized
	if req.Idempotent {
		inited, err := core.Initialized()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		if inited {
			sealConf, err := core.SealConfig()
			if err != nil {
				respondError(w, http.StatusInternalServerError, err)
				return
			}

			// Ensure that the parameters are the same
			switch {
			case sealConf.SecretShares != req.SecretShares:
				respondError(w, http.StatusBadRequest, fmt.Errorf(
					"requested init secret shares %d does not match current value of %d",
					req.SecretShares, sealConf.SecretShares,
				))
				return
			case sealConf.SecretThreshold != req.SecretThreshold:
				respondError(w, http.StatusBadRequest, fmt.Errorf(
					"requested init secret threshold %d does not match current value of %d",
					req.SecretThreshold, sealConf.SecretThreshold,
				))
				return
			case (sealConf.PGPKeys != nil && len(sealConf.PGPKeys) > 0) ||
				(req.PGPKeys != nil && len(req.PGPKeys) > 0):
				if req.PGPKeys == nil || sealConf.PGPKeys == nil {
					respondError(w, http.StatusBadRequest, fmt.Errorf(
						"requested init PGP keys does not match current set",
					))
					return
				}
				if len(req.PGPKeys) != len(sealConf.PGPKeys) {
					respondError(w, http.StatusBadRequest, fmt.Errorf(
						"requested init PGP keys does not match current set",
					))
					return
				}
				for i, v := range req.PGPKeys {
					if v != sealConf.PGPKeys[i] {
						respondError(w, http.StatusBadRequest, fmt.Errorf(
							"requested init PGP keys does not match current set",
						))
						return
					}
				}

				// It's all good, continue with the default
				fallthrough
			default:
				respondOk(w, &InitResponse{})
				return
			}
		}
	}

	// Initialize
	result, err := core.Initialize(&vault.SealConfig{
		SecretShares:    req.SecretShares,
		SecretThreshold: req.SecretThreshold,
		PGPKeys:         req.PGPKeys,
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
	PGPKeys         []string `json:"pgp_keys"`
	Idempotent      bool     `json:"idempotent"`
}

type InitResponse struct {
	Keys      []string `json:"keys"`
	RootToken string   `json:"root_token"`
}

type InitStatusResponse struct {
	Initialized bool `json:"initialized"`
}

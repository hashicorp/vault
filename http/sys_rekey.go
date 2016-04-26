package http

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/vault"
)

func handleSysRekeyInit(core *vault.Core, recovery bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case recovery && !core.SealAccess().RecoveryKeySupported():
			respondError(w, http.StatusBadRequest, fmt.Errorf("recovery rekeying not supported"))
		case r.Method == "GET":
			handleSysRekeyInitGet(core, recovery, w, r)
		case r.Method == "POST" || r.Method == "PUT":
			handleSysRekeyInitPut(core, recovery, w, r)
		case r.Method == "DELETE":
			handleSysRekeyInitDelete(core, recovery, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysRekeyInitGet(core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
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

	// Get the rekey configuration
	rekeyConf, err := core.RekeyConfig(recovery)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	// Get the progress
	progress, err := core.RekeyProgress(recovery)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}

	sealThreshold, err := core.RekeyThreshold(recovery)
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
		Required: sealThreshold,
	}
	if rekeyConf != nil {
		status.Nonce = rekeyConf.Nonce
		status.Started = true
		status.T = rekeyConf.SecretThreshold
		status.N = rekeyConf.SecretShares
		if rekeyConf.PGPKeys != nil && len(rekeyConf.PGPKeys) != 0 {
			pgpFingerprints, err := pgpkeys.GetFingerprints(rekeyConf.PGPKeys, nil)
			if err != nil {
				respondError(w, http.StatusInternalServerError, err)
				return
			}
			status.PGPFingerprints = pgpFingerprints
			status.Backup = rekeyConf.Backup
		}
	}
	respondOk(w, status)
}

func handleSysRekeyInitPut(core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req RekeyRequest
	if err := parseRequest(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if req.Backup && len(req.PGPKeys) == 0 {
		respondError(w, http.StatusBadRequest, fmt.Errorf("cannot request a backup of the new keys without providing PGP keys for encryption"))
		return
	}

	// Right now we don't support this, but the rest of the code is ready for
	// when we do, hence the check below for this to be false if
	// StoredShares is greater than zero
	if core.SealAccess().StoredKeysSupported() {
		respondError(w, http.StatusBadRequest, fmt.Errorf("rekeying of barrier not supported when stored key support is available"))
		return
	}

	// Initialize the rekey
	err := core.RekeyInit(&vault.SealConfig{
		SecretShares:    req.SecretShares,
		SecretThreshold: req.SecretThreshold,
		StoredShares:    req.StoredShares,
		PGPKeys:         req.PGPKeys,
		Backup:          req.Backup,
	}, recovery)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	handleSysRekeyInitGet(core, recovery, w, r)
}

func handleSysRekeyInitDelete(core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	err := core.RekeyCancel(recovery)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, nil)
}

func handleSysRekeyUpdate(core *vault.Core, recovery bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		result, err := core.RekeyUpdate(key, req.Nonce, recovery)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		// Format the response
		resp := &RekeyUpdateResponse{}
		if result != nil {
			resp.Complete = true
			resp.Nonce = req.Nonce

			// Encode the keys
			keys := make([]string, 0, len(result.SecretShares))
			for _, k := range result.SecretShares {
				keys = append(keys, hex.EncodeToString(k))
			}
			resp.Keys = keys

			resp.Backup = result.Backup
			resp.PGPFingerprints = result.PGPFingerprints
		}
		respondOk(w, resp)
	})
}

type RekeyRequest struct {
	SecretShares    int      `json:"secret_shares"`
	SecretThreshold int      `json:"secret_threshold"`
	StoredShares    int      `json:"stored_shares"`
	PGPKeys         []string `json:"pgp_keys"`
	Backup          bool     `json:"backup"`
}

type RekeyStatusResponse struct {
	Nonce           string   `json:"nonce"`
	Started         bool     `json:"started"`
	T               int      `json:"t"`
	N               int      `json:"n"`
	Progress        int      `json:"progress"`
	Required        int      `json:"required"`
	PGPFingerprints []string `json:"pgp_fingerprints"`
	Backup          bool     `json:"backup"`
}

type RekeyUpdateRequest struct {
	Nonce string
	Key   string
}

type RekeyUpdateResponse struct {
	Nonce           string   `json:"nonce"`
	Complete        bool     `json:"complete"`
	Keys            []string `json:"keys"`
	PGPFingerprints []string `json:"pgp_fingerprints"`
	Backup          bool     `json:"backup"`
}

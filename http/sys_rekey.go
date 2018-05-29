package http

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/vault"
)

func handleSysRekeyInit(core *vault.Core, recovery bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		standby, _ := core.Standby()
		if standby {
			respondStandby(core, w, r.URL)
			return
		}

		repState := core.ReplicationState()
		if repState.HasState(consts.ReplicationPerformanceSecondary) {
			respondError(w, http.StatusBadRequest,
				fmt.Errorf("rekeying can only be performed on the primary cluster when replication is activated"))
			return
		}

		ctx, cancel := core.GetContext()
		defer cancel()

		switch {
		case recovery && !core.SealAccess().RecoveryKeySupported():
			respondError(w, http.StatusBadRequest, fmt.Errorf("recovery rekeying not supported"))
		case r.Method == "GET":
			handleSysRekeyInitGet(ctx, core, recovery, w, r)
		case r.Method == "POST" || r.Method == "PUT":
			handleSysRekeyInitPut(ctx, core, recovery, w, r)
		case r.Method == "DELETE":
			handleSysRekeyInitDelete(ctx, core, recovery, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysRekeyInitGet(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	barrierConfig, barrierConfErr := core.SealAccess().BarrierConfig(ctx)
	if barrierConfErr != nil {
		respondError(w, http.StatusInternalServerError, barrierConfErr)
		return
	}
	if barrierConfig == nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("server is not yet initialized"))
		return
	}

	// Get the rekey configuration
	rekeyConf, err := core.RekeyConfig(recovery)
	if err != nil {
		respondError(w, err.Code(), err)
		return
	}

	sealThreshold, err := core.RekeyThreshold(ctx, recovery)
	if err != nil {
		respondError(w, err.Code(), err)
		return
	}

	// Format the status
	status := &RekeyStatusResponse{
		Started:  false,
		T:        0,
		N:        0,
		Required: sealThreshold,
	}
	if rekeyConf != nil {
		// Get the progress
		started, progress, err := core.RekeyProgress(recovery, false)
		if err != nil {
			respondError(w, err.Code(), err)
			return
		}

		status.Nonce = rekeyConf.Nonce
		status.Started = started
		status.T = rekeyConf.SecretThreshold
		status.N = rekeyConf.SecretShares
		status.Progress = progress
		status.VerificationRequired = rekeyConf.VerificationRequired
		status.VerificationNonce = rekeyConf.VerificationNonce
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

func handleSysRekeyInitPut(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req RekeyRequest
	if err := parseRequest(r, w, &req); err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	if req.Backup && len(req.PGPKeys) == 0 {
		respondError(w, http.StatusBadRequest, fmt.Errorf("cannot request a backup of the new keys without providing PGP keys for encryption"))
		return
	}

	if len(req.PGPKeys) > 0 && len(req.PGPKeys) != req.SecretShares {
		respondError(w, http.StatusBadRequest, fmt.Errorf("incorrect number of PGP keys for rekey"))
		return
	}

	// Initialize the rekey
	err := core.RekeyInit(&vault.SealConfig{
		SecretShares:         req.SecretShares,
		SecretThreshold:      req.SecretThreshold,
		StoredShares:         req.StoredShares,
		PGPKeys:              req.PGPKeys,
		Backup:               req.Backup,
		VerificationRequired: req.RequireVerification,
	}, recovery)
	if err != nil {
		respondError(w, err.Code(), err)
		return
	}

	handleSysRekeyInitGet(ctx, core, recovery, w, r)
}

func handleSysRekeyInitDelete(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	if err := core.RekeyCancel(recovery); err != nil {
		respondError(w, err.Code(), err)
		return
	}
	respondOk(w, nil)
}

func handleSysRekeyUpdate(core *vault.Core, recovery bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		standby, _ := core.Standby()
		if standby {
			respondStandby(core, w, r.URL)
			return
		}

		// Parse the request
		var req RekeyUpdateRequest
		if err := parseRequest(r, w, &req); err != nil {
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

		// Use the key to make progress on rekey
		result, rekeyErr := core.RekeyUpdate(ctx, key, req.Nonce, recovery)
		if rekeyErr != nil {
			respondError(w, rekeyErr.Code(), rekeyErr)
			return
		}

		// Format the response
		resp := &RekeyUpdateResponse{}
		if result != nil {
			resp.Complete = true
			resp.Nonce = req.Nonce
			resp.Backup = result.Backup
			resp.PGPFingerprints = result.PGPFingerprints
			resp.VerificationRequired = result.VerificationRequired
			resp.VerificationNonce = result.VerificationNonce

			// Encode the keys
			keys := make([]string, 0, len(result.SecretShares))
			keysB64 := make([]string, 0, len(result.SecretShares))
			for _, k := range result.SecretShares {
				keys = append(keys, hex.EncodeToString(k))
				keysB64 = append(keysB64, base64.StdEncoding.EncodeToString(k))
			}
			resp.Keys = keys
			resp.KeysB64 = keysB64
			respondOk(w, resp)
		} else {
			handleSysRekeyInitGet(ctx, core, recovery, w, r)
		}
	})
}

func handleSysRekeyVerify(core *vault.Core, recovery bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		standby, _ := core.Standby()
		if standby {
			respondStandby(core, w, r.URL)
			return
		}

		repState := core.ReplicationState()
		if repState.HasState(consts.ReplicationPerformanceSecondary) {
			respondError(w, http.StatusBadRequest,
				fmt.Errorf("rekeying can only be performed on the primary cluster when replication is activated"))
			return
		}

		ctx, cancel := core.GetContext()
		defer cancel()

		switch {
		case recovery && !core.SealAccess().RecoveryKeySupported():
			respondError(w, http.StatusBadRequest, fmt.Errorf("recovery rekeying not supported"))
		case r.Method == "GET":
			handleSysRekeyVerifyGet(ctx, core, recovery, w, r)
		case r.Method == "POST" || r.Method == "PUT":
			handleSysRekeyVerifyPut(ctx, core, recovery, w, r)
		case r.Method == "DELETE":
			handleSysRekeyVerifyDelete(ctx, core, recovery, w, r)
		default:
			respondError(w, http.StatusMethodNotAllowed, nil)
		}
	})
}

func handleSysRekeyVerifyGet(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	barrierConfig, barrierConfErr := core.SealAccess().BarrierConfig(ctx)
	if barrierConfErr != nil {
		respondError(w, http.StatusInternalServerError, barrierConfErr)
		return
	}
	if barrierConfig == nil {
		respondError(w, http.StatusBadRequest, fmt.Errorf("server is not yet initialized"))
		return
	}

	// Get the rekey configuration
	rekeyConf, err := core.RekeyConfig(recovery)
	if err != nil {
		respondError(w, err.Code(), err)
		return
	}
	if rekeyConf == nil {
		respondError(w, http.StatusBadRequest, errors.New("no rekey configuration found"))
		return
	}

	// Get the progress
	started, progress, err := core.RekeyProgress(recovery, true)
	if err != nil {
		respondError(w, err.Code(), err)
		return
	}

	// Format the status
	status := &RekeyVerificationStatusResponse{
		Started:  started,
		Nonce:    rekeyConf.VerificationNonce,
		T:        rekeyConf.SecretThreshold,
		N:        rekeyConf.SecretShares,
		Progress: progress,
	}
	respondOk(w, status)
}

func handleSysRekeyVerifyDelete(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	if err := core.RekeyVerifyRestart(recovery); err != nil {
		respondError(w, err.Code(), err)
		return
	}

	handleSysRekeyVerifyGet(ctx, core, recovery, w, r)
}

func handleSysRekeyVerifyPut(ctx context.Context, core *vault.Core, recovery bool, w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req RekeyVerificationUpdateRequest
	if err := parseRequest(r, w, &req); err != nil {
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

	// Use the key to make progress on rekey
	result, rekeyErr := core.RekeyVerify(ctx, key, req.Nonce, recovery)
	if rekeyErr != nil {
		respondError(w, rekeyErr.Code(), rekeyErr)
		return
	}

	// Format the response
	resp := &RekeyVerificationUpdateResponse{}
	if result != nil {
		resp.Complete = true
		resp.Nonce = result.Nonce
		respondOk(w, resp)
	} else {
		handleSysRekeyVerifyGet(ctx, core, recovery, w, r)
	}
}

type RekeyRequest struct {
	SecretShares        int      `json:"secret_shares"`
	SecretThreshold     int      `json:"secret_threshold"`
	StoredShares        int      `json:"stored_shares"`
	PGPKeys             []string `json:"pgp_keys"`
	Backup              bool     `json:"backup"`
	RequireVerification bool     `json:"require_verification"`
}

type RekeyStatusResponse struct {
	Nonce                string   `json:"nonce"`
	Started              bool     `json:"started"`
	T                    int      `json:"t"`
	N                    int      `json:"n"`
	Progress             int      `json:"progress"`
	Required             int      `json:"required"`
	PGPFingerprints      []string `json:"pgp_fingerprints"`
	Backup               bool     `json:"backup"`
	VerificationRequired bool     `json:"verification_required"`
	VerificationNonce    string   `json:"verification_nonce,omitempty"`
}

type RekeyUpdateRequest struct {
	Nonce string
	Key   string
}

type RekeyUpdateResponse struct {
	Nonce                string   `json:"nonce"`
	Complete             bool     `json:"complete"`
	Keys                 []string `json:"keys"`
	KeysB64              []string `json:"keys_base64"`
	PGPFingerprints      []string `json:"pgp_fingerprints"`
	Backup               bool     `json:"backup"`
	VerificationRequired bool     `json:"verification_required"`
	VerificationNonce    string   `json:"verification_nonce,omitempty"`
}

type RekeyVerificationUpdateRequest struct {
	Nonce string `json:"nonce"`
	Key   string `json:"key"`
}

type RekeyVerificationStatusResponse struct {
	Nonce    string `json:"nonce"`
	Started  bool   `json:"started"`
	T        int    `json:"t"`
	N        int    `json:"n"`
	Progress int    `json:"progress"`
}

type RekeyVerificationUpdateResponse struct {
	Nonce    string `json:"nonce"`
	Complete bool   `json:"complete"`
}

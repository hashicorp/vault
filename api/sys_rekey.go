package api

import (
	"context"
	"net/http"
)

func (c *Sys) RekeyStatus() (*RekeyStatusResponse, error) {
	return c.RekeyStatusWithContext(context.Background())
}

func (c *Sys) RekeyStatusWithContext(ctx context.Context) (*RekeyStatusResponse, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/rekey/init")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRecoveryKeyStatus() (*RekeyStatusResponse, error) {
	return c.RekeyRecoveryKeyStatusWithContext(context.Background())
}

func (c *Sys) RekeyRecoveryKeyStatusWithContext(ctx context.Context) (*RekeyStatusResponse, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/rekey-recovery-key/init")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyInit(config *RekeyInitRequest) (*RekeyStatusResponse, error) {
	return c.RekeyInitWithContext(context.Background(), config)
}

func (c *Sys) RekeyInitWithContext(ctx context.Context, config *RekeyInitRequest) (*RekeyStatusResponse, error) {
	req := c.c.NewRequest(http.MethodPut, "/v1/sys/rekey/init")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(config); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRecoveryKeyInit(config *RekeyInitRequest) (*RekeyStatusResponse, error) {
	return c.RekeyRecoveryKeyInitWithContext(context.Background(), config)
}

func (c *Sys) RekeyRecoveryKeyInitWithContext(ctx context.Context, config *RekeyInitRequest) (*RekeyStatusResponse, error) {
	req := c.c.NewRequest(http.MethodPut, "/v1/sys/rekey-recovery-key/init")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(config); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyCancel() error {
	return c.RekeyCancelWithContext(context.Background())
}

func (c *Sys) RekeyCancelWithContext(ctx context.Context) error {
	req := c.c.NewRequest(http.MethodDelete, "/v1/sys/rekey/init")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RekeyRecoveryKeyCancel() error {
	return c.RekeyRecoveryKeyCancelWithContext(context.Background())
}

func (c *Sys) RekeyRecoveryKeyCancelWithContext(ctx context.Context) error {
	req := c.c.NewRequest(http.MethodDelete, "/v1/sys/rekey-recovery-key/init")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RekeyUpdate(shard, nonce string) (*RekeyUpdateResponse, error) {
	return c.RekeyUpdateWithContext(context.Background(), shard, nonce)
}

func (c *Sys) RekeyUpdateWithContext(ctx context.Context, shard, nonce string) (*RekeyUpdateResponse, error) {
	body := map[string]interface{}{
		"key":   shard,
		"nonce": nonce,
	}

	req := c.c.NewRequest(http.MethodPut, "/v1/sys/rekey/update")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyUpdateResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRecoveryKeyUpdate(shard, nonce string) (*RekeyUpdateResponse, error) {
	return c.RekeyRecoveryKeyUpdateWithContext(context.Background(), shard, nonce)
}

func (c *Sys) RekeyRecoveryKeyUpdateWithContext(ctx context.Context, shard, nonce string) (*RekeyUpdateResponse, error) {
	body := map[string]interface{}{
		"key":   shard,
		"nonce": nonce,
	}

	req := c.c.NewRequest(http.MethodPut, "/v1/sys/rekey-recovery-key/update")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyUpdateResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRetrieveBackup() (*RekeyRetrieveResponse, error) {
	return c.RekeyRetrieveBackupWithContext(context.Background())
}

func (c *Sys) RekeyRetrieveBackupWithContext(ctx context.Context) (*RekeyRetrieveResponse, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/rekey/backup")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyRetrieveResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyRetrieveRecoveryBackup() (*RekeyRetrieveResponse, error) {
	return c.RekeyRetrieveRecoveryBackupWithContext(context.Background())
}

func (c *Sys) RekeyRetrieveRecoveryBackupWithContext(ctx context.Context) (*RekeyRetrieveResponse, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/rekey/recovery-backup")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RekeyRetrieveResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) RekeyDeleteBackup() error {
	return c.RekeyDeleteBackupWithContext(context.Background())
}

func (c *Sys) RekeyDeleteBackupWithContext(ctx context.Context) error {
	req := c.c.NewRequest(http.MethodDelete, "/v1/sys/rekey/backup")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}

	return err
}

func (c *Sys) RekeyDeleteRecoveryBackup() error {
	return c.RekeyDeleteRecoveryBackupWithContext(context.Background())
}

func (c *Sys) RekeyDeleteRecoveryBackupWithContext(ctx context.Context) error {
	req := c.c.NewRequest(http.MethodDelete, "/v1/sys/rekey/recovery-backup")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}

	return err
}

// RekeyInitRequest is used as input to the rekey request.
type RekeyInitRequest struct {
	// SecretShares is the new number of secret shares to use.
	SecretShares int `json:"secret_shares"`

	// SecretThreshold is the new number of secret threshold to use/
	SecretThreshold int `json:"secret_threshold"`

	// StoredShares is the new number of shares to store on the HSM.
	StoredShares int `json:"stored_shares"`

	// PGPKeys is the list of PGP keys to encrypt the new keys.
	PGPKeys []string `json:"pgp_keys"`

	// Backup is a boolean indicating if Vault should backup the current keys in
	// case of rekeying failure.
	Backup bool
}

// RekeyStatusResponse is the response from a rekey operation.
type RekeyStatusResponse struct {
	// Nonce is the operation nonce.
	Nonce string

	// Started is a boolean indicating if the rekeying operation is started.
	Started bool

	// T is the new threshold.
	T int

	// N is the new number of shares.
	N int

	// Progress is the current rekey progress, if any.
	Progress int

	// Required is the number of required shards to complete the rekey.
	Required int

	// PGPFingerprints is the list of PGP fingerprints, if given.
	PGPFingerprints []string `json:"pgp_fingerprints"`

	// Backup is a bool indicating whether the rekey operation is backed up.
	Backup bool
}

// RekeyUpdateResponse is the response when a user submits an unseal key.
type RekeyUpdateResponse struct {
	// Nonce is the operation nonce.
	Nonce string

	// Complete is a boolean indicating if the rekey is complete.
	Complete bool

	// Keys is the list of unseal keys hex-encoded.
	Keys []string

	// KeysB64 is the list of unseal keys base64-encoded.
	KeysB64 []string `json:"keys_base64"`

	// PGPFingerprints is the list of PGP fingerprints, if given.
	PGPFingerprints []string `json:"pgp_fingerprints"`

	// Backup is a bool indicating whether the rekey operation is backed up.
	Backup bool
}

// RekeyRetrieveResponse is the response for getting the backup.
type RekeyRetrieveResponse struct {
	// Nonce is the operation nonce.
	Nonce string

	Keys    map[string][]string
	KeysB64 map[string][]string `json:"keys_base64"`
}

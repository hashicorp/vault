package api

import (
	"context"
	"net/http"
)

// InitStatus gets the current status of initialization.
func (c *Sys) InitStatus() (bool, error) {
	return c.InitStatusWithContext(context.Background())
}

// InitStatusWithContext gets the current status of initialization, with a
// context.
func (c *Sys) InitStatusWithContext(ctx context.Context) (bool, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/init")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result InitStatusResponse
	err = resp.DecodeJSON(&result)
	return result.Initialized, err
}

// Init initializes the Vault.
func (c *Sys) Init(opts *InitRequest) (*InitResponse, error) {
	return c.InitWithContext(context.Background(), opts)
}

// Init initializes the Vault, with a context.
func (c *Sys) InitWithContext(ctx context.Context, opts *InitRequest) (*InitResponse, error) {
	req := c.c.NewRequest(http.MethodPut, "/v1/sys/init")
	req = req.WithContext(ctx)
	if err := req.SetJSONBody(opts); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result InitResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// InitRequest is used as input to the Init and InitWithContext functions.
type InitRequest struct {
	// SecretShares is the number of keys to generate.
	SecretShares int `json:"secret_shares"`

	// SecretThreshold is the number of keys that must be supplied to re-generate
	// the master key. This must be <= SecretShares.
	SecretThreshold int `json:"secret_threshold"`

	// StoredShares is the number of shares that should be stored on the HSM.
	// Vault Enterprise only.
	StoredShares int `json:"stored_shares"`

	// PGPKeys is a list of PGP public keys or keybase usernames to encrypt the
	// generated key shares with. This must be equal to the number of shares.
	PGPKeys []string `json:"pgp_keys"`

	// RecoveryShares is the number of recovery shares to generate.
	RecoveryShares int `json:"recovery_shares"`

	// RecoveryThreshold is the number of keys that must be supplied to
	// re-generate the recovery key. This must be <= RecoveryShares.
	RecoveryThreshold int `json:"recovery_threshold"`

	// RecoveryPGPKeys is a list of PGP public keys or keybase usernames to
	// encrypt the generated key shares with. This must be equal to the number of
	// recovery shares.
	RecoveryPGPKeys []string `json:"recovery_pgp_keys"`

	// RootTokenPGPKey is the PGP public key or keybase username to encrypt the
	// initial root token with.
	RootTokenPGPKey string `json:"root_token_pgp_key"`
}

// InitStatusResponse is the response for an init status request.
type InitStatusResponse struct {
	// Initialized will be true if the Vault is currently initialized.
	Initialized bool
}

// InitResponse is the response for the init request. If any PGP options were
// given to the init request, the results will be encrypted with the associated
// PGP key.
type InitResponse struct {
	// Keys is the list of keys hex-encoded.
	Keys []string `json:"keys"`

	// KeysB64 is the list of keys base64-encoded.
	KeysB64 []string `json:"keys_base64"`

	// RecoveryKeys is the list of recovery keys hex-encoded.
	RecoveryKeys []string `json:"recovery_keys"`

	// RecoveryKeysB64 is the list of recovery keys base64-encoded.
	RecoveryKeysB64 []string `json:"recovery_keys_base64"`

	// RootToken is the initial root token.
	RootToken string `json:"root_token"`
}

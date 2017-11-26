package api

import (
	"context"
	"net/http"
	"time"
)

// Rotate triggers a rotation of Vault's underlying encryption key.
func (c *Sys) Rotate() error {
	return c.RotateWithContext(context.Background())
}

// RotateWithContext triggers a rotation of Vault's underlying encryption key,
// with a context.
func (c *Sys) RotateWithContext(ctx context.Context) error {
	req := c.c.NewRequest(http.MethodPost, "/v1/sys/rotate")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// KeyStatus returns the current key status of Vault.
func (c *Sys) KeyStatus() (*KeyStatus, error) {
	return c.KeyStatusWithContext(context.Background())
}

// KeyStatusWithContext returns the current key status of Vault, with a context.
func (c *Sys) KeyStatusWithContext(ctx context.Context) (*KeyStatus, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/key-status")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := new(KeyStatus)
	err = resp.DecodeJSON(result)
	return result, err
}

// KeyStatus is the result of a key status query.
type KeyStatus struct {
	// Term is the current key term.
	Term int `json:"term"`

	// InstallTime is the time where the key was installed.
	InstallTime time.Time `json:"install_time"`
}

package api

import (
	"context"
	"net/http"
)

// SealStatus returns the seal status for the Vault server the client is
// pointing to.
func (c *Sys) SealStatus() (*SealStatusResponse, error) {
	return c.SealStatusWithContext(context.Background())
}

// SealStatusWithContext returns the seal status for the Vault server the client
// is pointing to, with a context.
func (c *Sys) SealStatusWithContext(ctx context.Context) (*SealStatusResponse, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/seal-status")
	req = req.WithContext(ctx)
	return sealStatusRequest(c, req)
}

// Seal seals the Vault server.
func (c *Sys) Seal() error {
	return c.SealWithContext(context.Background())
}

// SealWithContext seals the Vault server, with a context.
func (c *Sys) SealWithContext(ctx context.Context) error {
	req := c.c.NewRequest(http.MethodPut, "/v1/sys/seal")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// ResetUnsealProcess resets an unseal process back to zero, discarding any
// entered shares.
func (c *Sys) ResetUnsealProcess() (*SealStatusResponse, error) {
	return c.ResetUnsealProcessWithContext(context.Background())
}

// ResetUnsealProcessWithContext resets an unseal process back to zero,
// discarding any entered shares, with a context.
func (c *Sys) ResetUnsealProcessWithContext(ctx context.Context) (*SealStatusResponse, error) {
	body := map[string]interface{}{"reset": true}

	req := c.c.NewRequest(http.MethodPut, "/v1/sys/unseal")
	req = req.WithContext(ctx)
	if err := req.SetJSONBody(body); err != nil {
		return nil, err
	}

	return sealStatusRequest(c, req)
}

// Unseal is used to supply one piece of the unseal key.
func (c *Sys) Unseal(shard string) (*SealStatusResponse, error) {
	return c.UnsealWithContext(context.Background(), shard)
}

// UnsealWithContext is used to supply one piece of the unseal key, with a
// context.
func (c *Sys) UnsealWithContext(ctx context.Context, shard string) (*SealStatusResponse, error) {
	body := map[string]interface{}{"key": shard}

	req := c.c.NewRequest(http.MethodPut, "/v1/sys/unseal")
	req = req.WithContext(ctx)
	if err := req.SetJSONBody(body); err != nil {
		return nil, err
	}

	return sealStatusRequest(c, req)
}

func sealStatusRequest(c *Sys, r *Request) (*SealStatusResponse, error) {
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SealStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// SealStatusResponse is the response returned from a seal status request.
type SealStatusResponse struct {
	// Type is the type of seal implored.
	Type string `json:"type"`

	// Sealed indicates if the Vault is sealed.
	Sealed bool `json:"sealed"`

	// T is the threshold of shares.
	T int `json:"t"`

	// N is the total number of shares.
	N int `json:"n"`

	// Progress is the current unseal progress, if any.
	Progress int `json:"progress"`

	// None is the current nonce operation.
	Nonce string `json:"nonce"`

	// Version is the Vault version.
	Version string `json:"version"`

	// ClusterName is the user-supplied or auto-generated cluster name.
	ClusterName string `json:"cluster_name,omitempty"`

	// ClusterID is the ID of the cluster.
	ClusterID string `json:"cluster_id,omitempty"`
}

package api

import (
	"context"
	"net/http"
)

// Health gets the current health of the cluster.
func (c *Sys) Health() (*HealthResponse, error) {
	return c.HealthWithContext(context.Background())
}

// HealthWithContext gets the current health of the cluster, with a context.
func (c *Sys) HealthWithContext(ctx context.Context) (*HealthResponse, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/health")
	req = req.WithContext(ctx)

	// If the code is 400 or above it will automatically turn into an error,
	// but the sys/health API defaults to returning 5xx when not sealed or
	// inited, so we force this code to be something else so we parse correctly
	req.Params.Add("sealedcode", "299")
	req.Params.Add("uninitcode", "299")
	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result HealthResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// HealthResponse is the result of the health query.
type HealthResponse struct {
	// Initialized returns true if the Vault is currently initialized.
	Initialized bool `json:"initialized"`

	// Sealed returns true if the Vault is currently sealed.
	Sealed bool `json:"sealed"`

	// Standby returns true if the Vault is currently in standby (if running in HA
	// mode).
	Standby bool `json:"standby"`

	// ServerTimeUTC is the current time of the server in UTC>
	ServerTimeUTC int64 `json:"server_time_utc"`

	// Version is the Vault version.
	Version string `json:"version"`

	// ClusterName is the user-supplied or auto-generated cluster name.
	ClusterName string `json:"cluster_name,omitempty"`

	// ClusterID is the ID of the cluster.
	ClusterID string `json:"cluster_id,omitempty"`
}

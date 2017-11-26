package api

import (
	"context"
	"net/http"
)

// Leader returns information about the current leader.
func (c *Sys) Leader() (*LeaderResponse, error) {
	return c.LeaderWithContext(context.Background())
}

// LeaderWithContext returns information about the current leader, with a
// context.
func (c *Sys) LeaderWithContext(ctx context.Context) (*LeaderResponse, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/leader")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result LeaderResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// LeaderResponse is the response from a request for leader information.
type LeaderResponse struct {
	// HAEnabled indicates whether high-availability is enabled.
	HAEnabled bool `json:"ha_enabled"`

	// IsSelf is a boolean which will be true if the leader is the Vault server
	// which was queried.
	IsSelf bool `json:"is_self"`

	// LeaderAddress is the address of the current cluster leader.
	LeaderAddress string `json:"leader_address"`

	// LeaderClusterAddress is the address to use when addressing the entire
	// cluster.
	LeaderClusterAddress string `json:"leader_cluster_address"`
}

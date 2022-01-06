package api

import (
	"context"
	"time"
)

func (c *Sys) HAStatus() (*HAStatusResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/ha-status")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result HAStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

type HAStatusResponse struct {
	Nodes []HANode
}

type HANode struct {
	Hostname       string     `json:"hostname"`
	APIAddress     string     `json:"api_address"`
	ClusterAddress string     `json:"cluster_address"`
	ActiveNode     bool       `json:"active_node"`
	LastEcho       *time.Time `json:"last_echo"`
}

package linodego

import (
	"context"
)

// NodeBalancerStats represents a nodebalancer stats object
type NodeBalancerStats struct {
	Title string                `json:"title"`
	Data  NodeBalancerStatsData `json:"data"`
}

// NodeBalancerStatsData represents a nodebalancer stats data object
type NodeBalancerStatsData struct {
	Connections [][]float64  `json:"connections"`
	Traffic     StatsTraffic `json:"traffic"`
}

// StatsTraffic represents a Traffic stats object
type StatsTraffic struct {
	In  [][]float64 `json:"in"`
	Out [][]float64 `json:"out"`
}

// GetNodeBalancerStats gets the template with the provided ID
func (c *Client) GetNodeBalancerStats(ctx context.Context, nodebalancerID int) (*NodeBalancerStats, error) {
	e := formatAPIPath("nodebalancers/%d/stats", nodebalancerID)
	response, err := doGETRequest[NodeBalancerStats](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

package linodego

import (
	"context"
)

// StatsNet represents a network stats object
type StatsNet struct {
	In         [][]float64 `json:"in"`
	Out        [][]float64 `json:"out"`
	PrivateIn  [][]float64 `json:"private_in"`
	PrivateOut [][]float64 `json:"private_out"`
}

// StatsIO represents an IO stats object
type StatsIO struct {
	IO   [][]float64 `json:"io"`
	Swap [][]float64 `json:"swap"`
}

// InstanceStatsData represents an instance stats data object
type InstanceStatsData struct {
	CPU   [][]float64 `json:"cpu"`
	IO    StatsIO     `json:"io"`
	NetV4 StatsNet    `json:"netv4"`
	NetV6 StatsNet    `json:"netv6"`
}

// InstanceStats represents an instance stats object
type InstanceStats struct {
	Title string            `json:"title"`
	Data  InstanceStatsData `json:"data"`
}

// GetInstanceStats gets the template with the provided ID
func (c *Client) GetInstanceStats(ctx context.Context, linodeID int) (*InstanceStats, error) {
	e := formatAPIPath("linode/instances/%d/stats", linodeID)
	response, err := doGETRequest[InstanceStats](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetInstanceStatsByDate gets the template with the provided ID, year, and month
func (c *Client) GetInstanceStatsByDate(ctx context.Context, linodeID int, year int, month int) (*InstanceStats, error) {
	e := formatAPIPath("linode/instances/%d/stats/%d/%d", linodeID, year, month)
	response, err := doGETRequest[InstanceStats](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

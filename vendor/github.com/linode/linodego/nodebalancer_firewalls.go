package linodego

import (
	"context"
)

// ListNodeBalancerFirewalls returns a paginated list of Cloud Firewalls for nodebalancerID
func (c *Client) ListNodeBalancerFirewalls(ctx context.Context, nodebalancerID int, opts *ListOptions) ([]Firewall, error) {
	response, err := getPaginatedResults[Firewall](ctx, c, formatAPIPath("nodebalancers/%d/firewalls", nodebalancerID), opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

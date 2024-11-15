package linodego

import (
	"context"
)

// ListIPv6Pools lists IPv6Pools
func (c *Client) ListIPv6Pools(ctx context.Context, opts *ListOptions) ([]IPv6Range, error) {
	response, err := getPaginatedResults[IPv6Range](ctx, c, "networking/ipv6/pools", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetIPv6Pool gets the template with the provided ID
func (c *Client) GetIPv6Pool(ctx context.Context, id string) (*IPv6Range, error) {
	e := formatAPIPath("networking/ipv6/pools/%s", id)
	response, err := doGETRequest[IPv6Range](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

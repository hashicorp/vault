package linodego

import (
	"context"
)

// ListInstanceVolumes lists InstanceVolumes
func (c *Client) ListInstanceVolumes(ctx context.Context, linodeID int, opts *ListOptions) ([]Volume, error) {
	response, err := getPaginatedResults[Volume](ctx, c, formatAPIPath("linode/instances/%d/volumes", linodeID), opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

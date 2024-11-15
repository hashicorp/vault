package linodego

import (
	"context"
)

type GrantsListResponse = UserGrants

func (c *Client) GrantsList(ctx context.Context) (*GrantsListResponse, error) {
	e := "profile/grants"
	response, err := doGETRequest[GrantsListResponse](ctx, c, e)
	return response, err
}

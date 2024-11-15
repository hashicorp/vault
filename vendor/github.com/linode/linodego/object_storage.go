package linodego

import (
	"context"
)

// ObjectStorageTransfer is an object matching the response of object-storage/transfer
type ObjectStorageTransfer struct {
	AmmountUsed int `json:"used"`
}

// CancelObjectStorage cancels and removes all object storage from the Account
func (c *Client) CancelObjectStorage(ctx context.Context) error {
	e := "object-storage/cancel"
	_, err := doPOSTRequest[any, any](ctx, c, e)
	return err
}

// GetObjectStorageTransfer returns the amount of outbound data transferred used by the Account
func (c *Client) GetObjectStorageTransfer(ctx context.Context) (*ObjectStorageTransfer, error) {
	e := "object-storage/transfer"
	response, err := doGETRequest[ObjectStorageTransfer](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

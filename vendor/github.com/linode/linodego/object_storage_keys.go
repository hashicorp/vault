package linodego

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/linode/linodego/pkg/errors"
)

// ObjectStorageKey represents a linode object storage key object
type ObjectStorageKey struct {
	ID        int    `json:"id"`
	Label     string `json:"label"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// ObjectStorageKeyCreateOptions fields are those accepted by CreateObjectStorageKey
type ObjectStorageKeyCreateOptions struct {
	Label string `json:"label"`
}

// ObjectStorageKeyUpdateOptions fields are those accepted by UpdateObjectStorageKey
type ObjectStorageKeyUpdateOptions struct {
	Label string `json:"label"`
}

// ObjectStorageKeysPagedResponse represents a linode API response for listing
type ObjectStorageKeysPagedResponse struct {
	*PageOptions
	Data []ObjectStorageKey `json:"data"`
}

// endpoint gets the endpoint URL for Object Storage keys
func (ObjectStorageKeysPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.ObjectStorageKeys.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends ObjectStorageKeys when processing paginated Objkey responses
func (resp *ObjectStorageKeysPagedResponse) appendData(r *ObjectStorageKeysPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListObjectStorageKeys lists ObjectStorageKeys
func (c *Client) ListObjectStorageKeys(ctx context.Context, opts *ListOptions) ([]ObjectStorageKey, error) {
	response := ObjectStorageKeysPagedResponse{}
	err := c.listHelper(ctx, &response, opts)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// CreateObjectStorageKey creates a ObjectStorageKey
func (c *Client) CreateObjectStorageKey(ctx context.Context, createOpts ObjectStorageKeyCreateOptions) (*ObjectStorageKey, error) {
	var body string
	e, err := c.ObjectStorageKeys.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&ObjectStorageKey{})

	if bodyData, err := json.Marshal(createOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, errors.New(err)
	}

	r, err := errors.CoupleAPIErrors(req.
		SetBody(body).
		Post(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*ObjectStorageKey), nil
}

// GetObjectStorageKey gets the object storage key with the provided ID
func (c *Client) GetObjectStorageKey(ctx context.Context, id int) (*ObjectStorageKey, error) {
	e, err := c.ObjectStorageKeys.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)
	r, err := errors.CoupleAPIErrors(c.R(ctx).SetResult(&ObjectStorageKey{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*ObjectStorageKey), nil
}

// UpdateObjectStorageKey updates the object storage key with the specified id
func (c *Client) UpdateObjectStorageKey(ctx context.Context, id int, updateOpts ObjectStorageKeyUpdateOptions) (*ObjectStorageKey, error) {
	var body string
	e, err := c.ObjectStorageKeys.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	req := c.R(ctx).SetResult(&ObjectStorageKey{})

	if bodyData, err := json.Marshal(updateOpts); err == nil {
		body = string(bodyData)
	} else {
		return nil, errors.New(err)
	}

	r, err := errors.CoupleAPIErrors(req.
		SetBody(body).
		Put(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*ObjectStorageKey), nil
}

// DeleteObjectStorageKey deletes the ObjectStorageKey with the specified id
func (c *Client) DeleteObjectStorageKey(ctx context.Context, id int) error {
	e, err := c.ObjectStorageKeys.Endpoint()
	if err != nil {
		return err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	_, err = errors.CoupleAPIErrors(c.R(ctx).Delete(e))
	return err
}

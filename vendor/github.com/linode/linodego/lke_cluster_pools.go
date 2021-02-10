package linodego

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/linode/linodego/pkg/errors"
)

// LKELinodeStatus constants start with LKELinode and include
// Linode API LKEClusterPool Linode Status values
type LKELinodeStatus string

// LKEClusterPoolStatus constants reflect the current status of an LKEClusterPool
const (
	LKELinodeReady    LKELinodeStatus = "ready"
	LKELinodeNotReady LKELinodeStatus = "not_ready"
)

// LKEClusterPoolLinode represents a LKEClusterPoolLinode object
type LKEClusterPoolLinode struct {
	ID         string          `json:"id"`
	InstanceID int             `json:"instance_id"`
	Status     LKELinodeStatus `json:"status"`
}

// LKEClusterPool represents a LKEClusterPool object
type LKEClusterPool struct {
	ID      int                    `json:"id"`
	Count   int                    `json:"count"`
	Type    string                 `json:"type"`
	Linodes []LKEClusterPoolLinode `json:"nodes"`
}

// LKEClusterPoolCreateOptions fields are those accepted by CreateLKEClusterPool
type LKEClusterPoolCreateOptions struct {
	Count int    `json:"count"`
	Type  string `json:"type"`
}

// LKEClusterPoolUpdateOptions fields are those accepted by UpdateLKEClusterPool
type LKEClusterPoolUpdateOptions struct {
	Count int `json:"count"`
}

// GetCreateOptions converts a LKEClusterPool to LKEClusterPoolCreateOptions for
// use in CreateLKEClusterPool
func (l LKEClusterPool) GetCreateOptions() (o LKEClusterPoolCreateOptions) {
	o.Count = l.Count
	return
}

// GetUpdateOptions converts a LKEClusterPool to LKEClusterPoolUpdateOptions for use in UpdateLKEClusterPool
func (l LKEClusterPool) GetUpdateOptions() (o LKEClusterPoolUpdateOptions) {
	o.Count = l.Count
	return
}

// LKEClusterPoolsPagedResponse represents a paginated LKEClusterPool API response
type LKEClusterPoolsPagedResponse struct {
	*PageOptions
	Data []LKEClusterPool `json:"data"`
}

// endpointWithID gets the endpoint URL for InstanceConfigs of a given Instance
func (LKEClusterPoolsPagedResponse) endpointWithID(c *Client, id int) string {
	endpoint, err := c.LKEClusterPools.endpointWithID(id)
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends LKEClusterPools when processing paginated LKEClusterPool responses
func (resp *LKEClusterPoolsPagedResponse) appendData(r *LKEClusterPoolsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListLKEClusterPools lists LKEClusterPools
func (c *Client) ListLKEClusterPools(ctx context.Context, clusterID int, opts *ListOptions) ([]LKEClusterPool, error) {
	response := LKEClusterPoolsPagedResponse{}
	err := c.listHelperWithID(ctx, &response, clusterID, opts)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// GetLKEClusterPool gets the lkeClusterPool with the provided ID
func (c *Client) GetLKEClusterPool(ctx context.Context, clusterID, id int) (*LKEClusterPool, error) {
	e, err := c.LKEClusterPools.endpointWithID(clusterID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)
	r, err := errors.CoupleAPIErrors(c.R(ctx).SetResult(&LKEClusterPool{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*LKEClusterPool), nil
}

// CreateLKEClusterPool creates a LKEClusterPool
func (c *Client) CreateLKEClusterPool(ctx context.Context, clusterID int, createOpts LKEClusterPoolCreateOptions) (*LKEClusterPool, error) {
	var body string
	e, err := c.LKEClusterPools.endpointWithID(clusterID)
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&LKEClusterPool{})

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
	return r.Result().(*LKEClusterPool), nil
}

// UpdateLKEClusterPool updates the LKEClusterPool with the specified id
func (c *Client) UpdateLKEClusterPool(ctx context.Context, clusterID, id int, updateOpts LKEClusterPoolUpdateOptions) (*LKEClusterPool, error) {
	var body string
	e, err := c.LKEClusterPools.endpointWithID(clusterID)
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	req := c.R(ctx).SetResult(&LKEClusterPool{})

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
	return r.Result().(*LKEClusterPool), nil
}

// DeleteLKEClusterPool deletes the LKEClusterPool with the specified id
func (c *Client) DeleteLKEClusterPool(ctx context.Context,
	clusterID, id int) error {
	e, err := c.LKEClusterPools.endpointWithID(clusterID)
	if err != nil {
		return err
	}
	e = fmt.Sprintf("%s/%d", e, id)

	_, err = errors.CoupleAPIErrors(c.R(ctx).Delete(e))
	return err
}

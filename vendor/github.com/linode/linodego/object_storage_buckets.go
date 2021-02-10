package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
	"github.com/linode/linodego/pkg/errors"
)

// ObjectStorageBucket represents a ObjectStorage object
type ObjectStorageBucket struct {
	Label   string `json:"label"`
	Cluster string `json:"cluster"`

	Created  *time.Time `json:"-"`
	Hostname string     `json:"hostname"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *ObjectStorageBucket) UnmarshalJSON(b []byte) error {
	type Mask ObjectStorageBucket

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)

	return nil
}

// ObjectStorageBucketCreateOptions fields are those accepted by CreateObjectStorageBucket
type ObjectStorageBucketCreateOptions struct {
	Cluster string `json:"cluster"`
	Label   string `json:"label"`
}

// ObjectStorageBucketsPagedResponse represents a paginated ObjectStorageBucket API response
type ObjectStorageBucketsPagedResponse struct {
	*PageOptions
	Data []ObjectStorageBucket `json:"data"`
}

// endpoint gets the endpoint URL for ObjectStorageBucket
func (ObjectStorageBucketsPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.ObjectStorageBuckets.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

// appendData appends ObjectStorageBuckets when processing paginated ObjectStorageBucket responses
func (resp *ObjectStorageBucketsPagedResponse) appendData(r *ObjectStorageBucketsPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListObjectStorageBuckets lists ObjectStorageBuckets
func (c *Client) ListObjectStorageBuckets(ctx context.Context, opts *ListOptions) ([]ObjectStorageBucket, error) {
	response := ObjectStorageBucketsPagedResponse{}
	err := c.listHelper(ctx, &response, opts)

	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetObjectStorageBucket gets the ObjectStorageBucket with the provided label
func (c *Client) GetObjectStorageBucket(ctx context.Context, clusterID, label string) (*ObjectStorageBucket, error) {
	e, err := c.ObjectStorageBuckets.Endpoint()
	if err != nil {
		return nil, err
	}
	e = fmt.Sprintf("%s/%s/%s", e, clusterID, label)
	r, err := errors.CoupleAPIErrors(c.R(ctx).SetResult(&ObjectStorageBucket{}).Get(e))
	if err != nil {
		return nil, err
	}
	return r.Result().(*ObjectStorageBucket), nil
}

// CreateObjectStorageBucket creates an ObjectStorageBucket
func (c *Client) CreateObjectStorageBucket(ctx context.Context, createOpts ObjectStorageBucketCreateOptions) (*ObjectStorageBucket, error) {
	var body string
	e, err := c.ObjectStorageBuckets.Endpoint()
	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&ObjectStorageBucket{})

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
	return r.Result().(*ObjectStorageBucket), nil
}

// DeleteObjectStorageBucket deletes the ObjectStorageBucket with the specified label
func (c *Client) DeleteObjectStorageBucket(ctx context.Context, clusterID, label string) error {
	e, err := c.ObjectStorageBuckets.Endpoint()
	if err != nil {
		return err
	}
	e = fmt.Sprintf("%s/%s/%s", e, clusterID, label)

	_, err = errors.CoupleAPIErrors(c.R(ctx).Delete(e))
	return err
}

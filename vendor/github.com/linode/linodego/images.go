package linodego

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
	"github.com/linode/linodego/pkg/errors"
)

// Image represents a deployable Image object for use with Linode Instances
type Image struct {
	ID          string     `json:"id"`
	CreatedBy   string     `json:"created_by"`
	Label       string     `json:"label"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	Vendor      string     `json:"vendor"`
	Size        int        `json:"size"`
	IsPublic    bool       `json:"is_public"`
	Deprecated  bool       `json:"deprecated"`
	Created     *time.Time `json:"-"`
	Expiry      *time.Time `json:"-"`
}

// ImageCreateOptions fields are those accepted by CreateImage
type ImageCreateOptions struct {
	DiskID      int    `json:"disk_id"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

// ImageUpdateOptions fields are those accepted by UpdateImage
type ImageUpdateOptions struct {
	Label       string  `json:"label,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *Image) UnmarshalJSON(b []byte) error {
	type Mask Image

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
		Expiry  *parseabletime.ParseableTime `json:"expiry"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)
	i.Expiry = (*time.Time)(p.Expiry)

	return nil
}

// GetUpdateOptions converts an Image to ImageUpdateOptions for use in UpdateImage
func (i Image) GetUpdateOptions() (iu ImageUpdateOptions) {
	iu.Label = i.Label
	iu.Description = copyString(&i.Description)
	return
}

// ImagesPagedResponse represents a linode API response for listing of images
type ImagesPagedResponse struct {
	*PageOptions
	Data []Image `json:"data"`
}

func (ImagesPagedResponse) endpoint(c *Client) string {
	endpoint, err := c.Images.Endpoint()
	if err != nil {
		panic(err)
	}
	return endpoint
}

func (resp *ImagesPagedResponse) appendData(r *ImagesPagedResponse) {
	resp.Data = append(resp.Data, r.Data...)
}

// ListImages lists Images
func (c *Client) ListImages(ctx context.Context, opts *ListOptions) ([]Image, error) {
	response := ImagesPagedResponse{}
	err := c.listHelper(ctx, &response, opts)

	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// GetImage gets the Image with the provided ID
func (c *Client) GetImage(ctx context.Context, id string) (*Image, error) {
	e, err := c.Images.Endpoint()
	if err != nil {
		return nil, err
	}

	e = fmt.Sprintf("%s/%s", e, id)
	r, err := errors.CoupleAPIErrors(c.Images.R(ctx).Get(e))

	if err != nil {
		return nil, err
	}
	return r.Result().(*Image), nil
}

// CreateImage creates a Image
func (c *Client) CreateImage(ctx context.Context, createOpts ImageCreateOptions) (*Image, error) {
	var body string

	e, err := c.Images.Endpoint()

	if err != nil {
		return nil, err
	}

	req := c.R(ctx).SetResult(&Image{})

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
	return r.Result().(*Image), nil
}

// UpdateImage updates the Image with the specified id
func (c *Client) UpdateImage(ctx context.Context, id string, updateOpts ImageUpdateOptions) (*Image, error) {
	var body string

	e, err := c.Images.Endpoint()
	if err != nil {
		return nil, err
	}

	e = fmt.Sprintf("%s/%s", e, id)

	req := c.R(ctx).SetResult(&Image{})

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
	return r.Result().(*Image), nil
}

// DeleteImage deletes the Image with the specified id
func (c *Client) DeleteImage(ctx context.Context, id string) error {
	e, err := c.Images.Endpoint()
	if err != nil {
		return err
	}

	e = fmt.Sprintf("%s/%s", e, id)

	_, err = errors.CoupleAPIErrors(c.R(ctx).Delete(e))
	return err
}

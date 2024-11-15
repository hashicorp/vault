package linodego

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/linode/linodego/internal/parseabletime"
)

// ImageStatus represents the status of an Image.
type ImageStatus string

// ImageStatus options start with ImageStatus and include all Image statuses
const (
	ImageStatusCreating      ImageStatus = "creating"
	ImageStatusPendingUpload ImageStatus = "pending_upload"
	ImageStatusAvailable     ImageStatus = "available"
)

// ImageRegionStatus represents the status of an Image's replica.
type ImageRegionStatus string

// ImageRegionStatus options start with ImageRegionStatus and
// include all Image replica statuses
const (
	ImageRegionStatusAvailable          ImageRegionStatus = "available"
	ImageRegionStatusCreating           ImageRegionStatus = "creating"
	ImageRegionStatusPending            ImageRegionStatus = "pending"
	ImageRegionStatusPendingReplication ImageRegionStatus = "pending replication"
	ImageRegionStatusPendingDeletion    ImageRegionStatus = "pending deletion"
	ImageRegionStatusReplicating        ImageRegionStatus = "replicating"
)

// ImageRegion represents the status of an Image object in a given Region.
type ImageRegion struct {
	Region string            `json:"region"`
	Status ImageRegionStatus `json:"status"`
}

// Image represents a deployable Image object for use with Linode Instances
type Image struct {
	ID           string        `json:"id"`
	CreatedBy    string        `json:"created_by"`
	Capabilities []string      `json:"capabilities"`
	Label        string        `json:"label"`
	Description  string        `json:"description"`
	Type         string        `json:"type"`
	Vendor       string        `json:"vendor"`
	Status       ImageStatus   `json:"status"`
	Size         int           `json:"size"`
	TotalSize    int           `json:"total_size"`
	IsPublic     bool          `json:"is_public"`
	Deprecated   bool          `json:"deprecated"`
	Regions      []ImageRegion `json:"regions"`
	Tags         []string      `json:"tags"`

	Updated *time.Time `json:"-"`
	Created *time.Time `json:"-"`
	Expiry  *time.Time `json:"-"`
	EOL     *time.Time `json:"-"`
}

// ImageCreateOptions fields are those accepted by CreateImage
type ImageCreateOptions struct {
	DiskID      int       `json:"disk_id"`
	Label       string    `json:"label"`
	Description string    `json:"description,omitempty"`
	CloudInit   bool      `json:"cloud_init,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
}

// ImageUpdateOptions fields are those accepted by UpdateImage
type ImageUpdateOptions struct {
	Label       string    `json:"label,omitempty"`
	Description *string   `json:"description,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
}

// ImageReplicateOptions represents the options accepted by the
// ReplicateImage(...) function.
type ImageReplicateOptions struct {
	Regions []string `json:"regions"`
}

// ImageCreateUploadResponse fields are those returned by CreateImageUpload
type ImageCreateUploadResponse struct {
	Image    *Image `json:"image"`
	UploadTo string `json:"upload_to"`
}

// ImageCreateUploadOptions fields are those accepted by CreateImageUpload
type ImageCreateUploadOptions struct {
	Region      string    `json:"region"`
	Label       string    `json:"label"`
	Description string    `json:"description,omitempty"`
	CloudInit   bool      `json:"cloud_init,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
}

// ImageUploadOptions fields are those accepted by UploadImage
type ImageUploadOptions struct {
	Region      string    `json:"region"`
	Label       string    `json:"label"`
	Description string    `json:"description,omitempty"`
	CloudInit   bool      `json:"cloud_init"`
	Tags        *[]string `json:"tags,omitempty"`
	Image       io.Reader
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *Image) UnmarshalJSON(b []byte) error {
	type Mask Image

	p := struct {
		*Mask
		Updated *parseabletime.ParseableTime `json:"updated"`
		Created *parseabletime.ParseableTime `json:"created"`
		Expiry  *parseabletime.ParseableTime `json:"expiry"`
		EOL     *parseabletime.ParseableTime `json:"eol"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Updated = (*time.Time)(p.Updated)
	i.Created = (*time.Time)(p.Created)
	i.Expiry = (*time.Time)(p.Expiry)
	i.EOL = (*time.Time)(p.EOL)

	return nil
}

// GetUpdateOptions converts an Image to ImageUpdateOptions for use in UpdateImage
func (i Image) GetUpdateOptions() (iu ImageUpdateOptions) {
	iu.Label = i.Label
	iu.Description = copyString(&i.Description)
	return
}

// ListImages lists Images.
func (c *Client) ListImages(ctx context.Context, opts *ListOptions) ([]Image, error) {
	return getPaginatedResults[Image](
		ctx,
		c,
		"images",
		opts,
	)
}

// GetImage gets the Image with the provided ID.
func (c *Client) GetImage(ctx context.Context, imageID string) (*Image, error) {
	return doGETRequest[Image](
		ctx,
		c,
		formatAPIPath("images/%s", imageID),
	)
}

// CreateImage creates an Image.
func (c *Client) CreateImage(ctx context.Context, opts ImageCreateOptions) (*Image, error) {
	return doPOSTRequest[Image](
		ctx,
		c,
		"images",
		opts,
	)
}

// UpdateImage updates the Image with the specified id.
func (c *Client) UpdateImage(ctx context.Context, imageID string, opts ImageUpdateOptions) (*Image, error) {
	return doPUTRequest[Image](
		ctx,
		c,
		formatAPIPath("images/%s", imageID),
		opts,
	)
}

// ReplicateImage replicates an image to a given set of regions.
// NOTE: Image replication may not currently be available to all users.
func (c *Client) ReplicateImage(ctx context.Context, imageID string, opts ImageReplicateOptions) (*Image, error) {
	return doPOSTRequest[Image](
		ctx,
		c,
		formatAPIPath("images/%s/regions", imageID),
		opts,
	)
}

// DeleteImage deletes the Image with the specified id.
func (c *Client) DeleteImage(ctx context.Context, imageID string) error {
	return doDELETERequest(
		ctx,
		c,
		formatAPIPath("images/%s", imageID),
	)
}

// CreateImageUpload creates an Image and an upload URL.
func (c *Client) CreateImageUpload(ctx context.Context, opts ImageCreateUploadOptions) (*Image, string, error) {
	result, err := doPOSTRequest[ImageCreateUploadResponse](
		ctx,
		c,
		"images/upload",
		opts,
	)
	if err != nil {
		return nil, "", err
	}

	return result.Image, result.UploadTo, nil
}

// UploadImageToURL uploads the given image to the given upload URL.
func (c *Client) UploadImageToURL(ctx context.Context, uploadURL string, image io.Reader) error {
	// Linode-specific headers do not need to be sent to this endpoint
	req := resty.New().SetDebug(c.resty.Debug).R().
		SetContext(ctx).
		SetContentLength(true).
		SetHeader("Content-Type", "application/octet-stream").
		SetBody(image)

	_, err := coupleAPIErrors(req.
		Put(uploadURL))

	return err
}

// UploadImage creates and uploads an image.
func (c *Client) UploadImage(ctx context.Context, opts ImageUploadOptions) (*Image, error) {
	image, uploadURL, err := c.CreateImageUpload(ctx, ImageCreateUploadOptions{
		Label:       opts.Label,
		Region:      opts.Region,
		Description: opts.Description,
		CloudInit:   opts.CloudInit,
		Tags:        opts.Tags,
	})
	if err != nil {
		return nil, err
	}

	return image, c.UploadImageToURL(ctx, uploadURL, opts.Image)
}

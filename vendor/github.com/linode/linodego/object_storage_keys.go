package linodego

import (
	"context"
)

type ObjectStorageKeyRegion struct {
	ID         string `json:"id"`
	S3Endpoint string `json:"s3_endpoint"`
}

// ObjectStorageKey represents a linode object storage key object
type ObjectStorageKey struct {
	ID           int                             `json:"id"`
	Label        string                          `json:"label"`
	AccessKey    string                          `json:"access_key"`
	SecretKey    string                          `json:"secret_key"`
	Limited      bool                            `json:"limited"`
	BucketAccess *[]ObjectStorageKeyBucketAccess `json:"bucket_access"`
	Regions      []ObjectStorageKeyRegion        `json:"regions"`
}

// ObjectStorageKeyBucketAccess represents a linode limited object storage key's bucket access
type ObjectStorageKeyBucketAccess struct {
	// Deprecated: Cluster field has been deprecated.
	// Please consider switching to use the 'Region' field.
	// If your Cluster is `us-mia-1`, then the region would be `us-mia`.
	Cluster string `json:"cluster,omitempty"`
	Region  string `json:"region,omitempty"`

	BucketName  string `json:"bucket_name"`
	Permissions string `json:"permissions"`
}

// ObjectStorageKeyCreateOptions fields are those accepted by CreateObjectStorageKey
type ObjectStorageKeyCreateOptions struct {
	Label        string                          `json:"label"`
	BucketAccess *[]ObjectStorageKeyBucketAccess `json:"bucket_access,omitempty"`
	Regions      []string                        `json:"regions,omitempty"`
}

// ObjectStorageKeyUpdateOptions fields are those accepted by UpdateObjectStorageKey
type ObjectStorageKeyUpdateOptions struct {
	Label   string   `json:"label,omitempty"`
	Regions []string `json:"regions,omitempty"`
}

// ListObjectStorageKeys lists ObjectStorageKeys
func (c *Client) ListObjectStorageKeys(ctx context.Context, opts *ListOptions) ([]ObjectStorageKey, error) {
	response, err := getPaginatedResults[ObjectStorageKey](ctx, c, "object-storage/keys", opts)
	return response, err
}

// CreateObjectStorageKey creates a ObjectStorageKey
func (c *Client) CreateObjectStorageKey(ctx context.Context, opts ObjectStorageKeyCreateOptions) (*ObjectStorageKey, error) {
	e := "object-storage/keys"
	response, err := doPOSTRequest[ObjectStorageKey](ctx, c, e, opts)
	return response, err
}

// GetObjectStorageKey gets the object storage key with the provided ID
func (c *Client) GetObjectStorageKey(ctx context.Context, keyID int) (*ObjectStorageKey, error) {
	e := formatAPIPath("object-storage/keys/%d", keyID)
	response, err := doGETRequest[ObjectStorageKey](ctx, c, e)
	return response, err
}

// UpdateObjectStorageKey updates the object storage key with the specified id
func (c *Client) UpdateObjectStorageKey(ctx context.Context, keyID int, opts ObjectStorageKeyUpdateOptions) (*ObjectStorageKey, error) {
	e := formatAPIPath("object-storage/keys/%d", keyID)
	response, err := doPUTRequest[ObjectStorageKey](ctx, c, e, opts)
	return response, err
}

// DeleteObjectStorageKey deletes the ObjectStorageKey with the specified id
func (c *Client) DeleteObjectStorageKey(ctx context.Context, keyID int) error {
	e := formatAPIPath("object-storage/keys/%d", keyID)
	err := doDELETERequest(ctx, c, e)
	return err
}

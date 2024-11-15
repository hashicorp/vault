package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// ObjectStorageBucket represents a ObjectStorage object
type ObjectStorageBucket struct {
	Label string `json:"label"`

	// Deprecated: The 'Cluster' field has been deprecated in favor of the 'Region' field.
	// For example, a Cluster value of `us-mia-1` will translate to a Region value of `us-mia`.
	//
	// This is necessary because there are now multiple Object Storage clusters to a region.
	//
	// NOTE: The 'Cluster' field will always return a value similar to `<REGION>-1` (e.g., `us-mia-1`)
	// for backward compatibility purposes.
	Cluster string `json:"cluster"`
	Region  string `json:"region"`

	Created  *time.Time `json:"-"`
	Hostname string     `json:"hostname"`
	Objects  int        `json:"objects"`
	Size     int        `json:"size"`
}

// ObjectStorageBucketAccess holds Object Storage access info
type ObjectStorageBucketAccess struct {
	ACL         ObjectStorageACL `json:"acl"`
	CorsEnabled bool             `json:"cors_enabled"`
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
	// Deprecated: The 'Cluster' field has been deprecated.
	//
	// Going forward, the 'Region' field will be the supported way to designate where an
	// Object Storage Bucket should be created. For example, a 'Cluster' value of `us-mia-1`
	// will translate to a Region value of `us-mia`.
	Cluster string `json:"cluster,omitempty"`
	Region  string `json:"region,omitempty"`

	Label string `json:"label"`

	ACL         ObjectStorageACL `json:"acl,omitempty"`
	CorsEnabled *bool            `json:"cors_enabled,omitempty"`
}

// ObjectStorageBucketUpdateAccessOptions fields are those accepted by UpdateObjectStorageBucketAccess
type ObjectStorageBucketUpdateAccessOptions struct {
	ACL         ObjectStorageACL `json:"acl,omitempty"`
	CorsEnabled *bool            `json:"cors_enabled,omitempty"`
}

// ObjectStorageACL options start with ACL and include all known ACL types
type ObjectStorageACL string

// ObjectStorageACL options represent the access control level of a bucket.
const (
	ACLPrivate           ObjectStorageACL = "private"
	ACLPublicRead        ObjectStorageACL = "public-read"
	ACLAuthenticatedRead ObjectStorageACL = "authenticated-read"
	ACLPublicReadWrite   ObjectStorageACL = "public-read-write"
)

// ListObjectStorageBuckets lists ObjectStorageBuckets
func (c *Client) ListObjectStorageBuckets(ctx context.Context, opts *ListOptions) ([]ObjectStorageBucket, error) {
	response, err := getPaginatedResults[ObjectStorageBucket](ctx, c, "object-storage/buckets", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ListObjectStorageBucketsInCluster lists all ObjectStorageBuckets of a cluster
func (c *Client) ListObjectStorageBucketsInCluster(ctx context.Context, opts *ListOptions, clusterOrRegionID string) ([]ObjectStorageBucket, error) {
	response, err := getPaginatedResults[ObjectStorageBucket](ctx, c, formatAPIPath("object-storage/buckets/%s", clusterOrRegionID), opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetObjectStorageBucket gets the ObjectStorageBucket with the provided label
func (c *Client) GetObjectStorageBucket(ctx context.Context, clusterOrRegionID, label string) (*ObjectStorageBucket, error) {
	e := formatAPIPath("object-storage/buckets/%s/%s", clusterOrRegionID, label)
	response, err := doGETRequest[ObjectStorageBucket](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateObjectStorageBucket creates an ObjectStorageBucket
func (c *Client) CreateObjectStorageBucket(ctx context.Context, opts ObjectStorageBucketCreateOptions) (*ObjectStorageBucket, error) {
	e := "object-storage/buckets"
	response, err := doPOSTRequest[ObjectStorageBucket](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetObjectStorageBucketAccess gets the current access config for a bucket
func (c *Client) GetObjectStorageBucketAccess(ctx context.Context, clusterOrRegionID, label string) (*ObjectStorageBucketAccess, error) {
	e := formatAPIPath("object-storage/buckets/%s/%s/access", clusterOrRegionID, label)
	response, err := doGETRequest[ObjectStorageBucketAccess](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateObjectStorageBucketAccess updates the access configuration for an ObjectStorageBucket
func (c *Client) UpdateObjectStorageBucketAccess(ctx context.Context, clusterOrRegionID, label string, opts ObjectStorageBucketUpdateAccessOptions) error {
	e := formatAPIPath("object-storage/buckets/%s/%s/access", clusterOrRegionID, label)
	_, err := doPOSTRequest[ObjectStorageBucketAccess](ctx, c, e, opts)

	return err
}

// DeleteObjectStorageBucket deletes the ObjectStorageBucket with the specified label
func (c *Client) DeleteObjectStorageBucket(ctx context.Context, clusterOrRegionID, label string) error {
	e := formatAPIPath("object-storage/buckets/%s/%s", clusterOrRegionID, label)
	err := doDELETERequest(ctx, c, e)
	return err
}

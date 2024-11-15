package linodego

import (
	"context"
)

type ObjectStorageObjectURLCreateOptions struct {
	Name               string `json:"name"`
	Method             string `json:"method"`
	ContentType        string `json:"content_type,omitempty"`
	ContentDisposition string `json:"content_disposition,omitempty"`
	ExpiresIn          *int   `json:"expires_in,omitempty"`
}

type ObjectStorageObjectURL struct {
	URL    string `json:"url"`
	Exists bool   `json:"exists"`
}

type ObjectStorageObjectACLConfig struct {
	ACL    string `json:"acl"`
	ACLXML string `json:"acl_xml"`
}

type ObjectStorageObjectACLConfigUpdateOptions struct {
	Name string `json:"name"`
	ACL  string `json:"acl"`
}

func (c *Client) CreateObjectStorageObjectURL(ctx context.Context, objectID, label string, opts ObjectStorageObjectURLCreateOptions) (*ObjectStorageObjectURL, error) {
	e := formatAPIPath("object-storage/buckets/%s/%s/object-url", objectID, label)
	response, err := doPOSTRequest[ObjectStorageObjectURL](ctx, c, e, opts)
	return response, err
}

func (c *Client) GetObjectStorageObjectACLConfig(ctx context.Context, objectID, label, object string) (*ObjectStorageObjectACLConfig, error) {
	e := formatAPIPath("object-storage/buckets/%s/%s/object-acl?name=%s", objectID, label, object)
	response, err := doGETRequest[ObjectStorageObjectACLConfig](ctx, c, e)
	return response, err
}

func (c *Client) UpdateObjectStorageObjectACLConfig(ctx context.Context, objectID, label string, opts ObjectStorageObjectACLConfigUpdateOptions) (*ObjectStorageObjectACLConfig, error) {
	e := formatAPIPath("object-storage/buckets/%s/%s/object-acl", objectID, label)
	response, err := doPUTRequest[ObjectStorageObjectACLConfig](ctx, c, e, opts)
	return response, err
}

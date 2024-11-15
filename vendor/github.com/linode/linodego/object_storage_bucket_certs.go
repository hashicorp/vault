package linodego

import (
	"context"
)

type ObjectStorageBucketCert struct {
	SSL bool `json:"ssl"`
}

type ObjectStorageBucketCertUploadOptions struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"private_key"`
}

// UploadObjectStorageBucketCert uploads a TLS/SSL Cert to be used with an Object Storage Bucket.
func (c *Client) UploadObjectStorageBucketCert(ctx context.Context, clusterOrRegionID, bucket string, opts ObjectStorageBucketCertUploadOptions) (*ObjectStorageBucketCert, error) {
	e := formatAPIPath("object-storage/buckets/%s/%s/ssl", clusterOrRegionID, bucket)
	response, err := doPOSTRequest[ObjectStorageBucketCert](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetObjectStorageBucketCert gets an ObjectStorageBucketCert
func (c *Client) GetObjectStorageBucketCert(ctx context.Context, clusterOrRegionID, bucket string) (*ObjectStorageBucketCert, error) {
	e := formatAPIPath("object-storage/buckets/%s/%s/ssl", clusterOrRegionID, bucket)
	response, err := doGETRequest[ObjectStorageBucketCert](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteObjectStorageBucketCert deletes an ObjectStorageBucketCert
func (c *Client) DeleteObjectStorageBucketCert(ctx context.Context, clusterOrRegionID, bucket string) error {
	e := formatAPIPath("object-storage/buckets/%s/%s/ssl", clusterOrRegionID, bucket)
	err := doDELETERequest(ctx, c, e)
	return err
}

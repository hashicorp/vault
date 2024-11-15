package gocb

type bucketManagementProvider interface {
	GetBucket(bucketName string, opts *GetBucketOptions) (*BucketSettings, error)
	GetAllBuckets(opts *GetAllBucketsOptions) (map[string]BucketSettings, error)
	CreateBucket(settings CreateBucketSettings, opts *CreateBucketOptions) error
	UpdateBucket(settings BucketSettings, opts *UpdateBucketOptions) error
	DropBucket(name string, opts *DropBucketOptions) error
	FlushBucket(name string, opts *FlushBucketOptions) error
}

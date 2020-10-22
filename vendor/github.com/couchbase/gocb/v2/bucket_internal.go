package gocb

import gocbcore "github.com/couchbase/gocbcore/v9"

// InternalBucket is used for internal functionality.
// Internal: This should never be used and is not supported.
type InternalBucket struct {
	bucket *Bucket
}

// Internal returns a CollectionInternal.
// Internal: This should never be used and is not supported.
func (b *Bucket) Internal() *InternalBucket {
	return &InternalBucket{bucket: b}
}

// IORouter returns the collection's internal core router.
func (ib *InternalBucket) IORouter() (*gocbcore.Agent, error) {
	return ib.bucket.connectionManager.connection(ib.bucket.Name())
}

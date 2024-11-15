package gocb

import "github.com/couchbase/gocbcore/v10"

type kvCapabilityVerifier interface {
	BucketCapabilityStatus(cap gocbcore.BucketCapability) gocbcore.CapabilityStatus
}

// InternalBucket is used for internal functionality.
// Internal: This should never be used and is not supported.
type InternalBucket struct {
	bucket *Bucket
}

// Internal returns an InternalBucket.
// Internal: This should never be used and is not supported.
func (b *Bucket) Internal() *InternalBucket {
	return &InternalBucket{bucket: b}
}

// IORouter returns the collection's internal core router.
func (ib *InternalBucket) IORouter() (*gocbcore.Agent, error) {
	return ib.bucket.connectionManager.connection(ib.bucket.Name())
}

// HasCapabilityStatus verifies whether support for a server capability is in a given state.
func (ib *InternalBucket) CapabilityStatus(cap Capability) (CapabilityStatus, error) {
	switch cap {
	case CapabilityCreateAsDeleted:
		return ib.bucketCapabilityStatus(gocbcore.BucketCapabilityCreateAsDeleted)
	case CapabilityDurableWrites:
		return ib.bucketCapabilityStatus(gocbcore.BucketCapabilityDurableWrites)
	case CapabilityReplaceBodyWithXattr:
		return ib.bucketCapabilityStatus(gocbcore.BucketCapabilityReplaceBodyWithXattr)
	default:
		return CapabilityStatusUnsupported, nil
	}
}

func (ib *InternalBucket) bucketCapabilityStatus(capability gocbcore.BucketCapability) (CapabilityStatus, error) {
	provider, err := ib.bucket.getKvCapabilitiesProvider()
	if err != nil {
		return 0, err
	}

	return CapabilityStatus(provider.BucketCapabilityStatus(capability)), nil
}

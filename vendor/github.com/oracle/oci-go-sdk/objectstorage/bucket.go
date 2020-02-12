// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object Storage and Archive Storage APIs for managing buckets, objects, and related resources.
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// Bucket A bucket is a container for storing objects in a compartment within a namespace. A bucket is associated with a single compartment.
// The compartment has policies that indicate what actions a user can perform on a bucket and all the objects in the bucket. For more
// information, see Managing Buckets (https://docs.cloud.oracle.com/Content/Object/Tasks/managingbuckets.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type Bucket struct {

	// The Object Storage namespace in which the bucket lives.
	Namespace *string `mandatory:"true" json:"namespace"`

	// The name of the bucket. Avoid entering confidential information.
	// Example: my-new-bucket1
	Name *string `mandatory:"true" json:"name"`

	// The compartment ID in which the bucket is authorized.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// Arbitrary string keys and values for user-defined metadata.
	Metadata map[string]string `mandatory:"true" json:"metadata"`

	// The OCID of the user who created the bucket.
	CreatedBy *string `mandatory:"true" json:"createdBy"`

	// The date and time the bucket was created, as described in RFC 2616 (https://tools.ietf.org/rfc/rfc2616), section 14.29.
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The entity tag (ETag) for the bucket.
	Etag *string `mandatory:"true" json:"etag"`

	// The type of public access enabled on this bucket.
	// A bucket is set to `NoPublicAccess` by default, which only allows an authenticated caller to access the
	// bucket and its contents. When `ObjectRead` is enabled on the bucket, public access is allowed for the
	// `GetObject`, `HeadObject`, and `ListObjects` operations. When `ObjectReadWithoutList` is enabled on the
	// bucket, public access is allowed for the `GetObject` and `HeadObject` operations.
	PublicAccessType BucketPublicAccessTypeEnum `mandatory:"false" json:"publicAccessType,omitempty"`

	// The storage tier type assigned to the bucket. A bucket is set to 'Standard' tier by default, which means
	// objects uploaded or copied to the bucket will be in the standard storage tier. When the 'Archive' tier type
	// is set explicitly for a bucket, objects uploaded or copied to the bucket will be stored in archive storage.
	// The 'storageTier' property is immutable after bucket is created.
	StorageTier BucketStorageTierEnum `mandatory:"false" json:"storageTier,omitempty"`

	// A property that determines whether events will be generated for operations on objects in this bucket.
	// This is false by default.
	ObjectEventsEnabled *bool `mandatory:"false" json:"objectEventsEnabled"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The OCID of a KMS key id used to call KMS to generate data key or decrypt the encrypted data key.
	KmsKeyId *string `mandatory:"false" json:"kmsKeyId"`

	// The entity tag (ETag) for the live object lifecycle policy on the bucket.
	ObjectLifecyclePolicyEtag *string `mandatory:"false" json:"objectLifecyclePolicyEtag"`

	// The approximate number of objects in the bucket. Count statistics are reported periodically. You will see a
	// lag between what is displayed and the actual object count.
	ApproximateCount *int64 `mandatory:"false" json:"approximateCount"`

	// The approximate total size in bytes of all objects in the bucket. Size statistics are reported periodically. You will
	// see a lag between what is displayed and the actual size of the bucket.
	ApproximateSize *int64 `mandatory:"false" json:"approximateSize"`
}

func (m Bucket) String() string {
	return common.PointerString(m)
}

// BucketPublicAccessTypeEnum Enum with underlying type: string
type BucketPublicAccessTypeEnum string

// Set of constants representing the allowable values for BucketPublicAccessTypeEnum
const (
	BucketPublicAccessTypeNopublicaccess        BucketPublicAccessTypeEnum = "NoPublicAccess"
	BucketPublicAccessTypeObjectread            BucketPublicAccessTypeEnum = "ObjectRead"
	BucketPublicAccessTypeObjectreadwithoutlist BucketPublicAccessTypeEnum = "ObjectReadWithoutList"
)

var mappingBucketPublicAccessType = map[string]BucketPublicAccessTypeEnum{
	"NoPublicAccess":        BucketPublicAccessTypeNopublicaccess,
	"ObjectRead":            BucketPublicAccessTypeObjectread,
	"ObjectReadWithoutList": BucketPublicAccessTypeObjectreadwithoutlist,
}

// GetBucketPublicAccessTypeEnumValues Enumerates the set of values for BucketPublicAccessTypeEnum
func GetBucketPublicAccessTypeEnumValues() []BucketPublicAccessTypeEnum {
	values := make([]BucketPublicAccessTypeEnum, 0)
	for _, v := range mappingBucketPublicAccessType {
		values = append(values, v)
	}
	return values
}

// BucketStorageTierEnum Enum with underlying type: string
type BucketStorageTierEnum string

// Set of constants representing the allowable values for BucketStorageTierEnum
const (
	BucketStorageTierStandard BucketStorageTierEnum = "Standard"
	BucketStorageTierArchive  BucketStorageTierEnum = "Archive"
)

var mappingBucketStorageTier = map[string]BucketStorageTierEnum{
	"Standard": BucketStorageTierStandard,
	"Archive":  BucketStorageTierArchive,
}

// GetBucketStorageTierEnumValues Enumerates the set of values for BucketStorageTierEnum
func GetBucketStorageTierEnumValues() []BucketStorageTierEnum {
	values := make([]BucketStorageTierEnum, 0)
	for _, v := range mappingBucketStorageTier {
		values = append(values, v)
	}
	return values
}

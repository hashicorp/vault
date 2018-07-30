// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object and Archive Storage APIs for managing buckets and objects.
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// CreateBucketDetails To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.us-phoenix-1.oraclecloud.com/Content/Identity/Concepts/policygetstarted.htm).
type CreateBucketDetails struct {

	// The name of the bucket. Valid characters are uppercase or lowercase letters,
	// numbers, and dashes. Bucket names must be unique within the namespace. Avoid entering confidential information.
	// example: Example: my-new-bucket1
	Name *string `mandatory:"true" json:"name"`

	// The ID of the compartment in which to create the bucket.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// Arbitrary string, up to 4KB, of keys and values for user-defined metadata.
	Metadata map[string]string `mandatory:"false" json:"metadata"`

	// The type of public access enabled on this bucket.
	// A bucket is set to `NoPublicAccess` by default, which only allows an authenticated caller to access the
	// bucket and its contents. When `ObjectRead` is enabled on the bucket, public access is allowed for the
	// `GetObject`, `HeadObject`, and `ListObjects` operations. When `ObjectReadWithoutList` is enabled on the bucket,
	// public access is allowed for the `GetObject` and `HeadObject` operations.
	PublicAccessType CreateBucketDetailsPublicAccessTypeEnum `mandatory:"false" json:"publicAccessType,omitempty"`

	// The type of storage tier of this bucket.
	// A bucket is set to 'Standard' tier by default, which means the bucket will be put in the standard storage tier.
	// When 'Archive' tier type is set explicitly, the bucket is put in the Archive Storage tier. The 'storageTier'
	// property is immutable after bucket is created.
	StorageTier CreateBucketDetailsStorageTierEnum `mandatory:"false" json:"storageTier,omitempty"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.us-phoenix-1.oraclecloud.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.us-phoenix-1.oraclecloud.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m CreateBucketDetails) String() string {
	return common.PointerString(m)
}

// CreateBucketDetailsPublicAccessTypeEnum Enum with underlying type: string
type CreateBucketDetailsPublicAccessTypeEnum string

// Set of constants representing the allowable values for CreateBucketDetailsPublicAccessType
const (
	CreateBucketDetailsPublicAccessTypeNopublicaccess        CreateBucketDetailsPublicAccessTypeEnum = "NoPublicAccess"
	CreateBucketDetailsPublicAccessTypeObjectread            CreateBucketDetailsPublicAccessTypeEnum = "ObjectRead"
	CreateBucketDetailsPublicAccessTypeObjectreadwithoutlist CreateBucketDetailsPublicAccessTypeEnum = "ObjectReadWithoutList"
)

var mappingCreateBucketDetailsPublicAccessType = map[string]CreateBucketDetailsPublicAccessTypeEnum{
	"NoPublicAccess":        CreateBucketDetailsPublicAccessTypeNopublicaccess,
	"ObjectRead":            CreateBucketDetailsPublicAccessTypeObjectread,
	"ObjectReadWithoutList": CreateBucketDetailsPublicAccessTypeObjectreadwithoutlist,
}

// GetCreateBucketDetailsPublicAccessTypeEnumValues Enumerates the set of values for CreateBucketDetailsPublicAccessType
func GetCreateBucketDetailsPublicAccessTypeEnumValues() []CreateBucketDetailsPublicAccessTypeEnum {
	values := make([]CreateBucketDetailsPublicAccessTypeEnum, 0)
	for _, v := range mappingCreateBucketDetailsPublicAccessType {
		values = append(values, v)
	}
	return values
}

// CreateBucketDetailsStorageTierEnum Enum with underlying type: string
type CreateBucketDetailsStorageTierEnum string

// Set of constants representing the allowable values for CreateBucketDetailsStorageTier
const (
	CreateBucketDetailsStorageTierStandard CreateBucketDetailsStorageTierEnum = "Standard"
	CreateBucketDetailsStorageTierArchive  CreateBucketDetailsStorageTierEnum = "Archive"
)

var mappingCreateBucketDetailsStorageTier = map[string]CreateBucketDetailsStorageTierEnum{
	"Standard": CreateBucketDetailsStorageTierStandard,
	"Archive":  CreateBucketDetailsStorageTierArchive,
}

// GetCreateBucketDetailsStorageTierEnumValues Enumerates the set of values for CreateBucketDetailsStorageTier
func GetCreateBucketDetailsStorageTierEnumValues() []CreateBucketDetailsStorageTierEnum {
	values := make([]CreateBucketDetailsStorageTierEnum, 0)
	for _, v := range mappingCreateBucketDetailsStorageTier {
		values = append(values, v)
	}
	return values
}

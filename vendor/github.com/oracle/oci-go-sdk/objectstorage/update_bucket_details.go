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

// UpdateBucketDetails To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.us-phoenix-1.oraclecloud.com/Content/Identity/Concepts/policygetstarted.htm).
type UpdateBucketDetails struct {

	// The namespace in which the bucket lives.
	Namespace *string `mandatory:"false" json:"namespace"`

	// The compartmentId for the compartment to which the bucket is targeted to move to.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The name of the bucket. Avoid entering confidential information.
	// Example: my-new-bucket1
	Name *string `mandatory:"false" json:"name"`

	// Arbitrary string, up to 4KB, of keys and values for user-defined metadata.
	Metadata map[string]string `mandatory:"false" json:"metadata"`

	// The type of public access enabled on this bucket. A bucket is set to `NoPublicAccess` by default, which only allows an
	// authenticated caller to access the bucket and its contents. When `ObjectRead` is enabled on the bucket, public access
	// is allowed for the `GetObject`, `HeadObject`, and `ListObjects` operations. When `ObjectReadWithoutList` is enabled
	// on the bucket, public access is allowed for the `GetObject` and `HeadObject` operations.
	PublicAccessType UpdateBucketDetailsPublicAccessTypeEnum `mandatory:"false" json:"publicAccessType,omitempty"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.us-phoenix-1.oraclecloud.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.us-phoenix-1.oraclecloud.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m UpdateBucketDetails) String() string {
	return common.PointerString(m)
}

// UpdateBucketDetailsPublicAccessTypeEnum Enum with underlying type: string
type UpdateBucketDetailsPublicAccessTypeEnum string

// Set of constants representing the allowable values for UpdateBucketDetailsPublicAccessType
const (
	UpdateBucketDetailsPublicAccessTypeNopublicaccess        UpdateBucketDetailsPublicAccessTypeEnum = "NoPublicAccess"
	UpdateBucketDetailsPublicAccessTypeObjectread            UpdateBucketDetailsPublicAccessTypeEnum = "ObjectRead"
	UpdateBucketDetailsPublicAccessTypeObjectreadwithoutlist UpdateBucketDetailsPublicAccessTypeEnum = "ObjectReadWithoutList"
)

var mappingUpdateBucketDetailsPublicAccessType = map[string]UpdateBucketDetailsPublicAccessTypeEnum{
	"NoPublicAccess":        UpdateBucketDetailsPublicAccessTypeNopublicaccess,
	"ObjectRead":            UpdateBucketDetailsPublicAccessTypeObjectread,
	"ObjectReadWithoutList": UpdateBucketDetailsPublicAccessTypeObjectreadwithoutlist,
}

// GetUpdateBucketDetailsPublicAccessTypeEnumValues Enumerates the set of values for UpdateBucketDetailsPublicAccessType
func GetUpdateBucketDetailsPublicAccessTypeEnumValues() []UpdateBucketDetailsPublicAccessTypeEnum {
	values := make([]UpdateBucketDetailsPublicAccessTypeEnum, 0)
	for _, v := range mappingUpdateBucketDetailsPublicAccessType {
		values = append(values, v)
	}
	return values
}

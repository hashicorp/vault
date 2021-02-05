// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Object Storage Service API
//
// Common set of Object Storage and Archive Storage APIs for managing buckets, objects, and related resources.
// For more information, see Overview of Object Storage (https://docs.cloud.oracle.com/Content/Object/Concepts/objectstorageoverview.htm) and
// Overview of Archive Storage (https://docs.cloud.oracle.com/Content/Archive/Concepts/archivestorageoverview.htm).
//

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
)

// UpdateBucketDetails To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type UpdateBucketDetails struct {

	// The Object Storage namespace in which the bucket lives.
	Namespace *string `mandatory:"false" json:"namespace"`

	// The compartmentId for the compartment to move the bucket to.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The name of the bucket. Valid characters are uppercase or lowercase letters, numbers, hyphens, underscores, and periods.
	// Bucket names must be unique within an Object Storage namespace. Avoid entering confidential information.
	// Example: my-new-bucket1
	Name *string `mandatory:"false" json:"name"`

	// Arbitrary string, up to 4KB, of keys and values for user-defined metadata.
	Metadata map[string]string `mandatory:"false" json:"metadata"`

	// The type of public access enabled on this bucket. A bucket is set to `NoPublicAccess` by default, which only allows an
	// authenticated caller to access the bucket and its contents. When `ObjectRead` is enabled on the bucket, public access
	// is allowed for the `GetObject`, `HeadObject`, and `ListObjects` operations. When `ObjectReadWithoutList` is enabled
	// on the bucket, public access is allowed for the `GetObject` and `HeadObject` operations.
	PublicAccessType UpdateBucketDetailsPublicAccessTypeEnum `mandatory:"false" json:"publicAccessType,omitempty"`

	// Whether or not events are emitted for object state changes in this bucket. By default, `objectEventsEnabled` is
	// set to `false`. Set `objectEventsEnabled` to `true` to emit events for object state changes. For more information
	// about events, see Overview of Events (https://docs.cloud.oracle.com/Content/Events/Concepts/eventsoverview.htm).
	ObjectEventsEnabled *bool `mandatory:"false" json:"objectEventsEnabled"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the Key Management master encryption key to associate
	// with the specified bucket. If this value is empty, the Update operation will remove the associated key, if
	// there is one, from the bucket. (The bucket will continue to be encrypted, but with an encryption key managed
	// by Oracle.)
	KmsKeyId *string `mandatory:"false" json:"kmsKeyId"`

	// The versioning status on the bucket. If in state `Enabled`, multiple versions of the same object can be kept in the bucket.
	// When the object is overwritten or deleted, previous versions will still be available. When versioning is `Suspended`, the previous versions will still remain but new versions will no longer be created when overwitten or deleted.
	// Versioning cannot be disabled on a bucket once enabled.
	Versioning UpdateBucketDetailsVersioningEnum `mandatory:"false" json:"versioning,omitempty"`
}

func (m UpdateBucketDetails) String() string {
	return common.PointerString(m)
}

// UpdateBucketDetailsPublicAccessTypeEnum Enum with underlying type: string
type UpdateBucketDetailsPublicAccessTypeEnum string

// Set of constants representing the allowable values for UpdateBucketDetailsPublicAccessTypeEnum
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

// GetUpdateBucketDetailsPublicAccessTypeEnumValues Enumerates the set of values for UpdateBucketDetailsPublicAccessTypeEnum
func GetUpdateBucketDetailsPublicAccessTypeEnumValues() []UpdateBucketDetailsPublicAccessTypeEnum {
	values := make([]UpdateBucketDetailsPublicAccessTypeEnum, 0)
	for _, v := range mappingUpdateBucketDetailsPublicAccessType {
		values = append(values, v)
	}
	return values
}

// UpdateBucketDetailsVersioningEnum Enum with underlying type: string
type UpdateBucketDetailsVersioningEnum string

// Set of constants representing the allowable values for UpdateBucketDetailsVersioningEnum
const (
	UpdateBucketDetailsVersioningEnabled   UpdateBucketDetailsVersioningEnum = "Enabled"
	UpdateBucketDetailsVersioningSuspended UpdateBucketDetailsVersioningEnum = "Suspended"
)

var mappingUpdateBucketDetailsVersioning = map[string]UpdateBucketDetailsVersioningEnum{
	"Enabled":   UpdateBucketDetailsVersioningEnabled,
	"Suspended": UpdateBucketDetailsVersioningSuspended,
}

// GetUpdateBucketDetailsVersioningEnumValues Enumerates the set of values for UpdateBucketDetailsVersioningEnum
func GetUpdateBucketDetailsVersioningEnumValues() []UpdateBucketDetailsVersioningEnum {
	values := make([]UpdateBucketDetailsVersioningEnum, 0)
	for _, v := range mappingUpdateBucketDetailsVersioning {
		values = append(values, v)
	}
	return values
}

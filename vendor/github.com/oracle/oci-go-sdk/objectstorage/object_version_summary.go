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

// ObjectVersionSummary To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type ObjectVersionSummary struct {

	// The name of the object. Avoid entering confidential information.
	// Example: test/object1.log
	Name *string `mandatory:"true" json:"name"`

	// The date and time the object was modified, as described in RFC 2616 (https://tools.ietf.org/rfc/rfc2616#section-14.29).
	TimeModified *common.SDKTime `mandatory:"true" json:"timeModified"`

	// VersionId of the object.
	VersionId *string `mandatory:"true" json:"versionId"`

	// This flag will indicate if the version is deleted or not.
	IsDeleteMarker *bool `mandatory:"true" json:"isDeleteMarker"`

	// Size of the object in bytes.
	Size *int64 `mandatory:"false" json:"size"`

	// Base64-encoded MD5 hash of the object data.
	Md5 *string `mandatory:"false" json:"md5"`

	// The date and time the object was created, as described in RFC 2616 (https://tools.ietf.org/html/rfc2616#section-14.29).
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The current entity tag (ETag) for the object.
	Etag *string `mandatory:"false" json:"etag"`
}

func (m ObjectVersionSummary) String() string {
	return common.PointerString(m)
}

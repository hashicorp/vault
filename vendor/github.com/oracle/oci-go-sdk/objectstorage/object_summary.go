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

// ObjectSummary To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type ObjectSummary struct {

	// The name of the object. Avoid entering confidential information.
	// Example: test/object1.log
	Name *string `mandatory:"true" json:"name"`

	// Size of the object in bytes.
	Size *int64 `mandatory:"false" json:"size"`

	// Base64-encoded MD5 hash of the object data.
	Md5 *string `mandatory:"false" json:"md5"`

	// The date and time the object was created, as described in RFC 2616 (https://tools.ietf.org/rfc/rfc2616), section 14.29.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m ObjectSummary) String() string {
	return common.PointerString(m)
}

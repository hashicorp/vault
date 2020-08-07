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

// ObjectLifecyclePolicy The collection of lifecycle policy rules that together form the object lifecycle policy of a given bucket.
type ObjectLifecyclePolicy struct {

	// The date and time the object lifecycle policy was created, as described in
	// RFC 3339 (https://tools.ietf.org/rfc/rfc3339), section 14.29.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The live lifecycle policy on the bucket.
	// For an example of this value, see the
	// PutObjectLifecyclePolicy API documentation (https://docs.cloud.oracle.com/iaas/api/#/en/objectstorage/20160918/ObjectLifecyclePolicy/PutObjectLifecyclePolicy).
	Items []ObjectLifecycleRule `mandatory:"false" json:"items"`
}

func (m ObjectLifecyclePolicy) String() string {
	return common.PointerString(m)
}

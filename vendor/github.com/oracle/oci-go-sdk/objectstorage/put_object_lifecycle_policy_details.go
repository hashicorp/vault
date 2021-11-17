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

// PutObjectLifecyclePolicyDetails Creates a new object lifecycle policy for a bucket.
type PutObjectLifecyclePolicyDetails struct {

	// The bucket's set of lifecycle policy rules.
	Items []ObjectLifecycleRule `mandatory:"false" json:"items"`
}

func (m PutObjectLifecyclePolicyDetails) String() string {
	return common.PointerString(m)
}

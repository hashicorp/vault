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

// ReplicationSource The details of a replication source bucket that replicates to a target destination bucket.
type ReplicationSource struct {

	// The name of the policy.
	PolicyName *string `mandatory:"true" json:"policyName"`

	// The source region replicating data from, for example "us-ashburn-1".
	SourceRegionName *string `mandatory:"true" json:"sourceRegionName"`

	// The source bucket replicating data from.
	SourceBucketName *string `mandatory:"true" json:"sourceBucketName"`
}

func (m ReplicationSource) String() string {
	return common.PointerString(m)
}

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

// CommitMultipartUploadDetails To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.us-phoenix-1.oraclecloud.com/Content/Identity/Concepts/policygetstarted.htm).
type CommitMultipartUploadDetails struct {

	// The part numbers and ETags for the parts to be committed.
	PartsToCommit []CommitMultipartUploadPartDetails `mandatory:"true" json:"partsToCommit"`

	// The part numbers for the parts to be excluded from the completed object.
	// Each part created for this upload must be in either partsToExclude or partsToCommit, but cannot be in both.
	PartsToExclude []int `mandatory:"false" json:"partsToExclude"`
}

func (m CommitMultipartUploadDetails) String() string {
	return common.PointerString(m)
}

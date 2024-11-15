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

// CommitMultipartUploadDetails To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type CommitMultipartUploadDetails struct {

	// The part numbers and entity tags (ETags) for the parts to be committed.
	PartsToCommit []CommitMultipartUploadPartDetails `mandatory:"true" json:"partsToCommit"`

	// The part numbers for the parts to be excluded from the completed object.
	// Each part created for this upload must be in either partsToExclude or partsToCommit, but cannot be in both.
	PartsToExclude []int `mandatory:"false" json:"partsToExclude"`
}

func (m CommitMultipartUploadDetails) String() string {
	return common.PointerString(m)
}

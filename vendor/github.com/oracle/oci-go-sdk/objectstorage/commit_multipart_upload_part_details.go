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

// CommitMultipartUploadPartDetails To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type CommitMultipartUploadPartDetails struct {

	// The part number for this part.
	PartNum *int `mandatory:"true" json:"partNum"`

	// The entity tag (ETag) returned when this part was uploaded.
	Etag *string `mandatory:"true" json:"etag"`
}

func (m CommitMultipartUploadPartDetails) String() string {
	return common.PointerString(m)
}

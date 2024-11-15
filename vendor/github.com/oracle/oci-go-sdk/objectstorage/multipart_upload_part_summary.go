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

// MultipartUploadPartSummary Gets summary information about multipart uploads.
// To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access,
// see Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type MultipartUploadPartSummary struct {

	// The current entity tag (ETag) for the part.
	Etag *string `mandatory:"true" json:"etag"`

	// The MD5 hash of the bytes of the part.
	Md5 *string `mandatory:"true" json:"md5"`

	// The size of the part in bytes.
	Size *int64 `mandatory:"true" json:"size"`

	// The part number for this part.
	PartNumber *int `mandatory:"true" json:"partNumber"`
}

func (m MultipartUploadPartSummary) String() string {
	return common.PointerString(m)
}

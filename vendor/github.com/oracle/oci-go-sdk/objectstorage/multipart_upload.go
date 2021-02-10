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

// MultipartUpload Multipart uploads provide efficient and resilient uploads, especially for large objects. Multipart uploads also accommodate
// objects that are too large for a single upload operation. With multipart uploads, individual parts of an object can be
// uploaded in parallel to reduce the amount of time you spend uploading. Multipart uploads can also minimize the impact
// of network failures by letting you retry a failed part upload instead of requiring you to retry an entire object upload.
// See Using Multipart Uploads (https://docs.cloud.oracle.com/Content/Object/Tasks/usingmultipartuploads.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you are not authorized,
// talk to an administrator. If you are an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
type MultipartUpload struct {

	// The Object Storage namespace in which the in-progress multipart upload is stored.
	Namespace *string `mandatory:"true" json:"namespace"`

	// The bucket in which the in-progress multipart upload is stored.
	Bucket *string `mandatory:"true" json:"bucket"`

	// The object name of the in-progress multipart upload.
	Object *string `mandatory:"true" json:"object"`

	// The unique identifier for the in-progress multipart upload.
	UploadId *string `mandatory:"true" json:"uploadId"`

	// The date and time the upload was created, as described in RFC 2616 (https://tools.ietf.org/html/rfc2616#section-14.29).
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`
}

func (m MultipartUpload) String() string {
	return common.PointerString(m)
}

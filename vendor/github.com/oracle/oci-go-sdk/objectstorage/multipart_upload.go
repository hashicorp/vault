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

// MultipartUpload Multipart uploads provide efficient and resilient uploads, especially for large objects. Multipart uploads also accommodate
// objects that are too large for a single upload operation. With multipart uploads, individual parts of an object can be
// uploaded in parallel to reduce the amount of time you spend uploading. Multipart uploads can also minimize the impact
// of network failures by letting you retry a failed part upload instead of requiring you to retry an entire object upload.
// See Managing Multipart Uploads (https://docs.us-phoenix-1.oraclecloud.com/Content/Object/Tasks/managingmultipartuploads.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.us-phoenix-1.oraclecloud.com/Content/Identity/Concepts/policygetstarted.htm).
type MultipartUpload struct {

	// The namespace in which the in-progress multipart upload is stored.
	Namespace *string `mandatory:"true" json:"namespace"`

	// The bucket in which the in-progress multipart upload is stored.
	Bucket *string `mandatory:"true" json:"bucket"`

	// The object name of the in-progress multipart upload.
	Object *string `mandatory:"true" json:"object"`

	// The unique identifier for the in-progress multipart upload.
	UploadId *string `mandatory:"true" json:"uploadId"`

	// The date and time the upload was created, as described in RFC 2616 (https://tools.ietf.org/rfc/rfc2616), section 14.29.
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`
}

func (m MultipartUpload) String() string {
	return common.PointerString(m)
}

// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

package transfer

import (
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
)

// requestValidator validate user's input and assign default values if not defined
type requestValidator interface {
	// validate inputs, return error if request is not valid
	validate() error

	// assign default values
	assignDefaultValues() error
}

// UploadRequest defines the input parameters for UploadFile method
type UploadRequest struct {
	// The top-level namespace used for the request.
	NamespaceName *string `mandatory:"true"`

	// The name of the bucket. Avoid entering confidential information. Example: my-new-bucket1
	BucketName *string `mandatory:"true"`

	// The name of the object. Avoid entering confidential information. Example: test/object1.log
	ObjectName *string `mandatory:"true"`

	// [Optional] Override the default part size of 128 MiB, value is in bytes.
	PartSize *int64 `mandatory:"false"`

	// [Optional] Whether or not this UploadManager supports performing mulitpart uploads. Defaults to True.
	AllowMultipartUploads *bool `mandatory:"false"`

	// [Optional] Whether or not this UploadManager supports uploading individual parts of a multipart upload in parallel.
	// This setting has no effect on uploads that are performed using a single put_object call. Defaults to True.
	AllowParrallelUploads *bool `mandatory:"false"`

	// The number of go routines for uploading individual parts of a multipart upload.
	// This setting is only used if allow_parallel_uploads is set to True. Defaults to 10.
	// TODO: check with service team for upper bounds of the concurrent number of uploads
	// and update the document here
	NumberOfGoroutines *int `mandatory:"false"`

	// A configured object storage client to use for interacting with the Object Storage service.
	ObjectStorageClient *objectstorage.ObjectStorageClient `mandatory:"false"`

	// [Optional] The entity tag of the object to match.
	IfMatch *string `mandatory:"false"`

	// [Optional] The entity tag of the object to avoid matching. The only valid value is ‘*’,
	// which indicates that the request should fail if the object already exists.
	IfNoneMatch *string `mandatory:"false"`

	// [Optional] The base-64 encoded MD5 hash of the body. This parameter is only used if the object is uploaded in a single part.
	ContentMD5 *string `mandatory:"false"`

	// [Optional] The content type of the object to upload.
	ContentType *string `mandatory:"false"`

	// [Optional] The content language of the object to upload.
	ContentLanguage *string `mandatory:"false"`

	// [Optional] The content encoding of the object to upload.
	ContentEncoding *string `mandatory:"false"`

	// [Optional] Arbitrary string keys and values for the user-defined metadata for the object.
	// Keys must be in "opc-meta-*" format.
	Metadata map[string]string `mandatory:"false"`

	// [Optional] The client request ID for tracing.
	OpcClientRequestID *string `mandatory:"false"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request UploadRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

var (
	errorInvalidNamespace  = errors.New("namespaceName is required")
	errorInvalidBucketName = errors.New("bucketName is required")
	errorInvalidObjectName = errors.New("objectName is required")
)

const defaultNumberOfGoroutines = 5 // increase the value might cause 409 error form service and client timeout

func (request UploadRequest) validate() error {
	if request.NamespaceName == nil {
		return errorInvalidNamespace
	}

	if request.BucketName == nil {
		return errorInvalidBucketName
	}

	if request.ObjectName == nil {
		return errorInvalidObjectName
	}

	return nil
}

func (request *UploadRequest) initDefaultValues() error {
	if request.ObjectStorageClient == nil {
		client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(common.DefaultConfigProvider())

		// default timeout is 60s which includes the time for reading the body
		// default timeout doesn't work for big file, here will use the default
		// 0s which means no timeout
		client.HTTPClient = &http.Client{}

		if err != nil {
			return err
		}

		request.ObjectStorageClient = &client
	}

	if request.NumberOfGoroutines == nil ||
		*request.NumberOfGoroutines <= 0 {
		request.NumberOfGoroutines = common.Int(defaultNumberOfGoroutines)
	}

	if request.AllowMultipartUploads == nil {
		request.AllowMultipartUploads = common.Bool(true)
	}

	if request.AllowParrallelUploads == nil {
		request.AllowParrallelUploads = common.Bool(true)
	}

	if !*request.AllowParrallelUploads {
		request.NumberOfGoroutines = common.Int(1) // one go routine for upload
	}

	if request.RetryPolicy() == nil {
		// default retry policy
		request.RequestMetadata = common.RequestMetadata{RetryPolicy: getUploadManagerDefaultRetryPolicy()}
	}

	return nil
}

// UploadResponseType with underlying type: string
type UploadResponseType string

// Set of constants representing the allowable values for VolumeAttachmentLifecycleState
const (
	MultipartUpload  UploadResponseType = "MULTIPARTUPLOAD"
	SinglepartUpload UploadResponseType = "SINGLEPARTUPLOAD"
)

// UploadResponse is the response from commitMultipart or the putObject API operation.
type UploadResponse struct {

	// Polymorphic reponse type indicates the response type
	Type UploadResponseType

	// response for putObject API response (single part upload), will be nil if the operation is multiPart upload
	*SinglepartUploadResponse

	// response for commitMultipart API response (multipart upload), will be nil if the operation is singlePart upload
	*MultipartUploadResponse
}

// IsResumable is a function to check is previous failed upload resumable or not
func (resp UploadResponse) IsResumable() bool {
	if resp.Type == SinglepartUpload {
		return false
	}

	return *resp.MultipartUploadResponse.isResumable
}

// SinglepartUploadResponse is the response from putObject API operation.
type SinglepartUploadResponse struct {
	objectstorage.PutObjectResponse
}

// MultipartUploadResponse is the response from commitMultipart API operation.
type MultipartUploadResponse struct {
	objectstorage.CommitMultipartUploadResponse

	// The upload ID for a multipart upload.
	UploadID *string

	// The value indicates is the operation IsResumable or not, call the resume function if is true
	isResumable *bool
}

func getUploadManagerDefaultRetryPolicy() *common.RetryPolicy {
	attempts := uint(3)
	retryOnAllNon200ResponseCodes := func(r common.OCIOperationResponse) bool {
		return !(r.Error == nil && 199 < r.Response.HTTPResponse().StatusCode && r.Response.HTTPResponse().StatusCode < 300)
	}

	exponentialBackoff := func(r common.OCIOperationResponse) time.Duration {
		return time.Duration(math.Pow(float64(2), float64(r.AttemptNumber-1))) * time.Second
	}
	policy := common.NewRetryPolicy(attempts, retryOnAllNon200ResponseCodes, exponentialBackoff)

	return &policy
}

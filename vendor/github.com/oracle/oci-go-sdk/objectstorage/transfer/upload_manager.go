// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

// Package transfer simplifies interaction with the Object Storage service by abstracting away the method used
// to upload objects.  Depending on the configuration parameters, UploadManager may choose to do a single
// put_object request, or break up the upload into multiple parts and utilize multi-part uploads.
//
// An advantage of using multi-part uploads is the ability to retry individual failed parts, as well as being
// able to upload parts in parallel to reduce upload time.
//
// To use this package, you must be authorized in an IAM policy. If you're not authorized, talk to an administrator.
package transfer

import (
	"context"
	"errors"
	"math"
	"os"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/common"
)

// UploadManager is the interface that groups the upload methods
type UploadManager struct {
	FileUploader   FileUploader
	StreamUploader StreamUploader
}

var (
	errorInvalidStreamUploader = errors.New("streamUploader is required, use NewUploadManager for default implementation")
	errorInvalidFileUploader   = errors.New("fileUploader is required, use NewUploadManager for default implementation")
)

// NewUploadManager return a pointer to UploadManager
func NewUploadManager() *UploadManager {
	return &UploadManager{
		FileUploader:   &fileUpload{multipartUploader: &multipartUpload{}},
		StreamUploader: &streamUpload{multipartUploader: &multipartUpload{}},
	}
}

// UploadFile uploads an object to Object Storage. Depending on the options provided and the
// size of the object, the object may be uploaded in multiple parts or just an signle object.
func (uploadManager *UploadManager) UploadFile(ctx context.Context, request UploadFileRequest) (response UploadResponse, err error) {
	if err = request.validate(); err != nil {
		return
	}

	if err = request.initDefaultValues(); err != nil {
		return
	}

	if uploadManager.FileUploader == nil {
		err = errorInvalidFileUploader
		return
	}

	file, err := os.Open(request.FilePath)
	defer file.Close()

	if err != nil {
		return
	}

	fi, err := file.Stat()
	if err != nil {
		return
	}

	fileSize := fi.Size()

	// parrallel upload disabled by user or the file size smaller than or equal to partSize
	// use UploadFilePutObject
	if !*request.AllowMultipartUploads ||
		int64(fileSize) <= *request.PartSize {
		response, err = uploadManager.FileUploader.UploadFilePutObject(ctx, request)
		return
	}

	response, err = uploadManager.FileUploader.UploadFileMultiparts(ctx, request)
	return
}

// ResumeUploadFile resumes a multipart file upload.
func (uploadManager *UploadManager) ResumeUploadFile(ctx context.Context, uploadID string) (response UploadResponse, err error) {
	if len(strings.TrimSpace(uploadID)) == 0 {
		err = errors.New("uploadID is required to resume a multipart file upload")
		return
	}
	response, err = uploadManager.FileUploader.ResumeUploadFile(ctx, uploadID)
	return
}

// UploadStream uploads streaming data to Object Storage. If the stream is non-empty, this will always perform a
// multipart upload, splitting parts based on the part size (10 MiB if none specified). If the stream is empty,
// this will upload a single empty object to Object Storage.
// Stream uploads are not currently resumable.
func (uploadManager *UploadManager) UploadStream(ctx context.Context, request UploadStreamRequest) (response UploadResponse, err error) {
	if err = request.validate(); err != nil {
		return
	}

	if err = request.initDefaultValues(); err != nil {
		return
	}

	if uploadManager.StreamUploader == nil {
		err = errorInvalidStreamUploader
		return
	}

	response, err = uploadManager.StreamUploader.UploadStream(ctx, request)
	return
}

func getUploadManagerRetryPolicy() *common.RetryPolicy {
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

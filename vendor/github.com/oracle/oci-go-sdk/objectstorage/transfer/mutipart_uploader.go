// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

package transfer

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/objectstorage"
)

// multipartUploader is an interface wrap the methods talk to object storage service
type multipartUploader interface {
	createMultipartUpload(ctx context.Context, request UploadRequest) (string, error)
	uploadParts(ctx context.Context, done <-chan struct{}, parts <-chan uploadPart, result chan<- uploadPart, request UploadRequest, uploadID string)
	uploadPart(ctx context.Context, request UploadRequest, part uploadPart, uploadID string) (objectstorage.UploadPartResponse, error)
	commit(ctx context.Context, request UploadRequest, parts map[int]uploadPart, uploadID string) (resp objectstorage.CommitMultipartUploadResponse, err error)
}

// multipartUpload implements multipartUploader interface
type multipartUpload struct{}

// createMultipartUpload creates a new multipart upload in Object Storage and return the uploadId
func (uploader *multipartUpload) createMultipartUpload(ctx context.Context, request UploadRequest) (string, error) {
	multipartUploadRequest := objectstorage.CreateMultipartUploadRequest{
		NamespaceName:      request.NamespaceName,
		BucketName:         request.BucketName,
		IfMatch:            request.IfMatch,
		IfNoneMatch:        request.IfNoneMatch,
		OpcClientRequestId: request.OpcClientRequestID,
	}

	multipartUploadRequest.Object = request.ObjectName
	multipartUploadRequest.ContentType = request.ContentType
	multipartUploadRequest.ContentEncoding = request.ContentEncoding
	multipartUploadRequest.ContentLanguage = request.ContentLanguage
	multipartUploadRequest.Metadata = request.Metadata

	resp, err := request.ObjectStorageClient.CreateMultipartUpload(ctx, multipartUploadRequest)
	return *resp.UploadId, err
}

func (uploader *multipartUpload) uploadParts(ctx context.Context, done <-chan struct{}, parts <-chan uploadPart, result chan<- uploadPart, request UploadRequest, uploadID string) {
	// loop through the part from parts channel created by splitFileParts method
	for part := range parts {
		if part.err != nil {
			// ignore this part which contains error from split function
			result <- part
			return
		}

		resp, err := uploader.uploadPart(ctx, request, part, uploadID)
		if err != nil {
			common.Debugf("upload error %v\n", err)
			part.err = err
		}
		part.etag = resp.ETag
		select {
		case result <- part:
			common.Debugf("uploadParts resp %v, %v\n", part.partNum, resp.ETag)
		case <-done:
			common.Debugln("uploadParts received Done")
			return
		}
	}
}

// send request to upload part to object storage
func (uploader *multipartUpload) uploadPart(ctx context.Context, request UploadRequest, part uploadPart, uploadID string) (objectstorage.UploadPartResponse, error) {
	req := objectstorage.UploadPartRequest{
		NamespaceName:      request.NamespaceName,
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		UploadId:           common.String(uploadID),
		UploadPartNum:      common.Int(part.partNum),
		UploadPartBody:     ioutil.NopCloser(bytes.NewReader(part.partBody)),
		ContentLength:      common.Int64(part.size),
		IfMatch:            request.IfMatch,
		IfNoneMatch:        request.IfNoneMatch,
		OpcClientRequestId: request.OpcClientRequestID,
		RequestMetadata:    request.RequestMetadata,
	}

	resp, err := request.ObjectStorageClient.UploadPart(ctx, req)

	return resp, err
}

// commits the multipart upload
func (uploader *multipartUpload) commit(ctx context.Context, request UploadRequest, parts map[int]uploadPart, uploadID string) (resp objectstorage.CommitMultipartUploadResponse, err error) {
	req := objectstorage.CommitMultipartUploadRequest{
		NamespaceName:      request.NamespaceName,
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		UploadId:           common.String(uploadID),
		IfMatch:            request.IfMatch,
		IfNoneMatch:        request.IfNoneMatch,
		OpcClientRequestId: request.OpcClientRequestID,
		RequestMetadata:    request.RequestMetadata,
	}

	partsToCommit := []objectstorage.CommitMultipartUploadPartDetails{}

	for _, part := range parts {
		if part.etag != nil {
			detail := objectstorage.CommitMultipartUploadPartDetails{
				Etag:    part.etag,
				PartNum: common.Int(part.partNum),
			}

			// update the parts to commit
			partsToCommit = append(partsToCommit, detail)
		} else {
			// some parts failed, return error for resume
			common.Debugf("uploadPart has error: %v\n", part.err)
			err = part.err
			return
		}
	}

	req.PartsToCommit = partsToCommit
	resp, err = request.ObjectStorageClient.CommitMultipartUpload(ctx, req)
	return
}

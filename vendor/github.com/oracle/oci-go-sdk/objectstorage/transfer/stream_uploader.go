// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

package transfer

import (
	"context"
	"sync"

	"github.com/oracle/oci-go-sdk/common"
)

// StreamUploader is an interface for upload a stream
type StreamUploader interface {
	// uploads a stream to blob storage
	UploadStream(ctx context.Context, request UploadStreamRequest) (response UploadResponse, err error)
}

type streamUpload struct {
	uploadID          string
	manifest          *multipartManifest
	multipartUploader multipartUploader
	request           UploadStreamRequest
}

func (streamUpload *streamUpload) UploadStream(ctx context.Context, request UploadStreamRequest) (response UploadResponse, err error) {

	uploadID, err := streamUpload.multipartUploader.createMultipartUpload(ctx, request.UploadRequest)
	streamUpload.uploadID = uploadID

	if err != nil {
		return UploadResponse{}, err
	}

	if streamUpload.manifest == nil {
		streamUpload.manifest = &multipartManifest{parts: make(map[string]map[int]uploadPart)}
	}

	// UploadFileMultiparts closes the done channel when it returns
	done := make(chan struct{})
	defer close(done)
	parts := streamUpload.manifest.splitStreamToParts(done, *request.PartSize, request.StreamReader)

	return streamUpload.startConcurrentUpload(ctx, done, parts, request)
}

func (streamUpload *streamUpload) startConcurrentUpload(ctx context.Context, done <-chan struct{}, parts <-chan uploadPart, request UploadStreamRequest) (response UploadResponse, err error) {
	result := make(chan uploadPart)
	numUploads := *request.NumberOfGoroutines
	var wg sync.WaitGroup
	wg.Add(numUploads)

	// start fixed number of goroutines to upload parts
	for i := 0; i < numUploads; i++ {
		go func() {
			streamUpload.multipartUploader.uploadParts(ctx, done, parts, result, request.UploadRequest, streamUpload.uploadID)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	streamUpload.manifest.updateManifest(result, streamUpload.uploadID)
	resp, err := streamUpload.multipartUploader.commit(ctx, request.UploadRequest, streamUpload.manifest.parts[streamUpload.uploadID], streamUpload.uploadID)

	if err != nil {
		common.Debugf("failed to commit with error: %v\n", err)
		return UploadResponse{
				Type: MultipartUpload,
				MultipartUploadResponse: &MultipartUploadResponse{UploadID: common.String(streamUpload.uploadID)}},
			err
	}

	response = UploadResponse{
		Type: MultipartUpload,
		MultipartUploadResponse: &MultipartUploadResponse{CommitMultipartUploadResponse: resp},
	}
	return
}

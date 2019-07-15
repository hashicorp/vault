// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

package transfer

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/oracle/oci-go-sdk/objectstorage"
)

type fake struct {
	upLoadID              string
	failedPartNumbers     []int // use to simulate the part has error to upload for retry and resume
	numberOfUploadedParts *int
	numberOfCommitedParts *int
	resumedPartNumbers    []int // parts which are uploaded via resume
}

func (fake *fake) createMultipartUpload(ctx context.Context, request UploadRequest) (string, error) {
	return fake.upLoadID, nil
}

func (fake *fake) uploadParts(ctx context.Context, done <-chan struct{}, parts <-chan uploadPart, result chan<- uploadPart, request UploadRequest, uploadID string) {
	// loop through the part from parts channel created by splitFileParts method
	for part := range parts {
		resp, err := fake.uploadPart(ctx, request, part, uploadID)
		part.etag = resp.ETag
		part.err = err

		select {
		case result <- part:
		case <-done:
			return
		}
	}
}

func (fake *fake) uploadPart(ctx context.Context, request UploadRequest, part uploadPart, uploadID string) (objectstorage.UploadPartResponse, error) {
	// mark parts as failed
	for _, failedPartNum := range fake.failedPartNumbers {
		if failedPartNum == part.partNum {
			// simulate part upload failed
			return objectstorage.UploadPartResponse{}, errors.New("upload failed")
		}
	}

	*fake.numberOfUploadedParts++
	return objectstorage.UploadPartResponse{ETag: common.String("etag")}, nil
}

func (fake *fake) commit(ctx context.Context, request UploadRequest, parts map[int]uploadPart, uploadID string) (resp objectstorage.CommitMultipartUploadResponse, err error) {
	*fake.numberOfCommitedParts = 0
	for _, part := range parts {
		if part.err == nil {
			*fake.numberOfCommitedParts++
		} else {
			*fake.numberOfCommitedParts = 0
			break
		}
	}

	return
}

func TestUploadFileMultiparts(t *testing.T) {
	type testStruct struct {
		uploadID                   string
		failedPartNumbers          []int
		expectedUploadID           string
		expectedCachedNumOfRequest int
		expectedNumOfUploadParts   int
		expectedNumOfCommitParts   int
	}

	testDataSet := []testStruct{
		{
			uploadID:                   "id1",
			failedPartNumbers:          []int{},
			expectedUploadID:           "id1",
			expectedCachedNumOfRequest: 1,
			expectedNumOfUploadParts:   10,
			expectedNumOfCommitParts:   10,
		},
		{
			uploadID:                   "id2",
			failedPartNumbers:          []int{1, 2},
			expectedUploadID:           "id2",
			expectedCachedNumOfRequest: 2,
			expectedNumOfUploadParts:   8,
			expectedNumOfCommitParts:   10,
		},
	}

	fileUpload := fileUpload{}
	var partSize, fileSize int64
	fileSize = 100
	partSize = 10

	ctx := context.Background()
	for _, testData := range testDataSet {
		fake := fake{upLoadID: testData.uploadID, failedPartNumbers: testData.failedPartNumbers, numberOfCommitedParts: common.Int(0), numberOfUploadedParts: common.Int(0)}
		fileUpload.multipartUploader = &fake
		filePath, _ := helpers.WriteTempFileOfSize(int64(fileSize))
		request := UploadFileRequest{
			UploadRequest: UploadRequest{PartSize: common.Int64(partSize)},
			FilePath:      filePath,
		}

		request.initDefaultValues()
		resp, err := fileUpload.UploadFileMultiparts(ctx, request)
		assert.Equal(t, testData.expectedUploadID, *resp.UploadID)
		assert.NotEmpty(t, fileUpload.fileUploadReqs[testData.uploadID])
		assert.Equal(t, testData.expectedCachedNumOfRequest, len(fileUpload.fileUploadReqs))
		assert.NotEmpty(t, fileUpload.manifest.parts)

		// all parts have been committed
		totalParts := int(fileSize / partSize)
		failedParts := len(fake.failedPartNumbers)
		assert.Equal(t, testData.expectedNumOfUploadParts, totalParts-failedParts)
		assert.NoError(t, err)

		if failedParts != 0 {
			// no parts should be commit as there are failed parts
			assert.Equal(t, 0, *fake.numberOfCommitedParts)

			// parts gonna to be resumed
			fake.resumedPartNumbers = fake.failedPartNumbers

			// empty the failed part numbers array
			fake.failedPartNumbers = []int{}
			fileUpload.ResumeUploadFile(ctx, testData.uploadID)
			assert.Equal(t, totalParts, *fake.numberOfCommitedParts)
		}

		// finally, all parts should be committed
		assert.Equal(t, testData.expectedNumOfCommitParts, *fake.numberOfCommitedParts)
	}
}

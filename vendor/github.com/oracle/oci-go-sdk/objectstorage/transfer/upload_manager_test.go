// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

package transfer

import (
	"context"
	"testing"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/example/helpers"
	"github.com/stretchr/testify/assert"
)

type fakeFileUpload struct{}

// split file into multiple parts and uploads them to blob storage, then merge
func (fake fakeFileUpload) UploadFileMultiparts(ctx context.Context, request UploadFileRequest) (response UploadResponse, err error) {
	response = UploadResponse{
		Type: MultipartUpload,
	}

	return
}

// uploads a file to blob storage via PutObject API
func (fake fakeFileUpload) UploadFilePutObject(ctx context.Context, request UploadFileRequest) (response UploadResponse, err error) {
	response = UploadResponse{
		Type: SinglepartUpload,
	}

	return
}

// resume a file upload, use it when UploadFile failed
func (fake fakeFileUpload) ResumeUploadFile(ctx context.Context, uploadID string) (response UploadResponse, err error) {
	return
}

func TestUploadManager_UploadFile(t *testing.T) {
	type testData struct {
		FileSize            int
		PartSize            *int64
		FileUploader        FileUploader
		ExpectedResponsType UploadResponseType
		ExpectedError       error
	}

	testDataSet := []testData{
		{
			FileSize:            100,
			PartSize:            common.Int64(1000),
			FileUploader:        fakeFileUpload{},
			ExpectedResponsType: SinglepartUpload,
		}, {
			FileSize:            60,
			PartSize:            common.Int64(50),
			FileUploader:        fakeFileUpload{},
			ExpectedResponsType: MultipartUpload,
		}, {
			FileSize:      60,
			PartSize:      common.Int64(50),
			FileUploader:  nil,
			ExpectedError: errorInvalidFileUploader,
		},
	}

	for _, testData := range testDataSet {
		uploadManager := UploadManager{FileUploader: testData.FileUploader}
		// small file fits into one part
		filePath, _ := helpers.WriteTempFileOfSize(int64(testData.FileSize))
		req := UploadFileRequest{
			UploadRequest: UploadRequest{
				NamespaceName: common.String("namespace"),
				BucketName:    common.String("bname"),
				ObjectName:    common.String("objectName"),
				PartSize:      testData.PartSize,
			},
			FilePath: filePath,
		}

		resp, err := uploadManager.UploadFile(context.Background(), req)
		assert.Equal(t, err, testData.ExpectedError)
		assert.Equal(t, testData.ExpectedResponsType, resp.Type)
	}

}

func TestUploadManager_ResumeUploadFile(t *testing.T) {
	fileUploader := fakeFileUpload{}
	uploadManager := UploadManager{FileUploader: fileUploader}
	_, err := uploadManager.ResumeUploadFile(context.Background(), "")
	assert.Error(t, err)
}

type fakeReader struct{}

func (fr fakeReader) Read(p []byte) (n int, err error) {
	return
}

func TestUploadManager_UploadStream(t *testing.T) {

	req := UploadStreamRequest{
		UploadRequest: UploadRequest{
			NamespaceName: common.String("namespace"),
			BucketName:    common.String("bname"),
			ObjectName:    common.String("objectName"),
		},
		StreamReader: fakeReader{},
	}
	uploadManager := UploadManager{StreamUploader: nil}
	_, err := uploadManager.UploadStream(context.Background(), req)
	assert.Equal(t, errorInvalidStreamUploader, err)
}

// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

package transfer

import (
	"testing"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/stretchr/testify/assert"
)

func TestUploadReqest_validate(t *testing.T) {
	type testData struct {
		NamespaceName *string
		BucketName    *string
		ObjectName    *string
		ExpectedError error
	}

	testDataSet := []testData{
		{
			NamespaceName: nil,
			BucketName:    common.String("test"),
			ObjectName:    common.String("test"),
			ExpectedError: errorInvalidNamespace,
		}, {
			NamespaceName: common.String("test"),
			BucketName:    nil,
			ObjectName:    common.String("test"),
			ExpectedError: errorInvalidBucketName,
		}, {
			NamespaceName: common.String("test"),
			BucketName:    common.String("test"),
			ObjectName:    nil,
			ExpectedError: errorInvalidObjectName,
		},
	}

	for _, testData := range testDataSet {
		req := UploadRequest{
			NamespaceName: testData.NamespaceName,
			BucketName:    testData.BucketName,
			ObjectName:    testData.ObjectName,
		}

		err := req.validate()
		assert.Equal(t, testData.ExpectedError, err)
	}
}

func TestUploadReqest_initDefaultValues(t *testing.T) {
	req := UploadRequest{}
	err := req.initDefaultValues()
	assert.NoError(t, err)
	assert.Equal(t, defaultNumberOfGoroutines, *req.NumberOfGoroutines)
	assert.Equal(t, true, *req.AllowMultipartUploads)
	assert.Equal(t, true, *req.AllowParrallelUploads)
	assert.NotEmpty(t, req.ObjectStorageClient)
}

func TestUploadFileReqest_initDefaultValues(t *testing.T) {
	req := UploadFileRequest{}
	err := req.initDefaultValues()
	assert.NoError(t, err)
	assert.Equal(t, int64(defaultFilePartSize), *req.PartSize)
}

func TestUploadFileReqest_validate(t *testing.T) {
	req := UploadFileRequest{
		UploadRequest: UploadRequest{
			NamespaceName: common.String("test"),
			BucketName:    common.String("test"),
			ObjectName:    common.String("test"),
		},
		FilePath: "",
	}

	err := req.validate()
	assert.Equal(t, errorInvalidFilePath, err)
}

func TestUploadStreamReqest_validate(t *testing.T) {
	req := UploadStreamRequest{
		UploadRequest: UploadRequest{
			NamespaceName: common.String("test"),
			BucketName:    common.String("test"),
			ObjectName:    common.String("test"),
		},
		StreamReader: nil,
	}

	err := req.validate()
	assert.Equal(t, errorInvalidStream, err)
}

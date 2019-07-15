// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

package transfer

import (
	"errors"

	"github.com/oracle/oci-go-sdk/common"
)

// UploadFileRequest defines the input parameters for UploadFile method
type UploadFileRequest struct {
	UploadRequest

	// The path of the file to be uploaded (includs file name)
	FilePath string
}

var errorInvalidFilePath = errors.New("filePath is required")

const defaultFilePartSize = 128 * 1024 * 1024 // 128MB

func (request UploadFileRequest) validate() error {
	err := request.UploadRequest.validate()

	if err != nil {
		return err
	}

	if len(request.FilePath) == 0 {
		return errorInvalidFilePath
	}

	return nil
}

func (request *UploadFileRequest) initDefaultValues() error {
	if request.PartSize == nil {
		request.PartSize = common.Int64(defaultFilePartSize)
	}

	return request.UploadRequest.initDefaultValues()
}

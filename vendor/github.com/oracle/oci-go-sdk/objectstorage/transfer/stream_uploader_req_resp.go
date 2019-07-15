// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

package transfer

import (
	"errors"
	"io"

	"github.com/oracle/oci-go-sdk/common"
)

// UploadStreamRequest defines the input parameters for UploadFile method
type UploadStreamRequest struct {
	UploadRequest

	// The reader of input stream
	StreamReader io.Reader
}

var errorInvalidStream = errors.New("uploadStream is required")

const defaultStreamPartSize = 10 * 1024 * 1024 // 10MB

func (request UploadStreamRequest) validate() error {
	err := request.UploadRequest.validate()

	if err != nil {
		return err
	}

	if request.StreamReader == nil {
		return errorInvalidStream
	}

	return nil
}

func (request *UploadStreamRequest) initDefaultValues() error {
	if request.PartSize == nil {
		request.PartSize = common.Int64(defaultStreamPartSize)
	}

	return request.UploadRequest.initDefaultValues()
}

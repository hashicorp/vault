// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.

package common

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors_ServiceFailureFromResponse(t *testing.T) {
	header := http.Header{}
	opcID := "111"
	header.Set("opc-request-id", opcID)
	sampleResponse := `{"key" : "test"}`

	httpResponse := http.Response{Header: header}
	bodyBuffer := bytes.NewBufferString(sampleResponse)
	httpResponse.Body = ioutil.NopCloser(bodyBuffer)
	httpResponse.StatusCode = 200
	err := newServiceFailureFromResponse(&httpResponse)
	assert.Equal(t, err.(ServiceError).GetOpcRequestID(), opcID)
	assert.Equal(t, strings.Contains(err.Error(), "Service error:"), true)

	failure, ok := IsServiceError(err)
	assert.Equal(t, ok, true)
	assert.Equal(t, failure.GetHTTPStatusCode(), httpResponse.StatusCode)
}

func TestErrors_FailedToParseJson(t *testing.T) {
	// invalid json
	sampleResponse := `{"key" : test"}`

	httpResponse := http.Response{}
	bodyBuffer := bytes.NewBufferString(sampleResponse)
	httpResponse.Body = ioutil.NopCloser(bodyBuffer)
	err := newServiceFailureFromResponse(&httpResponse)

	failure, ok := IsServiceError(err)
	assert.Equal(t, ok, true)
	assert.Equal(t, failure.GetCode(), "BadErrorResponse")
	assert.Equal(t, strings.Contains(failure.GetMessage(), "Failed to parse json from response body due to"), true)
}

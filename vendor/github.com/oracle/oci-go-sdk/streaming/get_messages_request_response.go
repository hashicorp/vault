// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package streaming

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetMessagesRequest wrapper for the GetMessages operation
type GetMessagesRequest struct {

	// The OCID of the stream to get messages from.
	StreamId *string `mandatory:"true" contributesTo:"path" name:"streamId"`

	// The cursor used to consume the stream.
	Cursor *string `mandatory:"true" contributesTo:"query" name:"cursor"`

	// The maximum number of messages to return. You can specify any value up to 10000. By default, the service returns as many messages as possible.
	// Consider your average message size to help avoid exceeding throughput on the stream.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetMessagesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetMessagesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetMessagesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetMessagesResponse wrapper for the GetMessages operation
type GetMessagesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The []Message instance
	Items []Message `presentIn:"body"`

	// The cursor to use to get the next batch of messages.
	OpcNextCursor *string `presentIn:"header" name:"opc-next-cursor"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response GetMessagesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetMessagesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

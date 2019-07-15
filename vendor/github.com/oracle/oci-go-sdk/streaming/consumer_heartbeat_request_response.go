// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package streaming

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ConsumerHeartbeatRequest wrapper for the ConsumerHeartbeat operation
type ConsumerHeartbeatRequest struct {

	// The OCID of the stream for which the group is committing offsets.
	StreamId *string `mandatory:"true" contributesTo:"path" name:"streamId"`

	// The group-cursor representing the offsets of the group. This cursor is retrieved from the CreateGroupCursor API call.
	Cursor *string `mandatory:"true" contributesTo:"query" name:"cursor"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ConsumerHeartbeatRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ConsumerHeartbeatRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ConsumerHeartbeatRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ConsumerHeartbeatResponse wrapper for the ConsumerHeartbeat operation
type ConsumerHeartbeatResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The Cursor instance
	Cursor `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ConsumerHeartbeatResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ConsumerHeartbeatResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package containerengine

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetNodePoolOptionsRequest wrapper for the GetNodePoolOptions operation
type GetNodePoolOptionsRequest struct {

	// The id of the option set to retrieve. Use "all" get all options, or use a cluster ID to get options specific to the provided cluster.
	NodePoolOptionId *string `mandatory:"true" contributesTo:"path" name:"nodePoolOptionId"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetNodePoolOptionsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetNodePoolOptionsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetNodePoolOptionsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetNodePoolOptionsResponse wrapper for the GetNodePoolOptions operation
type GetNodePoolOptionsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The NodePoolOptions instance
	NodePoolOptions `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response GetNodePoolOptionsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetNodePoolOptionsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

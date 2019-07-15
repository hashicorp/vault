// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package resourcesearch

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetResourceTypeRequest wrapper for the GetResourceType operation
type GetResourceTypeRequest struct {

	// The name of the resource type.
	Name *string `mandatory:"true" contributesTo:"path" name:"name"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular
	// request, please provide the complete request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetResourceTypeRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetResourceTypeRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetResourceTypeRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetResourceTypeResponse wrapper for the GetResourceType operation
type GetResourceTypeResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The ResourceType instance
	ResourceType `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response GetResourceTypeResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetResourceTypeResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

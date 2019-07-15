// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package resourcesearch

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// SearchResourcesRequest wrapper for the SearchResources operation
type SearchResourcesRequest struct {

	// Request parameters that describe query criteria.
	SearchDetails `contributesTo:"body"`

	// The maximum number of items to return. The value must be between 1 and 1000.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The page at which to start retrieving results.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular
	// request, please provide the complete request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request SearchResourcesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request SearchResourcesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request SearchResourcesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// SearchResourcesResponse wrapper for the SearchResources operation
type SearchResourcesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of ResourceSummaryCollection instances
	ResourceSummaryCollection `presentIn:"body"`

	// Pagination token
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response SearchResourcesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response SearchResourcesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

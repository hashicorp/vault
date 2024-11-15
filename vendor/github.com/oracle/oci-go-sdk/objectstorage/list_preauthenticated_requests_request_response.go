// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListPreauthenticatedRequestsRequest wrapper for the ListPreauthenticatedRequests operation
type ListPreauthenticatedRequestsRequest struct {

	// The Object Storage namespace used for the request.
	NamespaceName *string `mandatory:"true" contributesTo:"path" name:"namespaceName"`

	// The name of the bucket. Avoid entering confidential information.
	// Example: `my-new-bucket1`
	BucketName *string `mandatory:"true" contributesTo:"path" name:"bucketName"`

	// User-specified object name prefixes can be used to query and return a list of pre-authenticated requests.
	ObjectNamePrefix *string `mandatory:"false" contributesTo:"query" name:"objectNamePrefix"`

	// The maximum number of items to return.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The page at which to start retrieving results.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The client request ID for tracing.
	OpcClientRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-client-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListPreauthenticatedRequestsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListPreauthenticatedRequestsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListPreauthenticatedRequestsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListPreauthenticatedRequestsResponse wrapper for the ListPreauthenticatedRequests operation
type ListPreauthenticatedRequestsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []PreauthenticatedRequestSummary instances
	Items []PreauthenticatedRequestSummary `presentIn:"body"`

	// Echoes back the value passed in the opc-client-request-id header, for use by clients when debugging.
	OpcClientRequestId *string `presentIn:"header" name:"opc-client-request-id"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular
	// request, provide this request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// Paginating a list of pre-authenticated requests.
	// In the GET request, set the limit to the number of pre-authenticated requests that you want returned in
	// the response. If the opc-next-page header appears in the response, then this is a partial list and there
	// are additional pre-authenticated requests to get. Include the header's value as the `page` parameter in
	// the subsequent GET request to get the next batch of pre-authenticated requests. Repeat this process to
	// retrieve the entire list of pre-authenticated requests.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListPreauthenticatedRequestsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListPreauthenticatedRequestsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

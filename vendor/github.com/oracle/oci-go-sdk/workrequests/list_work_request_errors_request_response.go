// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package workrequests

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListWorkRequestErrorsRequest wrapper for the ListWorkRequestErrors operation
type ListWorkRequestErrorsRequest struct {

	// The OCID of the work request.
	WorkRequestId *string `mandatory:"true" contributesTo:"path" name:"workRequestId"`

	// The maximum number of items to return in a paginated "List" call.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header from the previous "List" call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListWorkRequestErrorsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListWorkRequestErrorsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListWorkRequestErrorsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListWorkRequestErrorsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListWorkRequestErrorsResponse wrapper for the ListWorkRequestErrors operation
type ListWorkRequestErrorsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []WorkRequestError instances
	Items []WorkRequestError `presentIn:"body"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there might be additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListWorkRequestErrorsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListWorkRequestErrorsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListWorkRequestErrorsSortOrderEnum Enum with underlying type: string
type ListWorkRequestErrorsSortOrderEnum string

// Set of constants representing the allowable values for ListWorkRequestErrorsSortOrderEnum
const (
	ListWorkRequestErrorsSortOrderAsc  ListWorkRequestErrorsSortOrderEnum = "ASC"
	ListWorkRequestErrorsSortOrderDesc ListWorkRequestErrorsSortOrderEnum = "DESC"
)

var mappingListWorkRequestErrorsSortOrder = map[string]ListWorkRequestErrorsSortOrderEnum{
	"ASC":  ListWorkRequestErrorsSortOrderAsc,
	"DESC": ListWorkRequestErrorsSortOrderDesc,
}

// GetListWorkRequestErrorsSortOrderEnumValues Enumerates the set of values for ListWorkRequestErrorsSortOrderEnum
func GetListWorkRequestErrorsSortOrderEnumValues() []ListWorkRequestErrorsSortOrderEnum {
	values := make([]ListWorkRequestErrorsSortOrderEnum, 0)
	for _, v := range mappingListWorkRequestErrorsSortOrder {
		values = append(values, v)
	}
	return values
}

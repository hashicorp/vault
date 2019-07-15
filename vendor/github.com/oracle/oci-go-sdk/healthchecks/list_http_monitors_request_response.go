// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package healthchecks

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListHttpMonitorsRequest wrapper for the ListHttpMonitors operation
type ListHttpMonitorsRequest struct {

	// Filters results by compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The maximum number of items to return in a paginated "List" call.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header
	// from the previous "List" call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by when listing monitors.
	SortBy ListHttpMonitorsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// Controls the sort order of results.
	SortOrder ListHttpMonitorsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Filters results that exactly match the `displayName` field.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListHttpMonitorsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListHttpMonitorsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListHttpMonitorsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListHttpMonitorsResponse wrapper for the ListHttpMonitors operation
type ListHttpMonitorsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []HttpMonitorSummary instances
	Items []HttpMonitorSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to
	// contact Oracle about a particular request, please provide
	// the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this
	// header appears in the response, then a partial list might have been
	// returned. Include this value as the `page` parameter for the
	// subsequent GET request to get the next batch of items.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListHttpMonitorsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListHttpMonitorsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListHttpMonitorsSortByEnum Enum with underlying type: string
type ListHttpMonitorsSortByEnum string

// Set of constants representing the allowable values for ListHttpMonitorsSortByEnum
const (
	ListHttpMonitorsSortById          ListHttpMonitorsSortByEnum = "id"
	ListHttpMonitorsSortByDisplayname ListHttpMonitorsSortByEnum = "displayName"
)

var mappingListHttpMonitorsSortBy = map[string]ListHttpMonitorsSortByEnum{
	"id":          ListHttpMonitorsSortById,
	"displayName": ListHttpMonitorsSortByDisplayname,
}

// GetListHttpMonitorsSortByEnumValues Enumerates the set of values for ListHttpMonitorsSortByEnum
func GetListHttpMonitorsSortByEnumValues() []ListHttpMonitorsSortByEnum {
	values := make([]ListHttpMonitorsSortByEnum, 0)
	for _, v := range mappingListHttpMonitorsSortBy {
		values = append(values, v)
	}
	return values
}

// ListHttpMonitorsSortOrderEnum Enum with underlying type: string
type ListHttpMonitorsSortOrderEnum string

// Set of constants representing the allowable values for ListHttpMonitorsSortOrderEnum
const (
	ListHttpMonitorsSortOrderAsc  ListHttpMonitorsSortOrderEnum = "ASC"
	ListHttpMonitorsSortOrderDesc ListHttpMonitorsSortOrderEnum = "DESC"
)

var mappingListHttpMonitorsSortOrder = map[string]ListHttpMonitorsSortOrderEnum{
	"ASC":  ListHttpMonitorsSortOrderAsc,
	"DESC": ListHttpMonitorsSortOrderDesc,
}

// GetListHttpMonitorsSortOrderEnumValues Enumerates the set of values for ListHttpMonitorsSortOrderEnum
func GetListHttpMonitorsSortOrderEnumValues() []ListHttpMonitorsSortOrderEnum {
	values := make([]ListHttpMonitorsSortOrderEnum, 0)
	for _, v := range mappingListHttpMonitorsSortOrder {
		values = append(values, v)
	}
	return values
}

// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package healthchecks

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListPingMonitorsRequest wrapper for the ListPingMonitors operation
type ListPingMonitorsRequest struct {

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
	SortBy ListPingMonitorsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// Controls the sort order of results.
	SortOrder ListPingMonitorsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Filters results that exactly match the `displayName` field.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListPingMonitorsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListPingMonitorsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListPingMonitorsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListPingMonitorsResponse wrapper for the ListPingMonitors operation
type ListPingMonitorsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []PingMonitorSummary instances
	Items []PingMonitorSummary `presentIn:"body"`

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

func (response ListPingMonitorsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListPingMonitorsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListPingMonitorsSortByEnum Enum with underlying type: string
type ListPingMonitorsSortByEnum string

// Set of constants representing the allowable values for ListPingMonitorsSortByEnum
const (
	ListPingMonitorsSortById          ListPingMonitorsSortByEnum = "id"
	ListPingMonitorsSortByDisplayname ListPingMonitorsSortByEnum = "displayName"
)

var mappingListPingMonitorsSortBy = map[string]ListPingMonitorsSortByEnum{
	"id":          ListPingMonitorsSortById,
	"displayName": ListPingMonitorsSortByDisplayname,
}

// GetListPingMonitorsSortByEnumValues Enumerates the set of values for ListPingMonitorsSortByEnum
func GetListPingMonitorsSortByEnumValues() []ListPingMonitorsSortByEnum {
	values := make([]ListPingMonitorsSortByEnum, 0)
	for _, v := range mappingListPingMonitorsSortBy {
		values = append(values, v)
	}
	return values
}

// ListPingMonitorsSortOrderEnum Enum with underlying type: string
type ListPingMonitorsSortOrderEnum string

// Set of constants representing the allowable values for ListPingMonitorsSortOrderEnum
const (
	ListPingMonitorsSortOrderAsc  ListPingMonitorsSortOrderEnum = "ASC"
	ListPingMonitorsSortOrderDesc ListPingMonitorsSortOrderEnum = "DESC"
)

var mappingListPingMonitorsSortOrder = map[string]ListPingMonitorsSortOrderEnum{
	"ASC":  ListPingMonitorsSortOrderAsc,
	"DESC": ListPingMonitorsSortOrderDesc,
}

// GetListPingMonitorsSortOrderEnumValues Enumerates the set of values for ListPingMonitorsSortOrderEnum
func GetListPingMonitorsSortOrderEnumValues() []ListPingMonitorsSortOrderEnum {
	values := make([]ListPingMonitorsSortOrderEnum, 0)
	for _, v := range mappingListPingMonitorsSortOrder {
		values = append(values, v)
	}
	return values
}

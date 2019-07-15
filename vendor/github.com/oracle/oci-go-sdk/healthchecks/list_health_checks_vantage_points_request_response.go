// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package healthchecks

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListHealthChecksVantagePointsRequest wrapper for the ListHealthChecksVantagePoints operation
type ListHealthChecksVantagePointsRequest struct {

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The maximum number of items to return in a paginated "List" call.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header
	// from the previous "List" call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by when listing vantage points.
	SortBy ListHealthChecksVantagePointsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// Controls the sort order of results.
	SortOrder ListHealthChecksVantagePointsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Filters results that exactly match the `name` field.
	Name *string `mandatory:"false" contributesTo:"query" name:"name"`

	// Filters results that exactly match the `displayName` field.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListHealthChecksVantagePointsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListHealthChecksVantagePointsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListHealthChecksVantagePointsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListHealthChecksVantagePointsResponse wrapper for the ListHealthChecksVantagePoints operation
type ListHealthChecksVantagePointsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []HealthChecksVantagePointSummary instances
	Items []HealthChecksVantagePointSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to
	// contact Oracle about a particular request, please provide
	// the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if
	// this header appears in the response, then there may be additional
	// items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#List_Pagination).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListHealthChecksVantagePointsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListHealthChecksVantagePointsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListHealthChecksVantagePointsSortByEnum Enum with underlying type: string
type ListHealthChecksVantagePointsSortByEnum string

// Set of constants representing the allowable values for ListHealthChecksVantagePointsSortByEnum
const (
	ListHealthChecksVantagePointsSortByName        ListHealthChecksVantagePointsSortByEnum = "name"
	ListHealthChecksVantagePointsSortByDisplayname ListHealthChecksVantagePointsSortByEnum = "displayName"
)

var mappingListHealthChecksVantagePointsSortBy = map[string]ListHealthChecksVantagePointsSortByEnum{
	"name":        ListHealthChecksVantagePointsSortByName,
	"displayName": ListHealthChecksVantagePointsSortByDisplayname,
}

// GetListHealthChecksVantagePointsSortByEnumValues Enumerates the set of values for ListHealthChecksVantagePointsSortByEnum
func GetListHealthChecksVantagePointsSortByEnumValues() []ListHealthChecksVantagePointsSortByEnum {
	values := make([]ListHealthChecksVantagePointsSortByEnum, 0)
	for _, v := range mappingListHealthChecksVantagePointsSortBy {
		values = append(values, v)
	}
	return values
}

// ListHealthChecksVantagePointsSortOrderEnum Enum with underlying type: string
type ListHealthChecksVantagePointsSortOrderEnum string

// Set of constants representing the allowable values for ListHealthChecksVantagePointsSortOrderEnum
const (
	ListHealthChecksVantagePointsSortOrderAsc  ListHealthChecksVantagePointsSortOrderEnum = "ASC"
	ListHealthChecksVantagePointsSortOrderDesc ListHealthChecksVantagePointsSortOrderEnum = "DESC"
)

var mappingListHealthChecksVantagePointsSortOrder = map[string]ListHealthChecksVantagePointsSortOrderEnum{
	"ASC":  ListHealthChecksVantagePointsSortOrderAsc,
	"DESC": ListHealthChecksVantagePointsSortOrderDesc,
}

// GetListHealthChecksVantagePointsSortOrderEnumValues Enumerates the set of values for ListHealthChecksVantagePointsSortOrderEnum
func GetListHealthChecksVantagePointsSortOrderEnumValues() []ListHealthChecksVantagePointsSortOrderEnum {
	values := make([]ListHealthChecksVantagePointsSortOrderEnum, 0)
	for _, v := range mappingListHealthChecksVantagePointsSortOrder {
		values = append(values, v)
	}
	return values
}

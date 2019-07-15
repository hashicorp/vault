// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAutonomousDbPreviewVersionsRequest wrapper for the ListAutonomousDbPreviewVersions operation
type ListAutonomousDbPreviewVersionsRequest struct {

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Unique identifier for the request.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The field to sort by.  You can provide one sort order (`sortOrder`).  Default order for DBWORKLOAD is ascending.
	// **Note:** If you do not include the availability domain filter, the resources are grouped by availability domain, then sorted.
	SortBy ListAutonomousDbPreviewVersionsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListAutonomousDbPreviewVersionsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAutonomousDbPreviewVersionsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAutonomousDbPreviewVersionsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAutonomousDbPreviewVersionsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAutonomousDbPreviewVersionsResponse wrapper for the ListAutonomousDbPreviewVersions operation
type ListAutonomousDbPreviewVersionsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AutonomousDbPreviewVersionSummary instances
	Items []AutonomousDbPreviewVersionSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListAutonomousDbPreviewVersionsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAutonomousDbPreviewVersionsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAutonomousDbPreviewVersionsSortByEnum Enum with underlying type: string
type ListAutonomousDbPreviewVersionsSortByEnum string

// Set of constants representing the allowable values for ListAutonomousDbPreviewVersionsSortByEnum
const (
	ListAutonomousDbPreviewVersionsSortByDbworkload ListAutonomousDbPreviewVersionsSortByEnum = "DBWORKLOAD"
)

var mappingListAutonomousDbPreviewVersionsSortBy = map[string]ListAutonomousDbPreviewVersionsSortByEnum{
	"DBWORKLOAD": ListAutonomousDbPreviewVersionsSortByDbworkload,
}

// GetListAutonomousDbPreviewVersionsSortByEnumValues Enumerates the set of values for ListAutonomousDbPreviewVersionsSortByEnum
func GetListAutonomousDbPreviewVersionsSortByEnumValues() []ListAutonomousDbPreviewVersionsSortByEnum {
	values := make([]ListAutonomousDbPreviewVersionsSortByEnum, 0)
	for _, v := range mappingListAutonomousDbPreviewVersionsSortBy {
		values = append(values, v)
	}
	return values
}

// ListAutonomousDbPreviewVersionsSortOrderEnum Enum with underlying type: string
type ListAutonomousDbPreviewVersionsSortOrderEnum string

// Set of constants representing the allowable values for ListAutonomousDbPreviewVersionsSortOrderEnum
const (
	ListAutonomousDbPreviewVersionsSortOrderAsc  ListAutonomousDbPreviewVersionsSortOrderEnum = "ASC"
	ListAutonomousDbPreviewVersionsSortOrderDesc ListAutonomousDbPreviewVersionsSortOrderEnum = "DESC"
)

var mappingListAutonomousDbPreviewVersionsSortOrder = map[string]ListAutonomousDbPreviewVersionsSortOrderEnum{
	"ASC":  ListAutonomousDbPreviewVersionsSortOrderAsc,
	"DESC": ListAutonomousDbPreviewVersionsSortOrderDesc,
}

// GetListAutonomousDbPreviewVersionsSortOrderEnumValues Enumerates the set of values for ListAutonomousDbPreviewVersionsSortOrderEnum
func GetListAutonomousDbPreviewVersionsSortOrderEnumValues() []ListAutonomousDbPreviewVersionsSortOrderEnum {
	values := make([]ListAutonomousDbPreviewVersionsSortOrderEnum, 0)
	for _, v := range mappingListAutonomousDbPreviewVersionsSortOrder {
		values = append(values, v)
	}
	return values
}

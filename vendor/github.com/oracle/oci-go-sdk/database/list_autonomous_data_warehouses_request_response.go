// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAutonomousDataWarehousesRequest wrapper for the ListAutonomousDataWarehouses operation
type ListAutonomousDataWarehousesRequest struct {

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by.  You can provide one sort order (`sortOrder`).  Default order for TIMECREATED is descending.  Default order for DISPLAYNAME is ascending. The DISPLAYNAME sort order is case sensitive.
	// **Note:** If you do not include the availability domain filter, the resources are grouped by availability domain, then sorted.
	SortBy ListAutonomousDataWarehousesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListAutonomousDataWarehousesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the given lifecycle state exactly.
	LifecycleState AutonomousDataWarehouseSummaryLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// A filter to return only resources that match the entire display name given. The match is not case sensitive.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAutonomousDataWarehousesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAutonomousDataWarehousesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAutonomousDataWarehousesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAutonomousDataWarehousesResponse wrapper for the ListAutonomousDataWarehouses operation
type ListAutonomousDataWarehousesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AutonomousDataWarehouseSummary instances
	Items []AutonomousDataWarehouseSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListAutonomousDataWarehousesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAutonomousDataWarehousesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAutonomousDataWarehousesSortByEnum Enum with underlying type: string
type ListAutonomousDataWarehousesSortByEnum string

// Set of constants representing the allowable values for ListAutonomousDataWarehousesSortByEnum
const (
	ListAutonomousDataWarehousesSortByTimecreated ListAutonomousDataWarehousesSortByEnum = "TIMECREATED"
	ListAutonomousDataWarehousesSortByDisplayname ListAutonomousDataWarehousesSortByEnum = "DISPLAYNAME"
)

var mappingListAutonomousDataWarehousesSortBy = map[string]ListAutonomousDataWarehousesSortByEnum{
	"TIMECREATED": ListAutonomousDataWarehousesSortByTimecreated,
	"DISPLAYNAME": ListAutonomousDataWarehousesSortByDisplayname,
}

// GetListAutonomousDataWarehousesSortByEnumValues Enumerates the set of values for ListAutonomousDataWarehousesSortByEnum
func GetListAutonomousDataWarehousesSortByEnumValues() []ListAutonomousDataWarehousesSortByEnum {
	values := make([]ListAutonomousDataWarehousesSortByEnum, 0)
	for _, v := range mappingListAutonomousDataWarehousesSortBy {
		values = append(values, v)
	}
	return values
}

// ListAutonomousDataWarehousesSortOrderEnum Enum with underlying type: string
type ListAutonomousDataWarehousesSortOrderEnum string

// Set of constants representing the allowable values for ListAutonomousDataWarehousesSortOrderEnum
const (
	ListAutonomousDataWarehousesSortOrderAsc  ListAutonomousDataWarehousesSortOrderEnum = "ASC"
	ListAutonomousDataWarehousesSortOrderDesc ListAutonomousDataWarehousesSortOrderEnum = "DESC"
)

var mappingListAutonomousDataWarehousesSortOrder = map[string]ListAutonomousDataWarehousesSortOrderEnum{
	"ASC":  ListAutonomousDataWarehousesSortOrderAsc,
	"DESC": ListAutonomousDataWarehousesSortOrderDesc,
}

// GetListAutonomousDataWarehousesSortOrderEnumValues Enumerates the set of values for ListAutonomousDataWarehousesSortOrderEnum
func GetListAutonomousDataWarehousesSortOrderEnumValues() []ListAutonomousDataWarehousesSortOrderEnum {
	values := make([]ListAutonomousDataWarehousesSortOrderEnum, 0)
	for _, v := range mappingListAutonomousDataWarehousesSortOrder {
		values = append(values, v)
	}
	return values
}

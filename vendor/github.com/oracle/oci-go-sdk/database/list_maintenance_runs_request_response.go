// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListMaintenanceRunsRequest wrapper for the ListMaintenanceRuns operation
type ListMaintenanceRunsRequest struct {

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The target resource ID.
	TargetResourceId *string `mandatory:"false" contributesTo:"query" name:"targetResourceId"`

	// The type of the target resource.
	TargetResourceType MaintenanceRunSummaryTargetResourceTypeEnum `mandatory:"false" contributesTo:"query" name:"targetResourceType" omitEmpty:"true"`

	// The maintenance type.
	MaintenanceType MaintenanceRunSummaryMaintenanceTypeEnum `mandatory:"false" contributesTo:"query" name:"maintenanceType" omitEmpty:"true"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by.  You can provide one sort order (`sortOrder`).  Default order for TIME_SCHEDULED and TIME_ENDED is descending. Default order for DISPLAYNAME is ascending. The DISPLAYNAME sort order is case sensitive.
	// **Note:** If you do not include the availability domain filter, the resources are grouped by availability domain, then sorted.
	SortBy ListMaintenanceRunsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListMaintenanceRunsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the given lifecycle state exactly.
	LifecycleState MaintenanceRunSummaryLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// A filter to return only resources that match the given availability domain exactly.
	AvailabilityDomain *string `mandatory:"false" contributesTo:"query" name:"availabilityDomain"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListMaintenanceRunsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListMaintenanceRunsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListMaintenanceRunsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListMaintenanceRunsResponse wrapper for the ListMaintenanceRuns operation
type ListMaintenanceRunsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []MaintenanceRunSummary instances
	Items []MaintenanceRunSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListMaintenanceRunsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListMaintenanceRunsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListMaintenanceRunsSortByEnum Enum with underlying type: string
type ListMaintenanceRunsSortByEnum string

// Set of constants representing the allowable values for ListMaintenanceRunsSortByEnum
const (
	ListMaintenanceRunsSortByTimeScheduled ListMaintenanceRunsSortByEnum = "TIME_SCHEDULED"
	ListMaintenanceRunsSortByTimeEnded     ListMaintenanceRunsSortByEnum = "TIME_ENDED"
	ListMaintenanceRunsSortByDisplayname   ListMaintenanceRunsSortByEnum = "DISPLAYNAME"
)

var mappingListMaintenanceRunsSortBy = map[string]ListMaintenanceRunsSortByEnum{
	"TIME_SCHEDULED": ListMaintenanceRunsSortByTimeScheduled,
	"TIME_ENDED":     ListMaintenanceRunsSortByTimeEnded,
	"DISPLAYNAME":    ListMaintenanceRunsSortByDisplayname,
}

// GetListMaintenanceRunsSortByEnumValues Enumerates the set of values for ListMaintenanceRunsSortByEnum
func GetListMaintenanceRunsSortByEnumValues() []ListMaintenanceRunsSortByEnum {
	values := make([]ListMaintenanceRunsSortByEnum, 0)
	for _, v := range mappingListMaintenanceRunsSortBy {
		values = append(values, v)
	}
	return values
}

// ListMaintenanceRunsSortOrderEnum Enum with underlying type: string
type ListMaintenanceRunsSortOrderEnum string

// Set of constants representing the allowable values for ListMaintenanceRunsSortOrderEnum
const (
	ListMaintenanceRunsSortOrderAsc  ListMaintenanceRunsSortOrderEnum = "ASC"
	ListMaintenanceRunsSortOrderDesc ListMaintenanceRunsSortOrderEnum = "DESC"
)

var mappingListMaintenanceRunsSortOrder = map[string]ListMaintenanceRunsSortOrderEnum{
	"ASC":  ListMaintenanceRunsSortOrderAsc,
	"DESC": ListMaintenanceRunsSortOrderDesc,
}

// GetListMaintenanceRunsSortOrderEnumValues Enumerates the set of values for ListMaintenanceRunsSortOrderEnum
func GetListMaintenanceRunsSortOrderEnumValues() []ListMaintenanceRunsSortOrderEnum {
	values := make([]ListMaintenanceRunsSortOrderEnum, 0)
	for _, v := range mappingListMaintenanceRunsSortOrder {
		values = append(values, v)
	}
	return values
}

// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAutonomousContainerDatabasesRequest wrapper for the ListAutonomousContainerDatabases operation
type ListAutonomousContainerDatabasesRequest struct {

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The Autonomous Exadata Infrastructure OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	AutonomousExadataInfrastructureId *string `mandatory:"false" contributesTo:"query" name:"autonomousExadataInfrastructureId"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by.  You can provide one sort order (`sortOrder`).  Default order for TIMECREATED is descending.  Default order for DISPLAYNAME is ascending. The DISPLAYNAME sort order is case sensitive.
	// **Note:** If you do not include the availability domain filter, the resources are grouped by availability domain, then sorted.
	SortBy ListAutonomousContainerDatabasesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListAutonomousContainerDatabasesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the given lifecycle state exactly.
	LifecycleState AutonomousContainerDatabaseSummaryLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// A filter to return only resources that match the given availability domain exactly.
	AvailabilityDomain *string `mandatory:"false" contributesTo:"query" name:"availabilityDomain"`

	// A filter to return only resources that match the entire display name given. The match is not case sensitive.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAutonomousContainerDatabasesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAutonomousContainerDatabasesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAutonomousContainerDatabasesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAutonomousContainerDatabasesResponse wrapper for the ListAutonomousContainerDatabases operation
type ListAutonomousContainerDatabasesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AutonomousContainerDatabaseSummary instances
	Items []AutonomousContainerDatabaseSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListAutonomousContainerDatabasesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAutonomousContainerDatabasesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAutonomousContainerDatabasesSortByEnum Enum with underlying type: string
type ListAutonomousContainerDatabasesSortByEnum string

// Set of constants representing the allowable values for ListAutonomousContainerDatabasesSortByEnum
const (
	ListAutonomousContainerDatabasesSortByTimecreated ListAutonomousContainerDatabasesSortByEnum = "TIMECREATED"
	ListAutonomousContainerDatabasesSortByDisplayname ListAutonomousContainerDatabasesSortByEnum = "DISPLAYNAME"
)

var mappingListAutonomousContainerDatabasesSortBy = map[string]ListAutonomousContainerDatabasesSortByEnum{
	"TIMECREATED": ListAutonomousContainerDatabasesSortByTimecreated,
	"DISPLAYNAME": ListAutonomousContainerDatabasesSortByDisplayname,
}

// GetListAutonomousContainerDatabasesSortByEnumValues Enumerates the set of values for ListAutonomousContainerDatabasesSortByEnum
func GetListAutonomousContainerDatabasesSortByEnumValues() []ListAutonomousContainerDatabasesSortByEnum {
	values := make([]ListAutonomousContainerDatabasesSortByEnum, 0)
	for _, v := range mappingListAutonomousContainerDatabasesSortBy {
		values = append(values, v)
	}
	return values
}

// ListAutonomousContainerDatabasesSortOrderEnum Enum with underlying type: string
type ListAutonomousContainerDatabasesSortOrderEnum string

// Set of constants representing the allowable values for ListAutonomousContainerDatabasesSortOrderEnum
const (
	ListAutonomousContainerDatabasesSortOrderAsc  ListAutonomousContainerDatabasesSortOrderEnum = "ASC"
	ListAutonomousContainerDatabasesSortOrderDesc ListAutonomousContainerDatabasesSortOrderEnum = "DESC"
)

var mappingListAutonomousContainerDatabasesSortOrder = map[string]ListAutonomousContainerDatabasesSortOrderEnum{
	"ASC":  ListAutonomousContainerDatabasesSortOrderAsc,
	"DESC": ListAutonomousContainerDatabasesSortOrderDesc,
}

// GetListAutonomousContainerDatabasesSortOrderEnumValues Enumerates the set of values for ListAutonomousContainerDatabasesSortOrderEnum
func GetListAutonomousContainerDatabasesSortOrderEnumValues() []ListAutonomousContainerDatabasesSortOrderEnum {
	values := make([]ListAutonomousContainerDatabasesSortOrderEnum, 0)
	for _, v := range mappingListAutonomousContainerDatabasesSortOrder {
		values = append(values, v)
	}
	return values
}

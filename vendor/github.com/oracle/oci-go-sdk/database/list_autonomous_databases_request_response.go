// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAutonomousDatabasesRequest wrapper for the ListAutonomousDatabases operation
type ListAutonomousDatabasesRequest struct {

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The Autonomous Container Database OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	AutonomousContainerDatabaseId *string `mandatory:"false" contributesTo:"query" name:"autonomousContainerDatabaseId"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by.  You can provide one sort order (`sortOrder`).  Default order for TIMECREATED is descending.  Default order for DISPLAYNAME is ascending. The DISPLAYNAME sort order is case sensitive.
	// **Note:** If you do not include the availability domain filter, the resources are grouped by availability domain, then sorted.
	SortBy ListAutonomousDatabasesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListAutonomousDatabasesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the given lifecycle state exactly.
	LifecycleState AutonomousDatabaseSummaryLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// A filter to return only autonomous database resources that match the specified workload type.
	DbWorkload AutonomousDatabaseSummaryDbWorkloadEnum `mandatory:"false" contributesTo:"query" name:"dbWorkload" omitEmpty:"true"`

	// A filter to return only resources that match the entire display name given. The match is not case sensitive.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Unique identifier for the request.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAutonomousDatabasesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAutonomousDatabasesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAutonomousDatabasesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAutonomousDatabasesResponse wrapper for the ListAutonomousDatabases operation
type ListAutonomousDatabasesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AutonomousDatabaseSummary instances
	Items []AutonomousDatabaseSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListAutonomousDatabasesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAutonomousDatabasesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAutonomousDatabasesSortByEnum Enum with underlying type: string
type ListAutonomousDatabasesSortByEnum string

// Set of constants representing the allowable values for ListAutonomousDatabasesSortByEnum
const (
	ListAutonomousDatabasesSortByTimecreated ListAutonomousDatabasesSortByEnum = "TIMECREATED"
	ListAutonomousDatabasesSortByDisplayname ListAutonomousDatabasesSortByEnum = "DISPLAYNAME"
)

var mappingListAutonomousDatabasesSortBy = map[string]ListAutonomousDatabasesSortByEnum{
	"TIMECREATED": ListAutonomousDatabasesSortByTimecreated,
	"DISPLAYNAME": ListAutonomousDatabasesSortByDisplayname,
}

// GetListAutonomousDatabasesSortByEnumValues Enumerates the set of values for ListAutonomousDatabasesSortByEnum
func GetListAutonomousDatabasesSortByEnumValues() []ListAutonomousDatabasesSortByEnum {
	values := make([]ListAutonomousDatabasesSortByEnum, 0)
	for _, v := range mappingListAutonomousDatabasesSortBy {
		values = append(values, v)
	}
	return values
}

// ListAutonomousDatabasesSortOrderEnum Enum with underlying type: string
type ListAutonomousDatabasesSortOrderEnum string

// Set of constants representing the allowable values for ListAutonomousDatabasesSortOrderEnum
const (
	ListAutonomousDatabasesSortOrderAsc  ListAutonomousDatabasesSortOrderEnum = "ASC"
	ListAutonomousDatabasesSortOrderDesc ListAutonomousDatabasesSortOrderEnum = "DESC"
)

var mappingListAutonomousDatabasesSortOrder = map[string]ListAutonomousDatabasesSortOrderEnum{
	"ASC":  ListAutonomousDatabasesSortOrderAsc,
	"DESC": ListAutonomousDatabasesSortOrderDesc,
}

// GetListAutonomousDatabasesSortOrderEnumValues Enumerates the set of values for ListAutonomousDatabasesSortOrderEnum
func GetListAutonomousDatabasesSortOrderEnumValues() []ListAutonomousDatabasesSortOrderEnum {
	values := make([]ListAutonomousDatabasesSortOrderEnum, 0)
	for _, v := range mappingListAutonomousDatabasesSortOrder {
		values = append(values, v)
	}
	return values
}

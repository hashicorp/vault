// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListDatabasesRequest wrapper for the ListDatabases operation
type ListDatabasesRequest struct {

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// A database home OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	DbHomeId *string `mandatory:"true" contributesTo:"query" name:"dbHomeId"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by.  You can provide one sort order (`sortOrder`).  Default order for TIMECREATED is descending.  Default order for DBNAME is ascending. The DBNAME sort order is case sensitive.
	SortBy ListDatabasesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListDatabasesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the given lifecycle state exactly.
	LifecycleState DatabaseSummaryLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// A filter to return only resources that match the entire database name given. The match is not case sensitive.
	DbName *string `mandatory:"false" contributesTo:"query" name:"dbName"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListDatabasesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListDatabasesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListDatabasesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListDatabasesResponse wrapper for the ListDatabases operation
type ListDatabasesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []DatabaseSummary instances
	Items []DatabaseSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListDatabasesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListDatabasesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListDatabasesSortByEnum Enum with underlying type: string
type ListDatabasesSortByEnum string

// Set of constants representing the allowable values for ListDatabasesSortByEnum
const (
	ListDatabasesSortByDbname      ListDatabasesSortByEnum = "DBNAME"
	ListDatabasesSortByTimecreated ListDatabasesSortByEnum = "TIMECREATED"
)

var mappingListDatabasesSortBy = map[string]ListDatabasesSortByEnum{
	"DBNAME":      ListDatabasesSortByDbname,
	"TIMECREATED": ListDatabasesSortByTimecreated,
}

// GetListDatabasesSortByEnumValues Enumerates the set of values for ListDatabasesSortByEnum
func GetListDatabasesSortByEnumValues() []ListDatabasesSortByEnum {
	values := make([]ListDatabasesSortByEnum, 0)
	for _, v := range mappingListDatabasesSortBy {
		values = append(values, v)
	}
	return values
}

// ListDatabasesSortOrderEnum Enum with underlying type: string
type ListDatabasesSortOrderEnum string

// Set of constants representing the allowable values for ListDatabasesSortOrderEnum
const (
	ListDatabasesSortOrderAsc  ListDatabasesSortOrderEnum = "ASC"
	ListDatabasesSortOrderDesc ListDatabasesSortOrderEnum = "DESC"
)

var mappingListDatabasesSortOrder = map[string]ListDatabasesSortOrderEnum{
	"ASC":  ListDatabasesSortOrderAsc,
	"DESC": ListDatabasesSortOrderDesc,
}

// GetListDatabasesSortOrderEnumValues Enumerates the set of values for ListDatabasesSortOrderEnum
func GetListDatabasesSortOrderEnumValues() []ListDatabasesSortOrderEnum {
	values := make([]ListDatabasesSortOrderEnum, 0)
	for _, v := range mappingListDatabasesSortOrder {
		values = append(values, v)
	}
	return values
}

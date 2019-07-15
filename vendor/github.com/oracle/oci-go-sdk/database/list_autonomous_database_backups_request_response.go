// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAutonomousDatabaseBackupsRequest wrapper for the ListAutonomousDatabaseBackups operation
type ListAutonomousDatabaseBackupsRequest struct {

	// The database OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	AutonomousDatabaseId *string `mandatory:"false" contributesTo:"query" name:"autonomousDatabaseId"`

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"false" contributesTo:"query" name:"compartmentId"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by.  You can provide one sort order (`sortOrder`).  Default order for TIMECREATED is descending.  Default order for DISPLAYNAME is ascending. The DISPLAYNAME sort order is case sensitive.
	// **Note:** If you do not include the availability domain filter, the resources are grouped by availability domain, then sorted.
	SortBy ListAutonomousDatabaseBackupsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListAutonomousDatabaseBackupsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the given lifecycle state exactly.
	LifecycleState AutonomousDatabaseBackupSummaryLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// A filter to return only resources that match the entire display name given. The match is not case sensitive.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Unique identifier for the request.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAutonomousDatabaseBackupsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAutonomousDatabaseBackupsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAutonomousDatabaseBackupsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAutonomousDatabaseBackupsResponse wrapper for the ListAutonomousDatabaseBackups operation
type ListAutonomousDatabaseBackupsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AutonomousDatabaseBackupSummary instances
	Items []AutonomousDatabaseBackupSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListAutonomousDatabaseBackupsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAutonomousDatabaseBackupsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAutonomousDatabaseBackupsSortByEnum Enum with underlying type: string
type ListAutonomousDatabaseBackupsSortByEnum string

// Set of constants representing the allowable values for ListAutonomousDatabaseBackupsSortByEnum
const (
	ListAutonomousDatabaseBackupsSortByTimecreated ListAutonomousDatabaseBackupsSortByEnum = "TIMECREATED"
	ListAutonomousDatabaseBackupsSortByDisplayname ListAutonomousDatabaseBackupsSortByEnum = "DISPLAYNAME"
)

var mappingListAutonomousDatabaseBackupsSortBy = map[string]ListAutonomousDatabaseBackupsSortByEnum{
	"TIMECREATED": ListAutonomousDatabaseBackupsSortByTimecreated,
	"DISPLAYNAME": ListAutonomousDatabaseBackupsSortByDisplayname,
}

// GetListAutonomousDatabaseBackupsSortByEnumValues Enumerates the set of values for ListAutonomousDatabaseBackupsSortByEnum
func GetListAutonomousDatabaseBackupsSortByEnumValues() []ListAutonomousDatabaseBackupsSortByEnum {
	values := make([]ListAutonomousDatabaseBackupsSortByEnum, 0)
	for _, v := range mappingListAutonomousDatabaseBackupsSortBy {
		values = append(values, v)
	}
	return values
}

// ListAutonomousDatabaseBackupsSortOrderEnum Enum with underlying type: string
type ListAutonomousDatabaseBackupsSortOrderEnum string

// Set of constants representing the allowable values for ListAutonomousDatabaseBackupsSortOrderEnum
const (
	ListAutonomousDatabaseBackupsSortOrderAsc  ListAutonomousDatabaseBackupsSortOrderEnum = "ASC"
	ListAutonomousDatabaseBackupsSortOrderDesc ListAutonomousDatabaseBackupsSortOrderEnum = "DESC"
)

var mappingListAutonomousDatabaseBackupsSortOrder = map[string]ListAutonomousDatabaseBackupsSortOrderEnum{
	"ASC":  ListAutonomousDatabaseBackupsSortOrderAsc,
	"DESC": ListAutonomousDatabaseBackupsSortOrderDesc,
}

// GetListAutonomousDatabaseBackupsSortOrderEnumValues Enumerates the set of values for ListAutonomousDatabaseBackupsSortOrderEnum
func GetListAutonomousDatabaseBackupsSortOrderEnumValues() []ListAutonomousDatabaseBackupsSortOrderEnum {
	values := make([]ListAutonomousDatabaseBackupsSortOrderEnum, 0)
	for _, v := range mappingListAutonomousDatabaseBackupsSortOrder {
		values = append(values, v)
	}
	return values
}

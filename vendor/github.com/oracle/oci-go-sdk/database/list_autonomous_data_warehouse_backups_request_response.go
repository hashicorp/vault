// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAutonomousDataWarehouseBackupsRequest wrapper for the ListAutonomousDataWarehouseBackups operation
type ListAutonomousDataWarehouseBackupsRequest struct {

	// The database OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	AutonomousDataWarehouseId *string `mandatory:"false" contributesTo:"query" name:"autonomousDataWarehouseId"`

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"false" contributesTo:"query" name:"compartmentId"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by.  You can provide one sort order (`sortOrder`).  Default order for TIMECREATED is descending.  Default order for DISPLAYNAME is ascending. The DISPLAYNAME sort order is case sensitive.
	// **Note:** If you do not include the availability domain filter, the resources are grouped by availability domain, then sorted.
	SortBy ListAutonomousDataWarehouseBackupsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListAutonomousDataWarehouseBackupsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the given lifecycle state exactly.
	LifecycleState AutonomousDataWarehouseBackupSummaryLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// A filter to return only resources that match the entire display name given. The match is not case sensitive.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAutonomousDataWarehouseBackupsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAutonomousDataWarehouseBackupsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAutonomousDataWarehouseBackupsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAutonomousDataWarehouseBackupsResponse wrapper for the ListAutonomousDataWarehouseBackups operation
type ListAutonomousDataWarehouseBackupsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AutonomousDataWarehouseBackupSummary instances
	Items []AutonomousDataWarehouseBackupSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListAutonomousDataWarehouseBackupsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAutonomousDataWarehouseBackupsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAutonomousDataWarehouseBackupsSortByEnum Enum with underlying type: string
type ListAutonomousDataWarehouseBackupsSortByEnum string

// Set of constants representing the allowable values for ListAutonomousDataWarehouseBackupsSortByEnum
const (
	ListAutonomousDataWarehouseBackupsSortByTimecreated ListAutonomousDataWarehouseBackupsSortByEnum = "TIMECREATED"
	ListAutonomousDataWarehouseBackupsSortByDisplayname ListAutonomousDataWarehouseBackupsSortByEnum = "DISPLAYNAME"
)

var mappingListAutonomousDataWarehouseBackupsSortBy = map[string]ListAutonomousDataWarehouseBackupsSortByEnum{
	"TIMECREATED": ListAutonomousDataWarehouseBackupsSortByTimecreated,
	"DISPLAYNAME": ListAutonomousDataWarehouseBackupsSortByDisplayname,
}

// GetListAutonomousDataWarehouseBackupsSortByEnumValues Enumerates the set of values for ListAutonomousDataWarehouseBackupsSortByEnum
func GetListAutonomousDataWarehouseBackupsSortByEnumValues() []ListAutonomousDataWarehouseBackupsSortByEnum {
	values := make([]ListAutonomousDataWarehouseBackupsSortByEnum, 0)
	for _, v := range mappingListAutonomousDataWarehouseBackupsSortBy {
		values = append(values, v)
	}
	return values
}

// ListAutonomousDataWarehouseBackupsSortOrderEnum Enum with underlying type: string
type ListAutonomousDataWarehouseBackupsSortOrderEnum string

// Set of constants representing the allowable values for ListAutonomousDataWarehouseBackupsSortOrderEnum
const (
	ListAutonomousDataWarehouseBackupsSortOrderAsc  ListAutonomousDataWarehouseBackupsSortOrderEnum = "ASC"
	ListAutonomousDataWarehouseBackupsSortOrderDesc ListAutonomousDataWarehouseBackupsSortOrderEnum = "DESC"
)

var mappingListAutonomousDataWarehouseBackupsSortOrder = map[string]ListAutonomousDataWarehouseBackupsSortOrderEnum{
	"ASC":  ListAutonomousDataWarehouseBackupsSortOrderAsc,
	"DESC": ListAutonomousDataWarehouseBackupsSortOrderDesc,
}

// GetListAutonomousDataWarehouseBackupsSortOrderEnumValues Enumerates the set of values for ListAutonomousDataWarehouseBackupsSortOrderEnum
func GetListAutonomousDataWarehouseBackupsSortOrderEnumValues() []ListAutonomousDataWarehouseBackupsSortOrderEnum {
	values := make([]ListAutonomousDataWarehouseBackupsSortOrderEnum, 0)
	for _, v := range mappingListAutonomousDataWarehouseBackupsSortOrder {
		values = append(values, v)
	}
	return values
}

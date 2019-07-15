// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListDbSystemsRequest wrapper for the ListDbSystems operation
type ListDbSystemsRequest struct {

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the backup. Specify a backupId to list only the DB systems that support creating a database using this backup in this compartment.
	BackupId *string `mandatory:"false" contributesTo:"query" name:"backupId"`

	// The field to sort by.  You can provide one sort order (`sortOrder`).  Default order for TIMECREATED is descending.  Default order for DISPLAYNAME is ascending. The DISPLAYNAME sort order is case sensitive.
	// **Note:** If you do not include the availability domain filter, the resources are grouped by availability domain, then sorted.
	SortBy ListDbSystemsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListDbSystemsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the given lifecycle state exactly.
	LifecycleState DbSystemSummaryLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

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

func (request ListDbSystemsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListDbSystemsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListDbSystemsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListDbSystemsResponse wrapper for the ListDbSystems operation
type ListDbSystemsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []DbSystemSummary instances
	Items []DbSystemSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListDbSystemsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListDbSystemsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListDbSystemsSortByEnum Enum with underlying type: string
type ListDbSystemsSortByEnum string

// Set of constants representing the allowable values for ListDbSystemsSortByEnum
const (
	ListDbSystemsSortByTimecreated ListDbSystemsSortByEnum = "TIMECREATED"
	ListDbSystemsSortByDisplayname ListDbSystemsSortByEnum = "DISPLAYNAME"
)

var mappingListDbSystemsSortBy = map[string]ListDbSystemsSortByEnum{
	"TIMECREATED": ListDbSystemsSortByTimecreated,
	"DISPLAYNAME": ListDbSystemsSortByDisplayname,
}

// GetListDbSystemsSortByEnumValues Enumerates the set of values for ListDbSystemsSortByEnum
func GetListDbSystemsSortByEnumValues() []ListDbSystemsSortByEnum {
	values := make([]ListDbSystemsSortByEnum, 0)
	for _, v := range mappingListDbSystemsSortBy {
		values = append(values, v)
	}
	return values
}

// ListDbSystemsSortOrderEnum Enum with underlying type: string
type ListDbSystemsSortOrderEnum string

// Set of constants representing the allowable values for ListDbSystemsSortOrderEnum
const (
	ListDbSystemsSortOrderAsc  ListDbSystemsSortOrderEnum = "ASC"
	ListDbSystemsSortOrderDesc ListDbSystemsSortOrderEnum = "DESC"
)

var mappingListDbSystemsSortOrder = map[string]ListDbSystemsSortOrderEnum{
	"ASC":  ListDbSystemsSortOrderAsc,
	"DESC": ListDbSystemsSortOrderDesc,
}

// GetListDbSystemsSortOrderEnumValues Enumerates the set of values for ListDbSystemsSortOrderEnum
func GetListDbSystemsSortOrderEnumValues() []ListDbSystemsSortOrderEnum {
	values := make([]ListDbSystemsSortOrderEnum, 0)
	for _, v := range mappingListDbSystemsSortOrder {
		values = append(values, v)
	}
	return values
}

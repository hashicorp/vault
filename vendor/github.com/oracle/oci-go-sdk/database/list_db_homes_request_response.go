// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListDbHomesRequest wrapper for the ListDbHomes operation
type ListDbHomesRequest struct {

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the DB system.
	DbSystemId *string `mandatory:"false" contributesTo:"query" name:"dbSystemId"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by.  You can provide one sort order (`sortOrder`).  Default order for TIMECREATED is descending.  Default order for DISPLAYNAME is ascending. The DISPLAYNAME sort order is case sensitive.
	SortBy ListDbHomesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListDbHomesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the given lifecycle state exactly.
	LifecycleState DbHomeSummaryLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// A filter to return only resources that match the entire display name given. The match is not case sensitive.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListDbHomesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListDbHomesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListDbHomesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListDbHomesResponse wrapper for the ListDbHomes operation
type ListDbHomesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []DbHomeSummary instances
	Items []DbHomeSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListDbHomesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListDbHomesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListDbHomesSortByEnum Enum with underlying type: string
type ListDbHomesSortByEnum string

// Set of constants representing the allowable values for ListDbHomesSortByEnum
const (
	ListDbHomesSortByTimecreated ListDbHomesSortByEnum = "TIMECREATED"
	ListDbHomesSortByDisplayname ListDbHomesSortByEnum = "DISPLAYNAME"
)

var mappingListDbHomesSortBy = map[string]ListDbHomesSortByEnum{
	"TIMECREATED": ListDbHomesSortByTimecreated,
	"DISPLAYNAME": ListDbHomesSortByDisplayname,
}

// GetListDbHomesSortByEnumValues Enumerates the set of values for ListDbHomesSortByEnum
func GetListDbHomesSortByEnumValues() []ListDbHomesSortByEnum {
	values := make([]ListDbHomesSortByEnum, 0)
	for _, v := range mappingListDbHomesSortBy {
		values = append(values, v)
	}
	return values
}

// ListDbHomesSortOrderEnum Enum with underlying type: string
type ListDbHomesSortOrderEnum string

// Set of constants representing the allowable values for ListDbHomesSortOrderEnum
const (
	ListDbHomesSortOrderAsc  ListDbHomesSortOrderEnum = "ASC"
	ListDbHomesSortOrderDesc ListDbHomesSortOrderEnum = "DESC"
)

var mappingListDbHomesSortOrder = map[string]ListDbHomesSortOrderEnum{
	"ASC":  ListDbHomesSortOrderAsc,
	"DESC": ListDbHomesSortOrderDesc,
}

// GetListDbHomesSortOrderEnumValues Enumerates the set of values for ListDbHomesSortOrderEnum
func GetListDbHomesSortOrderEnumValues() []ListDbHomesSortOrderEnum {
	values := make([]ListDbHomesSortOrderEnum, 0)
	for _, v := range mappingListDbHomesSortOrder {
		values = append(values, v)
	}
	return values
}

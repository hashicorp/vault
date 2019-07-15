// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListDbNodesRequest wrapper for the ListDbNodes operation
type ListDbNodesRequest struct {

	// The compartment OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the DB system.
	DbSystemId *string `mandatory:"false" contributesTo:"query" name:"dbSystemId"`

	// The maximum number of items to return per page.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The pagination token to continue listing from.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Sort by TIMECREATED.  Default order for TIMECREATED is descending.
	SortBy ListDbNodesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`).
	SortOrder ListDbNodesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the given lifecycle state exactly.
	LifecycleState DbNodeSummaryLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListDbNodesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListDbNodesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListDbNodesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListDbNodesResponse wrapper for the ListDbNodes operation
type ListDbNodesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []DbNodeSummary instances
	Items []DbNodeSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then there are additional items still to get. Include this value as the `page` parameter for the
	// subsequent GET request. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListDbNodesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListDbNodesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListDbNodesSortByEnum Enum with underlying type: string
type ListDbNodesSortByEnum string

// Set of constants representing the allowable values for ListDbNodesSortByEnum
const (
	ListDbNodesSortByTimecreated ListDbNodesSortByEnum = "TIMECREATED"
)

var mappingListDbNodesSortBy = map[string]ListDbNodesSortByEnum{
	"TIMECREATED": ListDbNodesSortByTimecreated,
}

// GetListDbNodesSortByEnumValues Enumerates the set of values for ListDbNodesSortByEnum
func GetListDbNodesSortByEnumValues() []ListDbNodesSortByEnum {
	values := make([]ListDbNodesSortByEnum, 0)
	for _, v := range mappingListDbNodesSortBy {
		values = append(values, v)
	}
	return values
}

// ListDbNodesSortOrderEnum Enum with underlying type: string
type ListDbNodesSortOrderEnum string

// Set of constants representing the allowable values for ListDbNodesSortOrderEnum
const (
	ListDbNodesSortOrderAsc  ListDbNodesSortOrderEnum = "ASC"
	ListDbNodesSortOrderDesc ListDbNodesSortOrderEnum = "DESC"
)

var mappingListDbNodesSortOrder = map[string]ListDbNodesSortOrderEnum{
	"ASC":  ListDbNodesSortOrderAsc,
	"DESC": ListDbNodesSortOrderDesc,
}

// GetListDbNodesSortOrderEnumValues Enumerates the set of values for ListDbNodesSortOrderEnum
func GetListDbNodesSortOrderEnumValues() []ListDbNodesSortOrderEnum {
	values := make([]ListDbNodesSortOrderEnum, 0)
	for _, v := range mappingListDbNodesSortOrder {
		values = append(values, v)
	}
	return values
}

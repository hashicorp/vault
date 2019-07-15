// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package resourcemanager

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListStacksRequest wrapper for the ListStacks operation
type ListStacksRequest struct {

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The compartment OCID on which to filter.
	CompartmentId *string `mandatory:"false" contributesTo:"query" name:"compartmentId"`

	// The OCID on which to query for a stack.
	Id *string `mandatory:"false" contributesTo:"query" name:"id"`

	// A filter that returns only those resources that match the specified
	// lifecycle state. The state value is case-insensitive.
	// Allowable values:
	// - CREATING
	// - ACTIVE
	// - DELETING
	// - DELETED
	//
	LifecycleState StackLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// Display name on which to query.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Specifies the field on which to sort.
	// By default, `TIMECREATED` is ordered descending.
	// By default, `DISPLAYNAME` is ordered ascending. Note that you can sort only on one field.
	SortBy ListStacksSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order, either `ASC` (ascending) or `DESC` (descending).
	SortOrder ListStacksSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// The number of items returned in a paginated `List` call. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header from the preceding `List` call.
	// For information about pagination, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListStacksRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListStacksRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListStacksRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListStacksResponse wrapper for the ListStacks operation
type ListStacksResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []StackSummary instances
	Items []StackSummary `presentIn:"body"`

	// Unique identifier for the request.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// Retrieves the next page of paginated list items. If the `opc-next-page`
	// header appears in the response, additional pages of results remain.
	// To receive the next page, include the header value in the `page` param.
	// If the `opc-next-page` header does not appear in the response, there
	// are no more list items to get. For more information about list pagination,
	// see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListStacksResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListStacksResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListStacksSortByEnum Enum with underlying type: string
type ListStacksSortByEnum string

// Set of constants representing the allowable values for ListStacksSortByEnum
const (
	ListStacksSortByTimecreated ListStacksSortByEnum = "TIMECREATED"
	ListStacksSortByDisplayname ListStacksSortByEnum = "DISPLAYNAME"
)

var mappingListStacksSortBy = map[string]ListStacksSortByEnum{
	"TIMECREATED": ListStacksSortByTimecreated,
	"DISPLAYNAME": ListStacksSortByDisplayname,
}

// GetListStacksSortByEnumValues Enumerates the set of values for ListStacksSortByEnum
func GetListStacksSortByEnumValues() []ListStacksSortByEnum {
	values := make([]ListStacksSortByEnum, 0)
	for _, v := range mappingListStacksSortBy {
		values = append(values, v)
	}
	return values
}

// ListStacksSortOrderEnum Enum with underlying type: string
type ListStacksSortOrderEnum string

// Set of constants representing the allowable values for ListStacksSortOrderEnum
const (
	ListStacksSortOrderAsc  ListStacksSortOrderEnum = "ASC"
	ListStacksSortOrderDesc ListStacksSortOrderEnum = "DESC"
)

var mappingListStacksSortOrder = map[string]ListStacksSortOrderEnum{
	"ASC":  ListStacksSortOrderAsc,
	"DESC": ListStacksSortOrderDesc,
}

// GetListStacksSortOrderEnumValues Enumerates the set of values for ListStacksSortOrderEnum
func GetListStacksSortOrderEnumValues() []ListStacksSortOrderEnum {
	values := make([]ListStacksSortOrderEnum, 0)
	for _, v := range mappingListStacksSortOrder {
		values = append(values, v)
	}
	return values
}

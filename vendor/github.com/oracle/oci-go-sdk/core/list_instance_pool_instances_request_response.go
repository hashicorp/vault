// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListInstancePoolInstancesRequest wrapper for the ListInstancePoolInstances operation
type ListInstancePoolInstancesRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The OCID of the instance pool.
	InstancePoolId *string `mandatory:"true" contributesTo:"path" name:"instancePoolId"`

	// A filter to return only resources that match the given display name exactly.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated
	// "List" call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `50`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List"
	// call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by. You can provide one sort order (`sortOrder`). Default order for
	// TIMECREATED is descending. Default order for DISPLAYNAME is ascending. The DISPLAYNAME
	// sort order is case sensitive.
	// **Note:** In general, some "List" operations (for example, `ListInstances`) let you
	// optionally filter by availability domain if the scope of the resource type is within a
	// single availability domain. If you call one of these "List" operations without specifying
	// an availability domain, the resources are grouped by availability domain, then sorted.
	SortBy ListInstancePoolInstancesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListInstancePoolInstancesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListInstancePoolInstancesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListInstancePoolInstancesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListInstancePoolInstancesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListInstancePoolInstancesResponse wrapper for the ListInstancePoolInstances operation
type ListInstancePoolInstancesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []InstanceSummary instances
	Items []InstanceSummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListInstancePoolInstancesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListInstancePoolInstancesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListInstancePoolInstancesSortByEnum Enum with underlying type: string
type ListInstancePoolInstancesSortByEnum string

// Set of constants representing the allowable values for ListInstancePoolInstancesSortByEnum
const (
	ListInstancePoolInstancesSortByTimecreated ListInstancePoolInstancesSortByEnum = "TIMECREATED"
	ListInstancePoolInstancesSortByDisplayname ListInstancePoolInstancesSortByEnum = "DISPLAYNAME"
)

var mappingListInstancePoolInstancesSortBy = map[string]ListInstancePoolInstancesSortByEnum{
	"TIMECREATED": ListInstancePoolInstancesSortByTimecreated,
	"DISPLAYNAME": ListInstancePoolInstancesSortByDisplayname,
}

// GetListInstancePoolInstancesSortByEnumValues Enumerates the set of values for ListInstancePoolInstancesSortByEnum
func GetListInstancePoolInstancesSortByEnumValues() []ListInstancePoolInstancesSortByEnum {
	values := make([]ListInstancePoolInstancesSortByEnum, 0)
	for _, v := range mappingListInstancePoolInstancesSortBy {
		values = append(values, v)
	}
	return values
}

// ListInstancePoolInstancesSortOrderEnum Enum with underlying type: string
type ListInstancePoolInstancesSortOrderEnum string

// Set of constants representing the allowable values for ListInstancePoolInstancesSortOrderEnum
const (
	ListInstancePoolInstancesSortOrderAsc  ListInstancePoolInstancesSortOrderEnum = "ASC"
	ListInstancePoolInstancesSortOrderDesc ListInstancePoolInstancesSortOrderEnum = "DESC"
)

var mappingListInstancePoolInstancesSortOrder = map[string]ListInstancePoolInstancesSortOrderEnum{
	"ASC":  ListInstancePoolInstancesSortOrderAsc,
	"DESC": ListInstancePoolInstancesSortOrderDesc,
}

// GetListInstancePoolInstancesSortOrderEnumValues Enumerates the set of values for ListInstancePoolInstancesSortOrderEnum
func GetListInstancePoolInstancesSortOrderEnumValues() []ListInstancePoolInstancesSortOrderEnum {
	values := make([]ListInstancePoolInstancesSortOrderEnum, 0)
	for _, v := range mappingListInstancePoolInstancesSortOrder {
		values = append(values, v)
	}
	return values
}

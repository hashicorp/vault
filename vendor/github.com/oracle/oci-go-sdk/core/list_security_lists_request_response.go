// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListSecurityListsRequest wrapper for the ListSecurityLists operation
type ListSecurityListsRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the VCN.
	VcnId *string `mandatory:"true" contributesTo:"query" name:"vcnId"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated
	// "List" call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `50`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List"
	// call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// A filter to return only resources that match the given display name exactly.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// The field to sort by. You can provide one sort order (`sortOrder`). Default order for
	// TIMECREATED is descending. Default order for DISPLAYNAME is ascending. The DISPLAYNAME
	// sort order is case sensitive.
	// **Note:** In general, some "List" operations (for example, `ListInstances`) let you
	// optionally filter by availability domain if the scope of the resource type is within a
	// single availability domain. If you call one of these "List" operations without specifying
	// an availability domain, the resources are grouped by availability domain, then sorted.
	SortBy ListSecurityListsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListSecurityListsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to only return resources that match the given lifecycle state.  The state value is case-insensitive.
	LifecycleState SecurityListLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListSecurityListsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListSecurityListsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListSecurityListsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListSecurityListsResponse wrapper for the ListSecurityLists operation
type ListSecurityListsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []SecurityList instances
	Items []SecurityList `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListSecurityListsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListSecurityListsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListSecurityListsSortByEnum Enum with underlying type: string
type ListSecurityListsSortByEnum string

// Set of constants representing the allowable values for ListSecurityListsSortByEnum
const (
	ListSecurityListsSortByTimecreated ListSecurityListsSortByEnum = "TIMECREATED"
	ListSecurityListsSortByDisplayname ListSecurityListsSortByEnum = "DISPLAYNAME"
)

var mappingListSecurityListsSortBy = map[string]ListSecurityListsSortByEnum{
	"TIMECREATED": ListSecurityListsSortByTimecreated,
	"DISPLAYNAME": ListSecurityListsSortByDisplayname,
}

// GetListSecurityListsSortByEnumValues Enumerates the set of values for ListSecurityListsSortByEnum
func GetListSecurityListsSortByEnumValues() []ListSecurityListsSortByEnum {
	values := make([]ListSecurityListsSortByEnum, 0)
	for _, v := range mappingListSecurityListsSortBy {
		values = append(values, v)
	}
	return values
}

// ListSecurityListsSortOrderEnum Enum with underlying type: string
type ListSecurityListsSortOrderEnum string

// Set of constants representing the allowable values for ListSecurityListsSortOrderEnum
const (
	ListSecurityListsSortOrderAsc  ListSecurityListsSortOrderEnum = "ASC"
	ListSecurityListsSortOrderDesc ListSecurityListsSortOrderEnum = "DESC"
)

var mappingListSecurityListsSortOrder = map[string]ListSecurityListsSortOrderEnum{
	"ASC":  ListSecurityListsSortOrderAsc,
	"DESC": ListSecurityListsSortOrderDesc,
}

// GetListSecurityListsSortOrderEnumValues Enumerates the set of values for ListSecurityListsSortOrderEnum
func GetListSecurityListsSortOrderEnumValues() []ListSecurityListsSortOrderEnum {
	values := make([]ListSecurityListsSortOrderEnum, 0)
	for _, v := range mappingListSecurityListsSortOrder {
		values = append(values, v)
	}
	return values
}

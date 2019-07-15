// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package filestorage

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListFileSystemsRequest wrapper for the ListFileSystems operation
type ListFileSystemsRequest struct {

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The name of the availability domain.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"true" contributesTo:"query" name:"availabilityDomain"`

	// For list pagination. The maximum number of results per page,
	// or items to return in a paginated "List" call.
	// 1 is the minimum, 1000 is the maximum.
	// For important details about how pagination works,
	// see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `500`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response
	// header from the previous "List" call.
	// For important details about how pagination works,
	// see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// A user-friendly name. It does not have to be unique, and it is changeable.
	// Example: `My resource`
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Filter results by the specified lifecycle state. Must be a valid
	// state for the resource type.
	LifecycleState ListFileSystemsLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// Filter results by OCID. Must be an OCID of the correct type for
	// the resouce type.
	Id *string `mandatory:"false" contributesTo:"query" name:"id"`

	// The field to sort by. You can provide either value, but not both.
	// By default, when you sort by time created, results are shown
	// in descending order. When you sort by display name, results are
	// shown in ascending order.
	SortBy ListFileSystemsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either 'asc' or 'desc', where 'asc' is
	// ascending and 'desc' is descending.
	SortOrder ListFileSystemsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Unique identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListFileSystemsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListFileSystemsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListFileSystemsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListFileSystemsResponse wrapper for the ListFileSystems operation
type ListFileSystemsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []FileSystemSummary instances
	Items []FileSystemSummary `presentIn:"body"`

	// For list pagination. When this header appears in the response,
	// additional pages of results remain.
	// For important details about how pagination works,
	// see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If
	// you need to contact Oracle about a particular request,
	// please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListFileSystemsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListFileSystemsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListFileSystemsLifecycleStateEnum Enum with underlying type: string
type ListFileSystemsLifecycleStateEnum string

// Set of constants representing the allowable values for ListFileSystemsLifecycleStateEnum
const (
	ListFileSystemsLifecycleStateCreating ListFileSystemsLifecycleStateEnum = "CREATING"
	ListFileSystemsLifecycleStateActive   ListFileSystemsLifecycleStateEnum = "ACTIVE"
	ListFileSystemsLifecycleStateDeleting ListFileSystemsLifecycleStateEnum = "DELETING"
	ListFileSystemsLifecycleStateDeleted  ListFileSystemsLifecycleStateEnum = "DELETED"
	ListFileSystemsLifecycleStateFailed   ListFileSystemsLifecycleStateEnum = "FAILED"
)

var mappingListFileSystemsLifecycleState = map[string]ListFileSystemsLifecycleStateEnum{
	"CREATING": ListFileSystemsLifecycleStateCreating,
	"ACTIVE":   ListFileSystemsLifecycleStateActive,
	"DELETING": ListFileSystemsLifecycleStateDeleting,
	"DELETED":  ListFileSystemsLifecycleStateDeleted,
	"FAILED":   ListFileSystemsLifecycleStateFailed,
}

// GetListFileSystemsLifecycleStateEnumValues Enumerates the set of values for ListFileSystemsLifecycleStateEnum
func GetListFileSystemsLifecycleStateEnumValues() []ListFileSystemsLifecycleStateEnum {
	values := make([]ListFileSystemsLifecycleStateEnum, 0)
	for _, v := range mappingListFileSystemsLifecycleState {
		values = append(values, v)
	}
	return values
}

// ListFileSystemsSortByEnum Enum with underlying type: string
type ListFileSystemsSortByEnum string

// Set of constants representing the allowable values for ListFileSystemsSortByEnum
const (
	ListFileSystemsSortByTimecreated ListFileSystemsSortByEnum = "TIMECREATED"
	ListFileSystemsSortByDisplayname ListFileSystemsSortByEnum = "DISPLAYNAME"
)

var mappingListFileSystemsSortBy = map[string]ListFileSystemsSortByEnum{
	"TIMECREATED": ListFileSystemsSortByTimecreated,
	"DISPLAYNAME": ListFileSystemsSortByDisplayname,
}

// GetListFileSystemsSortByEnumValues Enumerates the set of values for ListFileSystemsSortByEnum
func GetListFileSystemsSortByEnumValues() []ListFileSystemsSortByEnum {
	values := make([]ListFileSystemsSortByEnum, 0)
	for _, v := range mappingListFileSystemsSortBy {
		values = append(values, v)
	}
	return values
}

// ListFileSystemsSortOrderEnum Enum with underlying type: string
type ListFileSystemsSortOrderEnum string

// Set of constants representing the allowable values for ListFileSystemsSortOrderEnum
const (
	ListFileSystemsSortOrderAsc  ListFileSystemsSortOrderEnum = "ASC"
	ListFileSystemsSortOrderDesc ListFileSystemsSortOrderEnum = "DESC"
)

var mappingListFileSystemsSortOrder = map[string]ListFileSystemsSortOrderEnum{
	"ASC":  ListFileSystemsSortOrderAsc,
	"DESC": ListFileSystemsSortOrderDesc,
}

// GetListFileSystemsSortOrderEnumValues Enumerates the set of values for ListFileSystemsSortOrderEnum
func GetListFileSystemsSortOrderEnumValues() []ListFileSystemsSortOrderEnum {
	values := make([]ListFileSystemsSortOrderEnum, 0)
	for _, v := range mappingListFileSystemsSortOrder {
		values = append(values, v)
	}
	return values
}

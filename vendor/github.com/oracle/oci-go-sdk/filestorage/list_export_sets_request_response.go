// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package filestorage

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListExportSetsRequest wrapper for the ListExportSets operation
type ListExportSetsRequest struct {

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
	LifecycleState ListExportSetsLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// Filter results by OCID. Must be an OCID of the correct type for
	// the resouce type.
	Id *string `mandatory:"false" contributesTo:"query" name:"id"`

	// The field to sort by. You can provide either value, but not both.
	// By default, when you sort by time created, results are shown
	// in descending order. When you sort by display name, results are
	// shown in ascending order.
	SortBy ListExportSetsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either 'asc' or 'desc', where 'asc' is
	// ascending and 'desc' is descending.
	SortOrder ListExportSetsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Unique identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListExportSetsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListExportSetsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListExportSetsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListExportSetsResponse wrapper for the ListExportSets operation
type ListExportSetsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []ExportSetSummary instances
	Items []ExportSetSummary `presentIn:"body"`

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

func (response ListExportSetsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListExportSetsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListExportSetsLifecycleStateEnum Enum with underlying type: string
type ListExportSetsLifecycleStateEnum string

// Set of constants representing the allowable values for ListExportSetsLifecycleStateEnum
const (
	ListExportSetsLifecycleStateCreating ListExportSetsLifecycleStateEnum = "CREATING"
	ListExportSetsLifecycleStateActive   ListExportSetsLifecycleStateEnum = "ACTIVE"
	ListExportSetsLifecycleStateDeleting ListExportSetsLifecycleStateEnum = "DELETING"
	ListExportSetsLifecycleStateDeleted  ListExportSetsLifecycleStateEnum = "DELETED"
	ListExportSetsLifecycleStateFailed   ListExportSetsLifecycleStateEnum = "FAILED"
)

var mappingListExportSetsLifecycleState = map[string]ListExportSetsLifecycleStateEnum{
	"CREATING": ListExportSetsLifecycleStateCreating,
	"ACTIVE":   ListExportSetsLifecycleStateActive,
	"DELETING": ListExportSetsLifecycleStateDeleting,
	"DELETED":  ListExportSetsLifecycleStateDeleted,
	"FAILED":   ListExportSetsLifecycleStateFailed,
}

// GetListExportSetsLifecycleStateEnumValues Enumerates the set of values for ListExportSetsLifecycleStateEnum
func GetListExportSetsLifecycleStateEnumValues() []ListExportSetsLifecycleStateEnum {
	values := make([]ListExportSetsLifecycleStateEnum, 0)
	for _, v := range mappingListExportSetsLifecycleState {
		values = append(values, v)
	}
	return values
}

// ListExportSetsSortByEnum Enum with underlying type: string
type ListExportSetsSortByEnum string

// Set of constants representing the allowable values for ListExportSetsSortByEnum
const (
	ListExportSetsSortByTimecreated ListExportSetsSortByEnum = "TIMECREATED"
	ListExportSetsSortByDisplayname ListExportSetsSortByEnum = "DISPLAYNAME"
)

var mappingListExportSetsSortBy = map[string]ListExportSetsSortByEnum{
	"TIMECREATED": ListExportSetsSortByTimecreated,
	"DISPLAYNAME": ListExportSetsSortByDisplayname,
}

// GetListExportSetsSortByEnumValues Enumerates the set of values for ListExportSetsSortByEnum
func GetListExportSetsSortByEnumValues() []ListExportSetsSortByEnum {
	values := make([]ListExportSetsSortByEnum, 0)
	for _, v := range mappingListExportSetsSortBy {
		values = append(values, v)
	}
	return values
}

// ListExportSetsSortOrderEnum Enum with underlying type: string
type ListExportSetsSortOrderEnum string

// Set of constants representing the allowable values for ListExportSetsSortOrderEnum
const (
	ListExportSetsSortOrderAsc  ListExportSetsSortOrderEnum = "ASC"
	ListExportSetsSortOrderDesc ListExportSetsSortOrderEnum = "DESC"
)

var mappingListExportSetsSortOrder = map[string]ListExportSetsSortOrderEnum{
	"ASC":  ListExportSetsSortOrderAsc,
	"DESC": ListExportSetsSortOrderDesc,
}

// GetListExportSetsSortOrderEnumValues Enumerates the set of values for ListExportSetsSortOrderEnum
func GetListExportSetsSortOrderEnumValues() []ListExportSetsSortOrderEnum {
	values := make([]ListExportSetsSortOrderEnum, 0)
	for _, v := range mappingListExportSetsSortOrder {
		values = append(values, v)
	}
	return values
}

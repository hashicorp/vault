// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package filestorage

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListExportsRequest wrapper for the ListExports operation
type ListExportsRequest struct {

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"false" contributesTo:"query" name:"compartmentId"`

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

	// The OCID of the export set.
	ExportSetId *string `mandatory:"false" contributesTo:"query" name:"exportSetId"`

	// The OCID of the file system.
	FileSystemId *string `mandatory:"false" contributesTo:"query" name:"fileSystemId"`

	// Filter results by the specified lifecycle state. Must be a valid
	// state for the resource type.
	LifecycleState ListExportsLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// Filter results by OCID. Must be an OCID of the correct type for
	// the resouce type.
	Id *string `mandatory:"false" contributesTo:"query" name:"id"`

	// The field to sort by. You can provide either value, but not both.
	// By default, when you sort by time created, results are shown
	// in descending order. When you sort by path, results are
	// shown in ascending alphanumeric order.
	SortBy ListExportsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either 'asc' or 'desc', where 'asc' is
	// ascending and 'desc' is descending.
	SortOrder ListExportsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Unique identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListExportsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListExportsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListExportsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListExportsResponse wrapper for the ListExports operation
type ListExportsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []ExportSummary instances
	Items []ExportSummary `presentIn:"body"`

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

func (response ListExportsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListExportsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListExportsLifecycleStateEnum Enum with underlying type: string
type ListExportsLifecycleStateEnum string

// Set of constants representing the allowable values for ListExportsLifecycleStateEnum
const (
	ListExportsLifecycleStateCreating ListExportsLifecycleStateEnum = "CREATING"
	ListExportsLifecycleStateActive   ListExportsLifecycleStateEnum = "ACTIVE"
	ListExportsLifecycleStateDeleting ListExportsLifecycleStateEnum = "DELETING"
	ListExportsLifecycleStateDeleted  ListExportsLifecycleStateEnum = "DELETED"
	ListExportsLifecycleStateFailed   ListExportsLifecycleStateEnum = "FAILED"
)

var mappingListExportsLifecycleState = map[string]ListExportsLifecycleStateEnum{
	"CREATING": ListExportsLifecycleStateCreating,
	"ACTIVE":   ListExportsLifecycleStateActive,
	"DELETING": ListExportsLifecycleStateDeleting,
	"DELETED":  ListExportsLifecycleStateDeleted,
	"FAILED":   ListExportsLifecycleStateFailed,
}

// GetListExportsLifecycleStateEnumValues Enumerates the set of values for ListExportsLifecycleStateEnum
func GetListExportsLifecycleStateEnumValues() []ListExportsLifecycleStateEnum {
	values := make([]ListExportsLifecycleStateEnum, 0)
	for _, v := range mappingListExportsLifecycleState {
		values = append(values, v)
	}
	return values
}

// ListExportsSortByEnum Enum with underlying type: string
type ListExportsSortByEnum string

// Set of constants representing the allowable values for ListExportsSortByEnum
const (
	ListExportsSortByTimecreated ListExportsSortByEnum = "TIMECREATED"
	ListExportsSortByPath        ListExportsSortByEnum = "PATH"
)

var mappingListExportsSortBy = map[string]ListExportsSortByEnum{
	"TIMECREATED": ListExportsSortByTimecreated,
	"PATH":        ListExportsSortByPath,
}

// GetListExportsSortByEnumValues Enumerates the set of values for ListExportsSortByEnum
func GetListExportsSortByEnumValues() []ListExportsSortByEnum {
	values := make([]ListExportsSortByEnum, 0)
	for _, v := range mappingListExportsSortBy {
		values = append(values, v)
	}
	return values
}

// ListExportsSortOrderEnum Enum with underlying type: string
type ListExportsSortOrderEnum string

// Set of constants representing the allowable values for ListExportsSortOrderEnum
const (
	ListExportsSortOrderAsc  ListExportsSortOrderEnum = "ASC"
	ListExportsSortOrderDesc ListExportsSortOrderEnum = "DESC"
)

var mappingListExportsSortOrder = map[string]ListExportsSortOrderEnum{
	"ASC":  ListExportsSortOrderAsc,
	"DESC": ListExportsSortOrderDesc,
}

// GetListExportsSortOrderEnumValues Enumerates the set of values for ListExportsSortOrderEnum
func GetListExportsSortOrderEnumValues() []ListExportsSortOrderEnum {
	values := make([]ListExportsSortOrderEnum, 0)
	for _, v := range mappingListExportsSortOrder {
		values = append(values, v)
	}
	return values
}

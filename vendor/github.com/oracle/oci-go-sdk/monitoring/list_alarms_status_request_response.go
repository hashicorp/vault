// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package monitoring

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAlarmsStatusRequest wrapper for the ListAlarmsStatus operation
type ListAlarmsStatusRequest struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment containing the
	// resources monitored by the metric that you are searching for. Use tenancyId to search in
	// the root compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// Customer part of the request identifier token. If you need to contact Oracle about a particular
	// request, please provide the complete request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// When true, returns resources from all compartments and subcompartments. The parameter can
	// only be set to true when compartmentId is the tenancy OCID (the tenancy is the root compartment).
	// A true value requires the user to have tenancy-level permissions. If this requirement is not met,
	// then the call is rejected. When false, returns resources from only the compartment specified in
	// compartmentId. Default is false.
	CompartmentIdInSubtree *bool `mandatory:"false" contributesTo:"query" name:"compartmentIdInSubtree"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List" call.
	// For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated "List" call.
	// 1 is the minimum, 1000 is the maximum.
	// For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Default: 1000
	// Example: 500
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// A filter to return only resources that match the given display name exactly.
	// Use this filter to list an alarm by name. Alternatively, when you know the alarm OCID, use the GetAlarm operation.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// The field to use when sorting returned alarm definitions. Only one sorting level is provided.
	// Example: `severity`
	SortBy ListAlarmsStatusSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use when sorting returned alarm definitions. Ascending (ASC) or descending (DESC).
	// Example: `ASC`
	SortOrder ListAlarmsStatusSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAlarmsStatusRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAlarmsStatusRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAlarmsStatusRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAlarmsStatusResponse wrapper for the ListAlarmsStatus operation
type ListAlarmsStatusResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AlarmStatusSummary instances
	Items []AlarmStatusSummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages of results remain.
	// For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListAlarmsStatusResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAlarmsStatusResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAlarmsStatusSortByEnum Enum with underlying type: string
type ListAlarmsStatusSortByEnum string

// Set of constants representing the allowable values for ListAlarmsStatusSortByEnum
const (
	ListAlarmsStatusSortByDisplayname ListAlarmsStatusSortByEnum = "displayName"
	ListAlarmsStatusSortBySeverity    ListAlarmsStatusSortByEnum = "severity"
)

var mappingListAlarmsStatusSortBy = map[string]ListAlarmsStatusSortByEnum{
	"displayName": ListAlarmsStatusSortByDisplayname,
	"severity":    ListAlarmsStatusSortBySeverity,
}

// GetListAlarmsStatusSortByEnumValues Enumerates the set of values for ListAlarmsStatusSortByEnum
func GetListAlarmsStatusSortByEnumValues() []ListAlarmsStatusSortByEnum {
	values := make([]ListAlarmsStatusSortByEnum, 0)
	for _, v := range mappingListAlarmsStatusSortBy {
		values = append(values, v)
	}
	return values
}

// ListAlarmsStatusSortOrderEnum Enum with underlying type: string
type ListAlarmsStatusSortOrderEnum string

// Set of constants representing the allowable values for ListAlarmsStatusSortOrderEnum
const (
	ListAlarmsStatusSortOrderAsc  ListAlarmsStatusSortOrderEnum = "ASC"
	ListAlarmsStatusSortOrderDesc ListAlarmsStatusSortOrderEnum = "DESC"
)

var mappingListAlarmsStatusSortOrder = map[string]ListAlarmsStatusSortOrderEnum{
	"ASC":  ListAlarmsStatusSortOrderAsc,
	"DESC": ListAlarmsStatusSortOrderDesc,
}

// GetListAlarmsStatusSortOrderEnumValues Enumerates the set of values for ListAlarmsStatusSortOrderEnum
func GetListAlarmsStatusSortOrderEnumValues() []ListAlarmsStatusSortOrderEnum {
	values := make([]ListAlarmsStatusSortOrderEnum, 0)
	for _, v := range mappingListAlarmsStatusSortOrder {
		values = append(values, v)
	}
	return values
}

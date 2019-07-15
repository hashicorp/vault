// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package autoscaling

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAutoScalingConfigurationsRequest wrapper for the ListAutoScalingConfigurations operation
type ListAutoScalingConfigurationsRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment containing the
	// resources monitored by the metric that you are searching for. Use tenancyId to search in
	// the root compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// A filter to return only resources that match the given display name exactly.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// For list pagination. The maximum number of items to return in a paginated "List" call. For important details
	// about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List" call. For important
	// details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by. You can provide one sort order (`sortOrder`). Default order for
	// TIMECREATED is descending. Default order for DISPLAYNAME is ascending. The DISPLAYNAME
	// sort order is case sensitive.
	SortBy ListAutoScalingConfigurationsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListAutoScalingConfigurationsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAutoScalingConfigurationsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAutoScalingConfigurationsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAutoScalingConfigurationsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAutoScalingConfigurationsResponse wrapper for the ListAutoScalingConfigurations operation
type ListAutoScalingConfigurationsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AutoScalingConfigurationSummary instances
	Items []AutoScalingConfigurationSummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages of results remain.
	// For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListAutoScalingConfigurationsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAutoScalingConfigurationsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAutoScalingConfigurationsSortByEnum Enum with underlying type: string
type ListAutoScalingConfigurationsSortByEnum string

// Set of constants representing the allowable values for ListAutoScalingConfigurationsSortByEnum
const (
	ListAutoScalingConfigurationsSortByTimecreated ListAutoScalingConfigurationsSortByEnum = "TIMECREATED"
	ListAutoScalingConfigurationsSortByDisplayname ListAutoScalingConfigurationsSortByEnum = "DISPLAYNAME"
)

var mappingListAutoScalingConfigurationsSortBy = map[string]ListAutoScalingConfigurationsSortByEnum{
	"TIMECREATED": ListAutoScalingConfigurationsSortByTimecreated,
	"DISPLAYNAME": ListAutoScalingConfigurationsSortByDisplayname,
}

// GetListAutoScalingConfigurationsSortByEnumValues Enumerates the set of values for ListAutoScalingConfigurationsSortByEnum
func GetListAutoScalingConfigurationsSortByEnumValues() []ListAutoScalingConfigurationsSortByEnum {
	values := make([]ListAutoScalingConfigurationsSortByEnum, 0)
	for _, v := range mappingListAutoScalingConfigurationsSortBy {
		values = append(values, v)
	}
	return values
}

// ListAutoScalingConfigurationsSortOrderEnum Enum with underlying type: string
type ListAutoScalingConfigurationsSortOrderEnum string

// Set of constants representing the allowable values for ListAutoScalingConfigurationsSortOrderEnum
const (
	ListAutoScalingConfigurationsSortOrderAsc  ListAutoScalingConfigurationsSortOrderEnum = "ASC"
	ListAutoScalingConfigurationsSortOrderDesc ListAutoScalingConfigurationsSortOrderEnum = "DESC"
)

var mappingListAutoScalingConfigurationsSortOrder = map[string]ListAutoScalingConfigurationsSortOrderEnum{
	"ASC":  ListAutoScalingConfigurationsSortOrderAsc,
	"DESC": ListAutoScalingConfigurationsSortOrderDesc,
}

// GetListAutoScalingConfigurationsSortOrderEnumValues Enumerates the set of values for ListAutoScalingConfigurationsSortOrderEnum
func GetListAutoScalingConfigurationsSortOrderEnumValues() []ListAutoScalingConfigurationsSortOrderEnum {
	values := make([]ListAutoScalingConfigurationsSortOrderEnum, 0)
	for _, v := range mappingListAutoScalingConfigurationsSortOrder {
		values = append(values, v)
	}
	return values
}

// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package autoscaling

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAutoScalingPoliciesRequest wrapper for the ListAutoScalingPolicies operation
type ListAutoScalingPoliciesRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the autoscaling configuration.
	AutoScalingConfigurationId *string `mandatory:"true" contributesTo:"path" name:"autoScalingConfigurationId"`

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
	SortBy ListAutoScalingPoliciesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListAutoScalingPoliciesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAutoScalingPoliciesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAutoScalingPoliciesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAutoScalingPoliciesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAutoScalingPoliciesResponse wrapper for the ListAutoScalingPolicies operation
type ListAutoScalingPoliciesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AutoScalingPolicySummary instances
	Items []AutoScalingPolicySummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages of results remain.
	// For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListAutoScalingPoliciesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAutoScalingPoliciesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAutoScalingPoliciesSortByEnum Enum with underlying type: string
type ListAutoScalingPoliciesSortByEnum string

// Set of constants representing the allowable values for ListAutoScalingPoliciesSortByEnum
const (
	ListAutoScalingPoliciesSortByTimecreated ListAutoScalingPoliciesSortByEnum = "TIMECREATED"
	ListAutoScalingPoliciesSortByDisplayname ListAutoScalingPoliciesSortByEnum = "DISPLAYNAME"
)

var mappingListAutoScalingPoliciesSortBy = map[string]ListAutoScalingPoliciesSortByEnum{
	"TIMECREATED": ListAutoScalingPoliciesSortByTimecreated,
	"DISPLAYNAME": ListAutoScalingPoliciesSortByDisplayname,
}

// GetListAutoScalingPoliciesSortByEnumValues Enumerates the set of values for ListAutoScalingPoliciesSortByEnum
func GetListAutoScalingPoliciesSortByEnumValues() []ListAutoScalingPoliciesSortByEnum {
	values := make([]ListAutoScalingPoliciesSortByEnum, 0)
	for _, v := range mappingListAutoScalingPoliciesSortBy {
		values = append(values, v)
	}
	return values
}

// ListAutoScalingPoliciesSortOrderEnum Enum with underlying type: string
type ListAutoScalingPoliciesSortOrderEnum string

// Set of constants representing the allowable values for ListAutoScalingPoliciesSortOrderEnum
const (
	ListAutoScalingPoliciesSortOrderAsc  ListAutoScalingPoliciesSortOrderEnum = "ASC"
	ListAutoScalingPoliciesSortOrderDesc ListAutoScalingPoliciesSortOrderEnum = "DESC"
)

var mappingListAutoScalingPoliciesSortOrder = map[string]ListAutoScalingPoliciesSortOrderEnum{
	"ASC":  ListAutoScalingPoliciesSortOrderAsc,
	"DESC": ListAutoScalingPoliciesSortOrderDesc,
}

// GetListAutoScalingPoliciesSortOrderEnumValues Enumerates the set of values for ListAutoScalingPoliciesSortOrderEnum
func GetListAutoScalingPoliciesSortOrderEnumValues() []ListAutoScalingPoliciesSortOrderEnum {
	values := make([]ListAutoScalingPoliciesSortOrderEnum, 0)
	for _, v := range mappingListAutoScalingPoliciesSortOrder {
		values = append(values, v)
	}
	return values
}

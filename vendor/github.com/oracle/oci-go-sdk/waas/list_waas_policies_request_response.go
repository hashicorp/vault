// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListWaasPoliciesRequest wrapper for the ListWaasPolicies operation
type ListWaasPoliciesRequest struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment. This number is generated when the compartment is created.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The maximum number of items to return in a paginated call. In unspecified, defaults to `10`.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header from the previous paginated call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The value by which policies are sorted in a paginated 'List' call.  If unspecified, defaults to `timeCreated`.
	SortBy ListWaasPoliciesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The value of the sorting direction of resources in a paginated 'List' call. If unspecified, defaults to `DESC`.
	SortOrder ListWaasPoliciesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Filter policies using a list of policy OCIDs.
	Id []string `contributesTo:"query" name:"id" collectionFormat:"multi"`

	// Filter policies using a list of display names.
	DisplayName []string `contributesTo:"query" name:"displayName" collectionFormat:"multi"`

	// Filter policies using a list of lifecycle states.
	LifecycleState []string `contributesTo:"query" name:"lifecycleState" collectionFormat:"multi"`

	// A filter that matches policies created on or after the specified date and time.
	TimeCreatedGreaterThanOrEqualTo *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timeCreatedGreaterThanOrEqualTo"`

	// A filter that matches policies created before the specified date-time.
	TimeCreatedLessThan *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timeCreatedLessThan"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListWaasPoliciesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListWaasPoliciesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListWaasPoliciesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListWaasPoliciesResponse wrapper for the ListWaasPolicies operation
type ListWaasPoliciesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []WaasPolicySummary instances
	Items []WaasPolicySummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages of results may remain. For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// A unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListWaasPoliciesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListWaasPoliciesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListWaasPoliciesSortByEnum Enum with underlying type: string
type ListWaasPoliciesSortByEnum string

// Set of constants representing the allowable values for ListWaasPoliciesSortByEnum
const (
	ListWaasPoliciesSortById          ListWaasPoliciesSortByEnum = "id"
	ListWaasPoliciesSortByDisplayname ListWaasPoliciesSortByEnum = "displayName"
	ListWaasPoliciesSortByTimecreated ListWaasPoliciesSortByEnum = "timeCreated"
)

var mappingListWaasPoliciesSortBy = map[string]ListWaasPoliciesSortByEnum{
	"id":          ListWaasPoliciesSortById,
	"displayName": ListWaasPoliciesSortByDisplayname,
	"timeCreated": ListWaasPoliciesSortByTimecreated,
}

// GetListWaasPoliciesSortByEnumValues Enumerates the set of values for ListWaasPoliciesSortByEnum
func GetListWaasPoliciesSortByEnumValues() []ListWaasPoliciesSortByEnum {
	values := make([]ListWaasPoliciesSortByEnum, 0)
	for _, v := range mappingListWaasPoliciesSortBy {
		values = append(values, v)
	}
	return values
}

// ListWaasPoliciesSortOrderEnum Enum with underlying type: string
type ListWaasPoliciesSortOrderEnum string

// Set of constants representing the allowable values for ListWaasPoliciesSortOrderEnum
const (
	ListWaasPoliciesSortOrderAsc  ListWaasPoliciesSortOrderEnum = "ASC"
	ListWaasPoliciesSortOrderDesc ListWaasPoliciesSortOrderEnum = "DESC"
)

var mappingListWaasPoliciesSortOrder = map[string]ListWaasPoliciesSortOrderEnum{
	"ASC":  ListWaasPoliciesSortOrderAsc,
	"DESC": ListWaasPoliciesSortOrderDesc,
}

// GetListWaasPoliciesSortOrderEnumValues Enumerates the set of values for ListWaasPoliciesSortOrderEnum
func GetListWaasPoliciesSortOrderEnumValues() []ListWaasPoliciesSortOrderEnum {
	values := make([]ListWaasPoliciesSortOrderEnum, 0)
	for _, v := range mappingListWaasPoliciesSortOrder {
		values = append(values, v)
	}
	return values
}

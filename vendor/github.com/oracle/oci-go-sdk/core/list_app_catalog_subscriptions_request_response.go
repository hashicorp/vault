// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAppCatalogSubscriptionsRequest wrapper for the ListAppCatalogSubscriptions operation
type ListAppCatalogSubscriptionsRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

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
	SortBy ListAppCatalogSubscriptionsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListAppCatalogSubscriptionsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only the listings that matches the given listing id.
	ListingId *string `mandatory:"false" contributesTo:"query" name:"listingId"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAppCatalogSubscriptionsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAppCatalogSubscriptionsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAppCatalogSubscriptionsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAppCatalogSubscriptionsResponse wrapper for the ListAppCatalogSubscriptions operation
type ListAppCatalogSubscriptionsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AppCatalogSubscriptionSummary instances
	Items []AppCatalogSubscriptionSummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListAppCatalogSubscriptionsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAppCatalogSubscriptionsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAppCatalogSubscriptionsSortByEnum Enum with underlying type: string
type ListAppCatalogSubscriptionsSortByEnum string

// Set of constants representing the allowable values for ListAppCatalogSubscriptionsSortByEnum
const (
	ListAppCatalogSubscriptionsSortByTimecreated ListAppCatalogSubscriptionsSortByEnum = "TIMECREATED"
	ListAppCatalogSubscriptionsSortByDisplayname ListAppCatalogSubscriptionsSortByEnum = "DISPLAYNAME"
)

var mappingListAppCatalogSubscriptionsSortBy = map[string]ListAppCatalogSubscriptionsSortByEnum{
	"TIMECREATED": ListAppCatalogSubscriptionsSortByTimecreated,
	"DISPLAYNAME": ListAppCatalogSubscriptionsSortByDisplayname,
}

// GetListAppCatalogSubscriptionsSortByEnumValues Enumerates the set of values for ListAppCatalogSubscriptionsSortByEnum
func GetListAppCatalogSubscriptionsSortByEnumValues() []ListAppCatalogSubscriptionsSortByEnum {
	values := make([]ListAppCatalogSubscriptionsSortByEnum, 0)
	for _, v := range mappingListAppCatalogSubscriptionsSortBy {
		values = append(values, v)
	}
	return values
}

// ListAppCatalogSubscriptionsSortOrderEnum Enum with underlying type: string
type ListAppCatalogSubscriptionsSortOrderEnum string

// Set of constants representing the allowable values for ListAppCatalogSubscriptionsSortOrderEnum
const (
	ListAppCatalogSubscriptionsSortOrderAsc  ListAppCatalogSubscriptionsSortOrderEnum = "ASC"
	ListAppCatalogSubscriptionsSortOrderDesc ListAppCatalogSubscriptionsSortOrderEnum = "DESC"
)

var mappingListAppCatalogSubscriptionsSortOrder = map[string]ListAppCatalogSubscriptionsSortOrderEnum{
	"ASC":  ListAppCatalogSubscriptionsSortOrderAsc,
	"DESC": ListAppCatalogSubscriptionsSortOrderDesc,
}

// GetListAppCatalogSubscriptionsSortOrderEnumValues Enumerates the set of values for ListAppCatalogSubscriptionsSortOrderEnum
func GetListAppCatalogSubscriptionsSortOrderEnumValues() []ListAppCatalogSubscriptionsSortOrderEnum {
	values := make([]ListAppCatalogSubscriptionsSortOrderEnum, 0)
	for _, v := range mappingListAppCatalogSubscriptionsSortOrder {
		values = append(values, v)
	}
	return values
}

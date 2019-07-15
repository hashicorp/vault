// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAppCatalogListingResourceVersionsRequest wrapper for the ListAppCatalogListingResourceVersions operation
type ListAppCatalogListingResourceVersionsRequest struct {

	// The OCID of the listing.
	ListingId *string `mandatory:"true" contributesTo:"path" name:"listingId"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated
	// "List" call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `50`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List"
	// call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListAppCatalogListingResourceVersionsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAppCatalogListingResourceVersionsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAppCatalogListingResourceVersionsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAppCatalogListingResourceVersionsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAppCatalogListingResourceVersionsResponse wrapper for the ListAppCatalogListingResourceVersions operation
type ListAppCatalogListingResourceVersionsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AppCatalogListingResourceVersionSummary instances
	Items []AppCatalogListingResourceVersionSummary `presentIn:"body"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `presentIn:"header" name:"etag"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListAppCatalogListingResourceVersionsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAppCatalogListingResourceVersionsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAppCatalogListingResourceVersionsSortOrderEnum Enum with underlying type: string
type ListAppCatalogListingResourceVersionsSortOrderEnum string

// Set of constants representing the allowable values for ListAppCatalogListingResourceVersionsSortOrderEnum
const (
	ListAppCatalogListingResourceVersionsSortOrderAsc  ListAppCatalogListingResourceVersionsSortOrderEnum = "ASC"
	ListAppCatalogListingResourceVersionsSortOrderDesc ListAppCatalogListingResourceVersionsSortOrderEnum = "DESC"
)

var mappingListAppCatalogListingResourceVersionsSortOrder = map[string]ListAppCatalogListingResourceVersionsSortOrderEnum{
	"ASC":  ListAppCatalogListingResourceVersionsSortOrderAsc,
	"DESC": ListAppCatalogListingResourceVersionsSortOrderDesc,
}

// GetListAppCatalogListingResourceVersionsSortOrderEnumValues Enumerates the set of values for ListAppCatalogListingResourceVersionsSortOrderEnum
func GetListAppCatalogListingResourceVersionsSortOrderEnumValues() []ListAppCatalogListingResourceVersionsSortOrderEnum {
	values := make([]ListAppCatalogListingResourceVersionsSortOrderEnum, 0)
	for _, v := range mappingListAppCatalogListingResourceVersionsSortOrder {
		values = append(values, v)
	}
	return values
}

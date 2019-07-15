// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAppCatalogListingsRequest wrapper for the ListAppCatalogListings operation
type ListAppCatalogListingsRequest struct {

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
	SortOrder ListAppCatalogListingsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only the publisher that matches the given publisher name exactly.
	PublisherName *string `mandatory:"false" contributesTo:"query" name:"publisherName"`

	// A filter to return only publishers that match the given publisher type exactly. Valid types are OCI, ORACLE, TRUSTED, STANDARD.
	PublisherType *string `mandatory:"false" contributesTo:"query" name:"publisherType"`

	// A filter to return only resources that match the given display name exactly.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAppCatalogListingsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAppCatalogListingsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAppCatalogListingsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAppCatalogListingsResponse wrapper for the ListAppCatalogListings operation
type ListAppCatalogListingsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AppCatalogListingSummary instances
	Items []AppCatalogListingSummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListAppCatalogListingsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAppCatalogListingsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAppCatalogListingsSortOrderEnum Enum with underlying type: string
type ListAppCatalogListingsSortOrderEnum string

// Set of constants representing the allowable values for ListAppCatalogListingsSortOrderEnum
const (
	ListAppCatalogListingsSortOrderAsc  ListAppCatalogListingsSortOrderEnum = "ASC"
	ListAppCatalogListingsSortOrderDesc ListAppCatalogListingsSortOrderEnum = "DESC"
)

var mappingListAppCatalogListingsSortOrder = map[string]ListAppCatalogListingsSortOrderEnum{
	"ASC":  ListAppCatalogListingsSortOrderAsc,
	"DESC": ListAppCatalogListingsSortOrderDesc,
}

// GetListAppCatalogListingsSortOrderEnumValues Enumerates the set of values for ListAppCatalogListingsSortOrderEnum
func GetListAppCatalogListingsSortOrderEnumValues() []ListAppCatalogListingsSortOrderEnum {
	values := make([]ListAppCatalogListingsSortOrderEnum, 0)
	for _, v := range mappingListAppCatalogListingsSortOrder {
		values = append(values, v)
	}
	return values
}

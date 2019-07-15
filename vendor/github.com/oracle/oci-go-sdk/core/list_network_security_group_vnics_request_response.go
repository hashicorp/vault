// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListNetworkSecurityGroupVnicsRequest wrapper for the ListNetworkSecurityGroupVnics operation
type ListNetworkSecurityGroupVnicsRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the network security group.
	NetworkSecurityGroupId *string `mandatory:"true" contributesTo:"path" name:"networkSecurityGroupId"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated
	// "List" call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `50`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List"
	// call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by.
	SortBy ListNetworkSecurityGroupVnicsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListNetworkSecurityGroupVnicsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListNetworkSecurityGroupVnicsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListNetworkSecurityGroupVnicsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListNetworkSecurityGroupVnicsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListNetworkSecurityGroupVnicsResponse wrapper for the ListNetworkSecurityGroupVnics operation
type ListNetworkSecurityGroupVnicsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []NetworkSecurityGroupVnic instances
	Items []NetworkSecurityGroupVnic `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListNetworkSecurityGroupVnicsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListNetworkSecurityGroupVnicsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListNetworkSecurityGroupVnicsSortByEnum Enum with underlying type: string
type ListNetworkSecurityGroupVnicsSortByEnum string

// Set of constants representing the allowable values for ListNetworkSecurityGroupVnicsSortByEnum
const (
	ListNetworkSecurityGroupVnicsSortByTimeassociated ListNetworkSecurityGroupVnicsSortByEnum = "TIMEASSOCIATED"
)

var mappingListNetworkSecurityGroupVnicsSortBy = map[string]ListNetworkSecurityGroupVnicsSortByEnum{
	"TIMEASSOCIATED": ListNetworkSecurityGroupVnicsSortByTimeassociated,
}

// GetListNetworkSecurityGroupVnicsSortByEnumValues Enumerates the set of values for ListNetworkSecurityGroupVnicsSortByEnum
func GetListNetworkSecurityGroupVnicsSortByEnumValues() []ListNetworkSecurityGroupVnicsSortByEnum {
	values := make([]ListNetworkSecurityGroupVnicsSortByEnum, 0)
	for _, v := range mappingListNetworkSecurityGroupVnicsSortBy {
		values = append(values, v)
	}
	return values
}

// ListNetworkSecurityGroupVnicsSortOrderEnum Enum with underlying type: string
type ListNetworkSecurityGroupVnicsSortOrderEnum string

// Set of constants representing the allowable values for ListNetworkSecurityGroupVnicsSortOrderEnum
const (
	ListNetworkSecurityGroupVnicsSortOrderAsc  ListNetworkSecurityGroupVnicsSortOrderEnum = "ASC"
	ListNetworkSecurityGroupVnicsSortOrderDesc ListNetworkSecurityGroupVnicsSortOrderEnum = "DESC"
)

var mappingListNetworkSecurityGroupVnicsSortOrder = map[string]ListNetworkSecurityGroupVnicsSortOrderEnum{
	"ASC":  ListNetworkSecurityGroupVnicsSortOrderAsc,
	"DESC": ListNetworkSecurityGroupVnicsSortOrderDesc,
}

// GetListNetworkSecurityGroupVnicsSortOrderEnumValues Enumerates the set of values for ListNetworkSecurityGroupVnicsSortOrderEnum
func GetListNetworkSecurityGroupVnicsSortOrderEnumValues() []ListNetworkSecurityGroupVnicsSortOrderEnum {
	values := make([]ListNetworkSecurityGroupVnicsSortOrderEnum, 0)
	for _, v := range mappingListNetworkSecurityGroupVnicsSortOrder {
		values = append(values, v)
	}
	return values
}

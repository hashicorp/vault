// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListEdgeSubnetsRequest wrapper for the ListEdgeSubnets operation
type ListEdgeSubnetsRequest struct {

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The maximum number of items to return in a paginated call. In unspecified, defaults to `10`.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header from the previous paginated call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The value by which edge node subnets are sorted in a paginated 'List' call. If unspecified, defaults to `timeModified`.
	SortBy ListEdgeSubnetsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The value of the sorting direction of resources in a paginated 'List' call. If unspecified, defaults to `DESC`.
	SortOrder ListEdgeSubnetsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListEdgeSubnetsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListEdgeSubnetsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListEdgeSubnetsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListEdgeSubnetsResponse wrapper for the ListEdgeSubnets operation
type ListEdgeSubnetsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []EdgeSubnet instances
	Items []EdgeSubnet `presentIn:"body"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response, then a partial list might have been returned. Include this value as the `page` parameter for the subsequent `GET` request to get the next batch of items.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// A unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListEdgeSubnetsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListEdgeSubnetsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListEdgeSubnetsSortByEnum Enum with underlying type: string
type ListEdgeSubnetsSortByEnum string

// Set of constants representing the allowable values for ListEdgeSubnetsSortByEnum
const (
	ListEdgeSubnetsSortByCidr         ListEdgeSubnetsSortByEnum = "cidr"
	ListEdgeSubnetsSortByRegion       ListEdgeSubnetsSortByEnum = "region"
	ListEdgeSubnetsSortByTimemodified ListEdgeSubnetsSortByEnum = "timeModified"
)

var mappingListEdgeSubnetsSortBy = map[string]ListEdgeSubnetsSortByEnum{
	"cidr":         ListEdgeSubnetsSortByCidr,
	"region":       ListEdgeSubnetsSortByRegion,
	"timeModified": ListEdgeSubnetsSortByTimemodified,
}

// GetListEdgeSubnetsSortByEnumValues Enumerates the set of values for ListEdgeSubnetsSortByEnum
func GetListEdgeSubnetsSortByEnumValues() []ListEdgeSubnetsSortByEnum {
	values := make([]ListEdgeSubnetsSortByEnum, 0)
	for _, v := range mappingListEdgeSubnetsSortBy {
		values = append(values, v)
	}
	return values
}

// ListEdgeSubnetsSortOrderEnum Enum with underlying type: string
type ListEdgeSubnetsSortOrderEnum string

// Set of constants representing the allowable values for ListEdgeSubnetsSortOrderEnum
const (
	ListEdgeSubnetsSortOrderAsc  ListEdgeSubnetsSortOrderEnum = "ASC"
	ListEdgeSubnetsSortOrderDesc ListEdgeSubnetsSortOrderEnum = "DESC"
)

var mappingListEdgeSubnetsSortOrder = map[string]ListEdgeSubnetsSortOrderEnum{
	"ASC":  ListEdgeSubnetsSortOrderAsc,
	"DESC": ListEdgeSubnetsSortOrderDesc,
}

// GetListEdgeSubnetsSortOrderEnumValues Enumerates the set of values for ListEdgeSubnetsSortOrderEnum
func GetListEdgeSubnetsSortOrderEnumValues() []ListEdgeSubnetsSortOrderEnum {
	values := make([]ListEdgeSubnetsSortOrderEnum, 0)
	for _, v := range mappingListEdgeSubnetsSortOrder {
		values = append(values, v)
	}
	return values
}

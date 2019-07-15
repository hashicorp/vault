// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListNatGatewaysRequest wrapper for the ListNatGateways operation
type ListNatGatewaysRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the VCN.
	VcnId *string `mandatory:"false" contributesTo:"query" name:"vcnId"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated
	// "List" call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `50`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List"
	// call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// A filter to return only resources that match the given display name exactly.
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// The field to sort by. You can provide one sort order (`sortOrder`). Default order for
	// TIMECREATED is descending. Default order for DISPLAYNAME is ascending. The DISPLAYNAME
	// sort order is case sensitive.
	// **Note:** In general, some "List" operations (for example, `ListInstances`) let you
	// optionally filter by availability domain if the scope of the resource type is within a
	// single availability domain. If you call one of these "List" operations without specifying
	// an availability domain, the resources are grouped by availability domain, then sorted.
	SortBy ListNatGatewaysSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListNatGatewaysSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to return only resources that match the specified lifecycle state. The value is case insensitive.
	LifecycleState NatGatewayLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListNatGatewaysRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListNatGatewaysRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListNatGatewaysRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListNatGatewaysResponse wrapper for the ListNatGateways operation
type ListNatGatewaysResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []NatGateway instances
	Items []NatGateway `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListNatGatewaysResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListNatGatewaysResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListNatGatewaysSortByEnum Enum with underlying type: string
type ListNatGatewaysSortByEnum string

// Set of constants representing the allowable values for ListNatGatewaysSortByEnum
const (
	ListNatGatewaysSortByTimecreated ListNatGatewaysSortByEnum = "TIMECREATED"
	ListNatGatewaysSortByDisplayname ListNatGatewaysSortByEnum = "DISPLAYNAME"
)

var mappingListNatGatewaysSortBy = map[string]ListNatGatewaysSortByEnum{
	"TIMECREATED": ListNatGatewaysSortByTimecreated,
	"DISPLAYNAME": ListNatGatewaysSortByDisplayname,
}

// GetListNatGatewaysSortByEnumValues Enumerates the set of values for ListNatGatewaysSortByEnum
func GetListNatGatewaysSortByEnumValues() []ListNatGatewaysSortByEnum {
	values := make([]ListNatGatewaysSortByEnum, 0)
	for _, v := range mappingListNatGatewaysSortBy {
		values = append(values, v)
	}
	return values
}

// ListNatGatewaysSortOrderEnum Enum with underlying type: string
type ListNatGatewaysSortOrderEnum string

// Set of constants representing the allowable values for ListNatGatewaysSortOrderEnum
const (
	ListNatGatewaysSortOrderAsc  ListNatGatewaysSortOrderEnum = "ASC"
	ListNatGatewaysSortOrderDesc ListNatGatewaysSortOrderEnum = "DESC"
)

var mappingListNatGatewaysSortOrder = map[string]ListNatGatewaysSortOrderEnum{
	"ASC":  ListNatGatewaysSortOrderAsc,
	"DESC": ListNatGatewaysSortOrderDesc,
}

// GetListNatGatewaysSortOrderEnumValues Enumerates the set of values for ListNatGatewaysSortOrderEnum
func GetListNatGatewaysSortOrderEnumValues() []ListNatGatewaysSortOrderEnum {
	values := make([]ListNatGatewaysSortOrderEnum, 0)
	for _, v := range mappingListNatGatewaysSortOrder {
		values = append(values, v)
	}
	return values
}

// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package containerengine

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListNodePoolsRequest wrapper for the ListNodePools operation
type ListNodePoolsRequest struct {

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The OCID of the cluster.
	ClusterId *string `mandatory:"false" contributesTo:"query" name:"clusterId"`

	// The name to filter on.
	Name *string `mandatory:"false" contributesTo:"query" name:"name"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated "List" call.
	// 1 is the minimum, 1000 is the maximum. For important details about how pagination works,
	// see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List" call.
	// For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The optional order in which to sort the results.
	SortOrder ListNodePoolsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// The optional field to sort the results by.
	SortBy ListNodePoolsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListNodePoolsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListNodePoolsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListNodePoolsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListNodePoolsResponse wrapper for the ListNodePools operation
type ListNodePoolsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []NodePoolSummary instances
	Items []NodePoolSummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages of results remain.
	// For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListNodePoolsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListNodePoolsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListNodePoolsSortOrderEnum Enum with underlying type: string
type ListNodePoolsSortOrderEnum string

// Set of constants representing the allowable values for ListNodePoolsSortOrderEnum
const (
	ListNodePoolsSortOrderAsc  ListNodePoolsSortOrderEnum = "ASC"
	ListNodePoolsSortOrderDesc ListNodePoolsSortOrderEnum = "DESC"
)

var mappingListNodePoolsSortOrder = map[string]ListNodePoolsSortOrderEnum{
	"ASC":  ListNodePoolsSortOrderAsc,
	"DESC": ListNodePoolsSortOrderDesc,
}

// GetListNodePoolsSortOrderEnumValues Enumerates the set of values for ListNodePoolsSortOrderEnum
func GetListNodePoolsSortOrderEnumValues() []ListNodePoolsSortOrderEnum {
	values := make([]ListNodePoolsSortOrderEnum, 0)
	for _, v := range mappingListNodePoolsSortOrder {
		values = append(values, v)
	}
	return values
}

// ListNodePoolsSortByEnum Enum with underlying type: string
type ListNodePoolsSortByEnum string

// Set of constants representing the allowable values for ListNodePoolsSortByEnum
const (
	ListNodePoolsSortById          ListNodePoolsSortByEnum = "ID"
	ListNodePoolsSortByName        ListNodePoolsSortByEnum = "NAME"
	ListNodePoolsSortByTimeCreated ListNodePoolsSortByEnum = "TIME_CREATED"
)

var mappingListNodePoolsSortBy = map[string]ListNodePoolsSortByEnum{
	"ID":           ListNodePoolsSortById,
	"NAME":         ListNodePoolsSortByName,
	"TIME_CREATED": ListNodePoolsSortByTimeCreated,
}

// GetListNodePoolsSortByEnumValues Enumerates the set of values for ListNodePoolsSortByEnum
func GetListNodePoolsSortByEnumValues() []ListNodePoolsSortByEnum {
	values := make([]ListNodePoolsSortByEnum, 0)
	for _, v := range mappingListNodePoolsSortBy {
		values = append(values, v)
	}
	return values
}

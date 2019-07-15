// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package streaming

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListStreamsRequest wrapper for the ListStreams operation
type ListStreamsRequest struct {

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// A filter to return only resources that match the given ID exactly.
	Id *string `mandatory:"false" contributesTo:"query" name:"id"`

	// A filter to return only resources that match the given name exactly.
	Name *string `mandatory:"false" contributesTo:"query" name:"name"`

	// The maximum number of items to return. The value must be between 1 and 50. The default is 10.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The page at which to start retrieving results.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by. You can provide no more than one sort order. By default, `TIMECREATED` sorts results in descending order and `NAME` sorts results in ascending order.
	SortBy ListStreamsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either 'asc' or 'desc'.
	SortOrder ListStreamsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// A filter to only return resources that match the given lifecycle state. The state value is case-insensitive.
	LifecycleState StreamLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListStreamsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListStreamsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListStreamsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListStreamsResponse wrapper for the ListStreams operation
type ListStreamsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []StreamSummary instances
	Items []StreamSummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages of results remain. For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// For list pagination. When this header appears in the response, previous pages of results exist. For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcPrevPage *string `presentIn:"header" name:"opc-prev-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListStreamsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListStreamsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListStreamsSortByEnum Enum with underlying type: string
type ListStreamsSortByEnum string

// Set of constants representing the allowable values for ListStreamsSortByEnum
const (
	ListStreamsSortByName        ListStreamsSortByEnum = "NAME"
	ListStreamsSortByTimecreated ListStreamsSortByEnum = "TIMECREATED"
)

var mappingListStreamsSortBy = map[string]ListStreamsSortByEnum{
	"NAME":        ListStreamsSortByName,
	"TIMECREATED": ListStreamsSortByTimecreated,
}

// GetListStreamsSortByEnumValues Enumerates the set of values for ListStreamsSortByEnum
func GetListStreamsSortByEnumValues() []ListStreamsSortByEnum {
	values := make([]ListStreamsSortByEnum, 0)
	for _, v := range mappingListStreamsSortBy {
		values = append(values, v)
	}
	return values
}

// ListStreamsSortOrderEnum Enum with underlying type: string
type ListStreamsSortOrderEnum string

// Set of constants representing the allowable values for ListStreamsSortOrderEnum
const (
	ListStreamsSortOrderAsc  ListStreamsSortOrderEnum = "ASC"
	ListStreamsSortOrderDesc ListStreamsSortOrderEnum = "DESC"
)

var mappingListStreamsSortOrder = map[string]ListStreamsSortOrderEnum{
	"ASC":  ListStreamsSortOrderAsc,
	"DESC": ListStreamsSortOrderDesc,
}

// GetListStreamsSortOrderEnumValues Enumerates the set of values for ListStreamsSortOrderEnum
func GetListStreamsSortOrderEnumValues() []ListStreamsSortOrderEnum {
	values := make([]ListStreamsSortOrderEnum, 0)
	for _, v := range mappingListStreamsSortOrder {
		values = append(values, v)
	}
	return values
}

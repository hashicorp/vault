// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package email

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListSendersRequest wrapper for the ListSenders operation
type ListSendersRequest struct {

	// The OCID for the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The request ID for tracing from the system
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The current state of a sender.
	LifecycleState SenderLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// The email address of the approved sender.
	EmailAddress *string `mandatory:"false" contributesTo:"query" name:"emailAddress"`

	// For list pagination. The value of the opc-next-page response header from the previous "List" call.
	// For important details about how pagination works,
	// see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// For list pagination. The maximum number of results per page, or items to return in a
	// paginated "List" call. `1` is the minimum, `1000` is the maximum. For important details about
	// how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The field to sort by. The `TIMECREATED` value returns the list in in
	// descending order by default. The `EMAILADDRESS` value returns the list in
	// ascending order by default. Use the `SortOrderQueryParam` to change the
	// direction of the returned list of items.
	SortBy ListSendersSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending or descending order.
	SortOrder ListSendersSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListSendersRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListSendersRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListSendersRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListSendersResponse wrapper for the ListSenders operation
type ListSendersResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []SenderSummary instances
	Items []SenderSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need
	// to contact Oracle about a particular request, please provide the
	// request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For list pagination. When this header appears in the response, additional
	// pages of results remain. For important details about how pagination works,
	// see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// For list pagination. When this header appears in the response, previous pages
	// of results exist. For important details about how pagination works,
	// see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcPrevPage *string `presentIn:"header" name:"opc-prev-page"`

	// The total number of items returned from the request.
	OpcTotalItems *int `presentIn:"header" name:"opc-total-items"`
}

func (response ListSendersResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListSendersResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListSendersSortByEnum Enum with underlying type: string
type ListSendersSortByEnum string

// Set of constants representing the allowable values for ListSendersSortByEnum
const (
	ListSendersSortByTimecreated  ListSendersSortByEnum = "TIMECREATED"
	ListSendersSortByEmailaddress ListSendersSortByEnum = "EMAILADDRESS"
)

var mappingListSendersSortBy = map[string]ListSendersSortByEnum{
	"TIMECREATED":  ListSendersSortByTimecreated,
	"EMAILADDRESS": ListSendersSortByEmailaddress,
}

// GetListSendersSortByEnumValues Enumerates the set of values for ListSendersSortByEnum
func GetListSendersSortByEnumValues() []ListSendersSortByEnum {
	values := make([]ListSendersSortByEnum, 0)
	for _, v := range mappingListSendersSortBy {
		values = append(values, v)
	}
	return values
}

// ListSendersSortOrderEnum Enum with underlying type: string
type ListSendersSortOrderEnum string

// Set of constants representing the allowable values for ListSendersSortOrderEnum
const (
	ListSendersSortOrderAsc  ListSendersSortOrderEnum = "ASC"
	ListSendersSortOrderDesc ListSendersSortOrderEnum = "DESC"
)

var mappingListSendersSortOrder = map[string]ListSendersSortOrderEnum{
	"ASC":  ListSendersSortOrderAsc,
	"DESC": ListSendersSortOrderDesc,
}

// GetListSendersSortOrderEnumValues Enumerates the set of values for ListSendersSortOrderEnum
func GetListSendersSortOrderEnumValues() []ListSendersSortOrderEnum {
	values := make([]ListSendersSortOrderEnum, 0)
	for _, v := range mappingListSendersSortOrder {
		values = append(values, v)
	}
	return values
}

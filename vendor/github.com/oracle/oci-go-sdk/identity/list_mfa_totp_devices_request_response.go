// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package identity

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListMfaTotpDevicesRequest wrapper for the ListMfaTotpDevices operation
type ListMfaTotpDevicesRequest struct {

	// The OCID of the user.
	UserId *string `mandatory:"true" contributesTo:"path" name:"userId"`

	// The value of the `opc-next-page` response header from the previous "List" call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The maximum number of items to return in a paginated "List" call.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The field to sort by. You can provide one sort order (`sortOrder`). Default order for
	// TIMECREATED is descending. Default order for NAME is ascending. The NAME
	// sort order is case sensitive.
	// **Note:** In general, some "List" operations (for example, `ListInstances`) let you
	// optionally filter by Availability Domain if the scope of the resource type is within a
	// single Availability Domain. If you call one of these "List" operations without specifying
	// an Availability Domain, the resources are grouped by Availability Domain, then sorted.
	SortBy ListMfaTotpDevicesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The NAME sort order
	// is case sensitive.
	SortOrder ListMfaTotpDevicesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListMfaTotpDevicesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListMfaTotpDevicesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListMfaTotpDevicesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListMfaTotpDevicesResponse wrapper for the ListMfaTotpDevices operation
type ListMfaTotpDevicesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []MfaTotpDeviceSummary instances
	Items []MfaTotpDeviceSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then a partial list might have been returned. Include this value as the `page` parameter for the
	// subsequent GET request to get the next batch of items.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListMfaTotpDevicesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListMfaTotpDevicesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListMfaTotpDevicesSortByEnum Enum with underlying type: string
type ListMfaTotpDevicesSortByEnum string

// Set of constants representing the allowable values for ListMfaTotpDevicesSortByEnum
const (
	ListMfaTotpDevicesSortByTimecreated ListMfaTotpDevicesSortByEnum = "TIMECREATED"
	ListMfaTotpDevicesSortByName        ListMfaTotpDevicesSortByEnum = "NAME"
)

var mappingListMfaTotpDevicesSortBy = map[string]ListMfaTotpDevicesSortByEnum{
	"TIMECREATED": ListMfaTotpDevicesSortByTimecreated,
	"NAME":        ListMfaTotpDevicesSortByName,
}

// GetListMfaTotpDevicesSortByEnumValues Enumerates the set of values for ListMfaTotpDevicesSortByEnum
func GetListMfaTotpDevicesSortByEnumValues() []ListMfaTotpDevicesSortByEnum {
	values := make([]ListMfaTotpDevicesSortByEnum, 0)
	for _, v := range mappingListMfaTotpDevicesSortBy {
		values = append(values, v)
	}
	return values
}

// ListMfaTotpDevicesSortOrderEnum Enum with underlying type: string
type ListMfaTotpDevicesSortOrderEnum string

// Set of constants representing the allowable values for ListMfaTotpDevicesSortOrderEnum
const (
	ListMfaTotpDevicesSortOrderAsc  ListMfaTotpDevicesSortOrderEnum = "ASC"
	ListMfaTotpDevicesSortOrderDesc ListMfaTotpDevicesSortOrderEnum = "DESC"
)

var mappingListMfaTotpDevicesSortOrder = map[string]ListMfaTotpDevicesSortOrderEnum{
	"ASC":  ListMfaTotpDevicesSortOrderAsc,
	"DESC": ListMfaTotpDevicesSortOrderDesc,
}

// GetListMfaTotpDevicesSortOrderEnumValues Enumerates the set of values for ListMfaTotpDevicesSortOrderEnum
func GetListMfaTotpDevicesSortOrderEnumValues() []ListMfaTotpDevicesSortOrderEnum {
	values := make([]ListMfaTotpDevicesSortOrderEnum, 0)
	for _, v := range mappingListMfaTotpDevicesSortOrder {
		values = append(values, v)
	}
	return values
}

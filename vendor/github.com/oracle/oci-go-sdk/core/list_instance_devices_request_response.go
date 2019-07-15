// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListInstanceDevicesRequest wrapper for the ListInstanceDevices operation
type ListInstanceDevicesRequest struct {

	// The OCID of the instance.
	InstanceId *string `mandatory:"true" contributesTo:"path" name:"instanceId"`

	// A filter to return only available devices or only used devices.
	IsAvailable *bool `mandatory:"false" contributesTo:"query" name:"isAvailable"`

	// A filter to return only devices that match the given name exactly.
	Name *string `mandatory:"false" contributesTo:"query" name:"name"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated
	// "List" call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `50`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List"
	// call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Unique identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The field to sort by. You can provide one sort order (`sortOrder`). Default order for
	// TIMECREATED is descending. Default order for DISPLAYNAME is ascending. The DISPLAYNAME
	// sort order is case sensitive.
	// **Note:** In general, some "List" operations (for example, `ListInstances`) let you
	// optionally filter by availability domain if the scope of the resource type is within a
	// single availability domain. If you call one of these "List" operations without specifying
	// an availability domain, the resources are grouped by availability domain, then sorted.
	SortBy ListInstanceDevicesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListInstanceDevicesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListInstanceDevicesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListInstanceDevicesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListInstanceDevicesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListInstanceDevicesResponse wrapper for the ListInstanceDevices operation
type ListInstanceDevicesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []Device instances
	Items []Device `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListInstanceDevicesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListInstanceDevicesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListInstanceDevicesSortByEnum Enum with underlying type: string
type ListInstanceDevicesSortByEnum string

// Set of constants representing the allowable values for ListInstanceDevicesSortByEnum
const (
	ListInstanceDevicesSortByTimecreated ListInstanceDevicesSortByEnum = "TIMECREATED"
	ListInstanceDevicesSortByDisplayname ListInstanceDevicesSortByEnum = "DISPLAYNAME"
)

var mappingListInstanceDevicesSortBy = map[string]ListInstanceDevicesSortByEnum{
	"TIMECREATED": ListInstanceDevicesSortByTimecreated,
	"DISPLAYNAME": ListInstanceDevicesSortByDisplayname,
}

// GetListInstanceDevicesSortByEnumValues Enumerates the set of values for ListInstanceDevicesSortByEnum
func GetListInstanceDevicesSortByEnumValues() []ListInstanceDevicesSortByEnum {
	values := make([]ListInstanceDevicesSortByEnum, 0)
	for _, v := range mappingListInstanceDevicesSortBy {
		values = append(values, v)
	}
	return values
}

// ListInstanceDevicesSortOrderEnum Enum with underlying type: string
type ListInstanceDevicesSortOrderEnum string

// Set of constants representing the allowable values for ListInstanceDevicesSortOrderEnum
const (
	ListInstanceDevicesSortOrderAsc  ListInstanceDevicesSortOrderEnum = "ASC"
	ListInstanceDevicesSortOrderDesc ListInstanceDevicesSortOrderEnum = "DESC"
)

var mappingListInstanceDevicesSortOrder = map[string]ListInstanceDevicesSortOrderEnum{
	"ASC":  ListInstanceDevicesSortOrderAsc,
	"DESC": ListInstanceDevicesSortOrderDesc,
}

// GetListInstanceDevicesSortOrderEnumValues Enumerates the set of values for ListInstanceDevicesSortOrderEnum
func GetListInstanceDevicesSortOrderEnumValues() []ListInstanceDevicesSortOrderEnum {
	values := make([]ListInstanceDevicesSortOrderEnum, 0)
	for _, v := range mappingListInstanceDevicesSortOrder {
		values = append(values, v)
	}
	return values
}

// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListInstanceConfigurationsRequest wrapper for the ListInstanceConfigurations operation
type ListInstanceConfigurationsRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated
	// "List" call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `50`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List"
	// call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by. You can provide one sort order (`sortOrder`). Default order for
	// TIMECREATED is descending. Default order for DISPLAYNAME is ascending. The DISPLAYNAME
	// sort order is case sensitive.
	// **Note:** In general, some "List" operations (for example, `ListInstances`) let you
	// optionally filter by availability domain if the scope of the resource type is within a
	// single availability domain. If you call one of these "List" operations without specifying
	// an availability domain, the resources are grouped by availability domain, then sorted.
	SortBy ListInstanceConfigurationsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListInstanceConfigurationsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListInstanceConfigurationsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListInstanceConfigurationsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListInstanceConfigurationsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListInstanceConfigurationsResponse wrapper for the ListInstanceConfigurations operation
type ListInstanceConfigurationsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []InstanceConfigurationSummary instances
	Items []InstanceConfigurationSummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListInstanceConfigurationsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListInstanceConfigurationsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListInstanceConfigurationsSortByEnum Enum with underlying type: string
type ListInstanceConfigurationsSortByEnum string

// Set of constants representing the allowable values for ListInstanceConfigurationsSortByEnum
const (
	ListInstanceConfigurationsSortByTimecreated ListInstanceConfigurationsSortByEnum = "TIMECREATED"
	ListInstanceConfigurationsSortByDisplayname ListInstanceConfigurationsSortByEnum = "DISPLAYNAME"
)

var mappingListInstanceConfigurationsSortBy = map[string]ListInstanceConfigurationsSortByEnum{
	"TIMECREATED": ListInstanceConfigurationsSortByTimecreated,
	"DISPLAYNAME": ListInstanceConfigurationsSortByDisplayname,
}

// GetListInstanceConfigurationsSortByEnumValues Enumerates the set of values for ListInstanceConfigurationsSortByEnum
func GetListInstanceConfigurationsSortByEnumValues() []ListInstanceConfigurationsSortByEnum {
	values := make([]ListInstanceConfigurationsSortByEnum, 0)
	for _, v := range mappingListInstanceConfigurationsSortBy {
		values = append(values, v)
	}
	return values
}

// ListInstanceConfigurationsSortOrderEnum Enum with underlying type: string
type ListInstanceConfigurationsSortOrderEnum string

// Set of constants representing the allowable values for ListInstanceConfigurationsSortOrderEnum
const (
	ListInstanceConfigurationsSortOrderAsc  ListInstanceConfigurationsSortOrderEnum = "ASC"
	ListInstanceConfigurationsSortOrderDesc ListInstanceConfigurationsSortOrderEnum = "DESC"
)

var mappingListInstanceConfigurationsSortOrder = map[string]ListInstanceConfigurationsSortOrderEnum{
	"ASC":  ListInstanceConfigurationsSortOrderAsc,
	"DESC": ListInstanceConfigurationsSortOrderDesc,
}

// GetListInstanceConfigurationsSortOrderEnumValues Enumerates the set of values for ListInstanceConfigurationsSortOrderEnum
func GetListInstanceConfigurationsSortOrderEnumValues() []ListInstanceConfigurationsSortOrderEnum {
	values := make([]ListInstanceConfigurationsSortOrderEnum, 0)
	for _, v := range mappingListInstanceConfigurationsSortOrder {
		values = append(values, v)
	}
	return values
}

// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package healthchecks

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListHttpProbeResultsRequest wrapper for the ListHttpProbeResults operation
type ListHttpProbeResultsRequest struct {

	// The OCID of a monitor or on-demand probe.
	ProbeConfigurationId *string `mandatory:"true" contributesTo:"path" name:"probeConfigurationId"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The maximum number of items to return in a paginated "List" call.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header
	// from the previous "List" call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Returns results with a `startTime` equal to or greater than the specified value.
	StartTimeGreaterThanOrEqualTo *float64 `mandatory:"false" contributesTo:"query" name:"startTimeGreaterThanOrEqualTo"`

	// Returns results with a `startTime` equal to or less than the specified value.
	StartTimeLessThanOrEqualTo *float64 `mandatory:"false" contributesTo:"query" name:"startTimeLessThanOrEqualTo"`

	// Controls the sort order of results.
	SortOrder ListHttpProbeResultsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Filters results that match the `target`.
	Target *string `mandatory:"false" contributesTo:"query" name:"target"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListHttpProbeResultsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListHttpProbeResultsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListHttpProbeResultsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListHttpProbeResultsResponse wrapper for the ListHttpProbeResults operation
type ListHttpProbeResultsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []HttpProbeResultSummary instances
	Items []HttpProbeResultSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to
	// contact Oracle about a particular request, please provide
	// the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this
	// header appears in the response, then a partial list might have been
	// returned. Include this value as the `page` parameter for the
	// subsequent GET request to get the next batch of items.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListHttpProbeResultsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListHttpProbeResultsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListHttpProbeResultsSortOrderEnum Enum with underlying type: string
type ListHttpProbeResultsSortOrderEnum string

// Set of constants representing the allowable values for ListHttpProbeResultsSortOrderEnum
const (
	ListHttpProbeResultsSortOrderAsc  ListHttpProbeResultsSortOrderEnum = "ASC"
	ListHttpProbeResultsSortOrderDesc ListHttpProbeResultsSortOrderEnum = "DESC"
)

var mappingListHttpProbeResultsSortOrder = map[string]ListHttpProbeResultsSortOrderEnum{
	"ASC":  ListHttpProbeResultsSortOrderAsc,
	"DESC": ListHttpProbeResultsSortOrderDesc,
}

// GetListHttpProbeResultsSortOrderEnumValues Enumerates the set of values for ListHttpProbeResultsSortOrderEnum
func GetListHttpProbeResultsSortOrderEnumValues() []ListHttpProbeResultsSortOrderEnum {
	values := make([]ListHttpProbeResultsSortOrderEnum, 0)
	for _, v := range mappingListHttpProbeResultsSortOrder {
		values = append(values, v)
	}
	return values
}

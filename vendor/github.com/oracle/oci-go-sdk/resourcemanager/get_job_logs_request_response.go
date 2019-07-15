// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package resourcemanager

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetJobLogsRequest wrapper for the GetJobLogs operation
type GetJobLogsRequest struct {

	// The job OCID.
	JobId *string `mandatory:"true" contributesTo:"path" name:"jobId"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// A filter that returns only logs of a specified type.
	Type []LogEntryTypeEnum `contributesTo:"query" name:"type" omitEmpty:"true" collectionFormat:"multi"`

	// A filter that returns only log entries that match a given severity level or greater.
	LevelGreaterThanOrEqualTo LogEntryLevelEnum `mandatory:"false" contributesTo:"query" name:"levelGreaterThanOrEqualTo" omitEmpty:"true"`

	// The sort order, either `ASC` (ascending) or `DESC` (descending).
	SortOrder GetJobLogsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// The number of items returned in a paginated `List` call. For information about pagination, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header from the preceding `List` call.
	// For information about pagination, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Time stamp specifying the lower time limit for which logs are returned in a query.
	TimestampGreaterThanOrEqualTo *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timestampGreaterThanOrEqualTo"`

	// Time stamp specifying the upper time limit for which logs are returned in a query.
	TimestampLessThanOrEqualTo *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timestampLessThanOrEqualTo"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetJobLogsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetJobLogsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetJobLogsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetJobLogsResponse wrapper for the GetJobLogs operation
type GetJobLogsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []LogEntry instances
	Items []LogEntry `presentIn:"body"`

	// Unique identifier for the request.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// Retrieves the next page of paginated list items. If the `opc-next-page`
	// header appears in the response, additional pages of results remain.
	// To receive the next page, include the header value in the `page` param.
	// If the `opc-next-page` header does not appear in the response, there
	// are no more list items to get. For more information about list pagination,
	// see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response GetJobLogsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetJobLogsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// GetJobLogsSortOrderEnum Enum with underlying type: string
type GetJobLogsSortOrderEnum string

// Set of constants representing the allowable values for GetJobLogsSortOrderEnum
const (
	GetJobLogsSortOrderAsc  GetJobLogsSortOrderEnum = "ASC"
	GetJobLogsSortOrderDesc GetJobLogsSortOrderEnum = "DESC"
)

var mappingGetJobLogsSortOrder = map[string]GetJobLogsSortOrderEnum{
	"ASC":  GetJobLogsSortOrderAsc,
	"DESC": GetJobLogsSortOrderDesc,
}

// GetGetJobLogsSortOrderEnumValues Enumerates the set of values for GetJobLogsSortOrderEnum
func GetGetJobLogsSortOrderEnumValues() []GetJobLogsSortOrderEnum {
	values := make([]GetJobLogsSortOrderEnum, 0)
	for _, v := range mappingGetJobLogsSortOrder {
		values = append(values, v)
	}
	return values
}

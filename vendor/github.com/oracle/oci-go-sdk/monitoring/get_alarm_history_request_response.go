// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package monitoring

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetAlarmHistoryRequest wrapper for the GetAlarmHistory operation
type GetAlarmHistoryRequest struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of an alarm.
	AlarmId *string `mandatory:"true" contributesTo:"path" name:"alarmId"`

	// Customer part of the request identifier token. If you need to contact Oracle about a particular
	// request, please provide the complete request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The type of history entries to retrieve. State history (STATE_HISTORY) or state transition history (STATE_TRANSITION_HISTORY).
	// If not specified, entries of both types are retrieved.
	// Example: `STATE_HISTORY`
	AlarmHistorytype GetAlarmHistoryAlarmHistorytypeEnum `mandatory:"false" contributesTo:"query" name:"alarmHistorytype" omitEmpty:"true"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List" call.
	// For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated "List" call.
	// 1 is the minimum, 1000 is the maximum.
	// For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Default: 1000
	// Example: 500
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// A filter to return only alarm history entries with timestamps occurring on or after the specified date and time. Format defined by RFC3339.
	// Example: `2019-01-01T01:00:00.789Z`
	TimestampGreaterThanOrEqualTo *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timestampGreaterThanOrEqualTo"`

	// A filter to return only alarm history entries with timestamps occurring before the specified date and time. Format defined by RFC3339.
	// Example: `2019-01-02T01:00:00.789Z`
	TimestampLessThan *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timestampLessThan"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetAlarmHistoryRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetAlarmHistoryRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetAlarmHistoryRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetAlarmHistoryResponse wrapper for the GetAlarmHistory operation
type GetAlarmHistoryResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of AlarmHistoryCollection instances
	AlarmHistoryCollection `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For list pagination. When this header appears in the response, additional pages of results remain.
	// For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response GetAlarmHistoryResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetAlarmHistoryResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// GetAlarmHistoryAlarmHistorytypeEnum Enum with underlying type: string
type GetAlarmHistoryAlarmHistorytypeEnum string

// Set of constants representing the allowable values for GetAlarmHistoryAlarmHistorytypeEnum
const (
	GetAlarmHistoryAlarmHistorytypeHistory           GetAlarmHistoryAlarmHistorytypeEnum = "STATE_HISTORY"
	GetAlarmHistoryAlarmHistorytypeTransitionHistory GetAlarmHistoryAlarmHistorytypeEnum = "STATE_TRANSITION_HISTORY"
)

var mappingGetAlarmHistoryAlarmHistorytype = map[string]GetAlarmHistoryAlarmHistorytypeEnum{
	"STATE_HISTORY":            GetAlarmHistoryAlarmHistorytypeHistory,
	"STATE_TRANSITION_HISTORY": GetAlarmHistoryAlarmHistorytypeTransitionHistory,
}

// GetGetAlarmHistoryAlarmHistorytypeEnumValues Enumerates the set of values for GetAlarmHistoryAlarmHistorytypeEnum
func GetGetAlarmHistoryAlarmHistorytypeEnumValues() []GetAlarmHistoryAlarmHistorytypeEnum {
	values := make([]GetAlarmHistoryAlarmHistorytypeEnum, 0)
	for _, v := range mappingGetAlarmHistoryAlarmHistorytype {
		values = append(values, v)
	}
	return values
}

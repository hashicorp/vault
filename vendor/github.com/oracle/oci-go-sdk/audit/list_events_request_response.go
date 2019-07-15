// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package audit

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListEventsRequest wrapper for the ListEvents operation
type ListEventsRequest struct {

	// The OCID of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// Returns events that were processed at or after this start date and time, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// For example, a start value of `2017-01-15T11:30:00Z` will retrieve a list of all events processed since 30 minutes after the 11th hour of January 15, 2017, in Coordinated Universal Time (UTC).
	// You can specify a value with granularity to the minute. Seconds (and milliseconds, if included) must be set to `0`.
	StartTime *common.SDKTime `mandatory:"true" contributesTo:"query" name:"startTime"`

	// Returns events that were processed before this end date and time, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format. For example, a start value of `2017-01-01T00:00:00Z` and an end value of `2017-01-02T00:00:00Z` will retrieve a list of all events processed on January 1, 2017.
	// Similarly, a start value of `2017-01-01T00:00:00Z` and an end value of `2017-02-01T00:00:00Z` will result in a list of all events processed between January 1, 2017 and January 31, 2017.
	// You can specify a value with granularity to the minute. Seconds (and milliseconds, if included) must be set to `0`.
	EndTime *common.SDKTime `mandatory:"true" contributesTo:"query" name:"endTime"`

	// The value of the `opc-next-page` response header from the previous list query.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListEventsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListEventsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListEventsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListEventsResponse wrapper for the ListEvents operation
type ListEventsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AuditEvent instances
	Items []AuditEvent `presentIn:"body"`

	// For pagination of a list of audit events. When this header appears in the response,
	// it means you received a partial list and there are more results.
	// Include this value as the `page` parameter for the subsequent ListEvents request to get the next batch of events.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListEventsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListEventsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package budget

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetAlertRuleRequest wrapper for the GetAlertRule operation
type GetAlertRuleRequest struct {

	// The unique Budget OCID
	BudgetId *string `mandatory:"true" contributesTo:"path" name:"budgetId"`

	// The unique Alert Rule OCID
	AlertRuleId *string `mandatory:"true" contributesTo:"path" name:"alertRuleId"`

	// The client request ID for tracing.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetAlertRuleRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetAlertRuleRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetAlertRuleRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetAlertRuleResponse wrapper for the GetAlertRule operation
type GetAlertRuleResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The AlertRule instance
	AlertRule `presentIn:"body"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `presentIn:"header" name:"etag"`

	// Unique Oracle-assigned identifier for the request. If
	// you need to contact Oracle about a particular request,
	// please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response GetAlertRuleResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetAlertRuleResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

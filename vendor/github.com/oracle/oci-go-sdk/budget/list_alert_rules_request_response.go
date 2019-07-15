// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package budget

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAlertRulesRequest wrapper for the ListAlertRules operation
type ListAlertRulesRequest struct {

	// The unique Budget OCID
	BudgetId *string `mandatory:"true" contributesTo:"path" name:"budgetId"`

	// The maximum number of items to return.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The page token representing the page at which to start retrieving results. This is usually retrieved from a previous list call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The sort order to use, either 'asc' or 'desc'.
	SortOrder ListAlertRulesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// The field to sort by. If not specified, the default is timeCreated.
	// The default sort order for timeCreated is DESC.
	// The default sort order for displayName is ASC in alphanumeric order.
	SortBy *string `mandatory:"false" contributesTo:"query" name:"sortBy"`

	// The current state of the resource to filter by.
	LifecycleState *string `mandatory:"false" contributesTo:"query" name:"lifecycleState"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Example: `My new resource`
	DisplayName *string `mandatory:"false" contributesTo:"query" name:"displayName"`

	// The client request ID for tracing.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAlertRulesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAlertRulesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAlertRulesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAlertRulesResponse wrapper for the ListAlertRules operation
type ListAlertRulesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []AlertRuleSummary instances
	Items []AlertRuleSummary `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If
	// you need to contact Oracle about a particular request,
	// please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of `AlertRuleSummary`s. If this header appears in the response, then this
	// is a partial list of AlertRuleSummaries. Include this value as the `page` parameter in a subsequent
	// GET request to get the next batch of AlertRuleSummaries.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListAlertRulesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAlertRulesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAlertRulesSortOrderEnum Enum with underlying type: string
type ListAlertRulesSortOrderEnum string

// Set of constants representing the allowable values for ListAlertRulesSortOrderEnum
const (
	ListAlertRulesSortOrderAsc  ListAlertRulesSortOrderEnum = "ASC"
	ListAlertRulesSortOrderDesc ListAlertRulesSortOrderEnum = "DESC"
)

var mappingListAlertRulesSortOrder = map[string]ListAlertRulesSortOrderEnum{
	"ASC":  ListAlertRulesSortOrderAsc,
	"DESC": ListAlertRulesSortOrderDesc,
}

// GetListAlertRulesSortOrderEnumValues Enumerates the set of values for ListAlertRulesSortOrderEnum
func GetListAlertRulesSortOrderEnumValues() []ListAlertRulesSortOrderEnum {
	values := make([]ListAlertRulesSortOrderEnum, 0)
	for _, v := range mappingListAlertRulesSortOrder {
		values = append(values, v)
	}
	return values
}

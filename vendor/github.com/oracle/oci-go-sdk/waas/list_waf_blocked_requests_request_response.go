// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListWafBlockedRequestsRequest wrapper for the ListWafBlockedRequests operation
type ListWafBlockedRequestsRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the WAAS policy.
	WaasPolicyId *string `mandatory:"true" contributesTo:"path" name:"waasPolicyId"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// A filter that limits returned events to those occurring on or after a date and time, specified in RFC 3339 format. If unspecified, defaults to 30 minutes before receipt of the request.
	TimeObservedGreaterThanOrEqualTo *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timeObservedGreaterThanOrEqualTo"`

	// A filter that limits returned events to those occurring before a date and time, specified in RFC 3339 format.
	TimeObservedLessThan *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timeObservedLessThan"`

	// The maximum number of items to return in a paginated call. In unspecified, defaults to `10`.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header from the previous paginated call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Filter stats by the Web Application Firewall feature that triggered the block action. If unspecified, data for all WAF features will be returned.
	WafFeature []ListWafBlockedRequestsWafFeatureEnum `contributesTo:"query" name:"wafFeature" omitEmpty:"true" collectionFormat:"multi"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListWafBlockedRequestsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListWafBlockedRequestsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListWafBlockedRequestsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListWafBlockedRequestsResponse wrapper for the ListWafBlockedRequests operation
type ListWafBlockedRequestsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []WafBlockedRequest instances
	Items []WafBlockedRequest `presentIn:"body"`

	// A unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response, then a partial list might have been returned. Include this value as the page parameter for the subsequent GET request to get the next batch of items.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListWafBlockedRequestsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListWafBlockedRequestsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListWafBlockedRequestsWafFeatureEnum Enum with underlying type: string
type ListWafBlockedRequestsWafFeatureEnum string

// Set of constants representing the allowable values for ListWafBlockedRequestsWafFeatureEnum
const (
	ListWafBlockedRequestsWafFeatureProtectionRules            ListWafBlockedRequestsWafFeatureEnum = "PROTECTION_RULES"
	ListWafBlockedRequestsWafFeatureJsChallenge                ListWafBlockedRequestsWafFeatureEnum = "JS_CHALLENGE"
	ListWafBlockedRequestsWafFeatureAccessRules                ListWafBlockedRequestsWafFeatureEnum = "ACCESS_RULES"
	ListWafBlockedRequestsWafFeatureThreatFeeds                ListWafBlockedRequestsWafFeatureEnum = "THREAT_FEEDS"
	ListWafBlockedRequestsWafFeatureHumanInteractionChallenge  ListWafBlockedRequestsWafFeatureEnum = "HUMAN_INTERACTION_CHALLENGE"
	ListWafBlockedRequestsWafFeatureDeviceFingerprintChallenge ListWafBlockedRequestsWafFeatureEnum = "DEVICE_FINGERPRINT_CHALLENGE"
	ListWafBlockedRequestsWafFeatureCaptcha                    ListWafBlockedRequestsWafFeatureEnum = "CAPTCHA"
	ListWafBlockedRequestsWafFeatureAddressRateLimiting        ListWafBlockedRequestsWafFeatureEnum = "ADDRESS_RATE_LIMITING"
)

var mappingListWafBlockedRequestsWafFeature = map[string]ListWafBlockedRequestsWafFeatureEnum{
	"PROTECTION_RULES":             ListWafBlockedRequestsWafFeatureProtectionRules,
	"JS_CHALLENGE":                 ListWafBlockedRequestsWafFeatureJsChallenge,
	"ACCESS_RULES":                 ListWafBlockedRequestsWafFeatureAccessRules,
	"THREAT_FEEDS":                 ListWafBlockedRequestsWafFeatureThreatFeeds,
	"HUMAN_INTERACTION_CHALLENGE":  ListWafBlockedRequestsWafFeatureHumanInteractionChallenge,
	"DEVICE_FINGERPRINT_CHALLENGE": ListWafBlockedRequestsWafFeatureDeviceFingerprintChallenge,
	"CAPTCHA":                      ListWafBlockedRequestsWafFeatureCaptcha,
	"ADDRESS_RATE_LIMITING":        ListWafBlockedRequestsWafFeatureAddressRateLimiting,
}

// GetListWafBlockedRequestsWafFeatureEnumValues Enumerates the set of values for ListWafBlockedRequestsWafFeatureEnum
func GetListWafBlockedRequestsWafFeatureEnumValues() []ListWafBlockedRequestsWafFeatureEnum {
	values := make([]ListWafBlockedRequestsWafFeatureEnum, 0)
	for _, v := range mappingListWafBlockedRequestsWafFeature {
		values = append(values, v)
	}
	return values
}

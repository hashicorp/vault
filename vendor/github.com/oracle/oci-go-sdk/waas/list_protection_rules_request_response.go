// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListProtectionRulesRequest wrapper for the ListProtectionRules operation
type ListProtectionRulesRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the WAAS policy.
	WaasPolicyId *string `mandatory:"true" contributesTo:"path" name:"waasPolicyId"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The maximum number of items to return in a paginated call. In unspecified, defaults to `10`.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header from the previous paginated call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Filter rules using a list of ModSecurity rule IDs.
	ModSecurityRuleId []string `contributesTo:"query" name:"modSecurityRuleId" collectionFormat:"multi"`

	// Filter rules using a list of actions.
	Action []ListProtectionRulesActionEnum `contributesTo:"query" name:"action" omitEmpty:"true" collectionFormat:"multi"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListProtectionRulesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListProtectionRulesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListProtectionRulesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListProtectionRulesResponse wrapper for the ListProtectionRules operation
type ListProtectionRulesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []ProtectionRule instances
	Items []ProtectionRule `presentIn:"body"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `presentIn:"header" name:"etag"`

	// For list pagination. When this header appears in the response, additional pages of results may remain. For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// A unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListProtectionRulesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListProtectionRulesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListProtectionRulesActionEnum Enum with underlying type: string
type ListProtectionRulesActionEnum string

// Set of constants representing the allowable values for ListProtectionRulesActionEnum
const (
	ListProtectionRulesActionOff    ListProtectionRulesActionEnum = "OFF"
	ListProtectionRulesActionDetect ListProtectionRulesActionEnum = "DETECT"
	ListProtectionRulesActionBlock  ListProtectionRulesActionEnum = "BLOCK"
)

var mappingListProtectionRulesAction = map[string]ListProtectionRulesActionEnum{
	"OFF":    ListProtectionRulesActionOff,
	"DETECT": ListProtectionRulesActionDetect,
	"BLOCK":  ListProtectionRulesActionBlock,
}

// GetListProtectionRulesActionEnumValues Enumerates the set of values for ListProtectionRulesActionEnum
func GetListProtectionRulesActionEnumValues() []ListProtectionRulesActionEnum {
	values := make([]ListProtectionRulesActionEnum, 0)
	for _, v := range mappingListProtectionRulesAction {
		values = append(values, v)
	}
	return values
}

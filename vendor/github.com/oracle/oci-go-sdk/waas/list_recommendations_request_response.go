// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListRecommendationsRequest wrapper for the ListRecommendations operation
type ListRecommendationsRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the WAAS policy.
	WaasPolicyId *string `mandatory:"true" contributesTo:"path" name:"waasPolicyId"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// A filter that matches recommended protection rules based on the selected action. If unspecified, rules with any action type are returned.
	RecommendedAction ListRecommendationsRecommendedActionEnum `mandatory:"false" contributesTo:"query" name:"recommendedAction" omitEmpty:"true"`

	// The maximum number of items to return in a paginated call. In unspecified, defaults to `10`.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header from the previous paginated call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListRecommendationsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListRecommendationsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListRecommendationsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListRecommendationsResponse wrapper for the ListRecommendations operation
type ListRecommendationsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []Recommendation instances
	Items []Recommendation `presentIn:"body"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `presentIn:"header" name:"etag"`

	// For list pagination. When this header appears in the response, additional pages of results may remain. For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// A unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListRecommendationsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListRecommendationsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListRecommendationsRecommendedActionEnum Enum with underlying type: string
type ListRecommendationsRecommendedActionEnum string

// Set of constants representing the allowable values for ListRecommendationsRecommendedActionEnum
const (
	ListRecommendationsRecommendedActionDetect ListRecommendationsRecommendedActionEnum = "DETECT"
	ListRecommendationsRecommendedActionBlock  ListRecommendationsRecommendedActionEnum = "BLOCK"
)

var mappingListRecommendationsRecommendedAction = map[string]ListRecommendationsRecommendedActionEnum{
	"DETECT": ListRecommendationsRecommendedActionDetect,
	"BLOCK":  ListRecommendationsRecommendedActionBlock,
}

// GetListRecommendationsRecommendedActionEnumValues Enumerates the set of values for ListRecommendationsRecommendedActionEnum
func GetListRecommendationsRecommendedActionEnumValues() []ListRecommendationsRecommendedActionEnum {
	values := make([]ListRecommendationsRecommendedActionEnum, 0)
	for _, v := range mappingListRecommendationsRecommendedAction {
		values = append(values, v)
	}
	return values
}

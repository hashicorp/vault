// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListNetworkSecurityGroupSecurityRulesRequest wrapper for the ListNetworkSecurityGroupSecurityRules operation
type ListNetworkSecurityGroupSecurityRulesRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the network security group.
	NetworkSecurityGroupId *string `mandatory:"true" contributesTo:"path" name:"networkSecurityGroupId"`

	// Direction of the security rule. Set to `EGRESS` for rules that allow outbound IP packets,
	// or `INGRESS` for rules that allow inbound IP packets.
	Direction ListNetworkSecurityGroupSecurityRulesDirectionEnum `mandatory:"false" contributesTo:"query" name:"direction" omitEmpty:"true"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated
	// "List" call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `50`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List"
	// call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The field to sort by.
	SortBy ListNetworkSecurityGroupSecurityRulesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use, either ascending (`ASC`) or descending (`DESC`). The DISPLAYNAME sort order
	// is case sensitive.
	SortOrder ListNetworkSecurityGroupSecurityRulesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListNetworkSecurityGroupSecurityRulesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListNetworkSecurityGroupSecurityRulesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListNetworkSecurityGroupSecurityRulesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListNetworkSecurityGroupSecurityRulesResponse wrapper for the ListNetworkSecurityGroupSecurityRules operation
type ListNetworkSecurityGroupSecurityRulesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []SecurityRule instances
	Items []SecurityRule `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListNetworkSecurityGroupSecurityRulesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListNetworkSecurityGroupSecurityRulesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListNetworkSecurityGroupSecurityRulesDirectionEnum Enum with underlying type: string
type ListNetworkSecurityGroupSecurityRulesDirectionEnum string

// Set of constants representing the allowable values for ListNetworkSecurityGroupSecurityRulesDirectionEnum
const (
	ListNetworkSecurityGroupSecurityRulesDirectionEgress  ListNetworkSecurityGroupSecurityRulesDirectionEnum = "EGRESS"
	ListNetworkSecurityGroupSecurityRulesDirectionIngress ListNetworkSecurityGroupSecurityRulesDirectionEnum = "INGRESS"
)

var mappingListNetworkSecurityGroupSecurityRulesDirection = map[string]ListNetworkSecurityGroupSecurityRulesDirectionEnum{
	"EGRESS":  ListNetworkSecurityGroupSecurityRulesDirectionEgress,
	"INGRESS": ListNetworkSecurityGroupSecurityRulesDirectionIngress,
}

// GetListNetworkSecurityGroupSecurityRulesDirectionEnumValues Enumerates the set of values for ListNetworkSecurityGroupSecurityRulesDirectionEnum
func GetListNetworkSecurityGroupSecurityRulesDirectionEnumValues() []ListNetworkSecurityGroupSecurityRulesDirectionEnum {
	values := make([]ListNetworkSecurityGroupSecurityRulesDirectionEnum, 0)
	for _, v := range mappingListNetworkSecurityGroupSecurityRulesDirection {
		values = append(values, v)
	}
	return values
}

// ListNetworkSecurityGroupSecurityRulesSortByEnum Enum with underlying type: string
type ListNetworkSecurityGroupSecurityRulesSortByEnum string

// Set of constants representing the allowable values for ListNetworkSecurityGroupSecurityRulesSortByEnum
const (
	ListNetworkSecurityGroupSecurityRulesSortByTimecreated ListNetworkSecurityGroupSecurityRulesSortByEnum = "TIMECREATED"
)

var mappingListNetworkSecurityGroupSecurityRulesSortBy = map[string]ListNetworkSecurityGroupSecurityRulesSortByEnum{
	"TIMECREATED": ListNetworkSecurityGroupSecurityRulesSortByTimecreated,
}

// GetListNetworkSecurityGroupSecurityRulesSortByEnumValues Enumerates the set of values for ListNetworkSecurityGroupSecurityRulesSortByEnum
func GetListNetworkSecurityGroupSecurityRulesSortByEnumValues() []ListNetworkSecurityGroupSecurityRulesSortByEnum {
	values := make([]ListNetworkSecurityGroupSecurityRulesSortByEnum, 0)
	for _, v := range mappingListNetworkSecurityGroupSecurityRulesSortBy {
		values = append(values, v)
	}
	return values
}

// ListNetworkSecurityGroupSecurityRulesSortOrderEnum Enum with underlying type: string
type ListNetworkSecurityGroupSecurityRulesSortOrderEnum string

// Set of constants representing the allowable values for ListNetworkSecurityGroupSecurityRulesSortOrderEnum
const (
	ListNetworkSecurityGroupSecurityRulesSortOrderAsc  ListNetworkSecurityGroupSecurityRulesSortOrderEnum = "ASC"
	ListNetworkSecurityGroupSecurityRulesSortOrderDesc ListNetworkSecurityGroupSecurityRulesSortOrderEnum = "DESC"
)

var mappingListNetworkSecurityGroupSecurityRulesSortOrder = map[string]ListNetworkSecurityGroupSecurityRulesSortOrderEnum{
	"ASC":  ListNetworkSecurityGroupSecurityRulesSortOrderAsc,
	"DESC": ListNetworkSecurityGroupSecurityRulesSortOrderDesc,
}

// GetListNetworkSecurityGroupSecurityRulesSortOrderEnumValues Enumerates the set of values for ListNetworkSecurityGroupSecurityRulesSortOrderEnum
func GetListNetworkSecurityGroupSecurityRulesSortOrderEnumValues() []ListNetworkSecurityGroupSecurityRulesSortOrderEnum {
	values := make([]ListNetworkSecurityGroupSecurityRulesSortOrderEnum, 0)
	for _, v := range mappingListNetworkSecurityGroupSecurityRulesSortOrder {
		values = append(values, v)
	}
	return values
}

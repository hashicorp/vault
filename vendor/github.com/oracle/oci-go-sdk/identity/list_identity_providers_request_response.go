// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package identity

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListIdentityProvidersRequest wrapper for the ListIdentityProviders operation
type ListIdentityProvidersRequest struct {

	// The protocol used for federation.
	Protocol ListIdentityProvidersProtocolEnum `mandatory:"true" contributesTo:"query" name:"protocol" omitEmpty:"true"`

	// The OCID of the compartment (remember that the tenancy is simply the root compartment).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The value of the `opc-next-page` response header from the previous "List" call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The maximum number of items to return in a paginated "List" call.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListIdentityProvidersRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListIdentityProvidersRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListIdentityProvidersRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListIdentityProvidersResponse wrapper for the ListIdentityProviders operation
type ListIdentityProvidersResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []IdentityProvider instances
	Items []IdentityProvider `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then a partial list might have been returned. Include this value as the `page` parameter for the
	// subsequent GET request to get the next batch of items.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListIdentityProvidersResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListIdentityProvidersResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListIdentityProvidersProtocolEnum Enum with underlying type: string
type ListIdentityProvidersProtocolEnum string

// Set of constants representing the allowable values for ListIdentityProvidersProtocolEnum
const (
	ListIdentityProvidersProtocolSaml2 ListIdentityProvidersProtocolEnum = "SAML2"
)

var mappingListIdentityProvidersProtocol = map[string]ListIdentityProvidersProtocolEnum{
	"SAML2": ListIdentityProvidersProtocolSaml2,
}

// GetListIdentityProvidersProtocolEnumValues Enumerates the set of values for ListIdentityProvidersProtocolEnum
func GetListIdentityProvidersProtocolEnumValues() []ListIdentityProvidersProtocolEnum {
	values := make([]ListIdentityProvidersProtocolEnum, 0)
	for _, v := range mappingListIdentityProvidersProtocol {
		values = append(values, v)
	}
	return values
}

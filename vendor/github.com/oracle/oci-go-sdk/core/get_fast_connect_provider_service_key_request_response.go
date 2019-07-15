// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetFastConnectProviderServiceKeyRequest wrapper for the GetFastConnectProviderServiceKey operation
type GetFastConnectProviderServiceKeyRequest struct {

	// The OCID of the provider service.
	ProviderServiceId *string `mandatory:"true" contributesTo:"path" name:"providerServiceId"`

	// The provider service key that the provider gives you when you set up a virtual circuit connection
	// from the provider to Oracle Cloud Infrastructure. You can set up that connection and get your
	// provider service key at the provider's website or portal. For the portal location, see the `description`
	// attribute of the FastConnectProviderService.
	ProviderServiceKeyName *string `mandatory:"true" contributesTo:"path" name:"providerServiceKeyName"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetFastConnectProviderServiceKeyRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetFastConnectProviderServiceKeyRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetFastConnectProviderServiceKeyRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetFastConnectProviderServiceKeyResponse wrapper for the GetFastConnectProviderServiceKey operation
type GetFastConnectProviderServiceKeyResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The FastConnectProviderServiceKey instance
	FastConnectProviderServiceKey `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response GetFastConnectProviderServiceKeyResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetFastConnectProviderServiceKeyResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

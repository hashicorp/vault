// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// AddNetworkSecurityGroupSecurityRulesRequest wrapper for the AddNetworkSecurityGroupSecurityRules operation
type AddNetworkSecurityGroupSecurityRulesRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the network security group.
	NetworkSecurityGroupId *string `mandatory:"true" contributesTo:"path" name:"networkSecurityGroupId"`

	// Request with one or more security rules to be associated with the network security group.
	AddNetworkSecurityGroupSecurityRulesDetails `contributesTo:"body"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request AddNetworkSecurityGroupSecurityRulesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request AddNetworkSecurityGroupSecurityRulesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request AddNetworkSecurityGroupSecurityRulesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// AddNetworkSecurityGroupSecurityRulesResponse wrapper for the AddNetworkSecurityGroupSecurityRules operation
type AddNetworkSecurityGroupSecurityRulesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The AddedNetworkSecurityGroupSecurityRules instance
	AddedNetworkSecurityGroupSecurityRules `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response AddNetworkSecurityGroupSecurityRulesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response AddNetworkSecurityGroupSecurityRulesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

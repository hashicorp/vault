// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package healthchecks

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// CreateOnDemandPingProbeRequest wrapper for the CreateOnDemandPingProbe operation
type CreateOnDemandPingProbeRequest struct {

	// Configuration details for creating an on-demand ping probe.
	CreateOnDemandPingProbeDetails `contributesTo:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request CreateOnDemandPingProbeRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request CreateOnDemandPingProbeRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request CreateOnDemandPingProbeRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// CreateOnDemandPingProbeResponse wrapper for the CreateOnDemandPingProbe operation
type CreateOnDemandPingProbeResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The PingProbe instance
	PingProbe `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to
	// contact Oracle about a particular request, please provide
	// the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// The URL for fetching probe results.
	Location *string `presentIn:"header" name:"location"`
}

func (response CreateOnDemandPingProbeResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response CreateOnDemandPingProbeResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package autoscaling

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetAutoScalingConfigurationRequest wrapper for the GetAutoScalingConfiguration operation
type GetAutoScalingConfigurationRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the autoscaling configuration.
	AutoScalingConfigurationId *string `mandatory:"true" contributesTo:"path" name:"autoScalingConfigurationId"`

	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetAutoScalingConfigurationRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetAutoScalingConfigurationRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetAutoScalingConfigurationRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetAutoScalingConfigurationResponse wrapper for the GetAutoScalingConfiguration operation
type GetAutoScalingConfigurationResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The AutoScalingConfiguration instance
	AutoScalingConfiguration `presentIn:"body"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `presentIn:"header" name:"etag"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response GetAutoScalingConfigurationResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetAutoScalingConfigurationResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

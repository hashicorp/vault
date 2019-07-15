// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package autoscaling

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ChangeAutoScalingConfigurationCompartmentRequest wrapper for the ChangeAutoScalingConfigurationCompartment operation
type ChangeAutoScalingConfigurationCompartmentRequest struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the autoscaling configuration.
	AutoScalingConfigurationId *string `mandatory:"true" contributesTo:"path" name:"autoScalingConfigurationId"`

	// Request to change the compartment of given autoscaling configuration.
	ChangeCompartmentDetails ChangeAutoScalingCompartmentDetails `contributesTo:"body"`

	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// For optimistic concurrency control. In the PUT or DELETE call for a resource, set the `if-match`
	// parameter to the value of the etag from a previous GET or POST response for that resource.  The resource
	// will be updated or deleted only if the etag you provide matches the resource's current etag value.
	IfMatch *string `mandatory:"false" contributesTo:"header" name:"if-match"`

	// A token that uniquely identifies a request so it can be retried in case of a timeout or
	// server error without risk of executing that same action again. Retry tokens expire after 24
	// hours, but can be invalidated before then due to conflicting operations (for example, if a resource
	// has been deleted and purged from the system, then a retry of the original creation request
	// may be rejected).
	OpcRetryToken *string `mandatory:"false" contributesTo:"header" name:"opc-retry-token"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ChangeAutoScalingConfigurationCompartmentRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ChangeAutoScalingConfigurationCompartmentRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ChangeAutoScalingConfigurationCompartmentRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ChangeAutoScalingConfigurationCompartmentResponse wrapper for the ChangeAutoScalingConfigurationCompartment operation
type ChangeAutoScalingConfigurationCompartmentResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The resulting etag of the autoscaling configuration affected by this operation.
	// For optimistic concurrency control. See `if-match`.
	Etag *string `presentIn:"header" name:"etag"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ChangeAutoScalingConfigurationCompartmentResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ChangeAutoScalingConfigurationCompartmentResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// TerminateAutonomousExadataInfrastructureRequest wrapper for the TerminateAutonomousExadataInfrastructure operation
type TerminateAutonomousExadataInfrastructureRequest struct {

	// The Autonomous Exadata Infrastructure  OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	AutonomousExadataInfrastructureId *string `mandatory:"true" contributesTo:"path" name:"autonomousExadataInfrastructureId"`

	// For optimistic concurrency control. In the PUT or DELETE call for a resource, set the `if-match`
	// parameter to the value of the etag from a previous GET or POST response for that resource.  The resource
	// will be updated or deleted only if the etag you provide matches the resource's current etag value.
	IfMatch *string `mandatory:"false" contributesTo:"header" name:"if-match"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request TerminateAutonomousExadataInfrastructureRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request TerminateAutonomousExadataInfrastructureRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request TerminateAutonomousExadataInfrastructureRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// TerminateAutonomousExadataInfrastructureResponse wrapper for the TerminateAutonomousExadataInfrastructure operation
type TerminateAutonomousExadataInfrastructureResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response TerminateAutonomousExadataInfrastructureResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response TerminateAutonomousExadataInfrastructureResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

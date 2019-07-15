// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package database

import (
	"github.com/oracle/oci-go-sdk/common"
	"io"
	"net/http"
)

// GenerateAutonomousDataWarehouseWalletRequest wrapper for the GenerateAutonomousDataWarehouseWallet operation
type GenerateAutonomousDataWarehouseWalletRequest struct {

	// The database OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm).
	AutonomousDataWarehouseId *string `mandatory:"true" contributesTo:"path" name:"autonomousDataWarehouseId"`

	// Request to create a new Autonomous Data Warehouse wallet.
	GenerateAutonomousDataWarehouseWalletDetails `contributesTo:"body"`

	// Unique identifier for the request.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

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

func (request GenerateAutonomousDataWarehouseWalletRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GenerateAutonomousDataWarehouseWalletRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GenerateAutonomousDataWarehouseWalletRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GenerateAutonomousDataWarehouseWalletResponse wrapper for the GenerateAutonomousDataWarehouseWallet operation
type GenerateAutonomousDataWarehouseWalletResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The io.ReadCloser instance
	Content io.ReadCloser `presentIn:"body" encoding:"binary"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `presentIn:"header" name:"etag"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// Size of the file.
	ContentLength *int64 `presentIn:"header" name:"content-length"`

	// The date and time the wallet was created, as described in RFC 3339 (https://tools.ietf.org/rfc/rfc3339), section 14.29.
	LastModified *common.SDKTime `presentIn:"header" name:"last-modified"`
}

func (response GenerateAutonomousDataWarehouseWalletResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GenerateAutonomousDataWarehouseWalletResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

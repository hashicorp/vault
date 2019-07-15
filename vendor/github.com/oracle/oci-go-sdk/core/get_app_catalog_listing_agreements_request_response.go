// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetAppCatalogListingAgreementsRequest wrapper for the GetAppCatalogListingAgreements operation
type GetAppCatalogListingAgreementsRequest struct {

	// The OCID of the listing.
	ListingId *string `mandatory:"true" contributesTo:"path" name:"listingId"`

	// Listing Resource Version.
	ResourceVersion *string `mandatory:"true" contributesTo:"path" name:"resourceVersion"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetAppCatalogListingAgreementsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetAppCatalogListingAgreementsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetAppCatalogListingAgreementsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetAppCatalogListingAgreementsResponse wrapper for the GetAppCatalogListingAgreements operation
type GetAppCatalogListingAgreementsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The AppCatalogListingResourceVersionAgreements instance
	AppCatalogListingResourceVersionAgreements `presentIn:"body"`

	// For optimistic concurrency control. See `if-match`.
	Etag *string `presentIn:"header" name:"etag"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response GetAppCatalogListingAgreementsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetAppCatalogListingAgreementsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

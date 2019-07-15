// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package announcementsservice

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetAnnouncementRequest wrapper for the GetAnnouncement operation
type GetAnnouncementRequest struct {

	// The OCID of the announcement.
	AnnouncementId *string `mandatory:"true" contributesTo:"path" name:"announcementId"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the complete request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetAnnouncementRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetAnnouncementRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetAnnouncementRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetAnnouncementResponse wrapper for the GetAnnouncement operation
type GetAnnouncementResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The Announcement instance
	Announcement `presentIn:"body"`

	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response GetAnnouncementResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetAnnouncementResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

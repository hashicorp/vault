// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package announcementsservice

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// UpdateAnnouncementUserStatusRequest wrapper for the UpdateAnnouncementUserStatus operation
type UpdateAnnouncementUserStatusRequest struct {

	// The OCID of the announcement.
	AnnouncementId *string `mandatory:"true" contributesTo:"path" name:"announcementId"`

	// The information to use to update the announcement's read status.
	StatusDetails AnnouncementUserStatusDetails `contributesTo:"body"`

	// The locking version, used for optimistic concurrency control.
	IfMatch *string `mandatory:"false" contributesTo:"header" name:"if-match"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the complete request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request UpdateAnnouncementUserStatusRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request UpdateAnnouncementUserStatusRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request UpdateAnnouncementUserStatusRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// UpdateAnnouncementUserStatusResponse wrapper for the UpdateAnnouncementUserStatus operation
type UpdateAnnouncementUserStatusResponse struct {

	// The underlying http response
	RawResponse *http.Response

	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	Etag *string `presentIn:"header" name:"etag"`
}

func (response UpdateAnnouncementUserStatusResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response UpdateAnnouncementUserStatusResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

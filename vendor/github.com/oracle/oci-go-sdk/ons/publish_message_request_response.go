// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package ons

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// PublishMessageRequest wrapper for the PublishMessage operation
type PublishMessageRequest struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the topic.
	TopicId *string `mandatory:"true" contributesTo:"path" name:"topicId"`

	// The message to publish.
	MessageDetails `contributesTo:"body"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Type of message body in the request. Default value: JSON.
	MessageType PublishMessageMessageTypeEnum `mandatory:"false" contributesTo:"header" name:"messageType"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request PublishMessageRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request PublishMessageRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request PublishMessageRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// PublishMessageResponse wrapper for the PublishMessage operation
type PublishMessageResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The PublishResult instance
	PublishResult `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response PublishMessageResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response PublishMessageResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// PublishMessageMessageTypeEnum Enum with underlying type: string
type PublishMessageMessageTypeEnum string

// Set of constants representing the allowable values for PublishMessageMessageTypeEnum
const (
	PublishMessageMessageTypeJson    PublishMessageMessageTypeEnum = "JSON"
	PublishMessageMessageTypeRawText PublishMessageMessageTypeEnum = "RAW_TEXT"
)

var mappingPublishMessageMessageType = map[string]PublishMessageMessageTypeEnum{
	"JSON":     PublishMessageMessageTypeJson,
	"RAW_TEXT": PublishMessageMessageTypeRawText,
}

// GetPublishMessageMessageTypeEnumValues Enumerates the set of values for PublishMessageMessageTypeEnum
func GetPublishMessageMessageTypeEnumValues() []PublishMessageMessageTypeEnum {
	values := make([]PublishMessageMessageTypeEnum, 0)
	for _, v := range mappingPublishMessageMessageType {
		values = append(values, v)
	}
	return values
}

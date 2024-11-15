// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// GetRetentionRuleRequest wrapper for the GetRetentionRule operation
type GetRetentionRuleRequest struct {

	// The Object Storage namespace used for the request.
	NamespaceName *string `mandatory:"true" contributesTo:"path" name:"namespaceName"`

	// The name of the bucket. Avoid entering confidential information.
	// Example: `my-new-bucket1`
	BucketName *string `mandatory:"true" contributesTo:"path" name:"bucketName"`

	// The ID of the retention rule.
	RetentionRuleId *string `mandatory:"true" contributesTo:"path" name:"retentionRuleId"`

	// The client request ID for tracing.
	OpcClientRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-client-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request GetRetentionRuleRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request GetRetentionRuleRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request GetRetentionRuleRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// GetRetentionRuleResponse wrapper for the GetRetentionRule operation
type GetRetentionRuleResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The RetentionRule instance
	RetentionRule `presentIn:"body"`

	// Echoes back the value passed in the opc-client-request-id header, for use by clients when debugging.
	OpcClientRequestId *string `presentIn:"header" name:"opc-client-request-id"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular
	// request, provide this request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// The entity tag (ETag) for the retention rule.
	Etag *string `presentIn:"header" name:"etag"`

	// The time the retention rule was last modified, as described in RFC 2616 (https://tools.ietf.org/html/rfc2616#section-14.29)
	LastModified *common.SDKTime `presentIn:"header" name:"last-modified"`
}

func (response GetRetentionRuleResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response GetRetentionRuleResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

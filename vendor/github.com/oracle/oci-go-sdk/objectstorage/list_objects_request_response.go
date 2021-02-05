// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListObjectsRequest wrapper for the ListObjects operation
type ListObjectsRequest struct {

	// The Object Storage namespace used for the request.
	NamespaceName *string `mandatory:"true" contributesTo:"path" name:"namespaceName"`

	// The name of the bucket. Avoid entering confidential information.
	// Example: `my-new-bucket1`
	BucketName *string `mandatory:"true" contributesTo:"path" name:"bucketName"`

	// The string to use for matching against the start of object names in a list query.
	Prefix *string `mandatory:"false" contributesTo:"query" name:"prefix"`

	// Object names returned by a list query must be greater or equal to this parameter.
	Start *string `mandatory:"false" contributesTo:"query" name:"start"`

	// Object names returned by a list query must be strictly less than this parameter.
	End *string `mandatory:"false" contributesTo:"query" name:"end"`

	// The maximum number of items to return.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// When this parameter is set, only objects whose names do not contain the delimiter character
	// (after an optionally specified prefix) are returned in the objects key of the response body.
	// Scanned objects whose names contain the delimiter have the part of their name up to the first
	// occurrence of the delimiter (including the optional prefix) returned as a set of prefixes.
	// Note that only '/' is a supported delimiter character at this time.
	Delimiter *string `mandatory:"false" contributesTo:"query" name:"delimiter"`

	// Object summary in list of objects includes the 'name' field. This parameter can also include 'size'
	// (object size in bytes), 'etag', 'md5', 'timeCreated' (object creation date and time) and 'timeModified'
	// (object modification date and time).
	// Value of this parameter should be a comma-separated, case-insensitive list of those field names.
	// For example 'name,etag,timeCreated,md5,timeModified'
	Fields *string `mandatory:"false" contributesTo:"query" name:"fields" omitEmpty:"true"`

	// The client request ID for tracing.
	OpcClientRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-client-request-id"`

	// Object names returned by a list query must be greater than this parameter.
	StartAfter *string `mandatory:"false" contributesTo:"query" name:"startAfter"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListObjectsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListObjectsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListObjectsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListObjectsResponse wrapper for the ListObjects operation
type ListObjectsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// The ListObjects instance
	ListObjects `presentIn:"body"`

	// Echoes back the value passed in the opc-client-request-id header, for use by clients when debugging.
	OpcClientRequestId *string `presentIn:"header" name:"opc-client-request-id"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular
	// request, provide this request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListObjectsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListObjectsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

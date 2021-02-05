// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// CopyObjectRequest wrapper for the CopyObject operation
type CopyObjectRequest struct {

	// The Object Storage namespace used for the request.
	NamespaceName *string `mandatory:"true" contributesTo:"path" name:"namespaceName"`

	// The name of the bucket. Avoid entering confidential information.
	// Example: `my-new-bucket1`
	BucketName *string `mandatory:"true" contributesTo:"path" name:"bucketName"`

	// The source and destination of the object to be copied.
	CopyObjectDetails `contributesTo:"body"`

	// The client request ID for tracing.
	OpcClientRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-client-request-id"`

	// The optional header that specifies "AES256" as the encryption algorithm. For more information, see
	// Using Your Own Keys for Server-Side Encryption (https://docs.cloud.oracle.com/Content/Object/Tasks/usingyourencryptionkeys.htm).
	OpcSseCustomerAlgorithm *string `mandatory:"false" contributesTo:"header" name:"opc-sse-customer-algorithm"`

	// The optional header that specifies the base64-encoded 256-bit encryption key to use to encrypt or
	// decrypt the data. For more information, see
	// Using Your Own Keys for Server-Side Encryption (https://docs.cloud.oracle.com/Content/Object/Tasks/usingyourencryptionkeys.htm).
	OpcSseCustomerKey *string `mandatory:"false" contributesTo:"header" name:"opc-sse-customer-key"`

	// The optional header that specifies the base64-encoded SHA256 hash of the encryption key. This
	// value is used to check the integrity of the encryption key. For more information, see
	// Using Your Own Keys for Server-Side Encryption (https://docs.cloud.oracle.com/Content/Object/Tasks/usingyourencryptionkeys.htm).
	OpcSseCustomerKeySha256 *string `mandatory:"false" contributesTo:"header" name:"opc-sse-customer-key-sha256"`

	// The optional header that specifies "AES256" as the encryption algorithm to use to decrypt the source
	// object. For more information, see
	// Using Your Own Keys for Server-Side Encryption (https://docs.cloud.oracle.com/Content/Object/Tasks/usingyourencryptionkeys.htm).
	OpcSourceSseCustomerAlgorithm *string `mandatory:"false" contributesTo:"header" name:"opc-source-sse-customer-algorithm"`

	// The optional header that specifies the base64-encoded 256-bit encryption key to use to decrypt
	// the source object. For more information, see
	// Using Your Own Keys for Server-Side Encryption (https://docs.cloud.oracle.com/Content/Object/Tasks/usingyourencryptionkeys.htm).
	OpcSourceSseCustomerKey *string `mandatory:"false" contributesTo:"header" name:"opc-source-sse-customer-key"`

	// The optional header that specifies the base64-encoded SHA256 hash of the encryption key used to
	// decrypt the source object. This value is used to check the integrity of the encryption key. For
	// more information, see
	// Using Your Own Keys for Server-Side Encryption (https://docs.cloud.oracle.com/Content/Object/Tasks/usingyourencryptionkeys.htm).
	OpcSourceSseCustomerKeySha256 *string `mandatory:"false" contributesTo:"header" name:"opc-source-sse-customer-key-sha256"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request CopyObjectRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request CopyObjectRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request CopyObjectRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// CopyObjectResponse wrapper for the CopyObject operation
type CopyObjectResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// Unique Oracle-assigned identifier for the asynchronous request. If you need to contact Oracle about a
	// particular request, provide this request ID.
	OpcWorkRequestId *string `presentIn:"header" name:"opc-work-request-id"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular
	// request, provide this request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// Echoes back the value passed in the opc-client-request-id header, for use by clients when debugging.
	OpcClientRequestId *string `presentIn:"header" name:"opc-client-request-id"`
}

func (response CopyObjectResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response CopyObjectResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

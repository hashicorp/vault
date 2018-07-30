// Copyright (c) 2016, 2018, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package objectstorage

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListBucketsRequest wrapper for the ListBuckets operation
type ListBucketsRequest struct {

	// The top-level namespace used for the request.
	NamespaceName *string `mandatory:"true" contributesTo:"path" name:"namespaceName"`

	// The ID of the compartment in which to list buckets.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The maximum number of items to return.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The page at which to start retrieving results.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// Bucket summary in list of buckets includes the 'namespace', 'name', 'compartmentId', 'createdBy', 'timeCreated',
	// and 'etag' fields. This parameter can also include 'tags' (freeformTags and definedTags). The only supported value
	// of this parameter is 'tags' for now. Example 'tags'.
	Fields []ListBucketsFieldsEnum `contributesTo:"query" name:"fields" omitEmpty:"true" collectionFormat:"csv"`

	// The client request ID for tracing.
	OpcClientRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-client-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListBucketsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListBucketsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListBucketsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListBucketsResponse wrapper for the ListBuckets operation
type ListBucketsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []BucketSummary instances
	Items []BucketSummary `presentIn:"body"`

	// Echoes back the value passed in the opc-client-request-id header, for use by clients when debugging.
	OpcClientRequestId *string `presentIn:"header" name:"opc-client-request-id"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular
	// request, provide this request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of `Bucket`s. If this header appears in the response, then this
	// is a partial list of buckets. Include this value as the `page` parameter in a subsequent
	// GET request to get the next batch of buckets. For information about pagination, see
	// List Pagination (https://docs.us-phoenix-1.oraclecloud.com/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListBucketsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListBucketsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListBucketsFieldsEnum Enum with underlying type: string
type ListBucketsFieldsEnum string

// Set of constants representing the allowable values for ListBucketsFields
const (
	ListBucketsFieldsTags ListBucketsFieldsEnum = "tags"
)

var mappingListBucketsFields = map[string]ListBucketsFieldsEnum{
	"tags": ListBucketsFieldsTags,
}

// GetListBucketsFieldsEnumValues Enumerates the set of values for ListBucketsFields
func GetListBucketsFieldsEnumValues() []ListBucketsFieldsEnum {
	values := make([]ListBucketsFieldsEnum, 0)
	for _, v := range mappingListBucketsFields {
		values = append(values, v)
	}
	return values
}

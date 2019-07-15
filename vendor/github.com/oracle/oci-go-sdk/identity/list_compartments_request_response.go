// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package identity

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListCompartmentsRequest wrapper for the ListCompartments operation
type ListCompartmentsRequest struct {

	// The OCID of the compartment (remember that the tenancy is simply the root compartment).
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The value of the `opc-next-page` response header from the previous "List" call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The maximum number of items to return in a paginated "List" call.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// Valid values are `ANY` and `ACCESSIBLE`. Default is `ANY`.
	// Setting this to `ACCESSIBLE` returns only those compartments for which the
	// user has INSPECT permissions directly or indirectly (permissions can be on a
	// resource in a subcompartment). For the compartments on which the user indirectly has
	// INSPECT permissions, a restricted set of fields is returned.
	// When set to `ANY` permissions are not checked.
	AccessLevel ListCompartmentsAccessLevelEnum `mandatory:"false" contributesTo:"query" name:"accessLevel" omitEmpty:"true"`

	// Default is false. Can only be set to true when performing
	// ListCompartments on the tenancy (root compartment).
	// When set to true, the hierarchy of compartments is traversed
	// and all compartments and subcompartments in the tenancy are
	// returned depending on the the setting of `accessLevel`.
	CompartmentIdInSubtree *bool `mandatory:"false" contributesTo:"query" name:"compartmentIdInSubtree"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListCompartmentsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListCompartmentsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListCompartmentsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListCompartmentsResponse wrapper for the ListCompartments operation
type ListCompartmentsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []Compartment instances
	Items []Compartment `presentIn:"body"`

	// Unique Oracle-assigned identifier for the request. If you need to contact Oracle about a
	// particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`

	// For pagination of a list of items. When paging through a list, if this header appears in the response,
	// then a partial list might have been returned. Include this value as the `page` parameter for the
	// subsequent GET request to get the next batch of items.
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`
}

func (response ListCompartmentsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListCompartmentsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListCompartmentsAccessLevelEnum Enum with underlying type: string
type ListCompartmentsAccessLevelEnum string

// Set of constants representing the allowable values for ListCompartmentsAccessLevelEnum
const (
	ListCompartmentsAccessLevelAny        ListCompartmentsAccessLevelEnum = "ANY"
	ListCompartmentsAccessLevelAccessible ListCompartmentsAccessLevelEnum = "ACCESSIBLE"
)

var mappingListCompartmentsAccessLevel = map[string]ListCompartmentsAccessLevelEnum{
	"ANY":        ListCompartmentsAccessLevelAny,
	"ACCESSIBLE": ListCompartmentsAccessLevelAccessible,
}

// GetListCompartmentsAccessLevelEnumValues Enumerates the set of values for ListCompartmentsAccessLevelEnum
func GetListCompartmentsAccessLevelEnumValues() []ListCompartmentsAccessLevelEnum {
	values := make([]ListCompartmentsAccessLevelEnum, 0)
	for _, v := range mappingListCompartmentsAccessLevel {
		values = append(values, v)
	}
	return values
}

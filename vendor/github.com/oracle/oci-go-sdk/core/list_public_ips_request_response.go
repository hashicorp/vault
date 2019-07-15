// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package core

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListPublicIpsRequest wrapper for the ListPublicIps operation
type ListPublicIpsRequest struct {

	// Whether the public IP is regional or specific to a particular availability domain.
	// * `REGION`: The public IP exists within a region and is assigned to a regional entity
	// (such as a NatGateway), or can be assigned to a private IP
	// in any availability domain in the region. Reserved public IPs have `scope` = `REGION`, as do
	// ephemeral public IPs assigned to a regional entity.
	// * `AVAILABILITY_DOMAIN`: The public IP exists within the availability domain of the entity
	// it's assigned to, which is specified by the `availabilityDomain` property of the public IP object.
	// Ephemeral public IPs that are assigned to private IPs have `scope` = `AVAILABILITY_DOMAIN`.
	Scope ListPublicIpsScopeEnum `mandatory:"true" contributesTo:"query" name:"scope" omitEmpty:"true"`

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// For list pagination. The maximum number of results per page, or items to return in a paginated
	// "List" call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	// Example: `50`
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// For list pagination. The value of the `opc-next-page` response header from the previous "List"
	// call. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The name of the availability domain.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"false" contributesTo:"query" name:"availabilityDomain"`

	// A filter to return only public IPs that match given lifetime.
	Lifetime ListPublicIpsLifetimeEnum `mandatory:"false" contributesTo:"query" name:"lifetime" omitEmpty:"true"`

	// Unique Oracle-assigned identifier for the request.
	// If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListPublicIpsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListPublicIpsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListPublicIpsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListPublicIpsResponse wrapper for the ListPublicIps operation
type ListPublicIpsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []PublicIp instances
	Items []PublicIp `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages
	// of results remain. For important details about how pagination works, see
	// List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// Unique Oracle-assigned identifier for the request. If you need to contact
	// Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListPublicIpsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListPublicIpsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListPublicIpsScopeEnum Enum with underlying type: string
type ListPublicIpsScopeEnum string

// Set of constants representing the allowable values for ListPublicIpsScopeEnum
const (
	ListPublicIpsScopeRegion             ListPublicIpsScopeEnum = "REGION"
	ListPublicIpsScopeAvailabilityDomain ListPublicIpsScopeEnum = "AVAILABILITY_DOMAIN"
)

var mappingListPublicIpsScope = map[string]ListPublicIpsScopeEnum{
	"REGION":              ListPublicIpsScopeRegion,
	"AVAILABILITY_DOMAIN": ListPublicIpsScopeAvailabilityDomain,
}

// GetListPublicIpsScopeEnumValues Enumerates the set of values for ListPublicIpsScopeEnum
func GetListPublicIpsScopeEnumValues() []ListPublicIpsScopeEnum {
	values := make([]ListPublicIpsScopeEnum, 0)
	for _, v := range mappingListPublicIpsScope {
		values = append(values, v)
	}
	return values
}

// ListPublicIpsLifetimeEnum Enum with underlying type: string
type ListPublicIpsLifetimeEnum string

// Set of constants representing the allowable values for ListPublicIpsLifetimeEnum
const (
	ListPublicIpsLifetimeEphemeral ListPublicIpsLifetimeEnum = "EPHEMERAL"
	ListPublicIpsLifetimeReserved  ListPublicIpsLifetimeEnum = "RESERVED"
)

var mappingListPublicIpsLifetime = map[string]ListPublicIpsLifetimeEnum{
	"EPHEMERAL": ListPublicIpsLifetimeEphemeral,
	"RESERVED":  ListPublicIpsLifetimeReserved,
}

// GetListPublicIpsLifetimeEnumValues Enumerates the set of values for ListPublicIpsLifetimeEnum
func GetListPublicIpsLifetimeEnumValues() []ListPublicIpsLifetimeEnum {
	values := make([]ListPublicIpsLifetimeEnum, 0)
	for _, v := range mappingListPublicIpsLifetime {
		values = append(values, v)
	}
	return values
}

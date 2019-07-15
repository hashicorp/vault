// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package waas

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListCertificatesRequest wrapper for the ListCertificates operation
type ListCertificatesRequest struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment. This number is generated when the compartment is created.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// The maximum number of items to return in a paginated call. In unspecified, defaults to `10`.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header from the previous paginated call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The value by which certificate summaries are sorted in a paginated 'List' call. If unspecified, defaults to `timeCreated`.
	SortBy ListCertificatesSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The value of the sorting direction of resources in a paginated 'List' call. If unspecified, defaults to `DESC`.
	SortOrder ListCertificatesSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// Filter certificates using a list of certificates OCIDs.
	Id []string `contributesTo:"query" name:"id" collectionFormat:"multi"`

	// Filter certificates using a list of display names.
	DisplayName []string `contributesTo:"query" name:"displayName" collectionFormat:"multi"`

	// Filter certificates using a list of lifecycle states.
	LifecycleState []string `contributesTo:"query" name:"lifecycleState" collectionFormat:"multi"`

	// A filter that matches certificates created on or after the specified date-time.
	TimeCreatedGreaterThanOrEqualTo *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timeCreatedGreaterThanOrEqualTo"`

	// A filter that matches certificates created before the specified date-time.
	TimeCreatedLessThan *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timeCreatedLessThan"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListCertificatesRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListCertificatesRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListCertificatesRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListCertificatesResponse wrapper for the ListCertificates operation
type ListCertificatesResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of []CertificateSummary instances
	Items []CertificateSummary `presentIn:"body"`

	// For list pagination. When this header appears in the response, additional pages of results may remain. For important details about how pagination works, see List Pagination (https://docs.cloud.oracle.com/iaas/Content/API/Concepts/usingapi.htm#nine).
	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	// A unique Oracle-assigned identifier for the request. If you need to contact Oracle about a particular request, please provide the request ID.
	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListCertificatesResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListCertificatesResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListCertificatesSortByEnum Enum with underlying type: string
type ListCertificatesSortByEnum string

// Set of constants representing the allowable values for ListCertificatesSortByEnum
const (
	ListCertificatesSortById            ListCertificatesSortByEnum = "id"
	ListCertificatesSortByCompartmentid ListCertificatesSortByEnum = "compartmentId"
	ListCertificatesSortByDisplayname   ListCertificatesSortByEnum = "displayName"
	ListCertificatesSortByNotvalidafter ListCertificatesSortByEnum = "notValidAfter"
	ListCertificatesSortByTimecreated   ListCertificatesSortByEnum = "timeCreated"
)

var mappingListCertificatesSortBy = map[string]ListCertificatesSortByEnum{
	"id":            ListCertificatesSortById,
	"compartmentId": ListCertificatesSortByCompartmentid,
	"displayName":   ListCertificatesSortByDisplayname,
	"notValidAfter": ListCertificatesSortByNotvalidafter,
	"timeCreated":   ListCertificatesSortByTimecreated,
}

// GetListCertificatesSortByEnumValues Enumerates the set of values for ListCertificatesSortByEnum
func GetListCertificatesSortByEnumValues() []ListCertificatesSortByEnum {
	values := make([]ListCertificatesSortByEnum, 0)
	for _, v := range mappingListCertificatesSortBy {
		values = append(values, v)
	}
	return values
}

// ListCertificatesSortOrderEnum Enum with underlying type: string
type ListCertificatesSortOrderEnum string

// Set of constants representing the allowable values for ListCertificatesSortOrderEnum
const (
	ListCertificatesSortOrderAsc  ListCertificatesSortOrderEnum = "ASC"
	ListCertificatesSortOrderDesc ListCertificatesSortOrderEnum = "DESC"
)

var mappingListCertificatesSortOrder = map[string]ListCertificatesSortOrderEnum{
	"ASC":  ListCertificatesSortOrderAsc,
	"DESC": ListCertificatesSortOrderDesc,
}

// GetListCertificatesSortOrderEnumValues Enumerates the set of values for ListCertificatesSortOrderEnum
func GetListCertificatesSortOrderEnumValues() []ListCertificatesSortOrderEnum {
	values := make([]ListCertificatesSortOrderEnum, 0)
	for _, v := range mappingListCertificatesSortOrder {
		values = append(values, v)
	}
	return values
}

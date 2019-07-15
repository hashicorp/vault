// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

package announcementsservice

import (
	"github.com/oracle/oci-go-sdk/common"
	"net/http"
)

// ListAnnouncementsRequest wrapper for the ListAnnouncements operation
type ListAnnouncementsRequest struct {

	// The OCID of the compartment. Because announcements are specific to a tenancy, this is the
	// OCID of the root compartment.
	CompartmentId *string `mandatory:"true" contributesTo:"query" name:"compartmentId"`

	// The maximum number of items to return in a paginated "List" call.
	Limit *int `mandatory:"false" contributesTo:"query" name:"limit"`

	// The value of the `opc-next-page` response header from the previous "List" call.
	Page *string `mandatory:"false" contributesTo:"query" name:"page"`

	// The type of announcement.
	AnnouncementType *string `mandatory:"false" contributesTo:"query" name:"announcementType"`

	// The announcement's current lifecycle state.
	LifecycleState ListAnnouncementsLifecycleStateEnum `mandatory:"false" contributesTo:"query" name:"lifecycleState" omitEmpty:"true"`

	// Whether the announcement is displayed as a console banner.
	IsBanner *bool `mandatory:"false" contributesTo:"query" name:"isBanner"`

	// The criteria to sort by. You can specify only one sort order.
	SortBy ListAnnouncementsSortByEnum `mandatory:"false" contributesTo:"query" name:"sortBy" omitEmpty:"true"`

	// The sort order to use. (Sorting by `announcementType` orders the announcements list according to importance.)
	SortOrder ListAnnouncementsSortOrderEnum `mandatory:"false" contributesTo:"query" name:"sortOrder" omitEmpty:"true"`

	// The boundary for the earliest `timeOneValue` date on announcements that you want to see.
	TimeOneEarliestTime *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timeOneEarliestTime"`

	// The boundary for the latest `timeOneValue` date on announcements that you want to see.
	TimeOneLatestTime *common.SDKTime `mandatory:"false" contributesTo:"query" name:"timeOneLatestTime"`

	// The unique Oracle-assigned identifier for the request. If you need to contact Oracle about
	// a particular request, please provide the complete request ID.
	OpcRequestId *string `mandatory:"false" contributesTo:"header" name:"opc-request-id"`

	// Metadata about the request. This information will not be transmitted to the service, but
	// represents information that the SDK will consume to drive retry behavior.
	RequestMetadata common.RequestMetadata
}

func (request ListAnnouncementsRequest) String() string {
	return common.PointerString(request)
}

// HTTPRequest implements the OCIRequest interface
func (request ListAnnouncementsRequest) HTTPRequest(method, path string) (http.Request, error) {
	return common.MakeDefaultHTTPRequestWithTaggedStruct(method, path, request)
}

// RetryPolicy implements the OCIRetryableRequest interface. This retrieves the specified retry policy.
func (request ListAnnouncementsRequest) RetryPolicy() *common.RetryPolicy {
	return request.RequestMetadata.RetryPolicy
}

// ListAnnouncementsResponse wrapper for the ListAnnouncements operation
type ListAnnouncementsResponse struct {

	// The underlying http response
	RawResponse *http.Response

	// A list of AnnouncementsCollection instances
	AnnouncementsCollection `presentIn:"body"`

	OpcNextPage *string `presentIn:"header" name:"opc-next-page"`

	OpcRequestId *string `presentIn:"header" name:"opc-request-id"`
}

func (response ListAnnouncementsResponse) String() string {
	return common.PointerString(response)
}

// HTTPResponse implements the OCIResponse interface
func (response ListAnnouncementsResponse) HTTPResponse() *http.Response {
	return response.RawResponse
}

// ListAnnouncementsLifecycleStateEnum Enum with underlying type: string
type ListAnnouncementsLifecycleStateEnum string

// Set of constants representing the allowable values for ListAnnouncementsLifecycleStateEnum
const (
	ListAnnouncementsLifecycleStateActive   ListAnnouncementsLifecycleStateEnum = "ACTIVE"
	ListAnnouncementsLifecycleStateInactive ListAnnouncementsLifecycleStateEnum = "INACTIVE"
)

var mappingListAnnouncementsLifecycleState = map[string]ListAnnouncementsLifecycleStateEnum{
	"ACTIVE":   ListAnnouncementsLifecycleStateActive,
	"INACTIVE": ListAnnouncementsLifecycleStateInactive,
}

// GetListAnnouncementsLifecycleStateEnumValues Enumerates the set of values for ListAnnouncementsLifecycleStateEnum
func GetListAnnouncementsLifecycleStateEnumValues() []ListAnnouncementsLifecycleStateEnum {
	values := make([]ListAnnouncementsLifecycleStateEnum, 0)
	for _, v := range mappingListAnnouncementsLifecycleState {
		values = append(values, v)
	}
	return values
}

// ListAnnouncementsSortByEnum Enum with underlying type: string
type ListAnnouncementsSortByEnum string

// Set of constants representing the allowable values for ListAnnouncementsSortByEnum
const (
	ListAnnouncementsSortByTimeonevalue          ListAnnouncementsSortByEnum = "timeOneValue"
	ListAnnouncementsSortByTimetwovalue          ListAnnouncementsSortByEnum = "timeTwoValue"
	ListAnnouncementsSortByTimecreated           ListAnnouncementsSortByEnum = "timeCreated"
	ListAnnouncementsSortByReferenceticketnumber ListAnnouncementsSortByEnum = "referenceTicketNumber"
	ListAnnouncementsSortBySummary               ListAnnouncementsSortByEnum = "summary"
	ListAnnouncementsSortByAnnouncementtype      ListAnnouncementsSortByEnum = "announcementType"
)

var mappingListAnnouncementsSortBy = map[string]ListAnnouncementsSortByEnum{
	"timeOneValue":          ListAnnouncementsSortByTimeonevalue,
	"timeTwoValue":          ListAnnouncementsSortByTimetwovalue,
	"timeCreated":           ListAnnouncementsSortByTimecreated,
	"referenceTicketNumber": ListAnnouncementsSortByReferenceticketnumber,
	"summary":               ListAnnouncementsSortBySummary,
	"announcementType":      ListAnnouncementsSortByAnnouncementtype,
}

// GetListAnnouncementsSortByEnumValues Enumerates the set of values for ListAnnouncementsSortByEnum
func GetListAnnouncementsSortByEnumValues() []ListAnnouncementsSortByEnum {
	values := make([]ListAnnouncementsSortByEnum, 0)
	for _, v := range mappingListAnnouncementsSortBy {
		values = append(values, v)
	}
	return values
}

// ListAnnouncementsSortOrderEnum Enum with underlying type: string
type ListAnnouncementsSortOrderEnum string

// Set of constants representing the allowable values for ListAnnouncementsSortOrderEnum
const (
	ListAnnouncementsSortOrderAsc  ListAnnouncementsSortOrderEnum = "ASC"
	ListAnnouncementsSortOrderDesc ListAnnouncementsSortOrderEnum = "DESC"
)

var mappingListAnnouncementsSortOrder = map[string]ListAnnouncementsSortOrderEnum{
	"ASC":  ListAnnouncementsSortOrderAsc,
	"DESC": ListAnnouncementsSortOrderDesc,
}

// GetListAnnouncementsSortOrderEnumValues Enumerates the set of values for ListAnnouncementsSortOrderEnum
func GetListAnnouncementsSortOrderEnumValues() []ListAnnouncementsSortOrderEnum {
	values := make([]ListAnnouncementsSortOrderEnum, 0)
	for _, v := range mappingListAnnouncementsSortOrder {
		values = append(values, v)
	}
	return values
}

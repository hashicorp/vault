package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemCalendarGroupsItemCalendarsItemGetSchedulePostResponseable instead.
type ItemCalendarGroupsItemCalendarsItemGetScheduleResponse struct {
    ItemCalendarGroupsItemCalendarsItemGetSchedulePostResponse
}
// NewItemCalendarGroupsItemCalendarsItemGetScheduleResponse instantiates a new ItemCalendarGroupsItemCalendarsItemGetScheduleResponse and sets the default values.
func NewItemCalendarGroupsItemCalendarsItemGetScheduleResponse()(*ItemCalendarGroupsItemCalendarsItemGetScheduleResponse) {
    m := &ItemCalendarGroupsItemCalendarsItemGetScheduleResponse{
        ItemCalendarGroupsItemCalendarsItemGetSchedulePostResponse: *NewItemCalendarGroupsItemCalendarsItemGetSchedulePostResponse(),
    }
    return m
}
// CreateItemCalendarGroupsItemCalendarsItemGetScheduleResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemCalendarGroupsItemCalendarsItemGetScheduleResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemCalendarGroupsItemCalendarsItemGetScheduleResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemCalendarGroupsItemCalendarsItemGetSchedulePostResponseable instead.
type ItemCalendarGroupsItemCalendarsItemGetScheduleResponseable interface {
    ItemCalendarGroupsItemCalendarsItemGetSchedulePostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

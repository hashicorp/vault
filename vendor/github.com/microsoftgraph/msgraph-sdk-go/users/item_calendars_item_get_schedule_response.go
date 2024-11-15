package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemCalendarsItemGetSchedulePostResponseable instead.
type ItemCalendarsItemGetScheduleResponse struct {
    ItemCalendarsItemGetSchedulePostResponse
}
// NewItemCalendarsItemGetScheduleResponse instantiates a new ItemCalendarsItemGetScheduleResponse and sets the default values.
func NewItemCalendarsItemGetScheduleResponse()(*ItemCalendarsItemGetScheduleResponse) {
    m := &ItemCalendarsItemGetScheduleResponse{
        ItemCalendarsItemGetSchedulePostResponse: *NewItemCalendarsItemGetSchedulePostResponse(),
    }
    return m
}
// CreateItemCalendarsItemGetScheduleResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemCalendarsItemGetScheduleResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemCalendarsItemGetScheduleResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemCalendarsItemGetSchedulePostResponseable instead.
type ItemCalendarsItemGetScheduleResponseable interface {
    ItemCalendarsItemGetSchedulePostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemCalendarGetSchedulePostResponseable instead.
type ItemCalendarGetScheduleResponse struct {
    ItemCalendarGetSchedulePostResponse
}
// NewItemCalendarGetScheduleResponse instantiates a new ItemCalendarGetScheduleResponse and sets the default values.
func NewItemCalendarGetScheduleResponse()(*ItemCalendarGetScheduleResponse) {
    m := &ItemCalendarGetScheduleResponse{
        ItemCalendarGetSchedulePostResponse: *NewItemCalendarGetSchedulePostResponse(),
    }
    return m
}
// CreateItemCalendarGetScheduleResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemCalendarGetScheduleResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemCalendarGetScheduleResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemCalendarGetSchedulePostResponseable instead.
type ItemCalendarGetScheduleResponseable interface {
    ItemCalendarGetSchedulePostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

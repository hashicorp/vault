package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemCalendarsItemCalendarViewDeltaGetResponseable instead.
type ItemCalendarsItemCalendarViewDeltaResponse struct {
    ItemCalendarsItemCalendarViewDeltaGetResponse
}
// NewItemCalendarsItemCalendarViewDeltaResponse instantiates a new ItemCalendarsItemCalendarViewDeltaResponse and sets the default values.
func NewItemCalendarsItemCalendarViewDeltaResponse()(*ItemCalendarsItemCalendarViewDeltaResponse) {
    m := &ItemCalendarsItemCalendarViewDeltaResponse{
        ItemCalendarsItemCalendarViewDeltaGetResponse: *NewItemCalendarsItemCalendarViewDeltaGetResponse(),
    }
    return m
}
// CreateItemCalendarsItemCalendarViewDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemCalendarsItemCalendarViewDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemCalendarsItemCalendarViewDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemCalendarsItemCalendarViewDeltaGetResponseable instead.
type ItemCalendarsItemCalendarViewDeltaResponseable interface {
    ItemCalendarsItemCalendarViewDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

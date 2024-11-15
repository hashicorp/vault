package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemCalendarsItemCalendarViewItemInstancesDeltaGetResponseable instead.
type ItemCalendarsItemCalendarViewItemInstancesDeltaResponse struct {
    ItemCalendarsItemCalendarViewItemInstancesDeltaGetResponse
}
// NewItemCalendarsItemCalendarViewItemInstancesDeltaResponse instantiates a new ItemCalendarsItemCalendarViewItemInstancesDeltaResponse and sets the default values.
func NewItemCalendarsItemCalendarViewItemInstancesDeltaResponse()(*ItemCalendarsItemCalendarViewItemInstancesDeltaResponse) {
    m := &ItemCalendarsItemCalendarViewItemInstancesDeltaResponse{
        ItemCalendarsItemCalendarViewItemInstancesDeltaGetResponse: *NewItemCalendarsItemCalendarViewItemInstancesDeltaGetResponse(),
    }
    return m
}
// CreateItemCalendarsItemCalendarViewItemInstancesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemCalendarsItemCalendarViewItemInstancesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemCalendarsItemCalendarViewItemInstancesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemCalendarsItemCalendarViewItemInstancesDeltaGetResponseable instead.
type ItemCalendarsItemCalendarViewItemInstancesDeltaResponseable interface {
    ItemCalendarsItemCalendarViewItemInstancesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemCalendarEventsItemInstancesDeltaGetResponseable instead.
type ItemCalendarEventsItemInstancesDeltaResponse struct {
    ItemCalendarEventsItemInstancesDeltaGetResponse
}
// NewItemCalendarEventsItemInstancesDeltaResponse instantiates a new ItemCalendarEventsItemInstancesDeltaResponse and sets the default values.
func NewItemCalendarEventsItemInstancesDeltaResponse()(*ItemCalendarEventsItemInstancesDeltaResponse) {
    m := &ItemCalendarEventsItemInstancesDeltaResponse{
        ItemCalendarEventsItemInstancesDeltaGetResponse: *NewItemCalendarEventsItemInstancesDeltaGetResponse(),
    }
    return m
}
// CreateItemCalendarEventsItemInstancesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemCalendarEventsItemInstancesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemCalendarEventsItemInstancesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemCalendarEventsItemInstancesDeltaGetResponseable instead.
type ItemCalendarEventsItemInstancesDeltaResponseable interface {
    ItemCalendarEventsItemInstancesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

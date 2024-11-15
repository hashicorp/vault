package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemCalendarsItemEventsDeltaGetResponseable instead.
type ItemCalendarsItemEventsDeltaResponse struct {
    ItemCalendarsItemEventsDeltaGetResponse
}
// NewItemCalendarsItemEventsDeltaResponse instantiates a new ItemCalendarsItemEventsDeltaResponse and sets the default values.
func NewItemCalendarsItemEventsDeltaResponse()(*ItemCalendarsItemEventsDeltaResponse) {
    m := &ItemCalendarsItemEventsDeltaResponse{
        ItemCalendarsItemEventsDeltaGetResponse: *NewItemCalendarsItemEventsDeltaGetResponse(),
    }
    return m
}
// CreateItemCalendarsItemEventsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemCalendarsItemEventsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemCalendarsItemEventsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemCalendarsItemEventsDeltaGetResponseable instead.
type ItemCalendarsItemEventsDeltaResponseable interface {
    ItemCalendarsItemEventsDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

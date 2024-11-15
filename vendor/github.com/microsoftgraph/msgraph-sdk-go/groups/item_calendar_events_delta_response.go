package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemCalendarEventsDeltaGetResponseable instead.
type ItemCalendarEventsDeltaResponse struct {
    ItemCalendarEventsDeltaGetResponse
}
// NewItemCalendarEventsDeltaResponse instantiates a new ItemCalendarEventsDeltaResponse and sets the default values.
func NewItemCalendarEventsDeltaResponse()(*ItemCalendarEventsDeltaResponse) {
    m := &ItemCalendarEventsDeltaResponse{
        ItemCalendarEventsDeltaGetResponse: *NewItemCalendarEventsDeltaGetResponse(),
    }
    return m
}
// CreateItemCalendarEventsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemCalendarEventsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemCalendarEventsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemCalendarEventsDeltaGetResponseable instead.
type ItemCalendarEventsDeltaResponseable interface {
    ItemCalendarEventsDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemEventsDeltaGetResponseable instead.
type ItemEventsDeltaResponse struct {
    ItemEventsDeltaGetResponse
}
// NewItemEventsDeltaResponse instantiates a new ItemEventsDeltaResponse and sets the default values.
func NewItemEventsDeltaResponse()(*ItemEventsDeltaResponse) {
    m := &ItemEventsDeltaResponse{
        ItemEventsDeltaGetResponse: *NewItemEventsDeltaGetResponse(),
    }
    return m
}
// CreateItemEventsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemEventsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemEventsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemEventsDeltaGetResponseable instead.
type ItemEventsDeltaResponseable interface {
    ItemEventsDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

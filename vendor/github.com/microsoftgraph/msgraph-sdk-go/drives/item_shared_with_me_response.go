package drives

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSharedWithMeGetResponseable instead.
type ItemSharedWithMeResponse struct {
    ItemSharedWithMeGetResponse
}
// NewItemSharedWithMeResponse instantiates a new ItemSharedWithMeResponse and sets the default values.
func NewItemSharedWithMeResponse()(*ItemSharedWithMeResponse) {
    m := &ItemSharedWithMeResponse{
        ItemSharedWithMeGetResponse: *NewItemSharedWithMeGetResponse(),
    }
    return m
}
// CreateItemSharedWithMeResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSharedWithMeResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSharedWithMeResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSharedWithMeGetResponseable instead.
type ItemSharedWithMeResponseable interface {
    ItemSharedWithMeGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

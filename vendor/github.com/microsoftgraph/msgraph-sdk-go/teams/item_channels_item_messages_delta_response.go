package teams

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemChannelsItemMessagesDeltaGetResponseable instead.
type ItemChannelsItemMessagesDeltaResponse struct {
    ItemChannelsItemMessagesDeltaGetResponse
}
// NewItemChannelsItemMessagesDeltaResponse instantiates a new ItemChannelsItemMessagesDeltaResponse and sets the default values.
func NewItemChannelsItemMessagesDeltaResponse()(*ItemChannelsItemMessagesDeltaResponse) {
    m := &ItemChannelsItemMessagesDeltaResponse{
        ItemChannelsItemMessagesDeltaGetResponse: *NewItemChannelsItemMessagesDeltaGetResponse(),
    }
    return m
}
// CreateItemChannelsItemMessagesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemChannelsItemMessagesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemChannelsItemMessagesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemChannelsItemMessagesDeltaGetResponseable instead.
type ItemChannelsItemMessagesDeltaResponseable interface {
    ItemChannelsItemMessagesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

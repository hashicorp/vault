package teams

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemPrimaryChannelMessagesDeltaGetResponseable instead.
type ItemPrimaryChannelMessagesDeltaResponse struct {
    ItemPrimaryChannelMessagesDeltaGetResponse
}
// NewItemPrimaryChannelMessagesDeltaResponse instantiates a new ItemPrimaryChannelMessagesDeltaResponse and sets the default values.
func NewItemPrimaryChannelMessagesDeltaResponse()(*ItemPrimaryChannelMessagesDeltaResponse) {
    m := &ItemPrimaryChannelMessagesDeltaResponse{
        ItemPrimaryChannelMessagesDeltaGetResponse: *NewItemPrimaryChannelMessagesDeltaGetResponse(),
    }
    return m
}
// CreateItemPrimaryChannelMessagesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemPrimaryChannelMessagesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemPrimaryChannelMessagesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemPrimaryChannelMessagesDeltaGetResponseable instead.
type ItemPrimaryChannelMessagesDeltaResponseable interface {
    ItemPrimaryChannelMessagesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

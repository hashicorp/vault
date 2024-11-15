package teams

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemPrimaryChannelMessagesItemRepliesDeltaGetResponseable instead.
type ItemPrimaryChannelMessagesItemRepliesDeltaResponse struct {
    ItemPrimaryChannelMessagesItemRepliesDeltaGetResponse
}
// NewItemPrimaryChannelMessagesItemRepliesDeltaResponse instantiates a new ItemPrimaryChannelMessagesItemRepliesDeltaResponse and sets the default values.
func NewItemPrimaryChannelMessagesItemRepliesDeltaResponse()(*ItemPrimaryChannelMessagesItemRepliesDeltaResponse) {
    m := &ItemPrimaryChannelMessagesItemRepliesDeltaResponse{
        ItemPrimaryChannelMessagesItemRepliesDeltaGetResponse: *NewItemPrimaryChannelMessagesItemRepliesDeltaGetResponse(),
    }
    return m
}
// CreateItemPrimaryChannelMessagesItemRepliesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemPrimaryChannelMessagesItemRepliesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemPrimaryChannelMessagesItemRepliesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemPrimaryChannelMessagesItemRepliesDeltaGetResponseable instead.
type ItemPrimaryChannelMessagesItemRepliesDeltaResponseable interface {
    ItemPrimaryChannelMessagesItemRepliesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemChatsItemMessagesItemRepliesDeltaGetResponseable instead.
type ItemChatsItemMessagesItemRepliesDeltaResponse struct {
    ItemChatsItemMessagesItemRepliesDeltaGetResponse
}
// NewItemChatsItemMessagesItemRepliesDeltaResponse instantiates a new ItemChatsItemMessagesItemRepliesDeltaResponse and sets the default values.
func NewItemChatsItemMessagesItemRepliesDeltaResponse()(*ItemChatsItemMessagesItemRepliesDeltaResponse) {
    m := &ItemChatsItemMessagesItemRepliesDeltaResponse{
        ItemChatsItemMessagesItemRepliesDeltaGetResponse: *NewItemChatsItemMessagesItemRepliesDeltaGetResponse(),
    }
    return m
}
// CreateItemChatsItemMessagesItemRepliesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemChatsItemMessagesItemRepliesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemChatsItemMessagesItemRepliesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemChatsItemMessagesItemRepliesDeltaGetResponseable instead.
type ItemChatsItemMessagesItemRepliesDeltaResponseable interface {
    ItemChatsItemMessagesItemRepliesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

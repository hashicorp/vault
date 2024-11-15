package chats

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemMessagesItemRepliesDeltaGetResponseable instead.
type ItemMessagesItemRepliesDeltaResponse struct {
    ItemMessagesItemRepliesDeltaGetResponse
}
// NewItemMessagesItemRepliesDeltaResponse instantiates a new ItemMessagesItemRepliesDeltaResponse and sets the default values.
func NewItemMessagesItemRepliesDeltaResponse()(*ItemMessagesItemRepliesDeltaResponse) {
    m := &ItemMessagesItemRepliesDeltaResponse{
        ItemMessagesItemRepliesDeltaGetResponse: *NewItemMessagesItemRepliesDeltaGetResponse(),
    }
    return m
}
// CreateItemMessagesItemRepliesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemMessagesItemRepliesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemMessagesItemRepliesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemMessagesItemRepliesDeltaGetResponseable instead.
type ItemMessagesItemRepliesDeltaResponseable interface {
    ItemMessagesItemRepliesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

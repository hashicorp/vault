package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemChatsGetAllRetainedMessagesGetResponseable instead.
type ItemChatsGetAllRetainedMessagesResponse struct {
    ItemChatsGetAllRetainedMessagesGetResponse
}
// NewItemChatsGetAllRetainedMessagesResponse instantiates a new ItemChatsGetAllRetainedMessagesResponse and sets the default values.
func NewItemChatsGetAllRetainedMessagesResponse()(*ItemChatsGetAllRetainedMessagesResponse) {
    m := &ItemChatsGetAllRetainedMessagesResponse{
        ItemChatsGetAllRetainedMessagesGetResponse: *NewItemChatsGetAllRetainedMessagesGetResponse(),
    }
    return m
}
// CreateItemChatsGetAllRetainedMessagesResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemChatsGetAllRetainedMessagesResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemChatsGetAllRetainedMessagesResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemChatsGetAllRetainedMessagesGetResponseable instead.
type ItemChatsGetAllRetainedMessagesResponseable interface {
    ItemChatsGetAllRetainedMessagesGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

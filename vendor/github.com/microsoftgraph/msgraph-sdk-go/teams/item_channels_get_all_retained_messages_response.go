package teams

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemChannelsGetAllRetainedMessagesGetResponseable instead.
type ItemChannelsGetAllRetainedMessagesResponse struct {
    ItemChannelsGetAllRetainedMessagesGetResponse
}
// NewItemChannelsGetAllRetainedMessagesResponse instantiates a new ItemChannelsGetAllRetainedMessagesResponse and sets the default values.
func NewItemChannelsGetAllRetainedMessagesResponse()(*ItemChannelsGetAllRetainedMessagesResponse) {
    m := &ItemChannelsGetAllRetainedMessagesResponse{
        ItemChannelsGetAllRetainedMessagesGetResponse: *NewItemChannelsGetAllRetainedMessagesGetResponse(),
    }
    return m
}
// CreateItemChannelsGetAllRetainedMessagesResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemChannelsGetAllRetainedMessagesResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemChannelsGetAllRetainedMessagesResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemChannelsGetAllRetainedMessagesGetResponseable instead.
type ItemChannelsGetAllRetainedMessagesResponseable interface {
    ItemChannelsGetAllRetainedMessagesGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

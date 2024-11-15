package teams

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemChannelsGetAllMessagesGetResponseable instead.
type ItemChannelsGetAllMessagesResponse struct {
    ItemChannelsGetAllMessagesGetResponse
}
// NewItemChannelsGetAllMessagesResponse instantiates a new ItemChannelsGetAllMessagesResponse and sets the default values.
func NewItemChannelsGetAllMessagesResponse()(*ItemChannelsGetAllMessagesResponse) {
    m := &ItemChannelsGetAllMessagesResponse{
        ItemChannelsGetAllMessagesGetResponse: *NewItemChannelsGetAllMessagesGetResponse(),
    }
    return m
}
// CreateItemChannelsGetAllMessagesResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemChannelsGetAllMessagesResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemChannelsGetAllMessagesResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemChannelsGetAllMessagesGetResponseable instead.
type ItemChannelsGetAllMessagesResponseable interface {
    ItemChannelsGetAllMessagesGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

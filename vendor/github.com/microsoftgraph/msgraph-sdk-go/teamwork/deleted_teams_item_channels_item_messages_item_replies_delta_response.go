package teamwork

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use DeletedTeamsItemChannelsItemMessagesItemRepliesDeltaGetResponseable instead.
type DeletedTeamsItemChannelsItemMessagesItemRepliesDeltaResponse struct {
    DeletedTeamsItemChannelsItemMessagesItemRepliesDeltaGetResponse
}
// NewDeletedTeamsItemChannelsItemMessagesItemRepliesDeltaResponse instantiates a new DeletedTeamsItemChannelsItemMessagesItemRepliesDeltaResponse and sets the default values.
func NewDeletedTeamsItemChannelsItemMessagesItemRepliesDeltaResponse()(*DeletedTeamsItemChannelsItemMessagesItemRepliesDeltaResponse) {
    m := &DeletedTeamsItemChannelsItemMessagesItemRepliesDeltaResponse{
        DeletedTeamsItemChannelsItemMessagesItemRepliesDeltaGetResponse: *NewDeletedTeamsItemChannelsItemMessagesItemRepliesDeltaGetResponse(),
    }
    return m
}
// CreateDeletedTeamsItemChannelsItemMessagesItemRepliesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeletedTeamsItemChannelsItemMessagesItemRepliesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeletedTeamsItemChannelsItemMessagesItemRepliesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use DeletedTeamsItemChannelsItemMessagesItemRepliesDeltaGetResponseable instead.
type DeletedTeamsItemChannelsItemMessagesItemRepliesDeltaResponseable interface {
    DeletedTeamsItemChannelsItemMessagesItemRepliesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

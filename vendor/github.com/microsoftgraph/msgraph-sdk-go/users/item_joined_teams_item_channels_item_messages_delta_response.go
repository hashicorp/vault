package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemJoinedTeamsItemChannelsItemMessagesDeltaGetResponseable instead.
type ItemJoinedTeamsItemChannelsItemMessagesDeltaResponse struct {
    ItemJoinedTeamsItemChannelsItemMessagesDeltaGetResponse
}
// NewItemJoinedTeamsItemChannelsItemMessagesDeltaResponse instantiates a new ItemJoinedTeamsItemChannelsItemMessagesDeltaResponse and sets the default values.
func NewItemJoinedTeamsItemChannelsItemMessagesDeltaResponse()(*ItemJoinedTeamsItemChannelsItemMessagesDeltaResponse) {
    m := &ItemJoinedTeamsItemChannelsItemMessagesDeltaResponse{
        ItemJoinedTeamsItemChannelsItemMessagesDeltaGetResponse: *NewItemJoinedTeamsItemChannelsItemMessagesDeltaGetResponse(),
    }
    return m
}
// CreateItemJoinedTeamsItemChannelsItemMessagesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemJoinedTeamsItemChannelsItemMessagesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemJoinedTeamsItemChannelsItemMessagesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemJoinedTeamsItemChannelsItemMessagesDeltaGetResponseable instead.
type ItemJoinedTeamsItemChannelsItemMessagesDeltaResponseable interface {
    ItemJoinedTeamsItemChannelsItemMessagesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

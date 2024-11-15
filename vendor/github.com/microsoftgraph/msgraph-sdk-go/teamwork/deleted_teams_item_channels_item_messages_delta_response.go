package teamwork

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use DeletedTeamsItemChannelsItemMessagesDeltaGetResponseable instead.
type DeletedTeamsItemChannelsItemMessagesDeltaResponse struct {
    DeletedTeamsItemChannelsItemMessagesDeltaGetResponse
}
// NewDeletedTeamsItemChannelsItemMessagesDeltaResponse instantiates a new DeletedTeamsItemChannelsItemMessagesDeltaResponse and sets the default values.
func NewDeletedTeamsItemChannelsItemMessagesDeltaResponse()(*DeletedTeamsItemChannelsItemMessagesDeltaResponse) {
    m := &DeletedTeamsItemChannelsItemMessagesDeltaResponse{
        DeletedTeamsItemChannelsItemMessagesDeltaGetResponse: *NewDeletedTeamsItemChannelsItemMessagesDeltaGetResponse(),
    }
    return m
}
// CreateDeletedTeamsItemChannelsItemMessagesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeletedTeamsItemChannelsItemMessagesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeletedTeamsItemChannelsItemMessagesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use DeletedTeamsItemChannelsItemMessagesDeltaGetResponseable instead.
type DeletedTeamsItemChannelsItemMessagesDeltaResponseable interface {
    DeletedTeamsItemChannelsItemMessagesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

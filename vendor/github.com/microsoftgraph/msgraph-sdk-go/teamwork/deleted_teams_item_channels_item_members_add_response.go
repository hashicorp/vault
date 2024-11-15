package teamwork

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use DeletedTeamsItemChannelsItemMembersAddPostResponseable instead.
type DeletedTeamsItemChannelsItemMembersAddResponse struct {
    DeletedTeamsItemChannelsItemMembersAddPostResponse
}
// NewDeletedTeamsItemChannelsItemMembersAddResponse instantiates a new DeletedTeamsItemChannelsItemMembersAddResponse and sets the default values.
func NewDeletedTeamsItemChannelsItemMembersAddResponse()(*DeletedTeamsItemChannelsItemMembersAddResponse) {
    m := &DeletedTeamsItemChannelsItemMembersAddResponse{
        DeletedTeamsItemChannelsItemMembersAddPostResponse: *NewDeletedTeamsItemChannelsItemMembersAddPostResponse(),
    }
    return m
}
// CreateDeletedTeamsItemChannelsItemMembersAddResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeletedTeamsItemChannelsItemMembersAddResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeletedTeamsItemChannelsItemMembersAddResponse(), nil
}
// Deprecated: This class is obsolete. Use DeletedTeamsItemChannelsItemMembersAddPostResponseable instead.
type DeletedTeamsItemChannelsItemMembersAddResponseable interface {
    DeletedTeamsItemChannelsItemMembersAddPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

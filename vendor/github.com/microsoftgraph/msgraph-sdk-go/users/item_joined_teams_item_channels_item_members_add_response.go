package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemJoinedTeamsItemChannelsItemMembersAddPostResponseable instead.
type ItemJoinedTeamsItemChannelsItemMembersAddResponse struct {
    ItemJoinedTeamsItemChannelsItemMembersAddPostResponse
}
// NewItemJoinedTeamsItemChannelsItemMembersAddResponse instantiates a new ItemJoinedTeamsItemChannelsItemMembersAddResponse and sets the default values.
func NewItemJoinedTeamsItemChannelsItemMembersAddResponse()(*ItemJoinedTeamsItemChannelsItemMembersAddResponse) {
    m := &ItemJoinedTeamsItemChannelsItemMembersAddResponse{
        ItemJoinedTeamsItemChannelsItemMembersAddPostResponse: *NewItemJoinedTeamsItemChannelsItemMembersAddPostResponse(),
    }
    return m
}
// CreateItemJoinedTeamsItemChannelsItemMembersAddResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemJoinedTeamsItemChannelsItemMembersAddResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemJoinedTeamsItemChannelsItemMembersAddResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemJoinedTeamsItemChannelsItemMembersAddPostResponseable instead.
type ItemJoinedTeamsItemChannelsItemMembersAddResponseable interface {
    ItemJoinedTeamsItemChannelsItemMembersAddPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

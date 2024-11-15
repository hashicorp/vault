package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemJoinedTeamsItemPrimaryChannelMembersAddPostResponseable instead.
type ItemJoinedTeamsItemPrimaryChannelMembersAddResponse struct {
    ItemJoinedTeamsItemPrimaryChannelMembersAddPostResponse
}
// NewItemJoinedTeamsItemPrimaryChannelMembersAddResponse instantiates a new ItemJoinedTeamsItemPrimaryChannelMembersAddResponse and sets the default values.
func NewItemJoinedTeamsItemPrimaryChannelMembersAddResponse()(*ItemJoinedTeamsItemPrimaryChannelMembersAddResponse) {
    m := &ItemJoinedTeamsItemPrimaryChannelMembersAddResponse{
        ItemJoinedTeamsItemPrimaryChannelMembersAddPostResponse: *NewItemJoinedTeamsItemPrimaryChannelMembersAddPostResponse(),
    }
    return m
}
// CreateItemJoinedTeamsItemPrimaryChannelMembersAddResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemJoinedTeamsItemPrimaryChannelMembersAddResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemJoinedTeamsItemPrimaryChannelMembersAddResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemJoinedTeamsItemPrimaryChannelMembersAddPostResponseable instead.
type ItemJoinedTeamsItemPrimaryChannelMembersAddResponseable interface {
    ItemJoinedTeamsItemPrimaryChannelMembersAddPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemChatsItemMembersAddPostResponseable instead.
type ItemChatsItemMembersAddResponse struct {
    ItemChatsItemMembersAddPostResponse
}
// NewItemChatsItemMembersAddResponse instantiates a new ItemChatsItemMembersAddResponse and sets the default values.
func NewItemChatsItemMembersAddResponse()(*ItemChatsItemMembersAddResponse) {
    m := &ItemChatsItemMembersAddResponse{
        ItemChatsItemMembersAddPostResponse: *NewItemChatsItemMembersAddPostResponse(),
    }
    return m
}
// CreateItemChatsItemMembersAddResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemChatsItemMembersAddResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemChatsItemMembersAddResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemChatsItemMembersAddPostResponseable instead.
type ItemChatsItemMembersAddResponseable interface {
    ItemChatsItemMembersAddPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

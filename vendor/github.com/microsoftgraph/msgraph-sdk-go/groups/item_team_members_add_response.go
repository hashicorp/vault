package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemTeamMembersAddPostResponseable instead.
type ItemTeamMembersAddResponse struct {
    ItemTeamMembersAddPostResponse
}
// NewItemTeamMembersAddResponse instantiates a new ItemTeamMembersAddResponse and sets the default values.
func NewItemTeamMembersAddResponse()(*ItemTeamMembersAddResponse) {
    m := &ItemTeamMembersAddResponse{
        ItemTeamMembersAddPostResponse: *NewItemTeamMembersAddPostResponse(),
    }
    return m
}
// CreateItemTeamMembersAddResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemTeamMembersAddResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemTeamMembersAddResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemTeamMembersAddPostResponseable instead.
type ItemTeamMembersAddResponseable interface {
    ItemTeamMembersAddPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

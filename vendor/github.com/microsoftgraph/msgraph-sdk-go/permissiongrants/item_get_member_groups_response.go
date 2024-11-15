package permissiongrants

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemGetMemberGroupsPostResponseable instead.
type ItemGetMemberGroupsResponse struct {
    ItemGetMemberGroupsPostResponse
}
// NewItemGetMemberGroupsResponse instantiates a new ItemGetMemberGroupsResponse and sets the default values.
func NewItemGetMemberGroupsResponse()(*ItemGetMemberGroupsResponse) {
    m := &ItemGetMemberGroupsResponse{
        ItemGetMemberGroupsPostResponse: *NewItemGetMemberGroupsPostResponse(),
    }
    return m
}
// CreateItemGetMemberGroupsResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemGetMemberGroupsResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemGetMemberGroupsResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemGetMemberGroupsPostResponseable instead.
type ItemGetMemberGroupsResponseable interface {
    ItemGetMemberGroupsPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

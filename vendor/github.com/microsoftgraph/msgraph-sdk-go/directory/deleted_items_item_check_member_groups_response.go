package directory

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use DeletedItemsItemCheckMemberGroupsPostResponseable instead.
type DeletedItemsItemCheckMemberGroupsResponse struct {
    DeletedItemsItemCheckMemberGroupsPostResponse
}
// NewDeletedItemsItemCheckMemberGroupsResponse instantiates a new DeletedItemsItemCheckMemberGroupsResponse and sets the default values.
func NewDeletedItemsItemCheckMemberGroupsResponse()(*DeletedItemsItemCheckMemberGroupsResponse) {
    m := &DeletedItemsItemCheckMemberGroupsResponse{
        DeletedItemsItemCheckMemberGroupsPostResponse: *NewDeletedItemsItemCheckMemberGroupsPostResponse(),
    }
    return m
}
// CreateDeletedItemsItemCheckMemberGroupsResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeletedItemsItemCheckMemberGroupsResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeletedItemsItemCheckMemberGroupsResponse(), nil
}
// Deprecated: This class is obsolete. Use DeletedItemsItemCheckMemberGroupsPostResponseable instead.
type DeletedItemsItemCheckMemberGroupsResponseable interface {
    DeletedItemsItemCheckMemberGroupsPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemFollowedSitesAddPostResponseable instead.
type ItemFollowedSitesAddResponse struct {
    ItemFollowedSitesAddPostResponse
}
// NewItemFollowedSitesAddResponse instantiates a new ItemFollowedSitesAddResponse and sets the default values.
func NewItemFollowedSitesAddResponse()(*ItemFollowedSitesAddResponse) {
    m := &ItemFollowedSitesAddResponse{
        ItemFollowedSitesAddPostResponse: *NewItemFollowedSitesAddPostResponse(),
    }
    return m
}
// CreateItemFollowedSitesAddResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemFollowedSitesAddResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemFollowedSitesAddResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemFollowedSitesAddPostResponseable instead.
type ItemFollowedSitesAddResponseable interface {
    ItemFollowedSitesAddPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

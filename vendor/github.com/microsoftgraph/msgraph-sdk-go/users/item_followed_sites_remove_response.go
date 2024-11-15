package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemFollowedSitesRemovePostResponseable instead.
type ItemFollowedSitesRemoveResponse struct {
    ItemFollowedSitesRemovePostResponse
}
// NewItemFollowedSitesRemoveResponse instantiates a new ItemFollowedSitesRemoveResponse and sets the default values.
func NewItemFollowedSitesRemoveResponse()(*ItemFollowedSitesRemoveResponse) {
    m := &ItemFollowedSitesRemoveResponse{
        ItemFollowedSitesRemovePostResponse: *NewItemFollowedSitesRemovePostResponse(),
    }
    return m
}
// CreateItemFollowedSitesRemoveResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemFollowedSitesRemoveResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemFollowedSitesRemoveResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemFollowedSitesRemovePostResponseable instead.
type ItemFollowedSitesRemoveResponseable interface {
    ItemFollowedSitesRemovePostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSitesRemovePostResponseable instead.
type ItemSitesRemoveResponse struct {
    ItemSitesRemovePostResponse
}
// NewItemSitesRemoveResponse instantiates a new ItemSitesRemoveResponse and sets the default values.
func NewItemSitesRemoveResponse()(*ItemSitesRemoveResponse) {
    m := &ItemSitesRemoveResponse{
        ItemSitesRemovePostResponse: *NewItemSitesRemovePostResponse(),
    }
    return m
}
// CreateItemSitesRemoveResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSitesRemoveResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSitesRemoveResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSitesRemovePostResponseable instead.
type ItemSitesRemoveResponseable interface {
    ItemSitesRemovePostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

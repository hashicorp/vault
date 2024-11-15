package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSitesAddPostResponseable instead.
type ItemSitesAddResponse struct {
    ItemSitesAddPostResponse
}
// NewItemSitesAddResponse instantiates a new ItemSitesAddResponse and sets the default values.
func NewItemSitesAddResponse()(*ItemSitesAddResponse) {
    m := &ItemSitesAddResponse{
        ItemSitesAddPostResponse: *NewItemSitesAddPostResponse(),
    }
    return m
}
// CreateItemSitesAddResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSitesAddResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSitesAddResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSitesAddPostResponseable instead.
type ItemSitesAddResponseable interface {
    ItemSitesAddPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

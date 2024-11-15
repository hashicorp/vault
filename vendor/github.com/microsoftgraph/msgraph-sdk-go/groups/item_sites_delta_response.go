package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSitesDeltaGetResponseable instead.
type ItemSitesDeltaResponse struct {
    ItemSitesDeltaGetResponse
}
// NewItemSitesDeltaResponse instantiates a new ItemSitesDeltaResponse and sets the default values.
func NewItemSitesDeltaResponse()(*ItemSitesDeltaResponse) {
    m := &ItemSitesDeltaResponse{
        ItemSitesDeltaGetResponse: *NewItemSitesDeltaGetResponse(),
    }
    return m
}
// CreateItemSitesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSitesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSitesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSitesDeltaGetResponseable instead.
type ItemSitesDeltaResponseable interface {
    ItemSitesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

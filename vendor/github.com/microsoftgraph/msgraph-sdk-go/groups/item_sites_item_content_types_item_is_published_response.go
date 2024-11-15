package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSitesItemContentTypesItemIsPublishedGetResponseable instead.
type ItemSitesItemContentTypesItemIsPublishedResponse struct {
    ItemSitesItemContentTypesItemIsPublishedGetResponse
}
// NewItemSitesItemContentTypesItemIsPublishedResponse instantiates a new ItemSitesItemContentTypesItemIsPublishedResponse and sets the default values.
func NewItemSitesItemContentTypesItemIsPublishedResponse()(*ItemSitesItemContentTypesItemIsPublishedResponse) {
    m := &ItemSitesItemContentTypesItemIsPublishedResponse{
        ItemSitesItemContentTypesItemIsPublishedGetResponse: *NewItemSitesItemContentTypesItemIsPublishedGetResponse(),
    }
    return m
}
// CreateItemSitesItemContentTypesItemIsPublishedResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSitesItemContentTypesItemIsPublishedResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSitesItemContentTypesItemIsPublishedResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSitesItemContentTypesItemIsPublishedGetResponseable instead.
type ItemSitesItemContentTypesItemIsPublishedResponseable interface {
    ItemSitesItemContentTypesItemIsPublishedGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSitesItemListsItemContentTypesItemIsPublishedGetResponseable instead.
type ItemSitesItemListsItemContentTypesItemIsPublishedResponse struct {
    ItemSitesItemListsItemContentTypesItemIsPublishedGetResponse
}
// NewItemSitesItemListsItemContentTypesItemIsPublishedResponse instantiates a new ItemSitesItemListsItemContentTypesItemIsPublishedResponse and sets the default values.
func NewItemSitesItemListsItemContentTypesItemIsPublishedResponse()(*ItemSitesItemListsItemContentTypesItemIsPublishedResponse) {
    m := &ItemSitesItemListsItemContentTypesItemIsPublishedResponse{
        ItemSitesItemListsItemContentTypesItemIsPublishedGetResponse: *NewItemSitesItemListsItemContentTypesItemIsPublishedGetResponse(),
    }
    return m
}
// CreateItemSitesItemListsItemContentTypesItemIsPublishedResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSitesItemListsItemContentTypesItemIsPublishedResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSitesItemListsItemContentTypesItemIsPublishedResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSitesItemListsItemContentTypesItemIsPublishedGetResponseable instead.
type ItemSitesItemListsItemContentTypesItemIsPublishedResponseable interface {
    ItemSitesItemListsItemContentTypesItemIsPublishedGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

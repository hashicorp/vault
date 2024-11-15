package sites

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemListsItemContentTypesItemIsPublishedGetResponseable instead.
type ItemListsItemContentTypesItemIsPublishedResponse struct {
    ItemListsItemContentTypesItemIsPublishedGetResponse
}
// NewItemListsItemContentTypesItemIsPublishedResponse instantiates a new ItemListsItemContentTypesItemIsPublishedResponse and sets the default values.
func NewItemListsItemContentTypesItemIsPublishedResponse()(*ItemListsItemContentTypesItemIsPublishedResponse) {
    m := &ItemListsItemContentTypesItemIsPublishedResponse{
        ItemListsItemContentTypesItemIsPublishedGetResponse: *NewItemListsItemContentTypesItemIsPublishedGetResponse(),
    }
    return m
}
// CreateItemListsItemContentTypesItemIsPublishedResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemListsItemContentTypesItemIsPublishedResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemListsItemContentTypesItemIsPublishedResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemListsItemContentTypesItemIsPublishedGetResponseable instead.
type ItemListsItemContentTypesItemIsPublishedResponseable interface {
    ItemListsItemContentTypesItemIsPublishedGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

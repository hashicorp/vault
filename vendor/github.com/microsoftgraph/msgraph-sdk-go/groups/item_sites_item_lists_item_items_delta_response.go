package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSitesItemListsItemItemsDeltaGetResponseable instead.
type ItemSitesItemListsItemItemsDeltaResponse struct {
    ItemSitesItemListsItemItemsDeltaGetResponse
}
// NewItemSitesItemListsItemItemsDeltaResponse instantiates a new ItemSitesItemListsItemItemsDeltaResponse and sets the default values.
func NewItemSitesItemListsItemItemsDeltaResponse()(*ItemSitesItemListsItemItemsDeltaResponse) {
    m := &ItemSitesItemListsItemItemsDeltaResponse{
        ItemSitesItemListsItemItemsDeltaGetResponse: *NewItemSitesItemListsItemItemsDeltaGetResponse(),
    }
    return m
}
// CreateItemSitesItemListsItemItemsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSitesItemListsItemItemsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSitesItemListsItemItemsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSitesItemListsItemItemsDeltaGetResponseable instead.
type ItemSitesItemListsItemItemsDeltaResponseable interface {
    ItemSitesItemListsItemItemsDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

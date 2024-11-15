package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSitesItemListsItemItemsDeltaWithTokenGetResponseable instead.
type ItemSitesItemListsItemItemsDeltaWithTokenResponse struct {
    ItemSitesItemListsItemItemsDeltaWithTokenGetResponse
}
// NewItemSitesItemListsItemItemsDeltaWithTokenResponse instantiates a new ItemSitesItemListsItemItemsDeltaWithTokenResponse and sets the default values.
func NewItemSitesItemListsItemItemsDeltaWithTokenResponse()(*ItemSitesItemListsItemItemsDeltaWithTokenResponse) {
    m := &ItemSitesItemListsItemItemsDeltaWithTokenResponse{
        ItemSitesItemListsItemItemsDeltaWithTokenGetResponse: *NewItemSitesItemListsItemItemsDeltaWithTokenGetResponse(),
    }
    return m
}
// CreateItemSitesItemListsItemItemsDeltaWithTokenResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSitesItemListsItemItemsDeltaWithTokenResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSitesItemListsItemItemsDeltaWithTokenResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSitesItemListsItemItemsDeltaWithTokenGetResponseable instead.
type ItemSitesItemListsItemItemsDeltaWithTokenResponseable interface {
    ItemSitesItemListsItemItemsDeltaWithTokenGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package sites

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemListsItemItemsDeltaGetResponseable instead.
type ItemListsItemItemsDeltaResponse struct {
    ItemListsItemItemsDeltaGetResponse
}
// NewItemListsItemItemsDeltaResponse instantiates a new ItemListsItemItemsDeltaResponse and sets the default values.
func NewItemListsItemItemsDeltaResponse()(*ItemListsItemItemsDeltaResponse) {
    m := &ItemListsItemItemsDeltaResponse{
        ItemListsItemItemsDeltaGetResponse: *NewItemListsItemItemsDeltaGetResponse(),
    }
    return m
}
// CreateItemListsItemItemsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemListsItemItemsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemListsItemItemsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemListsItemItemsDeltaGetResponseable instead.
type ItemListsItemItemsDeltaResponseable interface {
    ItemListsItemItemsDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

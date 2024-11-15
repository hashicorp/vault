package sites

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemListsItemItemsItemGetActivitiesByIntervalGetResponseable instead.
type ItemListsItemItemsItemGetActivitiesByIntervalResponse struct {
    ItemListsItemItemsItemGetActivitiesByIntervalGetResponse
}
// NewItemListsItemItemsItemGetActivitiesByIntervalResponse instantiates a new ItemListsItemItemsItemGetActivitiesByIntervalResponse and sets the default values.
func NewItemListsItemItemsItemGetActivitiesByIntervalResponse()(*ItemListsItemItemsItemGetActivitiesByIntervalResponse) {
    m := &ItemListsItemItemsItemGetActivitiesByIntervalResponse{
        ItemListsItemItemsItemGetActivitiesByIntervalGetResponse: *NewItemListsItemItemsItemGetActivitiesByIntervalGetResponse(),
    }
    return m
}
// CreateItemListsItemItemsItemGetActivitiesByIntervalResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemListsItemItemsItemGetActivitiesByIntervalResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemListsItemItemsItemGetActivitiesByIntervalResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemListsItemItemsItemGetActivitiesByIntervalGetResponseable instead.
type ItemListsItemItemsItemGetActivitiesByIntervalResponseable interface {
    ItemListsItemItemsItemGetActivitiesByIntervalGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

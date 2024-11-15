package drives

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemListItemsItemGetActivitiesByIntervalGetResponseable instead.
type ItemListItemsItemGetActivitiesByIntervalResponse struct {
    ItemListItemsItemGetActivitiesByIntervalGetResponse
}
// NewItemListItemsItemGetActivitiesByIntervalResponse instantiates a new ItemListItemsItemGetActivitiesByIntervalResponse and sets the default values.
func NewItemListItemsItemGetActivitiesByIntervalResponse()(*ItemListItemsItemGetActivitiesByIntervalResponse) {
    m := &ItemListItemsItemGetActivitiesByIntervalResponse{
        ItemListItemsItemGetActivitiesByIntervalGetResponse: *NewItemListItemsItemGetActivitiesByIntervalGetResponse(),
    }
    return m
}
// CreateItemListItemsItemGetActivitiesByIntervalResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemListItemsItemGetActivitiesByIntervalResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemListItemsItemGetActivitiesByIntervalResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemListItemsItemGetActivitiesByIntervalGetResponseable instead.
type ItemListItemsItemGetActivitiesByIntervalResponseable interface {
    ItemListItemsItemGetActivitiesByIntervalGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

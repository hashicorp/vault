package sites

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemListsItemContentTypesGetCompatibleHubContentTypesGetResponseable instead.
type ItemListsItemContentTypesGetCompatibleHubContentTypesResponse struct {
    ItemListsItemContentTypesGetCompatibleHubContentTypesGetResponse
}
// NewItemListsItemContentTypesGetCompatibleHubContentTypesResponse instantiates a new ItemListsItemContentTypesGetCompatibleHubContentTypesResponse and sets the default values.
func NewItemListsItemContentTypesGetCompatibleHubContentTypesResponse()(*ItemListsItemContentTypesGetCompatibleHubContentTypesResponse) {
    m := &ItemListsItemContentTypesGetCompatibleHubContentTypesResponse{
        ItemListsItemContentTypesGetCompatibleHubContentTypesGetResponse: *NewItemListsItemContentTypesGetCompatibleHubContentTypesGetResponse(),
    }
    return m
}
// CreateItemListsItemContentTypesGetCompatibleHubContentTypesResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemListsItemContentTypesGetCompatibleHubContentTypesResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemListsItemContentTypesGetCompatibleHubContentTypesResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemListsItemContentTypesGetCompatibleHubContentTypesGetResponseable instead.
type ItemListsItemContentTypesGetCompatibleHubContentTypesResponseable interface {
    ItemListsItemContentTypesGetCompatibleHubContentTypesGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package drives

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemItemsItemPermissionsItemGrantPostResponseable instead.
type ItemItemsItemPermissionsItemGrantResponse struct {
    ItemItemsItemPermissionsItemGrantPostResponse
}
// NewItemItemsItemPermissionsItemGrantResponse instantiates a new ItemItemsItemPermissionsItemGrantResponse and sets the default values.
func NewItemItemsItemPermissionsItemGrantResponse()(*ItemItemsItemPermissionsItemGrantResponse) {
    m := &ItemItemsItemPermissionsItemGrantResponse{
        ItemItemsItemPermissionsItemGrantPostResponse: *NewItemItemsItemPermissionsItemGrantPostResponse(),
    }
    return m
}
// CreateItemItemsItemPermissionsItemGrantResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemItemsItemPermissionsItemGrantResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemItemsItemPermissionsItemGrantResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemItemsItemPermissionsItemGrantPostResponseable instead.
type ItemItemsItemPermissionsItemGrantResponseable interface {
    ItemItemsItemPermissionsItemGrantPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

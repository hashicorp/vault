package sites

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemPermissionsItemGrantPostResponseable instead.
type ItemPermissionsItemGrantResponse struct {
    ItemPermissionsItemGrantPostResponse
}
// NewItemPermissionsItemGrantResponse instantiates a new ItemPermissionsItemGrantResponse and sets the default values.
func NewItemPermissionsItemGrantResponse()(*ItemPermissionsItemGrantResponse) {
    m := &ItemPermissionsItemGrantResponse{
        ItemPermissionsItemGrantPostResponse: *NewItemPermissionsItemGrantPostResponse(),
    }
    return m
}
// CreateItemPermissionsItemGrantResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemPermissionsItemGrantResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemPermissionsItemGrantResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemPermissionsItemGrantPostResponseable instead.
type ItemPermissionsItemGrantResponseable interface {
    ItemPermissionsItemGrantPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

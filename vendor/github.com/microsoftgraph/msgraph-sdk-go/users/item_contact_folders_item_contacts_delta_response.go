package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemContactFoldersItemContactsDeltaGetResponseable instead.
type ItemContactFoldersItemContactsDeltaResponse struct {
    ItemContactFoldersItemContactsDeltaGetResponse
}
// NewItemContactFoldersItemContactsDeltaResponse instantiates a new ItemContactFoldersItemContactsDeltaResponse and sets the default values.
func NewItemContactFoldersItemContactsDeltaResponse()(*ItemContactFoldersItemContactsDeltaResponse) {
    m := &ItemContactFoldersItemContactsDeltaResponse{
        ItemContactFoldersItemContactsDeltaGetResponse: *NewItemContactFoldersItemContactsDeltaGetResponse(),
    }
    return m
}
// CreateItemContactFoldersItemContactsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemContactFoldersItemContactsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemContactFoldersItemContactsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemContactFoldersItemContactsDeltaGetResponseable instead.
type ItemContactFoldersItemContactsDeltaResponseable interface {
    ItemContactFoldersItemContactsDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

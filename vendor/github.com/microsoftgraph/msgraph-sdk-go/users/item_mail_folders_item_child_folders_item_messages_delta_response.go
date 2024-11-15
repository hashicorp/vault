package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemMailFoldersItemChildFoldersItemMessagesDeltaGetResponseable instead.
type ItemMailFoldersItemChildFoldersItemMessagesDeltaResponse struct {
    ItemMailFoldersItemChildFoldersItemMessagesDeltaGetResponse
}
// NewItemMailFoldersItemChildFoldersItemMessagesDeltaResponse instantiates a new ItemMailFoldersItemChildFoldersItemMessagesDeltaResponse and sets the default values.
func NewItemMailFoldersItemChildFoldersItemMessagesDeltaResponse()(*ItemMailFoldersItemChildFoldersItemMessagesDeltaResponse) {
    m := &ItemMailFoldersItemChildFoldersItemMessagesDeltaResponse{
        ItemMailFoldersItemChildFoldersItemMessagesDeltaGetResponse: *NewItemMailFoldersItemChildFoldersItemMessagesDeltaGetResponse(),
    }
    return m
}
// CreateItemMailFoldersItemChildFoldersItemMessagesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemMailFoldersItemChildFoldersItemMessagesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemMailFoldersItemChildFoldersItemMessagesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemMailFoldersItemChildFoldersItemMessagesDeltaGetResponseable instead.
type ItemMailFoldersItemChildFoldersItemMessagesDeltaResponseable interface {
    ItemMailFoldersItemChildFoldersItemMessagesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

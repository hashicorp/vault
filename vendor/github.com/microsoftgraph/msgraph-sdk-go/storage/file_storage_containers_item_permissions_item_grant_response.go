package storage

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use FileStorageContainersItemPermissionsItemGrantPostResponseable instead.
type FileStorageContainersItemPermissionsItemGrantResponse struct {
    FileStorageContainersItemPermissionsItemGrantPostResponse
}
// NewFileStorageContainersItemPermissionsItemGrantResponse instantiates a new FileStorageContainersItemPermissionsItemGrantResponse and sets the default values.
func NewFileStorageContainersItemPermissionsItemGrantResponse()(*FileStorageContainersItemPermissionsItemGrantResponse) {
    m := &FileStorageContainersItemPermissionsItemGrantResponse{
        FileStorageContainersItemPermissionsItemGrantPostResponse: *NewFileStorageContainersItemPermissionsItemGrantPostResponse(),
    }
    return m
}
// CreateFileStorageContainersItemPermissionsItemGrantResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFileStorageContainersItemPermissionsItemGrantResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFileStorageContainersItemPermissionsItemGrantResponse(), nil
}
// Deprecated: This class is obsolete. Use FileStorageContainersItemPermissionsItemGrantPostResponseable instead.
type FileStorageContainersItemPermissionsItemGrantResponseable interface {
    FileStorageContainersItemPermissionsItemGrantPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

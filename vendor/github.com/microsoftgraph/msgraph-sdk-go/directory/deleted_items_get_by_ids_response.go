package directory

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use DeletedItemsGetByIdsPostResponseable instead.
type DeletedItemsGetByIdsResponse struct {
    DeletedItemsGetByIdsPostResponse
}
// NewDeletedItemsGetByIdsResponse instantiates a new DeletedItemsGetByIdsResponse and sets the default values.
func NewDeletedItemsGetByIdsResponse()(*DeletedItemsGetByIdsResponse) {
    m := &DeletedItemsGetByIdsResponse{
        DeletedItemsGetByIdsPostResponse: *NewDeletedItemsGetByIdsPostResponse(),
    }
    return m
}
// CreateDeletedItemsGetByIdsResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeletedItemsGetByIdsResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeletedItemsGetByIdsResponse(), nil
}
// Deprecated: This class is obsolete. Use DeletedItemsGetByIdsPostResponseable instead.
type DeletedItemsGetByIdsResponseable interface {
    DeletedItemsGetByIdsPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

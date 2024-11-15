package sites

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemContentTypesItemIsPublishedGetResponseable instead.
type ItemContentTypesItemIsPublishedResponse struct {
    ItemContentTypesItemIsPublishedGetResponse
}
// NewItemContentTypesItemIsPublishedResponse instantiates a new ItemContentTypesItemIsPublishedResponse and sets the default values.
func NewItemContentTypesItemIsPublishedResponse()(*ItemContentTypesItemIsPublishedResponse) {
    m := &ItemContentTypesItemIsPublishedResponse{
        ItemContentTypesItemIsPublishedGetResponse: *NewItemContentTypesItemIsPublishedGetResponse(),
    }
    return m
}
// CreateItemContentTypesItemIsPublishedResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemContentTypesItemIsPublishedResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemContentTypesItemIsPublishedResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemContentTypesItemIsPublishedGetResponseable instead.
type ItemContentTypesItemIsPublishedResponseable interface {
    ItemContentTypesItemIsPublishedGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

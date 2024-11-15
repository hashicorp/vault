package drives

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSearchWithQGetResponseable instead.
type ItemSearchWithQResponse struct {
    ItemSearchWithQGetResponse
}
// NewItemSearchWithQResponse instantiates a new ItemSearchWithQResponse and sets the default values.
func NewItemSearchWithQResponse()(*ItemSearchWithQResponse) {
    m := &ItemSearchWithQResponse{
        ItemSearchWithQGetResponse: *NewItemSearchWithQGetResponse(),
    }
    return m
}
// CreateItemSearchWithQResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSearchWithQResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSearchWithQResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSearchWithQGetResponseable instead.
type ItemSearchWithQResponseable interface {
    ItemSearchWithQGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

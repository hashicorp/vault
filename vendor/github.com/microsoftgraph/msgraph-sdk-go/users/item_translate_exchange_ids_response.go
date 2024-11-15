package users

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemTranslateExchangeIdsPostResponseable instead.
type ItemTranslateExchangeIdsResponse struct {
    ItemTranslateExchangeIdsPostResponse
}
// NewItemTranslateExchangeIdsResponse instantiates a new ItemTranslateExchangeIdsResponse and sets the default values.
func NewItemTranslateExchangeIdsResponse()(*ItemTranslateExchangeIdsResponse) {
    m := &ItemTranslateExchangeIdsResponse{
        ItemTranslateExchangeIdsPostResponse: *NewItemTranslateExchangeIdsPostResponse(),
    }
    return m
}
// CreateItemTranslateExchangeIdsResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemTranslateExchangeIdsResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemTranslateExchangeIdsResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemTranslateExchangeIdsPostResponseable instead.
type ItemTranslateExchangeIdsResponseable interface {
    ItemTranslateExchangeIdsPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

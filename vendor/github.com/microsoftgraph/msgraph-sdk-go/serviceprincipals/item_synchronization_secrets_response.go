package serviceprincipals

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSynchronizationSecretsPutResponseable instead.
type ItemSynchronizationSecretsResponse struct {
    ItemSynchronizationSecretsPutResponse
}
// NewItemSynchronizationSecretsResponse instantiates a new ItemSynchronizationSecretsResponse and sets the default values.
func NewItemSynchronizationSecretsResponse()(*ItemSynchronizationSecretsResponse) {
    m := &ItemSynchronizationSecretsResponse{
        ItemSynchronizationSecretsPutResponse: *NewItemSynchronizationSecretsPutResponse(),
    }
    return m
}
// CreateItemSynchronizationSecretsResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSynchronizationSecretsResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSynchronizationSecretsResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSynchronizationSecretsPutResponseable instead.
type ItemSynchronizationSecretsResponseable interface {
    ItemSynchronizationSecretsPutResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

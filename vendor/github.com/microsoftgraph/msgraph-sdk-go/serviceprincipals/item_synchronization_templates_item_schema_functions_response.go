package serviceprincipals

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSynchronizationTemplatesItemSchemaFunctionsGetResponseable instead.
type ItemSynchronizationTemplatesItemSchemaFunctionsResponse struct {
    ItemSynchronizationTemplatesItemSchemaFunctionsGetResponse
}
// NewItemSynchronizationTemplatesItemSchemaFunctionsResponse instantiates a new ItemSynchronizationTemplatesItemSchemaFunctionsResponse and sets the default values.
func NewItemSynchronizationTemplatesItemSchemaFunctionsResponse()(*ItemSynchronizationTemplatesItemSchemaFunctionsResponse) {
    m := &ItemSynchronizationTemplatesItemSchemaFunctionsResponse{
        ItemSynchronizationTemplatesItemSchemaFunctionsGetResponse: *NewItemSynchronizationTemplatesItemSchemaFunctionsGetResponse(),
    }
    return m
}
// CreateItemSynchronizationTemplatesItemSchemaFunctionsResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSynchronizationTemplatesItemSchemaFunctionsResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSynchronizationTemplatesItemSchemaFunctionsResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSynchronizationTemplatesItemSchemaFunctionsGetResponseable instead.
type ItemSynchronizationTemplatesItemSchemaFunctionsResponseable interface {
    ItemSynchronizationTemplatesItemSchemaFunctionsGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

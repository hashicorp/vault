package applications

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSynchronizationTemplatesItemSchemaFilterOperatorsGetResponseable instead.
type ItemSynchronizationTemplatesItemSchemaFilterOperatorsResponse struct {
    ItemSynchronizationTemplatesItemSchemaFilterOperatorsGetResponse
}
// NewItemSynchronizationTemplatesItemSchemaFilterOperatorsResponse instantiates a new ItemSynchronizationTemplatesItemSchemaFilterOperatorsResponse and sets the default values.
func NewItemSynchronizationTemplatesItemSchemaFilterOperatorsResponse()(*ItemSynchronizationTemplatesItemSchemaFilterOperatorsResponse) {
    m := &ItemSynchronizationTemplatesItemSchemaFilterOperatorsResponse{
        ItemSynchronizationTemplatesItemSchemaFilterOperatorsGetResponse: *NewItemSynchronizationTemplatesItemSchemaFilterOperatorsGetResponse(),
    }
    return m
}
// CreateItemSynchronizationTemplatesItemSchemaFilterOperatorsResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSynchronizationTemplatesItemSchemaFilterOperatorsResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSynchronizationTemplatesItemSchemaFilterOperatorsResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSynchronizationTemplatesItemSchemaFilterOperatorsGetResponseable instead.
type ItemSynchronizationTemplatesItemSchemaFilterOperatorsResponseable interface {
    ItemSynchronizationTemplatesItemSchemaFilterOperatorsGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

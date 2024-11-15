package serviceprincipals

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSynchronizationJobsItemSchemaFilterOperatorsGetResponseable instead.
type ItemSynchronizationJobsItemSchemaFilterOperatorsResponse struct {
    ItemSynchronizationJobsItemSchemaFilterOperatorsGetResponse
}
// NewItemSynchronizationJobsItemSchemaFilterOperatorsResponse instantiates a new ItemSynchronizationJobsItemSchemaFilterOperatorsResponse and sets the default values.
func NewItemSynchronizationJobsItemSchemaFilterOperatorsResponse()(*ItemSynchronizationJobsItemSchemaFilterOperatorsResponse) {
    m := &ItemSynchronizationJobsItemSchemaFilterOperatorsResponse{
        ItemSynchronizationJobsItemSchemaFilterOperatorsGetResponse: *NewItemSynchronizationJobsItemSchemaFilterOperatorsGetResponse(),
    }
    return m
}
// CreateItemSynchronizationJobsItemSchemaFilterOperatorsResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSynchronizationJobsItemSchemaFilterOperatorsResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSynchronizationJobsItemSchemaFilterOperatorsResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSynchronizationJobsItemSchemaFilterOperatorsGetResponseable instead.
type ItemSynchronizationJobsItemSchemaFilterOperatorsResponseable interface {
    ItemSynchronizationJobsItemSchemaFilterOperatorsGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

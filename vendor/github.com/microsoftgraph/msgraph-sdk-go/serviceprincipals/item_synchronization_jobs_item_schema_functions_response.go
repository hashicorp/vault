package serviceprincipals

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ItemSynchronizationJobsItemSchemaFunctionsGetResponseable instead.
type ItemSynchronizationJobsItemSchemaFunctionsResponse struct {
    ItemSynchronizationJobsItemSchemaFunctionsGetResponse
}
// NewItemSynchronizationJobsItemSchemaFunctionsResponse instantiates a new ItemSynchronizationJobsItemSchemaFunctionsResponse and sets the default values.
func NewItemSynchronizationJobsItemSchemaFunctionsResponse()(*ItemSynchronizationJobsItemSchemaFunctionsResponse) {
    m := &ItemSynchronizationJobsItemSchemaFunctionsResponse{
        ItemSynchronizationJobsItemSchemaFunctionsGetResponse: *NewItemSynchronizationJobsItemSchemaFunctionsGetResponse(),
    }
    return m
}
// CreateItemSynchronizationJobsItemSchemaFunctionsResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemSynchronizationJobsItemSchemaFunctionsResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSynchronizationJobsItemSchemaFunctionsResponse(), nil
}
// Deprecated: This class is obsolete. Use ItemSynchronizationJobsItemSchemaFunctionsGetResponseable instead.
type ItemSynchronizationJobsItemSchemaFunctionsResponseable interface {
    ItemSynchronizationJobsItemSchemaFunctionsGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

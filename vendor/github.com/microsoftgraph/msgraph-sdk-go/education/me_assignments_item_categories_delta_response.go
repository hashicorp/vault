package education

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use MeAssignmentsItemCategoriesDeltaGetResponseable instead.
type MeAssignmentsItemCategoriesDeltaResponse struct {
    MeAssignmentsItemCategoriesDeltaGetResponse
}
// NewMeAssignmentsItemCategoriesDeltaResponse instantiates a new MeAssignmentsItemCategoriesDeltaResponse and sets the default values.
func NewMeAssignmentsItemCategoriesDeltaResponse()(*MeAssignmentsItemCategoriesDeltaResponse) {
    m := &MeAssignmentsItemCategoriesDeltaResponse{
        MeAssignmentsItemCategoriesDeltaGetResponse: *NewMeAssignmentsItemCategoriesDeltaGetResponse(),
    }
    return m
}
// CreateMeAssignmentsItemCategoriesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMeAssignmentsItemCategoriesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMeAssignmentsItemCategoriesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use MeAssignmentsItemCategoriesDeltaGetResponseable instead.
type MeAssignmentsItemCategoriesDeltaResponseable interface {
    MeAssignmentsItemCategoriesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package education

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ClassesItemAssignmentsItemCategoriesDeltaGetResponseable instead.
type ClassesItemAssignmentsItemCategoriesDeltaResponse struct {
    ClassesItemAssignmentsItemCategoriesDeltaGetResponse
}
// NewClassesItemAssignmentsItemCategoriesDeltaResponse instantiates a new ClassesItemAssignmentsItemCategoriesDeltaResponse and sets the default values.
func NewClassesItemAssignmentsItemCategoriesDeltaResponse()(*ClassesItemAssignmentsItemCategoriesDeltaResponse) {
    m := &ClassesItemAssignmentsItemCategoriesDeltaResponse{
        ClassesItemAssignmentsItemCategoriesDeltaGetResponse: *NewClassesItemAssignmentsItemCategoriesDeltaGetResponse(),
    }
    return m
}
// CreateClassesItemAssignmentsItemCategoriesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateClassesItemAssignmentsItemCategoriesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewClassesItemAssignmentsItemCategoriesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ClassesItemAssignmentsItemCategoriesDeltaGetResponseable instead.
type ClassesItemAssignmentsItemCategoriesDeltaResponseable interface {
    ClassesItemAssignmentsItemCategoriesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

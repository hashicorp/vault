package education

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ClassesItemAssignmentCategoriesDeltaGetResponseable instead.
type ClassesItemAssignmentCategoriesDeltaResponse struct {
    ClassesItemAssignmentCategoriesDeltaGetResponse
}
// NewClassesItemAssignmentCategoriesDeltaResponse instantiates a new ClassesItemAssignmentCategoriesDeltaResponse and sets the default values.
func NewClassesItemAssignmentCategoriesDeltaResponse()(*ClassesItemAssignmentCategoriesDeltaResponse) {
    m := &ClassesItemAssignmentCategoriesDeltaResponse{
        ClassesItemAssignmentCategoriesDeltaGetResponse: *NewClassesItemAssignmentCategoriesDeltaGetResponse(),
    }
    return m
}
// CreateClassesItemAssignmentCategoriesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateClassesItemAssignmentCategoriesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewClassesItemAssignmentCategoriesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ClassesItemAssignmentCategoriesDeltaGetResponseable instead.
type ClassesItemAssignmentCategoriesDeltaResponseable interface {
    ClassesItemAssignmentCategoriesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

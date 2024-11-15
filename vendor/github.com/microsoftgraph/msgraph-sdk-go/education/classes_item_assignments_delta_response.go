package education

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ClassesItemAssignmentsDeltaGetResponseable instead.
type ClassesItemAssignmentsDeltaResponse struct {
    ClassesItemAssignmentsDeltaGetResponse
}
// NewClassesItemAssignmentsDeltaResponse instantiates a new ClassesItemAssignmentsDeltaResponse and sets the default values.
func NewClassesItemAssignmentsDeltaResponse()(*ClassesItemAssignmentsDeltaResponse) {
    m := &ClassesItemAssignmentsDeltaResponse{
        ClassesItemAssignmentsDeltaGetResponse: *NewClassesItemAssignmentsDeltaGetResponse(),
    }
    return m
}
// CreateClassesItemAssignmentsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateClassesItemAssignmentsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewClassesItemAssignmentsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ClassesItemAssignmentsDeltaGetResponseable instead.
type ClassesItemAssignmentsDeltaResponseable interface {
    ClassesItemAssignmentsDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

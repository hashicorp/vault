package education

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use ClassesDeltaGetResponseable instead.
type ClassesDeltaResponse struct {
    ClassesDeltaGetResponse
}
// NewClassesDeltaResponse instantiates a new ClassesDeltaResponse and sets the default values.
func NewClassesDeltaResponse()(*ClassesDeltaResponse) {
    m := &ClassesDeltaResponse{
        ClassesDeltaGetResponse: *NewClassesDeltaGetResponse(),
    }
    return m
}
// CreateClassesDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateClassesDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewClassesDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use ClassesDeltaGetResponseable instead.
type ClassesDeltaResponseable interface {
    ClassesDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

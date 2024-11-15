package education

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use MeAssignmentsDeltaGetResponseable instead.
type MeAssignmentsDeltaResponse struct {
    MeAssignmentsDeltaGetResponse
}
// NewMeAssignmentsDeltaResponse instantiates a new MeAssignmentsDeltaResponse and sets the default values.
func NewMeAssignmentsDeltaResponse()(*MeAssignmentsDeltaResponse) {
    m := &MeAssignmentsDeltaResponse{
        MeAssignmentsDeltaGetResponse: *NewMeAssignmentsDeltaGetResponse(),
    }
    return m
}
// CreateMeAssignmentsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMeAssignmentsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMeAssignmentsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use MeAssignmentsDeltaGetResponseable instead.
type MeAssignmentsDeltaResponseable interface {
    MeAssignmentsDeltaGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

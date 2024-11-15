package education

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use SchoolsDeltaGetResponseable instead.
type SchoolsDeltaResponse struct {
    SchoolsDeltaGetResponse
}
// NewSchoolsDeltaResponse instantiates a new SchoolsDeltaResponse and sets the default values.
func NewSchoolsDeltaResponse()(*SchoolsDeltaResponse) {
    m := &SchoolsDeltaResponse{
        SchoolsDeltaGetResponse: *NewSchoolsDeltaGetResponse(),
    }
    return m
}
// CreateSchoolsDeltaResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSchoolsDeltaResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSchoolsDeltaResponse(), nil
}
// Deprecated: This class is obsolete. Use SchoolsDeltaGetResponseable instead.
type SchoolsDeltaResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SchoolsDeltaGetResponseable
}

package reports

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use SecurityGetAttackSimulationTrainingUserCoverageGetResponseable instead.
type SecurityGetAttackSimulationTrainingUserCoverageResponse struct {
    SecurityGetAttackSimulationTrainingUserCoverageGetResponse
}
// NewSecurityGetAttackSimulationTrainingUserCoverageResponse instantiates a new SecurityGetAttackSimulationTrainingUserCoverageResponse and sets the default values.
func NewSecurityGetAttackSimulationTrainingUserCoverageResponse()(*SecurityGetAttackSimulationTrainingUserCoverageResponse) {
    m := &SecurityGetAttackSimulationTrainingUserCoverageResponse{
        SecurityGetAttackSimulationTrainingUserCoverageGetResponse: *NewSecurityGetAttackSimulationTrainingUserCoverageGetResponse(),
    }
    return m
}
// CreateSecurityGetAttackSimulationTrainingUserCoverageResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSecurityGetAttackSimulationTrainingUserCoverageResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSecurityGetAttackSimulationTrainingUserCoverageResponse(), nil
}
// Deprecated: This class is obsolete. Use SecurityGetAttackSimulationTrainingUserCoverageGetResponseable instead.
type SecurityGetAttackSimulationTrainingUserCoverageResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SecurityGetAttackSimulationTrainingUserCoverageGetResponseable
}

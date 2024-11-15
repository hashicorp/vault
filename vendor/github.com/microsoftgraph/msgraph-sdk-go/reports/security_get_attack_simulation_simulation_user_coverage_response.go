package reports

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use SecurityGetAttackSimulationSimulationUserCoverageGetResponseable instead.
type SecurityGetAttackSimulationSimulationUserCoverageResponse struct {
    SecurityGetAttackSimulationSimulationUserCoverageGetResponse
}
// NewSecurityGetAttackSimulationSimulationUserCoverageResponse instantiates a new SecurityGetAttackSimulationSimulationUserCoverageResponse and sets the default values.
func NewSecurityGetAttackSimulationSimulationUserCoverageResponse()(*SecurityGetAttackSimulationSimulationUserCoverageResponse) {
    m := &SecurityGetAttackSimulationSimulationUserCoverageResponse{
        SecurityGetAttackSimulationSimulationUserCoverageGetResponse: *NewSecurityGetAttackSimulationSimulationUserCoverageGetResponse(),
    }
    return m
}
// CreateSecurityGetAttackSimulationSimulationUserCoverageResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSecurityGetAttackSimulationSimulationUserCoverageResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSecurityGetAttackSimulationSimulationUserCoverageResponse(), nil
}
// Deprecated: This class is obsolete. Use SecurityGetAttackSimulationSimulationUserCoverageGetResponseable instead.
type SecurityGetAttackSimulationSimulationUserCoverageResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SecurityGetAttackSimulationSimulationUserCoverageGetResponseable
}

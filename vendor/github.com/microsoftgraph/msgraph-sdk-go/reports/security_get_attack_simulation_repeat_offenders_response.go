package reports

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use SecurityGetAttackSimulationRepeatOffendersGetResponseable instead.
type SecurityGetAttackSimulationRepeatOffendersResponse struct {
    SecurityGetAttackSimulationRepeatOffendersGetResponse
}
// NewSecurityGetAttackSimulationRepeatOffendersResponse instantiates a new SecurityGetAttackSimulationRepeatOffendersResponse and sets the default values.
func NewSecurityGetAttackSimulationRepeatOffendersResponse()(*SecurityGetAttackSimulationRepeatOffendersResponse) {
    m := &SecurityGetAttackSimulationRepeatOffendersResponse{
        SecurityGetAttackSimulationRepeatOffendersGetResponse: *NewSecurityGetAttackSimulationRepeatOffendersGetResponse(),
    }
    return m
}
// CreateSecurityGetAttackSimulationRepeatOffendersResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSecurityGetAttackSimulationRepeatOffendersResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSecurityGetAttackSimulationRepeatOffendersResponse(), nil
}
// Deprecated: This class is obsolete. Use SecurityGetAttackSimulationRepeatOffendersGetResponseable instead.
type SecurityGetAttackSimulationRepeatOffendersResponseable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    SecurityGetAttackSimulationRepeatOffendersGetResponseable
}

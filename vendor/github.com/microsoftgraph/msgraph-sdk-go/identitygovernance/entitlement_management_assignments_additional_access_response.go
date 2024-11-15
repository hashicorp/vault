package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use EntitlementManagementAssignmentsAdditionalAccessGetResponseable instead.
type EntitlementManagementAssignmentsAdditionalAccessResponse struct {
    EntitlementManagementAssignmentsAdditionalAccessGetResponse
}
// NewEntitlementManagementAssignmentsAdditionalAccessResponse instantiates a new EntitlementManagementAssignmentsAdditionalAccessResponse and sets the default values.
func NewEntitlementManagementAssignmentsAdditionalAccessResponse()(*EntitlementManagementAssignmentsAdditionalAccessResponse) {
    m := &EntitlementManagementAssignmentsAdditionalAccessResponse{
        EntitlementManagementAssignmentsAdditionalAccessGetResponse: *NewEntitlementManagementAssignmentsAdditionalAccessGetResponse(),
    }
    return m
}
// CreateEntitlementManagementAssignmentsAdditionalAccessResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEntitlementManagementAssignmentsAdditionalAccessResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEntitlementManagementAssignmentsAdditionalAccessResponse(), nil
}
// Deprecated: This class is obsolete. Use EntitlementManagementAssignmentsAdditionalAccessGetResponseable instead.
type EntitlementManagementAssignmentsAdditionalAccessResponseable interface {
    EntitlementManagementAssignmentsAdditionalAccessGetResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

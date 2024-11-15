package devicemanagement

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Deprecated: This class is obsolete. Use DeviceCompliancePoliciesItemAssignPostResponseable instead.
type DeviceCompliancePoliciesItemAssignResponse struct {
    DeviceCompliancePoliciesItemAssignPostResponse
}
// NewDeviceCompliancePoliciesItemAssignResponse instantiates a new DeviceCompliancePoliciesItemAssignResponse and sets the default values.
func NewDeviceCompliancePoliciesItemAssignResponse()(*DeviceCompliancePoliciesItemAssignResponse) {
    m := &DeviceCompliancePoliciesItemAssignResponse{
        DeviceCompliancePoliciesItemAssignPostResponse: *NewDeviceCompliancePoliciesItemAssignPostResponse(),
    }
    return m
}
// CreateDeviceCompliancePoliciesItemAssignResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceCompliancePoliciesItemAssignResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceCompliancePoliciesItemAssignResponse(), nil
}
// Deprecated: This class is obsolete. Use DeviceCompliancePoliciesItemAssignPostResponseable instead.
type DeviceCompliancePoliciesItemAssignResponseable interface {
    DeviceCompliancePoliciesItemAssignPostResponseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// AllLicensedUsersAssignmentTarget represents an assignment to all licensed users in the tenant.
type AllLicensedUsersAssignmentTarget struct {
    DeviceAndAppManagementAssignmentTarget
}
// NewAllLicensedUsersAssignmentTarget instantiates a new AllLicensedUsersAssignmentTarget and sets the default values.
func NewAllLicensedUsersAssignmentTarget()(*AllLicensedUsersAssignmentTarget) {
    m := &AllLicensedUsersAssignmentTarget{
        DeviceAndAppManagementAssignmentTarget: *NewDeviceAndAppManagementAssignmentTarget(),
    }
    odataTypeValue := "#microsoft.graph.allLicensedUsersAssignmentTarget"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAllLicensedUsersAssignmentTargetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAllLicensedUsersAssignmentTargetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAllLicensedUsersAssignmentTarget(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AllLicensedUsersAssignmentTarget) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceAndAppManagementAssignmentTarget.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *AllLicensedUsersAssignmentTarget) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceAndAppManagementAssignmentTarget.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type AllLicensedUsersAssignmentTargetable interface {
    DeviceAndAppManagementAssignmentTargetable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

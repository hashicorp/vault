package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceAndAppManagementRoleDefinition the Role Definition resource. The role definition is the foundation of role based access in Intune. The role combines an Intune resource such as a Mobile App and associated role permissions such as Create or Read for the resource. There are two types of roles, built-in and custom. Built-in roles cannot be modified. Both built-in roles and custom roles must have assignments to be enforced. Create custom roles if you want to define a role that allows any of the available resources and role permissions to be combined into a single role.
type DeviceAndAppManagementRoleDefinition struct {
    RoleDefinition
}
// NewDeviceAndAppManagementRoleDefinition instantiates a new DeviceAndAppManagementRoleDefinition and sets the default values.
func NewDeviceAndAppManagementRoleDefinition()(*DeviceAndAppManagementRoleDefinition) {
    m := &DeviceAndAppManagementRoleDefinition{
        RoleDefinition: *NewRoleDefinition(),
    }
    odataTypeValue := "#microsoft.graph.deviceAndAppManagementRoleDefinition"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateDeviceAndAppManagementRoleDefinitionFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceAndAppManagementRoleDefinitionFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceAndAppManagementRoleDefinition(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceAndAppManagementRoleDefinition) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.RoleDefinition.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *DeviceAndAppManagementRoleDefinition) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.RoleDefinition.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type DeviceAndAppManagementRoleDefinitionable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    RoleDefinitionable
}

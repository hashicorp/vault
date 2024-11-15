package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceAndAppManagementRoleAssignment the Role Assignment resource. Role assignments tie together a role definition with members and scopes. There can be one or more role assignments per role. This applies to custom and built-in roles.
type DeviceAndAppManagementRoleAssignment struct {
    RoleAssignment
}
// NewDeviceAndAppManagementRoleAssignment instantiates a new DeviceAndAppManagementRoleAssignment and sets the default values.
func NewDeviceAndAppManagementRoleAssignment()(*DeviceAndAppManagementRoleAssignment) {
    m := &DeviceAndAppManagementRoleAssignment{
        RoleAssignment: *NewRoleAssignment(),
    }
    return m
}
// CreateDeviceAndAppManagementRoleAssignmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceAndAppManagementRoleAssignmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceAndAppManagementRoleAssignment(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceAndAppManagementRoleAssignment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.RoleAssignment.GetFieldDeserializers()
    res["members"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetMembers(res)
        }
        return nil
    }
    return res
}
// GetMembers gets the members property value. The list of ids of role member security groups. These are IDs from Azure Active Directory.
// returns a []string when successful
func (m *DeviceAndAppManagementRoleAssignment) GetMembers()([]string) {
    val, err := m.GetBackingStore().Get("members")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceAndAppManagementRoleAssignment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.RoleAssignment.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetMembers() != nil {
        err = writer.WriteCollectionOfStringValues("members", m.GetMembers())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMembers sets the members property value. The list of ids of role member security groups. These are IDs from Azure Active Directory.
func (m *DeviceAndAppManagementRoleAssignment) SetMembers(value []string)() {
    err := m.GetBackingStore().Set("members", value)
    if err != nil {
        panic(err)
    }
}
type DeviceAndAppManagementRoleAssignmentable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    RoleAssignmentable
    GetMembers()([]string)
    SetMembers(value []string)()
}

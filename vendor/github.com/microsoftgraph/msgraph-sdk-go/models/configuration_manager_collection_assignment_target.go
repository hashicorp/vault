package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ConfigurationManagerCollectionAssignmentTarget represents an assignment to a Configuration Manager Collection.
type ConfigurationManagerCollectionAssignmentTarget struct {
    DeviceAndAppManagementAssignmentTarget
}
// NewConfigurationManagerCollectionAssignmentTarget instantiates a new ConfigurationManagerCollectionAssignmentTarget and sets the default values.
func NewConfigurationManagerCollectionAssignmentTarget()(*ConfigurationManagerCollectionAssignmentTarget) {
    m := &ConfigurationManagerCollectionAssignmentTarget{
        DeviceAndAppManagementAssignmentTarget: *NewDeviceAndAppManagementAssignmentTarget(),
    }
    odataTypeValue := "#microsoft.graph.configurationManagerCollectionAssignmentTarget"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateConfigurationManagerCollectionAssignmentTargetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConfigurationManagerCollectionAssignmentTargetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConfigurationManagerCollectionAssignmentTarget(), nil
}
// GetCollectionId gets the collectionId property value. The collection Id that is the target of the assignment.
// returns a *string when successful
func (m *ConfigurationManagerCollectionAssignmentTarget) GetCollectionId()(*string) {
    val, err := m.GetBackingStore().Get("collectionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConfigurationManagerCollectionAssignmentTarget) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceAndAppManagementAssignmentTarget.GetFieldDeserializers()
    res["collectionId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCollectionId(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *ConfigurationManagerCollectionAssignmentTarget) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceAndAppManagementAssignmentTarget.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("collectionId", m.GetCollectionId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCollectionId sets the collectionId property value. The collection Id that is the target of the assignment.
func (m *ConfigurationManagerCollectionAssignmentTarget) SetCollectionId(value *string)() {
    err := m.GetBackingStore().Set("collectionId", value)
    if err != nil {
        panic(err)
    }
}
type ConfigurationManagerCollectionAssignmentTargetable interface {
    DeviceAndAppManagementAssignmentTargetable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCollectionId()(*string)
    SetCollectionId(value *string)()
}

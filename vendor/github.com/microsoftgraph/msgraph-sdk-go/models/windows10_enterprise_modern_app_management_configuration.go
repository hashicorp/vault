package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// Windows10EnterpriseModernAppManagementConfiguration windows10 Enterprise Modern App Management Configuration.
type Windows10EnterpriseModernAppManagementConfiguration struct {
    DeviceConfiguration
}
// NewWindows10EnterpriseModernAppManagementConfiguration instantiates a new Windows10EnterpriseModernAppManagementConfiguration and sets the default values.
func NewWindows10EnterpriseModernAppManagementConfiguration()(*Windows10EnterpriseModernAppManagementConfiguration) {
    m := &Windows10EnterpriseModernAppManagementConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.windows10EnterpriseModernAppManagementConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindows10EnterpriseModernAppManagementConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindows10EnterpriseModernAppManagementConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindows10EnterpriseModernAppManagementConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Windows10EnterpriseModernAppManagementConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["uninstallBuiltInApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUninstallBuiltInApps(val)
        }
        return nil
    }
    return res
}
// GetUninstallBuiltInApps gets the uninstallBuiltInApps property value. Indicates whether or not to uninstall a fixed list of built-in Windows apps.
// returns a *bool when successful
func (m *Windows10EnterpriseModernAppManagementConfiguration) GetUninstallBuiltInApps()(*bool) {
    val, err := m.GetBackingStore().Get("uninstallBuiltInApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Windows10EnterpriseModernAppManagementConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("uninstallBuiltInApps", m.GetUninstallBuiltInApps())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetUninstallBuiltInApps sets the uninstallBuiltInApps property value. Indicates whether or not to uninstall a fixed list of built-in Windows apps.
func (m *Windows10EnterpriseModernAppManagementConfiguration) SetUninstallBuiltInApps(value *bool)() {
    err := m.GetBackingStore().Set("uninstallBuiltInApps", value)
    if err != nil {
        panic(err)
    }
}
type Windows10EnterpriseModernAppManagementConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetUninstallBuiltInApps()(*bool)
    SetUninstallBuiltInApps(value *bool)()
}

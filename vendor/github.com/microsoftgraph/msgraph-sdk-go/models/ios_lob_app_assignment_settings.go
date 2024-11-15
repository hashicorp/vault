package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosLobAppAssignmentSettings contains properties used to assign an iOS LOB mobile app to a group.
type IosLobAppAssignmentSettings struct {
    MobileAppAssignmentSettings
}
// NewIosLobAppAssignmentSettings instantiates a new IosLobAppAssignmentSettings and sets the default values.
func NewIosLobAppAssignmentSettings()(*IosLobAppAssignmentSettings) {
    m := &IosLobAppAssignmentSettings{
        MobileAppAssignmentSettings: *NewMobileAppAssignmentSettings(),
    }
    odataTypeValue := "#microsoft.graph.iosLobAppAssignmentSettings"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosLobAppAssignmentSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosLobAppAssignmentSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosLobAppAssignmentSettings(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IosLobAppAssignmentSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileAppAssignmentSettings.GetFieldDeserializers()
    res["isRemovable"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRemovable(val)
        }
        return nil
    }
    res["uninstallOnDeviceRemoval"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUninstallOnDeviceRemoval(val)
        }
        return nil
    }
    res["vpnConfigurationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVpnConfigurationId(val)
        }
        return nil
    }
    return res
}
// GetIsRemovable gets the isRemovable property value. When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to null which internally is treated as TRUE.
// returns a *bool when successful
func (m *IosLobAppAssignmentSettings) GetIsRemovable()(*bool) {
    val, err := m.GetBackingStore().Get("isRemovable")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUninstallOnDeviceRemoval gets the uninstallOnDeviceRemoval property value. When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, property is set to null which internally is treated as TRUE.
// returns a *bool when successful
func (m *IosLobAppAssignmentSettings) GetUninstallOnDeviceRemoval()(*bool) {
    val, err := m.GetBackingStore().Get("uninstallOnDeviceRemoval")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetVpnConfigurationId gets the vpnConfigurationId property value. This is the unique identifier (Id) of the VPN Configuration to apply to the app.
// returns a *string when successful
func (m *IosLobAppAssignmentSettings) GetVpnConfigurationId()(*string) {
    val, err := m.GetBackingStore().Get("vpnConfigurationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IosLobAppAssignmentSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileAppAssignmentSettings.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isRemovable", m.GetIsRemovable())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("uninstallOnDeviceRemoval", m.GetUninstallOnDeviceRemoval())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("vpnConfigurationId", m.GetVpnConfigurationId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsRemovable sets the isRemovable property value. When TRUE, indicates that the app can be uninstalled by the user. When FALSE, indicates that the app cannot be uninstalled by the user. By default, this property is set to null which internally is treated as TRUE.
func (m *IosLobAppAssignmentSettings) SetIsRemovable(value *bool)() {
    err := m.GetBackingStore().Set("isRemovable", value)
    if err != nil {
        panic(err)
    }
}
// SetUninstallOnDeviceRemoval sets the uninstallOnDeviceRemoval property value. When TRUE, indicates that the app should be uninstalled when the device is removed from Intune. When FALSE, indicates that the app will not be uninstalled when the device is removed from Intune. By default, property is set to null which internally is treated as TRUE.
func (m *IosLobAppAssignmentSettings) SetUninstallOnDeviceRemoval(value *bool)() {
    err := m.GetBackingStore().Set("uninstallOnDeviceRemoval", value)
    if err != nil {
        panic(err)
    }
}
// SetVpnConfigurationId sets the vpnConfigurationId property value. This is the unique identifier (Id) of the VPN Configuration to apply to the app.
func (m *IosLobAppAssignmentSettings) SetVpnConfigurationId(value *string)() {
    err := m.GetBackingStore().Set("vpnConfigurationId", value)
    if err != nil {
        panic(err)
    }
}
type IosLobAppAssignmentSettingsable interface {
    MobileAppAssignmentSettingsable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetIsRemovable()(*bool)
    GetUninstallOnDeviceRemoval()(*bool)
    GetVpnConfigurationId()(*string)
    SetIsRemovable(value *bool)()
    SetUninstallOnDeviceRemoval(value *bool)()
    SetVpnConfigurationId(value *string)()
}

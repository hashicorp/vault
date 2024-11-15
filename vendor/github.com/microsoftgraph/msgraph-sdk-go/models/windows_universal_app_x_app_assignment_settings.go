package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsUniversalAppXAppAssignmentSettings contains properties used when assigning a Windows Universal AppX mobile app to a group.
type WindowsUniversalAppXAppAssignmentSettings struct {
    MobileAppAssignmentSettings
}
// NewWindowsUniversalAppXAppAssignmentSettings instantiates a new WindowsUniversalAppXAppAssignmentSettings and sets the default values.
func NewWindowsUniversalAppXAppAssignmentSettings()(*WindowsUniversalAppXAppAssignmentSettings) {
    m := &WindowsUniversalAppXAppAssignmentSettings{
        MobileAppAssignmentSettings: *NewMobileAppAssignmentSettings(),
    }
    odataTypeValue := "#microsoft.graph.windowsUniversalAppXAppAssignmentSettings"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsUniversalAppXAppAssignmentSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsUniversalAppXAppAssignmentSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsUniversalAppXAppAssignmentSettings(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WindowsUniversalAppXAppAssignmentSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileAppAssignmentSettings.GetFieldDeserializers()
    res["useDeviceContext"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUseDeviceContext(val)
        }
        return nil
    }
    return res
}
// GetUseDeviceContext gets the useDeviceContext property value. If true, uses device execution context for Windows Universal AppX mobile app. Device-context install is not allowed when this type of app is targeted with Available intent. Defaults to false.
// returns a *bool when successful
func (m *WindowsUniversalAppXAppAssignmentSettings) GetUseDeviceContext()(*bool) {
    val, err := m.GetBackingStore().Get("useDeviceContext")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsUniversalAppXAppAssignmentSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileAppAssignmentSettings.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("useDeviceContext", m.GetUseDeviceContext())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetUseDeviceContext sets the useDeviceContext property value. If true, uses device execution context for Windows Universal AppX mobile app. Device-context install is not allowed when this type of app is targeted with Available intent. Defaults to false.
func (m *WindowsUniversalAppXAppAssignmentSettings) SetUseDeviceContext(value *bool)() {
    err := m.GetBackingStore().Set("useDeviceContext", value)
    if err != nil {
        panic(err)
    }
}
type WindowsUniversalAppXAppAssignmentSettingsable interface {
    MobileAppAssignmentSettingsable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetUseDeviceContext()(*bool)
    SetUseDeviceContext(value *bool)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsInformationProtectionDesktopApp desktop App for Windows information protection
type WindowsInformationProtectionDesktopApp struct {
    WindowsInformationProtectionApp
}
// NewWindowsInformationProtectionDesktopApp instantiates a new WindowsInformationProtectionDesktopApp and sets the default values.
func NewWindowsInformationProtectionDesktopApp()(*WindowsInformationProtectionDesktopApp) {
    m := &WindowsInformationProtectionDesktopApp{
        WindowsInformationProtectionApp: *NewWindowsInformationProtectionApp(),
    }
    odataTypeValue := "#microsoft.graph.windowsInformationProtectionDesktopApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsInformationProtectionDesktopAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsInformationProtectionDesktopAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsInformationProtectionDesktopApp(), nil
}
// GetBinaryName gets the binaryName property value. The binary name.
// returns a *string when successful
func (m *WindowsInformationProtectionDesktopApp) GetBinaryName()(*string) {
    val, err := m.GetBackingStore().Get("binaryName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBinaryVersionHigh gets the binaryVersionHigh property value. The high binary version.
// returns a *string when successful
func (m *WindowsInformationProtectionDesktopApp) GetBinaryVersionHigh()(*string) {
    val, err := m.GetBackingStore().Get("binaryVersionHigh")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBinaryVersionLow gets the binaryVersionLow property value. The lower binary version.
// returns a *string when successful
func (m *WindowsInformationProtectionDesktopApp) GetBinaryVersionLow()(*string) {
    val, err := m.GetBackingStore().Get("binaryVersionLow")
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
func (m *WindowsInformationProtectionDesktopApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WindowsInformationProtectionApp.GetFieldDeserializers()
    res["binaryName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBinaryName(val)
        }
        return nil
    }
    res["binaryVersionHigh"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBinaryVersionHigh(val)
        }
        return nil
    }
    res["binaryVersionLow"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBinaryVersionLow(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *WindowsInformationProtectionDesktopApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WindowsInformationProtectionApp.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("binaryName", m.GetBinaryName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("binaryVersionHigh", m.GetBinaryVersionHigh())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("binaryVersionLow", m.GetBinaryVersionLow())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBinaryName sets the binaryName property value. The binary name.
func (m *WindowsInformationProtectionDesktopApp) SetBinaryName(value *string)() {
    err := m.GetBackingStore().Set("binaryName", value)
    if err != nil {
        panic(err)
    }
}
// SetBinaryVersionHigh sets the binaryVersionHigh property value. The high binary version.
func (m *WindowsInformationProtectionDesktopApp) SetBinaryVersionHigh(value *string)() {
    err := m.GetBackingStore().Set("binaryVersionHigh", value)
    if err != nil {
        panic(err)
    }
}
// SetBinaryVersionLow sets the binaryVersionLow property value. The lower binary version.
func (m *WindowsInformationProtectionDesktopApp) SetBinaryVersionLow(value *string)() {
    err := m.GetBackingStore().Set("binaryVersionLow", value)
    if err != nil {
        panic(err)
    }
}
type WindowsInformationProtectionDesktopAppable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    WindowsInformationProtectionAppable
    GetBinaryName()(*string)
    GetBinaryVersionHigh()(*string)
    GetBinaryVersionLow()(*string)
    SetBinaryName(value *string)()
    SetBinaryVersionHigh(value *string)()
    SetBinaryVersionLow(value *string)()
}

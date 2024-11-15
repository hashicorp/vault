package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsUniversalAppXContainedApp a class that represents a contained app of a WindowsUniversalAppX app.
type WindowsUniversalAppXContainedApp struct {
    MobileContainedApp
}
// NewWindowsUniversalAppXContainedApp instantiates a new WindowsUniversalAppXContainedApp and sets the default values.
func NewWindowsUniversalAppXContainedApp()(*WindowsUniversalAppXContainedApp) {
    m := &WindowsUniversalAppXContainedApp{
        MobileContainedApp: *NewMobileContainedApp(),
    }
    odataTypeValue := "#microsoft.graph.windowsUniversalAppXContainedApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsUniversalAppXContainedAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsUniversalAppXContainedAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsUniversalAppXContainedApp(), nil
}
// GetAppUserModelId gets the appUserModelId property value. The app user model ID of the contained app of a WindowsUniversalAppX app.
// returns a *string when successful
func (m *WindowsUniversalAppXContainedApp) GetAppUserModelId()(*string) {
    val, err := m.GetBackingStore().Get("appUserModelId")
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
func (m *WindowsUniversalAppXContainedApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileContainedApp.GetFieldDeserializers()
    res["appUserModelId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppUserModelId(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *WindowsUniversalAppXContainedApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileContainedApp.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appUserModelId", m.GetAppUserModelId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppUserModelId sets the appUserModelId property value. The app user model ID of the contained app of a WindowsUniversalAppX app.
func (m *WindowsUniversalAppXContainedApp) SetAppUserModelId(value *string)() {
    err := m.GetBackingStore().Set("appUserModelId", value)
    if err != nil {
        panic(err)
    }
}
type WindowsUniversalAppXContainedAppable interface {
    MobileContainedAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppUserModelId()(*string)
    SetAppUserModelId(value *string)()
}

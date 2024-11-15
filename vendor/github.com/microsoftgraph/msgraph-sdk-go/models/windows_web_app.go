package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsWebApp contains properties and inherited properties for Windows web apps.
type WindowsWebApp struct {
    MobileApp
}
// NewWindowsWebApp instantiates a new WindowsWebApp and sets the default values.
func NewWindowsWebApp()(*WindowsWebApp) {
    m := &WindowsWebApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.windowsWebApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsWebAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsWebAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsWebApp(), nil
}
// GetAppUrl gets the appUrl property value. Indicates the Windows web app URL. Example: 'https://www.contoso.com'
// returns a *string when successful
func (m *WindowsWebApp) GetAppUrl()(*string) {
    val, err := m.GetBackingStore().Get("appUrl")
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
func (m *WindowsWebApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileApp.GetFieldDeserializers()
    res["appUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppUrl(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *WindowsWebApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileApp.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("appUrl", m.GetAppUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppUrl sets the appUrl property value. Indicates the Windows web app URL. Example: 'https://www.contoso.com'
func (m *WindowsWebApp) SetAppUrl(value *string)() {
    err := m.GetBackingStore().Set("appUrl", value)
    if err != nil {
        panic(err)
    }
}
type WindowsWebAppable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppUrl()(*string)
    SetAppUrl(value *string)()
}

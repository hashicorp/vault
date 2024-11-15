package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WebApp contains properties and inherited properties for web apps.
type WebApp struct {
    MobileApp
}
// NewWebApp instantiates a new WebApp and sets the default values.
func NewWebApp()(*WebApp) {
    m := &WebApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.webApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWebAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWebAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWebApp(), nil
}
// GetAppUrl gets the appUrl property value. The web app URL. This property cannot be PATCHed.
// returns a *string when successful
func (m *WebApp) GetAppUrl()(*string) {
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
func (m *WebApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["useManagedBrowser"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUseManagedBrowser(val)
        }
        return nil
    }
    return res
}
// GetUseManagedBrowser gets the useManagedBrowser property value. Whether or not to use managed browser. This property is only applicable for Android and IOS.
// returns a *bool when successful
func (m *WebApp) GetUseManagedBrowser()(*bool) {
    val, err := m.GetBackingStore().Get("useManagedBrowser")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WebApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    {
        err = writer.WriteBoolValue("useManagedBrowser", m.GetUseManagedBrowser())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAppUrl sets the appUrl property value. The web app URL. This property cannot be PATCHed.
func (m *WebApp) SetAppUrl(value *string)() {
    err := m.GetBackingStore().Set("appUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetUseManagedBrowser sets the useManagedBrowser property value. Whether or not to use managed browser. This property is only applicable for Android and IOS.
func (m *WebApp) SetUseManagedBrowser(value *bool)() {
    err := m.GetBackingStore().Set("useManagedBrowser", value)
    if err != nil {
        panic(err)
    }
}
type WebAppable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppUrl()(*string)
    GetUseManagedBrowser()(*bool)
    SetAppUrl(value *string)()
    SetUseManagedBrowser(value *bool)()
}

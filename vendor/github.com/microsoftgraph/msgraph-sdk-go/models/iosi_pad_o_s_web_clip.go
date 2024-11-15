package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosiPadOSWebClip contains properties and inherited properties for iOS web apps.
type IosiPadOSWebClip struct {
    MobileApp
}
// NewIosiPadOSWebClip instantiates a new IosiPadOSWebClip and sets the default values.
func NewIosiPadOSWebClip()(*IosiPadOSWebClip) {
    m := &IosiPadOSWebClip{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.iosiPadOSWebClip"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosiPadOSWebClipFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosiPadOSWebClipFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosiPadOSWebClip(), nil
}
// GetAppUrl gets the appUrl property value. Indicates iOS/iPadOS web clip app URL. Example: 'https://www.contoso.com'
// returns a *string when successful
func (m *IosiPadOSWebClip) GetAppUrl()(*string) {
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
func (m *IosiPadOSWebClip) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
// GetUseManagedBrowser gets the useManagedBrowser property value. Whether or not to use managed browser. When TRUE, the app will be required to be opened in Microsoft Edge. When FALSE, the app will not be required to be opened in Microsoft Edge. By default, this property is set to FALSE.
// returns a *bool when successful
func (m *IosiPadOSWebClip) GetUseManagedBrowser()(*bool) {
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
func (m *IosiPadOSWebClip) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
// SetAppUrl sets the appUrl property value. Indicates iOS/iPadOS web clip app URL. Example: 'https://www.contoso.com'
func (m *IosiPadOSWebClip) SetAppUrl(value *string)() {
    err := m.GetBackingStore().Set("appUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetUseManagedBrowser sets the useManagedBrowser property value. Whether or not to use managed browser. When TRUE, the app will be required to be opened in Microsoft Edge. When FALSE, the app will not be required to be opened in Microsoft Edge. By default, this property is set to FALSE.
func (m *IosiPadOSWebClip) SetUseManagedBrowser(value *bool)() {
    err := m.GetBackingStore().Set("useManagedBrowser", value)
    if err != nil {
        panic(err)
    }
}
type IosiPadOSWebClipable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppUrl()(*string)
    GetUseManagedBrowser()(*bool)
    SetAppUrl(value *string)()
    SetUseManagedBrowser(value *bool)()
}

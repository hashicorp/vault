package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsMicrosoftEdgeApp contains properties and inherited properties for the Microsoft Edge app on Windows.
type WindowsMicrosoftEdgeApp struct {
    MobileApp
}
// NewWindowsMicrosoftEdgeApp instantiates a new WindowsMicrosoftEdgeApp and sets the default values.
func NewWindowsMicrosoftEdgeApp()(*WindowsMicrosoftEdgeApp) {
    m := &WindowsMicrosoftEdgeApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.windowsMicrosoftEdgeApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsMicrosoftEdgeAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsMicrosoftEdgeAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsMicrosoftEdgeApp(), nil
}
// GetChannel gets the channel property value. The enum to specify the channels for Microsoft Edge apps.
// returns a *MicrosoftEdgeChannel when successful
func (m *WindowsMicrosoftEdgeApp) GetChannel()(*MicrosoftEdgeChannel) {
    val, err := m.GetBackingStore().Get("channel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MicrosoftEdgeChannel)
    }
    return nil
}
// GetDisplayLanguageLocale gets the displayLanguageLocale property value. The language locale to use when the Edge app displays text to the user.
// returns a *string when successful
func (m *WindowsMicrosoftEdgeApp) GetDisplayLanguageLocale()(*string) {
    val, err := m.GetBackingStore().Get("displayLanguageLocale")
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
func (m *WindowsMicrosoftEdgeApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileApp.GetFieldDeserializers()
    res["channel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMicrosoftEdgeChannel)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChannel(val.(*MicrosoftEdgeChannel))
        }
        return nil
    }
    res["displayLanguageLocale"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayLanguageLocale(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *WindowsMicrosoftEdgeApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileApp.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetChannel() != nil {
        cast := (*m.GetChannel()).String()
        err = writer.WriteStringValue("channel", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayLanguageLocale", m.GetDisplayLanguageLocale())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChannel sets the channel property value. The enum to specify the channels for Microsoft Edge apps.
func (m *WindowsMicrosoftEdgeApp) SetChannel(value *MicrosoftEdgeChannel)() {
    err := m.GetBackingStore().Set("channel", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayLanguageLocale sets the displayLanguageLocale property value. The language locale to use when the Edge app displays text to the user.
func (m *WindowsMicrosoftEdgeApp) SetDisplayLanguageLocale(value *string)() {
    err := m.GetBackingStore().Set("displayLanguageLocale", value)
    if err != nil {
        panic(err)
    }
}
type WindowsMicrosoftEdgeAppable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetChannel()(*MicrosoftEdgeChannel)
    GetDisplayLanguageLocale()(*string)
    SetChannel(value *MicrosoftEdgeChannel)()
    SetDisplayLanguageLocale(value *string)()
}

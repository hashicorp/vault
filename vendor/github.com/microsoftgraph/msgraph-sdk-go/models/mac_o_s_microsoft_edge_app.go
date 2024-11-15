package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MacOSMicrosoftEdgeApp contains properties and inherited properties for the macOS Microsoft Edge App.
type MacOSMicrosoftEdgeApp struct {
    MobileApp
}
// NewMacOSMicrosoftEdgeApp instantiates a new MacOSMicrosoftEdgeApp and sets the default values.
func NewMacOSMicrosoftEdgeApp()(*MacOSMicrosoftEdgeApp) {
    m := &MacOSMicrosoftEdgeApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.macOSMicrosoftEdgeApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMacOSMicrosoftEdgeAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMacOSMicrosoftEdgeAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMacOSMicrosoftEdgeApp(), nil
}
// GetChannel gets the channel property value. The enum to specify the channels for Microsoft Edge apps.
// returns a *MicrosoftEdgeChannel when successful
func (m *MacOSMicrosoftEdgeApp) GetChannel()(*MicrosoftEdgeChannel) {
    val, err := m.GetBackingStore().Get("channel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MicrosoftEdgeChannel)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MacOSMicrosoftEdgeApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    return res
}
// Serialize serializes information the current object
func (m *MacOSMicrosoftEdgeApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
    return nil
}
// SetChannel sets the channel property value. The enum to specify the channels for Microsoft Edge apps.
func (m *MacOSMicrosoftEdgeApp) SetChannel(value *MicrosoftEdgeChannel)() {
    err := m.GetBackingStore().Set("channel", value)
    if err != nil {
        panic(err)
    }
}
type MacOSMicrosoftEdgeAppable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetChannel()(*MicrosoftEdgeChannel)
    SetChannel(value *MicrosoftEdgeChannel)()
}

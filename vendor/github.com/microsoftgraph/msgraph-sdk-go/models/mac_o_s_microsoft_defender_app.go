package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MacOSMicrosoftDefenderApp contains properties and inherited properties for the macOS Microsoft Defender App.
type MacOSMicrosoftDefenderApp struct {
    MobileApp
}
// NewMacOSMicrosoftDefenderApp instantiates a new MacOSMicrosoftDefenderApp and sets the default values.
func NewMacOSMicrosoftDefenderApp()(*MacOSMicrosoftDefenderApp) {
    m := &MacOSMicrosoftDefenderApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.macOSMicrosoftDefenderApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMacOSMicrosoftDefenderAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMacOSMicrosoftDefenderAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMacOSMicrosoftDefenderApp(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MacOSMicrosoftDefenderApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileApp.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *MacOSMicrosoftDefenderApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileApp.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type MacOSMicrosoftDefenderAppable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MacOSOfficeSuiteApp contains properties and inherited properties for the MacOS Office Suite App.
type MacOSOfficeSuiteApp struct {
    MobileApp
}
// NewMacOSOfficeSuiteApp instantiates a new MacOSOfficeSuiteApp and sets the default values.
func NewMacOSOfficeSuiteApp()(*MacOSOfficeSuiteApp) {
    m := &MacOSOfficeSuiteApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.macOSOfficeSuiteApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMacOSOfficeSuiteAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMacOSOfficeSuiteAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMacOSOfficeSuiteApp(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MacOSOfficeSuiteApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileApp.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *MacOSOfficeSuiteApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileApp.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
type MacOSOfficeSuiteAppable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
}

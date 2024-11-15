package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosHomeScreenApp represents an icon for an app on the Home Screen
type IosHomeScreenApp struct {
    IosHomeScreenItem
}
// NewIosHomeScreenApp instantiates a new IosHomeScreenApp and sets the default values.
func NewIosHomeScreenApp()(*IosHomeScreenApp) {
    m := &IosHomeScreenApp{
        IosHomeScreenItem: *NewIosHomeScreenItem(),
    }
    odataTypeValue := "#microsoft.graph.iosHomeScreenApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosHomeScreenAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosHomeScreenAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosHomeScreenApp(), nil
}
// GetBundleID gets the bundleID property value. BundleID of the app if isWebClip is false or the URL of a web clip if isWebClip is true.
// returns a *string when successful
func (m *IosHomeScreenApp) GetBundleID()(*string) {
    val, err := m.GetBackingStore().Get("bundleID")
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
func (m *IosHomeScreenApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.IosHomeScreenItem.GetFieldDeserializers()
    res["bundleID"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBundleID(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *IosHomeScreenApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.IosHomeScreenItem.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("bundleID", m.GetBundleID())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBundleID sets the bundleID property value. BundleID of the app if isWebClip is false or the URL of a web clip if isWebClip is true.
func (m *IosHomeScreenApp) SetBundleID(value *string)() {
    err := m.GetBackingStore().Set("bundleID", value)
    if err != nil {
        panic(err)
    }
}
type IosHomeScreenAppable interface {
    IosHomeScreenItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBundleID()(*string)
    SetBundleID(value *string)()
}

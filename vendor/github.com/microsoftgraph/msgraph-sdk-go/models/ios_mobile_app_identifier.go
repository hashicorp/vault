package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosMobileAppIdentifier the identifier for an iOS app.
type IosMobileAppIdentifier struct {
    MobileAppIdentifier
}
// NewIosMobileAppIdentifier instantiates a new IosMobileAppIdentifier and sets the default values.
func NewIosMobileAppIdentifier()(*IosMobileAppIdentifier) {
    m := &IosMobileAppIdentifier{
        MobileAppIdentifier: *NewMobileAppIdentifier(),
    }
    odataTypeValue := "#microsoft.graph.iosMobileAppIdentifier"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosMobileAppIdentifierFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosMobileAppIdentifierFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosMobileAppIdentifier(), nil
}
// GetBundleId gets the bundleId property value. The identifier for an app, as specified in the app store.
// returns a *string when successful
func (m *IosMobileAppIdentifier) GetBundleId()(*string) {
    val, err := m.GetBackingStore().Get("bundleId")
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
func (m *IosMobileAppIdentifier) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileAppIdentifier.GetFieldDeserializers()
    res["bundleId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBundleId(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *IosMobileAppIdentifier) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileAppIdentifier.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("bundleId", m.GetBundleId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBundleId sets the bundleId property value. The identifier for an app, as specified in the app store.
func (m *IosMobileAppIdentifier) SetBundleId(value *string)() {
    err := m.GetBackingStore().Set("bundleId", value)
    if err != nil {
        panic(err)
    }
}
type IosMobileAppIdentifierable interface {
    MobileAppIdentifierable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBundleId()(*string)
    SetBundleId(value *string)()
}

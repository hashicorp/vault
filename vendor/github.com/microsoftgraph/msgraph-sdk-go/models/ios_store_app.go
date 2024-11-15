package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosStoreApp contains properties and inherited properties for iOS store apps.
type IosStoreApp struct {
    MobileApp
}
// NewIosStoreApp instantiates a new IosStoreApp and sets the default values.
func NewIosStoreApp()(*IosStoreApp) {
    m := &IosStoreApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.iosStoreApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosStoreAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosStoreAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosStoreApp(), nil
}
// GetApplicableDeviceType gets the applicableDeviceType property value. Contains properties of the possible iOS device types the mobile app can run on.
// returns a IosDeviceTypeable when successful
func (m *IosStoreApp) GetApplicableDeviceType()(IosDeviceTypeable) {
    val, err := m.GetBackingStore().Get("applicableDeviceType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IosDeviceTypeable)
    }
    return nil
}
// GetAppStoreUrl gets the appStoreUrl property value. The Apple App Store URL
// returns a *string when successful
func (m *IosStoreApp) GetAppStoreUrl()(*string) {
    val, err := m.GetBackingStore().Get("appStoreUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBundleId gets the bundleId property value. The Identity Name.
// returns a *string when successful
func (m *IosStoreApp) GetBundleId()(*string) {
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
func (m *IosStoreApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileApp.GetFieldDeserializers()
    res["applicableDeviceType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIosDeviceTypeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicableDeviceType(val.(IosDeviceTypeable))
        }
        return nil
    }
    res["appStoreUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppStoreUrl(val)
        }
        return nil
    }
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
    res["minimumSupportedOperatingSystem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIosMinimumOperatingSystemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumSupportedOperatingSystem(val.(IosMinimumOperatingSystemable))
        }
        return nil
    }
    return res
}
// GetMinimumSupportedOperatingSystem gets the minimumSupportedOperatingSystem property value. The value for the minimum applicable operating system.
// returns a IosMinimumOperatingSystemable when successful
func (m *IosStoreApp) GetMinimumSupportedOperatingSystem()(IosMinimumOperatingSystemable) {
    val, err := m.GetBackingStore().Get("minimumSupportedOperatingSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IosMinimumOperatingSystemable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IosStoreApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileApp.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("applicableDeviceType", m.GetApplicableDeviceType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("appStoreUrl", m.GetAppStoreUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("bundleId", m.GetBundleId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("minimumSupportedOperatingSystem", m.GetMinimumSupportedOperatingSystem())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicableDeviceType sets the applicableDeviceType property value. Contains properties of the possible iOS device types the mobile app can run on.
func (m *IosStoreApp) SetApplicableDeviceType(value IosDeviceTypeable)() {
    err := m.GetBackingStore().Set("applicableDeviceType", value)
    if err != nil {
        panic(err)
    }
}
// SetAppStoreUrl sets the appStoreUrl property value. The Apple App Store URL
func (m *IosStoreApp) SetAppStoreUrl(value *string)() {
    err := m.GetBackingStore().Set("appStoreUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetBundleId sets the bundleId property value. The Identity Name.
func (m *IosStoreApp) SetBundleId(value *string)() {
    err := m.GetBackingStore().Set("bundleId", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumSupportedOperatingSystem sets the minimumSupportedOperatingSystem property value. The value for the minimum applicable operating system.
func (m *IosStoreApp) SetMinimumSupportedOperatingSystem(value IosMinimumOperatingSystemable)() {
    err := m.GetBackingStore().Set("minimumSupportedOperatingSystem", value)
    if err != nil {
        panic(err)
    }
}
type IosStoreAppable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicableDeviceType()(IosDeviceTypeable)
    GetAppStoreUrl()(*string)
    GetBundleId()(*string)
    GetMinimumSupportedOperatingSystem()(IosMinimumOperatingSystemable)
    SetApplicableDeviceType(value IosDeviceTypeable)()
    SetAppStoreUrl(value *string)()
    SetBundleId(value *string)()
    SetMinimumSupportedOperatingSystem(value IosMinimumOperatingSystemable)()
}

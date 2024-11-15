package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// IosLobApp contains properties and inherited properties for iOS Line Of Business apps.
type IosLobApp struct {
    MobileLobApp
}
// NewIosLobApp instantiates a new IosLobApp and sets the default values.
func NewIosLobApp()(*IosLobApp) {
    m := &IosLobApp{
        MobileLobApp: *NewMobileLobApp(),
    }
    odataTypeValue := "#microsoft.graph.iosLobApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateIosLobAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateIosLobAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewIosLobApp(), nil
}
// GetApplicableDeviceType gets the applicableDeviceType property value. Contains properties of the possible iOS device types the mobile app can run on.
// returns a IosDeviceTypeable when successful
func (m *IosLobApp) GetApplicableDeviceType()(IosDeviceTypeable) {
    val, err := m.GetBackingStore().Get("applicableDeviceType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IosDeviceTypeable)
    }
    return nil
}
// GetBuildNumber gets the buildNumber property value. The build number of iOS Line of Business (LoB) app.
// returns a *string when successful
func (m *IosLobApp) GetBuildNumber()(*string) {
    val, err := m.GetBackingStore().Get("buildNumber")
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
func (m *IosLobApp) GetBundleId()(*string) {
    val, err := m.GetBackingStore().Get("bundleId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExpirationDateTime gets the expirationDateTime property value. The expiration time.
// returns a *Time when successful
func (m *IosLobApp) GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("expirationDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *IosLobApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileLobApp.GetFieldDeserializers()
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
    res["buildNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBuildNumber(val)
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
    res["expirationDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExpirationDateTime(val)
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
    res["versionNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersionNumber(val)
        }
        return nil
    }
    return res
}
// GetMinimumSupportedOperatingSystem gets the minimumSupportedOperatingSystem property value. The value for the minimum applicable operating system.
// returns a IosMinimumOperatingSystemable when successful
func (m *IosLobApp) GetMinimumSupportedOperatingSystem()(IosMinimumOperatingSystemable) {
    val, err := m.GetBackingStore().Get("minimumSupportedOperatingSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IosMinimumOperatingSystemable)
    }
    return nil
}
// GetVersionNumber gets the versionNumber property value. The version number of iOS Line of Business (LoB) app.
// returns a *string when successful
func (m *IosLobApp) GetVersionNumber()(*string) {
    val, err := m.GetBackingStore().Get("versionNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *IosLobApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileLobApp.Serialize(writer)
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
        err = writer.WriteStringValue("buildNumber", m.GetBuildNumber())
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
        err = writer.WriteTimeValue("expirationDateTime", m.GetExpirationDateTime())
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
    {
        err = writer.WriteStringValue("versionNumber", m.GetVersionNumber())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplicableDeviceType sets the applicableDeviceType property value. Contains properties of the possible iOS device types the mobile app can run on.
func (m *IosLobApp) SetApplicableDeviceType(value IosDeviceTypeable)() {
    err := m.GetBackingStore().Set("applicableDeviceType", value)
    if err != nil {
        panic(err)
    }
}
// SetBuildNumber sets the buildNumber property value. The build number of iOS Line of Business (LoB) app.
func (m *IosLobApp) SetBuildNumber(value *string)() {
    err := m.GetBackingStore().Set("buildNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetBundleId sets the bundleId property value. The Identity Name.
func (m *IosLobApp) SetBundleId(value *string)() {
    err := m.GetBackingStore().Set("bundleId", value)
    if err != nil {
        panic(err)
    }
}
// SetExpirationDateTime sets the expirationDateTime property value. The expiration time.
func (m *IosLobApp) SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("expirationDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumSupportedOperatingSystem sets the minimumSupportedOperatingSystem property value. The value for the minimum applicable operating system.
func (m *IosLobApp) SetMinimumSupportedOperatingSystem(value IosMinimumOperatingSystemable)() {
    err := m.GetBackingStore().Set("minimumSupportedOperatingSystem", value)
    if err != nil {
        panic(err)
    }
}
// SetVersionNumber sets the versionNumber property value. The version number of iOS Line of Business (LoB) app.
func (m *IosLobApp) SetVersionNumber(value *string)() {
    err := m.GetBackingStore().Set("versionNumber", value)
    if err != nil {
        panic(err)
    }
}
type IosLobAppable interface {
    MobileLobAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicableDeviceType()(IosDeviceTypeable)
    GetBuildNumber()(*string)
    GetBundleId()(*string)
    GetExpirationDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMinimumSupportedOperatingSystem()(IosMinimumOperatingSystemable)
    GetVersionNumber()(*string)
    SetApplicableDeviceType(value IosDeviceTypeable)()
    SetBuildNumber(value *string)()
    SetBundleId(value *string)()
    SetExpirationDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMinimumSupportedOperatingSystem(value IosMinimumOperatingSystemable)()
    SetVersionNumber(value *string)()
}

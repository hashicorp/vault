package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsUniversalAppX contains properties and inherited properties for Windows Universal AppX Line Of Business apps. Inherits from `mobileLobApp`.
type WindowsUniversalAppX struct {
    MobileLobApp
}
// NewWindowsUniversalAppX instantiates a new WindowsUniversalAppX and sets the default values.
func NewWindowsUniversalAppX()(*WindowsUniversalAppX) {
    m := &WindowsUniversalAppX{
        MobileLobApp: *NewMobileLobApp(),
    }
    odataTypeValue := "#microsoft.graph.windowsUniversalAppX"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateWindowsUniversalAppXFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWindowsUniversalAppXFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsUniversalAppX(), nil
}
// GetApplicableArchitectures gets the applicableArchitectures property value. Contains properties for Windows architecture.
// returns a *WindowsArchitecture when successful
func (m *WindowsUniversalAppX) GetApplicableArchitectures()(*WindowsArchitecture) {
    val, err := m.GetBackingStore().Get("applicableArchitectures")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsArchitecture)
    }
    return nil
}
// GetApplicableDeviceTypes gets the applicableDeviceTypes property value. Contains properties for Windows device type. Multiple values can be selected. Default value is `none`.
// returns a *WindowsDeviceType when successful
func (m *WindowsUniversalAppX) GetApplicableDeviceTypes()(*WindowsDeviceType) {
    val, err := m.GetBackingStore().Get("applicableDeviceTypes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*WindowsDeviceType)
    }
    return nil
}
// GetCommittedContainedApps gets the committedContainedApps property value. The collection of contained apps in the committed mobileAppContent of a windowsUniversalAppX app.
// returns a []MobileContainedAppable when successful
func (m *WindowsUniversalAppX) GetCommittedContainedApps()([]MobileContainedAppable) {
    val, err := m.GetBackingStore().Get("committedContainedApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MobileContainedAppable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WindowsUniversalAppX) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileLobApp.GetFieldDeserializers()
    res["applicableArchitectures"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsArchitecture)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicableArchitectures(val.(*WindowsArchitecture))
        }
        return nil
    }
    res["applicableDeviceTypes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindowsDeviceType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplicableDeviceTypes(val.(*WindowsDeviceType))
        }
        return nil
    }
    res["committedContainedApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMobileContainedAppFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MobileContainedAppable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MobileContainedAppable)
                }
            }
            m.SetCommittedContainedApps(res)
        }
        return nil
    }
    res["identityName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdentityName(val)
        }
        return nil
    }
    res["identityPublisherHash"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdentityPublisherHash(val)
        }
        return nil
    }
    res["identityResourceIdentifier"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdentityResourceIdentifier(val)
        }
        return nil
    }
    res["identityVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIdentityVersion(val)
        }
        return nil
    }
    res["isBundle"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsBundle(val)
        }
        return nil
    }
    res["minimumSupportedOperatingSystem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWindowsMinimumOperatingSystemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumSupportedOperatingSystem(val.(WindowsMinimumOperatingSystemable))
        }
        return nil
    }
    return res
}
// GetIdentityName gets the identityName property value. The Identity Name.
// returns a *string when successful
func (m *WindowsUniversalAppX) GetIdentityName()(*string) {
    val, err := m.GetBackingStore().Get("identityName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIdentityPublisherHash gets the identityPublisherHash property value. The Identity Publisher Hash.
// returns a *string when successful
func (m *WindowsUniversalAppX) GetIdentityPublisherHash()(*string) {
    val, err := m.GetBackingStore().Get("identityPublisherHash")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIdentityResourceIdentifier gets the identityResourceIdentifier property value. The Identity Resource Identifier.
// returns a *string when successful
func (m *WindowsUniversalAppX) GetIdentityResourceIdentifier()(*string) {
    val, err := m.GetBackingStore().Get("identityResourceIdentifier")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIdentityVersion gets the identityVersion property value. The identity version.
// returns a *string when successful
func (m *WindowsUniversalAppX) GetIdentityVersion()(*string) {
    val, err := m.GetBackingStore().Get("identityVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsBundle gets the isBundle property value. Whether or not the app is a bundle.
// returns a *bool when successful
func (m *WindowsUniversalAppX) GetIsBundle()(*bool) {
    val, err := m.GetBackingStore().Get("isBundle")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMinimumSupportedOperatingSystem gets the minimumSupportedOperatingSystem property value. The minimum operating system required for a Windows mobile app.
// returns a WindowsMinimumOperatingSystemable when successful
func (m *WindowsUniversalAppX) GetMinimumSupportedOperatingSystem()(WindowsMinimumOperatingSystemable) {
    val, err := m.GetBackingStore().Get("minimumSupportedOperatingSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WindowsMinimumOperatingSystemable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WindowsUniversalAppX) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileLobApp.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetApplicableArchitectures() != nil {
        cast := (*m.GetApplicableArchitectures()).String()
        err = writer.WriteStringValue("applicableArchitectures", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetApplicableDeviceTypes() != nil {
        cast := (*m.GetApplicableDeviceTypes()).String()
        err = writer.WriteStringValue("applicableDeviceTypes", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetCommittedContainedApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCommittedContainedApps()))
        for i, v := range m.GetCommittedContainedApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("committedContainedApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("identityName", m.GetIdentityName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("identityPublisherHash", m.GetIdentityPublisherHash())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("identityResourceIdentifier", m.GetIdentityResourceIdentifier())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("identityVersion", m.GetIdentityVersion())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isBundle", m.GetIsBundle())
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
// SetApplicableArchitectures sets the applicableArchitectures property value. Contains properties for Windows architecture.
func (m *WindowsUniversalAppX) SetApplicableArchitectures(value *WindowsArchitecture)() {
    err := m.GetBackingStore().Set("applicableArchitectures", value)
    if err != nil {
        panic(err)
    }
}
// SetApplicableDeviceTypes sets the applicableDeviceTypes property value. Contains properties for Windows device type. Multiple values can be selected. Default value is `none`.
func (m *WindowsUniversalAppX) SetApplicableDeviceTypes(value *WindowsDeviceType)() {
    err := m.GetBackingStore().Set("applicableDeviceTypes", value)
    if err != nil {
        panic(err)
    }
}
// SetCommittedContainedApps sets the committedContainedApps property value. The collection of contained apps in the committed mobileAppContent of a windowsUniversalAppX app.
func (m *WindowsUniversalAppX) SetCommittedContainedApps(value []MobileContainedAppable)() {
    err := m.GetBackingStore().Set("committedContainedApps", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentityName sets the identityName property value. The Identity Name.
func (m *WindowsUniversalAppX) SetIdentityName(value *string)() {
    err := m.GetBackingStore().Set("identityName", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentityPublisherHash sets the identityPublisherHash property value. The Identity Publisher Hash.
func (m *WindowsUniversalAppX) SetIdentityPublisherHash(value *string)() {
    err := m.GetBackingStore().Set("identityPublisherHash", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentityResourceIdentifier sets the identityResourceIdentifier property value. The Identity Resource Identifier.
func (m *WindowsUniversalAppX) SetIdentityResourceIdentifier(value *string)() {
    err := m.GetBackingStore().Set("identityResourceIdentifier", value)
    if err != nil {
        panic(err)
    }
}
// SetIdentityVersion sets the identityVersion property value. The identity version.
func (m *WindowsUniversalAppX) SetIdentityVersion(value *string)() {
    err := m.GetBackingStore().Set("identityVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetIsBundle sets the isBundle property value. Whether or not the app is a bundle.
func (m *WindowsUniversalAppX) SetIsBundle(value *bool)() {
    err := m.GetBackingStore().Set("isBundle", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumSupportedOperatingSystem sets the minimumSupportedOperatingSystem property value. The minimum operating system required for a Windows mobile app.
func (m *WindowsUniversalAppX) SetMinimumSupportedOperatingSystem(value WindowsMinimumOperatingSystemable)() {
    err := m.GetBackingStore().Set("minimumSupportedOperatingSystem", value)
    if err != nil {
        panic(err)
    }
}
type WindowsUniversalAppXable interface {
    MobileLobAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplicableArchitectures()(*WindowsArchitecture)
    GetApplicableDeviceTypes()(*WindowsDeviceType)
    GetCommittedContainedApps()([]MobileContainedAppable)
    GetIdentityName()(*string)
    GetIdentityPublisherHash()(*string)
    GetIdentityResourceIdentifier()(*string)
    GetIdentityVersion()(*string)
    GetIsBundle()(*bool)
    GetMinimumSupportedOperatingSystem()(WindowsMinimumOperatingSystemable)
    SetApplicableArchitectures(value *WindowsArchitecture)()
    SetApplicableDeviceTypes(value *WindowsDeviceType)()
    SetCommittedContainedApps(value []MobileContainedAppable)()
    SetIdentityName(value *string)()
    SetIdentityPublisherHash(value *string)()
    SetIdentityResourceIdentifier(value *string)()
    SetIdentityVersion(value *string)()
    SetIsBundle(value *bool)()
    SetMinimumSupportedOperatingSystem(value WindowsMinimumOperatingSystemable)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MacOSLobApp contains properties and inherited properties for the macOS LOB App.
type MacOSLobApp struct {
    MobileLobApp
}
// NewMacOSLobApp instantiates a new MacOSLobApp and sets the default values.
func NewMacOSLobApp()(*MacOSLobApp) {
    m := &MacOSLobApp{
        MobileLobApp: *NewMobileLobApp(),
    }
    odataTypeValue := "#microsoft.graph.macOSLobApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMacOSLobAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMacOSLobAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMacOSLobApp(), nil
}
// GetBuildNumber gets the buildNumber property value. The build number of the package. This should match the package CFBundleShortVersionString of the .pkg file.
// returns a *string when successful
func (m *MacOSLobApp) GetBuildNumber()(*string) {
    val, err := m.GetBackingStore().Get("buildNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBundleId gets the bundleId property value. The primary bundleId of the package.
// returns a *string when successful
func (m *MacOSLobApp) GetBundleId()(*string) {
    val, err := m.GetBackingStore().Get("bundleId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetChildApps gets the childApps property value. List of ComplexType macOSLobChildApp objects. Represents the apps expected to be installed by the package.
// returns a []MacOSLobChildAppable when successful
func (m *MacOSLobApp) GetChildApps()([]MacOSLobChildAppable) {
    val, err := m.GetBackingStore().Get("childApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MacOSLobChildAppable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MacOSLobApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileLobApp.GetFieldDeserializers()
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
    res["childApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMacOSLobChildAppFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MacOSLobChildAppable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MacOSLobChildAppable)
                }
            }
            m.SetChildApps(res)
        }
        return nil
    }
    res["ignoreVersionDetection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIgnoreVersionDetection(val)
        }
        return nil
    }
    res["installAsManaged"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInstallAsManaged(val)
        }
        return nil
    }
    res["md5Hash"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetMd5Hash(res)
        }
        return nil
    }
    res["md5HashChunkSize"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMd5HashChunkSize(val)
        }
        return nil
    }
    res["minimumSupportedOperatingSystem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMacOSMinimumOperatingSystemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumSupportedOperatingSystem(val.(MacOSMinimumOperatingSystemable))
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
// GetIgnoreVersionDetection gets the ignoreVersionDetection property value. When TRUE, indicates that the app's version will NOT be used to detect if the app is installed on a device. When FALSE, indicates that the app's version will be used to detect if the app is installed on a device. Set this to true for apps that use a self update feature.
// returns a *bool when successful
func (m *MacOSLobApp) GetIgnoreVersionDetection()(*bool) {
    val, err := m.GetBackingStore().Get("ignoreVersionDetection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetInstallAsManaged gets the installAsManaged property value. When TRUE, indicates that the app will be installed as managed (requires macOS 11.0 and other managed package restrictions). When FALSE, indicates that the app will be installed as unmanaged.
// returns a *bool when successful
func (m *MacOSLobApp) GetInstallAsManaged()(*bool) {
    val, err := m.GetBackingStore().Get("installAsManaged")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMd5Hash gets the md5Hash property value. The MD5 hash codes. This is empty if the package was uploaded directly. If the Intune App Wrapping Tool is used to create a .intunemac, this value can be found inside the Detection.xml file.
// returns a []string when successful
func (m *MacOSLobApp) GetMd5Hash()([]string) {
    val, err := m.GetBackingStore().Get("md5Hash")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetMd5HashChunkSize gets the md5HashChunkSize property value. The chunk size for MD5 hash. This is '0' or empty if the package was uploaded directly. If the Intune App Wrapping Tool is used to create a .intunemac, this value can be found inside the Detection.xml file.
// returns a *int32 when successful
func (m *MacOSLobApp) GetMd5HashChunkSize()(*int32) {
    val, err := m.GetBackingStore().Get("md5HashChunkSize")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetMinimumSupportedOperatingSystem gets the minimumSupportedOperatingSystem property value. ComplexType macOSMinimumOperatingSystem that indicates the minimum operating system applicable for the application.
// returns a MacOSMinimumOperatingSystemable when successful
func (m *MacOSLobApp) GetMinimumSupportedOperatingSystem()(MacOSMinimumOperatingSystemable) {
    val, err := m.GetBackingStore().Get("minimumSupportedOperatingSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MacOSMinimumOperatingSystemable)
    }
    return nil
}
// GetVersionNumber gets the versionNumber property value. The version number of the package. This should match the package CFBundleVersion in the packageinfo file.
// returns a *string when successful
func (m *MacOSLobApp) GetVersionNumber()(*string) {
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
func (m *MacOSLobApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileLobApp.Serialize(writer)
    if err != nil {
        return err
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
    if m.GetChildApps() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetChildApps()))
        for i, v := range m.GetChildApps() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("childApps", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("ignoreVersionDetection", m.GetIgnoreVersionDetection())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("installAsManaged", m.GetInstallAsManaged())
        if err != nil {
            return err
        }
    }
    if m.GetMd5Hash() != nil {
        err = writer.WriteCollectionOfStringValues("md5Hash", m.GetMd5Hash())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("md5HashChunkSize", m.GetMd5HashChunkSize())
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
// SetBuildNumber sets the buildNumber property value. The build number of the package. This should match the package CFBundleShortVersionString of the .pkg file.
func (m *MacOSLobApp) SetBuildNumber(value *string)() {
    err := m.GetBackingStore().Set("buildNumber", value)
    if err != nil {
        panic(err)
    }
}
// SetBundleId sets the bundleId property value. The primary bundleId of the package.
func (m *MacOSLobApp) SetBundleId(value *string)() {
    err := m.GetBackingStore().Set("bundleId", value)
    if err != nil {
        panic(err)
    }
}
// SetChildApps sets the childApps property value. List of ComplexType macOSLobChildApp objects. Represents the apps expected to be installed by the package.
func (m *MacOSLobApp) SetChildApps(value []MacOSLobChildAppable)() {
    err := m.GetBackingStore().Set("childApps", value)
    if err != nil {
        panic(err)
    }
}
// SetIgnoreVersionDetection sets the ignoreVersionDetection property value. When TRUE, indicates that the app's version will NOT be used to detect if the app is installed on a device. When FALSE, indicates that the app's version will be used to detect if the app is installed on a device. Set this to true for apps that use a self update feature.
func (m *MacOSLobApp) SetIgnoreVersionDetection(value *bool)() {
    err := m.GetBackingStore().Set("ignoreVersionDetection", value)
    if err != nil {
        panic(err)
    }
}
// SetInstallAsManaged sets the installAsManaged property value. When TRUE, indicates that the app will be installed as managed (requires macOS 11.0 and other managed package restrictions). When FALSE, indicates that the app will be installed as unmanaged.
func (m *MacOSLobApp) SetInstallAsManaged(value *bool)() {
    err := m.GetBackingStore().Set("installAsManaged", value)
    if err != nil {
        panic(err)
    }
}
// SetMd5Hash sets the md5Hash property value. The MD5 hash codes. This is empty if the package was uploaded directly. If the Intune App Wrapping Tool is used to create a .intunemac, this value can be found inside the Detection.xml file.
func (m *MacOSLobApp) SetMd5Hash(value []string)() {
    err := m.GetBackingStore().Set("md5Hash", value)
    if err != nil {
        panic(err)
    }
}
// SetMd5HashChunkSize sets the md5HashChunkSize property value. The chunk size for MD5 hash. This is '0' or empty if the package was uploaded directly. If the Intune App Wrapping Tool is used to create a .intunemac, this value can be found inside the Detection.xml file.
func (m *MacOSLobApp) SetMd5HashChunkSize(value *int32)() {
    err := m.GetBackingStore().Set("md5HashChunkSize", value)
    if err != nil {
        panic(err)
    }
}
// SetMinimumSupportedOperatingSystem sets the minimumSupportedOperatingSystem property value. ComplexType macOSMinimumOperatingSystem that indicates the minimum operating system applicable for the application.
func (m *MacOSLobApp) SetMinimumSupportedOperatingSystem(value MacOSMinimumOperatingSystemable)() {
    err := m.GetBackingStore().Set("minimumSupportedOperatingSystem", value)
    if err != nil {
        panic(err)
    }
}
// SetVersionNumber sets the versionNumber property value. The version number of the package. This should match the package CFBundleVersion in the packageinfo file.
func (m *MacOSLobApp) SetVersionNumber(value *string)() {
    err := m.GetBackingStore().Set("versionNumber", value)
    if err != nil {
        panic(err)
    }
}
type MacOSLobAppable interface {
    MobileLobAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBuildNumber()(*string)
    GetBundleId()(*string)
    GetChildApps()([]MacOSLobChildAppable)
    GetIgnoreVersionDetection()(*bool)
    GetInstallAsManaged()(*bool)
    GetMd5Hash()([]string)
    GetMd5HashChunkSize()(*int32)
    GetMinimumSupportedOperatingSystem()(MacOSMinimumOperatingSystemable)
    GetVersionNumber()(*string)
    SetBuildNumber(value *string)()
    SetBundleId(value *string)()
    SetChildApps(value []MacOSLobChildAppable)()
    SetIgnoreVersionDetection(value *bool)()
    SetInstallAsManaged(value *bool)()
    SetMd5Hash(value []string)()
    SetMd5HashChunkSize(value *int32)()
    SetMinimumSupportedOperatingSystem(value MacOSMinimumOperatingSystemable)()
    SetVersionNumber(value *string)()
}

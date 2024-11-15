package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ManagedAndroidLobApp contains properties and inherited properties for Managed Android Line Of Business apps.
type ManagedAndroidLobApp struct {
    ManagedMobileLobApp
}
// NewManagedAndroidLobApp instantiates a new ManagedAndroidLobApp and sets the default values.
func NewManagedAndroidLobApp()(*ManagedAndroidLobApp) {
    m := &ManagedAndroidLobApp{
        ManagedMobileLobApp: *NewManagedMobileLobApp(),
    }
    odataTypeValue := "#microsoft.graph.managedAndroidLobApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateManagedAndroidLobAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateManagedAndroidLobAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewManagedAndroidLobApp(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ManagedAndroidLobApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ManagedMobileLobApp.GetFieldDeserializers()
    res["minimumSupportedOperatingSystem"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAndroidMinimumOperatingSystemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMinimumSupportedOperatingSystem(val.(AndroidMinimumOperatingSystemable))
        }
        return nil
    }
    res["packageId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPackageId(val)
        }
        return nil
    }
    res["versionCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersionCode(val)
        }
        return nil
    }
    res["versionName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVersionName(val)
        }
        return nil
    }
    return res
}
// GetMinimumSupportedOperatingSystem gets the minimumSupportedOperatingSystem property value. The value for the minimum applicable operating system.
// returns a AndroidMinimumOperatingSystemable when successful
func (m *ManagedAndroidLobApp) GetMinimumSupportedOperatingSystem()(AndroidMinimumOperatingSystemable) {
    val, err := m.GetBackingStore().Get("minimumSupportedOperatingSystem")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AndroidMinimumOperatingSystemable)
    }
    return nil
}
// GetPackageId gets the packageId property value. The package identifier.
// returns a *string when successful
func (m *ManagedAndroidLobApp) GetPackageId()(*string) {
    val, err := m.GetBackingStore().Get("packageId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVersionCode gets the versionCode property value. The version code of managed Android Line of Business (LoB) app.
// returns a *string when successful
func (m *ManagedAndroidLobApp) GetVersionCode()(*string) {
    val, err := m.GetBackingStore().Get("versionCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVersionName gets the versionName property value. The version name of managed Android Line of Business (LoB) app.
// returns a *string when successful
func (m *ManagedAndroidLobApp) GetVersionName()(*string) {
    val, err := m.GetBackingStore().Get("versionName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ManagedAndroidLobApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ManagedMobileLobApp.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("minimumSupportedOperatingSystem", m.GetMinimumSupportedOperatingSystem())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("packageId", m.GetPackageId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("versionCode", m.GetVersionCode())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("versionName", m.GetVersionName())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMinimumSupportedOperatingSystem sets the minimumSupportedOperatingSystem property value. The value for the minimum applicable operating system.
func (m *ManagedAndroidLobApp) SetMinimumSupportedOperatingSystem(value AndroidMinimumOperatingSystemable)() {
    err := m.GetBackingStore().Set("minimumSupportedOperatingSystem", value)
    if err != nil {
        panic(err)
    }
}
// SetPackageId sets the packageId property value. The package identifier.
func (m *ManagedAndroidLobApp) SetPackageId(value *string)() {
    err := m.GetBackingStore().Set("packageId", value)
    if err != nil {
        panic(err)
    }
}
// SetVersionCode sets the versionCode property value. The version code of managed Android Line of Business (LoB) app.
func (m *ManagedAndroidLobApp) SetVersionCode(value *string)() {
    err := m.GetBackingStore().Set("versionCode", value)
    if err != nil {
        panic(err)
    }
}
// SetVersionName sets the versionName property value. The version name of managed Android Line of Business (LoB) app.
func (m *ManagedAndroidLobApp) SetVersionName(value *string)() {
    err := m.GetBackingStore().Set("versionName", value)
    if err != nil {
        panic(err)
    }
}
type ManagedAndroidLobAppable interface {
    ManagedMobileLobAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetMinimumSupportedOperatingSystem()(AndroidMinimumOperatingSystemable)
    GetPackageId()(*string)
    GetVersionCode()(*string)
    GetVersionName()(*string)
    SetMinimumSupportedOperatingSystem(value AndroidMinimumOperatingSystemable)()
    SetPackageId(value *string)()
    SetVersionCode(value *string)()
    SetVersionName(value *string)()
}

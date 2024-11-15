package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// EditionUpgradeConfiguration windows 10 Edition Upgrade configuration.
type EditionUpgradeConfiguration struct {
    DeviceConfiguration
}
// NewEditionUpgradeConfiguration instantiates a new EditionUpgradeConfiguration and sets the default values.
func NewEditionUpgradeConfiguration()(*EditionUpgradeConfiguration) {
    m := &EditionUpgradeConfiguration{
        DeviceConfiguration: *NewDeviceConfiguration(),
    }
    odataTypeValue := "#microsoft.graph.editionUpgradeConfiguration"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEditionUpgradeConfigurationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEditionUpgradeConfigurationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEditionUpgradeConfiguration(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EditionUpgradeConfiguration) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DeviceConfiguration.GetFieldDeserializers()
    res["license"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLicense(val)
        }
        return nil
    }
    res["licenseType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEditionUpgradeLicenseType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLicenseType(val.(*EditionUpgradeLicenseType))
        }
        return nil
    }
    res["productKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductKey(val)
        }
        return nil
    }
    res["targetEdition"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWindows10EditionType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTargetEdition(val.(*Windows10EditionType))
        }
        return nil
    }
    return res
}
// GetLicense gets the license property value. Edition Upgrade License File Content.
// returns a *string when successful
func (m *EditionUpgradeConfiguration) GetLicense()(*string) {
    val, err := m.GetBackingStore().Get("license")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetLicenseType gets the licenseType property value. Edition Upgrade License type
// returns a *EditionUpgradeLicenseType when successful
func (m *EditionUpgradeConfiguration) GetLicenseType()(*EditionUpgradeLicenseType) {
    val, err := m.GetBackingStore().Get("licenseType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EditionUpgradeLicenseType)
    }
    return nil
}
// GetProductKey gets the productKey property value. Edition Upgrade Product Key.
// returns a *string when successful
func (m *EditionUpgradeConfiguration) GetProductKey()(*string) {
    val, err := m.GetBackingStore().Get("productKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTargetEdition gets the targetEdition property value. Windows 10 Edition type.
// returns a *Windows10EditionType when successful
func (m *EditionUpgradeConfiguration) GetTargetEdition()(*Windows10EditionType) {
    val, err := m.GetBackingStore().Get("targetEdition")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Windows10EditionType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EditionUpgradeConfiguration) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DeviceConfiguration.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("license", m.GetLicense())
        if err != nil {
            return err
        }
    }
    if m.GetLicenseType() != nil {
        cast := (*m.GetLicenseType()).String()
        err = writer.WriteStringValue("licenseType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("productKey", m.GetProductKey())
        if err != nil {
            return err
        }
    }
    if m.GetTargetEdition() != nil {
        cast := (*m.GetTargetEdition()).String()
        err = writer.WriteStringValue("targetEdition", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLicense sets the license property value. Edition Upgrade License File Content.
func (m *EditionUpgradeConfiguration) SetLicense(value *string)() {
    err := m.GetBackingStore().Set("license", value)
    if err != nil {
        panic(err)
    }
}
// SetLicenseType sets the licenseType property value. Edition Upgrade License type
func (m *EditionUpgradeConfiguration) SetLicenseType(value *EditionUpgradeLicenseType)() {
    err := m.GetBackingStore().Set("licenseType", value)
    if err != nil {
        panic(err)
    }
}
// SetProductKey sets the productKey property value. Edition Upgrade Product Key.
func (m *EditionUpgradeConfiguration) SetProductKey(value *string)() {
    err := m.GetBackingStore().Set("productKey", value)
    if err != nil {
        panic(err)
    }
}
// SetTargetEdition sets the targetEdition property value. Windows 10 Edition type.
func (m *EditionUpgradeConfiguration) SetTargetEdition(value *Windows10EditionType)() {
    err := m.GetBackingStore().Set("targetEdition", value)
    if err != nil {
        panic(err)
    }
}
type EditionUpgradeConfigurationable interface {
    DeviceConfigurationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLicense()(*string)
    GetLicenseType()(*EditionUpgradeLicenseType)
    GetProductKey()(*string)
    GetTargetEdition()(*Windows10EditionType)
    SetLicense(value *string)()
    SetLicenseType(value *EditionUpgradeLicenseType)()
    SetProductKey(value *string)()
    SetTargetEdition(value *Windows10EditionType)()
}

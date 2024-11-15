package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MicrosoftStoreForBusinessApp microsoft Store for Business Apps. This class does not support Create, Delete, or Update.
type MicrosoftStoreForBusinessApp struct {
    MobileApp
}
// NewMicrosoftStoreForBusinessApp instantiates a new MicrosoftStoreForBusinessApp and sets the default values.
func NewMicrosoftStoreForBusinessApp()(*MicrosoftStoreForBusinessApp) {
    m := &MicrosoftStoreForBusinessApp{
        MobileApp: *NewMobileApp(),
    }
    odataTypeValue := "#microsoft.graph.microsoftStoreForBusinessApp"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMicrosoftStoreForBusinessAppFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMicrosoftStoreForBusinessAppFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMicrosoftStoreForBusinessApp(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MicrosoftStoreForBusinessApp) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MobileApp.GetFieldDeserializers()
    res["licenseType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMicrosoftStoreForBusinessLicenseType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLicenseType(val.(*MicrosoftStoreForBusinessLicenseType))
        }
        return nil
    }
    res["packageIdentityName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPackageIdentityName(val)
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
    res["totalLicenseCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalLicenseCount(val)
        }
        return nil
    }
    res["usedLicenseCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUsedLicenseCount(val)
        }
        return nil
    }
    return res
}
// GetLicenseType gets the licenseType property value. The licenseType property
// returns a *MicrosoftStoreForBusinessLicenseType when successful
func (m *MicrosoftStoreForBusinessApp) GetLicenseType()(*MicrosoftStoreForBusinessLicenseType) {
    val, err := m.GetBackingStore().Get("licenseType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MicrosoftStoreForBusinessLicenseType)
    }
    return nil
}
// GetPackageIdentityName gets the packageIdentityName property value. The app package identifier
// returns a *string when successful
func (m *MicrosoftStoreForBusinessApp) GetPackageIdentityName()(*string) {
    val, err := m.GetBackingStore().Get("packageIdentityName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProductKey gets the productKey property value. The app product key
// returns a *string when successful
func (m *MicrosoftStoreForBusinessApp) GetProductKey()(*string) {
    val, err := m.GetBackingStore().Get("productKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTotalLicenseCount gets the totalLicenseCount property value. The total number of Microsoft Store for Business licenses.
// returns a *int32 when successful
func (m *MicrosoftStoreForBusinessApp) GetTotalLicenseCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalLicenseCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUsedLicenseCount gets the usedLicenseCount property value. The number of Microsoft Store for Business licenses in use.
// returns a *int32 when successful
func (m *MicrosoftStoreForBusinessApp) GetUsedLicenseCount()(*int32) {
    val, err := m.GetBackingStore().Get("usedLicenseCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MicrosoftStoreForBusinessApp) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MobileApp.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetLicenseType() != nil {
        cast := (*m.GetLicenseType()).String()
        err = writer.WriteStringValue("licenseType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("packageIdentityName", m.GetPackageIdentityName())
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
    {
        err = writer.WriteInt32Value("totalLicenseCount", m.GetTotalLicenseCount())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("usedLicenseCount", m.GetUsedLicenseCount())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetLicenseType sets the licenseType property value. The licenseType property
func (m *MicrosoftStoreForBusinessApp) SetLicenseType(value *MicrosoftStoreForBusinessLicenseType)() {
    err := m.GetBackingStore().Set("licenseType", value)
    if err != nil {
        panic(err)
    }
}
// SetPackageIdentityName sets the packageIdentityName property value. The app package identifier
func (m *MicrosoftStoreForBusinessApp) SetPackageIdentityName(value *string)() {
    err := m.GetBackingStore().Set("packageIdentityName", value)
    if err != nil {
        panic(err)
    }
}
// SetProductKey sets the productKey property value. The app product key
func (m *MicrosoftStoreForBusinessApp) SetProductKey(value *string)() {
    err := m.GetBackingStore().Set("productKey", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalLicenseCount sets the totalLicenseCount property value. The total number of Microsoft Store for Business licenses.
func (m *MicrosoftStoreForBusinessApp) SetTotalLicenseCount(value *int32)() {
    err := m.GetBackingStore().Set("totalLicenseCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUsedLicenseCount sets the usedLicenseCount property value. The number of Microsoft Store for Business licenses in use.
func (m *MicrosoftStoreForBusinessApp) SetUsedLicenseCount(value *int32)() {
    err := m.GetBackingStore().Set("usedLicenseCount", value)
    if err != nil {
        panic(err)
    }
}
type MicrosoftStoreForBusinessAppable interface {
    MobileAppable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetLicenseType()(*MicrosoftStoreForBusinessLicenseType)
    GetPackageIdentityName()(*string)
    GetProductKey()(*string)
    GetTotalLicenseCount()(*int32)
    GetUsedLicenseCount()(*int32)
    SetLicenseType(value *MicrosoftStoreForBusinessLicenseType)()
    SetPackageIdentityName(value *string)()
    SetProductKey(value *string)()
    SetTotalLicenseCount(value *int32)()
    SetUsedLicenseCount(value *int32)()
}

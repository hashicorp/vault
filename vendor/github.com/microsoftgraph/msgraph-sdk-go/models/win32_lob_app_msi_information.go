package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// Win32LobAppMsiInformation contains MSI app properties for a Win32 App.
type Win32LobAppMsiInformation struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewWin32LobAppMsiInformation instantiates a new Win32LobAppMsiInformation and sets the default values.
func NewWin32LobAppMsiInformation()(*Win32LobAppMsiInformation) {
    m := &Win32LobAppMsiInformation{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateWin32LobAppMsiInformationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWin32LobAppMsiInformationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWin32LobAppMsiInformation(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Win32LobAppMsiInformation) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *Win32LobAppMsiInformation) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Win32LobAppMsiInformation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["packageType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseWin32LobAppMsiPackageType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPackageType(val.(*Win32LobAppMsiPackageType))
        }
        return nil
    }
    res["productCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductCode(val)
        }
        return nil
    }
    res["productName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductName(val)
        }
        return nil
    }
    res["productVersion"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProductVersion(val)
        }
        return nil
    }
    res["publisher"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPublisher(val)
        }
        return nil
    }
    res["requiresReboot"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequiresReboot(val)
        }
        return nil
    }
    res["upgradeCode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUpgradeCode(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Win32LobAppMsiInformation) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPackageType gets the packageType property value. Indicates the package type of an MSI Win32LobApp.
// returns a *Win32LobAppMsiPackageType when successful
func (m *Win32LobAppMsiInformation) GetPackageType()(*Win32LobAppMsiPackageType) {
    val, err := m.GetBackingStore().Get("packageType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Win32LobAppMsiPackageType)
    }
    return nil
}
// GetProductCode gets the productCode property value. The MSI product code.
// returns a *string when successful
func (m *Win32LobAppMsiInformation) GetProductCode()(*string) {
    val, err := m.GetBackingStore().Get("productCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProductName gets the productName property value. The MSI product name.
// returns a *string when successful
func (m *Win32LobAppMsiInformation) GetProductName()(*string) {
    val, err := m.GetBackingStore().Get("productName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProductVersion gets the productVersion property value. The MSI product version.
// returns a *string when successful
func (m *Win32LobAppMsiInformation) GetProductVersion()(*string) {
    val, err := m.GetBackingStore().Get("productVersion")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPublisher gets the publisher property value. The MSI publisher.
// returns a *string when successful
func (m *Win32LobAppMsiInformation) GetPublisher()(*string) {
    val, err := m.GetBackingStore().Get("publisher")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRequiresReboot gets the requiresReboot property value. Whether the MSI app requires the machine to reboot to complete installation.
// returns a *bool when successful
func (m *Win32LobAppMsiInformation) GetRequiresReboot()(*bool) {
    val, err := m.GetBackingStore().Get("requiresReboot")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetUpgradeCode gets the upgradeCode property value. The MSI upgrade code.
// returns a *string when successful
func (m *Win32LobAppMsiInformation) GetUpgradeCode()(*string) {
    val, err := m.GetBackingStore().Get("upgradeCode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Win32LobAppMsiInformation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetPackageType() != nil {
        cast := (*m.GetPackageType()).String()
        err := writer.WriteStringValue("packageType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("productCode", m.GetProductCode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("productName", m.GetProductName())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("productVersion", m.GetProductVersion())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("publisher", m.GetPublisher())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("requiresReboot", m.GetRequiresReboot())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("upgradeCode", m.GetUpgradeCode())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Win32LobAppMsiInformation) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Win32LobAppMsiInformation) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Win32LobAppMsiInformation) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPackageType sets the packageType property value. Indicates the package type of an MSI Win32LobApp.
func (m *Win32LobAppMsiInformation) SetPackageType(value *Win32LobAppMsiPackageType)() {
    err := m.GetBackingStore().Set("packageType", value)
    if err != nil {
        panic(err)
    }
}
// SetProductCode sets the productCode property value. The MSI product code.
func (m *Win32LobAppMsiInformation) SetProductCode(value *string)() {
    err := m.GetBackingStore().Set("productCode", value)
    if err != nil {
        panic(err)
    }
}
// SetProductName sets the productName property value. The MSI product name.
func (m *Win32LobAppMsiInformation) SetProductName(value *string)() {
    err := m.GetBackingStore().Set("productName", value)
    if err != nil {
        panic(err)
    }
}
// SetProductVersion sets the productVersion property value. The MSI product version.
func (m *Win32LobAppMsiInformation) SetProductVersion(value *string)() {
    err := m.GetBackingStore().Set("productVersion", value)
    if err != nil {
        panic(err)
    }
}
// SetPublisher sets the publisher property value. The MSI publisher.
func (m *Win32LobAppMsiInformation) SetPublisher(value *string)() {
    err := m.GetBackingStore().Set("publisher", value)
    if err != nil {
        panic(err)
    }
}
// SetRequiresReboot sets the requiresReboot property value. Whether the MSI app requires the machine to reboot to complete installation.
func (m *Win32LobAppMsiInformation) SetRequiresReboot(value *bool)() {
    err := m.GetBackingStore().Set("requiresReboot", value)
    if err != nil {
        panic(err)
    }
}
// SetUpgradeCode sets the upgradeCode property value. The MSI upgrade code.
func (m *Win32LobAppMsiInformation) SetUpgradeCode(value *string)() {
    err := m.GetBackingStore().Set("upgradeCode", value)
    if err != nil {
        panic(err)
    }
}
type Win32LobAppMsiInformationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    GetPackageType()(*Win32LobAppMsiPackageType)
    GetProductCode()(*string)
    GetProductName()(*string)
    GetProductVersion()(*string)
    GetPublisher()(*string)
    GetRequiresReboot()(*bool)
    GetUpgradeCode()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
    SetPackageType(value *Win32LobAppMsiPackageType)()
    SetProductCode(value *string)()
    SetProductName(value *string)()
    SetProductVersion(value *string)()
    SetPublisher(value *string)()
    SetRequiresReboot(value *bool)()
    SetUpgradeCode(value *string)()
}

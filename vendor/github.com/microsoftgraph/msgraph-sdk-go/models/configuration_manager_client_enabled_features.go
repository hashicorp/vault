package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// ConfigurationManagerClientEnabledFeatures configuration Manager client enabled features
type ConfigurationManagerClientEnabledFeatures struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewConfigurationManagerClientEnabledFeatures instantiates a new ConfigurationManagerClientEnabledFeatures and sets the default values.
func NewConfigurationManagerClientEnabledFeatures()(*ConfigurationManagerClientEnabledFeatures) {
    m := &ConfigurationManagerClientEnabledFeatures{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateConfigurationManagerClientEnabledFeaturesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConfigurationManagerClientEnabledFeaturesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConfigurationManagerClientEnabledFeatures(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ConfigurationManagerClientEnabledFeatures) GetAdditionalData()(map[string]any) {
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
func (m *ConfigurationManagerClientEnabledFeatures) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCompliancePolicy gets the compliancePolicy property value. Whether compliance policy is managed by Intune
// returns a *bool when successful
func (m *ConfigurationManagerClientEnabledFeatures) GetCompliancePolicy()(*bool) {
    val, err := m.GetBackingStore().Get("compliancePolicy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetDeviceConfiguration gets the deviceConfiguration property value. Whether device configuration is managed by Intune
// returns a *bool when successful
func (m *ConfigurationManagerClientEnabledFeatures) GetDeviceConfiguration()(*bool) {
    val, err := m.GetBackingStore().Get("deviceConfiguration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConfigurationManagerClientEnabledFeatures) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["compliancePolicy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCompliancePolicy(val)
        }
        return nil
    }
    res["deviceConfiguration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeviceConfiguration(val)
        }
        return nil
    }
    res["inventory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInventory(val)
        }
        return nil
    }
    res["modernApps"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetModernApps(val)
        }
        return nil
    }
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
    res["resourceAccess"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResourceAccess(val)
        }
        return nil
    }
    res["windowsUpdateForBusiness"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWindowsUpdateForBusiness(val)
        }
        return nil
    }
    return res
}
// GetInventory gets the inventory property value. Whether inventory is managed by Intune
// returns a *bool when successful
func (m *ConfigurationManagerClientEnabledFeatures) GetInventory()(*bool) {
    val, err := m.GetBackingStore().Get("inventory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetModernApps gets the modernApps property value. Whether modern application is managed by Intune
// returns a *bool when successful
func (m *ConfigurationManagerClientEnabledFeatures) GetModernApps()(*bool) {
    val, err := m.GetBackingStore().Get("modernApps")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ConfigurationManagerClientEnabledFeatures) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetResourceAccess gets the resourceAccess property value. Whether resource access is managed by Intune
// returns a *bool when successful
func (m *ConfigurationManagerClientEnabledFeatures) GetResourceAccess()(*bool) {
    val, err := m.GetBackingStore().Get("resourceAccess")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetWindowsUpdateForBusiness gets the windowsUpdateForBusiness property value. Whether Windows Update for Business is managed by Intune
// returns a *bool when successful
func (m *ConfigurationManagerClientEnabledFeatures) GetWindowsUpdateForBusiness()(*bool) {
    val, err := m.GetBackingStore().Get("windowsUpdateForBusiness")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ConfigurationManagerClientEnabledFeatures) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("compliancePolicy", m.GetCompliancePolicy())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("deviceConfiguration", m.GetDeviceConfiguration())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("inventory", m.GetInventory())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("modernApps", m.GetModernApps())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("resourceAccess", m.GetResourceAccess())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("windowsUpdateForBusiness", m.GetWindowsUpdateForBusiness())
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
func (m *ConfigurationManagerClientEnabledFeatures) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ConfigurationManagerClientEnabledFeatures) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCompliancePolicy sets the compliancePolicy property value. Whether compliance policy is managed by Intune
func (m *ConfigurationManagerClientEnabledFeatures) SetCompliancePolicy(value *bool)() {
    err := m.GetBackingStore().Set("compliancePolicy", value)
    if err != nil {
        panic(err)
    }
}
// SetDeviceConfiguration sets the deviceConfiguration property value. Whether device configuration is managed by Intune
func (m *ConfigurationManagerClientEnabledFeatures) SetDeviceConfiguration(value *bool)() {
    err := m.GetBackingStore().Set("deviceConfiguration", value)
    if err != nil {
        panic(err)
    }
}
// SetInventory sets the inventory property value. Whether inventory is managed by Intune
func (m *ConfigurationManagerClientEnabledFeatures) SetInventory(value *bool)() {
    err := m.GetBackingStore().Set("inventory", value)
    if err != nil {
        panic(err)
    }
}
// SetModernApps sets the modernApps property value. Whether modern application is managed by Intune
func (m *ConfigurationManagerClientEnabledFeatures) SetModernApps(value *bool)() {
    err := m.GetBackingStore().Set("modernApps", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ConfigurationManagerClientEnabledFeatures) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetResourceAccess sets the resourceAccess property value. Whether resource access is managed by Intune
func (m *ConfigurationManagerClientEnabledFeatures) SetResourceAccess(value *bool)() {
    err := m.GetBackingStore().Set("resourceAccess", value)
    if err != nil {
        panic(err)
    }
}
// SetWindowsUpdateForBusiness sets the windowsUpdateForBusiness property value. Whether Windows Update for Business is managed by Intune
func (m *ConfigurationManagerClientEnabledFeatures) SetWindowsUpdateForBusiness(value *bool)() {
    err := m.GetBackingStore().Set("windowsUpdateForBusiness", value)
    if err != nil {
        panic(err)
    }
}
type ConfigurationManagerClientEnabledFeaturesable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCompliancePolicy()(*bool)
    GetDeviceConfiguration()(*bool)
    GetInventory()(*bool)
    GetModernApps()(*bool)
    GetOdataType()(*string)
    GetResourceAccess()(*bool)
    GetWindowsUpdateForBusiness()(*bool)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCompliancePolicy(value *bool)()
    SetDeviceConfiguration(value *bool)()
    SetInventory(value *bool)()
    SetModernApps(value *bool)()
    SetOdataType(value *string)()
    SetResourceAccess(value *bool)()
    SetWindowsUpdateForBusiness(value *bool)()
}

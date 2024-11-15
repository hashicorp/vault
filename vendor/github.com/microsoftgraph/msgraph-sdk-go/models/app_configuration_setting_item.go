package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// AppConfigurationSettingItem contains properties for App configuration setting item.
type AppConfigurationSettingItem struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAppConfigurationSettingItem instantiates a new AppConfigurationSettingItem and sets the default values.
func NewAppConfigurationSettingItem()(*AppConfigurationSettingItem) {
    m := &AppConfigurationSettingItem{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAppConfigurationSettingItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAppConfigurationSettingItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAppConfigurationSettingItem(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AppConfigurationSettingItem) GetAdditionalData()(map[string]any) {
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
// GetAppConfigKey gets the appConfigKey property value. app configuration key.
// returns a *string when successful
func (m *AppConfigurationSettingItem) GetAppConfigKey()(*string) {
    val, err := m.GetBackingStore().Get("appConfigKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAppConfigKeyType gets the appConfigKeyType property value. App configuration key types.
// returns a *MdmAppConfigKeyType when successful
func (m *AppConfigurationSettingItem) GetAppConfigKeyType()(*MdmAppConfigKeyType) {
    val, err := m.GetBackingStore().Get("appConfigKeyType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MdmAppConfigKeyType)
    }
    return nil
}
// GetAppConfigKeyValue gets the appConfigKeyValue property value. app configuration key value.
// returns a *string when successful
func (m *AppConfigurationSettingItem) GetAppConfigKeyValue()(*string) {
    val, err := m.GetBackingStore().Get("appConfigKeyValue")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *AppConfigurationSettingItem) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AppConfigurationSettingItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["appConfigKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppConfigKey(val)
        }
        return nil
    }
    res["appConfigKeyType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMdmAppConfigKeyType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppConfigKeyType(val.(*MdmAppConfigKeyType))
        }
        return nil
    }
    res["appConfigKeyValue"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAppConfigKeyValue(val)
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
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *AppConfigurationSettingItem) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AppConfigurationSettingItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("appConfigKey", m.GetAppConfigKey())
        if err != nil {
            return err
        }
    }
    if m.GetAppConfigKeyType() != nil {
        cast := (*m.GetAppConfigKeyType()).String()
        err := writer.WriteStringValue("appConfigKeyType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("appConfigKeyValue", m.GetAppConfigKeyValue())
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
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *AppConfigurationSettingItem) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAppConfigKey sets the appConfigKey property value. app configuration key.
func (m *AppConfigurationSettingItem) SetAppConfigKey(value *string)() {
    err := m.GetBackingStore().Set("appConfigKey", value)
    if err != nil {
        panic(err)
    }
}
// SetAppConfigKeyType sets the appConfigKeyType property value. App configuration key types.
func (m *AppConfigurationSettingItem) SetAppConfigKeyType(value *MdmAppConfigKeyType)() {
    err := m.GetBackingStore().Set("appConfigKeyType", value)
    if err != nil {
        panic(err)
    }
}
// SetAppConfigKeyValue sets the appConfigKeyValue property value. app configuration key value.
func (m *AppConfigurationSettingItem) SetAppConfigKeyValue(value *string)() {
    err := m.GetBackingStore().Set("appConfigKeyValue", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AppConfigurationSettingItem) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AppConfigurationSettingItem) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type AppConfigurationSettingItemable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAppConfigKey()(*string)
    GetAppConfigKeyType()(*MdmAppConfigKeyType)
    GetAppConfigKeyValue()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    SetAppConfigKey(value *string)()
    SetAppConfigKeyType(value *MdmAppConfigKeyType)()
    SetAppConfigKeyValue(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
}

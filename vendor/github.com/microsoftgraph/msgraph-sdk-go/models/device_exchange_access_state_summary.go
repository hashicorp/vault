package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// DeviceExchangeAccessStateSummary device Exchange Access State summary
type DeviceExchangeAccessStateSummary struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewDeviceExchangeAccessStateSummary instantiates a new DeviceExchangeAccessStateSummary and sets the default values.
func NewDeviceExchangeAccessStateSummary()(*DeviceExchangeAccessStateSummary) {
    m := &DeviceExchangeAccessStateSummary{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateDeviceExchangeAccessStateSummaryFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeviceExchangeAccessStateSummaryFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeviceExchangeAccessStateSummary(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *DeviceExchangeAccessStateSummary) GetAdditionalData()(map[string]any) {
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
// GetAllowedDeviceCount gets the allowedDeviceCount property value. Total count of devices with Exchange Access State: Allowed.
// returns a *int32 when successful
func (m *DeviceExchangeAccessStateSummary) GetAllowedDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("allowedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *DeviceExchangeAccessStateSummary) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBlockedDeviceCount gets the blockedDeviceCount property value. Total count of devices with Exchange Access State: Blocked.
// returns a *int32 when successful
func (m *DeviceExchangeAccessStateSummary) GetBlockedDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("blockedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeviceExchangeAccessStateSummary) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowedDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowedDeviceCount(val)
        }
        return nil
    }
    res["blockedDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBlockedDeviceCount(val)
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
    res["quarantinedDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetQuarantinedDeviceCount(val)
        }
        return nil
    }
    res["unavailableDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnavailableDeviceCount(val)
        }
        return nil
    }
    res["unknownDeviceCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnknownDeviceCount(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *DeviceExchangeAccessStateSummary) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetQuarantinedDeviceCount gets the quarantinedDeviceCount property value. Total count of devices with Exchange Access State: Quarantined.
// returns a *int32 when successful
func (m *DeviceExchangeAccessStateSummary) GetQuarantinedDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("quarantinedDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUnavailableDeviceCount gets the unavailableDeviceCount property value. Total count of devices for which no Exchange Access State could be found.
// returns a *int32 when successful
func (m *DeviceExchangeAccessStateSummary) GetUnavailableDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("unavailableDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetUnknownDeviceCount gets the unknownDeviceCount property value. Total count of devices with Exchange Access State: Unknown.
// returns a *int32 when successful
func (m *DeviceExchangeAccessStateSummary) GetUnknownDeviceCount()(*int32) {
    val, err := m.GetBackingStore().Get("unknownDeviceCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeviceExchangeAccessStateSummary) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt32Value("allowedDeviceCount", m.GetAllowedDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("blockedDeviceCount", m.GetBlockedDeviceCount())
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
        err := writer.WriteInt32Value("quarantinedDeviceCount", m.GetQuarantinedDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("unavailableDeviceCount", m.GetUnavailableDeviceCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("unknownDeviceCount", m.GetUnknownDeviceCount())
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
func (m *DeviceExchangeAccessStateSummary) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowedDeviceCount sets the allowedDeviceCount property value. Total count of devices with Exchange Access State: Allowed.
func (m *DeviceExchangeAccessStateSummary) SetAllowedDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("allowedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *DeviceExchangeAccessStateSummary) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBlockedDeviceCount sets the blockedDeviceCount property value. Total count of devices with Exchange Access State: Blocked.
func (m *DeviceExchangeAccessStateSummary) SetBlockedDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("blockedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *DeviceExchangeAccessStateSummary) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetQuarantinedDeviceCount sets the quarantinedDeviceCount property value. Total count of devices with Exchange Access State: Quarantined.
func (m *DeviceExchangeAccessStateSummary) SetQuarantinedDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("quarantinedDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUnavailableDeviceCount sets the unavailableDeviceCount property value. Total count of devices for which no Exchange Access State could be found.
func (m *DeviceExchangeAccessStateSummary) SetUnavailableDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("unavailableDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
// SetUnknownDeviceCount sets the unknownDeviceCount property value. Total count of devices with Exchange Access State: Unknown.
func (m *DeviceExchangeAccessStateSummary) SetUnknownDeviceCount(value *int32)() {
    err := m.GetBackingStore().Set("unknownDeviceCount", value)
    if err != nil {
        panic(err)
    }
}
type DeviceExchangeAccessStateSummaryable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowedDeviceCount()(*int32)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBlockedDeviceCount()(*int32)
    GetOdataType()(*string)
    GetQuarantinedDeviceCount()(*int32)
    GetUnavailableDeviceCount()(*int32)
    GetUnknownDeviceCount()(*int32)
    SetAllowedDeviceCount(value *int32)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBlockedDeviceCount(value *int32)()
    SetOdataType(value *string)()
    SetQuarantinedDeviceCount(value *int32)()
    SetUnavailableDeviceCount(value *int32)()
    SetUnknownDeviceCount(value *int32)()
}

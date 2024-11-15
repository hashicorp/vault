package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Quota struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewQuota instantiates a new Quota and sets the default values.
func NewQuota()(*Quota) {
    m := &Quota{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateQuotaFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateQuotaFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewQuota(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Quota) GetAdditionalData()(map[string]any) {
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
func (m *Quota) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDeleted gets the deleted property value. Total space consumed by files in the recycle bin, in bytes. Read-only.
// returns a *int64 when successful
func (m *Quota) GetDeleted()(*int64) {
    val, err := m.GetBackingStore().Get("deleted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Quota) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["deleted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeleted(val)
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
    res["remaining"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemaining(val)
        }
        return nil
    }
    res["state"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetState(val)
        }
        return nil
    }
    res["storagePlanInformation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateStoragePlanInformationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStoragePlanInformation(val.(StoragePlanInformationable))
        }
        return nil
    }
    res["total"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotal(val)
        }
        return nil
    }
    res["used"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUsed(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Quota) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRemaining gets the remaining property value. Total space remaining before reaching the quota limit, in bytes. Read-only.
// returns a *int64 when successful
func (m *Quota) GetRemaining()(*int64) {
    val, err := m.GetBackingStore().Get("remaining")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetState gets the state property value. Enumeration value that indicates the state of the storage space. Read-only.
// returns a *string when successful
func (m *Quota) GetState()(*string) {
    val, err := m.GetBackingStore().Get("state")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStoragePlanInformation gets the storagePlanInformation property value. Information about the drive's storage quota plans. Only in Personal OneDrive.
// returns a StoragePlanInformationable when successful
func (m *Quota) GetStoragePlanInformation()(StoragePlanInformationable) {
    val, err := m.GetBackingStore().Get("storagePlanInformation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(StoragePlanInformationable)
    }
    return nil
}
// GetTotal gets the total property value. Total allowed storage space, in bytes. Read-only.
// returns a *int64 when successful
func (m *Quota) GetTotal()(*int64) {
    val, err := m.GetBackingStore().Get("total")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetUsed gets the used property value. Total space used, in bytes. Read-only.
// returns a *int64 when successful
func (m *Quota) GetUsed()(*int64) {
    val, err := m.GetBackingStore().Get("used")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Quota) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteInt64Value("deleted", m.GetDeleted())
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
        err := writer.WriteInt64Value("remaining", m.GetRemaining())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("state", m.GetState())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("storagePlanInformation", m.GetStoragePlanInformation())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("total", m.GetTotal())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("used", m.GetUsed())
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
func (m *Quota) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Quota) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDeleted sets the deleted property value. Total space consumed by files in the recycle bin, in bytes. Read-only.
func (m *Quota) SetDeleted(value *int64)() {
    err := m.GetBackingStore().Set("deleted", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Quota) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRemaining sets the remaining property value. Total space remaining before reaching the quota limit, in bytes. Read-only.
func (m *Quota) SetRemaining(value *int64)() {
    err := m.GetBackingStore().Set("remaining", value)
    if err != nil {
        panic(err)
    }
}
// SetState sets the state property value. Enumeration value that indicates the state of the storage space. Read-only.
func (m *Quota) SetState(value *string)() {
    err := m.GetBackingStore().Set("state", value)
    if err != nil {
        panic(err)
    }
}
// SetStoragePlanInformation sets the storagePlanInformation property value. Information about the drive's storage quota plans. Only in Personal OneDrive.
func (m *Quota) SetStoragePlanInformation(value StoragePlanInformationable)() {
    err := m.GetBackingStore().Set("storagePlanInformation", value)
    if err != nil {
        panic(err)
    }
}
// SetTotal sets the total property value. Total allowed storage space, in bytes. Read-only.
func (m *Quota) SetTotal(value *int64)() {
    err := m.GetBackingStore().Set("total", value)
    if err != nil {
        panic(err)
    }
}
// SetUsed sets the used property value. Total space used, in bytes. Read-only.
func (m *Quota) SetUsed(value *int64)() {
    err := m.GetBackingStore().Set("used", value)
    if err != nil {
        panic(err)
    }
}
type Quotaable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDeleted()(*int64)
    GetOdataType()(*string)
    GetRemaining()(*int64)
    GetState()(*string)
    GetStoragePlanInformation()(StoragePlanInformationable)
    GetTotal()(*int64)
    GetUsed()(*int64)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDeleted(value *int64)()
    SetOdataType(value *string)()
    SetRemaining(value *int64)()
    SetState(value *string)()
    SetStoragePlanInformation(value StoragePlanInformationable)()
    SetTotal(value *int64)()
    SetUsed(value *int64)()
}

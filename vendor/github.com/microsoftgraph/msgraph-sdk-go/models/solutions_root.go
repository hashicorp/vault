package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type SolutionsRoot struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSolutionsRoot instantiates a new SolutionsRoot and sets the default values.
func NewSolutionsRoot()(*SolutionsRoot) {
    m := &SolutionsRoot{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSolutionsRootFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSolutionsRootFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSolutionsRoot(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SolutionsRoot) GetAdditionalData()(map[string]any) {
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
func (m *SolutionsRoot) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBackupRestore gets the backupRestore property value. The backupRestore property
// returns a BackupRestoreRootable when successful
func (m *SolutionsRoot) GetBackupRestore()(BackupRestoreRootable) {
    val, err := m.GetBackingStore().Get("backupRestore")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(BackupRestoreRootable)
    }
    return nil
}
// GetBookingBusinesses gets the bookingBusinesses property value. The bookingBusinesses property
// returns a []BookingBusinessable when successful
func (m *SolutionsRoot) GetBookingBusinesses()([]BookingBusinessable) {
    val, err := m.GetBackingStore().Get("bookingBusinesses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BookingBusinessable)
    }
    return nil
}
// GetBookingCurrencies gets the bookingCurrencies property value. The bookingCurrencies property
// returns a []BookingCurrencyable when successful
func (m *SolutionsRoot) GetBookingCurrencies()([]BookingCurrencyable) {
    val, err := m.GetBackingStore().Get("bookingCurrencies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]BookingCurrencyable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *SolutionsRoot) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["backupRestore"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateBackupRestoreRootFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBackupRestore(val.(BackupRestoreRootable))
        }
        return nil
    }
    res["bookingBusinesses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBookingBusinessFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BookingBusinessable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BookingBusinessable)
                }
            }
            m.SetBookingBusinesses(res)
        }
        return nil
    }
    res["bookingCurrencies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateBookingCurrencyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]BookingCurrencyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(BookingCurrencyable)
                }
            }
            m.SetBookingCurrencies(res)
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
    res["virtualEvents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateVirtualEventsRootFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVirtualEvents(val.(VirtualEventsRootable))
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *SolutionsRoot) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVirtualEvents gets the virtualEvents property value. The virtualEvents property
// returns a VirtualEventsRootable when successful
func (m *SolutionsRoot) GetVirtualEvents()(VirtualEventsRootable) {
    val, err := m.GetBackingStore().Get("virtualEvents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(VirtualEventsRootable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SolutionsRoot) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("backupRestore", m.GetBackupRestore())
        if err != nil {
            return err
        }
    }
    if m.GetBookingBusinesses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetBookingBusinesses()))
        for i, v := range m.GetBookingBusinesses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("bookingBusinesses", cast)
        if err != nil {
            return err
        }
    }
    if m.GetBookingCurrencies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetBookingCurrencies()))
        for i, v := range m.GetBookingCurrencies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("bookingCurrencies", cast)
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
        err := writer.WriteObjectValue("virtualEvents", m.GetVirtualEvents())
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
func (m *SolutionsRoot) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SolutionsRoot) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBackupRestore sets the backupRestore property value. The backupRestore property
func (m *SolutionsRoot) SetBackupRestore(value BackupRestoreRootable)() {
    err := m.GetBackingStore().Set("backupRestore", value)
    if err != nil {
        panic(err)
    }
}
// SetBookingBusinesses sets the bookingBusinesses property value. The bookingBusinesses property
func (m *SolutionsRoot) SetBookingBusinesses(value []BookingBusinessable)() {
    err := m.GetBackingStore().Set("bookingBusinesses", value)
    if err != nil {
        panic(err)
    }
}
// SetBookingCurrencies sets the bookingCurrencies property value. The bookingCurrencies property
func (m *SolutionsRoot) SetBookingCurrencies(value []BookingCurrencyable)() {
    err := m.GetBackingStore().Set("bookingCurrencies", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SolutionsRoot) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetVirtualEvents sets the virtualEvents property value. The virtualEvents property
func (m *SolutionsRoot) SetVirtualEvents(value VirtualEventsRootable)() {
    err := m.GetBackingStore().Set("virtualEvents", value)
    if err != nil {
        panic(err)
    }
}
type SolutionsRootable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBackupRestore()(BackupRestoreRootable)
    GetBookingBusinesses()([]BookingBusinessable)
    GetBookingCurrencies()([]BookingCurrencyable)
    GetOdataType()(*string)
    GetVirtualEvents()(VirtualEventsRootable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBackupRestore(value BackupRestoreRootable)()
    SetBookingBusinesses(value []BookingBusinessable)()
    SetBookingCurrencies(value []BookingCurrencyable)()
    SetOdataType(value *string)()
    SetVirtualEvents(value VirtualEventsRootable)()
}

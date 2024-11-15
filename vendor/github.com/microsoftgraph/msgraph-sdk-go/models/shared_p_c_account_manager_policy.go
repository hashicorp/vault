package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// SharedPCAccountManagerPolicy sharedPC Account Manager Policy. Only applies when the account manager is enabled.
type SharedPCAccountManagerPolicy struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewSharedPCAccountManagerPolicy instantiates a new SharedPCAccountManagerPolicy and sets the default values.
func NewSharedPCAccountManagerPolicy()(*SharedPCAccountManagerPolicy) {
    m := &SharedPCAccountManagerPolicy{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateSharedPCAccountManagerPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateSharedPCAccountManagerPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewSharedPCAccountManagerPolicy(), nil
}
// GetAccountDeletionPolicy gets the accountDeletionPolicy property value. Possible values for when accounts are deleted on a shared PC.
// returns a *SharedPCAccountDeletionPolicyType when successful
func (m *SharedPCAccountManagerPolicy) GetAccountDeletionPolicy()(*SharedPCAccountDeletionPolicyType) {
    val, err := m.GetBackingStore().Get("accountDeletionPolicy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*SharedPCAccountDeletionPolicyType)
    }
    return nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *SharedPCAccountManagerPolicy) GetAdditionalData()(map[string]any) {
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
func (m *SharedPCAccountManagerPolicy) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCacheAccountsAboveDiskFreePercentage gets the cacheAccountsAboveDiskFreePercentage property value. Sets the percentage of available disk space a PC should have before it stops deleting cached shared PC accounts. Only applies when AccountDeletionPolicy is DiskSpaceThreshold or DiskSpaceThresholdOrInactiveThreshold. Valid values 0 to 100
// returns a *int32 when successful
func (m *SharedPCAccountManagerPolicy) GetCacheAccountsAboveDiskFreePercentage()(*int32) {
    val, err := m.GetBackingStore().Get("cacheAccountsAboveDiskFreePercentage")
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
func (m *SharedPCAccountManagerPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["accountDeletionPolicy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSharedPCAccountDeletionPolicyType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccountDeletionPolicy(val.(*SharedPCAccountDeletionPolicyType))
        }
        return nil
    }
    res["cacheAccountsAboveDiskFreePercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCacheAccountsAboveDiskFreePercentage(val)
        }
        return nil
    }
    res["inactiveThresholdDays"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInactiveThresholdDays(val)
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
    res["removeAccountsBelowDiskFreePercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemoveAccountsBelowDiskFreePercentage(val)
        }
        return nil
    }
    return res
}
// GetInactiveThresholdDays gets the inactiveThresholdDays property value. Specifies when the accounts will start being deleted when they have not been logged on during the specified period, given as number of days. Only applies when AccountDeletionPolicy is DiskSpaceThreshold or DiskSpaceThresholdOrInactiveThreshold.
// returns a *int32 when successful
func (m *SharedPCAccountManagerPolicy) GetInactiveThresholdDays()(*int32) {
    val, err := m.GetBackingStore().Get("inactiveThresholdDays")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *SharedPCAccountManagerPolicy) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRemoveAccountsBelowDiskFreePercentage gets the removeAccountsBelowDiskFreePercentage property value. Sets the percentage of disk space remaining on a PC before cached accounts will be deleted to free disk space. Accounts that have been inactive the longest will be deleted first. Only applies when AccountDeletionPolicy is DiskSpaceThresholdOrInactiveThreshold. Valid values 0 to 100
// returns a *int32 when successful
func (m *SharedPCAccountManagerPolicy) GetRemoveAccountsBelowDiskFreePercentage()(*int32) {
    val, err := m.GetBackingStore().Get("removeAccountsBelowDiskFreePercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *SharedPCAccountManagerPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetAccountDeletionPolicy() != nil {
        cast := (*m.GetAccountDeletionPolicy()).String()
        err := writer.WriteStringValue("accountDeletionPolicy", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("cacheAccountsAboveDiskFreePercentage", m.GetCacheAccountsAboveDiskFreePercentage())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("inactiveThresholdDays", m.GetInactiveThresholdDays())
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
        err := writer.WriteInt32Value("removeAccountsBelowDiskFreePercentage", m.GetRemoveAccountsBelowDiskFreePercentage())
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
// SetAccountDeletionPolicy sets the accountDeletionPolicy property value. Possible values for when accounts are deleted on a shared PC.
func (m *SharedPCAccountManagerPolicy) SetAccountDeletionPolicy(value *SharedPCAccountDeletionPolicyType)() {
    err := m.GetBackingStore().Set("accountDeletionPolicy", value)
    if err != nil {
        panic(err)
    }
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *SharedPCAccountManagerPolicy) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *SharedPCAccountManagerPolicy) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCacheAccountsAboveDiskFreePercentage sets the cacheAccountsAboveDiskFreePercentage property value. Sets the percentage of available disk space a PC should have before it stops deleting cached shared PC accounts. Only applies when AccountDeletionPolicy is DiskSpaceThreshold or DiskSpaceThresholdOrInactiveThreshold. Valid values 0 to 100
func (m *SharedPCAccountManagerPolicy) SetCacheAccountsAboveDiskFreePercentage(value *int32)() {
    err := m.GetBackingStore().Set("cacheAccountsAboveDiskFreePercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetInactiveThresholdDays sets the inactiveThresholdDays property value. Specifies when the accounts will start being deleted when they have not been logged on during the specified period, given as number of days. Only applies when AccountDeletionPolicy is DiskSpaceThreshold or DiskSpaceThresholdOrInactiveThreshold.
func (m *SharedPCAccountManagerPolicy) SetInactiveThresholdDays(value *int32)() {
    err := m.GetBackingStore().Set("inactiveThresholdDays", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *SharedPCAccountManagerPolicy) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRemoveAccountsBelowDiskFreePercentage sets the removeAccountsBelowDiskFreePercentage property value. Sets the percentage of disk space remaining on a PC before cached accounts will be deleted to free disk space. Accounts that have been inactive the longest will be deleted first. Only applies when AccountDeletionPolicy is DiskSpaceThresholdOrInactiveThreshold. Valid values 0 to 100
func (m *SharedPCAccountManagerPolicy) SetRemoveAccountsBelowDiskFreePercentage(value *int32)() {
    err := m.GetBackingStore().Set("removeAccountsBelowDiskFreePercentage", value)
    if err != nil {
        panic(err)
    }
}
type SharedPCAccountManagerPolicyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccountDeletionPolicy()(*SharedPCAccountDeletionPolicyType)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCacheAccountsAboveDiskFreePercentage()(*int32)
    GetInactiveThresholdDays()(*int32)
    GetOdataType()(*string)
    GetRemoveAccountsBelowDiskFreePercentage()(*int32)
    SetAccountDeletionPolicy(value *SharedPCAccountDeletionPolicyType)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCacheAccountsAboveDiskFreePercentage(value *int32)()
    SetInactiveThresholdDays(value *int32)()
    SetOdataType(value *string)()
    SetRemoveAccountsBelowDiskFreePercentage(value *int32)()
}

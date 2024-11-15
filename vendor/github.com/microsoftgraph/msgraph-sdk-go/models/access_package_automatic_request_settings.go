package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type AccessPackageAutomaticRequestSettings struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAccessPackageAutomaticRequestSettings instantiates a new AccessPackageAutomaticRequestSettings and sets the default values.
func NewAccessPackageAutomaticRequestSettings()(*AccessPackageAutomaticRequestSettings) {
    m := &AccessPackageAutomaticRequestSettings{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAccessPackageAutomaticRequestSettingsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAccessPackageAutomaticRequestSettingsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAccessPackageAutomaticRequestSettings(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AccessPackageAutomaticRequestSettings) GetAdditionalData()(map[string]any) {
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
func (m *AccessPackageAutomaticRequestSettings) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AccessPackageAutomaticRequestSettings) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["gracePeriodBeforeAccessRemoval"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGracePeriodBeforeAccessRemoval(val)
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
    res["removeAccessWhenTargetLeavesAllowedTargets"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRemoveAccessWhenTargetLeavesAllowedTargets(val)
        }
        return nil
    }
    res["requestAccessForAllowedTargets"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRequestAccessForAllowedTargets(val)
        }
        return nil
    }
    return res
}
// GetGracePeriodBeforeAccessRemoval gets the gracePeriodBeforeAccessRemoval property value. The duration for which access must be retained before the target's access is revoked once they leave the allowed target scope.
// returns a *ISODuration when successful
func (m *AccessPackageAutomaticRequestSettings) GetGracePeriodBeforeAccessRemoval()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("gracePeriodBeforeAccessRemoval")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *AccessPackageAutomaticRequestSettings) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRemoveAccessWhenTargetLeavesAllowedTargets gets the removeAccessWhenTargetLeavesAllowedTargets property value. Indicates whether automatic assignment must be removed for targets who move out of the allowed target scope.
// returns a *bool when successful
func (m *AccessPackageAutomaticRequestSettings) GetRemoveAccessWhenTargetLeavesAllowedTargets()(*bool) {
    val, err := m.GetBackingStore().Get("removeAccessWhenTargetLeavesAllowedTargets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetRequestAccessForAllowedTargets gets the requestAccessForAllowedTargets property value. If set to true, automatic assignments will be created for targets in the allowed target scope.
// returns a *bool when successful
func (m *AccessPackageAutomaticRequestSettings) GetRequestAccessForAllowedTargets()(*bool) {
    val, err := m.GetBackingStore().Get("requestAccessForAllowedTargets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AccessPackageAutomaticRequestSettings) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteISODurationValue("gracePeriodBeforeAccessRemoval", m.GetGracePeriodBeforeAccessRemoval())
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
        err := writer.WriteBoolValue("removeAccessWhenTargetLeavesAllowedTargets", m.GetRemoveAccessWhenTargetLeavesAllowedTargets())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("requestAccessForAllowedTargets", m.GetRequestAccessForAllowedTargets())
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
func (m *AccessPackageAutomaticRequestSettings) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AccessPackageAutomaticRequestSettings) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetGracePeriodBeforeAccessRemoval sets the gracePeriodBeforeAccessRemoval property value. The duration for which access must be retained before the target's access is revoked once they leave the allowed target scope.
func (m *AccessPackageAutomaticRequestSettings) SetGracePeriodBeforeAccessRemoval(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("gracePeriodBeforeAccessRemoval", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AccessPackageAutomaticRequestSettings) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRemoveAccessWhenTargetLeavesAllowedTargets sets the removeAccessWhenTargetLeavesAllowedTargets property value. Indicates whether automatic assignment must be removed for targets who move out of the allowed target scope.
func (m *AccessPackageAutomaticRequestSettings) SetRemoveAccessWhenTargetLeavesAllowedTargets(value *bool)() {
    err := m.GetBackingStore().Set("removeAccessWhenTargetLeavesAllowedTargets", value)
    if err != nil {
        panic(err)
    }
}
// SetRequestAccessForAllowedTargets sets the requestAccessForAllowedTargets property value. If set to true, automatic assignments will be created for targets in the allowed target scope.
func (m *AccessPackageAutomaticRequestSettings) SetRequestAccessForAllowedTargets(value *bool)() {
    err := m.GetBackingStore().Set("requestAccessForAllowedTargets", value)
    if err != nil {
        panic(err)
    }
}
type AccessPackageAutomaticRequestSettingsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetGracePeriodBeforeAccessRemoval()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetOdataType()(*string)
    GetRemoveAccessWhenTargetLeavesAllowedTargets()(*bool)
    GetRequestAccessForAllowedTargets()(*bool)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetGracePeriodBeforeAccessRemoval(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetOdataType(value *string)()
    SetRemoveAccessWhenTargetLeavesAllowedTargets(value *bool)()
    SetRequestAccessForAllowedTargets(value *bool)()
}

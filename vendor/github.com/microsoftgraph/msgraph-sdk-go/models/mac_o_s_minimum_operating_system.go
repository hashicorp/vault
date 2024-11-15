package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// MacOSMinimumOperatingSystem the minimum operating system required for a macOS app.
type MacOSMinimumOperatingSystem struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMacOSMinimumOperatingSystem instantiates a new MacOSMinimumOperatingSystem and sets the default values.
func NewMacOSMinimumOperatingSystem()(*MacOSMinimumOperatingSystem) {
    m := &MacOSMinimumOperatingSystem{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMacOSMinimumOperatingSystemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMacOSMinimumOperatingSystemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMacOSMinimumOperatingSystem(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *MacOSMinimumOperatingSystem) GetAdditionalData()(map[string]any) {
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
func (m *MacOSMinimumOperatingSystem) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MacOSMinimumOperatingSystem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["v10_10"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV1010(val)
        }
        return nil
    }
    res["v10_11"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV1011(val)
        }
        return nil
    }
    res["v10_12"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV1012(val)
        }
        return nil
    }
    res["v10_13"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV1013(val)
        }
        return nil
    }
    res["v10_14"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV1014(val)
        }
        return nil
    }
    res["v10_15"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV1015(val)
        }
        return nil
    }
    res["v10_7"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV107(val)
        }
        return nil
    }
    res["v10_8"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV108(val)
        }
        return nil
    }
    res["v10_9"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV109(val)
        }
        return nil
    }
    res["v11_0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV110(val)
        }
        return nil
    }
    res["v12_0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV120(val)
        }
        return nil
    }
    res["v13_0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV130(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *MacOSMinimumOperatingSystem) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetV1010 gets the v10_10 property value. When TRUE, indicates OS X 10.10 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV1010()(*bool) {
    val, err := m.GetBackingStore().Get("v10_10")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV1011 gets the v10_11 property value. When TRUE, indicates OS X 10.11 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV1011()(*bool) {
    val, err := m.GetBackingStore().Get("v10_11")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV1012 gets the v10_12 property value. When TRUE, indicates macOS 10.12 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV1012()(*bool) {
    val, err := m.GetBackingStore().Get("v10_12")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV1013 gets the v10_13 property value. When TRUE, indicates macOS 10.13 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV1013()(*bool) {
    val, err := m.GetBackingStore().Get("v10_13")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV1014 gets the v10_14 property value. When TRUE, indicates macOS 10.14 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV1014()(*bool) {
    val, err := m.GetBackingStore().Get("v10_14")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV1015 gets the v10_15 property value. When TRUE, indicates macOS 10.15 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV1015()(*bool) {
    val, err := m.GetBackingStore().Get("v10_15")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV107 gets the v10_7 property value. When TRUE, indicates Mac OS X 10.7 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV107()(*bool) {
    val, err := m.GetBackingStore().Get("v10_7")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV108 gets the v10_8 property value. When TRUE, indicates OS X 10.8 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV108()(*bool) {
    val, err := m.GetBackingStore().Get("v10_8")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV109 gets the v10_9 property value. When TRUE, indicates OS X 10.9 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV109()(*bool) {
    val, err := m.GetBackingStore().Get("v10_9")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV110 gets the v11_0 property value. When TRUE, indicates macOS 11.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV110()(*bool) {
    val, err := m.GetBackingStore().Get("v11_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV120 gets the v12_0 property value. When TRUE, indicates macOS 12.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV120()(*bool) {
    val, err := m.GetBackingStore().Get("v12_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV130 gets the v13_0 property value. When TRUE, indicates macOS 13.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
// returns a *bool when successful
func (m *MacOSMinimumOperatingSystem) GetV130()(*bool) {
    val, err := m.GetBackingStore().Get("v13_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MacOSMinimumOperatingSystem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v10_10", m.GetV1010())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v10_11", m.GetV1011())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v10_12", m.GetV1012())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v10_13", m.GetV1013())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v10_14", m.GetV1014())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v10_15", m.GetV1015())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v10_7", m.GetV107())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v10_8", m.GetV108())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v10_9", m.GetV109())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v11_0", m.GetV110())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v12_0", m.GetV120())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v13_0", m.GetV130())
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
func (m *MacOSMinimumOperatingSystem) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *MacOSMinimumOperatingSystem) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MacOSMinimumOperatingSystem) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetV1010 sets the v10_10 property value. When TRUE, indicates OS X 10.10 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV1010(value *bool)() {
    err := m.GetBackingStore().Set("v10_10", value)
    if err != nil {
        panic(err)
    }
}
// SetV1011 sets the v10_11 property value. When TRUE, indicates OS X 10.11 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV1011(value *bool)() {
    err := m.GetBackingStore().Set("v10_11", value)
    if err != nil {
        panic(err)
    }
}
// SetV1012 sets the v10_12 property value. When TRUE, indicates macOS 10.12 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV1012(value *bool)() {
    err := m.GetBackingStore().Set("v10_12", value)
    if err != nil {
        panic(err)
    }
}
// SetV1013 sets the v10_13 property value. When TRUE, indicates macOS 10.13 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV1013(value *bool)() {
    err := m.GetBackingStore().Set("v10_13", value)
    if err != nil {
        panic(err)
    }
}
// SetV1014 sets the v10_14 property value. When TRUE, indicates macOS 10.14 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV1014(value *bool)() {
    err := m.GetBackingStore().Set("v10_14", value)
    if err != nil {
        panic(err)
    }
}
// SetV1015 sets the v10_15 property value. When TRUE, indicates macOS 10.15 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV1015(value *bool)() {
    err := m.GetBackingStore().Set("v10_15", value)
    if err != nil {
        panic(err)
    }
}
// SetV107 sets the v10_7 property value. When TRUE, indicates Mac OS X 10.7 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV107(value *bool)() {
    err := m.GetBackingStore().Set("v10_7", value)
    if err != nil {
        panic(err)
    }
}
// SetV108 sets the v10_8 property value. When TRUE, indicates OS X 10.8 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV108(value *bool)() {
    err := m.GetBackingStore().Set("v10_8", value)
    if err != nil {
        panic(err)
    }
}
// SetV109 sets the v10_9 property value. When TRUE, indicates OS X 10.9 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV109(value *bool)() {
    err := m.GetBackingStore().Set("v10_9", value)
    if err != nil {
        panic(err)
    }
}
// SetV110 sets the v11_0 property value. When TRUE, indicates macOS 11.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV110(value *bool)() {
    err := m.GetBackingStore().Set("v11_0", value)
    if err != nil {
        panic(err)
    }
}
// SetV120 sets the v12_0 property value. When TRUE, indicates macOS 12.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV120(value *bool)() {
    err := m.GetBackingStore().Set("v12_0", value)
    if err != nil {
        panic(err)
    }
}
// SetV130 sets the v13_0 property value. When TRUE, indicates macOS 13.0 or later is required to install the app. When FALSE, indicates some other OS version is the minimum OS to install the app. Default value is FALSE.
func (m *MacOSMinimumOperatingSystem) SetV130(value *bool)() {
    err := m.GetBackingStore().Set("v13_0", value)
    if err != nil {
        panic(err)
    }
}
type MacOSMinimumOperatingSystemable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    GetV1010()(*bool)
    GetV1011()(*bool)
    GetV1012()(*bool)
    GetV1013()(*bool)
    GetV1014()(*bool)
    GetV1015()(*bool)
    GetV107()(*bool)
    GetV108()(*bool)
    GetV109()(*bool)
    GetV110()(*bool)
    GetV120()(*bool)
    GetV130()(*bool)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
    SetV1010(value *bool)()
    SetV1011(value *bool)()
    SetV1012(value *bool)()
    SetV1013(value *bool)()
    SetV1014(value *bool)()
    SetV1015(value *bool)()
    SetV107(value *bool)()
    SetV108(value *bool)()
    SetV109(value *bool)()
    SetV110(value *bool)()
    SetV120(value *bool)()
    SetV130(value *bool)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

// AndroidMinimumOperatingSystem contains properties for the minimum operating system required for an Android mobile app.
type AndroidMinimumOperatingSystem struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAndroidMinimumOperatingSystem instantiates a new AndroidMinimumOperatingSystem and sets the default values.
func NewAndroidMinimumOperatingSystem()(*AndroidMinimumOperatingSystem) {
    m := &AndroidMinimumOperatingSystem{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAndroidMinimumOperatingSystemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAndroidMinimumOperatingSystemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAndroidMinimumOperatingSystem(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *AndroidMinimumOperatingSystem) GetAdditionalData()(map[string]any) {
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
func (m *AndroidMinimumOperatingSystem) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AndroidMinimumOperatingSystem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["v10_0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV100(val)
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
    res["v4_0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV40(val)
        }
        return nil
    }
    res["v4_0_3"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV403(val)
        }
        return nil
    }
    res["v4_1"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV41(val)
        }
        return nil
    }
    res["v4_2"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV42(val)
        }
        return nil
    }
    res["v4_3"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV43(val)
        }
        return nil
    }
    res["v4_4"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV44(val)
        }
        return nil
    }
    res["v5_0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV50(val)
        }
        return nil
    }
    res["v5_1"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV51(val)
        }
        return nil
    }
    res["v6_0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV60(val)
        }
        return nil
    }
    res["v7_0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV70(val)
        }
        return nil
    }
    res["v7_1"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV71(val)
        }
        return nil
    }
    res["v8_0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV80(val)
        }
        return nil
    }
    res["v8_1"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV81(val)
        }
        return nil
    }
    res["v9_0"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetV90(val)
        }
        return nil
    }
    return res
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *AndroidMinimumOperatingSystem) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetV100 gets the v10_0 property value. When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV100()(*bool) {
    val, err := m.GetBackingStore().Get("v10_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV110 gets the v11_0 property value. When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV110()(*bool) {
    val, err := m.GetBackingStore().Get("v11_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV40 gets the v4_0 property value. When TRUE, only Version 4.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV40()(*bool) {
    val, err := m.GetBackingStore().Get("v4_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV403 gets the v4_0_3 property value. When TRUE, only Version 4.0.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV403()(*bool) {
    val, err := m.GetBackingStore().Get("v4_0_3")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV41 gets the v4_1 property value. When TRUE, only Version 4.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV41()(*bool) {
    val, err := m.GetBackingStore().Get("v4_1")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV42 gets the v4_2 property value. When TRUE, only Version 4.2 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV42()(*bool) {
    val, err := m.GetBackingStore().Get("v4_2")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV43 gets the v4_3 property value. When TRUE, only Version 4.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV43()(*bool) {
    val, err := m.GetBackingStore().Get("v4_3")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV44 gets the v4_4 property value. When TRUE, only Version 4.4 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV44()(*bool) {
    val, err := m.GetBackingStore().Get("v4_4")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV50 gets the v5_0 property value. When TRUE, only Version 5.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV50()(*bool) {
    val, err := m.GetBackingStore().Get("v5_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV51 gets the v5_1 property value. When TRUE, only Version 5.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV51()(*bool) {
    val, err := m.GetBackingStore().Get("v5_1")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV60 gets the v6_0 property value. When TRUE, only Version 6.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV60()(*bool) {
    val, err := m.GetBackingStore().Get("v6_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV70 gets the v7_0 property value. When TRUE, only Version 7.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV70()(*bool) {
    val, err := m.GetBackingStore().Get("v7_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV71 gets the v7_1 property value. When TRUE, only Version 7.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV71()(*bool) {
    val, err := m.GetBackingStore().Get("v7_1")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV80 gets the v8_0 property value. When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV80()(*bool) {
    val, err := m.GetBackingStore().Get("v8_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV81 gets the v8_1 property value. When TRUE, only Version 8.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV81()(*bool) {
    val, err := m.GetBackingStore().Get("v8_1")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetV90 gets the v9_0 property value. When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
// returns a *bool when successful
func (m *AndroidMinimumOperatingSystem) GetV90()(*bool) {
    val, err := m.GetBackingStore().Get("v9_0")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AndroidMinimumOperatingSystem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v10_0", m.GetV100())
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
        err := writer.WriteBoolValue("v4_0", m.GetV40())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v4_0_3", m.GetV403())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v4_1", m.GetV41())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v4_2", m.GetV42())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v4_3", m.GetV43())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v4_4", m.GetV44())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v5_0", m.GetV50())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v5_1", m.GetV51())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v6_0", m.GetV60())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v7_0", m.GetV70())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v7_1", m.GetV71())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v8_0", m.GetV80())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v8_1", m.GetV81())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("v9_0", m.GetV90())
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
func (m *AndroidMinimumOperatingSystem) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *AndroidMinimumOperatingSystem) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *AndroidMinimumOperatingSystem) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetV100 sets the v10_0 property value. When TRUE, only Version 10.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV100(value *bool)() {
    err := m.GetBackingStore().Set("v10_0", value)
    if err != nil {
        panic(err)
    }
}
// SetV110 sets the v11_0 property value. When TRUE, only Version 11.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV110(value *bool)() {
    err := m.GetBackingStore().Set("v11_0", value)
    if err != nil {
        panic(err)
    }
}
// SetV40 sets the v4_0 property value. When TRUE, only Version 4.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV40(value *bool)() {
    err := m.GetBackingStore().Set("v4_0", value)
    if err != nil {
        panic(err)
    }
}
// SetV403 sets the v4_0_3 property value. When TRUE, only Version 4.0.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV403(value *bool)() {
    err := m.GetBackingStore().Set("v4_0_3", value)
    if err != nil {
        panic(err)
    }
}
// SetV41 sets the v4_1 property value. When TRUE, only Version 4.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV41(value *bool)() {
    err := m.GetBackingStore().Set("v4_1", value)
    if err != nil {
        panic(err)
    }
}
// SetV42 sets the v4_2 property value. When TRUE, only Version 4.2 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV42(value *bool)() {
    err := m.GetBackingStore().Set("v4_2", value)
    if err != nil {
        panic(err)
    }
}
// SetV43 sets the v4_3 property value. When TRUE, only Version 4.3 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV43(value *bool)() {
    err := m.GetBackingStore().Set("v4_3", value)
    if err != nil {
        panic(err)
    }
}
// SetV44 sets the v4_4 property value. When TRUE, only Version 4.4 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV44(value *bool)() {
    err := m.GetBackingStore().Set("v4_4", value)
    if err != nil {
        panic(err)
    }
}
// SetV50 sets the v5_0 property value. When TRUE, only Version 5.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV50(value *bool)() {
    err := m.GetBackingStore().Set("v5_0", value)
    if err != nil {
        panic(err)
    }
}
// SetV51 sets the v5_1 property value. When TRUE, only Version 5.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV51(value *bool)() {
    err := m.GetBackingStore().Set("v5_1", value)
    if err != nil {
        panic(err)
    }
}
// SetV60 sets the v6_0 property value. When TRUE, only Version 6.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV60(value *bool)() {
    err := m.GetBackingStore().Set("v6_0", value)
    if err != nil {
        panic(err)
    }
}
// SetV70 sets the v7_0 property value. When TRUE, only Version 7.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV70(value *bool)() {
    err := m.GetBackingStore().Set("v7_0", value)
    if err != nil {
        panic(err)
    }
}
// SetV71 sets the v7_1 property value. When TRUE, only Version 7.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV71(value *bool)() {
    err := m.GetBackingStore().Set("v7_1", value)
    if err != nil {
        panic(err)
    }
}
// SetV80 sets the v8_0 property value. When TRUE, only Version 8.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV80(value *bool)() {
    err := m.GetBackingStore().Set("v8_0", value)
    if err != nil {
        panic(err)
    }
}
// SetV81 sets the v8_1 property value. When TRUE, only Version 8.1 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV81(value *bool)() {
    err := m.GetBackingStore().Set("v8_1", value)
    if err != nil {
        panic(err)
    }
}
// SetV90 sets the v9_0 property value. When TRUE, only Version 9.0 or later is supported. Default value is FALSE. Exactly one of the minimum operating system boolean values will be TRUE.
func (m *AndroidMinimumOperatingSystem) SetV90(value *bool)() {
    err := m.GetBackingStore().Set("v9_0", value)
    if err != nil {
        panic(err)
    }
}
type AndroidMinimumOperatingSystemable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    GetV100()(*bool)
    GetV110()(*bool)
    GetV40()(*bool)
    GetV403()(*bool)
    GetV41()(*bool)
    GetV42()(*bool)
    GetV43()(*bool)
    GetV44()(*bool)
    GetV50()(*bool)
    GetV51()(*bool)
    GetV60()(*bool)
    GetV70()(*bool)
    GetV71()(*bool)
    GetV80()(*bool)
    GetV81()(*bool)
    GetV90()(*bool)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
    SetV100(value *bool)()
    SetV110(value *bool)()
    SetV40(value *bool)()
    SetV403(value *bool)()
    SetV41(value *bool)()
    SetV42(value *bool)()
    SetV43(value *bool)()
    SetV44(value *bool)()
    SetV50(value *bool)()
    SetV51(value *bool)()
    SetV60(value *bool)()
    SetV70(value *bool)()
    SetV71(value *bool)()
    SetV80(value *bool)()
    SetV81(value *bool)()
    SetV90(value *bool)()
}

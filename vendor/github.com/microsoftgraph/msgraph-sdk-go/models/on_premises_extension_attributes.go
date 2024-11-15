package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type OnPremisesExtensionAttributes struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewOnPremisesExtensionAttributes instantiates a new OnPremisesExtensionAttributes and sets the default values.
func NewOnPremisesExtensionAttributes()(*OnPremisesExtensionAttributes) {
    m := &OnPremisesExtensionAttributes{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateOnPremisesExtensionAttributesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateOnPremisesExtensionAttributesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewOnPremisesExtensionAttributes(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *OnPremisesExtensionAttributes) GetAdditionalData()(map[string]any) {
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
func (m *OnPremisesExtensionAttributes) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetExtensionAttribute1 gets the extensionAttribute1 property value. First customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute1()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute1")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute10 gets the extensionAttribute10 property value. Tenth customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute10()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute10")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute11 gets the extensionAttribute11 property value. Eleventh customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute11()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute11")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute12 gets the extensionAttribute12 property value. Twelfth customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute12()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute12")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute13 gets the extensionAttribute13 property value. Thirteenth customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute13()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute13")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute14 gets the extensionAttribute14 property value. Fourteenth customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute14()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute14")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute15 gets the extensionAttribute15 property value. Fifteenth customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute15()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute15")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute2 gets the extensionAttribute2 property value. Second customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute2()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute2")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute3 gets the extensionAttribute3 property value. Third customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute3()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute3")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute4 gets the extensionAttribute4 property value. Fourth customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute4()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute4")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute5 gets the extensionAttribute5 property value. Fifth customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute5()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute5")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute6 gets the extensionAttribute6 property value. Sixth customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute6()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute6")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute7 gets the extensionAttribute7 property value. Seventh customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute7()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute7")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute8 gets the extensionAttribute8 property value. Eighth customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute8()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute8")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensionAttribute9 gets the extensionAttribute9 property value. Ninth customizable extension attribute.
// returns a *string when successful
func (m *OnPremisesExtensionAttributes) GetExtensionAttribute9()(*string) {
    val, err := m.GetBackingStore().Get("extensionAttribute9")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *OnPremisesExtensionAttributes) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["extensionAttribute1"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute1(val)
        }
        return nil
    }
    res["extensionAttribute10"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute10(val)
        }
        return nil
    }
    res["extensionAttribute11"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute11(val)
        }
        return nil
    }
    res["extensionAttribute12"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute12(val)
        }
        return nil
    }
    res["extensionAttribute13"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute13(val)
        }
        return nil
    }
    res["extensionAttribute14"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute14(val)
        }
        return nil
    }
    res["extensionAttribute15"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute15(val)
        }
        return nil
    }
    res["extensionAttribute2"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute2(val)
        }
        return nil
    }
    res["extensionAttribute3"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute3(val)
        }
        return nil
    }
    res["extensionAttribute4"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute4(val)
        }
        return nil
    }
    res["extensionAttribute5"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute5(val)
        }
        return nil
    }
    res["extensionAttribute6"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute6(val)
        }
        return nil
    }
    res["extensionAttribute7"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute7(val)
        }
        return nil
    }
    res["extensionAttribute8"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute8(val)
        }
        return nil
    }
    res["extensionAttribute9"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExtensionAttribute9(val)
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
func (m *OnPremisesExtensionAttributes) GetOdataType()(*string) {
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
func (m *OnPremisesExtensionAttributes) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("extensionAttribute1", m.GetExtensionAttribute1())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute10", m.GetExtensionAttribute10())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute11", m.GetExtensionAttribute11())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute12", m.GetExtensionAttribute12())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute13", m.GetExtensionAttribute13())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute14", m.GetExtensionAttribute14())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute15", m.GetExtensionAttribute15())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute2", m.GetExtensionAttribute2())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute3", m.GetExtensionAttribute3())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute4", m.GetExtensionAttribute4())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute5", m.GetExtensionAttribute5())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute6", m.GetExtensionAttribute6())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute7", m.GetExtensionAttribute7())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute8", m.GetExtensionAttribute8())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("extensionAttribute9", m.GetExtensionAttribute9())
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
func (m *OnPremisesExtensionAttributes) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *OnPremisesExtensionAttributes) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetExtensionAttribute1 sets the extensionAttribute1 property value. First customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute1(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute1", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute10 sets the extensionAttribute10 property value. Tenth customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute10(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute10", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute11 sets the extensionAttribute11 property value. Eleventh customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute11(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute11", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute12 sets the extensionAttribute12 property value. Twelfth customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute12(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute12", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute13 sets the extensionAttribute13 property value. Thirteenth customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute13(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute13", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute14 sets the extensionAttribute14 property value. Fourteenth customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute14(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute14", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute15 sets the extensionAttribute15 property value. Fifteenth customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute15(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute15", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute2 sets the extensionAttribute2 property value. Second customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute2(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute2", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute3 sets the extensionAttribute3 property value. Third customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute3(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute3", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute4 sets the extensionAttribute4 property value. Fourth customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute4(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute4", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute5 sets the extensionAttribute5 property value. Fifth customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute5(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute5", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute6 sets the extensionAttribute6 property value. Sixth customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute6(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute6", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute7 sets the extensionAttribute7 property value. Seventh customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute7(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute7", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute8 sets the extensionAttribute8 property value. Eighth customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute8(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute8", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensionAttribute9 sets the extensionAttribute9 property value. Ninth customizable extension attribute.
func (m *OnPremisesExtensionAttributes) SetExtensionAttribute9(value *string)() {
    err := m.GetBackingStore().Set("extensionAttribute9", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *OnPremisesExtensionAttributes) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type OnPremisesExtensionAttributesable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetExtensionAttribute1()(*string)
    GetExtensionAttribute10()(*string)
    GetExtensionAttribute11()(*string)
    GetExtensionAttribute12()(*string)
    GetExtensionAttribute13()(*string)
    GetExtensionAttribute14()(*string)
    GetExtensionAttribute15()(*string)
    GetExtensionAttribute2()(*string)
    GetExtensionAttribute3()(*string)
    GetExtensionAttribute4()(*string)
    GetExtensionAttribute5()(*string)
    GetExtensionAttribute6()(*string)
    GetExtensionAttribute7()(*string)
    GetExtensionAttribute8()(*string)
    GetExtensionAttribute9()(*string)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetExtensionAttribute1(value *string)()
    SetExtensionAttribute10(value *string)()
    SetExtensionAttribute11(value *string)()
    SetExtensionAttribute12(value *string)()
    SetExtensionAttribute13(value *string)()
    SetExtensionAttribute14(value *string)()
    SetExtensionAttribute15(value *string)()
    SetExtensionAttribute2(value *string)()
    SetExtensionAttribute3(value *string)()
    SetExtensionAttribute4(value *string)()
    SetExtensionAttribute5(value *string)()
    SetExtensionAttribute6(value *string)()
    SetExtensionAttribute7(value *string)()
    SetExtensionAttribute8(value *string)()
    SetExtensionAttribute9(value *string)()
    SetOdataType(value *string)()
}

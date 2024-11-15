package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type PlannerCategoryDescriptions struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewPlannerCategoryDescriptions instantiates a new PlannerCategoryDescriptions and sets the default values.
func NewPlannerCategoryDescriptions()(*PlannerCategoryDescriptions) {
    m := &PlannerCategoryDescriptions{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreatePlannerCategoryDescriptionsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePlannerCategoryDescriptionsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPlannerCategoryDescriptions(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *PlannerCategoryDescriptions) GetAdditionalData()(map[string]any) {
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
func (m *PlannerCategoryDescriptions) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCategory1 gets the category1 property value. The label associated with Category 1
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory1()(*string) {
    val, err := m.GetBackingStore().Get("category1")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory10 gets the category10 property value. The label associated with Category 10
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory10()(*string) {
    val, err := m.GetBackingStore().Get("category10")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory11 gets the category11 property value. The label associated with Category 11
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory11()(*string) {
    val, err := m.GetBackingStore().Get("category11")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory12 gets the category12 property value. The label associated with Category 12
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory12()(*string) {
    val, err := m.GetBackingStore().Get("category12")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory13 gets the category13 property value. The label associated with Category 13
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory13()(*string) {
    val, err := m.GetBackingStore().Get("category13")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory14 gets the category14 property value. The label associated with Category 14
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory14()(*string) {
    val, err := m.GetBackingStore().Get("category14")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory15 gets the category15 property value. The label associated with Category 15
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory15()(*string) {
    val, err := m.GetBackingStore().Get("category15")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory16 gets the category16 property value. The label associated with Category 16
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory16()(*string) {
    val, err := m.GetBackingStore().Get("category16")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory17 gets the category17 property value. The label associated with Category 17
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory17()(*string) {
    val, err := m.GetBackingStore().Get("category17")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory18 gets the category18 property value. The label associated with Category 18
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory18()(*string) {
    val, err := m.GetBackingStore().Get("category18")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory19 gets the category19 property value. The label associated with Category 19
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory19()(*string) {
    val, err := m.GetBackingStore().Get("category19")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory2 gets the category2 property value. The label associated with Category 2
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory2()(*string) {
    val, err := m.GetBackingStore().Get("category2")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory20 gets the category20 property value. The label associated with Category 20
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory20()(*string) {
    val, err := m.GetBackingStore().Get("category20")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory21 gets the category21 property value. The label associated with Category 21
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory21()(*string) {
    val, err := m.GetBackingStore().Get("category21")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory22 gets the category22 property value. The label associated with Category 22
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory22()(*string) {
    val, err := m.GetBackingStore().Get("category22")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory23 gets the category23 property value. The label associated with Category 23
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory23()(*string) {
    val, err := m.GetBackingStore().Get("category23")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory24 gets the category24 property value. The label associated with Category 24
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory24()(*string) {
    val, err := m.GetBackingStore().Get("category24")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory25 gets the category25 property value. The label associated with Category 25
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory25()(*string) {
    val, err := m.GetBackingStore().Get("category25")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory3 gets the category3 property value. The label associated with Category 3
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory3()(*string) {
    val, err := m.GetBackingStore().Get("category3")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory4 gets the category4 property value. The label associated with Category 4
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory4()(*string) {
    val, err := m.GetBackingStore().Get("category4")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory5 gets the category5 property value. The label associated with Category 5
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory5()(*string) {
    val, err := m.GetBackingStore().Get("category5")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory6 gets the category6 property value. The label associated with Category 6
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory6()(*string) {
    val, err := m.GetBackingStore().Get("category6")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory7 gets the category7 property value. The label associated with Category 7
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory7()(*string) {
    val, err := m.GetBackingStore().Get("category7")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory8 gets the category8 property value. The label associated with Category 8
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory8()(*string) {
    val, err := m.GetBackingStore().Get("category8")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCategory9 gets the category9 property value. The label associated with Category 9
// returns a *string when successful
func (m *PlannerCategoryDescriptions) GetCategory9()(*string) {
    val, err := m.GetBackingStore().Get("category9")
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
func (m *PlannerCategoryDescriptions) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["category1"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory1(val)
        }
        return nil
    }
    res["category10"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory10(val)
        }
        return nil
    }
    res["category11"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory11(val)
        }
        return nil
    }
    res["category12"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory12(val)
        }
        return nil
    }
    res["category13"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory13(val)
        }
        return nil
    }
    res["category14"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory14(val)
        }
        return nil
    }
    res["category15"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory15(val)
        }
        return nil
    }
    res["category16"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory16(val)
        }
        return nil
    }
    res["category17"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory17(val)
        }
        return nil
    }
    res["category18"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory18(val)
        }
        return nil
    }
    res["category19"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory19(val)
        }
        return nil
    }
    res["category2"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory2(val)
        }
        return nil
    }
    res["category20"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory20(val)
        }
        return nil
    }
    res["category21"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory21(val)
        }
        return nil
    }
    res["category22"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory22(val)
        }
        return nil
    }
    res["category23"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory23(val)
        }
        return nil
    }
    res["category24"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory24(val)
        }
        return nil
    }
    res["category25"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory25(val)
        }
        return nil
    }
    res["category3"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory3(val)
        }
        return nil
    }
    res["category4"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory4(val)
        }
        return nil
    }
    res["category5"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory5(val)
        }
        return nil
    }
    res["category6"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory6(val)
        }
        return nil
    }
    res["category7"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory7(val)
        }
        return nil
    }
    res["category8"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory8(val)
        }
        return nil
    }
    res["category9"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCategory9(val)
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
func (m *PlannerCategoryDescriptions) GetOdataType()(*string) {
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
func (m *PlannerCategoryDescriptions) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("category1", m.GetCategory1())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category10", m.GetCategory10())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category11", m.GetCategory11())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category12", m.GetCategory12())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category13", m.GetCategory13())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category14", m.GetCategory14())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category15", m.GetCategory15())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category16", m.GetCategory16())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category17", m.GetCategory17())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category18", m.GetCategory18())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category19", m.GetCategory19())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category2", m.GetCategory2())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category20", m.GetCategory20())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category21", m.GetCategory21())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category22", m.GetCategory22())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category23", m.GetCategory23())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category24", m.GetCategory24())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category25", m.GetCategory25())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category3", m.GetCategory3())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category4", m.GetCategory4())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category5", m.GetCategory5())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category6", m.GetCategory6())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category7", m.GetCategory7())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category8", m.GetCategory8())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("category9", m.GetCategory9())
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
func (m *PlannerCategoryDescriptions) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *PlannerCategoryDescriptions) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCategory1 sets the category1 property value. The label associated with Category 1
func (m *PlannerCategoryDescriptions) SetCategory1(value *string)() {
    err := m.GetBackingStore().Set("category1", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory10 sets the category10 property value. The label associated with Category 10
func (m *PlannerCategoryDescriptions) SetCategory10(value *string)() {
    err := m.GetBackingStore().Set("category10", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory11 sets the category11 property value. The label associated with Category 11
func (m *PlannerCategoryDescriptions) SetCategory11(value *string)() {
    err := m.GetBackingStore().Set("category11", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory12 sets the category12 property value. The label associated with Category 12
func (m *PlannerCategoryDescriptions) SetCategory12(value *string)() {
    err := m.GetBackingStore().Set("category12", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory13 sets the category13 property value. The label associated with Category 13
func (m *PlannerCategoryDescriptions) SetCategory13(value *string)() {
    err := m.GetBackingStore().Set("category13", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory14 sets the category14 property value. The label associated with Category 14
func (m *PlannerCategoryDescriptions) SetCategory14(value *string)() {
    err := m.GetBackingStore().Set("category14", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory15 sets the category15 property value. The label associated with Category 15
func (m *PlannerCategoryDescriptions) SetCategory15(value *string)() {
    err := m.GetBackingStore().Set("category15", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory16 sets the category16 property value. The label associated with Category 16
func (m *PlannerCategoryDescriptions) SetCategory16(value *string)() {
    err := m.GetBackingStore().Set("category16", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory17 sets the category17 property value. The label associated with Category 17
func (m *PlannerCategoryDescriptions) SetCategory17(value *string)() {
    err := m.GetBackingStore().Set("category17", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory18 sets the category18 property value. The label associated with Category 18
func (m *PlannerCategoryDescriptions) SetCategory18(value *string)() {
    err := m.GetBackingStore().Set("category18", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory19 sets the category19 property value. The label associated with Category 19
func (m *PlannerCategoryDescriptions) SetCategory19(value *string)() {
    err := m.GetBackingStore().Set("category19", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory2 sets the category2 property value. The label associated with Category 2
func (m *PlannerCategoryDescriptions) SetCategory2(value *string)() {
    err := m.GetBackingStore().Set("category2", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory20 sets the category20 property value. The label associated with Category 20
func (m *PlannerCategoryDescriptions) SetCategory20(value *string)() {
    err := m.GetBackingStore().Set("category20", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory21 sets the category21 property value. The label associated with Category 21
func (m *PlannerCategoryDescriptions) SetCategory21(value *string)() {
    err := m.GetBackingStore().Set("category21", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory22 sets the category22 property value. The label associated with Category 22
func (m *PlannerCategoryDescriptions) SetCategory22(value *string)() {
    err := m.GetBackingStore().Set("category22", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory23 sets the category23 property value. The label associated with Category 23
func (m *PlannerCategoryDescriptions) SetCategory23(value *string)() {
    err := m.GetBackingStore().Set("category23", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory24 sets the category24 property value. The label associated with Category 24
func (m *PlannerCategoryDescriptions) SetCategory24(value *string)() {
    err := m.GetBackingStore().Set("category24", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory25 sets the category25 property value. The label associated with Category 25
func (m *PlannerCategoryDescriptions) SetCategory25(value *string)() {
    err := m.GetBackingStore().Set("category25", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory3 sets the category3 property value. The label associated with Category 3
func (m *PlannerCategoryDescriptions) SetCategory3(value *string)() {
    err := m.GetBackingStore().Set("category3", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory4 sets the category4 property value. The label associated with Category 4
func (m *PlannerCategoryDescriptions) SetCategory4(value *string)() {
    err := m.GetBackingStore().Set("category4", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory5 sets the category5 property value. The label associated with Category 5
func (m *PlannerCategoryDescriptions) SetCategory5(value *string)() {
    err := m.GetBackingStore().Set("category5", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory6 sets the category6 property value. The label associated with Category 6
func (m *PlannerCategoryDescriptions) SetCategory6(value *string)() {
    err := m.GetBackingStore().Set("category6", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory7 sets the category7 property value. The label associated with Category 7
func (m *PlannerCategoryDescriptions) SetCategory7(value *string)() {
    err := m.GetBackingStore().Set("category7", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory8 sets the category8 property value. The label associated with Category 8
func (m *PlannerCategoryDescriptions) SetCategory8(value *string)() {
    err := m.GetBackingStore().Set("category8", value)
    if err != nil {
        panic(err)
    }
}
// SetCategory9 sets the category9 property value. The label associated with Category 9
func (m *PlannerCategoryDescriptions) SetCategory9(value *string)() {
    err := m.GetBackingStore().Set("category9", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *PlannerCategoryDescriptions) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type PlannerCategoryDescriptionsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCategory1()(*string)
    GetCategory10()(*string)
    GetCategory11()(*string)
    GetCategory12()(*string)
    GetCategory13()(*string)
    GetCategory14()(*string)
    GetCategory15()(*string)
    GetCategory16()(*string)
    GetCategory17()(*string)
    GetCategory18()(*string)
    GetCategory19()(*string)
    GetCategory2()(*string)
    GetCategory20()(*string)
    GetCategory21()(*string)
    GetCategory22()(*string)
    GetCategory23()(*string)
    GetCategory24()(*string)
    GetCategory25()(*string)
    GetCategory3()(*string)
    GetCategory4()(*string)
    GetCategory5()(*string)
    GetCategory6()(*string)
    GetCategory7()(*string)
    GetCategory8()(*string)
    GetCategory9()(*string)
    GetOdataType()(*string)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCategory1(value *string)()
    SetCategory10(value *string)()
    SetCategory11(value *string)()
    SetCategory12(value *string)()
    SetCategory13(value *string)()
    SetCategory14(value *string)()
    SetCategory15(value *string)()
    SetCategory16(value *string)()
    SetCategory17(value *string)()
    SetCategory18(value *string)()
    SetCategory19(value *string)()
    SetCategory2(value *string)()
    SetCategory20(value *string)()
    SetCategory21(value *string)()
    SetCategory22(value *string)()
    SetCategory23(value *string)()
    SetCategory24(value *string)()
    SetCategory25(value *string)()
    SetCategory3(value *string)()
    SetCategory4(value *string)()
    SetCategory5(value *string)()
    SetCategory6(value *string)()
    SetCategory7(value *string)()
    SetCategory8(value *string)()
    SetCategory9(value *string)()
    SetOdataType(value *string)()
}

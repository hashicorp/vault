package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type WorkbookWorksheetProtectionOptions struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewWorkbookWorksheetProtectionOptions instantiates a new WorkbookWorksheetProtectionOptions and sets the default values.
func NewWorkbookWorksheetProtectionOptions()(*WorkbookWorksheetProtectionOptions) {
    m := &WorkbookWorksheetProtectionOptions{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateWorkbookWorksheetProtectionOptionsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookWorksheetProtectionOptionsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookWorksheetProtectionOptions(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *WorkbookWorksheetProtectionOptions) GetAdditionalData()(map[string]any) {
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
// GetAllowAutoFilter gets the allowAutoFilter property value. Represents the worksheet protection option of allowing using auto filter feature.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowAutoFilter()(*bool) {
    val, err := m.GetBackingStore().Get("allowAutoFilter")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowDeleteColumns gets the allowDeleteColumns property value. Represents the worksheet protection option of allowing deleting columns.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowDeleteColumns()(*bool) {
    val, err := m.GetBackingStore().Get("allowDeleteColumns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowDeleteRows gets the allowDeleteRows property value. Represents the worksheet protection option of allowing deleting rows.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowDeleteRows()(*bool) {
    val, err := m.GetBackingStore().Get("allowDeleteRows")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowFormatCells gets the allowFormatCells property value. Represents the worksheet protection option of allowing formatting cells.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowFormatCells()(*bool) {
    val, err := m.GetBackingStore().Get("allowFormatCells")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowFormatColumns gets the allowFormatColumns property value. Represents the worksheet protection option of allowing formatting columns.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowFormatColumns()(*bool) {
    val, err := m.GetBackingStore().Get("allowFormatColumns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowFormatRows gets the allowFormatRows property value. Represents the worksheet protection option of allowing formatting rows.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowFormatRows()(*bool) {
    val, err := m.GetBackingStore().Get("allowFormatRows")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowInsertColumns gets the allowInsertColumns property value. Represents the worksheet protection option of allowing inserting columns.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowInsertColumns()(*bool) {
    val, err := m.GetBackingStore().Get("allowInsertColumns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowInsertHyperlinks gets the allowInsertHyperlinks property value. Represents the worksheet protection option of allowing inserting hyperlinks.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowInsertHyperlinks()(*bool) {
    val, err := m.GetBackingStore().Get("allowInsertHyperlinks")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowInsertRows gets the allowInsertRows property value. Represents the worksheet protection option of allowing inserting rows.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowInsertRows()(*bool) {
    val, err := m.GetBackingStore().Get("allowInsertRows")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowPivotTables gets the allowPivotTables property value. Represents the worksheet protection option of allowing using pivot table feature.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowPivotTables()(*bool) {
    val, err := m.GetBackingStore().Get("allowPivotTables")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetAllowSort gets the allowSort property value. Represents the worksheet protection option of allowing using sort feature.
// returns a *bool when successful
func (m *WorkbookWorksheetProtectionOptions) GetAllowSort()(*bool) {
    val, err := m.GetBackingStore().Get("allowSort")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *WorkbookWorksheetProtectionOptions) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookWorksheetProtectionOptions) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["allowAutoFilter"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowAutoFilter(val)
        }
        return nil
    }
    res["allowDeleteColumns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowDeleteColumns(val)
        }
        return nil
    }
    res["allowDeleteRows"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowDeleteRows(val)
        }
        return nil
    }
    res["allowFormatCells"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowFormatCells(val)
        }
        return nil
    }
    res["allowFormatColumns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowFormatColumns(val)
        }
        return nil
    }
    res["allowFormatRows"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowFormatRows(val)
        }
        return nil
    }
    res["allowInsertColumns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowInsertColumns(val)
        }
        return nil
    }
    res["allowInsertHyperlinks"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowInsertHyperlinks(val)
        }
        return nil
    }
    res["allowInsertRows"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowInsertRows(val)
        }
        return nil
    }
    res["allowPivotTables"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowPivotTables(val)
        }
        return nil
    }
    res["allowSort"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowSort(val)
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
func (m *WorkbookWorksheetProtectionOptions) GetOdataType()(*string) {
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
func (m *WorkbookWorksheetProtectionOptions) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("allowAutoFilter", m.GetAllowAutoFilter())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowDeleteColumns", m.GetAllowDeleteColumns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowDeleteRows", m.GetAllowDeleteRows())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowFormatCells", m.GetAllowFormatCells())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowFormatColumns", m.GetAllowFormatColumns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowFormatRows", m.GetAllowFormatRows())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowInsertColumns", m.GetAllowInsertColumns())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowInsertHyperlinks", m.GetAllowInsertHyperlinks())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowInsertRows", m.GetAllowInsertRows())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowPivotTables", m.GetAllowPivotTables())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("allowSort", m.GetAllowSort())
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
func (m *WorkbookWorksheetProtectionOptions) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowAutoFilter sets the allowAutoFilter property value. Represents the worksheet protection option of allowing using auto filter feature.
func (m *WorkbookWorksheetProtectionOptions) SetAllowAutoFilter(value *bool)() {
    err := m.GetBackingStore().Set("allowAutoFilter", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowDeleteColumns sets the allowDeleteColumns property value. Represents the worksheet protection option of allowing deleting columns.
func (m *WorkbookWorksheetProtectionOptions) SetAllowDeleteColumns(value *bool)() {
    err := m.GetBackingStore().Set("allowDeleteColumns", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowDeleteRows sets the allowDeleteRows property value. Represents the worksheet protection option of allowing deleting rows.
func (m *WorkbookWorksheetProtectionOptions) SetAllowDeleteRows(value *bool)() {
    err := m.GetBackingStore().Set("allowDeleteRows", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowFormatCells sets the allowFormatCells property value. Represents the worksheet protection option of allowing formatting cells.
func (m *WorkbookWorksheetProtectionOptions) SetAllowFormatCells(value *bool)() {
    err := m.GetBackingStore().Set("allowFormatCells", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowFormatColumns sets the allowFormatColumns property value. Represents the worksheet protection option of allowing formatting columns.
func (m *WorkbookWorksheetProtectionOptions) SetAllowFormatColumns(value *bool)() {
    err := m.GetBackingStore().Set("allowFormatColumns", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowFormatRows sets the allowFormatRows property value. Represents the worksheet protection option of allowing formatting rows.
func (m *WorkbookWorksheetProtectionOptions) SetAllowFormatRows(value *bool)() {
    err := m.GetBackingStore().Set("allowFormatRows", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowInsertColumns sets the allowInsertColumns property value. Represents the worksheet protection option of allowing inserting columns.
func (m *WorkbookWorksheetProtectionOptions) SetAllowInsertColumns(value *bool)() {
    err := m.GetBackingStore().Set("allowInsertColumns", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowInsertHyperlinks sets the allowInsertHyperlinks property value. Represents the worksheet protection option of allowing inserting hyperlinks.
func (m *WorkbookWorksheetProtectionOptions) SetAllowInsertHyperlinks(value *bool)() {
    err := m.GetBackingStore().Set("allowInsertHyperlinks", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowInsertRows sets the allowInsertRows property value. Represents the worksheet protection option of allowing inserting rows.
func (m *WorkbookWorksheetProtectionOptions) SetAllowInsertRows(value *bool)() {
    err := m.GetBackingStore().Set("allowInsertRows", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowPivotTables sets the allowPivotTables property value. Represents the worksheet protection option of allowing using pivot table feature.
func (m *WorkbookWorksheetProtectionOptions) SetAllowPivotTables(value *bool)() {
    err := m.GetBackingStore().Set("allowPivotTables", value)
    if err != nil {
        panic(err)
    }
}
// SetAllowSort sets the allowSort property value. Represents the worksheet protection option of allowing using sort feature.
func (m *WorkbookWorksheetProtectionOptions) SetAllowSort(value *bool)() {
    err := m.GetBackingStore().Set("allowSort", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *WorkbookWorksheetProtectionOptions) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *WorkbookWorksheetProtectionOptions) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookWorksheetProtectionOptionsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowAutoFilter()(*bool)
    GetAllowDeleteColumns()(*bool)
    GetAllowDeleteRows()(*bool)
    GetAllowFormatCells()(*bool)
    GetAllowFormatColumns()(*bool)
    GetAllowFormatRows()(*bool)
    GetAllowInsertColumns()(*bool)
    GetAllowInsertHyperlinks()(*bool)
    GetAllowInsertRows()(*bool)
    GetAllowPivotTables()(*bool)
    GetAllowSort()(*bool)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetOdataType()(*string)
    SetAllowAutoFilter(value *bool)()
    SetAllowDeleteColumns(value *bool)()
    SetAllowDeleteRows(value *bool)()
    SetAllowFormatCells(value *bool)()
    SetAllowFormatColumns(value *bool)()
    SetAllowFormatRows(value *bool)()
    SetAllowInsertColumns(value *bool)()
    SetAllowInsertHyperlinks(value *bool)()
    SetAllowInsertRows(value *bool)()
    SetAllowPivotTables(value *bool)()
    SetAllowSort(value *bool)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetOdataType(value *string)()
}

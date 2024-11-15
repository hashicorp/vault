package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookTable struct {
    Entity
}
// NewWorkbookTable instantiates a new WorkbookTable and sets the default values.
func NewWorkbookTable()(*WorkbookTable) {
    m := &WorkbookTable{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookTableFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookTableFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookTable(), nil
}
// GetColumns gets the columns property value. The list of all the columns in the table. Read-only.
// returns a []WorkbookTableColumnable when successful
func (m *WorkbookTable) GetColumns()([]WorkbookTableColumnable) {
    val, err := m.GetBackingStore().Get("columns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookTableColumnable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookTable) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["columns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookTableColumnFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookTableColumnable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookTableColumnable)
                }
            }
            m.SetColumns(res)
        }
        return nil
    }
    res["highlightFirstColumn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHighlightFirstColumn(val)
        }
        return nil
    }
    res["highlightLastColumn"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHighlightLastColumn(val)
        }
        return nil
    }
    res["legacyId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLegacyId(val)
        }
        return nil
    }
    res["name"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetName(val)
        }
        return nil
    }
    res["rows"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookTableRowFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookTableRowable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookTableRowable)
                }
            }
            m.SetRows(res)
        }
        return nil
    }
    res["showBandedColumns"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowBandedColumns(val)
        }
        return nil
    }
    res["showBandedRows"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowBandedRows(val)
        }
        return nil
    }
    res["showFilterButton"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowFilterButton(val)
        }
        return nil
    }
    res["showHeaders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowHeaders(val)
        }
        return nil
    }
    res["showTotals"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowTotals(val)
        }
        return nil
    }
    res["sort"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookTableSortFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSort(val.(WorkbookTableSortable))
        }
        return nil
    }
    res["style"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStyle(val)
        }
        return nil
    }
    res["worksheet"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookWorksheetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWorksheet(val.(WorkbookWorksheetable))
        }
        return nil
    }
    return res
}
// GetHighlightFirstColumn gets the highlightFirstColumn property value. Indicates whether the first column contains special formatting.
// returns a *bool when successful
func (m *WorkbookTable) GetHighlightFirstColumn()(*bool) {
    val, err := m.GetBackingStore().Get("highlightFirstColumn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHighlightLastColumn gets the highlightLastColumn property value. Indicates whether the last column contains special formatting.
// returns a *bool when successful
func (m *WorkbookTable) GetHighlightLastColumn()(*bool) {
    val, err := m.GetBackingStore().Get("highlightLastColumn")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLegacyId gets the legacyId property value. A legacy identifier used in older Excel clients. The value of the identifier remains the same even when the table is renamed. This property should be interpreted as an opaque string value and shouldn't be parsed to any other type. Read-only.
// returns a *string when successful
func (m *WorkbookTable) GetLegacyId()(*string) {
    val, err := m.GetBackingStore().Get("legacyId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetName gets the name property value. The name of the table.
// returns a *string when successful
func (m *WorkbookTable) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRows gets the rows property value. The list of all the rows in the table. Read-only.
// returns a []WorkbookTableRowable when successful
func (m *WorkbookTable) GetRows()([]WorkbookTableRowable) {
    val, err := m.GetBackingStore().Get("rows")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookTableRowable)
    }
    return nil
}
// GetShowBandedColumns gets the showBandedColumns property value. Indicates whether the columns show banded formatting in which odd columns are highlighted differently from even ones to make reading the table easier.
// returns a *bool when successful
func (m *WorkbookTable) GetShowBandedColumns()(*bool) {
    val, err := m.GetBackingStore().Get("showBandedColumns")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowBandedRows gets the showBandedRows property value. Indicates whether the rows show banded formatting in which odd rows are highlighted differently from even ones to make reading the table easier.
// returns a *bool when successful
func (m *WorkbookTable) GetShowBandedRows()(*bool) {
    val, err := m.GetBackingStore().Get("showBandedRows")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowFilterButton gets the showFilterButton property value. Indicates whether the filter buttons are visible at the top of each column header. Setting this is only allowed if the table contains a header row.
// returns a *bool when successful
func (m *WorkbookTable) GetShowFilterButton()(*bool) {
    val, err := m.GetBackingStore().Get("showFilterButton")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowHeaders gets the showHeaders property value. Indicates whether the header row is visible or not. This value can be set to show or remove the header row.
// returns a *bool when successful
func (m *WorkbookTable) GetShowHeaders()(*bool) {
    val, err := m.GetBackingStore().Get("showHeaders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowTotals gets the showTotals property value. Indicates whether the total row is visible or not. This value can be set to show or remove the total row.
// returns a *bool when successful
func (m *WorkbookTable) GetShowTotals()(*bool) {
    val, err := m.GetBackingStore().Get("showTotals")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSort gets the sort property value. The sorting for the table. Read-only.
// returns a WorkbookTableSortable when successful
func (m *WorkbookTable) GetSort()(WorkbookTableSortable) {
    val, err := m.GetBackingStore().Get("sort")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookTableSortable)
    }
    return nil
}
// GetStyle gets the style property value. A constant value that represents the Table style. Possible values are: TableStyleLight1 through TableStyleLight21, TableStyleMedium1 through TableStyleMedium28, TableStyleStyleDark1 through TableStyleStyleDark11. A custom user-defined style present in the workbook can also be specified.
// returns a *string when successful
func (m *WorkbookTable) GetStyle()(*string) {
    val, err := m.GetBackingStore().Get("style")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWorksheet gets the worksheet property value. The worksheet containing the current table. Read-only.
// returns a WorkbookWorksheetable when successful
func (m *WorkbookTable) GetWorksheet()(WorkbookWorksheetable) {
    val, err := m.GetBackingStore().Get("worksheet")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookWorksheetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookTable) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetColumns() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetColumns()))
        for i, v := range m.GetColumns() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("columns", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("highlightFirstColumn", m.GetHighlightFirstColumn())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("highlightLastColumn", m.GetHighlightLastColumn())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("legacyId", m.GetLegacyId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    if m.GetRows() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRows()))
        for i, v := range m.GetRows() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("rows", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showBandedColumns", m.GetShowBandedColumns())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showBandedRows", m.GetShowBandedRows())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showFilterButton", m.GetShowFilterButton())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showHeaders", m.GetShowHeaders())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showTotals", m.GetShowTotals())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sort", m.GetSort())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("style", m.GetStyle())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("worksheet", m.GetWorksheet())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetColumns sets the columns property value. The list of all the columns in the table. Read-only.
func (m *WorkbookTable) SetColumns(value []WorkbookTableColumnable)() {
    err := m.GetBackingStore().Set("columns", value)
    if err != nil {
        panic(err)
    }
}
// SetHighlightFirstColumn sets the highlightFirstColumn property value. Indicates whether the first column contains special formatting.
func (m *WorkbookTable) SetHighlightFirstColumn(value *bool)() {
    err := m.GetBackingStore().Set("highlightFirstColumn", value)
    if err != nil {
        panic(err)
    }
}
// SetHighlightLastColumn sets the highlightLastColumn property value. Indicates whether the last column contains special formatting.
func (m *WorkbookTable) SetHighlightLastColumn(value *bool)() {
    err := m.GetBackingStore().Set("highlightLastColumn", value)
    if err != nil {
        panic(err)
    }
}
// SetLegacyId sets the legacyId property value. A legacy identifier used in older Excel clients. The value of the identifier remains the same even when the table is renamed. This property should be interpreted as an opaque string value and shouldn't be parsed to any other type. Read-only.
func (m *WorkbookTable) SetLegacyId(value *string)() {
    err := m.GetBackingStore().Set("legacyId", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The name of the table.
func (m *WorkbookTable) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetRows sets the rows property value. The list of all the rows in the table. Read-only.
func (m *WorkbookTable) SetRows(value []WorkbookTableRowable)() {
    err := m.GetBackingStore().Set("rows", value)
    if err != nil {
        panic(err)
    }
}
// SetShowBandedColumns sets the showBandedColumns property value. Indicates whether the columns show banded formatting in which odd columns are highlighted differently from even ones to make reading the table easier.
func (m *WorkbookTable) SetShowBandedColumns(value *bool)() {
    err := m.GetBackingStore().Set("showBandedColumns", value)
    if err != nil {
        panic(err)
    }
}
// SetShowBandedRows sets the showBandedRows property value. Indicates whether the rows show banded formatting in which odd rows are highlighted differently from even ones to make reading the table easier.
func (m *WorkbookTable) SetShowBandedRows(value *bool)() {
    err := m.GetBackingStore().Set("showBandedRows", value)
    if err != nil {
        panic(err)
    }
}
// SetShowFilterButton sets the showFilterButton property value. Indicates whether the filter buttons are visible at the top of each column header. Setting this is only allowed if the table contains a header row.
func (m *WorkbookTable) SetShowFilterButton(value *bool)() {
    err := m.GetBackingStore().Set("showFilterButton", value)
    if err != nil {
        panic(err)
    }
}
// SetShowHeaders sets the showHeaders property value. Indicates whether the header row is visible or not. This value can be set to show or remove the header row.
func (m *WorkbookTable) SetShowHeaders(value *bool)() {
    err := m.GetBackingStore().Set("showHeaders", value)
    if err != nil {
        panic(err)
    }
}
// SetShowTotals sets the showTotals property value. Indicates whether the total row is visible or not. This value can be set to show or remove the total row.
func (m *WorkbookTable) SetShowTotals(value *bool)() {
    err := m.GetBackingStore().Set("showTotals", value)
    if err != nil {
        panic(err)
    }
}
// SetSort sets the sort property value. The sorting for the table. Read-only.
func (m *WorkbookTable) SetSort(value WorkbookTableSortable)() {
    err := m.GetBackingStore().Set("sort", value)
    if err != nil {
        panic(err)
    }
}
// SetStyle sets the style property value. A constant value that represents the Table style. Possible values are: TableStyleLight1 through TableStyleLight21, TableStyleMedium1 through TableStyleMedium28, TableStyleStyleDark1 through TableStyleStyleDark11. A custom user-defined style present in the workbook can also be specified.
func (m *WorkbookTable) SetStyle(value *string)() {
    err := m.GetBackingStore().Set("style", value)
    if err != nil {
        panic(err)
    }
}
// SetWorksheet sets the worksheet property value. The worksheet containing the current table. Read-only.
func (m *WorkbookTable) SetWorksheet(value WorkbookWorksheetable)() {
    err := m.GetBackingStore().Set("worksheet", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookTableable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetColumns()([]WorkbookTableColumnable)
    GetHighlightFirstColumn()(*bool)
    GetHighlightLastColumn()(*bool)
    GetLegacyId()(*string)
    GetName()(*string)
    GetRows()([]WorkbookTableRowable)
    GetShowBandedColumns()(*bool)
    GetShowBandedRows()(*bool)
    GetShowFilterButton()(*bool)
    GetShowHeaders()(*bool)
    GetShowTotals()(*bool)
    GetSort()(WorkbookTableSortable)
    GetStyle()(*string)
    GetWorksheet()(WorkbookWorksheetable)
    SetColumns(value []WorkbookTableColumnable)()
    SetHighlightFirstColumn(value *bool)()
    SetHighlightLastColumn(value *bool)()
    SetLegacyId(value *string)()
    SetName(value *string)()
    SetRows(value []WorkbookTableRowable)()
    SetShowBandedColumns(value *bool)()
    SetShowBandedRows(value *bool)()
    SetShowFilterButton(value *bool)()
    SetShowHeaders(value *bool)()
    SetShowTotals(value *bool)()
    SetSort(value WorkbookTableSortable)()
    SetStyle(value *string)()
    SetWorksheet(value WorkbookWorksheetable)()
}

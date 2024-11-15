package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookWorksheet struct {
    Entity
}
// NewWorkbookWorksheet instantiates a new WorkbookWorksheet and sets the default values.
func NewWorkbookWorksheet()(*WorkbookWorksheet) {
    m := &WorkbookWorksheet{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookWorksheetFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookWorksheetFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookWorksheet(), nil
}
// GetCharts gets the charts property value. The list of charts that are part of the worksheet. Read-only.
// returns a []WorkbookChartable when successful
func (m *WorkbookWorksheet) GetCharts()([]WorkbookChartable) {
    val, err := m.GetBackingStore().Get("charts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookChartable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookWorksheet) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["charts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookChartFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookChartable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookChartable)
                }
            }
            m.SetCharts(res)
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
    res["names"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookNamedItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookNamedItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookNamedItemable)
                }
            }
            m.SetNames(res)
        }
        return nil
    }
    res["pivotTables"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookPivotTableFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookPivotTableable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookPivotTableable)
                }
            }
            m.SetPivotTables(res)
        }
        return nil
    }
    res["position"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPosition(val)
        }
        return nil
    }
    res["protection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookWorksheetProtectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProtection(val.(WorkbookWorksheetProtectionable))
        }
        return nil
    }
    res["tables"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookTableFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookTableable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookTableable)
                }
            }
            m.SetTables(res)
        }
        return nil
    }
    res["visibility"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVisibility(val)
        }
        return nil
    }
    return res
}
// GetName gets the name property value. The display name of the worksheet.
// returns a *string when successful
func (m *WorkbookWorksheet) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNames gets the names property value. The list of names that are associated with the worksheet. Read-only.
// returns a []WorkbookNamedItemable when successful
func (m *WorkbookWorksheet) GetNames()([]WorkbookNamedItemable) {
    val, err := m.GetBackingStore().Get("names")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookNamedItemable)
    }
    return nil
}
// GetPivotTables gets the pivotTables property value. The list of piot tables that are part of the worksheet.
// returns a []WorkbookPivotTableable when successful
func (m *WorkbookWorksheet) GetPivotTables()([]WorkbookPivotTableable) {
    val, err := m.GetBackingStore().Get("pivotTables")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookPivotTableable)
    }
    return nil
}
// GetPosition gets the position property value. The zero-based position of the worksheet within the workbook.
// returns a *int32 when successful
func (m *WorkbookWorksheet) GetPosition()(*int32) {
    val, err := m.GetBackingStore().Get("position")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetProtection gets the protection property value. The sheet protection object for a worksheet. Read-only.
// returns a WorkbookWorksheetProtectionable when successful
func (m *WorkbookWorksheet) GetProtection()(WorkbookWorksheetProtectionable) {
    val, err := m.GetBackingStore().Get("protection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookWorksheetProtectionable)
    }
    return nil
}
// GetTables gets the tables property value. The list of tables that are part of the worksheet. Read-only.
// returns a []WorkbookTableable when successful
func (m *WorkbookWorksheet) GetTables()([]WorkbookTableable) {
    val, err := m.GetBackingStore().Get("tables")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookTableable)
    }
    return nil
}
// GetVisibility gets the visibility property value. The visibility of the worksheet. The possible values are: Visible, Hidden, VeryHidden.
// returns a *string when successful
func (m *WorkbookWorksheet) GetVisibility()(*string) {
    val, err := m.GetBackingStore().Get("visibility")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookWorksheet) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCharts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCharts()))
        for i, v := range m.GetCharts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("charts", cast)
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
    if m.GetNames() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetNames()))
        for i, v := range m.GetNames() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("names", cast)
        if err != nil {
            return err
        }
    }
    if m.GetPivotTables() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPivotTables()))
        for i, v := range m.GetPivotTables() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("pivotTables", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("position", m.GetPosition())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("protection", m.GetProtection())
        if err != nil {
            return err
        }
    }
    if m.GetTables() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTables()))
        for i, v := range m.GetTables() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("tables", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("visibility", m.GetVisibility())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCharts sets the charts property value. The list of charts that are part of the worksheet. Read-only.
func (m *WorkbookWorksheet) SetCharts(value []WorkbookChartable)() {
    err := m.GetBackingStore().Set("charts", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The display name of the worksheet.
func (m *WorkbookWorksheet) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetNames sets the names property value. The list of names that are associated with the worksheet. Read-only.
func (m *WorkbookWorksheet) SetNames(value []WorkbookNamedItemable)() {
    err := m.GetBackingStore().Set("names", value)
    if err != nil {
        panic(err)
    }
}
// SetPivotTables sets the pivotTables property value. The list of piot tables that are part of the worksheet.
func (m *WorkbookWorksheet) SetPivotTables(value []WorkbookPivotTableable)() {
    err := m.GetBackingStore().Set("pivotTables", value)
    if err != nil {
        panic(err)
    }
}
// SetPosition sets the position property value. The zero-based position of the worksheet within the workbook.
func (m *WorkbookWorksheet) SetPosition(value *int32)() {
    err := m.GetBackingStore().Set("position", value)
    if err != nil {
        panic(err)
    }
}
// SetProtection sets the protection property value. The sheet protection object for a worksheet. Read-only.
func (m *WorkbookWorksheet) SetProtection(value WorkbookWorksheetProtectionable)() {
    err := m.GetBackingStore().Set("protection", value)
    if err != nil {
        panic(err)
    }
}
// SetTables sets the tables property value. The list of tables that are part of the worksheet. Read-only.
func (m *WorkbookWorksheet) SetTables(value []WorkbookTableable)() {
    err := m.GetBackingStore().Set("tables", value)
    if err != nil {
        panic(err)
    }
}
// SetVisibility sets the visibility property value. The visibility of the worksheet. The possible values are: Visible, Hidden, VeryHidden.
func (m *WorkbookWorksheet) SetVisibility(value *string)() {
    err := m.GetBackingStore().Set("visibility", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookWorksheetable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCharts()([]WorkbookChartable)
    GetName()(*string)
    GetNames()([]WorkbookNamedItemable)
    GetPivotTables()([]WorkbookPivotTableable)
    GetPosition()(*int32)
    GetProtection()(WorkbookWorksheetProtectionable)
    GetTables()([]WorkbookTableable)
    GetVisibility()(*string)
    SetCharts(value []WorkbookChartable)()
    SetName(value *string)()
    SetNames(value []WorkbookNamedItemable)()
    SetPivotTables(value []WorkbookPivotTableable)()
    SetPosition(value *int32)()
    SetProtection(value WorkbookWorksheetProtectionable)()
    SetTables(value []WorkbookTableable)()
    SetVisibility(value *string)()
}

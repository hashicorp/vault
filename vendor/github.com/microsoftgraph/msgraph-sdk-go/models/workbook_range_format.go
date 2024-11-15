package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookRangeFormat struct {
    Entity
}
// NewWorkbookRangeFormat instantiates a new WorkbookRangeFormat and sets the default values.
func NewWorkbookRangeFormat()(*WorkbookRangeFormat) {
    m := &WorkbookRangeFormat{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookRangeFormatFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookRangeFormatFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookRangeFormat(), nil
}
// GetBorders gets the borders property value. Collection of border objects that apply to the overall range selected Read-only.
// returns a []WorkbookRangeBorderable when successful
func (m *WorkbookRangeFormat) GetBorders()([]WorkbookRangeBorderable) {
    val, err := m.GetBackingStore().Get("borders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookRangeBorderable)
    }
    return nil
}
// GetColumnWidth gets the columnWidth property value. The width of all columns within the range. If the column widths aren't uniform, null will be returned.
// returns a *float64 when successful
func (m *WorkbookRangeFormat) GetColumnWidth()(*float64) {
    val, err := m.GetBackingStore().Get("columnWidth")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookRangeFormat) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["borders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookRangeBorderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookRangeBorderable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookRangeBorderable)
                }
            }
            m.SetBorders(res)
        }
        return nil
    }
    res["columnWidth"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetColumnWidth(val)
        }
        return nil
    }
    res["fill"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookRangeFillFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFill(val.(WorkbookRangeFillable))
        }
        return nil
    }
    res["font"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookRangeFontFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFont(val.(WorkbookRangeFontable))
        }
        return nil
    }
    res["horizontalAlignment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHorizontalAlignment(val)
        }
        return nil
    }
    res["protection"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookFormatProtectionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProtection(val.(WorkbookFormatProtectionable))
        }
        return nil
    }
    res["rowHeight"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRowHeight(val)
        }
        return nil
    }
    res["verticalAlignment"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVerticalAlignment(val)
        }
        return nil
    }
    res["wrapText"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWrapText(val)
        }
        return nil
    }
    return res
}
// GetFill gets the fill property value. Returns the fill object defined on the overall range. Read-only.
// returns a WorkbookRangeFillable when successful
func (m *WorkbookRangeFormat) GetFill()(WorkbookRangeFillable) {
    val, err := m.GetBackingStore().Get("fill")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookRangeFillable)
    }
    return nil
}
// GetFont gets the font property value. Returns the font object defined on the overall range selected Read-only.
// returns a WorkbookRangeFontable when successful
func (m *WorkbookRangeFormat) GetFont()(WorkbookRangeFontable) {
    val, err := m.GetBackingStore().Get("font")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookRangeFontable)
    }
    return nil
}
// GetHorizontalAlignment gets the horizontalAlignment property value. The horizontal alignment for the specified object. Possible values are: General, Left, Center, Right, Fill, Justify, CenterAcrossSelection, Distributed.
// returns a *string when successful
func (m *WorkbookRangeFormat) GetHorizontalAlignment()(*string) {
    val, err := m.GetBackingStore().Get("horizontalAlignment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetProtection gets the protection property value. Returns the format protection object for a range. Read-only.
// returns a WorkbookFormatProtectionable when successful
func (m *WorkbookRangeFormat) GetProtection()(WorkbookFormatProtectionable) {
    val, err := m.GetBackingStore().Get("protection")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookFormatProtectionable)
    }
    return nil
}
// GetRowHeight gets the rowHeight property value. The height of all rows in the range. If the row heights aren't uniform null will be returned.
// returns a *float64 when successful
func (m *WorkbookRangeFormat) GetRowHeight()(*float64) {
    val, err := m.GetBackingStore().Get("rowHeight")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetVerticalAlignment gets the verticalAlignment property value. The vertical alignment for the specified object. Possible values are: Top, Center, Bottom, Justify, Distributed.
// returns a *string when successful
func (m *WorkbookRangeFormat) GetVerticalAlignment()(*string) {
    val, err := m.GetBackingStore().Get("verticalAlignment")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWrapText gets the wrapText property value. Indicates whether Excel wraps the text in the object. A null value indicates that the entire range doesn't have a uniform wrap setting.
// returns a *bool when successful
func (m *WorkbookRangeFormat) GetWrapText()(*bool) {
    val, err := m.GetBackingStore().Get("wrapText")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookRangeFormat) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetBorders() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetBorders()))
        for i, v := range m.GetBorders() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("borders", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("columnWidth", m.GetColumnWidth())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("fill", m.GetFill())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("font", m.GetFont())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("horizontalAlignment", m.GetHorizontalAlignment())
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
    {
        err = writer.WriteFloat64Value("rowHeight", m.GetRowHeight())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("verticalAlignment", m.GetVerticalAlignment())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("wrapText", m.GetWrapText())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBorders sets the borders property value. Collection of border objects that apply to the overall range selected Read-only.
func (m *WorkbookRangeFormat) SetBorders(value []WorkbookRangeBorderable)() {
    err := m.GetBackingStore().Set("borders", value)
    if err != nil {
        panic(err)
    }
}
// SetColumnWidth sets the columnWidth property value. The width of all columns within the range. If the column widths aren't uniform, null will be returned.
func (m *WorkbookRangeFormat) SetColumnWidth(value *float64)() {
    err := m.GetBackingStore().Set("columnWidth", value)
    if err != nil {
        panic(err)
    }
}
// SetFill sets the fill property value. Returns the fill object defined on the overall range. Read-only.
func (m *WorkbookRangeFormat) SetFill(value WorkbookRangeFillable)() {
    err := m.GetBackingStore().Set("fill", value)
    if err != nil {
        panic(err)
    }
}
// SetFont sets the font property value. Returns the font object defined on the overall range selected Read-only.
func (m *WorkbookRangeFormat) SetFont(value WorkbookRangeFontable)() {
    err := m.GetBackingStore().Set("font", value)
    if err != nil {
        panic(err)
    }
}
// SetHorizontalAlignment sets the horizontalAlignment property value. The horizontal alignment for the specified object. Possible values are: General, Left, Center, Right, Fill, Justify, CenterAcrossSelection, Distributed.
func (m *WorkbookRangeFormat) SetHorizontalAlignment(value *string)() {
    err := m.GetBackingStore().Set("horizontalAlignment", value)
    if err != nil {
        panic(err)
    }
}
// SetProtection sets the protection property value. Returns the format protection object for a range. Read-only.
func (m *WorkbookRangeFormat) SetProtection(value WorkbookFormatProtectionable)() {
    err := m.GetBackingStore().Set("protection", value)
    if err != nil {
        panic(err)
    }
}
// SetRowHeight sets the rowHeight property value. The height of all rows in the range. If the row heights aren't uniform null will be returned.
func (m *WorkbookRangeFormat) SetRowHeight(value *float64)() {
    err := m.GetBackingStore().Set("rowHeight", value)
    if err != nil {
        panic(err)
    }
}
// SetVerticalAlignment sets the verticalAlignment property value. The vertical alignment for the specified object. Possible values are: Top, Center, Bottom, Justify, Distributed.
func (m *WorkbookRangeFormat) SetVerticalAlignment(value *string)() {
    err := m.GetBackingStore().Set("verticalAlignment", value)
    if err != nil {
        panic(err)
    }
}
// SetWrapText sets the wrapText property value. Indicates whether Excel wraps the text in the object. A null value indicates that the entire range doesn't have a uniform wrap setting.
func (m *WorkbookRangeFormat) SetWrapText(value *bool)() {
    err := m.GetBackingStore().Set("wrapText", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookRangeFormatable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBorders()([]WorkbookRangeBorderable)
    GetColumnWidth()(*float64)
    GetFill()(WorkbookRangeFillable)
    GetFont()(WorkbookRangeFontable)
    GetHorizontalAlignment()(*string)
    GetProtection()(WorkbookFormatProtectionable)
    GetRowHeight()(*float64)
    GetVerticalAlignment()(*string)
    GetWrapText()(*bool)
    SetBorders(value []WorkbookRangeBorderable)()
    SetColumnWidth(value *float64)()
    SetFill(value WorkbookRangeFillable)()
    SetFont(value WorkbookRangeFontable)()
    SetHorizontalAlignment(value *string)()
    SetProtection(value WorkbookFormatProtectionable)()
    SetRowHeight(value *float64)()
    SetVerticalAlignment(value *string)()
    SetWrapText(value *bool)()
}

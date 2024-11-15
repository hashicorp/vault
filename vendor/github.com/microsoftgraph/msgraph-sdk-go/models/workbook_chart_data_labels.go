package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookChartDataLabels struct {
    Entity
}
// NewWorkbookChartDataLabels instantiates a new WorkbookChartDataLabels and sets the default values.
func NewWorkbookChartDataLabels()(*WorkbookChartDataLabels) {
    m := &WorkbookChartDataLabels{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartDataLabelsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookChartDataLabelsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartDataLabels(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookChartDataLabels) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["format"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartDataLabelFormatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFormat(val.(WorkbookChartDataLabelFormatable))
        }
        return nil
    }
    res["position"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPosition(val)
        }
        return nil
    }
    res["separator"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSeparator(val)
        }
        return nil
    }
    res["showBubbleSize"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowBubbleSize(val)
        }
        return nil
    }
    res["showCategoryName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowCategoryName(val)
        }
        return nil
    }
    res["showLegendKey"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowLegendKey(val)
        }
        return nil
    }
    res["showPercentage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowPercentage(val)
        }
        return nil
    }
    res["showSeriesName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowSeriesName(val)
        }
        return nil
    }
    res["showValue"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetShowValue(val)
        }
        return nil
    }
    return res
}
// GetFormat gets the format property value. Represents the format of chart data labels, which includes fill and font formatting. Read-only.
// returns a WorkbookChartDataLabelFormatable when successful
func (m *WorkbookChartDataLabels) GetFormat()(WorkbookChartDataLabelFormatable) {
    val, err := m.GetBackingStore().Get("format")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartDataLabelFormatable)
    }
    return nil
}
// GetPosition gets the position property value. DataLabelPosition value that represents the position of the data label. The possible values are: None, Center, InsideEnd, InsideBase, OutsideEnd, Left, Right, Top, Bottom, BestFit, Callout.
// returns a *string when successful
func (m *WorkbookChartDataLabels) GetPosition()(*string) {
    val, err := m.GetBackingStore().Get("position")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSeparator gets the separator property value. String that represents the separator used for the data labels on a chart.
// returns a *string when successful
func (m *WorkbookChartDataLabels) GetSeparator()(*string) {
    val, err := m.GetBackingStore().Get("separator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetShowBubbleSize gets the showBubbleSize property value. Boolean value that represents whether the data label bubble size is visible.
// returns a *bool when successful
func (m *WorkbookChartDataLabels) GetShowBubbleSize()(*bool) {
    val, err := m.GetBackingStore().Get("showBubbleSize")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowCategoryName gets the showCategoryName property value. Boolean value that represents whether the data label category name is visible.
// returns a *bool when successful
func (m *WorkbookChartDataLabels) GetShowCategoryName()(*bool) {
    val, err := m.GetBackingStore().Get("showCategoryName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowLegendKey gets the showLegendKey property value. Boolean value that represents whether the data label legend key is visible.
// returns a *bool when successful
func (m *WorkbookChartDataLabels) GetShowLegendKey()(*bool) {
    val, err := m.GetBackingStore().Get("showLegendKey")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowPercentage gets the showPercentage property value. Boolean value that represents whether the data label percentage is visible.
// returns a *bool when successful
func (m *WorkbookChartDataLabels) GetShowPercentage()(*bool) {
    val, err := m.GetBackingStore().Get("showPercentage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowSeriesName gets the showSeriesName property value. Boolean value that represents whether the data label series name is visible.
// returns a *bool when successful
func (m *WorkbookChartDataLabels) GetShowSeriesName()(*bool) {
    val, err := m.GetBackingStore().Get("showSeriesName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetShowValue gets the showValue property value. Boolean value that represents whether the data label value is visible.
// returns a *bool when successful
func (m *WorkbookChartDataLabels) GetShowValue()(*bool) {
    val, err := m.GetBackingStore().Get("showValue")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookChartDataLabels) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("format", m.GetFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("position", m.GetPosition())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("separator", m.GetSeparator())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showBubbleSize", m.GetShowBubbleSize())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showCategoryName", m.GetShowCategoryName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showLegendKey", m.GetShowLegendKey())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showPercentage", m.GetShowPercentage())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showSeriesName", m.GetShowSeriesName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("showValue", m.GetShowValue())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFormat sets the format property value. Represents the format of chart data labels, which includes fill and font formatting. Read-only.
func (m *WorkbookChartDataLabels) SetFormat(value WorkbookChartDataLabelFormatable)() {
    err := m.GetBackingStore().Set("format", value)
    if err != nil {
        panic(err)
    }
}
// SetPosition sets the position property value. DataLabelPosition value that represents the position of the data label. The possible values are: None, Center, InsideEnd, InsideBase, OutsideEnd, Left, Right, Top, Bottom, BestFit, Callout.
func (m *WorkbookChartDataLabels) SetPosition(value *string)() {
    err := m.GetBackingStore().Set("position", value)
    if err != nil {
        panic(err)
    }
}
// SetSeparator sets the separator property value. String that represents the separator used for the data labels on a chart.
func (m *WorkbookChartDataLabels) SetSeparator(value *string)() {
    err := m.GetBackingStore().Set("separator", value)
    if err != nil {
        panic(err)
    }
}
// SetShowBubbleSize sets the showBubbleSize property value. Boolean value that represents whether the data label bubble size is visible.
func (m *WorkbookChartDataLabels) SetShowBubbleSize(value *bool)() {
    err := m.GetBackingStore().Set("showBubbleSize", value)
    if err != nil {
        panic(err)
    }
}
// SetShowCategoryName sets the showCategoryName property value. Boolean value that represents whether the data label category name is visible.
func (m *WorkbookChartDataLabels) SetShowCategoryName(value *bool)() {
    err := m.GetBackingStore().Set("showCategoryName", value)
    if err != nil {
        panic(err)
    }
}
// SetShowLegendKey sets the showLegendKey property value. Boolean value that represents whether the data label legend key is visible.
func (m *WorkbookChartDataLabels) SetShowLegendKey(value *bool)() {
    err := m.GetBackingStore().Set("showLegendKey", value)
    if err != nil {
        panic(err)
    }
}
// SetShowPercentage sets the showPercentage property value. Boolean value that represents whether the data label percentage is visible.
func (m *WorkbookChartDataLabels) SetShowPercentage(value *bool)() {
    err := m.GetBackingStore().Set("showPercentage", value)
    if err != nil {
        panic(err)
    }
}
// SetShowSeriesName sets the showSeriesName property value. Boolean value that represents whether the data label series name is visible.
func (m *WorkbookChartDataLabels) SetShowSeriesName(value *bool)() {
    err := m.GetBackingStore().Set("showSeriesName", value)
    if err != nil {
        panic(err)
    }
}
// SetShowValue sets the showValue property value. Boolean value that represents whether the data label value is visible.
func (m *WorkbookChartDataLabels) SetShowValue(value *bool)() {
    err := m.GetBackingStore().Set("showValue", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookChartDataLabelsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFormat()(WorkbookChartDataLabelFormatable)
    GetPosition()(*string)
    GetSeparator()(*string)
    GetShowBubbleSize()(*bool)
    GetShowCategoryName()(*bool)
    GetShowLegendKey()(*bool)
    GetShowPercentage()(*bool)
    GetShowSeriesName()(*bool)
    GetShowValue()(*bool)
    SetFormat(value WorkbookChartDataLabelFormatable)()
    SetPosition(value *string)()
    SetSeparator(value *string)()
    SetShowBubbleSize(value *bool)()
    SetShowCategoryName(value *bool)()
    SetShowLegendKey(value *bool)()
    SetShowPercentage(value *bool)()
    SetShowSeriesName(value *bool)()
    SetShowValue(value *bool)()
}

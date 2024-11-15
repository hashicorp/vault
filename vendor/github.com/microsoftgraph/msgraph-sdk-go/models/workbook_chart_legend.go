package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookChartLegend struct {
    Entity
}
// NewWorkbookChartLegend instantiates a new WorkbookChartLegend and sets the default values.
func NewWorkbookChartLegend()(*WorkbookChartLegend) {
    m := &WorkbookChartLegend{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartLegendFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookChartLegendFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartLegend(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookChartLegend) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["format"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartLegendFormatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFormat(val.(WorkbookChartLegendFormatable))
        }
        return nil
    }
    res["overlay"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOverlay(val)
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
    res["visible"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVisible(val)
        }
        return nil
    }
    return res
}
// GetFormat gets the format property value. Represents the formatting of a chart legend, which includes fill and font formatting. Read-only.
// returns a WorkbookChartLegendFormatable when successful
func (m *WorkbookChartLegend) GetFormat()(WorkbookChartLegendFormatable) {
    val, err := m.GetBackingStore().Get("format")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartLegendFormatable)
    }
    return nil
}
// GetOverlay gets the overlay property value. Indicates whether the chart legend should overlap with the main body of the chart.
// returns a *bool when successful
func (m *WorkbookChartLegend) GetOverlay()(*bool) {
    val, err := m.GetBackingStore().Get("overlay")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetPosition gets the position property value. Represents the position of the legend on the chart. The possible values are: Top, Bottom, Left, Right, Corner, Custom.
// returns a *string when successful
func (m *WorkbookChartLegend) GetPosition()(*string) {
    val, err := m.GetBackingStore().Get("position")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetVisible gets the visible property value. Indicates whether the chart legend is visible.
// returns a *bool when successful
func (m *WorkbookChartLegend) GetVisible()(*bool) {
    val, err := m.GetBackingStore().Get("visible")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookChartLegend) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteBoolValue("overlay", m.GetOverlay())
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
        err = writer.WriteBoolValue("visible", m.GetVisible())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFormat sets the format property value. Represents the formatting of a chart legend, which includes fill and font formatting. Read-only.
func (m *WorkbookChartLegend) SetFormat(value WorkbookChartLegendFormatable)() {
    err := m.GetBackingStore().Set("format", value)
    if err != nil {
        panic(err)
    }
}
// SetOverlay sets the overlay property value. Indicates whether the chart legend should overlap with the main body of the chart.
func (m *WorkbookChartLegend) SetOverlay(value *bool)() {
    err := m.GetBackingStore().Set("overlay", value)
    if err != nil {
        panic(err)
    }
}
// SetPosition sets the position property value. Represents the position of the legend on the chart. The possible values are: Top, Bottom, Left, Right, Corner, Custom.
func (m *WorkbookChartLegend) SetPosition(value *string)() {
    err := m.GetBackingStore().Set("position", value)
    if err != nil {
        panic(err)
    }
}
// SetVisible sets the visible property value. Indicates whether the chart legend is visible.
func (m *WorkbookChartLegend) SetVisible(value *bool)() {
    err := m.GetBackingStore().Set("visible", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookChartLegendable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFormat()(WorkbookChartLegendFormatable)
    GetOverlay()(*bool)
    GetPosition()(*string)
    GetVisible()(*bool)
    SetFormat(value WorkbookChartLegendFormatable)()
    SetOverlay(value *bool)()
    SetPosition(value *string)()
    SetVisible(value *bool)()
}

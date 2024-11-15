package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookChartSeriesFormat struct {
    Entity
}
// NewWorkbookChartSeriesFormat instantiates a new WorkbookChartSeriesFormat and sets the default values.
func NewWorkbookChartSeriesFormat()(*WorkbookChartSeriesFormat) {
    m := &WorkbookChartSeriesFormat{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartSeriesFormatFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookChartSeriesFormatFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartSeriesFormat(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookChartSeriesFormat) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["fill"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartFillFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFill(val.(WorkbookChartFillable))
        }
        return nil
    }
    res["line"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartLineFormatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLine(val.(WorkbookChartLineFormatable))
        }
        return nil
    }
    return res
}
// GetFill gets the fill property value. Represents the fill format of a chart series, which includes background formatting information. Read-only.
// returns a WorkbookChartFillable when successful
func (m *WorkbookChartSeriesFormat) GetFill()(WorkbookChartFillable) {
    val, err := m.GetBackingStore().Get("fill")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartFillable)
    }
    return nil
}
// GetLine gets the line property value. Represents line formatting. Read-only.
// returns a WorkbookChartLineFormatable when successful
func (m *WorkbookChartSeriesFormat) GetLine()(WorkbookChartLineFormatable) {
    val, err := m.GetBackingStore().Get("line")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartLineFormatable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookChartSeriesFormat) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("fill", m.GetFill())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("line", m.GetLine())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFill sets the fill property value. Represents the fill format of a chart series, which includes background formatting information. Read-only.
func (m *WorkbookChartSeriesFormat) SetFill(value WorkbookChartFillable)() {
    err := m.GetBackingStore().Set("fill", value)
    if err != nil {
        panic(err)
    }
}
// SetLine sets the line property value. Represents line formatting. Read-only.
func (m *WorkbookChartSeriesFormat) SetLine(value WorkbookChartLineFormatable)() {
    err := m.GetBackingStore().Set("line", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookChartSeriesFormatable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFill()(WorkbookChartFillable)
    GetLine()(WorkbookChartLineFormatable)
    SetFill(value WorkbookChartFillable)()
    SetLine(value WorkbookChartLineFormatable)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookChartSeries struct {
    Entity
}
// NewWorkbookChartSeries instantiates a new WorkbookChartSeries and sets the default values.
func NewWorkbookChartSeries()(*WorkbookChartSeries) {
    m := &WorkbookChartSeries{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartSeriesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookChartSeriesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartSeries(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookChartSeries) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["format"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartSeriesFormatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFormat(val.(WorkbookChartSeriesFormatable))
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
    res["points"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookChartPointFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookChartPointable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookChartPointable)
                }
            }
            m.SetPoints(res)
        }
        return nil
    }
    return res
}
// GetFormat gets the format property value. The formatting of a chart series, which includes fill and line formatting. Read-only.
// returns a WorkbookChartSeriesFormatable when successful
func (m *WorkbookChartSeries) GetFormat()(WorkbookChartSeriesFormatable) {
    val, err := m.GetBackingStore().Get("format")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartSeriesFormatable)
    }
    return nil
}
// GetName gets the name property value. The name of a series in a chart.
// returns a *string when successful
func (m *WorkbookChartSeries) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPoints gets the points property value. A collection of all points in the series. Read-only.
// returns a []WorkbookChartPointable when successful
func (m *WorkbookChartSeries) GetPoints()([]WorkbookChartPointable) {
    val, err := m.GetBackingStore().Get("points")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookChartPointable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookChartSeries) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteStringValue("name", m.GetName())
        if err != nil {
            return err
        }
    }
    if m.GetPoints() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPoints()))
        for i, v := range m.GetPoints() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("points", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFormat sets the format property value. The formatting of a chart series, which includes fill and line formatting. Read-only.
func (m *WorkbookChartSeries) SetFormat(value WorkbookChartSeriesFormatable)() {
    err := m.GetBackingStore().Set("format", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The name of a series in a chart.
func (m *WorkbookChartSeries) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetPoints sets the points property value. A collection of all points in the series. Read-only.
func (m *WorkbookChartSeries) SetPoints(value []WorkbookChartPointable)() {
    err := m.GetBackingStore().Set("points", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookChartSeriesable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFormat()(WorkbookChartSeriesFormatable)
    GetName()(*string)
    GetPoints()([]WorkbookChartPointable)
    SetFormat(value WorkbookChartSeriesFormatable)()
    SetName(value *string)()
    SetPoints(value []WorkbookChartPointable)()
}

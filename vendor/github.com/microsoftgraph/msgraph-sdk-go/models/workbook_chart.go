package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookChart struct {
    Entity
}
// NewWorkbookChart instantiates a new WorkbookChart and sets the default values.
func NewWorkbookChart()(*WorkbookChart) {
    m := &WorkbookChart{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookChartFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChart(), nil
}
// GetAxes gets the axes property value. Represents chart axes. Read-only.
// returns a WorkbookChartAxesable when successful
func (m *WorkbookChart) GetAxes()(WorkbookChartAxesable) {
    val, err := m.GetBackingStore().Get("axes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartAxesable)
    }
    return nil
}
// GetDataLabels gets the dataLabels property value. Represents the data labels on the chart. Read-only.
// returns a WorkbookChartDataLabelsable when successful
func (m *WorkbookChart) GetDataLabels()(WorkbookChartDataLabelsable) {
    val, err := m.GetBackingStore().Get("dataLabels")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartDataLabelsable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookChart) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["axes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartAxesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAxes(val.(WorkbookChartAxesable))
        }
        return nil
    }
    res["dataLabels"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartDataLabelsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDataLabels(val.(WorkbookChartDataLabelsable))
        }
        return nil
    }
    res["format"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartAreaFormatFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFormat(val.(WorkbookChartAreaFormatable))
        }
        return nil
    }
    res["height"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHeight(val)
        }
        return nil
    }
    res["left"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLeft(val)
        }
        return nil
    }
    res["legend"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartLegendFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLegend(val.(WorkbookChartLegendable))
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
    res["series"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookChartSeriesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookChartSeriesable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookChartSeriesable)
                }
            }
            m.SetSeries(res)
        }
        return nil
    }
    res["title"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartTitleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTitle(val.(WorkbookChartTitleable))
        }
        return nil
    }
    res["top"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTop(val)
        }
        return nil
    }
    res["width"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWidth(val)
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
// GetFormat gets the format property value. Encapsulates the format properties for the chart area. Read-only.
// returns a WorkbookChartAreaFormatable when successful
func (m *WorkbookChart) GetFormat()(WorkbookChartAreaFormatable) {
    val, err := m.GetBackingStore().Get("format")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartAreaFormatable)
    }
    return nil
}
// GetHeight gets the height property value. Represents the height, in points, of the chart object.
// returns a *float64 when successful
func (m *WorkbookChart) GetHeight()(*float64) {
    val, err := m.GetBackingStore().Get("height")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetLeft gets the left property value. The distance, in points, from the left side of the chart to the worksheet origin.
// returns a *float64 when successful
func (m *WorkbookChart) GetLeft()(*float64) {
    val, err := m.GetBackingStore().Get("left")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetLegend gets the legend property value. Represents the legend for the chart. Read-only.
// returns a WorkbookChartLegendable when successful
func (m *WorkbookChart) GetLegend()(WorkbookChartLegendable) {
    val, err := m.GetBackingStore().Get("legend")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartLegendable)
    }
    return nil
}
// GetName gets the name property value. Represents the name of a chart object.
// returns a *string when successful
func (m *WorkbookChart) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSeries gets the series property value. Represents either a single series or collection of series in the chart. Read-only.
// returns a []WorkbookChartSeriesable when successful
func (m *WorkbookChart) GetSeries()([]WorkbookChartSeriesable) {
    val, err := m.GetBackingStore().Get("series")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookChartSeriesable)
    }
    return nil
}
// GetTitle gets the title property value. Represents the title of the specified chart, including the text, visibility, position and formatting of the title. Read-only.
// returns a WorkbookChartTitleable when successful
func (m *WorkbookChart) GetTitle()(WorkbookChartTitleable) {
    val, err := m.GetBackingStore().Get("title")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartTitleable)
    }
    return nil
}
// GetTop gets the top property value. Represents the distance, in points, from the top edge of the object to the top of row 1 (on a worksheet) or the top of the chart area (on a chart).
// returns a *float64 when successful
func (m *WorkbookChart) GetTop()(*float64) {
    val, err := m.GetBackingStore().Get("top")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetWidth gets the width property value. Represents the width, in points, of the chart object.
// returns a *float64 when successful
func (m *WorkbookChart) GetWidth()(*float64) {
    val, err := m.GetBackingStore().Get("width")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetWorksheet gets the worksheet property value. The worksheet containing the current chart. Read-only.
// returns a WorkbookWorksheetable when successful
func (m *WorkbookChart) GetWorksheet()(WorkbookWorksheetable) {
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
func (m *WorkbookChart) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("axes", m.GetAxes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("dataLabels", m.GetDataLabels())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("format", m.GetFormat())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("height", m.GetHeight())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("left", m.GetLeft())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("legend", m.GetLegend())
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
    if m.GetSeries() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSeries()))
        for i, v := range m.GetSeries() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("series", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("title", m.GetTitle())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("top", m.GetTop())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteFloat64Value("width", m.GetWidth())
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
// SetAxes sets the axes property value. Represents chart axes. Read-only.
func (m *WorkbookChart) SetAxes(value WorkbookChartAxesable)() {
    err := m.GetBackingStore().Set("axes", value)
    if err != nil {
        panic(err)
    }
}
// SetDataLabels sets the dataLabels property value. Represents the data labels on the chart. Read-only.
func (m *WorkbookChart) SetDataLabels(value WorkbookChartDataLabelsable)() {
    err := m.GetBackingStore().Set("dataLabels", value)
    if err != nil {
        panic(err)
    }
}
// SetFormat sets the format property value. Encapsulates the format properties for the chart area. Read-only.
func (m *WorkbookChart) SetFormat(value WorkbookChartAreaFormatable)() {
    err := m.GetBackingStore().Set("format", value)
    if err != nil {
        panic(err)
    }
}
// SetHeight sets the height property value. Represents the height, in points, of the chart object.
func (m *WorkbookChart) SetHeight(value *float64)() {
    err := m.GetBackingStore().Set("height", value)
    if err != nil {
        panic(err)
    }
}
// SetLeft sets the left property value. The distance, in points, from the left side of the chart to the worksheet origin.
func (m *WorkbookChart) SetLeft(value *float64)() {
    err := m.GetBackingStore().Set("left", value)
    if err != nil {
        panic(err)
    }
}
// SetLegend sets the legend property value. Represents the legend for the chart. Read-only.
func (m *WorkbookChart) SetLegend(value WorkbookChartLegendable)() {
    err := m.GetBackingStore().Set("legend", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. Represents the name of a chart object.
func (m *WorkbookChart) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetSeries sets the series property value. Represents either a single series or collection of series in the chart. Read-only.
func (m *WorkbookChart) SetSeries(value []WorkbookChartSeriesable)() {
    err := m.GetBackingStore().Set("series", value)
    if err != nil {
        panic(err)
    }
}
// SetTitle sets the title property value. Represents the title of the specified chart, including the text, visibility, position and formatting of the title. Read-only.
func (m *WorkbookChart) SetTitle(value WorkbookChartTitleable)() {
    err := m.GetBackingStore().Set("title", value)
    if err != nil {
        panic(err)
    }
}
// SetTop sets the top property value. Represents the distance, in points, from the top edge of the object to the top of row 1 (on a worksheet) or the top of the chart area (on a chart).
func (m *WorkbookChart) SetTop(value *float64)() {
    err := m.GetBackingStore().Set("top", value)
    if err != nil {
        panic(err)
    }
}
// SetWidth sets the width property value. Represents the width, in points, of the chart object.
func (m *WorkbookChart) SetWidth(value *float64)() {
    err := m.GetBackingStore().Set("width", value)
    if err != nil {
        panic(err)
    }
}
// SetWorksheet sets the worksheet property value. The worksheet containing the current chart. Read-only.
func (m *WorkbookChart) SetWorksheet(value WorkbookWorksheetable)() {
    err := m.GetBackingStore().Set("worksheet", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookChartable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAxes()(WorkbookChartAxesable)
    GetDataLabels()(WorkbookChartDataLabelsable)
    GetFormat()(WorkbookChartAreaFormatable)
    GetHeight()(*float64)
    GetLeft()(*float64)
    GetLegend()(WorkbookChartLegendable)
    GetName()(*string)
    GetSeries()([]WorkbookChartSeriesable)
    GetTitle()(WorkbookChartTitleable)
    GetTop()(*float64)
    GetWidth()(*float64)
    GetWorksheet()(WorkbookWorksheetable)
    SetAxes(value WorkbookChartAxesable)()
    SetDataLabels(value WorkbookChartDataLabelsable)()
    SetFormat(value WorkbookChartAreaFormatable)()
    SetHeight(value *float64)()
    SetLeft(value *float64)()
    SetLegend(value WorkbookChartLegendable)()
    SetName(value *string)()
    SetSeries(value []WorkbookChartSeriesable)()
    SetTitle(value WorkbookChartTitleable)()
    SetTop(value *float64)()
    SetWidth(value *float64)()
    SetWorksheet(value WorkbookWorksheetable)()
}

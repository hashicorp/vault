package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookChartDataLabelFormat struct {
    Entity
}
// NewWorkbookChartDataLabelFormat instantiates a new WorkbookChartDataLabelFormat and sets the default values.
func NewWorkbookChartDataLabelFormat()(*WorkbookChartDataLabelFormat) {
    m := &WorkbookChartDataLabelFormat{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookChartDataLabelFormatFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookChartDataLabelFormatFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookChartDataLabelFormat(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookChartDataLabelFormat) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["font"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookChartFontFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFont(val.(WorkbookChartFontable))
        }
        return nil
    }
    return res
}
// GetFill gets the fill property value. Represents the fill format of the current chart data label. Read-only.
// returns a WorkbookChartFillable when successful
func (m *WorkbookChartDataLabelFormat) GetFill()(WorkbookChartFillable) {
    val, err := m.GetBackingStore().Get("fill")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartFillable)
    }
    return nil
}
// GetFont gets the font property value. Represents the font attributes (font name, font size, color, etc.) for a chart data label. Read-only.
// returns a WorkbookChartFontable when successful
func (m *WorkbookChartDataLabelFormat) GetFont()(WorkbookChartFontable) {
    val, err := m.GetBackingStore().Get("font")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookChartFontable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookChartDataLabelFormat) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteObjectValue("font", m.GetFont())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFill sets the fill property value. Represents the fill format of the current chart data label. Read-only.
func (m *WorkbookChartDataLabelFormat) SetFill(value WorkbookChartFillable)() {
    err := m.GetBackingStore().Set("fill", value)
    if err != nil {
        panic(err)
    }
}
// SetFont sets the font property value. Represents the font attributes (font name, font size, color, etc.) for a chart data label. Read-only.
func (m *WorkbookChartDataLabelFormat) SetFont(value WorkbookChartFontable)() {
    err := m.GetBackingStore().Set("font", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookChartDataLabelFormatable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetFill()(WorkbookChartFillable)
    GetFont()(WorkbookChartFontable)
    SetFill(value WorkbookChartFillable)()
    SetFont(value WorkbookChartFontable)()
}

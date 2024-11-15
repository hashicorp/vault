package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type WorkbookRangeFont struct {
    Entity
}
// NewWorkbookRangeFont instantiates a new WorkbookRangeFont and sets the default values.
func NewWorkbookRangeFont()(*WorkbookRangeFont) {
    m := &WorkbookRangeFont{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookRangeFontFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookRangeFontFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbookRangeFont(), nil
}
// GetBold gets the bold property value. Inidicates whether the font is bold.
// returns a *bool when successful
func (m *WorkbookRangeFont) GetBold()(*bool) {
    val, err := m.GetBackingStore().Get("bold")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetColor gets the color property value. The HTML color code representation of the text color. For example, #FF0000 represents the color red.
// returns a *string when successful
func (m *WorkbookRangeFont) GetColor()(*string) {
    val, err := m.GetBackingStore().Get("color")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *WorkbookRangeFont) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["bold"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBold(val)
        }
        return nil
    }
    res["color"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetColor(val)
        }
        return nil
    }
    res["italic"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetItalic(val)
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
    res["size"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetFloat64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSize(val)
        }
        return nil
    }
    res["underline"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUnderline(val)
        }
        return nil
    }
    return res
}
// GetItalic gets the italic property value. Inidicates whether the font is italic.
// returns a *bool when successful
func (m *WorkbookRangeFont) GetItalic()(*bool) {
    val, err := m.GetBackingStore().Get("italic")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetName gets the name property value. The font name. For example, 'Calibri'.
// returns a *string when successful
func (m *WorkbookRangeFont) GetName()(*string) {
    val, err := m.GetBackingStore().Get("name")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSize gets the size property value. The font size.
// returns a *float64 when successful
func (m *WorkbookRangeFont) GetSize()(*float64) {
    val, err := m.GetBackingStore().Get("size")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*float64)
    }
    return nil
}
// GetUnderline gets the underline property value. The type of underlining applied to the font. The possible values are: None, Single, Double, SingleAccountant, DoubleAccountant.
// returns a *string when successful
func (m *WorkbookRangeFont) GetUnderline()(*string) {
    val, err := m.GetBackingStore().Get("underline")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *WorkbookRangeFont) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("bold", m.GetBold())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("color", m.GetColor())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("italic", m.GetItalic())
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
    {
        err = writer.WriteFloat64Value("size", m.GetSize())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("underline", m.GetUnderline())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBold sets the bold property value. Inidicates whether the font is bold.
func (m *WorkbookRangeFont) SetBold(value *bool)() {
    err := m.GetBackingStore().Set("bold", value)
    if err != nil {
        panic(err)
    }
}
// SetColor sets the color property value. The HTML color code representation of the text color. For example, #FF0000 represents the color red.
func (m *WorkbookRangeFont) SetColor(value *string)() {
    err := m.GetBackingStore().Set("color", value)
    if err != nil {
        panic(err)
    }
}
// SetItalic sets the italic property value. Inidicates whether the font is italic.
func (m *WorkbookRangeFont) SetItalic(value *bool)() {
    err := m.GetBackingStore().Set("italic", value)
    if err != nil {
        panic(err)
    }
}
// SetName sets the name property value. The font name. For example, 'Calibri'.
func (m *WorkbookRangeFont) SetName(value *string)() {
    err := m.GetBackingStore().Set("name", value)
    if err != nil {
        panic(err)
    }
}
// SetSize sets the size property value. The font size.
func (m *WorkbookRangeFont) SetSize(value *float64)() {
    err := m.GetBackingStore().Set("size", value)
    if err != nil {
        panic(err)
    }
}
// SetUnderline sets the underline property value. The type of underlining applied to the font. The possible values are: None, Single, Double, SingleAccountant, DoubleAccountant.
func (m *WorkbookRangeFont) SetUnderline(value *string)() {
    err := m.GetBackingStore().Set("underline", value)
    if err != nil {
        panic(err)
    }
}
type WorkbookRangeFontable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBold()(*bool)
    GetColor()(*string)
    GetItalic()(*bool)
    GetName()(*string)
    GetSize()(*float64)
    GetUnderline()(*string)
    SetBold(value *bool)()
    SetColor(value *string)()
    SetItalic(value *bool)()
    SetName(value *string)()
    SetSize(value *float64)()
    SetUnderline(value *string)()
}

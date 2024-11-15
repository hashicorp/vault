package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TextWebPart struct {
    WebPart
}
// NewTextWebPart instantiates a new TextWebPart and sets the default values.
func NewTextWebPart()(*TextWebPart) {
    m := &TextWebPart{
        WebPart: *NewWebPart(),
    }
    odataTypeValue := "#microsoft.graph.textWebPart"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTextWebPartFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTextWebPartFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTextWebPart(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TextWebPart) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WebPart.GetFieldDeserializers()
    res["innerHtml"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInnerHtml(val)
        }
        return nil
    }
    return res
}
// GetInnerHtml gets the innerHtml property value. The HTML string in text web part.
// returns a *string when successful
func (m *TextWebPart) GetInnerHtml()(*string) {
    val, err := m.GetBackingStore().Get("innerHtml")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TextWebPart) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WebPart.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("innerHtml", m.GetInnerHtml())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInnerHtml sets the innerHtml property value. The HTML string in text web part.
func (m *TextWebPart) SetInnerHtml(value *string)() {
    err := m.GetBackingStore().Set("innerHtml", value)
    if err != nil {
        panic(err)
    }
}
type TextWebPartable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    WebPartable
    GetInnerHtml()(*string)
    SetInnerHtml(value *string)()
}

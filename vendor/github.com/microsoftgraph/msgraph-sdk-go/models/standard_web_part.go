package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type StandardWebPart struct {
    WebPart
}
// NewStandardWebPart instantiates a new StandardWebPart and sets the default values.
func NewStandardWebPart()(*StandardWebPart) {
    m := &StandardWebPart{
        WebPart: *NewWebPart(),
    }
    odataTypeValue := "#microsoft.graph.standardWebPart"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateStandardWebPartFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateStandardWebPartFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewStandardWebPart(), nil
}
// GetContainerTextWebPartId gets the containerTextWebPartId property value. The instance identifier of the container text webPart. It only works for inline standard webPart in rich text webParts.
// returns a *string when successful
func (m *StandardWebPart) GetContainerTextWebPartId()(*string) {
    val, err := m.GetBackingStore().Get("containerTextWebPartId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetData gets the data property value. Data of the webPart.
// returns a WebPartDataable when successful
func (m *StandardWebPart) GetData()(WebPartDataable) {
    val, err := m.GetBackingStore().Get("data")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WebPartDataable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *StandardWebPart) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.WebPart.GetFieldDeserializers()
    res["containerTextWebPartId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContainerTextWebPartId(val)
        }
        return nil
    }
    res["data"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWebPartDataFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetData(val.(WebPartDataable))
        }
        return nil
    }
    res["webPartType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebPartType(val)
        }
        return nil
    }
    return res
}
// GetWebPartType gets the webPartType property value. A Guid that indicates the webPart type.
// returns a *string when successful
func (m *StandardWebPart) GetWebPartType()(*string) {
    val, err := m.GetBackingStore().Get("webPartType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *StandardWebPart) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.WebPart.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("containerTextWebPartId", m.GetContainerTextWebPartId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("data", m.GetData())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webPartType", m.GetWebPartType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContainerTextWebPartId sets the containerTextWebPartId property value. The instance identifier of the container text webPart. It only works for inline standard webPart in rich text webParts.
func (m *StandardWebPart) SetContainerTextWebPartId(value *string)() {
    err := m.GetBackingStore().Set("containerTextWebPartId", value)
    if err != nil {
        panic(err)
    }
}
// SetData sets the data property value. Data of the webPart.
func (m *StandardWebPart) SetData(value WebPartDataable)() {
    err := m.GetBackingStore().Set("data", value)
    if err != nil {
        panic(err)
    }
}
// SetWebPartType sets the webPartType property value. A Guid that indicates the webPart type.
func (m *StandardWebPart) SetWebPartType(value *string)() {
    err := m.GetBackingStore().Set("webPartType", value)
    if err != nil {
        panic(err)
    }
}
type StandardWebPartable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    WebPartable
    GetContainerTextWebPartId()(*string)
    GetData()(WebPartDataable)
    GetWebPartType()(*string)
    SetContainerTextWebPartId(value *string)()
    SetData(value WebPartDataable)()
    SetWebPartType(value *string)()
}

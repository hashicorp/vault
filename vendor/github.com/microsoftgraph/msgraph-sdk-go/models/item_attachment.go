package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ItemAttachment struct {
    Attachment
}
// NewItemAttachment instantiates a new ItemAttachment and sets the default values.
func NewItemAttachment()(*ItemAttachment) {
    m := &ItemAttachment{
        Attachment: *NewAttachment(),
    }
    odataTypeValue := "#microsoft.graph.itemAttachment"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateItemAttachmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateItemAttachmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemAttachment(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ItemAttachment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Attachment.GetFieldDeserializers()
    res["item"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateOutlookItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetItem(val.(OutlookItemable))
        }
        return nil
    }
    return res
}
// GetItem gets the item property value. The attached message or event. Navigation property.
// returns a OutlookItemable when successful
func (m *ItemAttachment) GetItem()(OutlookItemable) {
    val, err := m.GetBackingStore().Get("item")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(OutlookItemable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ItemAttachment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Attachment.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("item", m.GetItem())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetItem sets the item property value. The attached message or event. Navigation property.
func (m *ItemAttachment) SetItem(value OutlookItemable)() {
    err := m.GetBackingStore().Set("item", value)
    if err != nil {
        panic(err)
    }
}
type ItemAttachmentable interface {
    Attachmentable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetItem()(OutlookItemable)
    SetItem(value OutlookItemable)()
}

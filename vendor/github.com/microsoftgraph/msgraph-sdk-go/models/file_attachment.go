package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type FileAttachment struct {
    Attachment
}
// NewFileAttachment instantiates a new FileAttachment and sets the default values.
func NewFileAttachment()(*FileAttachment) {
    m := &FileAttachment{
        Attachment: *NewAttachment(),
    }
    odataTypeValue := "#microsoft.graph.fileAttachment"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateFileAttachmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateFileAttachmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewFileAttachment(), nil
}
// GetContentBytes gets the contentBytes property value. The base64-encoded contents of the file.
// returns a []byte when successful
func (m *FileAttachment) GetContentBytes()([]byte) {
    val, err := m.GetBackingStore().Get("contentBytes")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetContentId gets the contentId property value. The ID of the attachment in the Exchange store.
// returns a *string when successful
func (m *FileAttachment) GetContentId()(*string) {
    val, err := m.GetBackingStore().Get("contentId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetContentLocation gets the contentLocation property value. Don't use this property as it isn't supported.
// returns a *string when successful
func (m *FileAttachment) GetContentLocation()(*string) {
    val, err := m.GetBackingStore().Get("contentLocation")
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
func (m *FileAttachment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Attachment.GetFieldDeserializers()
    res["contentBytes"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentBytes(val)
        }
        return nil
    }
    res["contentId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentId(val)
        }
        return nil
    }
    res["contentLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetContentLocation(val)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *FileAttachment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Attachment.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteByteArrayValue("contentBytes", m.GetContentBytes())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("contentId", m.GetContentId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("contentLocation", m.GetContentLocation())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetContentBytes sets the contentBytes property value. The base64-encoded contents of the file.
func (m *FileAttachment) SetContentBytes(value []byte)() {
    err := m.GetBackingStore().Set("contentBytes", value)
    if err != nil {
        panic(err)
    }
}
// SetContentId sets the contentId property value. The ID of the attachment in the Exchange store.
func (m *FileAttachment) SetContentId(value *string)() {
    err := m.GetBackingStore().Set("contentId", value)
    if err != nil {
        panic(err)
    }
}
// SetContentLocation sets the contentLocation property value. Don't use this property as it isn't supported.
func (m *FileAttachment) SetContentLocation(value *string)() {
    err := m.GetBackingStore().Set("contentLocation", value)
    if err != nil {
        panic(err)
    }
}
type FileAttachmentable interface {
    Attachmentable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetContentBytes()([]byte)
    GetContentId()(*string)
    GetContentLocation()(*string)
    SetContentBytes(value []byte)()
    SetContentId(value *string)()
    SetContentLocation(value *string)()
}

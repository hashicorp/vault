package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EmailPayloadDetail struct {
    PayloadDetail
}
// NewEmailPayloadDetail instantiates a new EmailPayloadDetail and sets the default values.
func NewEmailPayloadDetail()(*EmailPayloadDetail) {
    m := &EmailPayloadDetail{
        PayloadDetail: *NewPayloadDetail(),
    }
    odataTypeValue := "#microsoft.graph.emailPayloadDetail"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEmailPayloadDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEmailPayloadDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEmailPayloadDetail(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EmailPayloadDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PayloadDetail.GetFieldDeserializers()
    res["fromEmail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFromEmail(val)
        }
        return nil
    }
    res["fromName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFromName(val)
        }
        return nil
    }
    res["isExternalSender"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsExternalSender(val)
        }
        return nil
    }
    res["subject"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSubject(val)
        }
        return nil
    }
    return res
}
// GetFromEmail gets the fromEmail property value. Email address of the user.
// returns a *string when successful
func (m *EmailPayloadDetail) GetFromEmail()(*string) {
    val, err := m.GetBackingStore().Get("fromEmail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFromName gets the fromName property value. Display name of the user.
// returns a *string when successful
func (m *EmailPayloadDetail) GetFromName()(*string) {
    val, err := m.GetBackingStore().Get("fromName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsExternalSender gets the isExternalSender property value. Indicates whether the sender isn't from the user's organization.
// returns a *bool when successful
func (m *EmailPayloadDetail) GetIsExternalSender()(*bool) {
    val, err := m.GetBackingStore().Get("isExternalSender")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSubject gets the subject property value. The subject of the email address sent to the user.
// returns a *string when successful
func (m *EmailPayloadDetail) GetSubject()(*string) {
    val, err := m.GetBackingStore().Get("subject")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EmailPayloadDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PayloadDetail.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("fromEmail", m.GetFromEmail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("fromName", m.GetFromName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isExternalSender", m.GetIsExternalSender())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("subject", m.GetSubject())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetFromEmail sets the fromEmail property value. Email address of the user.
func (m *EmailPayloadDetail) SetFromEmail(value *string)() {
    err := m.GetBackingStore().Set("fromEmail", value)
    if err != nil {
        panic(err)
    }
}
// SetFromName sets the fromName property value. Display name of the user.
func (m *EmailPayloadDetail) SetFromName(value *string)() {
    err := m.GetBackingStore().Set("fromName", value)
    if err != nil {
        panic(err)
    }
}
// SetIsExternalSender sets the isExternalSender property value. Indicates whether the sender isn't from the user's organization.
func (m *EmailPayloadDetail) SetIsExternalSender(value *bool)() {
    err := m.GetBackingStore().Set("isExternalSender", value)
    if err != nil {
        panic(err)
    }
}
// SetSubject sets the subject property value. The subject of the email address sent to the user.
func (m *EmailPayloadDetail) SetSubject(value *string)() {
    err := m.GetBackingStore().Set("subject", value)
    if err != nil {
        panic(err)
    }
}
type EmailPayloadDetailable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    PayloadDetailable
    GetFromEmail()(*string)
    GetFromName()(*string)
    GetIsExternalSender()(*bool)
    GetSubject()(*string)
    SetFromEmail(value *string)()
    SetFromName(value *string)()
    SetIsExternalSender(value *bool)()
    SetSubject(value *string)()
}

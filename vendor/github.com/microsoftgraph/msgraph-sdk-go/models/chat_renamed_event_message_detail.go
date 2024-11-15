package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ChatRenamedEventMessageDetail struct {
    EventMessageDetail
}
// NewChatRenamedEventMessageDetail instantiates a new ChatRenamedEventMessageDetail and sets the default values.
func NewChatRenamedEventMessageDetail()(*ChatRenamedEventMessageDetail) {
    m := &ChatRenamedEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.chatRenamedEventMessageDetail"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateChatRenamedEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateChatRenamedEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChatRenamedEventMessageDetail(), nil
}
// GetChatDisplayName gets the chatDisplayName property value. The updated name of the chat.
// returns a *string when successful
func (m *ChatRenamedEventMessageDetail) GetChatDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("chatDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetChatId gets the chatId property value. Unique identifier of the chat.
// returns a *string when successful
func (m *ChatRenamedEventMessageDetail) GetChatId()(*string) {
    val, err := m.GetBackingStore().Get("chatId")
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
func (m *ChatRenamedEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessageDetail.GetFieldDeserializers()
    res["chatDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChatDisplayName(val)
        }
        return nil
    }
    res["chatId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChatId(val)
        }
        return nil
    }
    res["initiator"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInitiator(val.(IdentitySetable))
        }
        return nil
    }
    return res
}
// GetInitiator gets the initiator property value. Initiator of the event.
// returns a IdentitySetable when successful
func (m *ChatRenamedEventMessageDetail) GetInitiator()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("initiator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ChatRenamedEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessageDetail.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("chatDisplayName", m.GetChatDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("chatId", m.GetChatId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("initiator", m.GetInitiator())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetChatDisplayName sets the chatDisplayName property value. The updated name of the chat.
func (m *ChatRenamedEventMessageDetail) SetChatDisplayName(value *string)() {
    err := m.GetBackingStore().Set("chatDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetChatId sets the chatId property value. Unique identifier of the chat.
func (m *ChatRenamedEventMessageDetail) SetChatId(value *string)() {
    err := m.GetBackingStore().Set("chatId", value)
    if err != nil {
        panic(err)
    }
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *ChatRenamedEventMessageDetail) SetInitiator(value IdentitySetable)() {
    err := m.GetBackingStore().Set("initiator", value)
    if err != nil {
        panic(err)
    }
}
type ChatRenamedEventMessageDetailable interface {
    EventMessageDetailable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetChatDisplayName()(*string)
    GetChatId()(*string)
    GetInitiator()(IdentitySetable)
    SetChatDisplayName(value *string)()
    SetChatId(value *string)()
    SetInitiator(value IdentitySetable)()
}

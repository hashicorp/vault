package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ChatMessageInfo struct {
    Entity
}
// NewChatMessageInfo instantiates a new ChatMessageInfo and sets the default values.
func NewChatMessageInfo()(*ChatMessageInfo) {
    m := &ChatMessageInfo{
        Entity: *NewEntity(),
    }
    return m
}
// CreateChatMessageInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateChatMessageInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChatMessageInfo(), nil
}
// GetBody gets the body property value. Body of the chatMessage. This will still contain markers for @mentions and attachments even though the object doesn't return @mentions and attachments.
// returns a ItemBodyable when successful
func (m *ChatMessageInfo) GetBody()(ItemBodyable) {
    val, err := m.GetBackingStore().Get("body")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemBodyable)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Date time object representing the time at which message was created.
// returns a *Time when successful
func (m *ChatMessageInfo) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetEventDetail gets the eventDetail property value. Read-only.  If present, represents details of an event that happened in a chat, a channel, or a team, for example, members were added, and so on. For event messages, the messageType property is set to systemEventMessage.
// returns a EventMessageDetailable when successful
func (m *ChatMessageInfo) GetEventDetail()(EventMessageDetailable) {
    val, err := m.GetBackingStore().Get("eventDetail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EventMessageDetailable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ChatMessageInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["body"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemBodyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBody(val.(ItemBodyable))
        }
        return nil
    }
    res["createdDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCreatedDateTime(val)
        }
        return nil
    }
    res["eventDetail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEventMessageDetailFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEventDetail(val.(EventMessageDetailable))
        }
        return nil
    }
    res["from"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChatMessageFromIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFrom(val.(ChatMessageFromIdentitySetable))
        }
        return nil
    }
    res["isDeleted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDeleted(val)
        }
        return nil
    }
    res["messageType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseChatMessageType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMessageType(val.(*ChatMessageType))
        }
        return nil
    }
    return res
}
// GetFrom gets the from property value. Information about the sender of the message.
// returns a ChatMessageFromIdentitySetable when successful
func (m *ChatMessageInfo) GetFrom()(ChatMessageFromIdentitySetable) {
    val, err := m.GetBackingStore().Get("from")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChatMessageFromIdentitySetable)
    }
    return nil
}
// GetIsDeleted gets the isDeleted property value. If set to true, the original message has been deleted.
// returns a *bool when successful
func (m *ChatMessageInfo) GetIsDeleted()(*bool) {
    val, err := m.GetBackingStore().Get("isDeleted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMessageType gets the messageType property value. The messageType property
// returns a *ChatMessageType when successful
func (m *ChatMessageInfo) GetMessageType()(*ChatMessageType) {
    val, err := m.GetBackingStore().Get("messageType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ChatMessageType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ChatMessageInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("body", m.GetBody())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("eventDetail", m.GetEventDetail())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("from", m.GetFrom())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isDeleted", m.GetIsDeleted())
        if err != nil {
            return err
        }
    }
    if m.GetMessageType() != nil {
        cast := (*m.GetMessageType()).String()
        err = writer.WriteStringValue("messageType", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetBody sets the body property value. Body of the chatMessage. This will still contain markers for @mentions and attachments even though the object doesn't return @mentions and attachments.
func (m *ChatMessageInfo) SetBody(value ItemBodyable)() {
    err := m.GetBackingStore().Set("body", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Date time object representing the time at which message was created.
func (m *ChatMessageInfo) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetEventDetail sets the eventDetail property value. Read-only.  If present, represents details of an event that happened in a chat, a channel, or a team, for example, members were added, and so on. For event messages, the messageType property is set to systemEventMessage.
func (m *ChatMessageInfo) SetEventDetail(value EventMessageDetailable)() {
    err := m.GetBackingStore().Set("eventDetail", value)
    if err != nil {
        panic(err)
    }
}
// SetFrom sets the from property value. Information about the sender of the message.
func (m *ChatMessageInfo) SetFrom(value ChatMessageFromIdentitySetable)() {
    err := m.GetBackingStore().Set("from", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDeleted sets the isDeleted property value. If set to true, the original message has been deleted.
func (m *ChatMessageInfo) SetIsDeleted(value *bool)() {
    err := m.GetBackingStore().Set("isDeleted", value)
    if err != nil {
        panic(err)
    }
}
// SetMessageType sets the messageType property value. The messageType property
func (m *ChatMessageInfo) SetMessageType(value *ChatMessageType)() {
    err := m.GetBackingStore().Set("messageType", value)
    if err != nil {
        panic(err)
    }
}
type ChatMessageInfoable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBody()(ItemBodyable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetEventDetail()(EventMessageDetailable)
    GetFrom()(ChatMessageFromIdentitySetable)
    GetIsDeleted()(*bool)
    GetMessageType()(*ChatMessageType)
    SetBody(value ItemBodyable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetEventDetail(value EventMessageDetailable)()
    SetFrom(value ChatMessageFromIdentitySetable)()
    SetIsDeleted(value *bool)()
    SetMessageType(value *ChatMessageType)()
}

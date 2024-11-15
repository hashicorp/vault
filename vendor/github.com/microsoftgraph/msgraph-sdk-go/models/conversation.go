package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Conversation struct {
    Entity
}
// NewConversation instantiates a new Conversation and sets the default values.
func NewConversation()(*Conversation) {
    m := &Conversation{
        Entity: *NewEntity(),
    }
    return m
}
// CreateConversationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConversationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConversation(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Conversation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["hasAttachments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHasAttachments(val)
        }
        return nil
    }
    res["lastDeliveredDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastDeliveredDateTime(val)
        }
        return nil
    }
    res["preview"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreview(val)
        }
        return nil
    }
    res["threads"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateConversationThreadFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ConversationThreadable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ConversationThreadable)
                }
            }
            m.SetThreads(res)
        }
        return nil
    }
    res["topic"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTopic(val)
        }
        return nil
    }
    res["uniqueSenders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfPrimitiveValues("string")
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]string, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = *(v.(*string))
                }
            }
            m.SetUniqueSenders(res)
        }
        return nil
    }
    return res
}
// GetHasAttachments gets the hasAttachments property value. Indicates whether any of the posts within this Conversation has at least one attachment. Supports $filter (eq, ne) and $search.
// returns a *bool when successful
func (m *Conversation) GetHasAttachments()(*bool) {
    val, err := m.GetBackingStore().Get("hasAttachments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLastDeliveredDateTime gets the lastDeliveredDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
// returns a *Time when successful
func (m *Conversation) GetLastDeliveredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastDeliveredDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetPreview gets the preview property value. A short summary from the body of the latest post in this conversation. Supports $filter (eq, ne, le, ge).
// returns a *string when successful
func (m *Conversation) GetPreview()(*string) {
    val, err := m.GetBackingStore().Get("preview")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetThreads gets the threads property value. A collection of all the conversation threads in the conversation. A navigation property. Read-only. Nullable.
// returns a []ConversationThreadable when successful
func (m *Conversation) GetThreads()([]ConversationThreadable) {
    val, err := m.GetBackingStore().Get("threads")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ConversationThreadable)
    }
    return nil
}
// GetTopic gets the topic property value. The topic of the conversation. This property can be set when the conversation is created, but it cannot be updated.
// returns a *string when successful
func (m *Conversation) GetTopic()(*string) {
    val, err := m.GetBackingStore().Get("topic")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetUniqueSenders gets the uniqueSenders property value. All the users that sent a message to this Conversation.
// returns a []string when successful
func (m *Conversation) GetUniqueSenders()([]string) {
    val, err := m.GetBackingStore().Get("uniqueSenders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Conversation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("hasAttachments", m.GetHasAttachments())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastDeliveredDateTime", m.GetLastDeliveredDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("preview", m.GetPreview())
        if err != nil {
            return err
        }
    }
    if m.GetThreads() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetThreads()))
        for i, v := range m.GetThreads() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("threads", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("topic", m.GetTopic())
        if err != nil {
            return err
        }
    }
    if m.GetUniqueSenders() != nil {
        err = writer.WriteCollectionOfStringValues("uniqueSenders", m.GetUniqueSenders())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetHasAttachments sets the hasAttachments property value. Indicates whether any of the posts within this Conversation has at least one attachment. Supports $filter (eq, ne) and $search.
func (m *Conversation) SetHasAttachments(value *bool)() {
    err := m.GetBackingStore().Set("hasAttachments", value)
    if err != nil {
        panic(err)
    }
}
// SetLastDeliveredDateTime sets the lastDeliveredDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z
func (m *Conversation) SetLastDeliveredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastDeliveredDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPreview sets the preview property value. A short summary from the body of the latest post in this conversation. Supports $filter (eq, ne, le, ge).
func (m *Conversation) SetPreview(value *string)() {
    err := m.GetBackingStore().Set("preview", value)
    if err != nil {
        panic(err)
    }
}
// SetThreads sets the threads property value. A collection of all the conversation threads in the conversation. A navigation property. Read-only. Nullable.
func (m *Conversation) SetThreads(value []ConversationThreadable)() {
    err := m.GetBackingStore().Set("threads", value)
    if err != nil {
        panic(err)
    }
}
// SetTopic sets the topic property value. The topic of the conversation. This property can be set when the conversation is created, but it cannot be updated.
func (m *Conversation) SetTopic(value *string)() {
    err := m.GetBackingStore().Set("topic", value)
    if err != nil {
        panic(err)
    }
}
// SetUniqueSenders sets the uniqueSenders property value. All the users that sent a message to this Conversation.
func (m *Conversation) SetUniqueSenders(value []string)() {
    err := m.GetBackingStore().Set("uniqueSenders", value)
    if err != nil {
        panic(err)
    }
}
type Conversationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetHasAttachments()(*bool)
    GetLastDeliveredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPreview()(*string)
    GetThreads()([]ConversationThreadable)
    GetTopic()(*string)
    GetUniqueSenders()([]string)
    SetHasAttachments(value *bool)()
    SetLastDeliveredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPreview(value *string)()
    SetThreads(value []ConversationThreadable)()
    SetTopic(value *string)()
    SetUniqueSenders(value []string)()
}

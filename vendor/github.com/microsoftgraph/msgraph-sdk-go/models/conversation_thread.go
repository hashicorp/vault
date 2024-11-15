package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ConversationThread struct {
    Entity
}
// NewConversationThread instantiates a new ConversationThread and sets the default values.
func NewConversationThread()(*ConversationThread) {
    m := &ConversationThread{
        Entity: *NewEntity(),
    }
    return m
}
// CreateConversationThreadFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateConversationThreadFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewConversationThread(), nil
}
// GetCcRecipients gets the ccRecipients property value. The Cc: recipients for the thread. Returned only on $select.
// returns a []Recipientable when successful
func (m *ConversationThread) GetCcRecipients()([]Recipientable) {
    val, err := m.GetBackingStore().Get("ccRecipients")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Recipientable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ConversationThread) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["ccRecipients"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRecipientFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Recipientable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Recipientable)
                }
            }
            m.SetCcRecipients(res)
        }
        return nil
    }
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
    res["isLocked"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsLocked(val)
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
    res["posts"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreatePostFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Postable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Postable)
                }
            }
            m.SetPosts(res)
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
    res["toRecipients"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateRecipientFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Recipientable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Recipientable)
                }
            }
            m.SetToRecipients(res)
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
// GetHasAttachments gets the hasAttachments property value. Indicates whether any of the posts within this thread has at least one attachment. Returned by default.
// returns a *bool when successful
func (m *ConversationThread) GetHasAttachments()(*bool) {
    val, err := m.GetBackingStore().Get("hasAttachments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsLocked gets the isLocked property value. Indicates if the thread is locked. Returned by default.
// returns a *bool when successful
func (m *ConversationThread) GetIsLocked()(*bool) {
    val, err := m.GetBackingStore().Get("isLocked")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLastDeliveredDateTime gets the lastDeliveredDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.Returned by default.
// returns a *Time when successful
func (m *ConversationThread) GetLastDeliveredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastDeliveredDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetPosts gets the posts property value. The posts property
// returns a []Postable when successful
func (m *ConversationThread) GetPosts()([]Postable) {
    val, err := m.GetBackingStore().Get("posts")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Postable)
    }
    return nil
}
// GetPreview gets the preview property value. A short summary from the body of the latest post in this conversation. Returned by default.
// returns a *string when successful
func (m *ConversationThread) GetPreview()(*string) {
    val, err := m.GetBackingStore().Get("preview")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTopic gets the topic property value. The topic of the conversation. This property can be set when the conversation is created, but it cannot be updated. Returned by default.
// returns a *string when successful
func (m *ConversationThread) GetTopic()(*string) {
    val, err := m.GetBackingStore().Get("topic")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetToRecipients gets the toRecipients property value. The To: recipients for the thread. Returned only on $select.
// returns a []Recipientable when successful
func (m *ConversationThread) GetToRecipients()([]Recipientable) {
    val, err := m.GetBackingStore().Get("toRecipients")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Recipientable)
    }
    return nil
}
// GetUniqueSenders gets the uniqueSenders property value. All the users that sent a message to this thread. Returned by default.
// returns a []string when successful
func (m *ConversationThread) GetUniqueSenders()([]string) {
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
func (m *ConversationThread) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetCcRecipients() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCcRecipients()))
        for i, v := range m.GetCcRecipients() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("ccRecipients", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hasAttachments", m.GetHasAttachments())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isLocked", m.GetIsLocked())
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
    if m.GetPosts() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetPosts()))
        for i, v := range m.GetPosts() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("posts", cast)
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
    {
        err = writer.WriteStringValue("topic", m.GetTopic())
        if err != nil {
            return err
        }
    }
    if m.GetToRecipients() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetToRecipients()))
        for i, v := range m.GetToRecipients() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("toRecipients", cast)
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
// SetCcRecipients sets the ccRecipients property value. The Cc: recipients for the thread. Returned only on $select.
func (m *ConversationThread) SetCcRecipients(value []Recipientable)() {
    err := m.GetBackingStore().Set("ccRecipients", value)
    if err != nil {
        panic(err)
    }
}
// SetHasAttachments sets the hasAttachments property value. Indicates whether any of the posts within this thread has at least one attachment. Returned by default.
func (m *ConversationThread) SetHasAttachments(value *bool)() {
    err := m.GetBackingStore().Set("hasAttachments", value)
    if err != nil {
        panic(err)
    }
}
// SetIsLocked sets the isLocked property value. Indicates if the thread is locked. Returned by default.
func (m *ConversationThread) SetIsLocked(value *bool)() {
    err := m.GetBackingStore().Set("isLocked", value)
    if err != nil {
        panic(err)
    }
}
// SetLastDeliveredDateTime sets the lastDeliveredDateTime property value. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.Returned by default.
func (m *ConversationThread) SetLastDeliveredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastDeliveredDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPosts sets the posts property value. The posts property
func (m *ConversationThread) SetPosts(value []Postable)() {
    err := m.GetBackingStore().Set("posts", value)
    if err != nil {
        panic(err)
    }
}
// SetPreview sets the preview property value. A short summary from the body of the latest post in this conversation. Returned by default.
func (m *ConversationThread) SetPreview(value *string)() {
    err := m.GetBackingStore().Set("preview", value)
    if err != nil {
        panic(err)
    }
}
// SetTopic sets the topic property value. The topic of the conversation. This property can be set when the conversation is created, but it cannot be updated. Returned by default.
func (m *ConversationThread) SetTopic(value *string)() {
    err := m.GetBackingStore().Set("topic", value)
    if err != nil {
        panic(err)
    }
}
// SetToRecipients sets the toRecipients property value. The To: recipients for the thread. Returned only on $select.
func (m *ConversationThread) SetToRecipients(value []Recipientable)() {
    err := m.GetBackingStore().Set("toRecipients", value)
    if err != nil {
        panic(err)
    }
}
// SetUniqueSenders sets the uniqueSenders property value. All the users that sent a message to this thread. Returned by default.
func (m *ConversationThread) SetUniqueSenders(value []string)() {
    err := m.GetBackingStore().Set("uniqueSenders", value)
    if err != nil {
        panic(err)
    }
}
type ConversationThreadable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCcRecipients()([]Recipientable)
    GetHasAttachments()(*bool)
    GetIsLocked()(*bool)
    GetLastDeliveredDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPosts()([]Postable)
    GetPreview()(*string)
    GetTopic()(*string)
    GetToRecipients()([]Recipientable)
    GetUniqueSenders()([]string)
    SetCcRecipients(value []Recipientable)()
    SetHasAttachments(value *bool)()
    SetIsLocked(value *bool)()
    SetLastDeliveredDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPosts(value []Postable)()
    SetPreview(value *string)()
    SetTopic(value *string)()
    SetToRecipients(value []Recipientable)()
    SetUniqueSenders(value []string)()
}

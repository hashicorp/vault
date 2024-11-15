package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ChatMessage struct {
    Entity
}
// NewChatMessage instantiates a new ChatMessage and sets the default values.
func NewChatMessage()(*ChatMessage) {
    m := &ChatMessage{
        Entity: *NewEntity(),
    }
    return m
}
// CreateChatMessageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateChatMessageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChatMessage(), nil
}
// GetAttachments gets the attachments property value. References to attached objects like files, tabs, meetings etc.
// returns a []ChatMessageAttachmentable when successful
func (m *ChatMessage) GetAttachments()([]ChatMessageAttachmentable) {
    val, err := m.GetBackingStore().Get("attachments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ChatMessageAttachmentable)
    }
    return nil
}
// GetBody gets the body property value. The body property
// returns a ItemBodyable when successful
func (m *ChatMessage) GetBody()(ItemBodyable) {
    val, err := m.GetBackingStore().Get("body")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemBodyable)
    }
    return nil
}
// GetChannelIdentity gets the channelIdentity property value. If the message was sent in a channel, represents identity of the channel.
// returns a ChannelIdentityable when successful
func (m *ChatMessage) GetChannelIdentity()(ChannelIdentityable) {
    val, err := m.GetBackingStore().Get("channelIdentity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChannelIdentityable)
    }
    return nil
}
// GetChatId gets the chatId property value. If the message was sent in a chat, represents the identity of the chat.
// returns a *string when successful
func (m *ChatMessage) GetChatId()(*string) {
    val, err := m.GetBackingStore().Get("chatId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCreatedDateTime gets the createdDateTime property value. Timestamp of when the chat message was created.
// returns a *Time when successful
func (m *ChatMessage) GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("createdDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetDeletedDateTime gets the deletedDateTime property value. Read only. Timestamp at which the chat message was deleted, or null if not deleted.
// returns a *Time when successful
func (m *ChatMessage) GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("deletedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetEtag gets the etag property value. Read-only. Version number of the chat message.
// returns a *string when successful
func (m *ChatMessage) GetEtag()(*string) {
    val, err := m.GetBackingStore().Get("etag")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEventDetail gets the eventDetail property value. Read-only. If present, represents details of an event that happened in a chat, a channel, or a team, for example, adding new members. For event messages, the messageType property will be set to systemEventMessage.
// returns a EventMessageDetailable when successful
func (m *ChatMessage) GetEventDetail()(EventMessageDetailable) {
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
func (m *ChatMessage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["attachments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChatMessageAttachmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ChatMessageAttachmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ChatMessageAttachmentable)
                }
            }
            m.SetAttachments(res)
        }
        return nil
    }
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
    res["channelIdentity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChannelIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetChannelIdentity(val.(ChannelIdentityable))
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
    res["deletedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeletedDateTime(val)
        }
        return nil
    }
    res["etag"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEtag(val)
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
    res["hostedContents"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChatMessageHostedContentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ChatMessageHostedContentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ChatMessageHostedContentable)
                }
            }
            m.SetHostedContents(res)
        }
        return nil
    }
    res["importance"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseChatMessageImportance)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImportance(val.(*ChatMessageImportance))
        }
        return nil
    }
    res["lastEditedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastEditedDateTime(val)
        }
        return nil
    }
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["locale"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocale(val)
        }
        return nil
    }
    res["mentions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChatMessageMentionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ChatMessageMentionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ChatMessageMentionable)
                }
            }
            m.SetMentions(res)
        }
        return nil
    }
    res["messageHistory"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChatMessageHistoryItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ChatMessageHistoryItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ChatMessageHistoryItemable)
                }
            }
            m.SetMessageHistory(res)
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
    res["policyViolation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChatMessagePolicyViolationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyViolation(val.(ChatMessagePolicyViolationable))
        }
        return nil
    }
    res["reactions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChatMessageReactionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ChatMessageReactionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ChatMessageReactionable)
                }
            }
            m.SetReactions(res)
        }
        return nil
    }
    res["replies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateChatMessageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ChatMessageable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ChatMessageable)
                }
            }
            m.SetReplies(res)
        }
        return nil
    }
    res["replyToId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReplyToId(val)
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
    res["summary"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSummary(val)
        }
        return nil
    }
    res["webUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebUrl(val)
        }
        return nil
    }
    return res
}
// GetFrom gets the from property value. Details of the sender of the chat message. Can only be set during migration.
// returns a ChatMessageFromIdentitySetable when successful
func (m *ChatMessage) GetFrom()(ChatMessageFromIdentitySetable) {
    val, err := m.GetBackingStore().Get("from")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChatMessageFromIdentitySetable)
    }
    return nil
}
// GetHostedContents gets the hostedContents property value. Content in a message hosted by Microsoft Teams - for example, images or code snippets.
// returns a []ChatMessageHostedContentable when successful
func (m *ChatMessage) GetHostedContents()([]ChatMessageHostedContentable) {
    val, err := m.GetBackingStore().Get("hostedContents")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ChatMessageHostedContentable)
    }
    return nil
}
// GetImportance gets the importance property value. The importance property
// returns a *ChatMessageImportance when successful
func (m *ChatMessage) GetImportance()(*ChatMessageImportance) {
    val, err := m.GetBackingStore().Get("importance")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ChatMessageImportance)
    }
    return nil
}
// GetLastEditedDateTime gets the lastEditedDateTime property value. Read only. Timestamp when edits to the chat message were made. Triggers an 'Edited' flag in the Teams UI. If no edits are made the value is null.
// returns a *Time when successful
func (m *ChatMessage) GetLastEditedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastEditedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. Read only. Timestamp when the chat message is created (initial setting) or modified, including when a reaction is added or removed.
// returns a *Time when successful
func (m *ChatMessage) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetLocale gets the locale property value. Locale of the chat message set by the client. Always set to en-us.
// returns a *string when successful
func (m *ChatMessage) GetLocale()(*string) {
    val, err := m.GetBackingStore().Get("locale")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMentions gets the mentions property value. List of entities mentioned in the chat message. Supported entities are: user, bot, team, and channel.
// returns a []ChatMessageMentionable when successful
func (m *ChatMessage) GetMentions()([]ChatMessageMentionable) {
    val, err := m.GetBackingStore().Get("mentions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ChatMessageMentionable)
    }
    return nil
}
// GetMessageHistory gets the messageHistory property value. List of activity history of a message item, including modification time and actions, such as reactionAdded, reactionRemoved, or reaction changes, on the message.
// returns a []ChatMessageHistoryItemable when successful
func (m *ChatMessage) GetMessageHistory()([]ChatMessageHistoryItemable) {
    val, err := m.GetBackingStore().Get("messageHistory")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ChatMessageHistoryItemable)
    }
    return nil
}
// GetMessageType gets the messageType property value. The messageType property
// returns a *ChatMessageType when successful
func (m *ChatMessage) GetMessageType()(*ChatMessageType) {
    val, err := m.GetBackingStore().Get("messageType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ChatMessageType)
    }
    return nil
}
// GetPolicyViolation gets the policyViolation property value. Defines the properties of a policy violation set by a data loss prevention (DLP) application.
// returns a ChatMessagePolicyViolationable when successful
func (m *ChatMessage) GetPolicyViolation()(ChatMessagePolicyViolationable) {
    val, err := m.GetBackingStore().Get("policyViolation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChatMessagePolicyViolationable)
    }
    return nil
}
// GetReactions gets the reactions property value. Reactions for this chat message (for example, Like).
// returns a []ChatMessageReactionable when successful
func (m *ChatMessage) GetReactions()([]ChatMessageReactionable) {
    val, err := m.GetBackingStore().Get("reactions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ChatMessageReactionable)
    }
    return nil
}
// GetReplies gets the replies property value. Replies for a specified message. Supports $expand for channel messages.
// returns a []ChatMessageable when successful
func (m *ChatMessage) GetReplies()([]ChatMessageable) {
    val, err := m.GetBackingStore().Get("replies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ChatMessageable)
    }
    return nil
}
// GetReplyToId gets the replyToId property value. Read-only. ID of the parent chat message or root chat message of the thread. (Only applies to chat messages in channels, not chats.)
// returns a *string when successful
func (m *ChatMessage) GetReplyToId()(*string) {
    val, err := m.GetBackingStore().Get("replyToId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSubject gets the subject property value. The subject of the chat message, in plaintext.
// returns a *string when successful
func (m *ChatMessage) GetSubject()(*string) {
    val, err := m.GetBackingStore().Get("subject")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSummary gets the summary property value. Summary text of the chat message that could be used for push notifications and summary views or fall back views. Only applies to channel chat messages, not chat messages in a chat.
// returns a *string when successful
func (m *ChatMessage) GetSummary()(*string) {
    val, err := m.GetBackingStore().Get("summary")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetWebUrl gets the webUrl property value. Read-only. Link to the message in Microsoft Teams.
// returns a *string when successful
func (m *ChatMessage) GetWebUrl()(*string) {
    val, err := m.GetBackingStore().Get("webUrl")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ChatMessage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAttachments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetAttachments()))
        for i, v := range m.GetAttachments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("attachments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("body", m.GetBody())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("channelIdentity", m.GetChannelIdentity())
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
        err = writer.WriteTimeValue("createdDateTime", m.GetCreatedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("deletedDateTime", m.GetDeletedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("etag", m.GetEtag())
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
    if m.GetHostedContents() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetHostedContents()))
        for i, v := range m.GetHostedContents() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("hostedContents", cast)
        if err != nil {
            return err
        }
    }
    if m.GetImportance() != nil {
        cast := (*m.GetImportance()).String()
        err = writer.WriteStringValue("importance", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastEditedDateTime", m.GetLastEditedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("locale", m.GetLocale())
        if err != nil {
            return err
        }
    }
    if m.GetMentions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMentions()))
        for i, v := range m.GetMentions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("mentions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMessageHistory() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMessageHistory()))
        for i, v := range m.GetMessageHistory() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("messageHistory", cast)
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
    {
        err = writer.WriteObjectValue("policyViolation", m.GetPolicyViolation())
        if err != nil {
            return err
        }
    }
    if m.GetReactions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetReactions()))
        for i, v := range m.GetReactions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("reactions", cast)
        if err != nil {
            return err
        }
    }
    if m.GetReplies() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetReplies()))
        for i, v := range m.GetReplies() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("replies", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("replyToId", m.GetReplyToId())
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
    {
        err = writer.WriteStringValue("summary", m.GetSummary())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webUrl", m.GetWebUrl())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAttachments sets the attachments property value. References to attached objects like files, tabs, meetings etc.
func (m *ChatMessage) SetAttachments(value []ChatMessageAttachmentable)() {
    err := m.GetBackingStore().Set("attachments", value)
    if err != nil {
        panic(err)
    }
}
// SetBody sets the body property value. The body property
func (m *ChatMessage) SetBody(value ItemBodyable)() {
    err := m.GetBackingStore().Set("body", value)
    if err != nil {
        panic(err)
    }
}
// SetChannelIdentity sets the channelIdentity property value. If the message was sent in a channel, represents identity of the channel.
func (m *ChatMessage) SetChannelIdentity(value ChannelIdentityable)() {
    err := m.GetBackingStore().Set("channelIdentity", value)
    if err != nil {
        panic(err)
    }
}
// SetChatId sets the chatId property value. If the message was sent in a chat, represents the identity of the chat.
func (m *ChatMessage) SetChatId(value *string)() {
    err := m.GetBackingStore().Set("chatId", value)
    if err != nil {
        panic(err)
    }
}
// SetCreatedDateTime sets the createdDateTime property value. Timestamp of when the chat message was created.
func (m *ChatMessage) SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("createdDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetDeletedDateTime sets the deletedDateTime property value. Read only. Timestamp at which the chat message was deleted, or null if not deleted.
func (m *ChatMessage) SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("deletedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetEtag sets the etag property value. Read-only. Version number of the chat message.
func (m *ChatMessage) SetEtag(value *string)() {
    err := m.GetBackingStore().Set("etag", value)
    if err != nil {
        panic(err)
    }
}
// SetEventDetail sets the eventDetail property value. Read-only. If present, represents details of an event that happened in a chat, a channel, or a team, for example, adding new members. For event messages, the messageType property will be set to systemEventMessage.
func (m *ChatMessage) SetEventDetail(value EventMessageDetailable)() {
    err := m.GetBackingStore().Set("eventDetail", value)
    if err != nil {
        panic(err)
    }
}
// SetFrom sets the from property value. Details of the sender of the chat message. Can only be set during migration.
func (m *ChatMessage) SetFrom(value ChatMessageFromIdentitySetable)() {
    err := m.GetBackingStore().Set("from", value)
    if err != nil {
        panic(err)
    }
}
// SetHostedContents sets the hostedContents property value. Content in a message hosted by Microsoft Teams - for example, images or code snippets.
func (m *ChatMessage) SetHostedContents(value []ChatMessageHostedContentable)() {
    err := m.GetBackingStore().Set("hostedContents", value)
    if err != nil {
        panic(err)
    }
}
// SetImportance sets the importance property value. The importance property
func (m *ChatMessage) SetImportance(value *ChatMessageImportance)() {
    err := m.GetBackingStore().Set("importance", value)
    if err != nil {
        panic(err)
    }
}
// SetLastEditedDateTime sets the lastEditedDateTime property value. Read only. Timestamp when edits to the chat message were made. Triggers an 'Edited' flag in the Teams UI. If no edits are made the value is null.
func (m *ChatMessage) SetLastEditedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastEditedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. Read only. Timestamp when the chat message is created (initial setting) or modified, including when a reaction is added or removed.
func (m *ChatMessage) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetLocale sets the locale property value. Locale of the chat message set by the client. Always set to en-us.
func (m *ChatMessage) SetLocale(value *string)() {
    err := m.GetBackingStore().Set("locale", value)
    if err != nil {
        panic(err)
    }
}
// SetMentions sets the mentions property value. List of entities mentioned in the chat message. Supported entities are: user, bot, team, and channel.
func (m *ChatMessage) SetMentions(value []ChatMessageMentionable)() {
    err := m.GetBackingStore().Set("mentions", value)
    if err != nil {
        panic(err)
    }
}
// SetMessageHistory sets the messageHistory property value. List of activity history of a message item, including modification time and actions, such as reactionAdded, reactionRemoved, or reaction changes, on the message.
func (m *ChatMessage) SetMessageHistory(value []ChatMessageHistoryItemable)() {
    err := m.GetBackingStore().Set("messageHistory", value)
    if err != nil {
        panic(err)
    }
}
// SetMessageType sets the messageType property value. The messageType property
func (m *ChatMessage) SetMessageType(value *ChatMessageType)() {
    err := m.GetBackingStore().Set("messageType", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyViolation sets the policyViolation property value. Defines the properties of a policy violation set by a data loss prevention (DLP) application.
func (m *ChatMessage) SetPolicyViolation(value ChatMessagePolicyViolationable)() {
    err := m.GetBackingStore().Set("policyViolation", value)
    if err != nil {
        panic(err)
    }
}
// SetReactions sets the reactions property value. Reactions for this chat message (for example, Like).
func (m *ChatMessage) SetReactions(value []ChatMessageReactionable)() {
    err := m.GetBackingStore().Set("reactions", value)
    if err != nil {
        panic(err)
    }
}
// SetReplies sets the replies property value. Replies for a specified message. Supports $expand for channel messages.
func (m *ChatMessage) SetReplies(value []ChatMessageable)() {
    err := m.GetBackingStore().Set("replies", value)
    if err != nil {
        panic(err)
    }
}
// SetReplyToId sets the replyToId property value. Read-only. ID of the parent chat message or root chat message of the thread. (Only applies to chat messages in channels, not chats.)
func (m *ChatMessage) SetReplyToId(value *string)() {
    err := m.GetBackingStore().Set("replyToId", value)
    if err != nil {
        panic(err)
    }
}
// SetSubject sets the subject property value. The subject of the chat message, in plaintext.
func (m *ChatMessage) SetSubject(value *string)() {
    err := m.GetBackingStore().Set("subject", value)
    if err != nil {
        panic(err)
    }
}
// SetSummary sets the summary property value. Summary text of the chat message that could be used for push notifications and summary views or fall back views. Only applies to channel chat messages, not chat messages in a chat.
func (m *ChatMessage) SetSummary(value *string)() {
    err := m.GetBackingStore().Set("summary", value)
    if err != nil {
        panic(err)
    }
}
// SetWebUrl sets the webUrl property value. Read-only. Link to the message in Microsoft Teams.
func (m *ChatMessage) SetWebUrl(value *string)() {
    err := m.GetBackingStore().Set("webUrl", value)
    if err != nil {
        panic(err)
    }
}
type ChatMessageable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttachments()([]ChatMessageAttachmentable)
    GetBody()(ItemBodyable)
    GetChannelIdentity()(ChannelIdentityable)
    GetChatId()(*string)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetDeletedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetEtag()(*string)
    GetEventDetail()(EventMessageDetailable)
    GetFrom()(ChatMessageFromIdentitySetable)
    GetHostedContents()([]ChatMessageHostedContentable)
    GetImportance()(*ChatMessageImportance)
    GetLastEditedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetLocale()(*string)
    GetMentions()([]ChatMessageMentionable)
    GetMessageHistory()([]ChatMessageHistoryItemable)
    GetMessageType()(*ChatMessageType)
    GetPolicyViolation()(ChatMessagePolicyViolationable)
    GetReactions()([]ChatMessageReactionable)
    GetReplies()([]ChatMessageable)
    GetReplyToId()(*string)
    GetSubject()(*string)
    GetSummary()(*string)
    GetWebUrl()(*string)
    SetAttachments(value []ChatMessageAttachmentable)()
    SetBody(value ItemBodyable)()
    SetChannelIdentity(value ChannelIdentityable)()
    SetChatId(value *string)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetDeletedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetEtag(value *string)()
    SetEventDetail(value EventMessageDetailable)()
    SetFrom(value ChatMessageFromIdentitySetable)()
    SetHostedContents(value []ChatMessageHostedContentable)()
    SetImportance(value *ChatMessageImportance)()
    SetLastEditedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetLocale(value *string)()
    SetMentions(value []ChatMessageMentionable)()
    SetMessageHistory(value []ChatMessageHistoryItemable)()
    SetMessageType(value *ChatMessageType)()
    SetPolicyViolation(value ChatMessagePolicyViolationable)()
    SetReactions(value []ChatMessageReactionable)()
    SetReplies(value []ChatMessageable)()
    SetReplyToId(value *string)()
    SetSubject(value *string)()
    SetSummary(value *string)()
    SetWebUrl(value *string)()
}

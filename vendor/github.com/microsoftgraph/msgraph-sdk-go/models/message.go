package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Message struct {
    OutlookItem
}
// NewMessage instantiates a new Message and sets the default values.
func NewMessage()(*Message) {
    m := &Message{
        OutlookItem: *NewOutlookItem(),
    }
    odataTypeValue := "#microsoft.graph.message"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMessageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMessageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.calendarSharingMessage":
                        return NewCalendarSharingMessage(), nil
                    case "#microsoft.graph.eventMessage":
                        return NewEventMessage(), nil
                    case "#microsoft.graph.eventMessageRequest":
                        return NewEventMessageRequest(), nil
                    case "#microsoft.graph.eventMessageResponse":
                        return NewEventMessageResponse(), nil
                }
            }
        }
    }
    return NewMessage(), nil
}
// GetAttachments gets the attachments property value. The fileAttachment and itemAttachment attachments for the message.
// returns a []Attachmentable when successful
func (m *Message) GetAttachments()([]Attachmentable) {
    val, err := m.GetBackingStore().Get("attachments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Attachmentable)
    }
    return nil
}
// GetBccRecipients gets the bccRecipients property value. The Bcc: recipients for the message.
// returns a []Recipientable when successful
func (m *Message) GetBccRecipients()([]Recipientable) {
    val, err := m.GetBackingStore().Get("bccRecipients")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Recipientable)
    }
    return nil
}
// GetBody gets the body property value. The body of the message. It can be in HTML or text format. Find out about safe HTML in a message body.
// returns a ItemBodyable when successful
func (m *Message) GetBody()(ItemBodyable) {
    val, err := m.GetBackingStore().Get("body")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemBodyable)
    }
    return nil
}
// GetBodyPreview gets the bodyPreview property value. The first 255 characters of the message body. It is in text format.
// returns a *string when successful
func (m *Message) GetBodyPreview()(*string) {
    val, err := m.GetBackingStore().Get("bodyPreview")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCcRecipients gets the ccRecipients property value. The Cc: recipients for the message.
// returns a []Recipientable when successful
func (m *Message) GetCcRecipients()([]Recipientable) {
    val, err := m.GetBackingStore().Get("ccRecipients")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Recipientable)
    }
    return nil
}
// GetConversationId gets the conversationId property value. The ID of the conversation the email belongs to.
// returns a *string when successful
func (m *Message) GetConversationId()(*string) {
    val, err := m.GetBackingStore().Get("conversationId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetConversationIndex gets the conversationIndex property value. Indicates the position of the message within the conversation.
// returns a []byte when successful
func (m *Message) GetConversationIndex()([]byte) {
    val, err := m.GetBackingStore().Get("conversationIndex")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]byte)
    }
    return nil
}
// GetExtensions gets the extensions property value. The collection of open extensions defined for the message. Nullable.
// returns a []Extensionable when successful
func (m *Message) GetExtensions()([]Extensionable) {
    val, err := m.GetBackingStore().Get("extensions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Extensionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Message) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.OutlookItem.GetFieldDeserializers()
    res["attachments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateAttachmentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Attachmentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Attachmentable)
                }
            }
            m.SetAttachments(res)
        }
        return nil
    }
    res["bccRecipients"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetBccRecipients(res)
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
    res["bodyPreview"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBodyPreview(val)
        }
        return nil
    }
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
    res["conversationId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConversationId(val)
        }
        return nil
    }
    res["conversationIndex"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetByteArrayValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConversationIndex(val)
        }
        return nil
    }
    res["extensions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateExtensionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Extensionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Extensionable)
                }
            }
            m.SetExtensions(res)
        }
        return nil
    }
    res["flag"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateFollowupFlagFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFlag(val.(FollowupFlagable))
        }
        return nil
    }
    res["from"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRecipientFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFrom(val.(Recipientable))
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
    res["importance"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseImportance)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetImportance(val.(*Importance))
        }
        return nil
    }
    res["inferenceClassification"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseInferenceClassificationType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInferenceClassification(val.(*InferenceClassificationType))
        }
        return nil
    }
    res["internetMessageHeaders"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateInternetMessageHeaderFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]InternetMessageHeaderable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(InternetMessageHeaderable)
                }
            }
            m.SetInternetMessageHeaders(res)
        }
        return nil
    }
    res["internetMessageId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInternetMessageId(val)
        }
        return nil
    }
    res["isDeliveryReceiptRequested"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDeliveryReceiptRequested(val)
        }
        return nil
    }
    res["isDraft"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDraft(val)
        }
        return nil
    }
    res["isRead"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsRead(val)
        }
        return nil
    }
    res["isReadReceiptRequested"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsReadReceiptRequested(val)
        }
        return nil
    }
    res["multiValueExtendedProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMultiValueLegacyExtendedPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MultiValueLegacyExtendedPropertyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MultiValueLegacyExtendedPropertyable)
                }
            }
            m.SetMultiValueExtendedProperties(res)
        }
        return nil
    }
    res["parentFolderId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParentFolderId(val)
        }
        return nil
    }
    res["receivedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetReceivedDateTime(val)
        }
        return nil
    }
    res["replyTo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetReplyTo(res)
        }
        return nil
    }
    res["sender"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRecipientFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSender(val.(Recipientable))
        }
        return nil
    }
    res["sentDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSentDateTime(val)
        }
        return nil
    }
    res["singleValueExtendedProperties"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateSingleValueLegacyExtendedPropertyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]SingleValueLegacyExtendedPropertyable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(SingleValueLegacyExtendedPropertyable)
                }
            }
            m.SetSingleValueExtendedProperties(res)
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
    res["uniqueBody"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateItemBodyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUniqueBody(val.(ItemBodyable))
        }
        return nil
    }
    res["webLink"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWebLink(val)
        }
        return nil
    }
    return res
}
// GetFlag gets the flag property value. The flag value that indicates the status, start date, due date, or completion date for the message.
// returns a FollowupFlagable when successful
func (m *Message) GetFlag()(FollowupFlagable) {
    val, err := m.GetBackingStore().Get("flag")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(FollowupFlagable)
    }
    return nil
}
// GetFrom gets the from property value. The owner of the mailbox from which the message is sent. In most cases, this value is the same as the sender property, except for sharing or delegation scenarios. The value must correspond to the actual mailbox used. Find out more about setting the from and sender properties of a message.
// returns a Recipientable when successful
func (m *Message) GetFrom()(Recipientable) {
    val, err := m.GetBackingStore().Get("from")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Recipientable)
    }
    return nil
}
// GetHasAttachments gets the hasAttachments property value. Indicates whether the message has attachments. This property doesn't include inline attachments, so if a message contains only inline attachments, this property is false. To verify the existence of inline attachments, parse the body property to look for a src attribute, such as <IMG src='cid:image001.jpg@01D26CD8.6C05F070'>.
// returns a *bool when successful
func (m *Message) GetHasAttachments()(*bool) {
    val, err := m.GetBackingStore().Get("hasAttachments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetImportance gets the importance property value. The importance of the message. The possible values are: low, normal, and high.
// returns a *Importance when successful
func (m *Message) GetImportance()(*Importance) {
    val, err := m.GetBackingStore().Get("importance")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Importance)
    }
    return nil
}
// GetInferenceClassification gets the inferenceClassification property value. The classification of the message for the user, based on inferred relevance or importance, or on an explicit override. The possible values are: focused or other.
// returns a *InferenceClassificationType when successful
func (m *Message) GetInferenceClassification()(*InferenceClassificationType) {
    val, err := m.GetBackingStore().Get("inferenceClassification")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*InferenceClassificationType)
    }
    return nil
}
// GetInternetMessageHeaders gets the internetMessageHeaders property value. A collection of message headers defined by RFC5322. The set includes message headers indicating the network path taken by a message from the sender to the recipient. It can also contain custom message headers that hold app data for the message.  Returned only on applying a $select query option. Read-only.
// returns a []InternetMessageHeaderable when successful
func (m *Message) GetInternetMessageHeaders()([]InternetMessageHeaderable) {
    val, err := m.GetBackingStore().Get("internetMessageHeaders")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]InternetMessageHeaderable)
    }
    return nil
}
// GetInternetMessageId gets the internetMessageId property value. The message ID in the format specified by RFC2822.
// returns a *string when successful
func (m *Message) GetInternetMessageId()(*string) {
    val, err := m.GetBackingStore().Get("internetMessageId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetIsDeliveryReceiptRequested gets the isDeliveryReceiptRequested property value. Indicates whether a read receipt is requested for the message.
// returns a *bool when successful
func (m *Message) GetIsDeliveryReceiptRequested()(*bool) {
    val, err := m.GetBackingStore().Get("isDeliveryReceiptRequested")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsDraft gets the isDraft property value. Indicates whether the message is a draft. A message is a draft if it hasn't been sent yet.
// returns a *bool when successful
func (m *Message) GetIsDraft()(*bool) {
    val, err := m.GetBackingStore().Get("isDraft")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsRead gets the isRead property value. Indicates whether the message has been read.
// returns a *bool when successful
func (m *Message) GetIsRead()(*bool) {
    val, err := m.GetBackingStore().Get("isRead")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsReadReceiptRequested gets the isReadReceiptRequested property value. Indicates whether a read receipt is requested for the message.
// returns a *bool when successful
func (m *Message) GetIsReadReceiptRequested()(*bool) {
    val, err := m.GetBackingStore().Get("isReadReceiptRequested")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMultiValueExtendedProperties gets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the message. Nullable.
// returns a []MultiValueLegacyExtendedPropertyable when successful
func (m *Message) GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable) {
    val, err := m.GetBackingStore().Get("multiValueExtendedProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MultiValueLegacyExtendedPropertyable)
    }
    return nil
}
// GetParentFolderId gets the parentFolderId property value. The unique identifier for the message's parent mailFolder.
// returns a *string when successful
func (m *Message) GetParentFolderId()(*string) {
    val, err := m.GetBackingStore().Get("parentFolderId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetReceivedDateTime gets the receivedDateTime property value. The date and time the message was received.  The date and time information uses ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Message) GetReceivedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("receivedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetReplyTo gets the replyTo property value. The email addresses to use when replying.
// returns a []Recipientable when successful
func (m *Message) GetReplyTo()([]Recipientable) {
    val, err := m.GetBackingStore().Get("replyTo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Recipientable)
    }
    return nil
}
// GetSender gets the sender property value. The account that is actually used to generate the message. In most cases, this value is the same as the from property. You can set this property to a different value when sending a message from a shared mailbox, for a shared calendar, or as a delegate. In any case, the value must correspond to the actual mailbox used. Find out more about setting the from and sender properties of a message.
// returns a Recipientable when successful
func (m *Message) GetSender()(Recipientable) {
    val, err := m.GetBackingStore().Get("sender")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Recipientable)
    }
    return nil
}
// GetSentDateTime gets the sentDateTime property value. The date and time the message was sent.  The date and time information uses ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
// returns a *Time when successful
func (m *Message) GetSentDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("sentDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetSingleValueExtendedProperties gets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the message. Nullable.
// returns a []SingleValueLegacyExtendedPropertyable when successful
func (m *Message) GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable) {
    val, err := m.GetBackingStore().Get("singleValueExtendedProperties")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]SingleValueLegacyExtendedPropertyable)
    }
    return nil
}
// GetSubject gets the subject property value. The subject of the message.
// returns a *string when successful
func (m *Message) GetSubject()(*string) {
    val, err := m.GetBackingStore().Get("subject")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetToRecipients gets the toRecipients property value. The To: recipients for the message.
// returns a []Recipientable when successful
func (m *Message) GetToRecipients()([]Recipientable) {
    val, err := m.GetBackingStore().Get("toRecipients")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Recipientable)
    }
    return nil
}
// GetUniqueBody gets the uniqueBody property value. The part of the body of the message that is unique to the current message. uniqueBody is not returned by default but can be retrieved for a given message by use of the ?$select=uniqueBody query. It can be in HTML or text format.
// returns a ItemBodyable when successful
func (m *Message) GetUniqueBody()(ItemBodyable) {
    val, err := m.GetBackingStore().Get("uniqueBody")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ItemBodyable)
    }
    return nil
}
// GetWebLink gets the webLink property value. The URL to open the message in Outlook on the web.You can append an ispopout argument to the end of the URL to change how the message is displayed. If ispopout is not present or if it is set to 1, then the message is shown in a popout window. If ispopout is set to 0, the browser shows the message in the Outlook on the web review pane.The message opens in the browser if you are signed in to your mailbox via Outlook on the web. You are prompted to sign in if you are not already signed in with the browser.This URL cannot be accessed from within an iFrame.
// returns a *string when successful
func (m *Message) GetWebLink()(*string) {
    val, err := m.GetBackingStore().Get("webLink")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Message) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.OutlookItem.Serialize(writer)
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
    if m.GetBccRecipients() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetBccRecipients()))
        for i, v := range m.GetBccRecipients() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("bccRecipients", cast)
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
        err = writer.WriteStringValue("bodyPreview", m.GetBodyPreview())
        if err != nil {
            return err
        }
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
        err = writer.WriteStringValue("conversationId", m.GetConversationId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteByteArrayValue("conversationIndex", m.GetConversationIndex())
        if err != nil {
            return err
        }
    }
    if m.GetExtensions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExtensions()))
        for i, v := range m.GetExtensions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("extensions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("flag", m.GetFlag())
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
        err = writer.WriteBoolValue("hasAttachments", m.GetHasAttachments())
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
    if m.GetInferenceClassification() != nil {
        cast := (*m.GetInferenceClassification()).String()
        err = writer.WriteStringValue("inferenceClassification", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetInternetMessageHeaders() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetInternetMessageHeaders()))
        for i, v := range m.GetInternetMessageHeaders() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("internetMessageHeaders", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("internetMessageId", m.GetInternetMessageId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isDeliveryReceiptRequested", m.GetIsDeliveryReceiptRequested())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isDraft", m.GetIsDraft())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isRead", m.GetIsRead())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isReadReceiptRequested", m.GetIsReadReceiptRequested())
        if err != nil {
            return err
        }
    }
    if m.GetMultiValueExtendedProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMultiValueExtendedProperties()))
        for i, v := range m.GetMultiValueExtendedProperties() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("multiValueExtendedProperties", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("parentFolderId", m.GetParentFolderId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("receivedDateTime", m.GetReceivedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetReplyTo() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetReplyTo()))
        for i, v := range m.GetReplyTo() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("replyTo", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sender", m.GetSender())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("sentDateTime", m.GetSentDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetSingleValueExtendedProperties() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSingleValueExtendedProperties()))
        for i, v := range m.GetSingleValueExtendedProperties() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("singleValueExtendedProperties", cast)
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
    {
        err = writer.WriteObjectValue("uniqueBody", m.GetUniqueBody())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("webLink", m.GetWebLink())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAttachments sets the attachments property value. The fileAttachment and itemAttachment attachments for the message.
func (m *Message) SetAttachments(value []Attachmentable)() {
    err := m.GetBackingStore().Set("attachments", value)
    if err != nil {
        panic(err)
    }
}
// SetBccRecipients sets the bccRecipients property value. The Bcc: recipients for the message.
func (m *Message) SetBccRecipients(value []Recipientable)() {
    err := m.GetBackingStore().Set("bccRecipients", value)
    if err != nil {
        panic(err)
    }
}
// SetBody sets the body property value. The body of the message. It can be in HTML or text format. Find out about safe HTML in a message body.
func (m *Message) SetBody(value ItemBodyable)() {
    err := m.GetBackingStore().Set("body", value)
    if err != nil {
        panic(err)
    }
}
// SetBodyPreview sets the bodyPreview property value. The first 255 characters of the message body. It is in text format.
func (m *Message) SetBodyPreview(value *string)() {
    err := m.GetBackingStore().Set("bodyPreview", value)
    if err != nil {
        panic(err)
    }
}
// SetCcRecipients sets the ccRecipients property value. The Cc: recipients for the message.
func (m *Message) SetCcRecipients(value []Recipientable)() {
    err := m.GetBackingStore().Set("ccRecipients", value)
    if err != nil {
        panic(err)
    }
}
// SetConversationId sets the conversationId property value. The ID of the conversation the email belongs to.
func (m *Message) SetConversationId(value *string)() {
    err := m.GetBackingStore().Set("conversationId", value)
    if err != nil {
        panic(err)
    }
}
// SetConversationIndex sets the conversationIndex property value. Indicates the position of the message within the conversation.
func (m *Message) SetConversationIndex(value []byte)() {
    err := m.GetBackingStore().Set("conversationIndex", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensions sets the extensions property value. The collection of open extensions defined for the message. Nullable.
func (m *Message) SetExtensions(value []Extensionable)() {
    err := m.GetBackingStore().Set("extensions", value)
    if err != nil {
        panic(err)
    }
}
// SetFlag sets the flag property value. The flag value that indicates the status, start date, due date, or completion date for the message.
func (m *Message) SetFlag(value FollowupFlagable)() {
    err := m.GetBackingStore().Set("flag", value)
    if err != nil {
        panic(err)
    }
}
// SetFrom sets the from property value. The owner of the mailbox from which the message is sent. In most cases, this value is the same as the sender property, except for sharing or delegation scenarios. The value must correspond to the actual mailbox used. Find out more about setting the from and sender properties of a message.
func (m *Message) SetFrom(value Recipientable)() {
    err := m.GetBackingStore().Set("from", value)
    if err != nil {
        panic(err)
    }
}
// SetHasAttachments sets the hasAttachments property value. Indicates whether the message has attachments. This property doesn't include inline attachments, so if a message contains only inline attachments, this property is false. To verify the existence of inline attachments, parse the body property to look for a src attribute, such as <IMG src='cid:image001.jpg@01D26CD8.6C05F070'>.
func (m *Message) SetHasAttachments(value *bool)() {
    err := m.GetBackingStore().Set("hasAttachments", value)
    if err != nil {
        panic(err)
    }
}
// SetImportance sets the importance property value. The importance of the message. The possible values are: low, normal, and high.
func (m *Message) SetImportance(value *Importance)() {
    err := m.GetBackingStore().Set("importance", value)
    if err != nil {
        panic(err)
    }
}
// SetInferenceClassification sets the inferenceClassification property value. The classification of the message for the user, based on inferred relevance or importance, or on an explicit override. The possible values are: focused or other.
func (m *Message) SetInferenceClassification(value *InferenceClassificationType)() {
    err := m.GetBackingStore().Set("inferenceClassification", value)
    if err != nil {
        panic(err)
    }
}
// SetInternetMessageHeaders sets the internetMessageHeaders property value. A collection of message headers defined by RFC5322. The set includes message headers indicating the network path taken by a message from the sender to the recipient. It can also contain custom message headers that hold app data for the message.  Returned only on applying a $select query option. Read-only.
func (m *Message) SetInternetMessageHeaders(value []InternetMessageHeaderable)() {
    err := m.GetBackingStore().Set("internetMessageHeaders", value)
    if err != nil {
        panic(err)
    }
}
// SetInternetMessageId sets the internetMessageId property value. The message ID in the format specified by RFC2822.
func (m *Message) SetInternetMessageId(value *string)() {
    err := m.GetBackingStore().Set("internetMessageId", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDeliveryReceiptRequested sets the isDeliveryReceiptRequested property value. Indicates whether a read receipt is requested for the message.
func (m *Message) SetIsDeliveryReceiptRequested(value *bool)() {
    err := m.GetBackingStore().Set("isDeliveryReceiptRequested", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDraft sets the isDraft property value. Indicates whether the message is a draft. A message is a draft if it hasn't been sent yet.
func (m *Message) SetIsDraft(value *bool)() {
    err := m.GetBackingStore().Set("isDraft", value)
    if err != nil {
        panic(err)
    }
}
// SetIsRead sets the isRead property value. Indicates whether the message has been read.
func (m *Message) SetIsRead(value *bool)() {
    err := m.GetBackingStore().Set("isRead", value)
    if err != nil {
        panic(err)
    }
}
// SetIsReadReceiptRequested sets the isReadReceiptRequested property value. Indicates whether a read receipt is requested for the message.
func (m *Message) SetIsReadReceiptRequested(value *bool)() {
    err := m.GetBackingStore().Set("isReadReceiptRequested", value)
    if err != nil {
        panic(err)
    }
}
// SetMultiValueExtendedProperties sets the multiValueExtendedProperties property value. The collection of multi-value extended properties defined for the message. Nullable.
func (m *Message) SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)() {
    err := m.GetBackingStore().Set("multiValueExtendedProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetParentFolderId sets the parentFolderId property value. The unique identifier for the message's parent mailFolder.
func (m *Message) SetParentFolderId(value *string)() {
    err := m.GetBackingStore().Set("parentFolderId", value)
    if err != nil {
        panic(err)
    }
}
// SetReceivedDateTime sets the receivedDateTime property value. The date and time the message was received.  The date and time information uses ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Message) SetReceivedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("receivedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetReplyTo sets the replyTo property value. The email addresses to use when replying.
func (m *Message) SetReplyTo(value []Recipientable)() {
    err := m.GetBackingStore().Set("replyTo", value)
    if err != nil {
        panic(err)
    }
}
// SetSender sets the sender property value. The account that is actually used to generate the message. In most cases, this value is the same as the from property. You can set this property to a different value when sending a message from a shared mailbox, for a shared calendar, or as a delegate. In any case, the value must correspond to the actual mailbox used. Find out more about setting the from and sender properties of a message.
func (m *Message) SetSender(value Recipientable)() {
    err := m.GetBackingStore().Set("sender", value)
    if err != nil {
        panic(err)
    }
}
// SetSentDateTime sets the sentDateTime property value. The date and time the message was sent.  The date and time information uses ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
func (m *Message) SetSentDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("sentDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetSingleValueExtendedProperties sets the singleValueExtendedProperties property value. The collection of single-value extended properties defined for the message. Nullable.
func (m *Message) SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)() {
    err := m.GetBackingStore().Set("singleValueExtendedProperties", value)
    if err != nil {
        panic(err)
    }
}
// SetSubject sets the subject property value. The subject of the message.
func (m *Message) SetSubject(value *string)() {
    err := m.GetBackingStore().Set("subject", value)
    if err != nil {
        panic(err)
    }
}
// SetToRecipients sets the toRecipients property value. The To: recipients for the message.
func (m *Message) SetToRecipients(value []Recipientable)() {
    err := m.GetBackingStore().Set("toRecipients", value)
    if err != nil {
        panic(err)
    }
}
// SetUniqueBody sets the uniqueBody property value. The part of the body of the message that is unique to the current message. uniqueBody is not returned by default but can be retrieved for a given message by use of the ?$select=uniqueBody query. It can be in HTML or text format.
func (m *Message) SetUniqueBody(value ItemBodyable)() {
    err := m.GetBackingStore().Set("uniqueBody", value)
    if err != nil {
        panic(err)
    }
}
// SetWebLink sets the webLink property value. The URL to open the message in Outlook on the web.You can append an ispopout argument to the end of the URL to change how the message is displayed. If ispopout is not present or if it is set to 1, then the message is shown in a popout window. If ispopout is set to 0, the browser shows the message in the Outlook on the web review pane.The message opens in the browser if you are signed in to your mailbox via Outlook on the web. You are prompted to sign in if you are not already signed in with the browser.This URL cannot be accessed from within an iFrame.
func (m *Message) SetWebLink(value *string)() {
    err := m.GetBackingStore().Set("webLink", value)
    if err != nil {
        panic(err)
    }
}
type Messageable interface {
    OutlookItemable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAttachments()([]Attachmentable)
    GetBccRecipients()([]Recipientable)
    GetBody()(ItemBodyable)
    GetBodyPreview()(*string)
    GetCcRecipients()([]Recipientable)
    GetConversationId()(*string)
    GetConversationIndex()([]byte)
    GetExtensions()([]Extensionable)
    GetFlag()(FollowupFlagable)
    GetFrom()(Recipientable)
    GetHasAttachments()(*bool)
    GetImportance()(*Importance)
    GetInferenceClassification()(*InferenceClassificationType)
    GetInternetMessageHeaders()([]InternetMessageHeaderable)
    GetInternetMessageId()(*string)
    GetIsDeliveryReceiptRequested()(*bool)
    GetIsDraft()(*bool)
    GetIsRead()(*bool)
    GetIsReadReceiptRequested()(*bool)
    GetMultiValueExtendedProperties()([]MultiValueLegacyExtendedPropertyable)
    GetParentFolderId()(*string)
    GetReceivedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetReplyTo()([]Recipientable)
    GetSender()(Recipientable)
    GetSentDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetSingleValueExtendedProperties()([]SingleValueLegacyExtendedPropertyable)
    GetSubject()(*string)
    GetToRecipients()([]Recipientable)
    GetUniqueBody()(ItemBodyable)
    GetWebLink()(*string)
    SetAttachments(value []Attachmentable)()
    SetBccRecipients(value []Recipientable)()
    SetBody(value ItemBodyable)()
    SetBodyPreview(value *string)()
    SetCcRecipients(value []Recipientable)()
    SetConversationId(value *string)()
    SetConversationIndex(value []byte)()
    SetExtensions(value []Extensionable)()
    SetFlag(value FollowupFlagable)()
    SetFrom(value Recipientable)()
    SetHasAttachments(value *bool)()
    SetImportance(value *Importance)()
    SetInferenceClassification(value *InferenceClassificationType)()
    SetInternetMessageHeaders(value []InternetMessageHeaderable)()
    SetInternetMessageId(value *string)()
    SetIsDeliveryReceiptRequested(value *bool)()
    SetIsDraft(value *bool)()
    SetIsRead(value *bool)()
    SetIsReadReceiptRequested(value *bool)()
    SetMultiValueExtendedProperties(value []MultiValueLegacyExtendedPropertyable)()
    SetParentFolderId(value *string)()
    SetReceivedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetReplyTo(value []Recipientable)()
    SetSender(value Recipientable)()
    SetSentDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetSingleValueExtendedProperties(value []SingleValueLegacyExtendedPropertyable)()
    SetSubject(value *string)()
    SetToRecipients(value []Recipientable)()
    SetUniqueBody(value ItemBodyable)()
    SetWebLink(value *string)()
}

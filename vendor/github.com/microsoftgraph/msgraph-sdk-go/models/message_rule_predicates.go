package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type MessageRulePredicates struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMessageRulePredicates instantiates a new MessageRulePredicates and sets the default values.
func NewMessageRulePredicates()(*MessageRulePredicates) {
    m := &MessageRulePredicates{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMessageRulePredicatesFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMessageRulePredicatesFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMessageRulePredicates(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *MessageRulePredicates) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *MessageRulePredicates) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBodyContains gets the bodyContains property value. Represents the strings that should appear in the body of an incoming message in order for the condition or exception to apply.
// returns a []string when successful
func (m *MessageRulePredicates) GetBodyContains()([]string) {
    val, err := m.GetBackingStore().Get("bodyContains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetBodyOrSubjectContains gets the bodyOrSubjectContains property value. Represents the strings that should appear in the body or subject of an incoming message in order for the condition or exception to apply.
// returns a []string when successful
func (m *MessageRulePredicates) GetBodyOrSubjectContains()([]string) {
    val, err := m.GetBackingStore().Get("bodyOrSubjectContains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetCategories gets the categories property value. Represents the categories that an incoming message should be labeled with in order for the condition or exception to apply.
// returns a []string when successful
func (m *MessageRulePredicates) GetCategories()([]string) {
    val, err := m.GetBackingStore().Get("categories")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MessageRulePredicates) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["bodyContains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetBodyContains(res)
        }
        return nil
    }
    res["bodyOrSubjectContains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetBodyOrSubjectContains(res)
        }
        return nil
    }
    res["categories"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetCategories(res)
        }
        return nil
    }
    res["fromAddresses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetFromAddresses(res)
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
    res["headerContains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetHeaderContains(res)
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
    res["isApprovalRequest"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsApprovalRequest(val)
        }
        return nil
    }
    res["isAutomaticForward"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAutomaticForward(val)
        }
        return nil
    }
    res["isAutomaticReply"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAutomaticReply(val)
        }
        return nil
    }
    res["isEncrypted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEncrypted(val)
        }
        return nil
    }
    res["isMeetingRequest"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsMeetingRequest(val)
        }
        return nil
    }
    res["isMeetingResponse"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsMeetingResponse(val)
        }
        return nil
    }
    res["isNonDeliveryReport"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsNonDeliveryReport(val)
        }
        return nil
    }
    res["isPermissionControlled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsPermissionControlled(val)
        }
        return nil
    }
    res["isReadReceipt"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsReadReceipt(val)
        }
        return nil
    }
    res["isSigned"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsSigned(val)
        }
        return nil
    }
    res["isVoicemail"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsVoicemail(val)
        }
        return nil
    }
    res["messageActionFlag"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMessageActionFlag)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMessageActionFlag(val.(*MessageActionFlag))
        }
        return nil
    }
    res["notSentToMe"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotSentToMe(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["recipientContains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetRecipientContains(res)
        }
        return nil
    }
    res["senderContains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSenderContains(res)
        }
        return nil
    }
    res["sensitivity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseSensitivity)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSensitivity(val.(*Sensitivity))
        }
        return nil
    }
    res["sentCcMe"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSentCcMe(val)
        }
        return nil
    }
    res["sentOnlyToMe"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSentOnlyToMe(val)
        }
        return nil
    }
    res["sentToAddresses"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSentToAddresses(res)
        }
        return nil
    }
    res["sentToMe"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSentToMe(val)
        }
        return nil
    }
    res["sentToOrCcMe"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSentToOrCcMe(val)
        }
        return nil
    }
    res["subjectContains"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetSubjectContains(res)
        }
        return nil
    }
    res["withinSizeRange"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateSizeRangeFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetWithinSizeRange(val.(SizeRangeable))
        }
        return nil
    }
    return res
}
// GetFromAddresses gets the fromAddresses property value. Represents the specific sender email addresses of an incoming message in order for the condition or exception to apply.
// returns a []Recipientable when successful
func (m *MessageRulePredicates) GetFromAddresses()([]Recipientable) {
    val, err := m.GetBackingStore().Get("fromAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Recipientable)
    }
    return nil
}
// GetHasAttachments gets the hasAttachments property value. Indicates whether an incoming message must have attachments in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetHasAttachments()(*bool) {
    val, err := m.GetBackingStore().Get("hasAttachments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetHeaderContains gets the headerContains property value. Represents the strings that appear in the headers of an incoming message in order for the condition or exception to apply.
// returns a []string when successful
func (m *MessageRulePredicates) GetHeaderContains()([]string) {
    val, err := m.GetBackingStore().Get("headerContains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetImportance gets the importance property value. The importance that is stamped on an incoming message in order for the condition or exception to apply: low, normal, high.
// returns a *Importance when successful
func (m *MessageRulePredicates) GetImportance()(*Importance) {
    val, err := m.GetBackingStore().Get("importance")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Importance)
    }
    return nil
}
// GetIsApprovalRequest gets the isApprovalRequest property value. Indicates whether an incoming message must be an approval request in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsApprovalRequest()(*bool) {
    val, err := m.GetBackingStore().Get("isApprovalRequest")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsAutomaticForward gets the isAutomaticForward property value. Indicates whether an incoming message must be automatically forwarded in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsAutomaticForward()(*bool) {
    val, err := m.GetBackingStore().Get("isAutomaticForward")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsAutomaticReply gets the isAutomaticReply property value. Indicates whether an incoming message must be an auto reply in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsAutomaticReply()(*bool) {
    val, err := m.GetBackingStore().Get("isAutomaticReply")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsEncrypted gets the isEncrypted property value. Indicates whether an incoming message must be encrypted in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsEncrypted()(*bool) {
    val, err := m.GetBackingStore().Get("isEncrypted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsMeetingRequest gets the isMeetingRequest property value. Indicates whether an incoming message must be a meeting request in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsMeetingRequest()(*bool) {
    val, err := m.GetBackingStore().Get("isMeetingRequest")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsMeetingResponse gets the isMeetingResponse property value. Indicates whether an incoming message must be a meeting response in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsMeetingResponse()(*bool) {
    val, err := m.GetBackingStore().Get("isMeetingResponse")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsNonDeliveryReport gets the isNonDeliveryReport property value. Indicates whether an incoming message must be a non-delivery report in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsNonDeliveryReport()(*bool) {
    val, err := m.GetBackingStore().Get("isNonDeliveryReport")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsPermissionControlled gets the isPermissionControlled property value. Indicates whether an incoming message must be permission controlled (RMS-protected) in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsPermissionControlled()(*bool) {
    val, err := m.GetBackingStore().Get("isPermissionControlled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsReadReceipt gets the isReadReceipt property value. Indicates whether an incoming message must be a read receipt in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsReadReceipt()(*bool) {
    val, err := m.GetBackingStore().Get("isReadReceipt")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsSigned gets the isSigned property value. Indicates whether an incoming message must be S/MIME-signed in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsSigned()(*bool) {
    val, err := m.GetBackingStore().Get("isSigned")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsVoicemail gets the isVoicemail property value. Indicates whether an incoming message must be a voice mail in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetIsVoicemail()(*bool) {
    val, err := m.GetBackingStore().Get("isVoicemail")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMessageActionFlag gets the messageActionFlag property value. Represents the flag-for-action value that appears on an incoming message in order for the condition or exception to apply. The possible values are: any, call, doNotForward, followUp, fyi, forward, noResponseNecessary, read, reply, replyToAll, review.
// returns a *MessageActionFlag when successful
func (m *MessageRulePredicates) GetMessageActionFlag()(*MessageActionFlag) {
    val, err := m.GetBackingStore().Get("messageActionFlag")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MessageActionFlag)
    }
    return nil
}
// GetNotSentToMe gets the notSentToMe property value. Indicates whether the owner of the mailbox must not be a recipient of an incoming message in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetNotSentToMe()(*bool) {
    val, err := m.GetBackingStore().Get("notSentToMe")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *MessageRulePredicates) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecipientContains gets the recipientContains property value. Represents the strings that appear in either the toRecipients or ccRecipients properties of an incoming message in order for the condition or exception to apply.
// returns a []string when successful
func (m *MessageRulePredicates) GetRecipientContains()([]string) {
    val, err := m.GetBackingStore().Get("recipientContains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSenderContains gets the senderContains property value. Represents the strings that appear in the from property of an incoming message in order for the condition or exception to apply.
// returns a []string when successful
func (m *MessageRulePredicates) GetSenderContains()([]string) {
    val, err := m.GetBackingStore().Get("senderContains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetSensitivity gets the sensitivity property value. Represents the sensitivity level that must be stamped on an incoming message in order for the condition or exception to apply. The possible values are: normal, personal, private, confidential.
// returns a *Sensitivity when successful
func (m *MessageRulePredicates) GetSensitivity()(*Sensitivity) {
    val, err := m.GetBackingStore().Get("sensitivity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*Sensitivity)
    }
    return nil
}
// GetSentCcMe gets the sentCcMe property value. Indicates whether the owner of the mailbox must be in the ccRecipients property of an incoming message in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetSentCcMe()(*bool) {
    val, err := m.GetBackingStore().Get("sentCcMe")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSentOnlyToMe gets the sentOnlyToMe property value. Indicates whether the owner of the mailbox must be the only recipient in an incoming message in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetSentOnlyToMe()(*bool) {
    val, err := m.GetBackingStore().Get("sentOnlyToMe")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSentToAddresses gets the sentToAddresses property value. Represents the email addresses that an incoming message must have been sent to in order for the condition or exception to apply.
// returns a []Recipientable when successful
func (m *MessageRulePredicates) GetSentToAddresses()([]Recipientable) {
    val, err := m.GetBackingStore().Get("sentToAddresses")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Recipientable)
    }
    return nil
}
// GetSentToMe gets the sentToMe property value. Indicates whether the owner of the mailbox must be in the toRecipients property of an incoming message in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetSentToMe()(*bool) {
    val, err := m.GetBackingStore().Get("sentToMe")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSentToOrCcMe gets the sentToOrCcMe property value. Indicates whether the owner of the mailbox must be in either a toRecipients or ccRecipients property of an incoming message in order for the condition or exception to apply.
// returns a *bool when successful
func (m *MessageRulePredicates) GetSentToOrCcMe()(*bool) {
    val, err := m.GetBackingStore().Get("sentToOrCcMe")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSubjectContains gets the subjectContains property value. Represents the strings that appear in the subject of an incoming message in order for the condition or exception to apply.
// returns a []string when successful
func (m *MessageRulePredicates) GetSubjectContains()([]string) {
    val, err := m.GetBackingStore().Get("subjectContains")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetWithinSizeRange gets the withinSizeRange property value. Represents the minimum and maximum sizes (in kilobytes) that an incoming message must fall in between in order for the condition or exception to apply.
// returns a SizeRangeable when successful
func (m *MessageRulePredicates) GetWithinSizeRange()(SizeRangeable) {
    val, err := m.GetBackingStore().Get("withinSizeRange")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(SizeRangeable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MessageRulePredicates) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetBodyContains() != nil {
        err := writer.WriteCollectionOfStringValues("bodyContains", m.GetBodyContains())
        if err != nil {
            return err
        }
    }
    if m.GetBodyOrSubjectContains() != nil {
        err := writer.WriteCollectionOfStringValues("bodyOrSubjectContains", m.GetBodyOrSubjectContains())
        if err != nil {
            return err
        }
    }
    if m.GetCategories() != nil {
        err := writer.WriteCollectionOfStringValues("categories", m.GetCategories())
        if err != nil {
            return err
        }
    }
    if m.GetFromAddresses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetFromAddresses()))
        for i, v := range m.GetFromAddresses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("fromAddresses", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("hasAttachments", m.GetHasAttachments())
        if err != nil {
            return err
        }
    }
    if m.GetHeaderContains() != nil {
        err := writer.WriteCollectionOfStringValues("headerContains", m.GetHeaderContains())
        if err != nil {
            return err
        }
    }
    if m.GetImportance() != nil {
        cast := (*m.GetImportance()).String()
        err := writer.WriteStringValue("importance", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isApprovalRequest", m.GetIsApprovalRequest())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isAutomaticForward", m.GetIsAutomaticForward())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isAutomaticReply", m.GetIsAutomaticReply())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isEncrypted", m.GetIsEncrypted())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isMeetingRequest", m.GetIsMeetingRequest())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isMeetingResponse", m.GetIsMeetingResponse())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isNonDeliveryReport", m.GetIsNonDeliveryReport())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isPermissionControlled", m.GetIsPermissionControlled())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isReadReceipt", m.GetIsReadReceipt())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isSigned", m.GetIsSigned())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isVoicemail", m.GetIsVoicemail())
        if err != nil {
            return err
        }
    }
    if m.GetMessageActionFlag() != nil {
        cast := (*m.GetMessageActionFlag()).String()
        err := writer.WriteStringValue("messageActionFlag", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("notSentToMe", m.GetNotSentToMe())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    if m.GetRecipientContains() != nil {
        err := writer.WriteCollectionOfStringValues("recipientContains", m.GetRecipientContains())
        if err != nil {
            return err
        }
    }
    if m.GetSenderContains() != nil {
        err := writer.WriteCollectionOfStringValues("senderContains", m.GetSenderContains())
        if err != nil {
            return err
        }
    }
    if m.GetSensitivity() != nil {
        cast := (*m.GetSensitivity()).String()
        err := writer.WriteStringValue("sensitivity", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("sentCcMe", m.GetSentCcMe())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("sentOnlyToMe", m.GetSentOnlyToMe())
        if err != nil {
            return err
        }
    }
    if m.GetSentToAddresses() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetSentToAddresses()))
        for i, v := range m.GetSentToAddresses() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("sentToAddresses", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("sentToMe", m.GetSentToMe())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("sentToOrCcMe", m.GetSentToOrCcMe())
        if err != nil {
            return err
        }
    }
    if m.GetSubjectContains() != nil {
        err := writer.WriteCollectionOfStringValues("subjectContains", m.GetSubjectContains())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("withinSizeRange", m.GetWithinSizeRange())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *MessageRulePredicates) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *MessageRulePredicates) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBodyContains sets the bodyContains property value. Represents the strings that should appear in the body of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetBodyContains(value []string)() {
    err := m.GetBackingStore().Set("bodyContains", value)
    if err != nil {
        panic(err)
    }
}
// SetBodyOrSubjectContains sets the bodyOrSubjectContains property value. Represents the strings that should appear in the body or subject of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetBodyOrSubjectContains(value []string)() {
    err := m.GetBackingStore().Set("bodyOrSubjectContains", value)
    if err != nil {
        panic(err)
    }
}
// SetCategories sets the categories property value. Represents the categories that an incoming message should be labeled with in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetCategories(value []string)() {
    err := m.GetBackingStore().Set("categories", value)
    if err != nil {
        panic(err)
    }
}
// SetFromAddresses sets the fromAddresses property value. Represents the specific sender email addresses of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetFromAddresses(value []Recipientable)() {
    err := m.GetBackingStore().Set("fromAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetHasAttachments sets the hasAttachments property value. Indicates whether an incoming message must have attachments in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetHasAttachments(value *bool)() {
    err := m.GetBackingStore().Set("hasAttachments", value)
    if err != nil {
        panic(err)
    }
}
// SetHeaderContains sets the headerContains property value. Represents the strings that appear in the headers of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetHeaderContains(value []string)() {
    err := m.GetBackingStore().Set("headerContains", value)
    if err != nil {
        panic(err)
    }
}
// SetImportance sets the importance property value. The importance that is stamped on an incoming message in order for the condition or exception to apply: low, normal, high.
func (m *MessageRulePredicates) SetImportance(value *Importance)() {
    err := m.GetBackingStore().Set("importance", value)
    if err != nil {
        panic(err)
    }
}
// SetIsApprovalRequest sets the isApprovalRequest property value. Indicates whether an incoming message must be an approval request in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsApprovalRequest(value *bool)() {
    err := m.GetBackingStore().Set("isApprovalRequest", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAutomaticForward sets the isAutomaticForward property value. Indicates whether an incoming message must be automatically forwarded in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsAutomaticForward(value *bool)() {
    err := m.GetBackingStore().Set("isAutomaticForward", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAutomaticReply sets the isAutomaticReply property value. Indicates whether an incoming message must be an auto reply in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsAutomaticReply(value *bool)() {
    err := m.GetBackingStore().Set("isAutomaticReply", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEncrypted sets the isEncrypted property value. Indicates whether an incoming message must be encrypted in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsEncrypted(value *bool)() {
    err := m.GetBackingStore().Set("isEncrypted", value)
    if err != nil {
        panic(err)
    }
}
// SetIsMeetingRequest sets the isMeetingRequest property value. Indicates whether an incoming message must be a meeting request in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsMeetingRequest(value *bool)() {
    err := m.GetBackingStore().Set("isMeetingRequest", value)
    if err != nil {
        panic(err)
    }
}
// SetIsMeetingResponse sets the isMeetingResponse property value. Indicates whether an incoming message must be a meeting response in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsMeetingResponse(value *bool)() {
    err := m.GetBackingStore().Set("isMeetingResponse", value)
    if err != nil {
        panic(err)
    }
}
// SetIsNonDeliveryReport sets the isNonDeliveryReport property value. Indicates whether an incoming message must be a non-delivery report in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsNonDeliveryReport(value *bool)() {
    err := m.GetBackingStore().Set("isNonDeliveryReport", value)
    if err != nil {
        panic(err)
    }
}
// SetIsPermissionControlled sets the isPermissionControlled property value. Indicates whether an incoming message must be permission controlled (RMS-protected) in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsPermissionControlled(value *bool)() {
    err := m.GetBackingStore().Set("isPermissionControlled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsReadReceipt sets the isReadReceipt property value. Indicates whether an incoming message must be a read receipt in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsReadReceipt(value *bool)() {
    err := m.GetBackingStore().Set("isReadReceipt", value)
    if err != nil {
        panic(err)
    }
}
// SetIsSigned sets the isSigned property value. Indicates whether an incoming message must be S/MIME-signed in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsSigned(value *bool)() {
    err := m.GetBackingStore().Set("isSigned", value)
    if err != nil {
        panic(err)
    }
}
// SetIsVoicemail sets the isVoicemail property value. Indicates whether an incoming message must be a voice mail in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetIsVoicemail(value *bool)() {
    err := m.GetBackingStore().Set("isVoicemail", value)
    if err != nil {
        panic(err)
    }
}
// SetMessageActionFlag sets the messageActionFlag property value. Represents the flag-for-action value that appears on an incoming message in order for the condition or exception to apply. The possible values are: any, call, doNotForward, followUp, fyi, forward, noResponseNecessary, read, reply, replyToAll, review.
func (m *MessageRulePredicates) SetMessageActionFlag(value *MessageActionFlag)() {
    err := m.GetBackingStore().Set("messageActionFlag", value)
    if err != nil {
        panic(err)
    }
}
// SetNotSentToMe sets the notSentToMe property value. Indicates whether the owner of the mailbox must not be a recipient of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetNotSentToMe(value *bool)() {
    err := m.GetBackingStore().Set("notSentToMe", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MessageRulePredicates) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRecipientContains sets the recipientContains property value. Represents the strings that appear in either the toRecipients or ccRecipients properties of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetRecipientContains(value []string)() {
    err := m.GetBackingStore().Set("recipientContains", value)
    if err != nil {
        panic(err)
    }
}
// SetSenderContains sets the senderContains property value. Represents the strings that appear in the from property of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetSenderContains(value []string)() {
    err := m.GetBackingStore().Set("senderContains", value)
    if err != nil {
        panic(err)
    }
}
// SetSensitivity sets the sensitivity property value. Represents the sensitivity level that must be stamped on an incoming message in order for the condition or exception to apply. The possible values are: normal, personal, private, confidential.
func (m *MessageRulePredicates) SetSensitivity(value *Sensitivity)() {
    err := m.GetBackingStore().Set("sensitivity", value)
    if err != nil {
        panic(err)
    }
}
// SetSentCcMe sets the sentCcMe property value. Indicates whether the owner of the mailbox must be in the ccRecipients property of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetSentCcMe(value *bool)() {
    err := m.GetBackingStore().Set("sentCcMe", value)
    if err != nil {
        panic(err)
    }
}
// SetSentOnlyToMe sets the sentOnlyToMe property value. Indicates whether the owner of the mailbox must be the only recipient in an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetSentOnlyToMe(value *bool)() {
    err := m.GetBackingStore().Set("sentOnlyToMe", value)
    if err != nil {
        panic(err)
    }
}
// SetSentToAddresses sets the sentToAddresses property value. Represents the email addresses that an incoming message must have been sent to in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetSentToAddresses(value []Recipientable)() {
    err := m.GetBackingStore().Set("sentToAddresses", value)
    if err != nil {
        panic(err)
    }
}
// SetSentToMe sets the sentToMe property value. Indicates whether the owner of the mailbox must be in the toRecipients property of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetSentToMe(value *bool)() {
    err := m.GetBackingStore().Set("sentToMe", value)
    if err != nil {
        panic(err)
    }
}
// SetSentToOrCcMe sets the sentToOrCcMe property value. Indicates whether the owner of the mailbox must be in either a toRecipients or ccRecipients property of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetSentToOrCcMe(value *bool)() {
    err := m.GetBackingStore().Set("sentToOrCcMe", value)
    if err != nil {
        panic(err)
    }
}
// SetSubjectContains sets the subjectContains property value. Represents the strings that appear in the subject of an incoming message in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetSubjectContains(value []string)() {
    err := m.GetBackingStore().Set("subjectContains", value)
    if err != nil {
        panic(err)
    }
}
// SetWithinSizeRange sets the withinSizeRange property value. Represents the minimum and maximum sizes (in kilobytes) that an incoming message must fall in between in order for the condition or exception to apply.
func (m *MessageRulePredicates) SetWithinSizeRange(value SizeRangeable)() {
    err := m.GetBackingStore().Set("withinSizeRange", value)
    if err != nil {
        panic(err)
    }
}
type MessageRulePredicatesable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBodyContains()([]string)
    GetBodyOrSubjectContains()([]string)
    GetCategories()([]string)
    GetFromAddresses()([]Recipientable)
    GetHasAttachments()(*bool)
    GetHeaderContains()([]string)
    GetImportance()(*Importance)
    GetIsApprovalRequest()(*bool)
    GetIsAutomaticForward()(*bool)
    GetIsAutomaticReply()(*bool)
    GetIsEncrypted()(*bool)
    GetIsMeetingRequest()(*bool)
    GetIsMeetingResponse()(*bool)
    GetIsNonDeliveryReport()(*bool)
    GetIsPermissionControlled()(*bool)
    GetIsReadReceipt()(*bool)
    GetIsSigned()(*bool)
    GetIsVoicemail()(*bool)
    GetMessageActionFlag()(*MessageActionFlag)
    GetNotSentToMe()(*bool)
    GetOdataType()(*string)
    GetRecipientContains()([]string)
    GetSenderContains()([]string)
    GetSensitivity()(*Sensitivity)
    GetSentCcMe()(*bool)
    GetSentOnlyToMe()(*bool)
    GetSentToAddresses()([]Recipientable)
    GetSentToMe()(*bool)
    GetSentToOrCcMe()(*bool)
    GetSubjectContains()([]string)
    GetWithinSizeRange()(SizeRangeable)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBodyContains(value []string)()
    SetBodyOrSubjectContains(value []string)()
    SetCategories(value []string)()
    SetFromAddresses(value []Recipientable)()
    SetHasAttachments(value *bool)()
    SetHeaderContains(value []string)()
    SetImportance(value *Importance)()
    SetIsApprovalRequest(value *bool)()
    SetIsAutomaticForward(value *bool)()
    SetIsAutomaticReply(value *bool)()
    SetIsEncrypted(value *bool)()
    SetIsMeetingRequest(value *bool)()
    SetIsMeetingResponse(value *bool)()
    SetIsNonDeliveryReport(value *bool)()
    SetIsPermissionControlled(value *bool)()
    SetIsReadReceipt(value *bool)()
    SetIsSigned(value *bool)()
    SetIsVoicemail(value *bool)()
    SetMessageActionFlag(value *MessageActionFlag)()
    SetNotSentToMe(value *bool)()
    SetOdataType(value *string)()
    SetRecipientContains(value []string)()
    SetSenderContains(value []string)()
    SetSensitivity(value *Sensitivity)()
    SetSentCcMe(value *bool)()
    SetSentOnlyToMe(value *bool)()
    SetSentToAddresses(value []Recipientable)()
    SetSentToMe(value *bool)()
    SetSentToOrCcMe(value *bool)()
    SetSubjectContains(value []string)()
    SetWithinSizeRange(value SizeRangeable)()
}

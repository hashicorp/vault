package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type MailTips struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewMailTips instantiates a new MailTips and sets the default values.
func NewMailTips()(*MailTips) {
    m := &MailTips{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateMailTipsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMailTipsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMailTips(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *MailTips) GetAdditionalData()(map[string]any) {
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
// GetAutomaticReplies gets the automaticReplies property value. Mail tips for automatic reply if it has been set up by the recipient.
// returns a AutomaticRepliesMailTipsable when successful
func (m *MailTips) GetAutomaticReplies()(AutomaticRepliesMailTipsable) {
    val, err := m.GetBackingStore().Get("automaticReplies")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(AutomaticRepliesMailTipsable)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *MailTips) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetCustomMailTip gets the customMailTip property value. A custom mail tip that can be set on the recipient's mailbox.
// returns a *string when successful
func (m *MailTips) GetCustomMailTip()(*string) {
    val, err := m.GetBackingStore().Get("customMailTip")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDeliveryRestricted gets the deliveryRestricted property value. Whether the recipient's mailbox is restricted, for example, accepting messages from only a predefined list of senders, rejecting messages from a predefined list of senders, or accepting messages from only authenticated senders.
// returns a *bool when successful
func (m *MailTips) GetDeliveryRestricted()(*bool) {
    val, err := m.GetBackingStore().Get("deliveryRestricted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetEmailAddress gets the emailAddress property value. The email address of the recipient to get mailtips for.
// returns a EmailAddressable when successful
func (m *MailTips) GetEmailAddress()(EmailAddressable) {
    val, err := m.GetBackingStore().Get("emailAddress")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(EmailAddressable)
    }
    return nil
}
// GetError gets the error property value. Errors that occur during the getMailTips action.
// returns a MailTipsErrorable when successful
func (m *MailTips) GetError()(MailTipsErrorable) {
    val, err := m.GetBackingStore().Get("error")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MailTipsErrorable)
    }
    return nil
}
// GetExternalMemberCount gets the externalMemberCount property value. The number of external members if the recipient is a distribution list.
// returns a *int32 when successful
func (m *MailTips) GetExternalMemberCount()(*int32) {
    val, err := m.GetBackingStore().Get("externalMemberCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MailTips) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["automaticReplies"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateAutomaticRepliesMailTipsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAutomaticReplies(val.(AutomaticRepliesMailTipsable))
        }
        return nil
    }
    res["customMailTip"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCustomMailTip(val)
        }
        return nil
    }
    res["deliveryRestricted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDeliveryRestricted(val)
        }
        return nil
    }
    res["emailAddress"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEmailAddressFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEmailAddress(val.(EmailAddressable))
        }
        return nil
    }
    res["error"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMailTipsErrorFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetError(val.(MailTipsErrorable))
        }
        return nil
    }
    res["externalMemberCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExternalMemberCount(val)
        }
        return nil
    }
    res["isModerated"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsModerated(val)
        }
        return nil
    }
    res["mailboxFull"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMailboxFull(val)
        }
        return nil
    }
    res["maxMessageSize"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaxMessageSize(val)
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
    res["recipientScope"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseRecipientScopeType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecipientScope(val.(*RecipientScopeType))
        }
        return nil
    }
    res["recipientSuggestions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetRecipientSuggestions(res)
        }
        return nil
    }
    res["totalMemberCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTotalMemberCount(val)
        }
        return nil
    }
    return res
}
// GetIsModerated gets the isModerated property value. Whether sending messages to the recipient requires approval. For example, if the recipient is a large distribution list and a moderator has been set up to approve messages sent to that distribution list, or if sending messages to a recipient requires approval of the recipient's manager.
// returns a *bool when successful
func (m *MailTips) GetIsModerated()(*bool) {
    val, err := m.GetBackingStore().Get("isModerated")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMailboxFull gets the mailboxFull property value. The mailbox full status of the recipient.
// returns a *bool when successful
func (m *MailTips) GetMailboxFull()(*bool) {
    val, err := m.GetBackingStore().Get("mailboxFull")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMaxMessageSize gets the maxMessageSize property value. The maximum message size that has been configured for the recipient's organization or mailbox.
// returns a *int32 when successful
func (m *MailTips) GetMaxMessageSize()(*int32) {
    val, err := m.GetBackingStore().Get("maxMessageSize")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *MailTips) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecipientScope gets the recipientScope property value. The scope of the recipient. Possible values are: none, internal, external, externalPartner, externalNonParther. For example, an administrator can set another organization to be its 'partner'. The scope is useful if an administrator wants certain mailtips to be accessible to certain scopes. It's also useful to senders to inform them that their message may leave the organization, helping them make the correct decisions about wording, tone and content.
// returns a *RecipientScopeType when successful
func (m *MailTips) GetRecipientScope()(*RecipientScopeType) {
    val, err := m.GetBackingStore().Get("recipientScope")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*RecipientScopeType)
    }
    return nil
}
// GetRecipientSuggestions gets the recipientSuggestions property value. Recipients suggested based on previous contexts where they appear in the same message.
// returns a []Recipientable when successful
func (m *MailTips) GetRecipientSuggestions()([]Recipientable) {
    val, err := m.GetBackingStore().Get("recipientSuggestions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Recipientable)
    }
    return nil
}
// GetTotalMemberCount gets the totalMemberCount property value. The number of members if the recipient is a distribution list.
// returns a *int32 when successful
func (m *MailTips) GetTotalMemberCount()(*int32) {
    val, err := m.GetBackingStore().Get("totalMemberCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MailTips) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteObjectValue("automaticReplies", m.GetAutomaticReplies())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("customMailTip", m.GetCustomMailTip())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("deliveryRestricted", m.GetDeliveryRestricted())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("emailAddress", m.GetEmailAddress())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteObjectValue("error", m.GetError())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("externalMemberCount", m.GetExternalMemberCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isModerated", m.GetIsModerated())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("mailboxFull", m.GetMailboxFull())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("maxMessageSize", m.GetMaxMessageSize())
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
    if m.GetRecipientScope() != nil {
        cast := (*m.GetRecipientScope()).String()
        err := writer.WriteStringValue("recipientScope", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetRecipientSuggestions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRecipientSuggestions()))
        for i, v := range m.GetRecipientSuggestions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err := writer.WriteCollectionOfObjectValues("recipientSuggestions", cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("totalMemberCount", m.GetTotalMemberCount())
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
func (m *MailTips) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAutomaticReplies sets the automaticReplies property value. Mail tips for automatic reply if it has been set up by the recipient.
func (m *MailTips) SetAutomaticReplies(value AutomaticRepliesMailTipsable)() {
    err := m.GetBackingStore().Set("automaticReplies", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *MailTips) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetCustomMailTip sets the customMailTip property value. A custom mail tip that can be set on the recipient's mailbox.
func (m *MailTips) SetCustomMailTip(value *string)() {
    err := m.GetBackingStore().Set("customMailTip", value)
    if err != nil {
        panic(err)
    }
}
// SetDeliveryRestricted sets the deliveryRestricted property value. Whether the recipient's mailbox is restricted, for example, accepting messages from only a predefined list of senders, rejecting messages from a predefined list of senders, or accepting messages from only authenticated senders.
func (m *MailTips) SetDeliveryRestricted(value *bool)() {
    err := m.GetBackingStore().Set("deliveryRestricted", value)
    if err != nil {
        panic(err)
    }
}
// SetEmailAddress sets the emailAddress property value. The email address of the recipient to get mailtips for.
func (m *MailTips) SetEmailAddress(value EmailAddressable)() {
    err := m.GetBackingStore().Set("emailAddress", value)
    if err != nil {
        panic(err)
    }
}
// SetError sets the error property value. Errors that occur during the getMailTips action.
func (m *MailTips) SetError(value MailTipsErrorable)() {
    err := m.GetBackingStore().Set("error", value)
    if err != nil {
        panic(err)
    }
}
// SetExternalMemberCount sets the externalMemberCount property value. The number of external members if the recipient is a distribution list.
func (m *MailTips) SetExternalMemberCount(value *int32)() {
    err := m.GetBackingStore().Set("externalMemberCount", value)
    if err != nil {
        panic(err)
    }
}
// SetIsModerated sets the isModerated property value. Whether sending messages to the recipient requires approval. For example, if the recipient is a large distribution list and a moderator has been set up to approve messages sent to that distribution list, or if sending messages to a recipient requires approval of the recipient's manager.
func (m *MailTips) SetIsModerated(value *bool)() {
    err := m.GetBackingStore().Set("isModerated", value)
    if err != nil {
        panic(err)
    }
}
// SetMailboxFull sets the mailboxFull property value. The mailbox full status of the recipient.
func (m *MailTips) SetMailboxFull(value *bool)() {
    err := m.GetBackingStore().Set("mailboxFull", value)
    if err != nil {
        panic(err)
    }
}
// SetMaxMessageSize sets the maxMessageSize property value. The maximum message size that has been configured for the recipient's organization or mailbox.
func (m *MailTips) SetMaxMessageSize(value *int32)() {
    err := m.GetBackingStore().Set("maxMessageSize", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *MailTips) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetRecipientScope sets the recipientScope property value. The scope of the recipient. Possible values are: none, internal, external, externalPartner, externalNonParther. For example, an administrator can set another organization to be its 'partner'. The scope is useful if an administrator wants certain mailtips to be accessible to certain scopes. It's also useful to senders to inform them that their message may leave the organization, helping them make the correct decisions about wording, tone and content.
func (m *MailTips) SetRecipientScope(value *RecipientScopeType)() {
    err := m.GetBackingStore().Set("recipientScope", value)
    if err != nil {
        panic(err)
    }
}
// SetRecipientSuggestions sets the recipientSuggestions property value. Recipients suggested based on previous contexts where they appear in the same message.
func (m *MailTips) SetRecipientSuggestions(value []Recipientable)() {
    err := m.GetBackingStore().Set("recipientSuggestions", value)
    if err != nil {
        panic(err)
    }
}
// SetTotalMemberCount sets the totalMemberCount property value. The number of members if the recipient is a distribution list.
func (m *MailTips) SetTotalMemberCount(value *int32)() {
    err := m.GetBackingStore().Set("totalMemberCount", value)
    if err != nil {
        panic(err)
    }
}
type MailTipsable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAutomaticReplies()(AutomaticRepliesMailTipsable)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetCustomMailTip()(*string)
    GetDeliveryRestricted()(*bool)
    GetEmailAddress()(EmailAddressable)
    GetError()(MailTipsErrorable)
    GetExternalMemberCount()(*int32)
    GetIsModerated()(*bool)
    GetMailboxFull()(*bool)
    GetMaxMessageSize()(*int32)
    GetOdataType()(*string)
    GetRecipientScope()(*RecipientScopeType)
    GetRecipientSuggestions()([]Recipientable)
    GetTotalMemberCount()(*int32)
    SetAutomaticReplies(value AutomaticRepliesMailTipsable)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetCustomMailTip(value *string)()
    SetDeliveryRestricted(value *bool)()
    SetEmailAddress(value EmailAddressable)()
    SetError(value MailTipsErrorable)()
    SetExternalMemberCount(value *int32)()
    SetIsModerated(value *bool)()
    SetMailboxFull(value *bool)()
    SetMaxMessageSize(value *int32)()
    SetOdataType(value *string)()
    SetRecipientScope(value *RecipientScopeType)()
    SetRecipientSuggestions(value []Recipientable)()
    SetTotalMemberCount(value *int32)()
}

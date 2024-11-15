package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type ChatMessagePolicyViolation struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewChatMessagePolicyViolation instantiates a new ChatMessagePolicyViolation and sets the default values.
func NewChatMessagePolicyViolation()(*ChatMessagePolicyViolation) {
    m := &ChatMessagePolicyViolation{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateChatMessagePolicyViolationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateChatMessagePolicyViolationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewChatMessagePolicyViolation(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *ChatMessagePolicyViolation) GetAdditionalData()(map[string]any) {
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
func (m *ChatMessagePolicyViolation) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetDlpAction gets the dlpAction property value. The action taken by the DLP provider on the message with sensitive content. Supported values are: NoneNotifySender -- Inform the sender of the violation but allow readers to read the message.BlockAccess -- Block readers from reading the message.BlockAccessExternal -- Block users outside the organization from reading the message, while allowing users within the organization to read the message.
// returns a *ChatMessagePolicyViolationDlpActionTypes when successful
func (m *ChatMessagePolicyViolation) GetDlpAction()(*ChatMessagePolicyViolationDlpActionTypes) {
    val, err := m.GetBackingStore().Get("dlpAction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ChatMessagePolicyViolationDlpActionTypes)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ChatMessagePolicyViolation) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["dlpAction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseChatMessagePolicyViolationDlpActionTypes)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDlpAction(val.(*ChatMessagePolicyViolationDlpActionTypes))
        }
        return nil
    }
    res["justificationText"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetJustificationText(val)
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
    res["policyTip"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateChatMessagePolicyViolationPolicyTipFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyTip(val.(ChatMessagePolicyViolationPolicyTipable))
        }
        return nil
    }
    res["userAction"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseChatMessagePolicyViolationUserActionTypes)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserAction(val.(*ChatMessagePolicyViolationUserActionTypes))
        }
        return nil
    }
    res["verdictDetails"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseChatMessagePolicyViolationVerdictDetailsTypes)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVerdictDetails(val.(*ChatMessagePolicyViolationVerdictDetailsTypes))
        }
        return nil
    }
    return res
}
// GetJustificationText gets the justificationText property value. Justification text provided by the sender of the message when overriding a policy violation.
// returns a *string when successful
func (m *ChatMessagePolicyViolation) GetJustificationText()(*string) {
    val, err := m.GetBackingStore().Get("justificationText")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *ChatMessagePolicyViolation) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPolicyTip gets the policyTip property value. Information to display to the message sender about why the message was flagged as a violation.
// returns a ChatMessagePolicyViolationPolicyTipable when successful
func (m *ChatMessagePolicyViolation) GetPolicyTip()(ChatMessagePolicyViolationPolicyTipable) {
    val, err := m.GetBackingStore().Get("policyTip")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ChatMessagePolicyViolationPolicyTipable)
    }
    return nil
}
// GetUserAction gets the userAction property value. Indicates the action taken by the user on a message blocked by the DLP provider. Supported values are: NoneOverrideReportFalsePositiveWhen the DLP provider is updating the message for blocking sensitive content, userAction isn't required.
// returns a *ChatMessagePolicyViolationUserActionTypes when successful
func (m *ChatMessagePolicyViolation) GetUserAction()(*ChatMessagePolicyViolationUserActionTypes) {
    val, err := m.GetBackingStore().Get("userAction")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ChatMessagePolicyViolationUserActionTypes)
    }
    return nil
}
// GetVerdictDetails gets the verdictDetails property value. Indicates what actions the sender may take in response to the policy violation. Supported values are: NoneAllowFalsePositiveOverride -- Allows the sender to declare the policyViolation to be an error in the DLP app and its rules, and allow readers to see the message again if the dlpAction hides it.AllowOverrideWithoutJustification -- Allows the sender to override the DLP violation and allow readers to see the message again if the dlpAction hides it, without needing to provide an explanation for doing so. AllowOverrideWithJustification -- Allows the sender to override the DLP violation and allow readers to see the message again if the dlpAction hides it, after providing an explanation for doing so.AllowOverrideWithoutJustification and AllowOverrideWithJustification are mutually exclusive.
// returns a *ChatMessagePolicyViolationVerdictDetailsTypes when successful
func (m *ChatMessagePolicyViolation) GetVerdictDetails()(*ChatMessagePolicyViolationVerdictDetailsTypes) {
    val, err := m.GetBackingStore().Get("verdictDetails")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*ChatMessagePolicyViolationVerdictDetailsTypes)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ChatMessagePolicyViolation) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    if m.GetDlpAction() != nil {
        cast := (*m.GetDlpAction()).String()
        err := writer.WriteStringValue("dlpAction", &cast)
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("justificationText", m.GetJustificationText())
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
    {
        err := writer.WriteObjectValue("policyTip", m.GetPolicyTip())
        if err != nil {
            return err
        }
    }
    if m.GetUserAction() != nil {
        cast := (*m.GetUserAction()).String()
        err := writer.WriteStringValue("userAction", &cast)
        if err != nil {
            return err
        }
    }
    if m.GetVerdictDetails() != nil {
        cast := (*m.GetVerdictDetails()).String()
        err := writer.WriteStringValue("verdictDetails", &cast)
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
func (m *ChatMessagePolicyViolation) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *ChatMessagePolicyViolation) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetDlpAction sets the dlpAction property value. The action taken by the DLP provider on the message with sensitive content. Supported values are: NoneNotifySender -- Inform the sender of the violation but allow readers to read the message.BlockAccess -- Block readers from reading the message.BlockAccessExternal -- Block users outside the organization from reading the message, while allowing users within the organization to read the message.
func (m *ChatMessagePolicyViolation) SetDlpAction(value *ChatMessagePolicyViolationDlpActionTypes)() {
    err := m.GetBackingStore().Set("dlpAction", value)
    if err != nil {
        panic(err)
    }
}
// SetJustificationText sets the justificationText property value. Justification text provided by the sender of the message when overriding a policy violation.
func (m *ChatMessagePolicyViolation) SetJustificationText(value *string)() {
    err := m.GetBackingStore().Set("justificationText", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *ChatMessagePolicyViolation) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyTip sets the policyTip property value. Information to display to the message sender about why the message was flagged as a violation.
func (m *ChatMessagePolicyViolation) SetPolicyTip(value ChatMessagePolicyViolationPolicyTipable)() {
    err := m.GetBackingStore().Set("policyTip", value)
    if err != nil {
        panic(err)
    }
}
// SetUserAction sets the userAction property value. Indicates the action taken by the user on a message blocked by the DLP provider. Supported values are: NoneOverrideReportFalsePositiveWhen the DLP provider is updating the message for blocking sensitive content, userAction isn't required.
func (m *ChatMessagePolicyViolation) SetUserAction(value *ChatMessagePolicyViolationUserActionTypes)() {
    err := m.GetBackingStore().Set("userAction", value)
    if err != nil {
        panic(err)
    }
}
// SetVerdictDetails sets the verdictDetails property value. Indicates what actions the sender may take in response to the policy violation. Supported values are: NoneAllowFalsePositiveOverride -- Allows the sender to declare the policyViolation to be an error in the DLP app and its rules, and allow readers to see the message again if the dlpAction hides it.AllowOverrideWithoutJustification -- Allows the sender to override the DLP violation and allow readers to see the message again if the dlpAction hides it, without needing to provide an explanation for doing so. AllowOverrideWithJustification -- Allows the sender to override the DLP violation and allow readers to see the message again if the dlpAction hides it, after providing an explanation for doing so.AllowOverrideWithoutJustification and AllowOverrideWithJustification are mutually exclusive.
func (m *ChatMessagePolicyViolation) SetVerdictDetails(value *ChatMessagePolicyViolationVerdictDetailsTypes)() {
    err := m.GetBackingStore().Set("verdictDetails", value)
    if err != nil {
        panic(err)
    }
}
type ChatMessagePolicyViolationable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetDlpAction()(*ChatMessagePolicyViolationDlpActionTypes)
    GetJustificationText()(*string)
    GetOdataType()(*string)
    GetPolicyTip()(ChatMessagePolicyViolationPolicyTipable)
    GetUserAction()(*ChatMessagePolicyViolationUserActionTypes)
    GetVerdictDetails()(*ChatMessagePolicyViolationVerdictDetailsTypes)
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetDlpAction(value *ChatMessagePolicyViolationDlpActionTypes)()
    SetJustificationText(value *string)()
    SetOdataType(value *string)()
    SetPolicyTip(value ChatMessagePolicyViolationPolicyTipable)()
    SetUserAction(value *ChatMessagePolicyViolationUserActionTypes)()
    SetVerdictDetails(value *ChatMessagePolicyViolationVerdictDetailsTypes)()
}

package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MessageRule struct {
    Entity
}
// NewMessageRule instantiates a new MessageRule and sets the default values.
func NewMessageRule()(*MessageRule) {
    m := &MessageRule{
        Entity: *NewEntity(),
    }
    return m
}
// CreateMessageRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMessageRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMessageRule(), nil
}
// GetActions gets the actions property value. Actions to be taken on a message when the corresponding conditions are fulfilled.
// returns a MessageRuleActionsable when successful
func (m *MessageRule) GetActions()(MessageRuleActionsable) {
    val, err := m.GetBackingStore().Get("actions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MessageRuleActionsable)
    }
    return nil
}
// GetConditions gets the conditions property value. Conditions that when fulfilled trigger the corresponding actions for that rule.
// returns a MessageRulePredicatesable when successful
func (m *MessageRule) GetConditions()(MessageRulePredicatesable) {
    val, err := m.GetBackingStore().Get("conditions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MessageRulePredicatesable)
    }
    return nil
}
// GetDisplayName gets the displayName property value. The display name of the rule.
// returns a *string when successful
func (m *MessageRule) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExceptions gets the exceptions property value. Exception conditions for the rule.
// returns a MessageRulePredicatesable when successful
func (m *MessageRule) GetExceptions()(MessageRulePredicatesable) {
    val, err := m.GetBackingStore().Get("exceptions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(MessageRulePredicatesable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MessageRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["actions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMessageRuleActionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActions(val.(MessageRuleActionsable))
        }
        return nil
    }
    res["conditions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMessageRulePredicatesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetConditions(val.(MessageRulePredicatesable))
        }
        return nil
    }
    res["displayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisplayName(val)
        }
        return nil
    }
    res["exceptions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateMessageRulePredicatesFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetExceptions(val.(MessageRulePredicatesable))
        }
        return nil
    }
    res["hasError"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHasError(val)
        }
        return nil
    }
    res["isEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsEnabled(val)
        }
        return nil
    }
    res["isReadOnly"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsReadOnly(val)
        }
        return nil
    }
    res["sequence"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSequence(val)
        }
        return nil
    }
    return res
}
// GetHasError gets the hasError property value. Indicates whether the rule is in an error condition. Read-only.
// returns a *bool when successful
func (m *MessageRule) GetHasError()(*bool) {
    val, err := m.GetBackingStore().Get("hasError")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsEnabled gets the isEnabled property value. Indicates whether the rule is enabled to be applied to messages.
// returns a *bool when successful
func (m *MessageRule) GetIsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsReadOnly gets the isReadOnly property value. Indicates if the rule is read-only and cannot be modified or deleted by the rules REST API.
// returns a *bool when successful
func (m *MessageRule) GetIsReadOnly()(*bool) {
    val, err := m.GetBackingStore().Get("isReadOnly")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetSequence gets the sequence property value. Indicates the order in which the rule is executed, among other rules.
// returns a *int32 when successful
func (m *MessageRule) GetSequence()(*int32) {
    val, err := m.GetBackingStore().Get("sequence")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MessageRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("actions", m.GetActions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("conditions", m.GetConditions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("exceptions", m.GetExceptions())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("hasError", m.GetHasError())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isEnabled", m.GetIsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isReadOnly", m.GetIsReadOnly())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt32Value("sequence", m.GetSequence())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActions sets the actions property value. Actions to be taken on a message when the corresponding conditions are fulfilled.
func (m *MessageRule) SetActions(value MessageRuleActionsable)() {
    err := m.GetBackingStore().Set("actions", value)
    if err != nil {
        panic(err)
    }
}
// SetConditions sets the conditions property value. Conditions that when fulfilled trigger the corresponding actions for that rule.
func (m *MessageRule) SetConditions(value MessageRulePredicatesable)() {
    err := m.GetBackingStore().Set("conditions", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. The display name of the rule.
func (m *MessageRule) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExceptions sets the exceptions property value. Exception conditions for the rule.
func (m *MessageRule) SetExceptions(value MessageRulePredicatesable)() {
    err := m.GetBackingStore().Set("exceptions", value)
    if err != nil {
        panic(err)
    }
}
// SetHasError sets the hasError property value. Indicates whether the rule is in an error condition. Read-only.
func (m *MessageRule) SetHasError(value *bool)() {
    err := m.GetBackingStore().Set("hasError", value)
    if err != nil {
        panic(err)
    }
}
// SetIsEnabled sets the isEnabled property value. Indicates whether the rule is enabled to be applied to messages.
func (m *MessageRule) SetIsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetIsReadOnly sets the isReadOnly property value. Indicates if the rule is read-only and cannot be modified or deleted by the rules REST API.
func (m *MessageRule) SetIsReadOnly(value *bool)() {
    err := m.GetBackingStore().Set("isReadOnly", value)
    if err != nil {
        panic(err)
    }
}
// SetSequence sets the sequence property value. Indicates the order in which the rule is executed, among other rules.
func (m *MessageRule) SetSequence(value *int32)() {
    err := m.GetBackingStore().Set("sequence", value)
    if err != nil {
        panic(err)
    }
}
type MessageRuleable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActions()(MessageRuleActionsable)
    GetConditions()(MessageRulePredicatesable)
    GetDisplayName()(*string)
    GetExceptions()(MessageRulePredicatesable)
    GetHasError()(*bool)
    GetIsEnabled()(*bool)
    GetIsReadOnly()(*bool)
    GetSequence()(*int32)
    SetActions(value MessageRuleActionsable)()
    SetConditions(value MessageRulePredicatesable)()
    SetDisplayName(value *string)()
    SetExceptions(value MessageRulePredicatesable)()
    SetHasError(value *bool)()
    SetIsEnabled(value *bool)()
    SetIsReadOnly(value *bool)()
    SetSequence(value *int32)()
}

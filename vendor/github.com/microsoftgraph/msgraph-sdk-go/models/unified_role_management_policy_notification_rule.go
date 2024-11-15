package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnifiedRoleManagementPolicyNotificationRule struct {
    UnifiedRoleManagementPolicyRule
}
// NewUnifiedRoleManagementPolicyNotificationRule instantiates a new UnifiedRoleManagementPolicyNotificationRule and sets the default values.
func NewUnifiedRoleManagementPolicyNotificationRule()(*UnifiedRoleManagementPolicyNotificationRule) {
    m := &UnifiedRoleManagementPolicyNotificationRule{
        UnifiedRoleManagementPolicyRule: *NewUnifiedRoleManagementPolicyRule(),
    }
    odataTypeValue := "#microsoft.graph.unifiedRoleManagementPolicyNotificationRule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateUnifiedRoleManagementPolicyNotificationRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedRoleManagementPolicyNotificationRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleManagementPolicyNotificationRule(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UnifiedRoleManagementPolicyNotificationRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.UnifiedRoleManagementPolicyRule.GetFieldDeserializers()
    res["isDefaultRecipientsEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDefaultRecipientsEnabled(val)
        }
        return nil
    }
    res["notificationLevel"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotificationLevel(val)
        }
        return nil
    }
    res["notificationRecipients"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetNotificationRecipients(res)
        }
        return nil
    }
    res["notificationType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetNotificationType(val)
        }
        return nil
    }
    res["recipientType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecipientType(val)
        }
        return nil
    }
    return res
}
// GetIsDefaultRecipientsEnabled gets the isDefaultRecipientsEnabled property value. Indicates whether a default recipient will receive the notification email.
// returns a *bool when successful
func (m *UnifiedRoleManagementPolicyNotificationRule) GetIsDefaultRecipientsEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("isDefaultRecipientsEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetNotificationLevel gets the notificationLevel property value. The level of notification. The possible values are None, Critical, All.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicyNotificationRule) GetNotificationLevel()(*string) {
    val, err := m.GetBackingStore().Get("notificationLevel")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetNotificationRecipients gets the notificationRecipients property value. The list of recipients of the email notifications.
// returns a []string when successful
func (m *UnifiedRoleManagementPolicyNotificationRule) GetNotificationRecipients()([]string) {
    val, err := m.GetBackingStore().Get("notificationRecipients")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]string)
    }
    return nil
}
// GetNotificationType gets the notificationType property value. The type of notification. Only Email is supported.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicyNotificationRule) GetNotificationType()(*string) {
    val, err := m.GetBackingStore().Get("notificationType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRecipientType gets the recipientType property value. The type of recipient of the notification. The possible values are Requestor, Approver, Admin.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicyNotificationRule) GetRecipientType()(*string) {
    val, err := m.GetBackingStore().Get("recipientType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UnifiedRoleManagementPolicyNotificationRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.UnifiedRoleManagementPolicyRule.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isDefaultRecipientsEnabled", m.GetIsDefaultRecipientsEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notificationLevel", m.GetNotificationLevel())
        if err != nil {
            return err
        }
    }
    if m.GetNotificationRecipients() != nil {
        err = writer.WriteCollectionOfStringValues("notificationRecipients", m.GetNotificationRecipients())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("notificationType", m.GetNotificationType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("recipientType", m.GetRecipientType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsDefaultRecipientsEnabled sets the isDefaultRecipientsEnabled property value. Indicates whether a default recipient will receive the notification email.
func (m *UnifiedRoleManagementPolicyNotificationRule) SetIsDefaultRecipientsEnabled(value *bool)() {
    err := m.GetBackingStore().Set("isDefaultRecipientsEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetNotificationLevel sets the notificationLevel property value. The level of notification. The possible values are None, Critical, All.
func (m *UnifiedRoleManagementPolicyNotificationRule) SetNotificationLevel(value *string)() {
    err := m.GetBackingStore().Set("notificationLevel", value)
    if err != nil {
        panic(err)
    }
}
// SetNotificationRecipients sets the notificationRecipients property value. The list of recipients of the email notifications.
func (m *UnifiedRoleManagementPolicyNotificationRule) SetNotificationRecipients(value []string)() {
    err := m.GetBackingStore().Set("notificationRecipients", value)
    if err != nil {
        panic(err)
    }
}
// SetNotificationType sets the notificationType property value. The type of notification. Only Email is supported.
func (m *UnifiedRoleManagementPolicyNotificationRule) SetNotificationType(value *string)() {
    err := m.GetBackingStore().Set("notificationType", value)
    if err != nil {
        panic(err)
    }
}
// SetRecipientType sets the recipientType property value. The type of recipient of the notification. The possible values are Requestor, Approver, Admin.
func (m *UnifiedRoleManagementPolicyNotificationRule) SetRecipientType(value *string)() {
    err := m.GetBackingStore().Set("recipientType", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedRoleManagementPolicyNotificationRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    UnifiedRoleManagementPolicyRuleable
    GetIsDefaultRecipientsEnabled()(*bool)
    GetNotificationLevel()(*string)
    GetNotificationRecipients()([]string)
    GetNotificationType()(*string)
    GetRecipientType()(*string)
    SetIsDefaultRecipientsEnabled(value *bool)()
    SetNotificationLevel(value *string)()
    SetNotificationRecipients(value []string)()
    SetNotificationType(value *string)()
    SetRecipientType(value *string)()
}

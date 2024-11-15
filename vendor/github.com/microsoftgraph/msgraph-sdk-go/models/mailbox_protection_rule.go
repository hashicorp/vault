package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MailboxProtectionRule struct {
    ProtectionRuleBase
}
// NewMailboxProtectionRule instantiates a new MailboxProtectionRule and sets the default values.
func NewMailboxProtectionRule()(*MailboxProtectionRule) {
    m := &MailboxProtectionRule{
        ProtectionRuleBase: *NewProtectionRuleBase(),
    }
    odataTypeValue := "#microsoft.graph.mailboxProtectionRule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMailboxProtectionRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMailboxProtectionRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMailboxProtectionRule(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MailboxProtectionRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ProtectionRuleBase.GetFieldDeserializers()
    res["mailboxExpression"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMailboxExpression(val)
        }
        return nil
    }
    return res
}
// GetMailboxExpression gets the mailboxExpression property value. Contains a mailbox expression. For examples, see mailboxExpression examples.
// returns a *string when successful
func (m *MailboxProtectionRule) GetMailboxExpression()(*string) {
    val, err := m.GetBackingStore().Get("mailboxExpression")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MailboxProtectionRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ProtectionRuleBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("mailboxExpression", m.GetMailboxExpression())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMailboxExpression sets the mailboxExpression property value. Contains a mailbox expression. For examples, see mailboxExpression examples.
func (m *MailboxProtectionRule) SetMailboxExpression(value *string)() {
    err := m.GetBackingStore().Set("mailboxExpression", value)
    if err != nil {
        panic(err)
    }
}
type MailboxProtectionRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ProtectionRuleBaseable
    GetMailboxExpression()(*string)
    SetMailboxExpression(value *string)()
}

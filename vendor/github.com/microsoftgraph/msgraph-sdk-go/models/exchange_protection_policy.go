package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ExchangeProtectionPolicy struct {
    ProtectionPolicyBase
}
// NewExchangeProtectionPolicy instantiates a new ExchangeProtectionPolicy and sets the default values.
func NewExchangeProtectionPolicy()(*ExchangeProtectionPolicy) {
    m := &ExchangeProtectionPolicy{
        ProtectionPolicyBase: *NewProtectionPolicyBase(),
    }
    odataTypeValue := "#microsoft.graph.exchangeProtectionPolicy"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateExchangeProtectionPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateExchangeProtectionPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewExchangeProtectionPolicy(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ExchangeProtectionPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ProtectionPolicyBase.GetFieldDeserializers()
    res["mailboxInclusionRules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMailboxProtectionRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MailboxProtectionRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MailboxProtectionRuleable)
                }
            }
            m.SetMailboxInclusionRules(res)
        }
        return nil
    }
    res["mailboxProtectionUnits"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateMailboxProtectionUnitFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]MailboxProtectionUnitable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(MailboxProtectionUnitable)
                }
            }
            m.SetMailboxProtectionUnits(res)
        }
        return nil
    }
    return res
}
// GetMailboxInclusionRules gets the mailboxInclusionRules property value. The rules associated with the Exchange protection policy.
// returns a []MailboxProtectionRuleable when successful
func (m *ExchangeProtectionPolicy) GetMailboxInclusionRules()([]MailboxProtectionRuleable) {
    val, err := m.GetBackingStore().Get("mailboxInclusionRules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MailboxProtectionRuleable)
    }
    return nil
}
// GetMailboxProtectionUnits gets the mailboxProtectionUnits property value. The protection units (mailboxes) that are  protected under the Exchange protection policy.
// returns a []MailboxProtectionUnitable when successful
func (m *ExchangeProtectionPolicy) GetMailboxProtectionUnits()([]MailboxProtectionUnitable) {
    val, err := m.GetBackingStore().Get("mailboxProtectionUnits")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]MailboxProtectionUnitable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ExchangeProtectionPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ProtectionPolicyBase.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetMailboxInclusionRules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMailboxInclusionRules()))
        for i, v := range m.GetMailboxInclusionRules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("mailboxInclusionRules", cast)
        if err != nil {
            return err
        }
    }
    if m.GetMailboxProtectionUnits() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMailboxProtectionUnits()))
        for i, v := range m.GetMailboxProtectionUnits() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("mailboxProtectionUnits", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMailboxInclusionRules sets the mailboxInclusionRules property value. The rules associated with the Exchange protection policy.
func (m *ExchangeProtectionPolicy) SetMailboxInclusionRules(value []MailboxProtectionRuleable)() {
    err := m.GetBackingStore().Set("mailboxInclusionRules", value)
    if err != nil {
        panic(err)
    }
}
// SetMailboxProtectionUnits sets the mailboxProtectionUnits property value. The protection units (mailboxes) that are  protected under the Exchange protection policy.
func (m *ExchangeProtectionPolicy) SetMailboxProtectionUnits(value []MailboxProtectionUnitable)() {
    err := m.GetBackingStore().Set("mailboxProtectionUnits", value)
    if err != nil {
        panic(err)
    }
}
type ExchangeProtectionPolicyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    ProtectionPolicyBaseable
    GetMailboxInclusionRules()([]MailboxProtectionRuleable)
    GetMailboxProtectionUnits()([]MailboxProtectionUnitable)
    SetMailboxInclusionRules(value []MailboxProtectionRuleable)()
    SetMailboxProtectionUnits(value []MailboxProtectionUnitable)()
}

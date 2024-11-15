package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnifiedRoleManagementPolicyEnablementRule struct {
    UnifiedRoleManagementPolicyRule
}
// NewUnifiedRoleManagementPolicyEnablementRule instantiates a new UnifiedRoleManagementPolicyEnablementRule and sets the default values.
func NewUnifiedRoleManagementPolicyEnablementRule()(*UnifiedRoleManagementPolicyEnablementRule) {
    m := &UnifiedRoleManagementPolicyEnablementRule{
        UnifiedRoleManagementPolicyRule: *NewUnifiedRoleManagementPolicyRule(),
    }
    odataTypeValue := "#microsoft.graph.unifiedRoleManagementPolicyEnablementRule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateUnifiedRoleManagementPolicyEnablementRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedRoleManagementPolicyEnablementRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleManagementPolicyEnablementRule(), nil
}
// GetEnabledRules gets the enabledRules property value. The collection of rules that are enabled for this policy rule. For example, MultiFactorAuthentication, Ticketing, and Justification.
// returns a []string when successful
func (m *UnifiedRoleManagementPolicyEnablementRule) GetEnabledRules()([]string) {
    val, err := m.GetBackingStore().Get("enabledRules")
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
func (m *UnifiedRoleManagementPolicyEnablementRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.UnifiedRoleManagementPolicyRule.GetFieldDeserializers()
    res["enabledRules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
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
            m.SetEnabledRules(res)
        }
        return nil
    }
    return res
}
// Serialize serializes information the current object
func (m *UnifiedRoleManagementPolicyEnablementRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.UnifiedRoleManagementPolicyRule.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetEnabledRules() != nil {
        err = writer.WriteCollectionOfStringValues("enabledRules", m.GetEnabledRules())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEnabledRules sets the enabledRules property value. The collection of rules that are enabled for this policy rule. For example, MultiFactorAuthentication, Ticketing, and Justification.
func (m *UnifiedRoleManagementPolicyEnablementRule) SetEnabledRules(value []string)() {
    err := m.GetBackingStore().Set("enabledRules", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedRoleManagementPolicyEnablementRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    UnifiedRoleManagementPolicyRuleable
    GetEnabledRules()([]string)
    SetEnabledRules(value []string)()
}

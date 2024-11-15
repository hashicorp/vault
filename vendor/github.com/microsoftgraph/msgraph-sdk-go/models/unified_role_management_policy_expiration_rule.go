package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnifiedRoleManagementPolicyExpirationRule struct {
    UnifiedRoleManagementPolicyRule
}
// NewUnifiedRoleManagementPolicyExpirationRule instantiates a new UnifiedRoleManagementPolicyExpirationRule and sets the default values.
func NewUnifiedRoleManagementPolicyExpirationRule()(*UnifiedRoleManagementPolicyExpirationRule) {
    m := &UnifiedRoleManagementPolicyExpirationRule{
        UnifiedRoleManagementPolicyRule: *NewUnifiedRoleManagementPolicyRule(),
    }
    odataTypeValue := "#microsoft.graph.unifiedRoleManagementPolicyExpirationRule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateUnifiedRoleManagementPolicyExpirationRuleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedRoleManagementPolicyExpirationRuleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleManagementPolicyExpirationRule(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UnifiedRoleManagementPolicyExpirationRule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.UnifiedRoleManagementPolicyRule.GetFieldDeserializers()
    res["isExpirationRequired"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsExpirationRequired(val)
        }
        return nil
    }
    res["maximumDuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMaximumDuration(val)
        }
        return nil
    }
    return res
}
// GetIsExpirationRequired gets the isExpirationRequired property value. Indicates whether expiration is required or if it's a permanently active assignment or eligibility.
// returns a *bool when successful
func (m *UnifiedRoleManagementPolicyExpirationRule) GetIsExpirationRequired()(*bool) {
    val, err := m.GetBackingStore().Get("isExpirationRequired")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMaximumDuration gets the maximumDuration property value. The maximum duration allowed for eligibility or assignment that isn't permanent. Required when isExpirationRequired is true.
// returns a *ISODuration when successful
func (m *UnifiedRoleManagementPolicyExpirationRule) GetMaximumDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("maximumDuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UnifiedRoleManagementPolicyExpirationRule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.UnifiedRoleManagementPolicyRule.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("isExpirationRequired", m.GetIsExpirationRequired())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteISODurationValue("maximumDuration", m.GetMaximumDuration())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetIsExpirationRequired sets the isExpirationRequired property value. Indicates whether expiration is required or if it's a permanently active assignment or eligibility.
func (m *UnifiedRoleManagementPolicyExpirationRule) SetIsExpirationRequired(value *bool)() {
    err := m.GetBackingStore().Set("isExpirationRequired", value)
    if err != nil {
        panic(err)
    }
}
// SetMaximumDuration sets the maximumDuration property value. The maximum duration allowed for eligibility or assignment that isn't permanent. Required when isExpirationRequired is true.
func (m *UnifiedRoleManagementPolicyExpirationRule) SetMaximumDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("maximumDuration", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedRoleManagementPolicyExpirationRuleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    UnifiedRoleManagementPolicyRuleable
    GetIsExpirationRequired()(*bool)
    GetMaximumDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    SetIsExpirationRequired(value *bool)()
    SetMaximumDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
}

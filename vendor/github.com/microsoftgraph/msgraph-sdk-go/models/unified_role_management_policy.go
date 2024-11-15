package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnifiedRoleManagementPolicy struct {
    Entity
}
// NewUnifiedRoleManagementPolicy instantiates a new UnifiedRoleManagementPolicy and sets the default values.
func NewUnifiedRoleManagementPolicy()(*UnifiedRoleManagementPolicy) {
    m := &UnifiedRoleManagementPolicy{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUnifiedRoleManagementPolicyFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedRoleManagementPolicyFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleManagementPolicy(), nil
}
// GetDescription gets the description property value. Description for the policy.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicy) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Display name for the policy.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicy) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEffectiveRules gets the effectiveRules property value. The list of effective rules like approval rules and expiration rules evaluated based on inherited referenced rules. For example, if there is a tenant-wide policy to enforce enabling an approval rule, the effective rule will be to enable approval even if the policy has a rule to disable approval. Supports $expand.
// returns a []UnifiedRoleManagementPolicyRuleable when successful
func (m *UnifiedRoleManagementPolicy) GetEffectiveRules()([]UnifiedRoleManagementPolicyRuleable) {
    val, err := m.GetBackingStore().Get("effectiveRules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleManagementPolicyRuleable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UnifiedRoleManagementPolicy) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["description"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDescription(val)
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
    res["effectiveRules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRoleManagementPolicyRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleManagementPolicyRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRoleManagementPolicyRuleable)
                }
            }
            m.SetEffectiveRules(res)
        }
        return nil
    }
    res["isOrganizationDefault"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsOrganizationDefault(val)
        }
        return nil
    }
    res["lastModifiedBy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentityFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedBy(val.(Identityable))
        }
        return nil
    }
    res["lastModifiedDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLastModifiedDateTime(val)
        }
        return nil
    }
    res["rules"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateUnifiedRoleManagementPolicyRuleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]UnifiedRoleManagementPolicyRuleable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(UnifiedRoleManagementPolicyRuleable)
                }
            }
            m.SetRules(res)
        }
        return nil
    }
    res["scopeId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScopeId(val)
        }
        return nil
    }
    res["scopeType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScopeType(val)
        }
        return nil
    }
    return res
}
// GetIsOrganizationDefault gets the isOrganizationDefault property value. This can only be set to true for a single tenant-wide policy which will apply to all scopes and roles. Set the scopeId to / and scopeType to Directory. Supports $filter (eq, ne).
// returns a *bool when successful
func (m *UnifiedRoleManagementPolicy) GetIsOrganizationDefault()(*bool) {
    val, err := m.GetBackingStore().Get("isOrganizationDefault")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLastModifiedBy gets the lastModifiedBy property value. The identity who last modified the role setting.
// returns a Identityable when successful
func (m *UnifiedRoleManagementPolicy) GetLastModifiedBy()(Identityable) {
    val, err := m.GetBackingStore().Get("lastModifiedBy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Identityable)
    }
    return nil
}
// GetLastModifiedDateTime gets the lastModifiedDateTime property value. The time when the role setting was last modified.
// returns a *Time when successful
func (m *UnifiedRoleManagementPolicy) GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("lastModifiedDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetRules gets the rules property value. The collection of rules like approval rules and expiration rules. Supports $expand.
// returns a []UnifiedRoleManagementPolicyRuleable when successful
func (m *UnifiedRoleManagementPolicy) GetRules()([]UnifiedRoleManagementPolicyRuleable) {
    val, err := m.GetBackingStore().Get("rules")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]UnifiedRoleManagementPolicyRuleable)
    }
    return nil
}
// GetScopeId gets the scopeId property value. The identifier of the scope where the policy is created. Can be / for the tenant or a group ID. Required.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicy) GetScopeId()(*string) {
    val, err := m.GetBackingStore().Get("scopeId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScopeType gets the scopeType property value. The type of the scope where the policy is created. One of Directory, DirectoryRole, Group. Required.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicy) GetScopeType()(*string) {
    val, err := m.GetBackingStore().Get("scopeType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UnifiedRoleManagementPolicy) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("description", m.GetDescription())
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
    if m.GetEffectiveRules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetEffectiveRules()))
        for i, v := range m.GetEffectiveRules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("effectiveRules", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isOrganizationDefault", m.GetIsOrganizationDefault())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("lastModifiedBy", m.GetLastModifiedBy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("lastModifiedDateTime", m.GetLastModifiedDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetRules() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetRules()))
        for i, v := range m.GetRules() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("rules", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("scopeId", m.GetScopeId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("scopeType", m.GetScopeType())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. Description for the policy.
func (m *UnifiedRoleManagementPolicy) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Display name for the policy.
func (m *UnifiedRoleManagementPolicy) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetEffectiveRules sets the effectiveRules property value. The list of effective rules like approval rules and expiration rules evaluated based on inherited referenced rules. For example, if there is a tenant-wide policy to enforce enabling an approval rule, the effective rule will be to enable approval even if the policy has a rule to disable approval. Supports $expand.
func (m *UnifiedRoleManagementPolicy) SetEffectiveRules(value []UnifiedRoleManagementPolicyRuleable)() {
    err := m.GetBackingStore().Set("effectiveRules", value)
    if err != nil {
        panic(err)
    }
}
// SetIsOrganizationDefault sets the isOrganizationDefault property value. This can only be set to true for a single tenant-wide policy which will apply to all scopes and roles. Set the scopeId to / and scopeType to Directory. Supports $filter (eq, ne).
func (m *UnifiedRoleManagementPolicy) SetIsOrganizationDefault(value *bool)() {
    err := m.GetBackingStore().Set("isOrganizationDefault", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedBy sets the lastModifiedBy property value. The identity who last modified the role setting.
func (m *UnifiedRoleManagementPolicy) SetLastModifiedBy(value Identityable)() {
    err := m.GetBackingStore().Set("lastModifiedBy", value)
    if err != nil {
        panic(err)
    }
}
// SetLastModifiedDateTime sets the lastModifiedDateTime property value. The time when the role setting was last modified.
func (m *UnifiedRoleManagementPolicy) SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("lastModifiedDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetRules sets the rules property value. The collection of rules like approval rules and expiration rules. Supports $expand.
func (m *UnifiedRoleManagementPolicy) SetRules(value []UnifiedRoleManagementPolicyRuleable)() {
    err := m.GetBackingStore().Set("rules", value)
    if err != nil {
        panic(err)
    }
}
// SetScopeId sets the scopeId property value. The identifier of the scope where the policy is created. Can be / for the tenant or a group ID. Required.
func (m *UnifiedRoleManagementPolicy) SetScopeId(value *string)() {
    err := m.GetBackingStore().Set("scopeId", value)
    if err != nil {
        panic(err)
    }
}
// SetScopeType sets the scopeType property value. The type of the scope where the policy is created. One of Directory, DirectoryRole, Group. Required.
func (m *UnifiedRoleManagementPolicy) SetScopeType(value *string)() {
    err := m.GetBackingStore().Set("scopeType", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedRoleManagementPolicyable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetEffectiveRules()([]UnifiedRoleManagementPolicyRuleable)
    GetIsOrganizationDefault()(*bool)
    GetLastModifiedBy()(Identityable)
    GetLastModifiedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetRules()([]UnifiedRoleManagementPolicyRuleable)
    GetScopeId()(*string)
    GetScopeType()(*string)
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetEffectiveRules(value []UnifiedRoleManagementPolicyRuleable)()
    SetIsOrganizationDefault(value *bool)()
    SetLastModifiedBy(value Identityable)()
    SetLastModifiedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetRules(value []UnifiedRoleManagementPolicyRuleable)()
    SetScopeId(value *string)()
    SetScopeType(value *string)()
}

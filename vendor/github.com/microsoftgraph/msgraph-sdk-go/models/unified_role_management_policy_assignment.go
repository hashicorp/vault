package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnifiedRoleManagementPolicyAssignment struct {
    Entity
}
// NewUnifiedRoleManagementPolicyAssignment instantiates a new UnifiedRoleManagementPolicyAssignment and sets the default values.
func NewUnifiedRoleManagementPolicyAssignment()(*UnifiedRoleManagementPolicyAssignment) {
    m := &UnifiedRoleManagementPolicyAssignment{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUnifiedRoleManagementPolicyAssignmentFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedRoleManagementPolicyAssignmentFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleManagementPolicyAssignment(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UnifiedRoleManagementPolicyAssignment) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["policy"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUnifiedRoleManagementPolicyFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicy(val.(UnifiedRoleManagementPolicyable))
        }
        return nil
    }
    res["policyId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPolicyId(val)
        }
        return nil
    }
    res["roleDefinitionId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoleDefinitionId(val)
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
// GetPolicy gets the policy property value. The policy that's associated with a policy assignment. Supports $expand and a nested $expand of the rules and effectiveRules relationships for the policy.
// returns a UnifiedRoleManagementPolicyable when successful
func (m *UnifiedRoleManagementPolicyAssignment) GetPolicy()(UnifiedRoleManagementPolicyable) {
    val, err := m.GetBackingStore().Get("policy")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UnifiedRoleManagementPolicyable)
    }
    return nil
}
// GetPolicyId gets the policyId property value. The id of the policy. Inherited from entity.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicyAssignment) GetPolicyId()(*string) {
    val, err := m.GetBackingStore().Get("policyId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRoleDefinitionId gets the roleDefinitionId property value. For Microsoft Entra roles policy, it's the identifier of the role definition object where the policy applies. For PIM for groups membership and ownership, it's either member or owner. Supports $filter (eq).
// returns a *string when successful
func (m *UnifiedRoleManagementPolicyAssignment) GetRoleDefinitionId()(*string) {
    val, err := m.GetBackingStore().Get("roleDefinitionId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScopeId gets the scopeId property value. The identifier of the scope where the policy is assigned.  Can be / for the tenant or a group ID. Required.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicyAssignment) GetScopeId()(*string) {
    val, err := m.GetBackingStore().Get("scopeId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScopeType gets the scopeType property value. The type of the scope where the policy is assigned. One of Directory, DirectoryRole, Group. Required.
// returns a *string when successful
func (m *UnifiedRoleManagementPolicyAssignment) GetScopeType()(*string) {
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
func (m *UnifiedRoleManagementPolicyAssignment) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("policy", m.GetPolicy())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("policyId", m.GetPolicyId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("roleDefinitionId", m.GetRoleDefinitionId())
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
// SetPolicy sets the policy property value. The policy that's associated with a policy assignment. Supports $expand and a nested $expand of the rules and effectiveRules relationships for the policy.
func (m *UnifiedRoleManagementPolicyAssignment) SetPolicy(value UnifiedRoleManagementPolicyable)() {
    err := m.GetBackingStore().Set("policy", value)
    if err != nil {
        panic(err)
    }
}
// SetPolicyId sets the policyId property value. The id of the policy. Inherited from entity.
func (m *UnifiedRoleManagementPolicyAssignment) SetPolicyId(value *string)() {
    err := m.GetBackingStore().Set("policyId", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleDefinitionId sets the roleDefinitionId property value. For Microsoft Entra roles policy, it's the identifier of the role definition object where the policy applies. For PIM for groups membership and ownership, it's either member or owner. Supports $filter (eq).
func (m *UnifiedRoleManagementPolicyAssignment) SetRoleDefinitionId(value *string)() {
    err := m.GetBackingStore().Set("roleDefinitionId", value)
    if err != nil {
        panic(err)
    }
}
// SetScopeId sets the scopeId property value. The identifier of the scope where the policy is assigned.  Can be / for the tenant or a group ID. Required.
func (m *UnifiedRoleManagementPolicyAssignment) SetScopeId(value *string)() {
    err := m.GetBackingStore().Set("scopeId", value)
    if err != nil {
        panic(err)
    }
}
// SetScopeType sets the scopeType property value. The type of the scope where the policy is assigned. One of Directory, DirectoryRole, Group. Required.
func (m *UnifiedRoleManagementPolicyAssignment) SetScopeType(value *string)() {
    err := m.GetBackingStore().Set("scopeType", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedRoleManagementPolicyAssignmentable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetPolicy()(UnifiedRoleManagementPolicyable)
    GetPolicyId()(*string)
    GetRoleDefinitionId()(*string)
    GetScopeId()(*string)
    GetScopeType()(*string)
    SetPolicy(value UnifiedRoleManagementPolicyable)()
    SetPolicyId(value *string)()
    SetRoleDefinitionId(value *string)()
    SetScopeId(value *string)()
    SetScopeType(value *string)()
}

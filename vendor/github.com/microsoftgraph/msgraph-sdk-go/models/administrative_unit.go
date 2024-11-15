package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type AdministrativeUnit struct {
    DirectoryObject
}
// NewAdministrativeUnit instantiates a new AdministrativeUnit and sets the default values.
func NewAdministrativeUnit()(*AdministrativeUnit) {
    m := &AdministrativeUnit{
        DirectoryObject: *NewDirectoryObject(),
    }
    odataTypeValue := "#microsoft.graph.administrativeUnit"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAdministrativeUnitFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAdministrativeUnitFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAdministrativeUnit(), nil
}
// GetDescription gets the description property value. An optional description for the administrative unit. Supports $filter (eq, ne, in, startsWith), $search.
// returns a *string when successful
func (m *AdministrativeUnit) GetDescription()(*string) {
    val, err := m.GetBackingStore().Get("description")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisplayName gets the displayName property value. Display name for the administrative unit. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.
// returns a *string when successful
func (m *AdministrativeUnit) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetExtensions gets the extensions property value. The collection of open extensions defined for this administrative unit. Nullable.
// returns a []Extensionable when successful
func (m *AdministrativeUnit) GetExtensions()([]Extensionable) {
    val, err := m.GetBackingStore().Get("extensions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Extensionable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *AdministrativeUnit) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.DirectoryObject.GetFieldDeserializers()
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
    res["extensions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateExtensionFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Extensionable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Extensionable)
                }
            }
            m.SetExtensions(res)
        }
        return nil
    }
    res["isMemberManagementRestricted"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsMemberManagementRestricted(val)
        }
        return nil
    }
    res["members"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]DirectoryObjectable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(DirectoryObjectable)
                }
            }
            m.SetMembers(res)
        }
        return nil
    }
    res["membershipRule"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMembershipRule(val)
        }
        return nil
    }
    res["membershipRuleProcessingState"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMembershipRuleProcessingState(val)
        }
        return nil
    }
    res["membershipType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMembershipType(val)
        }
        return nil
    }
    res["scopedRoleMembers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateScopedRoleMembershipFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]ScopedRoleMembershipable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(ScopedRoleMembershipable)
                }
            }
            m.SetScopedRoleMembers(res)
        }
        return nil
    }
    res["visibility"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetVisibility(val)
        }
        return nil
    }
    return res
}
// GetIsMemberManagementRestricted gets the isMemberManagementRestricted property value. The isMemberManagementRestricted property
// returns a *bool when successful
func (m *AdministrativeUnit) GetIsMemberManagementRestricted()(*bool) {
    val, err := m.GetBackingStore().Get("isMemberManagementRestricted")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMembers gets the members property value. Users and groups that are members of this administrative unit. Supports $expand.
// returns a []DirectoryObjectable when successful
func (m *AdministrativeUnit) GetMembers()([]DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("members")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]DirectoryObjectable)
    }
    return nil
}
// GetMembershipRule gets the membershipRule property value. The membershipRule property
// returns a *string when successful
func (m *AdministrativeUnit) GetMembershipRule()(*string) {
    val, err := m.GetBackingStore().Get("membershipRule")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMembershipRuleProcessingState gets the membershipRuleProcessingState property value. The membershipRuleProcessingState property
// returns a *string when successful
func (m *AdministrativeUnit) GetMembershipRuleProcessingState()(*string) {
    val, err := m.GetBackingStore().Get("membershipRuleProcessingState")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMembershipType gets the membershipType property value. The membershipType property
// returns a *string when successful
func (m *AdministrativeUnit) GetMembershipType()(*string) {
    val, err := m.GetBackingStore().Get("membershipType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScopedRoleMembers gets the scopedRoleMembers property value. Scoped-role members of this administrative unit.
// returns a []ScopedRoleMembershipable when successful
func (m *AdministrativeUnit) GetScopedRoleMembers()([]ScopedRoleMembershipable) {
    val, err := m.GetBackingStore().Get("scopedRoleMembers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]ScopedRoleMembershipable)
    }
    return nil
}
// GetVisibility gets the visibility property value. Controls whether the administrative unit and its members are hidden or public. Can be set to HiddenMembership. If not set (value is null), the default behavior is public. When set to HiddenMembership, only members of the administrative unit can list other members of the administrative unit.
// returns a *string when successful
func (m *AdministrativeUnit) GetVisibility()(*string) {
    val, err := m.GetBackingStore().Get("visibility")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *AdministrativeUnit) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.DirectoryObject.Serialize(writer)
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
    if m.GetExtensions() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetExtensions()))
        for i, v := range m.GetExtensions() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("extensions", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isMemberManagementRestricted", m.GetIsMemberManagementRestricted())
        if err != nil {
            return err
        }
    }
    if m.GetMembers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetMembers()))
        for i, v := range m.GetMembers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("members", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("membershipRule", m.GetMembershipRule())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("membershipRuleProcessingState", m.GetMembershipRuleProcessingState())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("membershipType", m.GetMembershipType())
        if err != nil {
            return err
        }
    }
    if m.GetScopedRoleMembers() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetScopedRoleMembers()))
        for i, v := range m.GetScopedRoleMembers() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("scopedRoleMembers", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("visibility", m.GetVisibility())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDescription sets the description property value. An optional description for the administrative unit. Supports $filter (eq, ne, in, startsWith), $search.
func (m *AdministrativeUnit) SetDescription(value *string)() {
    err := m.GetBackingStore().Set("description", value)
    if err != nil {
        panic(err)
    }
}
// SetDisplayName sets the displayName property value. Display name for the administrative unit. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.
func (m *AdministrativeUnit) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetExtensions sets the extensions property value. The collection of open extensions defined for this administrative unit. Nullable.
func (m *AdministrativeUnit) SetExtensions(value []Extensionable)() {
    err := m.GetBackingStore().Set("extensions", value)
    if err != nil {
        panic(err)
    }
}
// SetIsMemberManagementRestricted sets the isMemberManagementRestricted property value. The isMemberManagementRestricted property
func (m *AdministrativeUnit) SetIsMemberManagementRestricted(value *bool)() {
    err := m.GetBackingStore().Set("isMemberManagementRestricted", value)
    if err != nil {
        panic(err)
    }
}
// SetMembers sets the members property value. Users and groups that are members of this administrative unit. Supports $expand.
func (m *AdministrativeUnit) SetMembers(value []DirectoryObjectable)() {
    err := m.GetBackingStore().Set("members", value)
    if err != nil {
        panic(err)
    }
}
// SetMembershipRule sets the membershipRule property value. The membershipRule property
func (m *AdministrativeUnit) SetMembershipRule(value *string)() {
    err := m.GetBackingStore().Set("membershipRule", value)
    if err != nil {
        panic(err)
    }
}
// SetMembershipRuleProcessingState sets the membershipRuleProcessingState property value. The membershipRuleProcessingState property
func (m *AdministrativeUnit) SetMembershipRuleProcessingState(value *string)() {
    err := m.GetBackingStore().Set("membershipRuleProcessingState", value)
    if err != nil {
        panic(err)
    }
}
// SetMembershipType sets the membershipType property value. The membershipType property
func (m *AdministrativeUnit) SetMembershipType(value *string)() {
    err := m.GetBackingStore().Set("membershipType", value)
    if err != nil {
        panic(err)
    }
}
// SetScopedRoleMembers sets the scopedRoleMembers property value. Scoped-role members of this administrative unit.
func (m *AdministrativeUnit) SetScopedRoleMembers(value []ScopedRoleMembershipable)() {
    err := m.GetBackingStore().Set("scopedRoleMembers", value)
    if err != nil {
        panic(err)
    }
}
// SetVisibility sets the visibility property value. Controls whether the administrative unit and its members are hidden or public. Can be set to HiddenMembership. If not set (value is null), the default behavior is public. When set to HiddenMembership, only members of the administrative unit can list other members of the administrative unit.
func (m *AdministrativeUnit) SetVisibility(value *string)() {
    err := m.GetBackingStore().Set("visibility", value)
    if err != nil {
        panic(err)
    }
}
type AdministrativeUnitable interface {
    DirectoryObjectable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetExtensions()([]Extensionable)
    GetIsMemberManagementRestricted()(*bool)
    GetMembers()([]DirectoryObjectable)
    GetMembershipRule()(*string)
    GetMembershipRuleProcessingState()(*string)
    GetMembershipType()(*string)
    GetScopedRoleMembers()([]ScopedRoleMembershipable)
    GetVisibility()(*string)
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetExtensions(value []Extensionable)()
    SetIsMemberManagementRestricted(value *bool)()
    SetMembers(value []DirectoryObjectable)()
    SetMembershipRule(value *string)()
    SetMembershipRuleProcessingState(value *string)()
    SetMembershipType(value *string)()
    SetScopedRoleMembers(value []ScopedRoleMembershipable)()
    SetVisibility(value *string)()
}

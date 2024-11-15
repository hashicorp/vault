package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PrivilegedAccessGroupAssignmentSchedule struct {
    PrivilegedAccessSchedule
}
// NewPrivilegedAccessGroupAssignmentSchedule instantiates a new PrivilegedAccessGroupAssignmentSchedule and sets the default values.
func NewPrivilegedAccessGroupAssignmentSchedule()(*PrivilegedAccessGroupAssignmentSchedule) {
    m := &PrivilegedAccessGroupAssignmentSchedule{
        PrivilegedAccessSchedule: *NewPrivilegedAccessSchedule(),
    }
    odataTypeValue := "#microsoft.graph.privilegedAccessGroupAssignmentSchedule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreatePrivilegedAccessGroupAssignmentScheduleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrivilegedAccessGroupAssignmentScheduleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrivilegedAccessGroupAssignmentSchedule(), nil
}
// GetAccessId gets the accessId property value. The identifier of the membership or ownership assignment to the group that is governed through PIM. Required. The possible values are: owner, member, unknownFutureValue. Supports $filter (eq).
// returns a *PrivilegedAccessGroupRelationships when successful
func (m *PrivilegedAccessGroupAssignmentSchedule) GetAccessId()(*PrivilegedAccessGroupRelationships) {
    val, err := m.GetBackingStore().Get("accessId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrivilegedAccessGroupRelationships)
    }
    return nil
}
// GetActivatedUsing gets the activatedUsing property value. When the request activates an ownership or membership assignment in PIM for groups, this object represents the eligibility relationship. Otherwise, it's null. Supports $expand.
// returns a PrivilegedAccessGroupEligibilityScheduleable when successful
func (m *PrivilegedAccessGroupAssignmentSchedule) GetActivatedUsing()(PrivilegedAccessGroupEligibilityScheduleable) {
    val, err := m.GetBackingStore().Get("activatedUsing")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PrivilegedAccessGroupEligibilityScheduleable)
    }
    return nil
}
// GetAssignmentType gets the assignmentType property value. Indicates whether the membership or ownership assignment for the principal is granted through activation or direct assignment. Required. The possible values are: assigned, activated, unknownFutureValue. Supports $filter (eq).
// returns a *PrivilegedAccessGroupAssignmentType when successful
func (m *PrivilegedAccessGroupAssignmentSchedule) GetAssignmentType()(*PrivilegedAccessGroupAssignmentType) {
    val, err := m.GetBackingStore().Get("assignmentType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrivilegedAccessGroupAssignmentType)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrivilegedAccessGroupAssignmentSchedule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.PrivilegedAccessSchedule.GetFieldDeserializers()
    res["accessId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrivilegedAccessGroupRelationships)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAccessId(val.(*PrivilegedAccessGroupRelationships))
        }
        return nil
    }
    res["activatedUsing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePrivilegedAccessGroupEligibilityScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivatedUsing(val.(PrivilegedAccessGroupEligibilityScheduleable))
        }
        return nil
    }
    res["assignmentType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrivilegedAccessGroupAssignmentType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignmentType(val.(*PrivilegedAccessGroupAssignmentType))
        }
        return nil
    }
    res["group"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateGroupFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroup(val.(Groupable))
        }
        return nil
    }
    res["groupId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGroupId(val)
        }
        return nil
    }
    res["memberType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParsePrivilegedAccessGroupMemberType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMemberType(val.(*PrivilegedAccessGroupMemberType))
        }
        return nil
    }
    res["principal"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDirectoryObjectFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrincipal(val.(DirectoryObjectable))
        }
        return nil
    }
    res["principalId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPrincipalId(val)
        }
        return nil
    }
    return res
}
// GetGroup gets the group property value. References the group that is the scope of the membership or ownership assignment through PIM for groups. Supports $expand and $select nested in $expand for select properties like id, displayName, and mail.
// returns a Groupable when successful
func (m *PrivilegedAccessGroupAssignmentSchedule) GetGroup()(Groupable) {
    val, err := m.GetBackingStore().Get("group")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Groupable)
    }
    return nil
}
// GetGroupId gets the groupId property value. The identifier of the group representing the scope of the membership or ownership assignment through PIM for groups. Required. Supports $filter (eq).
// returns a *string when successful
func (m *PrivilegedAccessGroupAssignmentSchedule) GetGroupId()(*string) {
    val, err := m.GetBackingStore().Get("groupId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMemberType gets the memberType property value. Indicates whether the assignment is derived from a direct group assignment or through a transitive assignment. The possible values are: direct, group, unknownFutureValue. Supports $filter (eq).
// returns a *PrivilegedAccessGroupMemberType when successful
func (m *PrivilegedAccessGroupAssignmentSchedule) GetMemberType()(*PrivilegedAccessGroupMemberType) {
    val, err := m.GetBackingStore().Get("memberType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrivilegedAccessGroupMemberType)
    }
    return nil
}
// GetPrincipal gets the principal property value. References the principal that's in the scope of this membership or ownership assignment request to the group that's governed through PIM. Supports $expand and $select nested in $expand for id only.
// returns a DirectoryObjectable when successful
func (m *PrivilegedAccessGroupAssignmentSchedule) GetPrincipal()(DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("principal")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DirectoryObjectable)
    }
    return nil
}
// GetPrincipalId gets the principalId property value. The identifier of the principal whose membership or ownership assignment is granted through PIM for groups. Required. Supports $filter (eq).
// returns a *string when successful
func (m *PrivilegedAccessGroupAssignmentSchedule) GetPrincipalId()(*string) {
    val, err := m.GetBackingStore().Get("principalId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *PrivilegedAccessGroupAssignmentSchedule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.PrivilegedAccessSchedule.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetAccessId() != nil {
        cast := (*m.GetAccessId()).String()
        err = writer.WriteStringValue("accessId", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("activatedUsing", m.GetActivatedUsing())
        if err != nil {
            return err
        }
    }
    if m.GetAssignmentType() != nil {
        cast := (*m.GetAssignmentType()).String()
        err = writer.WriteStringValue("assignmentType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("group", m.GetGroup())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("groupId", m.GetGroupId())
        if err != nil {
            return err
        }
    }
    if m.GetMemberType() != nil {
        cast := (*m.GetMemberType()).String()
        err = writer.WriteStringValue("memberType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("principal", m.GetPrincipal())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("principalId", m.GetPrincipalId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAccessId sets the accessId property value. The identifier of the membership or ownership assignment to the group that is governed through PIM. Required. The possible values are: owner, member, unknownFutureValue. Supports $filter (eq).
func (m *PrivilegedAccessGroupAssignmentSchedule) SetAccessId(value *PrivilegedAccessGroupRelationships)() {
    err := m.GetBackingStore().Set("accessId", value)
    if err != nil {
        panic(err)
    }
}
// SetActivatedUsing sets the activatedUsing property value. When the request activates an ownership or membership assignment in PIM for groups, this object represents the eligibility relationship. Otherwise, it's null. Supports $expand.
func (m *PrivilegedAccessGroupAssignmentSchedule) SetActivatedUsing(value PrivilegedAccessGroupEligibilityScheduleable)() {
    err := m.GetBackingStore().Set("activatedUsing", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignmentType sets the assignmentType property value. Indicates whether the membership or ownership assignment for the principal is granted through activation or direct assignment. Required. The possible values are: assigned, activated, unknownFutureValue. Supports $filter (eq).
func (m *PrivilegedAccessGroupAssignmentSchedule) SetAssignmentType(value *PrivilegedAccessGroupAssignmentType)() {
    err := m.GetBackingStore().Set("assignmentType", value)
    if err != nil {
        panic(err)
    }
}
// SetGroup sets the group property value. References the group that is the scope of the membership or ownership assignment through PIM for groups. Supports $expand and $select nested in $expand for select properties like id, displayName, and mail.
func (m *PrivilegedAccessGroupAssignmentSchedule) SetGroup(value Groupable)() {
    err := m.GetBackingStore().Set("group", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupId sets the groupId property value. The identifier of the group representing the scope of the membership or ownership assignment through PIM for groups. Required. Supports $filter (eq).
func (m *PrivilegedAccessGroupAssignmentSchedule) SetGroupId(value *string)() {
    err := m.GetBackingStore().Set("groupId", value)
    if err != nil {
        panic(err)
    }
}
// SetMemberType sets the memberType property value. Indicates whether the assignment is derived from a direct group assignment or through a transitive assignment. The possible values are: direct, group, unknownFutureValue. Supports $filter (eq).
func (m *PrivilegedAccessGroupAssignmentSchedule) SetMemberType(value *PrivilegedAccessGroupMemberType)() {
    err := m.GetBackingStore().Set("memberType", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipal sets the principal property value. References the principal that's in the scope of this membership or ownership assignment request to the group that's governed through PIM. Supports $expand and $select nested in $expand for id only.
func (m *PrivilegedAccessGroupAssignmentSchedule) SetPrincipal(value DirectoryObjectable)() {
    err := m.GetBackingStore().Set("principal", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipalId sets the principalId property value. The identifier of the principal whose membership or ownership assignment is granted through PIM for groups. Required. Supports $filter (eq).
func (m *PrivilegedAccessGroupAssignmentSchedule) SetPrincipalId(value *string)() {
    err := m.GetBackingStore().Set("principalId", value)
    if err != nil {
        panic(err)
    }
}
type PrivilegedAccessGroupAssignmentScheduleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    PrivilegedAccessScheduleable
    GetAccessId()(*PrivilegedAccessGroupRelationships)
    GetActivatedUsing()(PrivilegedAccessGroupEligibilityScheduleable)
    GetAssignmentType()(*PrivilegedAccessGroupAssignmentType)
    GetGroup()(Groupable)
    GetGroupId()(*string)
    GetMemberType()(*PrivilegedAccessGroupMemberType)
    GetPrincipal()(DirectoryObjectable)
    GetPrincipalId()(*string)
    SetAccessId(value *PrivilegedAccessGroupRelationships)()
    SetActivatedUsing(value PrivilegedAccessGroupEligibilityScheduleable)()
    SetAssignmentType(value *PrivilegedAccessGroupAssignmentType)()
    SetGroup(value Groupable)()
    SetGroupId(value *string)()
    SetMemberType(value *PrivilegedAccessGroupMemberType)()
    SetPrincipal(value DirectoryObjectable)()
    SetPrincipalId(value *string)()
}

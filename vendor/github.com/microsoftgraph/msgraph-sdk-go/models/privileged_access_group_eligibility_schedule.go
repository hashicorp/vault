package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type PrivilegedAccessGroupEligibilitySchedule struct {
    PrivilegedAccessSchedule
}
// NewPrivilegedAccessGroupEligibilitySchedule instantiates a new PrivilegedAccessGroupEligibilitySchedule and sets the default values.
func NewPrivilegedAccessGroupEligibilitySchedule()(*PrivilegedAccessGroupEligibilitySchedule) {
    m := &PrivilegedAccessGroupEligibilitySchedule{
        PrivilegedAccessSchedule: *NewPrivilegedAccessSchedule(),
    }
    odataTypeValue := "#microsoft.graph.privilegedAccessGroupEligibilitySchedule"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreatePrivilegedAccessGroupEligibilityScheduleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePrivilegedAccessGroupEligibilityScheduleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPrivilegedAccessGroupEligibilitySchedule(), nil
}
// GetAccessId gets the accessId property value. The identifier of the membership or ownership eligibility to the group that is governed by PIM. Required. The possible values are: owner, member. Supports $filter (eq).
// returns a *PrivilegedAccessGroupRelationships when successful
func (m *PrivilegedAccessGroupEligibilitySchedule) GetAccessId()(*PrivilegedAccessGroupRelationships) {
    val, err := m.GetBackingStore().Get("accessId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrivilegedAccessGroupRelationships)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *PrivilegedAccessGroupEligibilitySchedule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
// GetGroup gets the group property value. References the group that is the scope of the membership or ownership eligibility through PIM for groups. Supports $expand.
// returns a Groupable when successful
func (m *PrivilegedAccessGroupEligibilitySchedule) GetGroup()(Groupable) {
    val, err := m.GetBackingStore().Get("group")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Groupable)
    }
    return nil
}
// GetGroupId gets the groupId property value. The identifier of the group representing the scope of the membership or ownership eligibility through PIM for groups. Required. Supports $filter (eq).
// returns a *string when successful
func (m *PrivilegedAccessGroupEligibilitySchedule) GetGroupId()(*string) {
    val, err := m.GetBackingStore().Get("groupId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetMemberType gets the memberType property value. Indicates whether the assignment is derived from a group assignment. It can further imply whether the caller can manage the schedule. Required. The possible values are: direct, group, unknownFutureValue. Supports $filter (eq).
// returns a *PrivilegedAccessGroupMemberType when successful
func (m *PrivilegedAccessGroupEligibilitySchedule) GetMemberType()(*PrivilegedAccessGroupMemberType) {
    val, err := m.GetBackingStore().Get("memberType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*PrivilegedAccessGroupMemberType)
    }
    return nil
}
// GetPrincipal gets the principal property value. References the principal that's in the scope of this membership or ownership eligibility request to the group that's governed by PIM. Supports $expand.
// returns a DirectoryObjectable when successful
func (m *PrivilegedAccessGroupEligibilitySchedule) GetPrincipal()(DirectoryObjectable) {
    val, err := m.GetBackingStore().Get("principal")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DirectoryObjectable)
    }
    return nil
}
// GetPrincipalId gets the principalId property value. The identifier of the principal whose membership or ownership eligibility is granted through PIM for groups. Required. Supports $filter (eq).
// returns a *string when successful
func (m *PrivilegedAccessGroupEligibilitySchedule) GetPrincipalId()(*string) {
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
func (m *PrivilegedAccessGroupEligibilitySchedule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
// SetAccessId sets the accessId property value. The identifier of the membership or ownership eligibility to the group that is governed by PIM. Required. The possible values are: owner, member. Supports $filter (eq).
func (m *PrivilegedAccessGroupEligibilitySchedule) SetAccessId(value *PrivilegedAccessGroupRelationships)() {
    err := m.GetBackingStore().Set("accessId", value)
    if err != nil {
        panic(err)
    }
}
// SetGroup sets the group property value. References the group that is the scope of the membership or ownership eligibility through PIM for groups. Supports $expand.
func (m *PrivilegedAccessGroupEligibilitySchedule) SetGroup(value Groupable)() {
    err := m.GetBackingStore().Set("group", value)
    if err != nil {
        panic(err)
    }
}
// SetGroupId sets the groupId property value. The identifier of the group representing the scope of the membership or ownership eligibility through PIM for groups. Required. Supports $filter (eq).
func (m *PrivilegedAccessGroupEligibilitySchedule) SetGroupId(value *string)() {
    err := m.GetBackingStore().Set("groupId", value)
    if err != nil {
        panic(err)
    }
}
// SetMemberType sets the memberType property value. Indicates whether the assignment is derived from a group assignment. It can further imply whether the caller can manage the schedule. Required. The possible values are: direct, group, unknownFutureValue. Supports $filter (eq).
func (m *PrivilegedAccessGroupEligibilitySchedule) SetMemberType(value *PrivilegedAccessGroupMemberType)() {
    err := m.GetBackingStore().Set("memberType", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipal sets the principal property value. References the principal that's in the scope of this membership or ownership eligibility request to the group that's governed by PIM. Supports $expand.
func (m *PrivilegedAccessGroupEligibilitySchedule) SetPrincipal(value DirectoryObjectable)() {
    err := m.GetBackingStore().Set("principal", value)
    if err != nil {
        panic(err)
    }
}
// SetPrincipalId sets the principalId property value. The identifier of the principal whose membership or ownership eligibility is granted through PIM for groups. Required. Supports $filter (eq).
func (m *PrivilegedAccessGroupEligibilitySchedule) SetPrincipalId(value *string)() {
    err := m.GetBackingStore().Set("principalId", value)
    if err != nil {
        panic(err)
    }
}
type PrivilegedAccessGroupEligibilityScheduleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    PrivilegedAccessScheduleable
    GetAccessId()(*PrivilegedAccessGroupRelationships)
    GetGroup()(Groupable)
    GetGroupId()(*string)
    GetMemberType()(*PrivilegedAccessGroupMemberType)
    GetPrincipal()(DirectoryObjectable)
    GetPrincipalId()(*string)
    SetAccessId(value *PrivilegedAccessGroupRelationships)()
    SetGroup(value Groupable)()
    SetGroupId(value *string)()
    SetMemberType(value *PrivilegedAccessGroupMemberType)()
    SetPrincipal(value DirectoryObjectable)()
    SetPrincipalId(value *string)()
}

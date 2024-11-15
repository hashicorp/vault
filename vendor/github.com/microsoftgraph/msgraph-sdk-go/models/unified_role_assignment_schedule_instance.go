package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnifiedRoleAssignmentScheduleInstance struct {
    UnifiedRoleScheduleInstanceBase
}
// NewUnifiedRoleAssignmentScheduleInstance instantiates a new UnifiedRoleAssignmentScheduleInstance and sets the default values.
func NewUnifiedRoleAssignmentScheduleInstance()(*UnifiedRoleAssignmentScheduleInstance) {
    m := &UnifiedRoleAssignmentScheduleInstance{
        UnifiedRoleScheduleInstanceBase: *NewUnifiedRoleScheduleInstanceBase(),
    }
    return m
}
// CreateUnifiedRoleAssignmentScheduleInstanceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedRoleAssignmentScheduleInstanceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleAssignmentScheduleInstance(), nil
}
// GetActivatedUsing gets the activatedUsing property value. If the request is from an eligible administrator to activate a role, this parameter shows the related eligible assignment for that activation. Otherwise, it's null. Supports $expand and $select nested in $expand.
// returns a UnifiedRoleEligibilityScheduleInstanceable when successful
func (m *UnifiedRoleAssignmentScheduleInstance) GetActivatedUsing()(UnifiedRoleEligibilityScheduleInstanceable) {
    val, err := m.GetBackingStore().Get("activatedUsing")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(UnifiedRoleEligibilityScheduleInstanceable)
    }
    return nil
}
// GetAssignmentType gets the assignmentType property value. The type of the assignment that can either be Assigned or Activated. Supports $filter (eq, ne).
// returns a *string when successful
func (m *UnifiedRoleAssignmentScheduleInstance) GetAssignmentType()(*string) {
    val, err := m.GetBackingStore().Get("assignmentType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetEndDateTime gets the endDateTime property value. The end date of the schedule instance.
// returns a *Time when successful
func (m *UnifiedRoleAssignmentScheduleInstance) GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("endDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UnifiedRoleAssignmentScheduleInstance) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.UnifiedRoleScheduleInstanceBase.GetFieldDeserializers()
    res["activatedUsing"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateUnifiedRoleEligibilityScheduleInstanceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivatedUsing(val.(UnifiedRoleEligibilityScheduleInstanceable))
        }
        return nil
    }
    res["assignmentType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAssignmentType(val)
        }
        return nil
    }
    res["endDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEndDateTime(val)
        }
        return nil
    }
    res["memberType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMemberType(val)
        }
        return nil
    }
    res["roleAssignmentOriginId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoleAssignmentOriginId(val)
        }
        return nil
    }
    res["roleAssignmentScheduleId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRoleAssignmentScheduleId(val)
        }
        return nil
    }
    res["startDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetTimeValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartDateTime(val)
        }
        return nil
    }
    return res
}
// GetMemberType gets the memberType property value. How the assignment is inherited. It can either be Inherited, Direct, or Group. It can further imply whether the unifiedRoleAssignmentSchedule can be managed by the caller. Supports $filter (eq, ne).
// returns a *string when successful
func (m *UnifiedRoleAssignmentScheduleInstance) GetMemberType()(*string) {
    val, err := m.GetBackingStore().Get("memberType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRoleAssignmentOriginId gets the roleAssignmentOriginId property value. The identifier of the role assignment in Microsoft Entra. Supports $filter (eq, ne).
// returns a *string when successful
func (m *UnifiedRoleAssignmentScheduleInstance) GetRoleAssignmentOriginId()(*string) {
    val, err := m.GetBackingStore().Get("roleAssignmentOriginId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetRoleAssignmentScheduleId gets the roleAssignmentScheduleId property value. The identifier of the unifiedRoleAssignmentSchedule object from which this instance was created. Supports $filter (eq, ne).
// returns a *string when successful
func (m *UnifiedRoleAssignmentScheduleInstance) GetRoleAssignmentScheduleId()(*string) {
    val, err := m.GetBackingStore().Get("roleAssignmentScheduleId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetStartDateTime gets the startDateTime property value. When this instance starts.
// returns a *Time when successful
func (m *UnifiedRoleAssignmentScheduleInstance) GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time) {
    val, err := m.GetBackingStore().Get("startDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UnifiedRoleAssignmentScheduleInstance) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.UnifiedRoleScheduleInstanceBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("activatedUsing", m.GetActivatedUsing())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("assignmentType", m.GetAssignmentType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("endDateTime", m.GetEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("memberType", m.GetMemberType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("roleAssignmentOriginId", m.GetRoleAssignmentOriginId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("roleAssignmentScheduleId", m.GetRoleAssignmentScheduleId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteTimeValue("startDateTime", m.GetStartDateTime())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivatedUsing sets the activatedUsing property value. If the request is from an eligible administrator to activate a role, this parameter shows the related eligible assignment for that activation. Otherwise, it's null. Supports $expand and $select nested in $expand.
func (m *UnifiedRoleAssignmentScheduleInstance) SetActivatedUsing(value UnifiedRoleEligibilityScheduleInstanceable)() {
    err := m.GetBackingStore().Set("activatedUsing", value)
    if err != nil {
        panic(err)
    }
}
// SetAssignmentType sets the assignmentType property value. The type of the assignment that can either be Assigned or Activated. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) SetAssignmentType(value *string)() {
    err := m.GetBackingStore().Set("assignmentType", value)
    if err != nil {
        panic(err)
    }
}
// SetEndDateTime sets the endDateTime property value. The end date of the schedule instance.
func (m *UnifiedRoleAssignmentScheduleInstance) SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("endDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetMemberType sets the memberType property value. How the assignment is inherited. It can either be Inherited, Direct, or Group. It can further imply whether the unifiedRoleAssignmentSchedule can be managed by the caller. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) SetMemberType(value *string)() {
    err := m.GetBackingStore().Set("memberType", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleAssignmentOriginId sets the roleAssignmentOriginId property value. The identifier of the role assignment in Microsoft Entra. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) SetRoleAssignmentOriginId(value *string)() {
    err := m.GetBackingStore().Set("roleAssignmentOriginId", value)
    if err != nil {
        panic(err)
    }
}
// SetRoleAssignmentScheduleId sets the roleAssignmentScheduleId property value. The identifier of the unifiedRoleAssignmentSchedule object from which this instance was created. Supports $filter (eq, ne).
func (m *UnifiedRoleAssignmentScheduleInstance) SetRoleAssignmentScheduleId(value *string)() {
    err := m.GetBackingStore().Set("roleAssignmentScheduleId", value)
    if err != nil {
        panic(err)
    }
}
// SetStartDateTime sets the startDateTime property value. When this instance starts.
func (m *UnifiedRoleAssignmentScheduleInstance) SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)() {
    err := m.GetBackingStore().Set("startDateTime", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedRoleAssignmentScheduleInstanceable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    UnifiedRoleScheduleInstanceBaseable
    GetActivatedUsing()(UnifiedRoleEligibilityScheduleInstanceable)
    GetAssignmentType()(*string)
    GetEndDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetMemberType()(*string)
    GetRoleAssignmentOriginId()(*string)
    GetRoleAssignmentScheduleId()(*string)
    GetStartDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    SetActivatedUsing(value UnifiedRoleEligibilityScheduleInstanceable)()
    SetAssignmentType(value *string)()
    SetEndDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetMemberType(value *string)()
    SetRoleAssignmentOriginId(value *string)()
    SetRoleAssignmentScheduleId(value *string)()
    SetStartDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
}

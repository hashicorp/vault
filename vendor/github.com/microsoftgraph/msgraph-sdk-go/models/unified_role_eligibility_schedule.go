package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type UnifiedRoleEligibilitySchedule struct {
    UnifiedRoleScheduleBase
}
// NewUnifiedRoleEligibilitySchedule instantiates a new UnifiedRoleEligibilitySchedule and sets the default values.
func NewUnifiedRoleEligibilitySchedule()(*UnifiedRoleEligibilitySchedule) {
    m := &UnifiedRoleEligibilitySchedule{
        UnifiedRoleScheduleBase: *NewUnifiedRoleScheduleBase(),
    }
    return m
}
// CreateUnifiedRoleEligibilityScheduleFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateUnifiedRoleEligibilityScheduleFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUnifiedRoleEligibilitySchedule(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *UnifiedRoleEligibilitySchedule) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.UnifiedRoleScheduleBase.GetFieldDeserializers()
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
    res["scheduleInfo"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateRequestScheduleFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetScheduleInfo(val.(RequestScheduleable))
        }
        return nil
    }
    return res
}
// GetMemberType gets the memberType property value. How the role eligibility is inherited. It can either be Inherited, Direct, or Group. It can further imply whether the unifiedRoleEligibilitySchedule can be managed by the caller. Supports $filter (eq, ne).
// returns a *string when successful
func (m *UnifiedRoleEligibilitySchedule) GetMemberType()(*string) {
    val, err := m.GetBackingStore().Get("memberType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetScheduleInfo gets the scheduleInfo property value. The period of the role eligibility.
// returns a RequestScheduleable when successful
func (m *UnifiedRoleEligibilitySchedule) GetScheduleInfo()(RequestScheduleable) {
    val, err := m.GetBackingStore().Get("scheduleInfo")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(RequestScheduleable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *UnifiedRoleEligibilitySchedule) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.UnifiedRoleScheduleBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("memberType", m.GetMemberType())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("scheduleInfo", m.GetScheduleInfo())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetMemberType sets the memberType property value. How the role eligibility is inherited. It can either be Inherited, Direct, or Group. It can further imply whether the unifiedRoleEligibilitySchedule can be managed by the caller. Supports $filter (eq, ne).
func (m *UnifiedRoleEligibilitySchedule) SetMemberType(value *string)() {
    err := m.GetBackingStore().Set("memberType", value)
    if err != nil {
        panic(err)
    }
}
// SetScheduleInfo sets the scheduleInfo property value. The period of the role eligibility.
func (m *UnifiedRoleEligibilitySchedule) SetScheduleInfo(value RequestScheduleable)() {
    err := m.GetBackingStore().Set("scheduleInfo", value)
    if err != nil {
        panic(err)
    }
}
type UnifiedRoleEligibilityScheduleable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    UnifiedRoleScheduleBaseable
    GetMemberType()(*string)
    GetScheduleInfo()(RequestScheduleable)
    SetMemberType(value *string)()
    SetScheduleInfo(value RequestScheduleable)()
}

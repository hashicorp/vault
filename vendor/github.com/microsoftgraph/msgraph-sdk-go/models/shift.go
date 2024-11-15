package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Shift struct {
    ChangeTrackedEntity
}
// NewShift instantiates a new Shift and sets the default values.
func NewShift()(*Shift) {
    m := &Shift{
        ChangeTrackedEntity: *NewChangeTrackedEntity(),
    }
    odataTypeValue := "#microsoft.graph.shift"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateShiftFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateShiftFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewShift(), nil
}
// GetDraftShift gets the draftShift property value. Draft changes in the shift. Draft changes are only visible to managers. The changes are visible to employees when they are shared, which copies the changes from the draftShift to the sharedShift property.
// returns a ShiftItemable when successful
func (m *Shift) GetDraftShift()(ShiftItemable) {
    val, err := m.GetBackingStore().Get("draftShift")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ShiftItemable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Shift) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ChangeTrackedEntity.GetFieldDeserializers()
    res["draftShift"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateShiftItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDraftShift(val.(ShiftItemable))
        }
        return nil
    }
    res["schedulingGroupId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSchedulingGroupId(val)
        }
        return nil
    }
    res["sharedShift"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateShiftItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharedShift(val.(ShiftItemable))
        }
        return nil
    }
    res["userId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetUserId(val)
        }
        return nil
    }
    return res
}
// GetSchedulingGroupId gets the schedulingGroupId property value. ID of the scheduling group the shift is part of. Required.
// returns a *string when successful
func (m *Shift) GetSchedulingGroupId()(*string) {
    val, err := m.GetBackingStore().Get("schedulingGroupId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetSharedShift gets the sharedShift property value. The shared version of this shift that is viewable by both employees and managers. Updates to the sharedShift property send notifications to users in the Teams client.
// returns a ShiftItemable when successful
func (m *Shift) GetSharedShift()(ShiftItemable) {
    val, err := m.GetBackingStore().Get("sharedShift")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ShiftItemable)
    }
    return nil
}
// GetUserId gets the userId property value. ID of the user assigned to the shift. Required.
// returns a *string when successful
func (m *Shift) GetUserId()(*string) {
    val, err := m.GetBackingStore().Get("userId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Shift) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ChangeTrackedEntity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("draftShift", m.GetDraftShift())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("schedulingGroupId", m.GetSchedulingGroupId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sharedShift", m.GetSharedShift())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("userId", m.GetUserId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDraftShift sets the draftShift property value. Draft changes in the shift. Draft changes are only visible to managers. The changes are visible to employees when they are shared, which copies the changes from the draftShift to the sharedShift property.
func (m *Shift) SetDraftShift(value ShiftItemable)() {
    err := m.GetBackingStore().Set("draftShift", value)
    if err != nil {
        panic(err)
    }
}
// SetSchedulingGroupId sets the schedulingGroupId property value. ID of the scheduling group the shift is part of. Required.
func (m *Shift) SetSchedulingGroupId(value *string)() {
    err := m.GetBackingStore().Set("schedulingGroupId", value)
    if err != nil {
        panic(err)
    }
}
// SetSharedShift sets the sharedShift property value. The shared version of this shift that is viewable by both employees and managers. Updates to the sharedShift property send notifications to users in the Teams client.
func (m *Shift) SetSharedShift(value ShiftItemable)() {
    err := m.GetBackingStore().Set("sharedShift", value)
    if err != nil {
        panic(err)
    }
}
// SetUserId sets the userId property value. ID of the user assigned to the shift. Required.
func (m *Shift) SetUserId(value *string)() {
    err := m.GetBackingStore().Set("userId", value)
    if err != nil {
        panic(err)
    }
}
type Shiftable interface {
    ChangeTrackedEntityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDraftShift()(ShiftItemable)
    GetSchedulingGroupId()(*string)
    GetSharedShift()(ShiftItemable)
    GetUserId()(*string)
    SetDraftShift(value ShiftItemable)()
    SetSchedulingGroupId(value *string)()
    SetSharedShift(value ShiftItemable)()
    SetUserId(value *string)()
}

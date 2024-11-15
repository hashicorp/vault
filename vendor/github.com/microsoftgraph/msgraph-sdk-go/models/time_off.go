package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TimeOff struct {
    ChangeTrackedEntity
}
// NewTimeOff instantiates a new TimeOff and sets the default values.
func NewTimeOff()(*TimeOff) {
    m := &TimeOff{
        ChangeTrackedEntity: *NewChangeTrackedEntity(),
    }
    odataTypeValue := "#microsoft.graph.timeOff"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTimeOffFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTimeOffFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTimeOff(), nil
}
// GetDraftTimeOff gets the draftTimeOff property value. The draft version of this timeOff item that is viewable by managers. It must be shared before it is visible to team members. Required.
// returns a TimeOffItemable when successful
func (m *TimeOff) GetDraftTimeOff()(TimeOffItemable) {
    val, err := m.GetBackingStore().Get("draftTimeOff")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TimeOffItemable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TimeOff) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ChangeTrackedEntity.GetFieldDeserializers()
    res["draftTimeOff"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTimeOffItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDraftTimeOff(val.(TimeOffItemable))
        }
        return nil
    }
    res["sharedTimeOff"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTimeOffItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSharedTimeOff(val.(TimeOffItemable))
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
// GetSharedTimeOff gets the sharedTimeOff property value. The shared version of this timeOff that is viewable by both employees and managers. Updates to the sharedTimeOff property send notifications to users in the Teams client. Required.
// returns a TimeOffItemable when successful
func (m *TimeOff) GetSharedTimeOff()(TimeOffItemable) {
    val, err := m.GetBackingStore().Get("sharedTimeOff")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TimeOffItemable)
    }
    return nil
}
// GetUserId gets the userId property value. ID of the user assigned to the timeOff. Required.
// returns a *string when successful
func (m *TimeOff) GetUserId()(*string) {
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
func (m *TimeOff) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ChangeTrackedEntity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("draftTimeOff", m.GetDraftTimeOff())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("sharedTimeOff", m.GetSharedTimeOff())
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
// SetDraftTimeOff sets the draftTimeOff property value. The draft version of this timeOff item that is viewable by managers. It must be shared before it is visible to team members. Required.
func (m *TimeOff) SetDraftTimeOff(value TimeOffItemable)() {
    err := m.GetBackingStore().Set("draftTimeOff", value)
    if err != nil {
        panic(err)
    }
}
// SetSharedTimeOff sets the sharedTimeOff property value. The shared version of this timeOff that is viewable by both employees and managers. Updates to the sharedTimeOff property send notifications to users in the Teams client. Required.
func (m *TimeOff) SetSharedTimeOff(value TimeOffItemable)() {
    err := m.GetBackingStore().Set("sharedTimeOff", value)
    if err != nil {
        panic(err)
    }
}
// SetUserId sets the userId property value. ID of the user assigned to the timeOff. Required.
func (m *TimeOff) SetUserId(value *string)() {
    err := m.GetBackingStore().Set("userId", value)
    if err != nil {
        panic(err)
    }
}
type TimeOffable interface {
    ChangeTrackedEntityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDraftTimeOff()(TimeOffItemable)
    GetSharedTimeOff()(TimeOffItemable)
    GetUserId()(*string)
    SetDraftTimeOff(value TimeOffItemable)()
    SetSharedTimeOff(value TimeOffItemable)()
    SetUserId(value *string)()
}

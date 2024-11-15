package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type TimeOffReason struct {
    ChangeTrackedEntity
}
// NewTimeOffReason instantiates a new TimeOffReason and sets the default values.
func NewTimeOffReason()(*TimeOffReason) {
    m := &TimeOffReason{
        ChangeTrackedEntity: *NewChangeTrackedEntity(),
    }
    odataTypeValue := "#microsoft.graph.timeOffReason"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateTimeOffReasonFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateTimeOffReasonFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewTimeOffReason(), nil
}
// GetDisplayName gets the displayName property value. The name of the timeOffReason. Required.
// returns a *string when successful
func (m *TimeOffReason) GetDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("displayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *TimeOffReason) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.ChangeTrackedEntity.GetFieldDeserializers()
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
    res["iconType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTimeOffReasonIconType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIconType(val.(*TimeOffReasonIconType))
        }
        return nil
    }
    res["isActive"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsActive(val)
        }
        return nil
    }
    return res
}
// GetIconType gets the iconType property value. Supported icon types are: none, car, calendar, running, plane, firstAid, doctor, notWorking, clock, juryDuty, globe, cup, phone, weather, umbrella, piggyBank, dog, cake, trafficCone, pin, sunny. Required.
// returns a *TimeOffReasonIconType when successful
func (m *TimeOffReason) GetIconType()(*TimeOffReasonIconType) {
    val, err := m.GetBackingStore().Get("iconType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TimeOffReasonIconType)
    }
    return nil
}
// GetIsActive gets the isActive property value. Indicates whether the timeOffReason can be used when creating new entities or updating existing ones. Required.
// returns a *bool when successful
func (m *TimeOffReason) GetIsActive()(*bool) {
    val, err := m.GetBackingStore().Get("isActive")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *TimeOffReason) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.ChangeTrackedEntity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("displayName", m.GetDisplayName())
        if err != nil {
            return err
        }
    }
    if m.GetIconType() != nil {
        cast := (*m.GetIconType()).String()
        err = writer.WriteStringValue("iconType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isActive", m.GetIsActive())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetDisplayName sets the displayName property value. The name of the timeOffReason. Required.
func (m *TimeOffReason) SetDisplayName(value *string)() {
    err := m.GetBackingStore().Set("displayName", value)
    if err != nil {
        panic(err)
    }
}
// SetIconType sets the iconType property value. Supported icon types are: none, car, calendar, running, plane, firstAid, doctor, notWorking, clock, juryDuty, globe, cup, phone, weather, umbrella, piggyBank, dog, cake, trafficCone, pin, sunny. Required.
func (m *TimeOffReason) SetIconType(value *TimeOffReasonIconType)() {
    err := m.GetBackingStore().Set("iconType", value)
    if err != nil {
        panic(err)
    }
}
// SetIsActive sets the isActive property value. Indicates whether the timeOffReason can be used when creating new entities or updating existing ones. Required.
func (m *TimeOffReason) SetIsActive(value *bool)() {
    err := m.GetBackingStore().Set("isActive", value)
    if err != nil {
        panic(err)
    }
}
type TimeOffReasonable interface {
    ChangeTrackedEntityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetDisplayName()(*string)
    GetIconType()(*TimeOffReasonIconType)
    GetIsActive()(*bool)
    SetDisplayName(value *string)()
    SetIconType(value *TimeOffReasonIconType)()
    SetIsActive(value *bool)()
}

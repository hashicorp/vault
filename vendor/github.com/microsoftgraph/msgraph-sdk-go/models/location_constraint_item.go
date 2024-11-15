package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type LocationConstraintItem struct {
    Location
}
// NewLocationConstraintItem instantiates a new LocationConstraintItem and sets the default values.
func NewLocationConstraintItem()(*LocationConstraintItem) {
    m := &LocationConstraintItem{
        Location: *NewLocation(),
    }
    odataTypeValue := "#microsoft.graph.locationConstraintItem"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateLocationConstraintItemFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateLocationConstraintItemFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewLocationConstraintItem(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *LocationConstraintItem) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Location.GetFieldDeserializers()
    res["resolveAvailability"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResolveAvailability(val)
        }
        return nil
    }
    return res
}
// GetResolveAvailability gets the resolveAvailability property value. If set to true and the specified resource is busy, findMeetingTimes looks for another resource that is free. If set to false and the specified resource is busy, findMeetingTimes returns the resource best ranked in the user's cache without checking if it's free. Default is true.
// returns a *bool when successful
func (m *LocationConstraintItem) GetResolveAvailability()(*bool) {
    val, err := m.GetBackingStore().Get("resolveAvailability")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *LocationConstraintItem) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Location.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("resolveAvailability", m.GetResolveAvailability())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetResolveAvailability sets the resolveAvailability property value. If set to true and the specified resource is busy, findMeetingTimes looks for another resource that is free. If set to false and the specified resource is busy, findMeetingTimes returns the resource best ranked in the user's cache without checking if it's free. Default is true.
func (m *LocationConstraintItem) SetResolveAvailability(value *bool)() {
    err := m.GetBackingStore().Set("resolveAvailability", value)
    if err != nil {
        panic(err)
    }
}
type LocationConstraintItemable interface {
    Locationable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetResolveAvailability()(*bool)
    SetResolveAvailability(value *bool)()
}

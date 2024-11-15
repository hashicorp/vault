package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Presence struct {
    Entity
}
// NewPresence instantiates a new Presence and sets the default values.
func NewPresence()(*Presence) {
    m := &Presence{
        Entity: *NewEntity(),
    }
    return m
}
// CreatePresenceFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreatePresenceFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewPresence(), nil
}
// GetActivity gets the activity property value. The supplemental information to a user's availability. Possible values are Available, Away, BeRightBack, Busy, DoNotDisturb, InACall, InAConferenceCall, Inactive, InAMeeting, Offline, OffWork, OutOfOffice, PresenceUnknown, Presenting, UrgentInterruptionsOnly.
// returns a *string when successful
func (m *Presence) GetActivity()(*string) {
    val, err := m.GetBackingStore().Get("activity")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAvailability gets the availability property value. The base presence information for a user. Possible values are Available, AvailableIdle,  Away, BeRightBack, Busy, BusyIdle, DoNotDisturb, Offline, PresenceUnknown
// returns a *string when successful
func (m *Presence) GetAvailability()(*string) {
    val, err := m.GetBackingStore().Get("availability")
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
func (m *Presence) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["activity"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetActivity(val)
        }
        return nil
    }
    res["availability"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAvailability(val)
        }
        return nil
    }
    res["statusMessage"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePresenceStatusMessageFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatusMessage(val.(PresenceStatusMessageable))
        }
        return nil
    }
    return res
}
// GetStatusMessage gets the statusMessage property value. The presence status message of a user.
// returns a PresenceStatusMessageable when successful
func (m *Presence) GetStatusMessage()(PresenceStatusMessageable) {
    val, err := m.GetBackingStore().Get("statusMessage")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PresenceStatusMessageable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Presence) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("activity", m.GetActivity())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("availability", m.GetAvailability())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("statusMessage", m.GetStatusMessage())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetActivity sets the activity property value. The supplemental information to a user's availability. Possible values are Available, Away, BeRightBack, Busy, DoNotDisturb, InACall, InAConferenceCall, Inactive, InAMeeting, Offline, OffWork, OutOfOffice, PresenceUnknown, Presenting, UrgentInterruptionsOnly.
func (m *Presence) SetActivity(value *string)() {
    err := m.GetBackingStore().Set("activity", value)
    if err != nil {
        panic(err)
    }
}
// SetAvailability sets the availability property value. The base presence information for a user. Possible values are Available, AvailableIdle,  Away, BeRightBack, Busy, BusyIdle, DoNotDisturb, Offline, PresenceUnknown
func (m *Presence) SetAvailability(value *string)() {
    err := m.GetBackingStore().Set("availability", value)
    if err != nil {
        panic(err)
    }
}
// SetStatusMessage sets the statusMessage property value. The presence status message of a user.
func (m *Presence) SetStatusMessage(value PresenceStatusMessageable)() {
    err := m.GetBackingStore().Set("statusMessage", value)
    if err != nil {
        panic(err)
    }
}
type Presenceable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetActivity()(*string)
    GetAvailability()(*string)
    GetStatusMessage()(PresenceStatusMessageable)
    SetActivity(value *string)()
    SetAvailability(value *string)()
    SetStatusMessage(value PresenceStatusMessageable)()
}

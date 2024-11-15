package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Attendee struct {
    AttendeeBase
}
// NewAttendee instantiates a new Attendee and sets the default values.
func NewAttendee()(*Attendee) {
    m := &Attendee{
        AttendeeBase: *NewAttendeeBase(),
    }
    odataTypeValue := "#microsoft.graph.attendee"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateAttendeeFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAttendeeFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAttendee(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Attendee) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.AttendeeBase.GetFieldDeserializers()
    res["proposedNewTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateTimeSlotFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetProposedNewTime(val.(TimeSlotable))
        }
        return nil
    }
    res["status"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateResponseStatusFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStatus(val.(ResponseStatusable))
        }
        return nil
    }
    return res
}
// GetProposedNewTime gets the proposedNewTime property value. An alternate date/time proposed by the attendee for a meeting request to start and end. If the attendee hasn't proposed another time, then this property isn't included in a response of a GET event.
// returns a TimeSlotable when successful
func (m *Attendee) GetProposedNewTime()(TimeSlotable) {
    val, err := m.GetBackingStore().Get("proposedNewTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(TimeSlotable)
    }
    return nil
}
// GetStatus gets the status property value. The attendee's response (none, accepted, declined, etc.) for the event and date-time that the response was sent.
// returns a ResponseStatusable when successful
func (m *Attendee) GetStatus()(ResponseStatusable) {
    val, err := m.GetBackingStore().Get("status")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(ResponseStatusable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Attendee) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.AttendeeBase.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("proposedNewTime", m.GetProposedNewTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("status", m.GetStatus())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetProposedNewTime sets the proposedNewTime property value. An alternate date/time proposed by the attendee for a meeting request to start and end. If the attendee hasn't proposed another time, then this property isn't included in a response of a GET event.
func (m *Attendee) SetProposedNewTime(value TimeSlotable)() {
    err := m.GetBackingStore().Set("proposedNewTime", value)
    if err != nil {
        panic(err)
    }
}
// SetStatus sets the status property value. The attendee's response (none, accepted, declined, etc.) for the event and date-time that the response was sent.
func (m *Attendee) SetStatus(value ResponseStatusable)() {
    err := m.GetBackingStore().Set("status", value)
    if err != nil {
        panic(err)
    }
}
type Attendeeable interface {
    AttendeeBaseable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetProposedNewTime()(TimeSlotable)
    GetStatus()(ResponseStatusable)
    SetProposedNewTime(value TimeSlotable)()
    SetStatus(value ResponseStatusable)()
}

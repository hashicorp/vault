package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EventMessageRequest struct {
    EventMessage
}
// NewEventMessageRequest instantiates a new EventMessageRequest and sets the default values.
func NewEventMessageRequest()(*EventMessageRequest) {
    m := &EventMessageRequest{
        EventMessage: *NewEventMessage(),
    }
    odataTypeValue := "#microsoft.graph.eventMessageRequest"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEventMessageRequestFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEventMessageRequestFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewEventMessageRequest(), nil
}
// GetAllowNewTimeProposals gets the allowNewTimeProposals property value. True if the meeting organizer allows invitees to propose a new time when responding, false otherwise. Optional. Default is true.
// returns a *bool when successful
func (m *EventMessageRequest) GetAllowNewTimeProposals()(*bool) {
    val, err := m.GetBackingStore().Get("allowNewTimeProposals")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EventMessageRequest) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessage.GetFieldDeserializers()
    res["allowNewTimeProposals"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAllowNewTimeProposals(val)
        }
        return nil
    }
    res["meetingRequestType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMeetingRequestType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMeetingRequestType(val.(*MeetingRequestType))
        }
        return nil
    }
    res["previousEndDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreviousEndDateTime(val.(DateTimeTimeZoneable))
        }
        return nil
    }
    res["previousLocation"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreviousLocation(val.(Locationable))
        }
        return nil
    }
    res["previousStartDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPreviousStartDateTime(val.(DateTimeTimeZoneable))
        }
        return nil
    }
    res["responseRequested"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetResponseRequested(val)
        }
        return nil
    }
    return res
}
// GetMeetingRequestType gets the meetingRequestType property value. The meetingRequestType property
// returns a *MeetingRequestType when successful
func (m *EventMessageRequest) GetMeetingRequestType()(*MeetingRequestType) {
    val, err := m.GetBackingStore().Get("meetingRequestType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MeetingRequestType)
    }
    return nil
}
// GetPreviousEndDateTime gets the previousEndDateTime property value. If the meeting update changes the meeting end time, this property specifies the previous meeting end time.
// returns a DateTimeTimeZoneable when successful
func (m *EventMessageRequest) GetPreviousEndDateTime()(DateTimeTimeZoneable) {
    val, err := m.GetBackingStore().Get("previousEndDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeTimeZoneable)
    }
    return nil
}
// GetPreviousLocation gets the previousLocation property value. If the meeting update changes the meeting location, this property specifies the previous meeting location.
// returns a Locationable when successful
func (m *EventMessageRequest) GetPreviousLocation()(Locationable) {
    val, err := m.GetBackingStore().Get("previousLocation")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Locationable)
    }
    return nil
}
// GetPreviousStartDateTime gets the previousStartDateTime property value. If the meeting update changes the meeting start time, this property specifies the previous meeting start time.
// returns a DateTimeTimeZoneable when successful
func (m *EventMessageRequest) GetPreviousStartDateTime()(DateTimeTimeZoneable) {
    val, err := m.GetBackingStore().Get("previousStartDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeTimeZoneable)
    }
    return nil
}
// GetResponseRequested gets the responseRequested property value. Set to true if the sender would like the invitee to send a response to the requested meeting.
// returns a *bool when successful
func (m *EventMessageRequest) GetResponseRequested()(*bool) {
    val, err := m.GetBackingStore().Get("responseRequested")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EventMessageRequest) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessage.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteBoolValue("allowNewTimeProposals", m.GetAllowNewTimeProposals())
        if err != nil {
            return err
        }
    }
    if m.GetMeetingRequestType() != nil {
        cast := (*m.GetMeetingRequestType()).String()
        err = writer.WriteStringValue("meetingRequestType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("previousEndDateTime", m.GetPreviousEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("previousLocation", m.GetPreviousLocation())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("previousStartDateTime", m.GetPreviousStartDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("responseRequested", m.GetResponseRequested())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAllowNewTimeProposals sets the allowNewTimeProposals property value. True if the meeting organizer allows invitees to propose a new time when responding, false otherwise. Optional. Default is true.
func (m *EventMessageRequest) SetAllowNewTimeProposals(value *bool)() {
    err := m.GetBackingStore().Set("allowNewTimeProposals", value)
    if err != nil {
        panic(err)
    }
}
// SetMeetingRequestType sets the meetingRequestType property value. The meetingRequestType property
func (m *EventMessageRequest) SetMeetingRequestType(value *MeetingRequestType)() {
    err := m.GetBackingStore().Set("meetingRequestType", value)
    if err != nil {
        panic(err)
    }
}
// SetPreviousEndDateTime sets the previousEndDateTime property value. If the meeting update changes the meeting end time, this property specifies the previous meeting end time.
func (m *EventMessageRequest) SetPreviousEndDateTime(value DateTimeTimeZoneable)() {
    err := m.GetBackingStore().Set("previousEndDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetPreviousLocation sets the previousLocation property value. If the meeting update changes the meeting location, this property specifies the previous meeting location.
func (m *EventMessageRequest) SetPreviousLocation(value Locationable)() {
    err := m.GetBackingStore().Set("previousLocation", value)
    if err != nil {
        panic(err)
    }
}
// SetPreviousStartDateTime sets the previousStartDateTime property value. If the meeting update changes the meeting start time, this property specifies the previous meeting start time.
func (m *EventMessageRequest) SetPreviousStartDateTime(value DateTimeTimeZoneable)() {
    err := m.GetBackingStore().Set("previousStartDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetResponseRequested sets the responseRequested property value. Set to true if the sender would like the invitee to send a response to the requested meeting.
func (m *EventMessageRequest) SetResponseRequested(value *bool)() {
    err := m.GetBackingStore().Set("responseRequested", value)
    if err != nil {
        panic(err)
    }
}
type EventMessageRequestable interface {
    EventMessageable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAllowNewTimeProposals()(*bool)
    GetMeetingRequestType()(*MeetingRequestType)
    GetPreviousEndDateTime()(DateTimeTimeZoneable)
    GetPreviousLocation()(Locationable)
    GetPreviousStartDateTime()(DateTimeTimeZoneable)
    GetResponseRequested()(*bool)
    SetAllowNewTimeProposals(value *bool)()
    SetMeetingRequestType(value *MeetingRequestType)()
    SetPreviousEndDateTime(value DateTimeTimeZoneable)()
    SetPreviousLocation(value Locationable)()
    SetPreviousStartDateTime(value DateTimeTimeZoneable)()
    SetResponseRequested(value *bool)()
}

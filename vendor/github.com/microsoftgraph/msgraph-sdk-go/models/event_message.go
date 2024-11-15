package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type EventMessage struct {
    Message
}
// NewEventMessage instantiates a new EventMessage and sets the default values.
func NewEventMessage()(*EventMessage) {
    m := &EventMessage{
        Message: *NewMessage(),
    }
    odataTypeValue := "#microsoft.graph.eventMessage"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateEventMessageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateEventMessageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    if parseNode != nil {
        mappingValueNode, err := parseNode.GetChildNode("@odata.type")
        if err != nil {
            return nil, err
        }
        if mappingValueNode != nil {
            mappingValue, err := mappingValueNode.GetStringValue()
            if err != nil {
                return nil, err
            }
            if mappingValue != nil {
                switch *mappingValue {
                    case "#microsoft.graph.eventMessageRequest":
                        return NewEventMessageRequest(), nil
                    case "#microsoft.graph.eventMessageResponse":
                        return NewEventMessageResponse(), nil
                }
            }
        }
    }
    return NewEventMessage(), nil
}
// GetEndDateTime gets the endDateTime property value. The endDateTime property
// returns a DateTimeTimeZoneable when successful
func (m *EventMessage) GetEndDateTime()(DateTimeTimeZoneable) {
    val, err := m.GetBackingStore().Get("endDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeTimeZoneable)
    }
    return nil
}
// GetEvent gets the event property value. The event associated with the event message. The assumption for attendees or room resources is that the Calendar Attendant is set to automatically update the calendar with an event when meeting request event messages arrive. Navigation property.  Read-only.
// returns a Eventable when successful
func (m *EventMessage) GetEvent()(Eventable) {
    val, err := m.GetBackingStore().Get("event")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Eventable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *EventMessage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Message.GetFieldDeserializers()
    res["endDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEndDateTime(val.(DateTimeTimeZoneable))
        }
        return nil
    }
    res["event"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateEventFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetEvent(val.(Eventable))
        }
        return nil
    }
    res["isAllDay"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsAllDay(val)
        }
        return nil
    }
    res["isDelegated"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsDelegated(val)
        }
        return nil
    }
    res["isOutOfDate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsOutOfDate(val)
        }
        return nil
    }
    res["location"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateLocationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetLocation(val.(Locationable))
        }
        return nil
    }
    res["meetingMessageType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseMeetingMessageType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMeetingMessageType(val.(*MeetingMessageType))
        }
        return nil
    }
    res["recurrence"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreatePatternedRecurrenceFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetRecurrence(val.(PatternedRecurrenceable))
        }
        return nil
    }
    res["startDateTime"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateDateTimeTimeZoneFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetStartDateTime(val.(DateTimeTimeZoneable))
        }
        return nil
    }
    res["type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseEventType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTypeEscaped(val.(*EventType))
        }
        return nil
    }
    return res
}
// GetIsAllDay gets the isAllDay property value. The isAllDay property
// returns a *bool when successful
func (m *EventMessage) GetIsAllDay()(*bool) {
    val, err := m.GetBackingStore().Get("isAllDay")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsDelegated gets the isDelegated property value. True if this meeting request is accessible to a delegate, false otherwise. Default is false.
// returns a *bool when successful
func (m *EventMessage) GetIsDelegated()(*bool) {
    val, err := m.GetBackingStore().Get("isDelegated")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsOutOfDate gets the isOutOfDate property value. The isOutOfDate property
// returns a *bool when successful
func (m *EventMessage) GetIsOutOfDate()(*bool) {
    val, err := m.GetBackingStore().Get("isOutOfDate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetLocation gets the location property value. The location property
// returns a Locationable when successful
func (m *EventMessage) GetLocation()(Locationable) {
    val, err := m.GetBackingStore().Get("location")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Locationable)
    }
    return nil
}
// GetMeetingMessageType gets the meetingMessageType property value. The type of event message: none, meetingRequest, meetingCancelled, meetingAccepted, meetingTenativelyAccepted, meetingDeclined.
// returns a *MeetingMessageType when successful
func (m *EventMessage) GetMeetingMessageType()(*MeetingMessageType) {
    val, err := m.GetBackingStore().Get("meetingMessageType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*MeetingMessageType)
    }
    return nil
}
// GetRecurrence gets the recurrence property value. The recurrence property
// returns a PatternedRecurrenceable when successful
func (m *EventMessage) GetRecurrence()(PatternedRecurrenceable) {
    val, err := m.GetBackingStore().Get("recurrence")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(PatternedRecurrenceable)
    }
    return nil
}
// GetStartDateTime gets the startDateTime property value. The startDateTime property
// returns a DateTimeTimeZoneable when successful
func (m *EventMessage) GetStartDateTime()(DateTimeTimeZoneable) {
    val, err := m.GetBackingStore().Get("startDateTime")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(DateTimeTimeZoneable)
    }
    return nil
}
// GetTypeEscaped gets the type property value. The type property
// returns a *EventType when successful
func (m *EventMessage) GetTypeEscaped()(*EventType) {
    val, err := m.GetBackingStore().Get("typeEscaped")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*EventType)
    }
    return nil
}
// Serialize serializes information the current object
func (m *EventMessage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Message.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("endDateTime", m.GetEndDateTime())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("event", m.GetEvent())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isAllDay", m.GetIsAllDay())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isDelegated", m.GetIsDelegated())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("isOutOfDate", m.GetIsOutOfDate())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("location", m.GetLocation())
        if err != nil {
            return err
        }
    }
    if m.GetMeetingMessageType() != nil {
        cast := (*m.GetMeetingMessageType()).String()
        err = writer.WriteStringValue("meetingMessageType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("recurrence", m.GetRecurrence())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("startDateTime", m.GetStartDateTime())
        if err != nil {
            return err
        }
    }
    if m.GetTypeEscaped() != nil {
        cast := (*m.GetTypeEscaped()).String()
        err = writer.WriteStringValue("type", &cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetEndDateTime sets the endDateTime property value. The endDateTime property
func (m *EventMessage) SetEndDateTime(value DateTimeTimeZoneable)() {
    err := m.GetBackingStore().Set("endDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetEvent sets the event property value. The event associated with the event message. The assumption for attendees or room resources is that the Calendar Attendant is set to automatically update the calendar with an event when meeting request event messages arrive. Navigation property.  Read-only.
func (m *EventMessage) SetEvent(value Eventable)() {
    err := m.GetBackingStore().Set("event", value)
    if err != nil {
        panic(err)
    }
}
// SetIsAllDay sets the isAllDay property value. The isAllDay property
func (m *EventMessage) SetIsAllDay(value *bool)() {
    err := m.GetBackingStore().Set("isAllDay", value)
    if err != nil {
        panic(err)
    }
}
// SetIsDelegated sets the isDelegated property value. True if this meeting request is accessible to a delegate, false otherwise. Default is false.
func (m *EventMessage) SetIsDelegated(value *bool)() {
    err := m.GetBackingStore().Set("isDelegated", value)
    if err != nil {
        panic(err)
    }
}
// SetIsOutOfDate sets the isOutOfDate property value. The isOutOfDate property
func (m *EventMessage) SetIsOutOfDate(value *bool)() {
    err := m.GetBackingStore().Set("isOutOfDate", value)
    if err != nil {
        panic(err)
    }
}
// SetLocation sets the location property value. The location property
func (m *EventMessage) SetLocation(value Locationable)() {
    err := m.GetBackingStore().Set("location", value)
    if err != nil {
        panic(err)
    }
}
// SetMeetingMessageType sets the meetingMessageType property value. The type of event message: none, meetingRequest, meetingCancelled, meetingAccepted, meetingTenativelyAccepted, meetingDeclined.
func (m *EventMessage) SetMeetingMessageType(value *MeetingMessageType)() {
    err := m.GetBackingStore().Set("meetingMessageType", value)
    if err != nil {
        panic(err)
    }
}
// SetRecurrence sets the recurrence property value. The recurrence property
func (m *EventMessage) SetRecurrence(value PatternedRecurrenceable)() {
    err := m.GetBackingStore().Set("recurrence", value)
    if err != nil {
        panic(err)
    }
}
// SetStartDateTime sets the startDateTime property value. The startDateTime property
func (m *EventMessage) SetStartDateTime(value DateTimeTimeZoneable)() {
    err := m.GetBackingStore().Set("startDateTime", value)
    if err != nil {
        panic(err)
    }
}
// SetTypeEscaped sets the type property value. The type property
func (m *EventMessage) SetTypeEscaped(value *EventType)() {
    err := m.GetBackingStore().Set("typeEscaped", value)
    if err != nil {
        panic(err)
    }
}
type EventMessageable interface {
    Messageable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetEndDateTime()(DateTimeTimeZoneable)
    GetEvent()(Eventable)
    GetIsAllDay()(*bool)
    GetIsDelegated()(*bool)
    GetIsOutOfDate()(*bool)
    GetLocation()(Locationable)
    GetMeetingMessageType()(*MeetingMessageType)
    GetRecurrence()(PatternedRecurrenceable)
    GetStartDateTime()(DateTimeTimeZoneable)
    GetTypeEscaped()(*EventType)
    SetEndDateTime(value DateTimeTimeZoneable)()
    SetEvent(value Eventable)()
    SetIsAllDay(value *bool)()
    SetIsDelegated(value *bool)()
    SetIsOutOfDate(value *bool)()
    SetLocation(value Locationable)()
    SetMeetingMessageType(value *MeetingMessageType)()
    SetRecurrence(value PatternedRecurrenceable)()
    SetStartDateTime(value DateTimeTimeZoneable)()
    SetTypeEscaped(value *EventType)()
}

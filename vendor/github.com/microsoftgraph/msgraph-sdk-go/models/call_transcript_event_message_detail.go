package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CallTranscriptEventMessageDetail struct {
    EventMessageDetail
}
// NewCallTranscriptEventMessageDetail instantiates a new CallTranscriptEventMessageDetail and sets the default values.
func NewCallTranscriptEventMessageDetail()(*CallTranscriptEventMessageDetail) {
    m := &CallTranscriptEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.callTranscriptEventMessageDetail"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCallTranscriptEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCallTranscriptEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCallTranscriptEventMessageDetail(), nil
}
// GetCallId gets the callId property value. Unique identifier of the call.
// returns a *string when successful
func (m *CallTranscriptEventMessageDetail) GetCallId()(*string) {
    val, err := m.GetBackingStore().Get("callId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCallTranscriptICalUid gets the callTranscriptICalUid property value. Unique identifier for a call transcript.
// returns a *string when successful
func (m *CallTranscriptEventMessageDetail) GetCallTranscriptICalUid()(*string) {
    val, err := m.GetBackingStore().Get("callTranscriptICalUid")
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
func (m *CallTranscriptEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessageDetail.GetFieldDeserializers()
    res["callId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallId(val)
        }
        return nil
    }
    res["callTranscriptICalUid"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallTranscriptICalUid(val)
        }
        return nil
    }
    res["meetingOrganizer"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMeetingOrganizer(val.(IdentitySetable))
        }
        return nil
    }
    return res
}
// GetMeetingOrganizer gets the meetingOrganizer property value. The organizer of the meeting.
// returns a IdentitySetable when successful
func (m *CallTranscriptEventMessageDetail) GetMeetingOrganizer()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("meetingOrganizer")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CallTranscriptEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessageDetail.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("callId", m.GetCallId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("callTranscriptICalUid", m.GetCallTranscriptICalUid())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("meetingOrganizer", m.GetMeetingOrganizer())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCallId sets the callId property value. Unique identifier of the call.
func (m *CallTranscriptEventMessageDetail) SetCallId(value *string)() {
    err := m.GetBackingStore().Set("callId", value)
    if err != nil {
        panic(err)
    }
}
// SetCallTranscriptICalUid sets the callTranscriptICalUid property value. Unique identifier for a call transcript.
func (m *CallTranscriptEventMessageDetail) SetCallTranscriptICalUid(value *string)() {
    err := m.GetBackingStore().Set("callTranscriptICalUid", value)
    if err != nil {
        panic(err)
    }
}
// SetMeetingOrganizer sets the meetingOrganizer property value. The organizer of the meeting.
func (m *CallTranscriptEventMessageDetail) SetMeetingOrganizer(value IdentitySetable)() {
    err := m.GetBackingStore().Set("meetingOrganizer", value)
    if err != nil {
        panic(err)
    }
}
type CallTranscriptEventMessageDetailable interface {
    EventMessageDetailable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCallId()(*string)
    GetCallTranscriptICalUid()(*string)
    GetMeetingOrganizer()(IdentitySetable)
    SetCallId(value *string)()
    SetCallTranscriptICalUid(value *string)()
    SetMeetingOrganizer(value IdentitySetable)()
}

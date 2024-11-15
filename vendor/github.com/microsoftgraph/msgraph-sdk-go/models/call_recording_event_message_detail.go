package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CallRecordingEventMessageDetail struct {
    EventMessageDetail
}
// NewCallRecordingEventMessageDetail instantiates a new CallRecordingEventMessageDetail and sets the default values.
func NewCallRecordingEventMessageDetail()(*CallRecordingEventMessageDetail) {
    m := &CallRecordingEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.callRecordingEventMessageDetail"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCallRecordingEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCallRecordingEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCallRecordingEventMessageDetail(), nil
}
// GetCallId gets the callId property value. Unique identifier of the call.
// returns a *string when successful
func (m *CallRecordingEventMessageDetail) GetCallId()(*string) {
    val, err := m.GetBackingStore().Get("callId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCallRecordingDisplayName gets the callRecordingDisplayName property value. Display name for the call recording.
// returns a *string when successful
func (m *CallRecordingEventMessageDetail) GetCallRecordingDisplayName()(*string) {
    val, err := m.GetBackingStore().Get("callRecordingDisplayName")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCallRecordingDuration gets the callRecordingDuration property value. Duration of the call recording.
// returns a *ISODuration when successful
func (m *CallRecordingEventMessageDetail) GetCallRecordingDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("callRecordingDuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetCallRecordingStatus gets the callRecordingStatus property value. Status of the call recording. Possible values are: success, failure, initial, chunkFinished, unknownFutureValue.
// returns a *CallRecordingStatus when successful
func (m *CallRecordingEventMessageDetail) GetCallRecordingStatus()(*CallRecordingStatus) {
    val, err := m.GetBackingStore().Get("callRecordingStatus")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*CallRecordingStatus)
    }
    return nil
}
// GetCallRecordingUrl gets the callRecordingUrl property value. Call recording URL.
// returns a *string when successful
func (m *CallRecordingEventMessageDetail) GetCallRecordingUrl()(*string) {
    val, err := m.GetBackingStore().Get("callRecordingUrl")
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
func (m *CallRecordingEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
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
    res["callRecordingDisplayName"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallRecordingDisplayName(val)
        }
        return nil
    }
    res["callRecordingDuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallRecordingDuration(val)
        }
        return nil
    }
    res["callRecordingStatus"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseCallRecordingStatus)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallRecordingStatus(val.(*CallRecordingStatus))
        }
        return nil
    }
    res["callRecordingUrl"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallRecordingUrl(val)
        }
        return nil
    }
    res["initiator"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateIdentitySetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetInitiator(val.(IdentitySetable))
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
// GetInitiator gets the initiator property value. Initiator of the event.
// returns a IdentitySetable when successful
func (m *CallRecordingEventMessageDetail) GetInitiator()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("initiator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetMeetingOrganizer gets the meetingOrganizer property value. Organizer of the meeting.
// returns a IdentitySetable when successful
func (m *CallRecordingEventMessageDetail) GetMeetingOrganizer()(IdentitySetable) {
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
func (m *CallRecordingEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
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
        err = writer.WriteStringValue("callRecordingDisplayName", m.GetCallRecordingDisplayName())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteISODurationValue("callRecordingDuration", m.GetCallRecordingDuration())
        if err != nil {
            return err
        }
    }
    if m.GetCallRecordingStatus() != nil {
        cast := (*m.GetCallRecordingStatus()).String()
        err = writer.WriteStringValue("callRecordingStatus", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("callRecordingUrl", m.GetCallRecordingUrl())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("initiator", m.GetInitiator())
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
func (m *CallRecordingEventMessageDetail) SetCallId(value *string)() {
    err := m.GetBackingStore().Set("callId", value)
    if err != nil {
        panic(err)
    }
}
// SetCallRecordingDisplayName sets the callRecordingDisplayName property value. Display name for the call recording.
func (m *CallRecordingEventMessageDetail) SetCallRecordingDisplayName(value *string)() {
    err := m.GetBackingStore().Set("callRecordingDisplayName", value)
    if err != nil {
        panic(err)
    }
}
// SetCallRecordingDuration sets the callRecordingDuration property value. Duration of the call recording.
func (m *CallRecordingEventMessageDetail) SetCallRecordingDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("callRecordingDuration", value)
    if err != nil {
        panic(err)
    }
}
// SetCallRecordingStatus sets the callRecordingStatus property value. Status of the call recording. Possible values are: success, failure, initial, chunkFinished, unknownFutureValue.
func (m *CallRecordingEventMessageDetail) SetCallRecordingStatus(value *CallRecordingStatus)() {
    err := m.GetBackingStore().Set("callRecordingStatus", value)
    if err != nil {
        panic(err)
    }
}
// SetCallRecordingUrl sets the callRecordingUrl property value. Call recording URL.
func (m *CallRecordingEventMessageDetail) SetCallRecordingUrl(value *string)() {
    err := m.GetBackingStore().Set("callRecordingUrl", value)
    if err != nil {
        panic(err)
    }
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *CallRecordingEventMessageDetail) SetInitiator(value IdentitySetable)() {
    err := m.GetBackingStore().Set("initiator", value)
    if err != nil {
        panic(err)
    }
}
// SetMeetingOrganizer sets the meetingOrganizer property value. Organizer of the meeting.
func (m *CallRecordingEventMessageDetail) SetMeetingOrganizer(value IdentitySetable)() {
    err := m.GetBackingStore().Set("meetingOrganizer", value)
    if err != nil {
        panic(err)
    }
}
type CallRecordingEventMessageDetailable interface {
    EventMessageDetailable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCallId()(*string)
    GetCallRecordingDisplayName()(*string)
    GetCallRecordingDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetCallRecordingStatus()(*CallRecordingStatus)
    GetCallRecordingUrl()(*string)
    GetInitiator()(IdentitySetable)
    GetMeetingOrganizer()(IdentitySetable)
    SetCallId(value *string)()
    SetCallRecordingDisplayName(value *string)()
    SetCallRecordingDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetCallRecordingStatus(value *CallRecordingStatus)()
    SetCallRecordingUrl(value *string)()
    SetInitiator(value IdentitySetable)()
    SetMeetingOrganizer(value IdentitySetable)()
}

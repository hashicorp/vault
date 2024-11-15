package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type MeetingPolicyUpdatedEventMessageDetail struct {
    EventMessageDetail
}
// NewMeetingPolicyUpdatedEventMessageDetail instantiates a new MeetingPolicyUpdatedEventMessageDetail and sets the default values.
func NewMeetingPolicyUpdatedEventMessageDetail()(*MeetingPolicyUpdatedEventMessageDetail) {
    m := &MeetingPolicyUpdatedEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.meetingPolicyUpdatedEventMessageDetail"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateMeetingPolicyUpdatedEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateMeetingPolicyUpdatedEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewMeetingPolicyUpdatedEventMessageDetail(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *MeetingPolicyUpdatedEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessageDetail.GetFieldDeserializers()
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
    res["meetingChatEnabled"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMeetingChatEnabled(val)
        }
        return nil
    }
    res["meetingChatId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetMeetingChatId(val)
        }
        return nil
    }
    return res
}
// GetInitiator gets the initiator property value. Initiator of the event.
// returns a IdentitySetable when successful
func (m *MeetingPolicyUpdatedEventMessageDetail) GetInitiator()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("initiator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// GetMeetingChatEnabled gets the meetingChatEnabled property value. Represents whether the meeting chat is enabled or not.
// returns a *bool when successful
func (m *MeetingPolicyUpdatedEventMessageDetail) GetMeetingChatEnabled()(*bool) {
    val, err := m.GetBackingStore().Get("meetingChatEnabled")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetMeetingChatId gets the meetingChatId property value. Unique identifier of the meeting chat.
// returns a *string when successful
func (m *MeetingPolicyUpdatedEventMessageDetail) GetMeetingChatId()(*string) {
    val, err := m.GetBackingStore().Get("meetingChatId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *MeetingPolicyUpdatedEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessageDetail.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("initiator", m.GetInitiator())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteBoolValue("meetingChatEnabled", m.GetMeetingChatEnabled())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("meetingChatId", m.GetMeetingChatId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *MeetingPolicyUpdatedEventMessageDetail) SetInitiator(value IdentitySetable)() {
    err := m.GetBackingStore().Set("initiator", value)
    if err != nil {
        panic(err)
    }
}
// SetMeetingChatEnabled sets the meetingChatEnabled property value. Represents whether the meeting chat is enabled or not.
func (m *MeetingPolicyUpdatedEventMessageDetail) SetMeetingChatEnabled(value *bool)() {
    err := m.GetBackingStore().Set("meetingChatEnabled", value)
    if err != nil {
        panic(err)
    }
}
// SetMeetingChatId sets the meetingChatId property value. Unique identifier of the meeting chat.
func (m *MeetingPolicyUpdatedEventMessageDetail) SetMeetingChatId(value *string)() {
    err := m.GetBackingStore().Set("meetingChatId", value)
    if err != nil {
        panic(err)
    }
}
type MeetingPolicyUpdatedEventMessageDetailable interface {
    EventMessageDetailable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetInitiator()(IdentitySetable)
    GetMeetingChatEnabled()(*bool)
    GetMeetingChatId()(*string)
    SetInitiator(value IdentitySetable)()
    SetMeetingChatEnabled(value *bool)()
    SetMeetingChatId(value *string)()
}

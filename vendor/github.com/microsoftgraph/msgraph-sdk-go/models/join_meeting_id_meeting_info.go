package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type JoinMeetingIdMeetingInfo struct {
    MeetingInfo
}
// NewJoinMeetingIdMeetingInfo instantiates a new JoinMeetingIdMeetingInfo and sets the default values.
func NewJoinMeetingIdMeetingInfo()(*JoinMeetingIdMeetingInfo) {
    m := &JoinMeetingIdMeetingInfo{
        MeetingInfo: *NewMeetingInfo(),
    }
    odataTypeValue := "#microsoft.graph.joinMeetingIdMeetingInfo"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateJoinMeetingIdMeetingInfoFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateJoinMeetingIdMeetingInfoFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewJoinMeetingIdMeetingInfo(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *JoinMeetingIdMeetingInfo) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.MeetingInfo.GetFieldDeserializers()
    res["joinMeetingId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetJoinMeetingId(val)
        }
        return nil
    }
    res["passcode"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetPasscode(val)
        }
        return nil
    }
    return res
}
// GetJoinMeetingId gets the joinMeetingId property value. The ID used to join the meeting.
// returns a *string when successful
func (m *JoinMeetingIdMeetingInfo) GetJoinMeetingId()(*string) {
    val, err := m.GetBackingStore().Get("joinMeetingId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetPasscode gets the passcode property value. The passcode used to join the meeting. Optional.
// returns a *string when successful
func (m *JoinMeetingIdMeetingInfo) GetPasscode()(*string) {
    val, err := m.GetBackingStore().Get("passcode")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *JoinMeetingIdMeetingInfo) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.MeetingInfo.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteStringValue("joinMeetingId", m.GetJoinMeetingId())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("passcode", m.GetPasscode())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetJoinMeetingId sets the joinMeetingId property value. The ID used to join the meeting.
func (m *JoinMeetingIdMeetingInfo) SetJoinMeetingId(value *string)() {
    err := m.GetBackingStore().Set("joinMeetingId", value)
    if err != nil {
        panic(err)
    }
}
// SetPasscode sets the passcode property value. The passcode used to join the meeting. Optional.
func (m *JoinMeetingIdMeetingInfo) SetPasscode(value *string)() {
    err := m.GetBackingStore().Set("passcode", value)
    if err != nil {
        panic(err)
    }
}
type JoinMeetingIdMeetingInfoable interface {
    MeetingInfoable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetJoinMeetingId()(*string)
    GetPasscode()(*string)
    SetJoinMeetingId(value *string)()
    SetPasscode(value *string)()
}

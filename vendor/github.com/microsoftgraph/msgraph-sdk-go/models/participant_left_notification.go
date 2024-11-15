package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type ParticipantLeftNotification struct {
    Entity
}
// NewParticipantLeftNotification instantiates a new ParticipantLeftNotification and sets the default values.
func NewParticipantLeftNotification()(*ParticipantLeftNotification) {
    m := &ParticipantLeftNotification{
        Entity: *NewEntity(),
    }
    return m
}
// CreateParticipantLeftNotificationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateParticipantLeftNotificationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewParticipantLeftNotification(), nil
}
// GetCall gets the call property value. The call property
// returns a Callable when successful
func (m *ParticipantLeftNotification) GetCall()(Callable) {
    val, err := m.GetBackingStore().Get("call")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(Callable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *ParticipantLeftNotification) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["call"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateCallFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCall(val.(Callable))
        }
        return nil
    }
    res["participantId"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetParticipantId(val)
        }
        return nil
    }
    return res
}
// GetParticipantId gets the participantId property value. ID of the participant under the policy who has left the meeting.
// returns a *string when successful
func (m *ParticipantLeftNotification) GetParticipantId()(*string) {
    val, err := m.GetBackingStore().Get("participantId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// Serialize serializes information the current object
func (m *ParticipantLeftNotification) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("call", m.GetCall())
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("participantId", m.GetParticipantId())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetCall sets the call property value. The call property
func (m *ParticipantLeftNotification) SetCall(value Callable)() {
    err := m.GetBackingStore().Set("call", value)
    if err != nil {
        panic(err)
    }
}
// SetParticipantId sets the participantId property value. ID of the participant under the policy who has left the meeting.
func (m *ParticipantLeftNotification) SetParticipantId(value *string)() {
    err := m.GetBackingStore().Set("participantId", value)
    if err != nil {
        panic(err)
    }
}
type ParticipantLeftNotificationable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCall()(Callable)
    GetParticipantId()(*string)
    SetCall(value Callable)()
    SetParticipantId(value *string)()
}

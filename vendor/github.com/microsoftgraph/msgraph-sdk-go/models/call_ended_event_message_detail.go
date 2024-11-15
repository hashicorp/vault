package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type CallEndedEventMessageDetail struct {
    EventMessageDetail
}
// NewCallEndedEventMessageDetail instantiates a new CallEndedEventMessageDetail and sets the default values.
func NewCallEndedEventMessageDetail()(*CallEndedEventMessageDetail) {
    m := &CallEndedEventMessageDetail{
        EventMessageDetail: *NewEventMessageDetail(),
    }
    odataTypeValue := "#microsoft.graph.callEndedEventMessageDetail"
    m.SetOdataType(&odataTypeValue)
    return m
}
// CreateCallEndedEventMessageDetailFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateCallEndedEventMessageDetailFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewCallEndedEventMessageDetail(), nil
}
// GetCallDuration gets the callDuration property value. Duration of the call.
// returns a *ISODuration when successful
func (m *CallEndedEventMessageDetail) GetCallDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration) {
    val, err := m.GetBackingStore().Get("callDuration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    }
    return nil
}
// GetCallEventType gets the callEventType property value. Represents the call event type. Possible values are: call, meeting, screenShare, unknownFutureValue.
// returns a *TeamworkCallEventType when successful
func (m *CallEndedEventMessageDetail) GetCallEventType()(*TeamworkCallEventType) {
    val, err := m.GetBackingStore().Get("callEventType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*TeamworkCallEventType)
    }
    return nil
}
// GetCallId gets the callId property value. Unique identifier of the call.
// returns a *string when successful
func (m *CallEndedEventMessageDetail) GetCallId()(*string) {
    val, err := m.GetBackingStore().Get("callId")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCallParticipants gets the callParticipants property value. List of call participants.
// returns a []CallParticipantInfoable when successful
func (m *CallEndedEventMessageDetail) GetCallParticipants()([]CallParticipantInfoable) {
    val, err := m.GetBackingStore().Get("callParticipants")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]CallParticipantInfoable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *CallEndedEventMessageDetail) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.EventMessageDetail.GetFieldDeserializers()
    res["callDuration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetISODurationValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallDuration(val)
        }
        return nil
    }
    res["callEventType"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetEnumValue(ParseTeamworkCallEventType)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCallEventType(val.(*TeamworkCallEventType))
        }
        return nil
    }
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
    res["callParticipants"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateCallParticipantInfoFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]CallParticipantInfoable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(CallParticipantInfoable)
                }
            }
            m.SetCallParticipants(res)
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
    return res
}
// GetInitiator gets the initiator property value. Initiator of the event.
// returns a IdentitySetable when successful
func (m *CallEndedEventMessageDetail) GetInitiator()(IdentitySetable) {
    val, err := m.GetBackingStore().Get("initiator")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(IdentitySetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *CallEndedEventMessageDetail) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.EventMessageDetail.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteISODurationValue("callDuration", m.GetCallDuration())
        if err != nil {
            return err
        }
    }
    if m.GetCallEventType() != nil {
        cast := (*m.GetCallEventType()).String()
        err = writer.WriteStringValue("callEventType", &cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteStringValue("callId", m.GetCallId())
        if err != nil {
            return err
        }
    }
    if m.GetCallParticipants() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetCallParticipants()))
        for i, v := range m.GetCallParticipants() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("callParticipants", cast)
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
    return nil
}
// SetCallDuration sets the callDuration property value. Duration of the call.
func (m *CallEndedEventMessageDetail) SetCallDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)() {
    err := m.GetBackingStore().Set("callDuration", value)
    if err != nil {
        panic(err)
    }
}
// SetCallEventType sets the callEventType property value. Represents the call event type. Possible values are: call, meeting, screenShare, unknownFutureValue.
func (m *CallEndedEventMessageDetail) SetCallEventType(value *TeamworkCallEventType)() {
    err := m.GetBackingStore().Set("callEventType", value)
    if err != nil {
        panic(err)
    }
}
// SetCallId sets the callId property value. Unique identifier of the call.
func (m *CallEndedEventMessageDetail) SetCallId(value *string)() {
    err := m.GetBackingStore().Set("callId", value)
    if err != nil {
        panic(err)
    }
}
// SetCallParticipants sets the callParticipants property value. List of call participants.
func (m *CallEndedEventMessageDetail) SetCallParticipants(value []CallParticipantInfoable)() {
    err := m.GetBackingStore().Set("callParticipants", value)
    if err != nil {
        panic(err)
    }
}
// SetInitiator sets the initiator property value. Initiator of the event.
func (m *CallEndedEventMessageDetail) SetInitiator(value IdentitySetable)() {
    err := m.GetBackingStore().Set("initiator", value)
    if err != nil {
        panic(err)
    }
}
type CallEndedEventMessageDetailable interface {
    EventMessageDetailable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCallDuration()(*i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)
    GetCallEventType()(*TeamworkCallEventType)
    GetCallId()(*string)
    GetCallParticipants()([]CallParticipantInfoable)
    GetInitiator()(IdentitySetable)
    SetCallDuration(value *i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ISODuration)()
    SetCallEventType(value *TeamworkCallEventType)()
    SetCallId(value *string)()
    SetCallParticipants(value []CallParticipantInfoable)()
    SetInitiator(value IdentitySetable)()
}

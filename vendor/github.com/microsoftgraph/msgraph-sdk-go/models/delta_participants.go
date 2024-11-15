package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type DeltaParticipants struct {
    Entity
}
// NewDeltaParticipants instantiates a new DeltaParticipants and sets the default values.
func NewDeltaParticipants()(*DeltaParticipants) {
    m := &DeltaParticipants{
        Entity: *NewEntity(),
    }
    return m
}
// CreateDeltaParticipantsFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateDeltaParticipantsFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewDeltaParticipants(), nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *DeltaParticipants) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["participants"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateParticipantFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]Participantable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(Participantable)
                }
            }
            m.SetParticipants(res)
        }
        return nil
    }
    res["sequenceNumber"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetSequenceNumber(val)
        }
        return nil
    }
    return res
}
// GetParticipants gets the participants property value. The collection of participants that were updated since the last roster update.
// returns a []Participantable when successful
func (m *DeltaParticipants) GetParticipants()([]Participantable) {
    val, err := m.GetBackingStore().Get("participants")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]Participantable)
    }
    return nil
}
// GetSequenceNumber gets the sequenceNumber property value. The sequence number for the roster update that is used to identify the notification order.
// returns a *int64 when successful
func (m *DeltaParticipants) GetSequenceNumber()(*int64) {
    val, err := m.GetBackingStore().Get("sequenceNumber")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// Serialize serializes information the current object
func (m *DeltaParticipants) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetParticipants() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetParticipants()))
        for i, v := range m.GetParticipants() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("participants", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteInt64Value("sequenceNumber", m.GetSequenceNumber())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetParticipants sets the participants property value. The collection of participants that were updated since the last roster update.
func (m *DeltaParticipants) SetParticipants(value []Participantable)() {
    err := m.GetBackingStore().Set("participants", value)
    if err != nil {
        panic(err)
    }
}
// SetSequenceNumber sets the sequenceNumber property value. The sequence number for the roster update that is used to identify the notification order.
func (m *DeltaParticipants) SetSequenceNumber(value *int64)() {
    err := m.GetBackingStore().Set("sequenceNumber", value)
    if err != nil {
        panic(err)
    }
}
type DeltaParticipantsable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetParticipants()([]Participantable)
    GetSequenceNumber()(*int64)
    SetParticipants(value []Participantable)()
    SetSequenceNumber(value *int64)()
}
